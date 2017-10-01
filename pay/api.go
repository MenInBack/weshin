package pay

import (
	"log"
	"strconv"
	"time"

	"github.com/MenInBack/weshin/crypto"
	"github.com/MenInBack/weshin/wx"
)

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=7_7&index=6
func (m *MerchantInfo) PrepareJSAPIPay(req *PreOrderRequest) (*JSPayRequest, error) {
	req.TradeType = JSAPIPay
	req.DeviceInfo = "WEB"

	if req.OpenID == "" {
		return nil, wx.ParameterError{"openID"}
	}

	resp, e := m.preOrder(req)
	if e != nil {
		return nil, e
	}

	pr := &JSPayRequest{
		AppID:     m.AppID,
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		Nonce:     string(crypto.RandString(NonceLength)),
		Package:   "prepay_id=" + resp.PrepayID.Data,
		SignType:  MD5,
	}

	fields := []field{
		field{"appId", m.AppID},
		field{"timeStamp", pr.TimeStamp},
		field{"nonceStr", pr.Nonce},
		field{"package", pr.Package},
	}

	s, e := sign(fields, m.PaymentKey, MD5)
	if e != nil {
		return nil, e
	}
	pr.PaySign = s
	if verbose {
		log.Println("jspai pay request: ", *pr)
	}

	return pr, nil
}

// https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_12&index=2
func PrepareAppPay(req *PreOrderRequest) (*AppPayRequest, error) {
	req.TradeType = AppPay
	req.DeviceInfo = "WEB"

	resp, e := m.preOrder(req)
	if e != nil {
		return nil, e
	}

	pr := &AppPayRequest{
		AppID:     m.AppID,
		PartnerID: m.MerchantID,
		PrePayID:  resp.PrepayID.Data,
		Package:   "Sign=WXPay",
		Nonce:     string(crypto.RandString(NonceLength)),
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
	}

	fields := []field{
		field{"appid", pr.AppID},
		field{"partnerid", pr.PartnerID},
		field{"prepayid", pr.PrePayID},
		field{"package", pr.Package},
		field{"noncestr", pr.Nonce},
		field{"timestamp", pr.TimeStamp},
	}

	s, e := sign(fields, m.PaymentKey, MD5)
	if e != nil {
		return nil, e
	}
	pr.Sign = s
	if verbose {
		log.Println("app pay request: ", *pr)
	}

	return pr, nil
}

func (m *MerchantInfo) PrepareQRPay(req *PreOrderRequest) (codeURL string, e error) {
	req.TradeType = QRPay
	req.DeviceInfo = "WEB"

	resp, e := m.preOrder(req)
	if e != nil {
		return "", e
	}
	if verbose {
		log.Println("code url for qr pay: ", resp.CodeURL)
	}

	return resp.CodeURL.Data, nil
}

// https://pay.weixin.qq.com/wiki/doc/api/H5.php?chapter=9_20&index=1
func (m *MerchantInfo) PrepareWebPay(req *PreOrderRequest) (url string, e error) {
	req.TradeType = WebPay
	req.DeviceInfo = "WEB"

	resp, e := m.preOrder(req)
	if e != nil {
		return "", e
	}
	if verbose {
		log.Println("web url for mobile web pay: ", resp.WebURL)
	}

	return resp.WebURL.Data, nil
}

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
const urlPreOrder = "https://api.mch.weixin.qq.com/pay/unifiedorder"

func (m *MerchantInfo) preOrder(req *PreOrderRequest) (*PreOrderResponse, error) {
	// check parameters
	if req.Description == "" {
		return nil, wx.ParameterError{"description"}
	}
	if req.TradeType == "" {
		return nil, wx.ParameterError{"tradeType"}
	}
	if req.TradeNo == "" {
		return nil, wx.ParameterError{"tradeNo"}
	}
	if req.TotalFee <= 0 {
		return nil, wx.ParameterError{"totalFee"}
	}
	if req.CreateIP == "" {
		return nil, wx.ParameterError{"createIP"}
	}
	if time.Time(req.TimeStart).IsZero() {
		return nil, wx.ParameterError{"timeStart"}
	}
	if time.Time(req.TimeExpire).IsZero() {
		return nil, wx.ParameterError{"timeExpire"}
	}
	if !time.Time(req.TimeStart).Before(time.Time(req.TimeExpire)) {
		return nil, wx.ParameterError{"timeExpire before timeStart"}
	}

	if req.NotifyURL == "" {
		req.NotifyURL = m.PayNotifyURL
	}
	if req.FeeType == "" {
		req.FeeType = CNY
	}

	resp := new(PreOrderResponse)
	if e := m.postXML(urlPreOrder, req, resp, false); e != nil {
		return nil, e
	}

	return resp, nil
}

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
const urlOrderQuery = "https://api.mch.weixin.qq.com/pay/orderquery"

func (m *MerchantInfo) QueryOrder(req *QueryOrderRequest) (*QueryOrderResponse, error) {
	// check parameters
	if req.TradeNo == "" && req.TransactionID == "" {
		return nil, wx.ParameterError{"no tradeNo nor transactionID"}
	}

	resp := new(QueryOrderResponse)
	if e := m.postXML(urlPreOrder, req, resp, false); e != nil {
		return nil, e
	}

	return resp, nil
}

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_3
const urlCloseOrder = "https://api.mch.weixin.qq.com/pay/closeorder"

func (m *MerchantInfo) CloseOrder(req *CloseOrderRequest) error {
	// check parameters
	if req.TradeNo == "" {
		return wx.ParameterError{"tradeNo"}
	}

	if e := m.postXML(urlPreOrder, req, nil, false); e != nil {
		return e
	}
	return nil
}

// need certification
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_4
const urlRefundOrder = "https://api.mch.weixin.qq.com/secapi/pay/refund"

func (m *MerchantInfo) RefundOrder(req RefundRequest) (*RefundResponse, error) {
	// check parameters
	if req.TransactionID == "" && req.TradeNo == "" {
		return nil, wx.ParameterError{"no tradeNo nor transactionID"}
	}
	if req.RefundNo == "" {
		return nil, wx.ParameterError{"refundNo"}
	}
	if req.TotalFee <= 0 {
		return nil, wx.ParameterError{"totalFee"}
	}
	if req.RefundFee <= 0 {
		return nil, wx.ParameterError{"refundFee"}
	}

	if req.RefundFeeType == "" {
		req.RefundFeeType = CNY
	}

	resp := new(RefundResponse)
	if e := m.postXML(urlPreOrder, req, resp, true); e != nil {
		return nil, e
	}
	return resp, nil
}

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_5
const urlQueryRefund = "https://api.mch.weixin.qq.com/pay/refundquery"

func (m *MerchantInfo) QueryRefund(req QueryRefundRequest) (*QueryRefundResponse, error) {
	// check parameters
	if req.TradeNo == "" && req.TransactionID == "" && req.RefundID == "" && req.RefundNo == "" {
		return nil, wx.ParameterError{"none of tradeNo, transactionID, refundID nor refundNo"}
	}

	resp := new(QueryRefundResponse)
	if e := m.postXML(urlQueryRefund, req, resp, false); e != nil {
		return nil, e
	}

	return resp, nil
}

const urlDownloadBill = "https://api.mch.weixin.qq.com/pay/downloadbill"

// func (m *MerchantInfo) DownloadBill(req DownloadBillRequest) error {}
