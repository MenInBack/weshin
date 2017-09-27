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

type CData struct {
	Data string `xml:",cdata"`
}

func (cd CData) String() string {
	return cd.Data
}

func (cd *CData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return d.DecodeElement(cd.Data, &start)
}

type Fee int32

type Time time.Time

func (t *Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(time.Time(*t).Format("20060102150405"), start)
	return nil
}

func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	s := new(string)
	e := d.DecodeElement(s, &start)
	if e != nil {
		return e
	}
	tt, e := time.Parse("20060102150405", *s)
	if e != nil {
		return e
	}
	*t = Time(tt)
	return nil
}

type TradeState string

const (
	NotPay     TradeState = "NOTPAY"
	PaySuccess TradeState = "SUCCESS"
	PayRefund  TradeState = "REFUND"
	PayClosed  TradeState = "CLOSED"
	PayError   TradeState = "PAYERROR"
	PayRevoked TradeState = "REVOKED"
	Paying     TradeState = "USERPAYING"
)

type RefundStatus string

const (
	RefundProcessing RefundStatus = "PROCESSING"
	RefundSuccess    RefundStatus = "SUCCESS"
	RefundClosed     RefundStatus = "REFUNDCLOSE"
	RefundChanged    RefundStatus = "CHANGE"
)

type SignType string

const (
	MD5  SignType = "MD5"
	HMAC SignType = "HMAC-SHA256"
)

type FeeType string

const (
	CNY FeeType = "CNY"
)

type TradeType string

const (
	JSAPI  TradeType = "JSAPI"
	NATIVE TradeType = "NATIVE"
	APP    TradeType = "APP"
)

type LimitPay string

const (
	NoCredit LimitPay = "no_credit"
)

type BankType string

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
	TradeNo     string     `xml:"out_trade_no,omitempty"`     // 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
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
	TradeType  `xml:"trade_type,omitempty"`
	PrepayID   CData `xml:"prepay_id,omitempty"`
	CodeURL    CData `xml:"code_url,omitempty"`
}

type QueryOrderRequest struct {
	TransactionID string `xml:"transaction_id,omitempty"` // 商户订单号 二选一
	TradeNo       string `xml:"out_trade_no,omitempty"`   //微信的订单号，建议优先使用
}

// not include coupon_type_$n,coupon_id_$n,coupon_fee_$n
type QueryOrderResponse struct {
	DeviceInfo     CData   `xml:"device_info,omitempty"`
	OpenID         CData   `xml:"openid,omitempty"`
	IsSubscribe    CData   `xml:"is_subscribe,omitempty"`
	TotalFee       Fee     `xml:"total_fee,omitempty"`
	SettlementFee  Fee     `xml:"settlement_total_fee,omitempty"` // 应结订单金额, 当订单使用了免充值型优惠券后返回该参数，应结订单金额=订单金额-免充值优惠券金额。
	CashFee        Fee     `xml:"cash_fee,omitempty"`
	CashFeeType    FeeType `xml:"cash_fee_type,omitempty"`
	CouponFee      Fee     `xml:"coupon_fee,omitempty"`
	CouponCount    int     `xml:"coupon_count,omitempty"`
	TransactionID  CData   `xml:"transaction_id,omitempty"`
	TradeNo        string  `xml:"out_trade_no,omitempty"`
	Attach         CData   `xml:"attach,omitempty"`
	TimeEnd        Time    `xml:"time_end,omitempty"`
	TradeStateDesc string  `xml:"trade_state_desc,omitempty"`
	TradeType      `xml:"trade_type,omitempty"`
	TradeState     `xml:"trade_state,omitempty"`
	BankType       `xml:"bank_type,omitempty"`
	FeeType        `xml:"fee_type,omitempty"`
}

type CloseOrderRequest struct {
	TradeNo string `xml:"out_trade_no,omitempty"`
}

type RefundRequest struct {
	TransactionID CData   `xml:"transaction_id,omitempty"`
	TradeNo       string  `xml:"out_trade_no,omitempty"`  // 商户系统内部订单号
	RefundNo      string  `xml:"out_refund_no,omitempty"` // 商户系统内部的退款单号
	TotalFee      Fee     `xml:"total_fee,omitempty"`
	RefundFee     Fee     `xml:"refund_fee,omitempty"`
	RefundFeeType FeeType `xml:"refund_fee_type,omitempty"`
	RefundDesc    string  `xml:"refund_desc,omitempty"`
}

// not include coupon_type_$n, coupon_refund_fee_$n, coupon_refund_id_$n
type RefundResponse struct {
	TransactionID       CData   `xml:"transaction_id,omitempty"`
	TradeNo             string  `xml:"out_trade_no,omitempty"`
	RefundNo            string  `xml:"out_refund_no,omitempty"`
	RefundID            string  `xml:"refund_id,omitempty"` // 微信退款单号
	RefundFee           Fee     `xml:"refund_fee,omitempty"`
	RefundFeeType       FeeType `xml:"refund_fee_type,omitempty"`
	SettlementRefundFee Fee     `xml:"settlement_refund_fee,omitempty"` // 应结退款金额, 去掉非充值代金券退款金额后的退款金额，退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额应结退款金额
	TotalFee            Fee     `xml:"total_fee,omitempty"`
	SettlementFee       Fee     `xml:"settlement_total_fee,omitempty"` // 应结订单金额, 去掉非充值代金券金额后的订单总金额，应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	FeeType             FeeType `xml:"fee_type,omitempty"`
	CashFee             Fee     `xml:"cash_fee,omitempty"`
	CashFeeType         FeeType `xml:"cash_fee_type,omitempty"`
	CashRefundFee       Fee     `xml:"cash_refund_fee,omitempty"`
	CouponRefundFee     Fee     `xml:"coupon_refund_fee,omitempty"`
	CouponRefundCount   int     `xml:"coupon_refund_count,omitempty"`
	RefundDesc          string  `xml:"refund_desc,omitempty"`
}

// 四选一 refund_id > out_refund_no > transaction_id > out_trade_no
type QueryRefundRequest struct {
	TransactionID CData  `xml:"transaction_id,omitempty"`
	TradeNo       string `xml:"out_trade_no,omitempty"`
	RefundNo      string `xml:"out_refund_no,omitempty"`
	RefundID      string `xml:"refund_id,omitempty"` // 微信退款单号
}

type QueryRefundResponse struct {
	TransactionID  CData   `xml:"transaction_id,omitempty"`
	TradeNo        string  `xml:"out_trade_no,omitempty"`
	TotalFee       Fee     `xml:"total_fee,omitempty"`
	SettlementFee  Fee     `xml:"settlement_total_fee,omitempty"` // 应结订单金额, 去掉非充值代金券金额后的订单总金额，应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	FeeType        FeeType `xml:"fee_type,omitempty"`
	CashFee        Fee     `xml:"cash_fee,omitempty"`
	RefundCount    int     `xml:"refund_count,omitempty"`
	RefundNos      []CData `xml:"out_refund_no,omitempty"`
	RefundIDs      []CData `xml:"refund_id,omitempty"`
	RefundChannels []CData `xml:"refund_channel,omitempty"`
}
