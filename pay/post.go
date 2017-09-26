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

	if resp.StatusCode != http.StatusOK {
		return wx.HttpError{
			State: resp.StatusCode,
		}
	}
	if e = m.parseResponse(resp.Body, response); e != nil {
		return e
	}

	return nil
}

// sign and marshal request
func (m *MerchantInfo) prepareRequest(req interface{}) ([]byte, error) {
	var reqAll struct {
		RequestBase
		Req interface{}
	}
	reqAll.Req = req
	reqAll.RequestBase = RequestBase{
		AppID:      m.AppID,
		MerchantID: m.MerchantID,
		Nonce:      randomString(NonceLength),
	}

	s, e := signRequest(&reqAll, m.PaymentKey, MD5)
	if e != nil {
		return nil, e
	}

	reqAll.RequestBase.Sign = s
	reqAll.RequestBase.SignType = MD5
	return xml.Marshal(reqAll)
}

// check signature and parse other fields of response
func (m *MerchantInfo) parseResponse(body io.Reader, response interface{}) error {
	var respAll struct {
		ResponseBase
		Resp interface{}
	}
	respAll.Resp = response

	data, e := ioutil.ReadAll(body)
	if e != nil {
		return e
	}
	if verbose {
		fmt.Println("response body: ", string(data))
	}

	if e := xml.Unmarshal(data, &respAll); e != nil {
		return e
	}
	if e := checkSignature(&respAll, m.PaymentKey, MD5); e != nil {
		return e
	}
	if e := checkResult(respAll.ResponseBase); e != nil {
		return e
	}
	return nil
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
