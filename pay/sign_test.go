package pay

import (
	"fmt"
	"reflect"
	"testing"
)

type typ struct {
	DeviceInfo string `xml:"device_info"`
	Body       string `xml:"body"`
}

const key = "192006250b4c09247ec02edce69f6a2d"

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=20_1
func TestSign(t *testing.T) {
	v := typ{
		DeviceInfo: "1000",
		Body:       "test",
	}

	fields := structToFields(reflect.ValueOf(v))
	s, err := sign(fields, key, MD5)
	if err != nil {
		t.Error(err)
	}

	s, e := sign(fields, key, HMAC)
	if e != nil {
		t.Error(e)
	}

	fmt.Println(s)
}

func TestDecodeMessage(t *testing.T) {
	info := `T87GAHG17TGAHG1TGHAHAHA1Y1CIOA9UGJH1GAHV871HAGAGQYQQPOOJMXNBCXBVNMNMAJAA`

	data, e := decodeNoticeMessage(info, key)
	if e != nil {
		t.Error(e)
		return
	}

	fmt.Println(data)

}
