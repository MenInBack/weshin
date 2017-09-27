package pay

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"testing"
)

type typ struct {
	RequestBase
	DeviceInfo string `xml:"device_info"`
	Body       string `xml:"body"`
}

const key = "192006250b4c09247ec02edce69f6a2d"

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=20_1
func TestSign(t *testing.T) {
	v := typ{
		RequestBase: RequestBase{
			AppID:      "wxd930ea5d5a258f4f",
			MerchantID: "10000100",
		},
		DeviceInfo: "1000",
		Body:       "test",
	}

	fields := parseStruct(reflect.ValueOf(v))
	s, err := sign(fields, key, MD5)
	if err != nil {
		t.Error(err)
	}

	s, e := sign(fields, key, HMAC)
	if e != nil {
		t.Error(e)
	}

	fmt.Println(s)

	d, e := xml.MarshalIndent(v, "", "  ")
	if e != nil {
		t.Error(e)
	}
	fmt.Println(string(d))
}

func TestRandomString(t *testing.T) {
	s := randomString(16)
	fmt.Println(s)
}

func BenchmarkRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = randomString(32)
	}
}
