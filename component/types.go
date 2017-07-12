package component

import (
	"encoding/xml"
)

type Storage interface {
	Set(content string, createAt, expiresIn int64)
	Get() string
}


// handle component verify ticket storage
// a package-wise global variable is used as default,
// user should implement custom storage.
type ComponentVerifyTicket struct {
	XMLName               xml.Name `xml:"xml"`
	APPID                 string   `xml:"AppId"`
	CreateTime            int64    `xml:"CreateTime"`
	InfoType              string   `xml:"InfoType"`
	ComponentVerifyTicket string   `xml:"ComponentVerifyTicket"`
}

var ticketStorage *Storage

func UseCustomTicketStorage(s *Storage) {
	ticketStorage = s
}

type defaultTicketStorage ComponentVerifyTicket
ticketStorage = new(defaultTicketStorage)

func (s *defaultTicketStorage) Set(ticket string, createAt, expiresIn int64) error {
	s.ComponentVerifyTicket = ticket
	s.CreateTime = createAt
}

func (s *defaultTicketStorage) Get() string  {
	return s.ComponentVerifyTicket
}

// handle component access token storage
// a package-wise global variable is used as default,
// user should implement custom storage.
type ComponentAccessToken struct {
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int64  `json:"expires_in"`
}

var tokenStorage *Storage 

func UseCustomTokenStorage(s *Storage) {
	tokenStorage = s
}

type defaultTokenStorage ComponentAccessToken
tokenStorage = new(defaultTokenStorage)

func (s *defaultTokenStorage) Set(token string, createAt, expiresIn int64)  {
	s.ComponentAccessToken = token
	s.ExpiresIn = expiresIn
}

func (s *defaultTokenStorage) Get()string{
	return s.ComponentAccessToken
}

// pre auth code
type PreAuthCode struct {
	PreAuthCode string
	ExpiresIn   int64
}

// authorization info
type Authorization struct {
	AuthorizationInfo `json:"authorization_info"`
	FuncInfo          []FunctionInfo `json:"func_info"`
}

type AuthorizationInfo struct {
	AuthorizerAppID        string `json:"authorizer_appid"`
	AuthorizerAccessToken  string `json:"authorizer_access_token"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
	ExpiresIn              int64  `json:"expires_in"`
}

type FunctionInfo struct {
	FuncScopeCategory IDInfo `json:"funcscope_category"`
}

type IDInfo struct {
	ID int32 `json:"id"`
}

// authorizer info
type Authorizer struct {
	AuthorizerInfo struct {
		NickName        string         `json:"nick_name"`
		HeadImg         string         `json:"head_img"`
		ServiceTypeInfo IDInfo         `json:"service_type_info"`
		VerifyTypeInfo  IDInfo         `json:"verify_type_info"`
		UserName        string         `json:"user_name"`
		PrincipalName   string         `json:"principal_name"`
		BusinessInfo    map[string]int `json:"business_info"`
		Alias           string         `json:"alias,omitempty"`
		QRCodeURL       string         `json:"qrcode_url"`
		Signature       string         `json:"signature,omitempty"`
		MiniProgramInfo struct {
			Network     struct {
				RequestDomain   []string `json:"requestdomain"`
				WsRequestDomain []string `json:"wsrequestdomain"`
				UploadDomain    []string `json:"uploaddomain"`
				DownloadDomain  []string `json:"downloaddomain"`
			}    `json:"network"`
			Categories  []MiniProgramCategory `json:"categories"`
			VisitStatus int32                 `json:"visit_status"`
		} `json:"miniprograminfo,omitempty"`
	} `json:"authorizer_info"`

	AuthorizationInfo struct {
		AppID    string         `json:"appid"`
		FuncInfo []FunctionInfo `json:"func_info"`
	} `json:"authorization_info"`
}

type MiniProgramCategory struct {
	First  string
	Second string
}


// authorizer option
type AuthorizerOption struct {
	AuthorizerAppID string
	OptionName      string
	OptionValue     string
}

// authorization notify request body
type authorizationNotifyBody struct {
	XMLName    xml.Name `xml:"xml,cdata" json:"xmlName,omitempty"`
	AppID      string   `json:"appId,omitempty" xml:"appId,cdata"`
	CreateTime int64    `json:"createTime,omitempty" xml:"createTime,cdata"`
	InfoType   string   `json:"infoType,omitempty" xml:"infoType,cdata"`
	AuthorizationCode
}

type AuthorizationCode struct {
	AuthorizerAppID              string `json:"authorizerAppid,omitempty" xml:"authorizerAppid,cdata"`
	AuthorizationCode            string `json:"authorizationCode,omitempty" xml:"authorizationCode,cdata"`
	AuthorizationCodeExpiredTime int64  `json:"authorizationCodeExpiredTime,omitempty" xml:"authorizationCodeExpiredTime,cdata"`
}

// AuthorizationStorage holds authorization code storage.
// A package-wise global map is used as default.
// User should implement a custom storage, and request authorization info and authorizer info after setting new code.
type AuthorizationStorage interface{
	Set(*AuthorizationCode)
	Get(authorizerAppID string)*AuthorizationCode
	Clear(authorizerAppID string)
}

func UseCustomAuthorizationStorage(s *AuthorizationStorage){
	authorizationStorage = s
}

type defaultAuthorizationStorage map([string]*AuthorizationCode)
var authorizationStorage = make(map[string]*AuthorizationCode)

func (s *defaultAuthorizationStorage)Set(code *AuthorizationCode){
	s[code.AuthorizerAppID] = code
}

func (s *defaultAuthorizationStorage)Get(authorizerAppID string)*AuthorizationCode{
	return s[authorizerAppID]
}

func (s *defaultAuthorizationStorage)Clear(authorizerAppID string){
	delete(s, authorizerAppID)
}

