package component

// https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1453779503&token=&lang=zh_CN

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"bytes"

	"github.com/MenInBack/weshin/wx"
)

type NotifyConfig struct{
	Address string
	VerifyTickPath string
	AuthorizationPath string
}

var notifyConfig NotifyConfig

func SetNotifyAddress(conf *NotifyConfig){
	notifyConfig = conf
}

func StartNotifyHandler(address, path string)error {
	if notifyConfig == nil{
		return wx.ParameterError{InvalidParameter: "notify config"}
	}
	log.Println("starting http service on: ", address)
	http.HandleFunc(notifyConfig.VerifyTickPath, verifyTicketHandler)
	http.HandleFunc(notifyConfig.AuthorizationPath, authorizationNotifyHandler)
	go http.ListenAndServe(address, nil)

	return nil
}

func verifyTicketHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("got verify ticket req")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error: ", err)
		return
	}
	// todo: decrypt
	var reqBody ComponentVerifyTicket
	err = xml.Unmarshal(body, &reqBody)
	if err != nil {
		log.Println("xml.Unmarshal error: ", err)
		return
	}

	log.Printf("request body: %+v\n", reqBody)
	w.Write("success")

	go ticketStorage.Set(reqBody.ComponentVerifyTicket, reqBody.CreateTime, 0)
}

func authorizationNotifyHandler(w http.ResponseWriter, req *http.Request){
	log.Println("got authorization notify")

	body, err := ioutile.ReadAll(req.Body)
	if err != nil{
		log.Println("ioutile.ReadAll error: ", err)
		return
	}

	// todo: decrypt
	var reqBody = new(authorizationNotifyBody)
	err = xml.Unmarshal(body, reqBody)
	if err != nil{
		log.Println("xml.Unmarshal error: " ,err)
		return
	}

	log.Printf("request body: %+v\n", reqBody)
	w.Write("success")

	switch reqBody.InfoType{
	case NotifyTypeAuthorized, NotifyTypeUpdateAuthorized:
		go authorizationStorage.Set(&reqBody.AuthorizationCode)
	case NotifyTypeUnauthorized:
		go authorizationStorage.Clear(reqBody.AuthorizerAppID)	
	}

}

// https://api.weixin.qq.com/cgi-bin/component/api_component_token
func GetComponentAccessToken(timeout int) (token *ComponentAccessToken, err error){
	req := wx.HttpClient{
		Path:    accessTokenURI,
		ContentType: "application/json",
		Timeout: timeout,
	}

	var body struct {
		ComponentAppID        string `json:"component_appid"`
		ComponentAppSecret    string `json:"component_appsecret"`
		ComponentVerifyTicket string `json:"component_verify_token"`
	}{
		ComponentAppID:componentConfig.AppID,
		ComponentAppSecret:componentConfig.AppSecret,
		ComponentVerifyTicket:ticketStorage.Get()
	}

	b, err := json.Marshal(body)
	if err != nil{
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	token = new(ComponentAccessToken)
	err := req.DoPost(buf, token)
	if err != nil{
		return nil, err
	}

	return token, nil
}

// https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=xxx
func GetPreAuthCode(timeout int)(code *PreAuthCode, err error){
	req := wx.HttpClient{
		Path: preAuthCodeURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{
			"component_access_token": tokenStorage.Get(),
		},
		Timeout: timeout,
	}

	var body struct{
		ComponentAppID        string `json:"component_appid"`
	}{
		componentConfig.AppID,		
	}

	b, err := json.Marshal(body)
	if err != nil{
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	code = new(PreAuthCode)
	err := req.DoPost(buf, code)
	if err != nil{
		return nil, err
	}

	return code, nil
}

// https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=xxxx
func GetAuthorizationInfo(authorizationCode string, timeout int)(auth *Authorization, err error){
	req := wx.HttpClient{
		Path: authorizationInfoURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{
			"component_access_token": tokenStorage.Get(),
		},
		Timeout: timeout,
	}

	var body struct{
		ComponentAppID        string `json:"component_appid"`
		AuthorizationCode string `json:"authorization_code"`
	}{
		componentConfig.AppID,		
		authorizationCode,
	}

	b, err := json.Marshal(body)
	if err != nil{
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	auth = new(Authorization)
	err := req.DoPost(buf, auth)
	if err != nil{
		return nil, err
	}

	return code, nil
}

// https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=xxxxx
func RefreshAuthorizerToken(authorizerAppID, refreshToken string, timeout int)(info *AuthorizationInfo, err error){
		req := wx.HttpClient{
		Path: authorizeInfoURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{
			"component_access_token": tokenStorage.Get(),
		},
		Timeout: timeout,
	}

	var body struct{
		ComponentAppID        string `json:"component_appid"`
		AuthorizerAppID string `json:"authorizer_appid"`
		AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
	}{
		componentConfig.AppID,
		authorizerAppID,
		refreshToken,
	}

	b, err := json.Marshal(body)
	if err != nil{
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	info = new(AuthorizationInfo)
	err := req.DoPost(buf, info)
	if err != nil{
		return nil, err
	}

	return info, nil
}

// https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?component_access_token=xxxx
func GetAuthorizerInfo(authorizerAppID string, timeout int)(info *Authorizer, err error){
	req := wx.HttpClient{
		Path: authorizerInfoURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{
			"component_access_token": tokenStorage.Get(),
		},
		Timeout: timeout,
	}

	var body struct{
		ComponentAppID        string `json:"component_appid"`
		AuthorizerAppID string `json:"authorizer_appid"`
	}{
		componentConfig.AppID,
		authorizerAppID,
	}

	b, err := json.Marshal(body)
	if err != nil{
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	info = new(Authorizer)
	err := req.DoPost(buf, info)
	if err != nil{
		return nil, err
	}

	return info, nil
}

// https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_option?component_access_token=xxxx
func GetAuthorizerOption(authorizerAppID, optionName string, timeout int)(option *AuthorizerOption, err error){
	req := wx.HttpClient{
		Path: getAuthorizerOptionURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{
			"component_access_token": tokenStorage.Get(),
		},
		Timeout: timeout,
	}

	var body struct{
		ComponentAppID        string `json:"component_appid"`
		AuthorizerAppID string `json:"authorizer_appid"`
		OptionName string `json:"option_name"`
	}{
		componentConfig.AppID,
		authorizerAppID,
		optionName,
	}

	b, err := json.Marshal(body)
	if err != nil{
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	option = new(AuthorizerOption)
	err := req.DoPost(buf, option)
	if err != nil{
		return nil, err
	}

	return option, nil
}

// https://api.weixin.qq.com/cgi-bin/component/api_set_authorizer_option?component_access_token=xxxx
func SetAuthorizerOption(option *AuthorizerOption, timeout int)(error){
	req := wx.HttpClient{
		Path: setAuthorizerOptionURI,
		ContentType: "application/json",
		Parameters: []wx.QueryParameter{
			"component_access_token": tokenStorage.Get(),
		},
		Timeout: timeout,
	}

	var body struct{
		ComponentAppID        string `json:"component_appid"`
		AuthorizerAppID string `json:"authorizer_appid"`
		OptionName string `json:"option_name"`
		OptionValue string `json:"option_value"`
	}{
		componentConfig.AppID,
		option.AuthorizerAppID,
		option.OptionName,
		option.OptionValue,
	}

	b, err := json.Marshal(body)
	if err != nil{
		return nil, fmt.Errorf("json.Marshal: ", err)
	}
	buf := bytes.NewBuffer(b)

	err := req.DoPost(buf, interface{})
	if err != nil{
		return nil, err
	}

	return option, nil
}