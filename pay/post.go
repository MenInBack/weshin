package pay

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/MenInBack/weshin/wx"
)

// for https connection
var transport *http.Transport

// SetCertificationFile initializes tls certification with:
// application_cert.pem &
// application_key.pem.
// should be called in the begining.
// see https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
func SetCertificationFile(cert, key string) error {
	keyPair, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return err
	}

	certData, err := ioutil.ReadFile(cert)
	if err != nil {
		return err
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(certData); !ok {
		return wx.WeshinError{Detail: "AppendCertsFromPEM error"}
	}

	config := &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{keyPair},
	}

	transport = &http.Transport{
		TLSClientConfig: config,
	}

	return nil
}

func (m *MerchantInfo) postXML(path string, request, response interface{}, safe bool) error {
	body, e := m.prepareRequest(request)
	if e != nil {
		return e
	}
	if verbose {
		fmt.Println("request path: ", path, " body: ", string(body))
	}

	req, err := http.NewRequest("POST", path, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/xml")

	c := http.Client{
		Timeout: 10 * time.Second,
	}

	if safe {
		// will supply certification with request and check server certification
		c.Transport = transport
	}

	resp, e := c.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	if verbose {
		fmt.Println("response: ", resp)
	}

	if resp.StatusCode != http.StatusOK {
		return wx.HttpError{
			State: resp.StatusCode,
		}
	}
	if e = m.handleResponse(resp.Body, response); e != nil {
		return e
	}

	return nil
}

// sign and marshal request
func (m *MerchantInfo) prepareRequest(req interface{}) ([]byte, error) {
	fields := structToFields(reflect.ValueOf(req))

	fields = append(fields,
		field{"app_id", m.AppID},
		field{"mch_id", m.MerchantID},
		field{"nonce", randomString(NonceLength)})

	s, err := sign(fields, m.PaymentKey, MD5)
	if err != nil {
		return nil, err
	}
	fields = append(fields, field{"sign", s}, field{"sign_type", "MD5"})

	if verbose {
		fmt.Println("request fields: ", fields)
	}

	return marshalRequest(fields)
}

func structToFields(val reflect.Value) []field {
	typ := val.Type()
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
		typ = typ.Elem()
	}

	// extract xml name and value for signing
	fields := make([]field, 0, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		var name, value string
		f := typ.Field(i)
		v := val.Field(i)

		// ignore zero field
		if !v.IsValid() {
			continue
		}

		switch f.Name {
		case "XMLName", "SignType", "Sign":
			continue
		}

		tags := strings.Split(f.Tag.Get("xml"), ",")
		if len(tags) > 0 {
			name = tags[0]
			if name == "" {
				name = f.Name
			}
		}

		if v.Type().Implements(reflect.ValueOf(new(fmt.Stringer)).Elem().Type()) {
			value = v.MethodByName("String").Call([]reflect.Value{})[0].Interface().(string)
		} else {
			switch v.Kind() {
			// no slice right now
			case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
				value = v.String()
			case reflect.Struct, reflect.Ptr:
				fields = append(fields, structToFields(v)...)
				continue
			}
		}

		if value == "" {
			continue
		}

		fields = append(fields, field{name, value})
	}

	return fields
}

func marshalRequest(fields []field) ([]byte, error) {
	buf := bytes.NewBufferString("<xml>")
	for _, f := range fields {
		if f.value == "" {
			continue
		}
		buf.WriteByte('<')
		buf.WriteString(f.name)
		buf.WriteByte('>')
		buf.WriteString(f.value)
		buf.WriteString("</")
		buf.WriteString(f.name)
		buf.WriteByte('>')
	}
	buf.WriteString("</xml>")
	return buf.Bytes(), nil
}

func checkResult(r ResponseBase) error {
	if r.ReturnCode.Data != "SUCCESS" {
		return wx.WeshinError{Detail: fmt.Sprintf("pay request failed: [%s]%s", r.ReturnCode.Data, r.ReturnMessage.Data)}
	}
	if r.ResultCode.Data != "SUCCESS" {
		return wx.WeshinError{Detail: fmt.Sprintf("pay response failed: [%s]%s", r.ErrorCode.Data, r.ErrorDescription.Data)}
	}
	return nil
}

func (m *MerchantInfo) handleResponse(body io.Reader, response interface{}) error {
	fields, e := parseToFields(body)
	if e != nil {
		return e
	}

	if verbose {
		fmt.Println("xml to fields: ", fields)
	}

	// check signature
	var signature string
	var signType SignType
	var ok bool
	if signature, ok = fields["sign"]; ok {
		delete(fields, "sign")
	}
	if st, ok := fields["sign_type"]; ok {
		signType = SignType(st)
		delete(fields, "sign_type")
	}
	if signType == "" {
		signType = MD5
	}

	fs := make([]field, 0, len(fields))
	for n, v := range fields {
		fs = append(fs, field{n, v})
	}

	s, e := sign(fs, m.PaymentKey, signType)
	if e != nil {
		return e
	}
	if s != signature {
		return wx.WeshinError{Detail: "response signature mismatch"}
	}

	// check result
	var resp struct {
		Base ResponseBase
		Resp interface{}
	}
	resp.Resp = response

	if e = composeStruct(fields, reflect.ValueOf(resp)); e != nil {
		return e
	}
	if e = checkResult(resp.Base); e != nil {
		return e
	}

	return nil
}

func parseToFields(body io.Reader) (fields map[string]string, e error) {
	tokens := make([]xml.Token, 0, 4) // use as stack
	fields = make(map[string]string)

	// parse a xml element
	parseXML := func(e xml.EndElement) error {
		var name, value string
		n := len(tokens)
		if n < 1 {
			return wx.WeshinError{Detail: "unexpected EndElement in response xml"}
		}

		if _, ok := tokens[n-1].(xml.StartElement); ok {
			tokens = tokens[:n-1]
			return nil // empty xml element
		}
		if t, ok := tokens[n-1].(xml.CharData); ok {
			value = string(t.Copy())
		} else if t, ok := tokens[n-1].(xml.Directive); ok {
			value = string(t.Copy())
		} else {
			return nil
			// return wx.WeshinError{Detail: "expect Directive or CharData before an EndElement"}
		}

		t, ok := tokens[n-2].(xml.StartElement)
		if !ok {
			return wx.WeshinError{Detail: "expect StartElement before an Directive or CharData"}
		}
		if t.Name != e.Name {
			return wx.WeshinError{Detail: fmt.Sprintf("mismatched StartElement %s with EndElement %s", t.Name, e.Name)}
		}
		name = t.Name.Local
		tokens = tokens[:n-2]

		switch name {
		case "xml":
			return nil
		}

		fields[name] = value
		return nil
	}

	decoder := xml.NewDecoder(body)
	for t, e := decoder.Token(); e == nil; t, e = decoder.Token() {
		switch t.(type) {
		case xml.StartElement:
			if t.(xml.StartElement).Name.Local == "xml" {
				continue
			}
			tokens = append(tokens, t)
		case xml.CharData:
			if len(tokens) == 0 {
				continue
			}
			if _, ok := tokens[len(tokens)-1].(xml.StartElement); ok {
				tokens = append(tokens, t.(xml.CharData).Copy())
			}
		case xml.Directive:
			if len(tokens) == 0 {
				continue
			}
			if _, ok := tokens[len(tokens)-1].(xml.StartElement); ok {
				tokens = append(tokens, t.(xml.Directive).Copy())
			}
		case xml.EndElement:
			if t.(xml.EndElement).Name.Local == "xml" {
				continue
			}
			if e = parseXML(t.(xml.EndElement)); e != nil {
				return nil, e
			}
		}
	}

	return
}

// fields to struct
func composeStruct(fields map[string]string, val reflect.Value) error {
	typ := val.Type()
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		t := typ.Field(i)
		v := val.Field(i)

		switch t.Name {
		case "XMLName":
			continue
		case "SignType", "Sign":
			continue
		case "RequestBase", "ResponseBase":
			if e := composeStruct(fields, v); e != nil {
				return e
			}
			continue
		}

		var name string
		tags := strings.Split(t.Tag.Get("xml"), ",")
		if len(tags) > 0 {
			name = tags[0]
			if name == "" {
				name = t.Name
			}
		}

		if e := parseField(name, fields, v); e != nil {
			return e
		}
	}

	return nil
}

func parseField(name string, fields map[string]string, v reflect.Value) error {
	val := v
	value, ok := fields[name]
	if !ok {
		return nil
	}

	// wechat specified slice first
	if val.Kind() == reflect.Slice {
		return parseSlice(name, fields, val)
	}

	// customized Unmarshaler next
	if val.Type().Implements(reflect.ValueOf(new(Unstringer)).Elem().Type()) {
		fmt.Println("unmarshalling: ", name)

		if val.Kind() == reflect.Ptr {
			val = reflect.New(val.Type().Elem())
			v.Set(val)
		}

		// start := xml.StartElement{
		// 	Name: xml.Name{
		// 		Local: name,
		// 	},
		// }

		// buf := bytes.NewBuffer([]byte("<"))
		// buf.WriteString(name)
		// buf.WriteByte('>')
		// buf.WriteString(value)
		// buf.WriteString("</")
		// buf.WriteString(name)
		// buf.WriteByte('>')
		// decoder := xml.NewDecoder(buf)

		// ret := val.MethodByName("UnmarshalXML").Call([]reflect.Value{reflect.ValueOf(decoder), reflect.ValueOf(start)})

		ret := val.MethodByName("Unstring").Call([]reflect.Value{reflect.ValueOf(value)})
		if len(ret) > 0 {
			if e, ok := ret[0].Interface().(error); ok && e != nil {
				return e
			}
		}
		return nil
	}

	// default Unmarshaler
	switch val.Kind() {
	case reflect.String:
		val.SetString(value)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		i, _ := strconv.ParseInt(value, 10, 64)
		val.Set(reflect.ValueOf(i))
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, _ := strconv.ParseUint(value, 10, 64)
		val.Set(reflect.ValueOf(u))
	case reflect.Float32, reflect.Float64:
		f, _ := strconv.ParseFloat(value, 64)
		val.Set(reflect.ValueOf(f))
	case reflect.Bool:
		b, _ := strconv.ParseBool(value)
		val.Set(reflect.ValueOf(b))
	default:
		return wx.WeshinError{Detail: "unknown type to unmarshal"}
	}

	return nil
}

// parse wechat pay specified slice
func parseSlice(name string, fields map[string]string, val reflect.Value) error {
	namer := func(i int) string {
		return fmt.Sprintf("%s_%d", name, i)
	}
	for i := 0; ; i++ {
		n := namer(i)
		if _, ok := fields[n]; !ok {
			break
		}

		v := reflect.New(val.Elem().Type())
		if e := parseField(n, fields, v); e != nil {
			return e
		}

		if val.Cap() < i+1 {
			if val.Cap() == 0 {
				val.SetCap(1)
			} else {
				val.SetCap(2 * val.Cap())
			}
		}
		val.SetLen(i + 1)
		val.Index(i).Set(v)
	}

	return nil
}
