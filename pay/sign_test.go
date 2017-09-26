package pay

import (
	"encoding/xml"
	"fmt"
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

	e := signRequest(&v, key, HMAC)
	if e != nil {
		t.Error(e)
	}

	fmt.Println(v)

	d, e := xml.MarshalIndent(v, "", "  ")
	if e != nil {
		t.Error(e)
	}
	fmt.Println(string(d))
}

// quite slow with reflect
// 3782 ns/op on Mac
func BenchmarkSign(b *testing.B) {
	v := typ{}

	for i := 0; i < b.N; i++ {
		signRequest(&v, key, MD5)
	}
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
