package pay

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strconv"
	"time"
)

var verbose bool
var donotCheckSign bool

type MerchantInfo struct {
	AppID           string
	MerchantID      string
	PaymentKey      string
	PayNotifyURL    string
	RefundNotifyURL string
	PayNoticeHander
	RefundNoticeHandler
}

type JSPayRequest struct {
	AppID     string `json:"appId,omitempty"`
	TimeStamp string `json:"timeStamp,omitempty"`
	Nonce     string `json:"nonceStr,omitempty"`
	Package   string `json:"package,omitempty"`
	SignType  `json:"signType,omitempty"`
	PaySign   string `json:"paySign,omitempty"`
}

type AppPayRequest struct {
	AppID     string `json:"appid,omitempty"`
	PartnerID string `json:"partnerid,omitempty"` // merchantID
	PrePayID  string `json:"prepayid,omitempty"`
	Package   string `json:"package,omitempty"`
	Nonce     string `json:"noncestr,omitempty"`
	TimeStamp string `json:"timestamp,omitempty"`
	Sign      string `json:"sign,omitempty"`
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

type Unstringer interface {
	Unstring(string) error
}

type CData struct {
	Data string `xml:",cdata"`
}

func (cd CData) String() string {
	return cd.Data
}

func (cd *CData) Unstring(s string) error {
	cd.Data = s
	return nil
}

func (cd *CData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return d.DecodeElement(cd.Data, &start)
}

type Fee int

func (f Fee) String() string {
	if f == 0 {
		return ""
	}
	return strconv.Itoa(int(f))
}

func (f *Fee) Unstring(s string) error {
	i, _ := strconv.Atoi(s)
	*f = Fee(i)
	return nil
}

func (f *Fee) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return d.DecodeElement((*int)(f), &start)
}

type Time time.Time

func (t Time) String() string {
	if time.Time(t).IsZero() {
		return ""
	}
	return time.Time(t).Format("20060102150405")
}

func (t *Time) Unstring(s string) error {
	tt, e := time.Parse("20060102150405", s)
	if e != nil {
		return e
	}
	*t = Time(tt)
	return nil
}

func (t Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if time.Time(t).IsZero() {
		return nil
	}
	e.EncodeElement(time.Time(t).Format("20060102150405"), start)
	return nil
}

func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if e := d.DecodeElement(&s, &start); e != nil {
		return e
	}
	if tt, e := time.Parse("20060102150405", s); e != nil {
		return e
	} else {
		*t = Time(tt)
	}
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

type CouponType string

const (
	CouponCash   CouponType = "CASH"
	CouponNoCash CouponType = "NO_CASH"
)

type BillType string

const (
	BillAll     BillType = "ALL"
	BillSuccess BillType = "SUCCESS"
	BillRefund  BillType = "REFUND"
)

type TarType string

const (
	GZIP TarType = "GZIP"
)

type ErrorCode string

const (
	SystemError ErrorCode = "SYSTEMERROR"
)

// request and response parameter of apis:

type PreOrderRequest struct {
	DeviceInfo  string     `xml:"device_info,omitempty"`      // 设备号
	Description string     `xml:"body,omitempty"`             // 商品描述 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_2
	Detail      string     `xml:"detail,CDATA,omitempty"`     // 商品详情https://pay.weixin.qq.com/wiki/doc/api/danpin.php?chapter=9_102&index=2
	Attach      string     `xml:"attach,omitempty"`           // 附加数据
	TradeNo     string     `xml:"out_trade_no,omitempty"`     // 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	FeeType     FeeType    `xml:"fee_type,omitempty"`         // 符合ISO 4217标准的三位字母代码，默认人民币：CNY
	TotalFee    Fee        `xml:"total_fee,omitempty"`        // 标价金额, 订单总金额，单位为分
	CreateIP    string     `xml:"spbill_create_ip,omitempty"` // APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP
	TimeStart   Time       `xml:"time_start,omitempty"`       // yyyyMMddHHmmss
	TimeExpire  Time       `xml:"time_expire,omitempty"`      // yyyyMMddHHmmss 最短失效时间间隔必须大于5分钟
	GoodsTag    string     `xml:"goods_tag,omitempty"`        // https://pay.weixin.qq.com/wiki/doc/api/tools/sp_coupon.php?chapter=12_1
	NotifyURL   string     `xml:"notify_url,omitempty"`       // 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
	TradeType   TradeType  `xml:"trade_type,omitempty"`       // 交易类型 取值如下：JSAPI，NATIVE，APP
	ProductID   string     `xml:"product_id,omitempty"`       // 扫码支付必填 此参数为二维码中包含的商品ID，商户自行定义。
	LimitPay    LimitPay   `xml:"limit_pay,omitempty"`        // 上传此参数no_credit--可限制用户不能使用信用卡支付
	OpenID      string     `xml:"openid,omitempty"`           // 公众号支付必填
	SceneInfo   *SceneInfo `xml:"scene_info,omitempty"`       // 上报场景信息，目前支持上报实际门店信息。该字段为JSON对象数据
}

type SceneInfo struct {
	ID       string `json:"id,omitempty"`        // 门店唯一标识
	Name     string `json:"name,omitempty"`      // 门店名称
	AreaCode string `json:"area_code,omitempty"` // 门店所在地行政区划码
	Address  string `json:"address,omitempty"`   // 门店详细地址
}

func (si *SceneInfo) String() string {
	if si == nil {
		return ""
	}
	if si.ID == "" && si.Name == "" && si.AreaCode == "" && si.Address == "" {
		return ""
	}
	body, e := json.Marshal(si)
	if e != nil {
		return ""
	}

	return string(body)
}

func (si *SceneInfo) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if si == nil {
		return nil
	}
	if si.ID == "" && si.Name == "" && si.AreaCode == "" && si.Address == "" {
		return nil
	}
	body, err := json.Marshal(si)
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString("[CDATA[")
	buf.Write(body)
	buf.WriteString("]]")

	e.EncodeToken(start)
	e.EncodeToken(xml.Directive(buf.Bytes()))
	e.EncodeToken(start.End())
	return nil
}

func (si *SceneInfo) Unstring(s string) error {
	return json.Unmarshal([]byte(s), si)
}

type PreOrderResponse struct {
	DeviceInfo CData `xml:"deviceInfo,omitempty"`
	TradeType  `xml:"trade_type,omitempty"`
	PrepayID   CData `xml:"prepay_id,omitempty"`
	CodeURL    CData `xml:"code_url,omitempty"`
}

type QueryOrderRequest struct {
	TransactionID string `xml:"transaction_id,omitempty"` // 商户订单号 二选一
	TradeNo       string `xml:"out_trade_no,omitempty"`   // 微信的订单号，建议优先使用
}

type QueryOrderResponse struct {
	DeviceInfo     CData      `xml:"device_info,omitempty"`  // 设备号
	OpenID         CData      `xml:"openid,omitempty"`       // 用户标识
	IsSubscribe    CData      `xml:"is_subscribe,omitempty"` // 是否关注公众账号
	TradeType      TradeType  `xml:"trade_type,omitempty"`
	TradeStatus    TradeState `xml:"trade_state,omitempty"`
	BankType       BankType   `xml:"bank_type,omitempty"`            // 付款银行
	TotalFee       Fee        `xml:"total_fee,omitempty"`            // 标价金额
	SettlementFee  Fee        `xml:"settlement_total_fee,omitempty"` // 应结订单金额, 当订单使用了免充值型优惠券后返回该参数，应结订单金额=订单金额-免充值优惠券金额。
	FeeType        FeeType    `xml:"fee_type,omitempty"`
	CashFee        Fee        `xml:"cash_fee,omitempty"` // 现金支付金额
	CashFeeType    FeeType    `xml:"cash_fee_type,omitempty"`
	CouponFee      Fee        `xml:"coupon_fee,omitempty"`     // 代金券金额
	CouponCount    int        `xml:"coupon_count,omitempty"`   // 代金券使用数量
	CouponTypes    []CData    `xml:"coupon_type,omitempty"`    // 代金券类型
	CouponIDs      []CData    `xml:"coupon_id,omitempty"`      // 代金券ID
	CouponFees     []CData    `xml:"coupon_fee,omitempty"`     // 单个代金券支付金额
	TransactionID  CData      `xml:"transaction_id,omitempty"` // 微信支付订单号
	TradeNo        CData      `xml:"out_trade_no,omitempty"`   // 商户订单号
	Attach         CData      `xml:"attach,omitempty"`         // 附加数据
	TimeEnd        Time       `xml:"time_end,omitempty"`
	TradeStateDesc CData      `xml:"trade_state_desc,omitempty"`
}

type CloseOrderRequest struct {
	TradeNo string `xml:"out_trade_no,omitempty"` // 商户订单号
}

type RefundRequest struct {
	TransactionID string  `xml:"transaction_id,omitempty"` // 微信订单号, 与 TradeNo 二选一
	TradeNo       string  `xml:"out_trade_no,omitempty"`   // 商户系统内部订单号, 与 TransactionID 二选一
	RefundNo      string  `xml:"out_refund_no,omitempty"`  // 商户系统内部的退款单号
	TotalFee      Fee     `xml:"total_fee,omitempty"`      // 订单金额
	RefundFee     Fee     `xml:"refund_fee,omitempty"`     // 退款金额
	RefundFeeType FeeType `xml:"refund_fee_type,omitempty"`
	RefundDesc    string  `xml:"refund_desc,omitempty"`
	RefundAccount string  `xml:"refund_account,omitempty"` // 退款资金来源
}

type RefundResponse struct {
	TransactionID       CData   `xml:"transaction_id,omitempty"`        // 微信订单号
	TradeNo             CData   `xml:"out_trade_no,omitempty"`          // 商户订单号
	RefundNo            CData   `xml:"out_refund_no,omitempty"`         // 商户退款单号
	RefundID            CData   `xml:"refund_id,omitempty"`             // 微信退款单号
	RefundFee           Fee     `xml:"refund_fee,omitempty"`            // 退款总金额,单位为分,可以做部分退款
	SettlementRefundFee Fee     `xml:"settlement_refund_fee,omitempty"` // 应结退款金额, 去掉非充值代金券退款金额后的退款金额，退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额应结退款金额
	TotalFee            Fee     `xml:"total_fee,omitempty"`             // 标价金额
	SettlementTotalFee  Fee     `xml:"settlement_total_fee,omitempty"`  // 应结订单金额, 去掉非充值代金券金额后的订单总金额，应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	RefundFeeType       FeeType `xml:"refund_fee_type,omitempty"`
	FeeType             FeeType `xml:"fee_type,omitempty"`
	CashFee             Fee     `xml:"cash_fee,omitempty"` // 现金支付金额
	CashFeeType         FeeType `xml:"cash_fee_type,omitempty"`
	CashRefundFee       Fee     `xml:"cash_refund_fee,omitempty"`   // 现金退款金额
	CouponRefundFee     Fee     `xml:"coupon_refund_fee,omitempty"` // 代金券退款总金额
	CouponRefundFees    []Fee   `xml:"coupon_refund_fee,omitempty"` // 单个代金券退款金额
	CouponRefundIDs     []CData `xml:"coupon_refund_id,omitempty"`  // 退款代金券ID
	CouponRefundCount   int     `xml:"coupon_refund_count,omitempty"`
}

// 四选一 refund_id > refund_no > transaction_id > trade_no
type QueryRefundRequest struct {
	TransactionID string `xml:"transaction_id,omitempty"` // 微信订单号
	TradeNo       string `xml:"out_trade_no,omitempty"`   // 商户订单号
	RefundNo      string `xml:"out_refund_no,omitempty"`  // 商户退款单号
	RefundID      string `xml:"refund_id,omitempty"`      // 微信退款单号
}

// not include coupon_type, coupon_refund_fee, coupon_refund_id
type QueryRefundResponse struct {
	TransactionID          CData          `xml:"transaction_id,omitempty"`       // 微信订单号
	TradeNo                CData          `xml:"out_trade_no,omitempty"`         // 商户订单号
	TotalFee               Fee            `xml:"total_fee,omitempty"`            // 订单金额
	SettlementFee          Fee            `xml:"settlement_total_fee,omitempty"` // 应结订单金额, 去掉非充值代金券金额后的订单总金额，应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	FeeType                FeeType        `xml:"fee_type,omitempty"`
	CashFee                Fee            `xml:"cash_fee,omitempty"`              // 现金支付金额
	RefundCount            int            `xml:"refund_count,omitempty"`          // 退款笔数
	RefundNos              []CData        `xml:"out_refund_no,omitempty"`         // 商户退款单号
	RefundIDs              []CData        `xml:"refund_id,omitempty"`             // 微信退款单号
	RefundChannels         []CData        `xml:"refund_channel,omitempty"`        // 退款渠道
	RefundFees             []Fee          `xml:"refund_fee,omitempty"`            // 申请退款金额
	SettlementRefundFees   []Fee          `xml:"settlement_refund_fee,omitempty"` // 退款金额
	CouponRefundFees       []Fee          `xml:"coupon_refund_fee,omitempty"`     // 总代金券退款金额, 代金券退款金额<=退款金额，退款金额-代金券或立减优惠退款金额为现金
	CouponRefundCount      []Fee          `xml:"coupon_refund_count,omitempty"`   // 退款代金券使用数量
	RefundStatus           []RefundStatus `xml:"refund_status,omitempty"`
	RefundAccounts         []CData        `xml:"refund_account,omitempty"`      // 退款资金来源
	RefundReceiverAccounts []CData        `xml:"refund_recv_account,omitempty"` // 退款入账账户
	RefundSuccessTime      []Time         `xml:"refund_success_time,omitempty"`
}

type DownloadBillRequest struct {
	DeviceInfo string `xml:"device_info,omitempty"`
	BillData   string `xml:"bill_data,omitempty"`
	BillType   `xml:"bill_type,omitempty"`
	TarType    `xml:"tar_type,omitempty"`
}

type NoticeResult struct {
	ReturnCode    CData `xml:"return_code,omitempty"`
	ReturnMessage CData `xml:"return_msg,omitempty"`
}

type PayNotice struct {
	TradeState
	DeviceInfo     CData     `xml:"device_info,omitempty"`  // 设备号
	OpenID         CData     `xml:"openid,omitempty"`       // 用户标识
	IsSubscribe    CData     `xml:"is_subscribe,omitempty"` // 是否关注公众账号
	TradeType      TradeType `xml:"trade_type,omitempty"`
	BankType       BankType  `xml:"bank_type,omitempty"`            // 付款银行
	TotalFee       Fee       `xml:"total_fee,omitempty"`            // 标价金额
	SettlementFee  Fee       `xml:"settlement_total_fee,omitempty"` // 应结订单金额, 当订单使用了免充值型优惠券后返回该参数，应结订单金额=订单金额-免充值优惠券金额。
	FeeType        FeeType   `xml:"fee_type,omitempty"`
	CashFee        Fee       `xml:"cash_fee,omitempty"` // 现金支付金额
	CashFeeType    FeeType   `xml:"cash_fee_type,omitempty"`
	CouponFee      Fee       `xml:"coupon_fee,omitempty"`     // 代金券金额
	CouponCount    int       `xml:"coupon_count,omitempty"`   // 代金券使用数量
	CouponTypes    []CData   `xml:"coupon_type,omitempty"`    // 代金券类型
	CouponIDs      []CData   `xml:"coupon_id,omitempty"`      // 代金券ID
	CouponFees     []CData   `xml:"coupon_fee,omitempty"`     // 单个代金券支付金额
	TransactionID  CData     `xml:"transaction_id,omitempty"` // 微信支付订单号
	TradeNo        CData     `xml:"out_trade_no,omitempty"`   // 商户订单号
	Attach         CData     `xml:"attach,omitempty"`         // 附加数据
	TimeEnd        Time      `xml:"time_end,omitempty"`
	TradeStateDesc CData     `xml:"trade_state_desc,omitempty"`
}

type RefundNotice struct {
	TransactionID         CData        `xml:"transaction_id,omitempty"`        // 微信订单号
	TradeNo               CData        `xml:"out_trade_no,omitempty"`          // 商户订单号
	RefundNo              CData        `xml:"out_refund_no,omitempty"`         // 商户退款单号
	RefundID              CData        `xml:"refund_id,omitempty"`             // 微信退款单号
	TotalFee              Fee          `xml:"total_fee,omitempty"`             // 订单金额
	SettlementFee         Fee          `xml:"settlement_total_fee,omitempty"`  // 应结订单金额, 去掉非充值代金券金额后的订单总金额，应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	RefundFees            Fee          `xml:"refund_fee,omitempty"`            // 申请退款金额
	SettlementRefundFees  Fee          `xml:"settlement_refund_fee,omitempty"` // 退款金额
	RefundStatus          RefundStatus `xml:"refund_status,omitempty"`
	SuccessTime           Time         `xml:"success_time,omitempty"`
	RefundReceiverAccount CData        `xml:"refund_recv_account,omitempty"` // 退款入账账户
	RefundAccount         CData        `xml:"refund_account,omitempty"`      // 退款资金来源
	RefundRequestSource   string       `xml:"refund_request_source,omitempty"`
}
