package pay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var m MerchantInfo = MerchantInfo{
	AppID:      "wx2421b1c4370ec43b",
	MerchantID: "10000100",
	PaymentKey: key,
}

func init() {
	verbose = true
}

func TestEncodeJsonInXML(t *testing.T) {
	v := PreOrderRequest{
		DeviceInfo: "1000",
		SceneInfo: &SceneInfo{
			ID:       "aaaa",
			Name:     "7-11",
			AreaCode: "310000",
			Address:  "Minsheng Rd.",
		},
	}

	data, e := m.prepareRequest(v)
	if e != nil {
		t.Error(e)
	}
	fmt.Println(string(data))
}

func TestDecodeJsonInXML(t *testing.T) {
	data := []byte(`<xml><app_id>appid001122</app_id><device_info>1000</device_info><mch_id>100100</mch_id><nonce>IWQOYXUmSO8l2GBmAYQo3HOfqMETSmaY</nonce><scene_info>{"id":"aaaa","name":"7-11","area_code":"310000","address":"Minsheng Rd."}</scene_info><sign>4AE62ED2F1D999F96B0E9A66D47C3770</sign><sign_type>MD5</sign_type></xml>`)

	n := new(PreOrderRequest)
	fields, e := parseToFields(bytes.NewBuffer(data))
	if e != nil {
		t.Error(e)
	}
	fmt.Println(fields)

	e = composeStruct(fields, reflect.ValueOf(n))
	if e != nil {
		t.Error(e)
	}
	fmt.Println(n)
}

func TestSceneInfo(t *testing.T) {
	s := SceneInfo{}
	data, e := json.Marshal(s)
	if e != nil {
		t.Error(e)
	}
	fmt.Println(string(data))
}

func TestUnmarshalSceneInfo(t *testing.T) {
	data := []byte(`{"id":"aaaa","name":"7-11","area_code":"310000","address":"Minsheng Rd."}`)
	s := new(SceneInfo)
	e := json.Unmarshal(data, s)
	if e != nil {
		t.Error(e)
	}
	fmt.Printf("value: %+v\n", s)
}
