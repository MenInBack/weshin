package component

// wechat thirdparty component api
// https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1453779503&token=&lang=zh_CN

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/MenInBack/weshin/crypto"
	"github.com/MenInBack/weshin/wx"
)

func (c *Component) StartNotifyHandler() {
	http.HandleFunc("/", c.AuthMessageHandler)
	c.NotifyErrors = make(chan error)
	go http.ListenAndServe(c.Address.Address, nil)
}

// AuthMessageHandler responses to messages from wechat for verify ticket and thirdparty authorization events
func (c *Component) AuthMessageHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		c.NotifyErrors <- wx.NotifyError{err}
		return
	}

	// decrypt
	p := getParameter(req)
	encoding, err := crypto.New(c.EncodingAESKey, c.GetAccessToken(), c.AppID)
	if err != nil {
		c.NotifyErrors <- wx.NotifyError{err}
		return
	}
	data, err := encoding.Decrypt(body, p.signature, p.nonce, p.timestamp)
	if err != nil {
		c.NotifyErrors <- wx.NotifyError{err}
		return
	}
	w.Write([]byte("success"))

	// unmarshal two different type int
	var reqBody struct {
		*ComponentVerifyTicket
		*AuthorizationNotifyBody
	}

	err = xml.Unmarshal(data, &reqBody)
	if err != nil {
		c.NotifyErrors <- wx.NotifyError{err}
		return
	}

	// ticket notify
	if reqBody.ComponentVerifyTicket != nil {
		go c.SetVerifyTicket(&wx.APITicket{
			Typ:      wx.TicketTypeVerify,
			Ticket:   reqBody.ComponentVerifyTicket.ComponentVerifyTicket,
			CreateAt: reqBody.ComponentVerifyTicket.CreateTime,
		})
		return
	}

	// authorization notify
	if reqBody.AuthorizationNotifyBody != nil {
		switch reqBody.AuthorizationNotifyBody.InfoType {
		case NotifyTypeAuthorized, NotifyTypeUpdateAuthorized:
			go func() {
				tokenInfo, err := c.MPAuthorize(reqBody.AuthorizationNotifyBody.AppID, 0)
				if err != nil {
					c.NotifyErrors <- wx.NotifyError{err}
					return
				}
				go c.SetAuthorizationInfo(tokenInfo)
			}()
		case NotifyTypeUnauthorized:
			go c.ClearAuthorizerToken(reqBody.AuthorizationCode.AppID)
		}
	}
}

type messageParameter struct {
	timestamp   string
	nonce       string
	encryptType string
	signature   string
}

func getParameter(req *http.Request) *messageParameter {
	queries := req.URL.Query()
	return &messageParameter{
		timestamp:   queries.Get("timestamp"),
		nonce:       queries.Get("nonce"),
		encryptType: queries.Get("encrypt_type"),
		signature:   queries.Get("msg_signature"),
	}
}

// https://api.weixin.qq.com/cgi-bin/component/api_component_token
func (c *Component) GrantComponentAccessToken(timeout int) (token *ComponentAccessToken, err error) {
	req := wx.HttpClient{
		Path:        accessTokenURI,
		ContentType: "application/json",
		Timeout:     timeout,
	}

	body := struct {
		ComponentAppID        string `json:"component_appid"`
		ComponentAppSecret    string `json:"component_appsecret"`
		ComponentVerifyTicket string `json:"component_verify_token"`
	}{
		c.AppID,
		c.Secret,
		c.GetVerifyTicket(),
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	token = new(ComponentAccessToken)
	err = req.DoPost(buf, token)
	if err != nil {
		return nil, err
	}

	go c.SetAccessToken(token.Token, token.ExpiresIn)

	return token, nil
}

// https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=xxx
func (c *Component) GetPreAuthCode(timeout int) (code *PreAuthCode, err error) {
	req := wx.HttpClient{
		Path:        preAuthCodeURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{{
			"component_access_token", c.GetAccessToken(),
		}},
		Timeout: timeout,
	}

	body := struct {
		ComponentAppID string `json:"component_appid"`
	}{c.AppID}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	code = new(PreAuthCode)
	err = req.DoPost(buf, code)
	if err != nil {
		return nil, err
	}

	return code, nil
}

//https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=xxxx&pre_auth_code=xxxxx&redirect_uri=xxxx
func (c *Component) JumpToOAuth(preAuthCode string) string {
	uri := bytes.NewBufferString(authorizeURI)
	uri.WriteString("?component_appid=")
	uri.WriteString(c.AppID)
	uri.WriteString("&pre_auth_code=")
	uri.WriteString(preAuthCode)
	uri.WriteString("&redirect_uri=")
	uri.WriteString(url.QueryEscape(c.Address.AuthRedirectURI))
	return uri.String()
}

// https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=xxxx
func (c *Component) MPAuthorize(authorizationCode string, timeout int) (auth *AuthorizationTokenInfo, err error) {
	req := wx.HttpClient{
		Path:        authorizationInfoURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{{
			"component_access_token", c.GetAccessToken(),
		}},
		Timeout: timeout,
	}

	body := struct {
		ComponentAppID    string `json:"component_appid"`
		AuthorizationCode string `json:"authorization_code"`
	}{
		c.AppID,
		authorizationCode,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	auth = new(AuthorizationTokenInfo)
	err = req.DoPost(buf, auth)
	if err != nil {
		return nil, err
	}

	go c.SetAuthorizerToken(&auth.AuthorizationToken)

	return auth, nil
}

// https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=xxxxx
func (c *Component) RefreshAuthorizerToken(authorizerAppID, refreshToken string, timeout int) (token *AuthorizerToken, err error) {
	req := wx.HttpClient{
		Path:        authorizerTokenURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{{
			"component_access_token", c.GetAccessToken(),
		}},
		Timeout: timeout,
	}

	body := struct {
		ComponentAppID         string `json:"component_appid"`
		AuthorizerAppID        string `json:"authorizer_appid"`
		AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
	}{
		c.AppID,
		authorizerAppID,
		refreshToken,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	token = new(AuthorizerToken)
	err = req.DoPost(buf, token)
	if err != nil {
		return nil, err
	}
	token.AppID = authorizerAppID

	go c.SetAuthorizerToken(token)

	return token, nil
}

// https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?component_access_token=xxxx
func (c *Component) GetAuthorizerInfo(authorizerAppID string, timeout int) (info *Authorizer, err error) {
	req := wx.HttpClient{
		Path:        authorizerInfoURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{{
			"component_access_token", c.GetAccessToken(),
		}},
		Timeout: timeout,
	}

	body := struct {
		ComponentAppID  string `json:"component_appid"`
		AuthorizerAppID string `json:"authorizer_appid"`
	}{
		c.AppID,
		authorizerAppID,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	info = new(Authorizer)
	err = req.DoPost(buf, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_option?component_access_token=xxxx
func (c *Component) GetAuthorizerOption(authorizerAppID, optionName string, timeout int) (option *AuthorizerOption, err error) {
	req := wx.HttpClient{
		Path:        getAuthorizerOptionURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{{
			"component_access_token", c.GetAccessToken(),
		}},
		Timeout: timeout,
	}

	body := struct {
		ComponentAppID  string `json:"component_appid"`
		AuthorizerAppID string `json:"authorizer_appid"`
		OptionName      string `json:"option_name"`
	}{
		c.AppID,
		authorizerAppID,
		optionName,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	option = new(AuthorizerOption)
	err = req.DoPost(buf, option)
	if err != nil {
		return nil, err
	}

	return option, nil
}

// https://api.weixin.qq.com/cgi-bin/component/api_set_authorizer_option?component_access_token=xxxx
func (c *Component) SetAuthorizerOption(option *AuthorizerOption, timeout int) error {
	req := wx.HttpClient{
		Path:        setAuthorizerOptionURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{{
			"component_access_token", c.GetAccessToken(),
		}},
		Timeout: timeout,
	}

	body := struct {
		ComponentAppID  string `json:"component_appid"`
		AuthorizerAppID string `json:"authorizer_appid"`
		OptionName      string `json:"option_name"`
		OptionValue     string `json:"option_value"`
	}{
		c.AppID,
		option.AuthorizerAppID,
		option.OptionName,
		option.OptionValue,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	err = req.DoPost(buf, nil)
	if err != nil {
		return err
	}

	return nil
}
