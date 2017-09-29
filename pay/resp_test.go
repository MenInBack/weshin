package pay

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"testing"
)

func TestResponsePreOrder(t *testing.T) {
	data := `<xml>
	<return_code><![CDATA[SUCCESS]]></return_code>
	<return_msg><![CDATA[OK]]></return_msg>
	<appid><![CDATA[wx2421b1c4370ec43b]]></appid>
	<mch_id><![CDATA[10000100]]></mch_id>
	<nonce_str><![CDATA[IITRi8Iabbblz1Jc]]></nonce_str>
	<openid><![CDATA[oUpF8uMuAJO_M2pxb1Q9zNjWeS6o]]></openid>
	<sign><![CDATA[7921E432F65EB8ED0CE9755F0E86D72F]]></sign>
	<result_code><![CDATA[SUCCESS]]></result_code>
	<prepay_id><![CDATA[wx201411101639507cbf6ffd8b0779950874]]></prepay_id>
	<trade_type><![CDATA[JSAPI]]></trade_type>
 </xml>`

	resp := struct {
		*ResponseBase
		*PreOrderResponse
	}{
		new(ResponseBase),
		new(PreOrderResponse),
	}
	e := xml.Unmarshal([]byte(data), &resp)
	if e != nil {
		t.Error(e)
	}
	fmt.Printf("base: %+v, resp: %+v", resp.ResponseBase, resp.PreOrderResponse)
}

func TestResponseQueryOrder(t *testing.T) {
	data := `<xml>
	<return_code><![CDATA[SUCCESS]]></return_code>
	<return_msg><![CDATA[OK]]></return_msg>
	<appid><![CDATA[wx2421b1c4370ec43b]]></appid>
	<mch_id><![CDATA[10000100]]></mch_id>
	<device_info><![CDATA[1000]]></device_info>
	<nonce_str><![CDATA[TN55wO9Pba5yENl8]]></nonce_str>
	<sign><![CDATA[BDF0099C15FF7BC6B1585FBB110AB635]]></sign>
	<result_code><![CDATA[SUCCESS]]></result_code>
	<openid><![CDATA[oUpF8uN95-Ptaags6E_roPHg7AG0]]></openid>
	<is_subscribe><![CDATA[Y]]></is_subscribe>
	<trade_type><![CDATA[MICROPAY]]></trade_type>
	<bank_type><![CDATA[CCB_DEBIT]]></bank_type>
	<total_fee>1</total_fee>
	<fee_type><![CDATA[CNY]]></fee_type>
	<transaction_id><![CDATA[1008450740201411110005820873]]></transaction_id>
	<out_trade_no><![CDATA[1415757673]]></out_trade_no>
	<attach><![CDATA[订单额外描述]]></attach>
	<time_end><![CDATA[20141111170043]]></time_end>
	<trade_state><![CDATA[SUCCESS]]></trade_state>
 </xml>`

	resp := struct {
		*ResponseBase
		*QueryOrderResponse
	}{
		new(ResponseBase),
		new(QueryOrderResponse),
	}
	e := xml.Unmarshal([]byte(data), &resp)
	if e != nil {
		t.Error(e)
	}
	fmt.Printf("base: %+v, resp: %+v", resp.ResponseBase, resp.QueryOrderResponse)
}

func TestResponseBase(t *testing.T) {
	data := `
	<xml>
	<return_code><![CDATA[SUCCESS]]></return_code>
	<return_msg><![CDATA[OK]]></return_msg>
	<appid><![CDATA[wx2421b1c4370ec43b]]></appid>
	<mch_id><![CDATA[10000100]]></mch_id>
	<nonce_str><![CDATA[BFK89FC6rxKCOjLX]]></nonce_str>
	<sign><![CDATA[72B321D92A7BFA0B2509F3D13C7B1631]]></sign>
	<result_code><![CDATA[SUCCESS]]></result_code>
	<result_msg><![CDATA[OK]]></result_msg>
 </xml>`
	resp := new(ResponseBase)
	e := xml.Unmarshal([]byte(data), &resp)
	if e != nil {
		t.Error(e)
	}
	fmt.Printf("resp: %+v", resp)
}

func TestParseResponse(t *testing.T) {
	data := `<xml>
	<appid><![CDATA[wx2421b1c4370ec43b]]></appid>
	<mch_id><![CDATA[10000100]]></mch_id>
	<nonce_str><![CDATA[TeqClE3i0mvn3DrK]]></nonce_str>
	<out_refund_no_0><![CDATA[1415701182]]></out_refund_no_0>
	<out_trade_no><![CDATA[1415757673]]></out_trade_no>
	<refund_count>1</refund_count>
	<refund_fee_0>1</refund_fee_0>
	<refund_id_0><![CDATA[2008450740201411110000174436]]></refund_id_0>
	<refund_status_0><![CDATA[PROCESSING]]></refund_status_0>
	<result_code><![CDATA[SUCCESS]]></result_code>
	<return_code><![CDATA[SUCCESS]]></return_code>
	<return_msg><![CDATA[OK]]></return_msg>
	<sign><![CDATA[1F2841558E233C33ABA71A961D27561C]]></sign>
	<transaction_id><![CDATA[1008450740201411110005820873]]></transaction_id>
	</xml>
	`

	fields, e := parseToFields(bytes.NewBuffer([]byte(data)))
	if e != nil {
		t.Error(t)
	}

	fmt.Println(fields)
}

func TestRefundResponse(t *testing.T) {
	data := `<xml>
<appid><![CDATA[wx2421b1c4370ec43b]]></appid>
<mch_id><![CDATA[10000100]]></mch_id>
<nonce_str><![CDATA[TeqClE3i0mvn3DrK]]></nonce_str>
<out_refund_no_0><![CDATA[1415701182]]></out_refund_no_0>
<out_trade_no><![CDATA[1415757673]]></out_trade_no>
<refund_count>1</refund_count>
<refund_fee_0>1</refund_fee_0>
<refund_id_0><![CDATA[2008450740201411110000174436]]></refund_id_0>
<refund_status_0><![CDATA[PROCESSING]]></refund_status_0>
<result_code><![CDATA[SUCCESS]]></result_code>
<return_code><![CDATA[SUCCESS]]></return_code>
<return_msg><![CDATA[OK]]></return_msg>
<sign><![CDATA[1F2841558E233C33ABA71A961D27561C]]></sign>
<transaction_id><![CDATA[1008450740201411110005820873]]></transaction_id>
</xml>
`

	resp := new(QueryRefundResponse)
	m := &MerchantInfo{
		PaymentKey: key,
	}
	e := m.handleResponse(bytes.NewBuffer([]byte(data)), resp)

	if e != nil {
		t.Error(e)
	}
	fmt.Printf("resp: %+v", resp)
}
