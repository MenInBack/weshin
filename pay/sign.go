// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
package pay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strings"

	"github.com/MenInBack/weshin/wx"
)

const (
	NonceLength = 32
)

// initialize random source with crypto/rand
func init() {
	buf := make([]byte, 8)
	n, e := crand.Read(buf)
	if n != 8 || e != nil {
		panic("init random source failed")
	}
	s := uint64(buf[0])<<56 | uint64(buf[1])<<48 | uint64(buf[2])<<40 | uint64(buf[3])<<32 |
		uint64(buf[4])<<24 | uint64(buf[5])<<16 | uint64(buf[6])<<8 | uint64(buf[7])
	rand.Seed(int64(s))
}

// Sign for payment, type of v must be a pointer to a struct with a "requestBase" field
func signRequest(value interface{}, key string, t SignType) (string, error) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	fields := parseStructFields(val)
	s, err := sign(fields, key, t)
	if err != nil {
		return "", err
	}

	return s, nil
}

func checkSignature(value interface{}, key string, typ SignType) error {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	base, ok := val.FieldByName("ResponseBase").Interface().(ResponseBase)
	if !ok {
		return wx.WeshinError{Detail: "empty ResponseBase"}
	}

	fields := parseStructFields(val)
	s, err := sign(fields, key, typ)
	if err != nil {
		return err
	}
	if s != base.Sign.Data {
		return wx.WeshinError{Detail: fmt.Sprint("mismatched signature, expected: ", s)}
	}
	return nil
}

func parseStructFields(val reflect.Value) []field {
	typ := val.Type()
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
		typ = typ.Elem()
	}

	// extract xml name and value for signature
	fields := make([]field, 0, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		var name, value string
		t := typ.Field(i)
		v := val.Field(i)

		switch t.Name {
		case "XMLName":
			continue
		case "SignType", "Sign":
			continue
		case "RequestBase", "ResponseBase":
			fields = append(fields, parseStructFields(v)...)
			continue
		}

		tags := strings.Split(t.Tag.Get("xml"), ",")
		if len(tags) > 0 {
			name = tags[0]
			if name == "" {
				name = t.Name
			}
		}
		switch v.Kind() {
		case reflect.String:
			value = v.Interface().(string)
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			value = v.String()
		case reflect.Struct:
			value = v.MethodByName("String").Call([]reflect.Value{})[0].Interface().(string)
		}

		if value == "" {
			continue
		}

		fields = append(fields, field{name, value})
	}

	return fields
}

func sign(fields []field, key string, typ SignType) (string, error) {
	if len(fields) <= 0 {
		return "", wx.WeshinError{Detail: "empty query parameter"}
	}
	sort.Slice(fields, func(i, j int) bool { return strings.Compare(fields[i].name, fields[j].name) < 0 })
	fields = append(fields, field{"key", key})

	var buf bytes.Buffer
	for _, f := range fields {
		buf.WriteString(f.name)
		buf.WriteByte('=')
		buf.WriteString(f.value)
		buf.WriteByte('&')
	}
	str := bytes.TrimSuffix(buf.Bytes(), []byte{'&'})

	switch typ {
	case MD5:
		return fmt.Sprintf("%X", md5.Sum(str)), nil
	case HMAC:
		mac := hmac.New(sha256.New, []byte(key))
		mac.Write(str)
		return fmt.Sprintf("%X", mac.Sum(nil)), nil
	}
	return "", wx.WeshinError{Detail: "unknown sign type"}
}

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func randomString(n int) string {
	s := make([]byte, 0, n)
	cache := rand.Uint64()
	remain := 64
	for i := 0; i < n; remain -= 6 {
		if remain < 6 {
			cache, remain = rand.Uint64(), 64
		}
		b := cache & (1<<6 - 1)
		if b < 62 {
			s = append(s, chars[b])
			i++
		}
		cache >>= 6
	}
	return string(s)
}

type field struct {
	name  string
	value string
}
