package pay

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

var verbose bool

type MerchantInfo struct {
	AppID      string
	MerchantID string
	PaymentKey string
}

type RequestBase struct {
	XMLName    xml.Name `xml:"xml"`
	AppID      string   `xml:"appid,omitempty"`
	MerchantID string   `xml:"mch_id,omitempty"`
	Nonce      string   `xml:"nonce_str,omitempty"`
	Sign       string   `xml:"sign,omitempty"`
	SignType   `xml:"sign_type,omitempty"`
}

type ResponseBase struct {
	XMLName          xml.Name `xml:"xml"`
	ReturnCode       CData    `xml:"return_code,omitempty"`
	ReturnMessage    CData    `xml:"return_msg,omitempty"`
	AppID            CData    `xml:"appid,omitempty"`
	MerchantID       CData    `xml:"mch_id,omitempty"`
	Nonce            CData    `xml:"nonce_str,omitempty"`
	Sign             CData    `xml:"sign,omitempty"`
	ResultCode       CData    `xml:"result_code,omitempty"`
	ResultMessage    CData    `xml:"result_msg,omitempty"`
	ErrorCode        CData    `xml:"err_code,omitempty"`
	ErrorDescription CData    `xml:"err_code_des,omitempty"`
}

type Fee int32

type Time time.Time

func (t Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(time.Time(t).Format("20060102150405"), start)
	return nil
}

//go:generate stringer -type=PayStatus,RefundStatus,SignType,FeeType,TradeType,LimitPay $GOFILE
type PayStatus int

const (
	NotPay     PayStatus = 1
	PaySuccess PayStatus = 2
	PayRefund  PayStatus = 3
	PayClosed  PayStatus = 4
	PayError   PayStatus = 5
)

type RefundStatus int

const (
	RefundProcessing RefundStatus = 1
)

type SignType int

const (
	MD5  SignType = 1
	HMAC SignType = 2
)

type FeeType int

const (
	CNY FeeType = 1
)

type TradeType int

const (
	JSAPI  TradeType = 1
	NATIVE TradeType = 2
	APP    TradeType = 3
)

type LimitPay int

const (
	NoCredit LimitPay = 1
)

type ErrorCode string

const (
	SystemError ErrorCode = "SYSTEMERROR"
)

// request and response parameter of apis:

type PreOrderRequest struct {
	DeviceInfo  string     `xml:"device_info,omitempty"`
	Description string     `xml:"body,omitempty"`         // https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_2
	Detail      string     `xml:"detail,CDATA,omitempty"` // https://pay.weixin.qq.com/wiki/doc/api/danpin.php?chapter=9_102&index=2
	Attach      string     `xml:"attach,omitempty"`
	OrderID     string     `xml:"out_trade_no,omitempty"`     // 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	TotalFee    Fee        `xml:"total_fee,omitempty"`        // 订单总金额，单位为分
	CreateIP    string     `xml:"spbill_create_ip,omitempty"` // APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP
	TimeStart   Time       `xml:"time_start,omitempty"`       // yyyyMMddHHmmss
	TimeExpire  Time       `xml:"time_expire,omitempty"`      // yyyyMMddHHmmss 最短失效时间间隔必须大于5分钟
	GoodsTag    string     `xml:"goods_tag,omitempty"`        // https://pay.weixin.qq.com/wiki/doc/api/tools/sp_coupon.php?chapter=12_1
	NotifyURL   string     `xml:"notify_url,omitempty"`
	ProductID   string     `xml:"product_id,omitempty"` // 扫码支付必填
	OpenID      string     `xml:"openid,omitempty"`     // 公众号支付必填
	SceneInfo   *SceneInfo `xml:"scene_info,omitempty"` // 上报场景信息，目前支持上报实际门店信息。该字段为JSON对象数据
	TradeType   `xml:"trade_type,omitempty"`
	LimitPay    `xml:"limit_pay,omitempty"`
	FeeType     `xml:"fee_type,omitempty"` // 符合ISO 4217标准的三位字母代码，默认人民币：CNY
}

type SceneInfo struct {
	ID       string `json:"id,omitempty"`        // 门店唯一标识
	Name     string `json:"name,omitempty"`      // 门店名称
	AreaCode string `json:"area_code,omitempty"` // 门店所在地行政区划码
	Address  string `json:"address,omitempty"`   // 门店详细地址
}

func (si *SceneInfo) String() string {
	body, e := json.Marshal(si)
	if e != nil {
		return ""
	}
	return string(body)
}

func (si *SceneInfo) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	body, err := json.Marshal(si)
	if err != nil {
		return err
	}
	e.EncodeElement(body, start)
	return nil
}

type PreOrderResponse struct {
	DeviceInfo CData `xml:"deviceInfo,omitempty"`
	TradeType  `xml:"tradeType,omitempty"`
	PrepayID   CData `xml:"prepayID,omitempty"`
	CodeURL    CData `xml:"codeURL,omitempty"`
}

type QueryOrderRequest struct {
	TransactionID string `xml:"transactionID,omitempty"` // 商户订单号 二选一
	OrderID       string `xml:"orderID,omitempty"`       //微信的订单号，建议优先使用
}

type QueryOrderResponse struct {
	DeviceInfo string
}

// helpers for xml marshalling:

type CData struct {
	Data string `xml:",cdata"`
}

func (d CData) String() string {
	return d.Data
}

var signTypeNames = []string{"MD5", "HMAC-SHA256"}

func (st SignType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if int(st) > len(signTypeNames) {
		e.EncodeElement(st.String(), start)
	} else {
		e.EncodeElement(signTypeNames[st-1], start)
	}
	return nil
}

func (tt TradeType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(tt.String(), start)
	return nil
}

var limitPayNames = []string{"no_credit"}

func (lp LimitPay) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if int(lp) > len(limitPayNames) {
		e.EncodeElement(lp.String(), start)
	} else {
		e.EncodeElement(limitPayNames[lp-1], start)
	}
	return nil
}
