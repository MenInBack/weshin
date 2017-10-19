package pay

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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
func (m *MerchantInfo) PrepareAppPay(req *PreOrderRequest) (*AppPayRequest, error) {
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

func (m *MerchantInfo) RefundOrder(req *RefundRequest) (*RefundResponse, error) {
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

func (m *MerchantInfo) QueryRefund(req *QueryRefundRequest) (*QueryRefundResponse, error) {
	// check parameters
	if req.TradeNo == "" && req.TransactionID == "" && req.RefundID == "" && req.RefundNo == "" {
		return nil, wx.ParameterError{"none of tradeNo, transactionID, refundID nor refundNo provided"}
	}

	resp := new(QueryRefundResponse)
	if e := m.postXML(urlQueryRefund, req, resp, false); e != nil {
		return nil, e
	}

	return resp, nil
}

// https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_6&index=8
const urlDownloadBill = "https://api.mch.weixin.qq.com/pay/downloadbill"

func (m *MerchantInfo) DownloadBill(req *DownloadBillRequest) ([]*Bill, *BillInTotal, error) {
	if time.Since(time.Time(req.BillData)) > time.Hour*24*92 {
		return nil, nil, wx.ParameterError{"billDate"}
	}
	body, e := m.prepareRequest(req)
	if e != nil {
		return nil, nil, e
	}
	if verbose {
		log.Println("request downloading bill: ", string(body))
	}

	resp, e := http.Post(urlDownloadBill, "application/xml", bytes.NewBuffer(body))
	if e != nil {
		return nil, nil, e
	}
	defer resp.Body.Close()
	if verbose {
		log.Println("response: ", resp)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, wx.HttpError{
			State: resp.StatusCode,
		}
	}

	return parseBillFile(resp.Body, req.BillType)
}

func parseBillFile(r io.Reader, t BillType) ([]*Bill, *BillInTotal, error) {
	bills := make([]*Bill, 0)

	s := bufio.NewScanner(r)
	s.Scan() // drop title
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 || line[0] != '`' {
			break // second title
		}
		fs := strings.Split(line[1:], ",`") // trim leading "`"
		tails := fs[len(fs)-4 : len(fs)]    // trailling common fields
		b := &Bill{
			Time:          parseBillTime(fs[0]),
			AppID:         fs[1],
			MerchantID:    fs[2],
			SubMerchantID: fs[3],
			DeviceInfo:    fs[4],
			TransactionID: fs[5],
			TradeNo:       fs[6],
			OpenID:        fs[7],
			TradeType:     TradeType(fs[8]),
			TradeState:    TradeState(fs[9]),
			BankType:      BankType(fs[10]),
			FeeType:       FeeType(fs[11]),
			TotalFee:      parseFee(fs[12]),
			CouponFee:     parseFee(fs[13]),
			ProductName:   tails[0],
			Attach:        tails[1],
			ServiceCharge: parseFee(tails[2]),
			ChargeRate:    tails[3],
		}

		switch t {
		case BillAll:
			b.RefundID = fs[14]
			b.RefundNo = fs[15]
			b.RefundFee = parseFee(fs[16])
			b.CouponRefundFee = parseFee(fs[17])
			b.RefundType = fs[18]
			b.RefundStatus = RefundStatus(fs[19])
		case BillSuccess:
		case BillRefund:
			b.RefundApplyTime = parseBillTime(fs[14])
			b.RefundSucceedTime = parseBillTime(fs[15])
			b.RefundID = fs[16]
			b.RefundNo = fs[17]
			b.RefundFee = parseFee(fs[18])
			b.CouponRefundFee = parseFee(fs[19])
			b.RefundType = fs[20]
			b.RefundStatus = RefundStatus(fs[21])
		}
		bills = append(bills, b)
	}
	if e := s.Err(); e != nil {
		return nil, nil, e
	}

	s.Scan()
	statics := s.Text()
	fs := strings.Split(statics[1:], ",`")
	total := &BillInTotal{
		Transactions:    parseInt(fs[0]),
		TradeFee:        parseFee(fs[1]),
		RefundFee:       parseFee(fs[2]),
		CouponRefundFee: parseFee(fs[3]),
		Charge:          parseFee(fs[4]),
	}

	return bills, total, nil
}

func parseBillTime(s string) time.Time {
	t, e := time.Parse("2096-01-0215：04：05", s)
	if e != nil {
		return time.Time{}
	}
	return t
}

func parseFee(s string) Fee {
	if s == "0" {
		return 0
	}

	if !strings.ContainsRune(s, '.') {
		f, _ := strconv.ParseInt(s, 10, 64)
		return Fee(f * 100) // to cents
	}

	ss := strings.Split(s, ".")
	if len(ss) > 2 {
		panic("invalid fee type to parse")
	}
	f, _ := strconv.ParseInt(ss[0], 10, 64)
	f *= 100 // to cents

	if len(ss[1]) > 2 {
		panic("invalid fee type to parse")
	}

	if len(ss[1]) == 1 {
		d, _ := strconv.ParseInt(ss[1], 10, 64)
		f += d * 10
	}
	if len(ss[1]) == 2 {
		d, _ := strconv.ParseInt(ss[1], 10, 64)
		f += d
	}

	return Fee(f)
}

func parseInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

// 当日所有订单
// 交易时间,公众账号ID,商户号,子商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,总金额,代金券或立减优惠金额,微信退款单号,商户退款单号,退款金额,代金券或立减优惠退款金额，退款类型，退款状态,商品名称,商户数据包,手续费,费率

// 当日成功支付的订单
// 交易时间,公众账号ID,商户号,子商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,总金额,代金券或立减优惠金额,商品名称,商户数据包,手续费,费率

// 当日退款的订单
// 交易时间,公众账号ID,商户号,子商户号,设备号,微信订单号,商户订单号,用户标识,交易类型,交易状态,付款银行,货币种类,总金额,代金券或立减优惠金额,退款申请时间,退款成功时间,微信退款单号,商户退款单号,退款金额,代金券或立减优惠退款金额,退款类型,退款状态,商品名称,商户数据包,手续费,费率
