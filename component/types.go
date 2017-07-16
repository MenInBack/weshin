package component

import (
	"encoding/xml"
	"github.com/MenInBack/weshin/wx"
)

// Component services in place of official accounts
type Component struct {
	AppID          string
	appSecret      string
	encodingAESKey string
	addresses      *NotifyConfig
	Storage
}

// New Component instance
func New(appid, appsecret, encodingAESKey string, storage Storage, address *NotifyConfig) *Component {
	if storage == nil {
		storage = newDefaultStorage()
	}
	return &Component{
		AppID:          appid,
		appSecret:      appsecret,
		encodingAESKey: encodingAESKey,
		addresses:      address,
		Storage:        storage,
	}
}

// Storage holds component ticket, access token, and authorizer codes,
// and should be responsible for token refreshing.
type Storage interface {
	// holds token
	wx.TokenStorage

	// holds ticket
	wx.TicketStorage

	// SetAuthorizerToken when authorized.
	SetAuthorizerToken(token *AuthorizerToken)
	// GetAuthorizerToken for querying authorizer info if authorized,
	// should refresh authorizer token if expired.
	GetAuthorizerToken(authorizerAppID string) string
	// ClearAuthorizertoken when authorization cancelled.
	ClearAuthorizertoken(authorizerAppID string)

	SetAuthorizationCode(code *AuthorizationCode)
}

// AuthorizationCode holds authorizer code
type AuthorizationCode struct {
	AppID       string `json:"authorizerAppid,omitempty" xml:"authorizerAppid,cdata"`
	Code        string `json:"authorizationCode,omitempty" xml:"authorizationCode,cdata"`
	ExpiredTime int64  `json:"authorizationCodeExpiredTime,omitempty" xml:"authorizationCodeExpiredTime,cdata"`
}

// NotifyConfig configures notify addresses for wechat message
type NotifyConfig struct {
	Address           string
	VerifyTicketPath  string
	AuthorizationPath string
}

// {
// "component_access_token":"61W3mEpU66027wgNZ_MhGHNQDHnFATkDa9-2llqrMBjUwxRSNPbVsMmyD-yq8wZETSoE5NQgecigDrSHkPtIYA",
// "expires_in":7200
// }
type ComponentAccessToken struct {
	Token     string `json:"component_access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

// pre auth code for authorization
type PreAuthCode struct {
	PreAuthCode string
	ExpiresIn   int64
}

// authorization token info
type AuthorizationTokenInfo struct {
	AuthorizationToken AuthorizerToken `json:"authorization_info"`
	FuncInfo           []FunctionInfo  `json:"func_info"`
}

type AuthorizerToken struct {
	AppID        string `json:"authorizer_appid"`
	AccessToken  string `json:"authorizer_access_token"`
	RefreshToken string `json:"authorizer_refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
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
			Network struct {
				RequestDomain   []string `json:"requestdomain"`
				WsRequestDomain []string `json:"wsrequestdomain"`
				UploadDomain    []string `json:"uploaddomain"`
				DownloadDomain  []string `json:"downloaddomain"`
			} `json:"network"`
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

// defaultStorage implements Storage using local variables
type defaultStorage struct {
	verifyTicket      string
	jsAPITicket       string
	token             string
	tokenExpireAt     int64
	authorizationCode map[string]*AuthorizationCode
	authorizerToken   map[string]*AuthorizerToken
}

func newDefaultStorage() *defaultStorage {
	return &defaultStorage{
		authorizationCode: make(map[string]*AuthorizationCode),
		authorizerToken:   make(map[string]*AuthorizerToken),
	}
}

func (s *defaultStorage) SetAccessToken(token string, expiresIn int64) {
	s.token = token
	s.tokenExpireAt = expiresIn
}

func (s *defaultStorage) GetAccessToken() string {
	return s.token
}

func (s *defaultStorage) SetAuthorizationCode(code *AuthorizationCode) {
	s.authorizationCode[code.AppID] = code
}

func (s *defaultStorage) GetAuthorizerToken(authorizerAppID string) string {
	return s.authorizerToken[authorizerAppID].AccessToken
}

func (s *defaultStorage) SetAuthorizerToken(token *AuthorizerToken) {
	s.authorizerToken[token.AppID] = token
}

func (s *defaultStorage) ClearAuthorizertoken(authorizerAppID string) {
	delete(s.authorizerToken, authorizerAppID)
}

func (s *defaultStorage) SetAPITicket(ticket *wx.APITicket) {
	switch ticket.Typ {
	case wx.TicketTypeJSPAI:
		s.jsAPITicket = ticket.Ticket
	case wx.TicketTypeVerify:
		s.verifyTicket = ticket.Ticket
	}
}

func (s *defaultStorage) GetAPITicket(typ string) string {
	switch typ {
	case wx.TicketTypeJSPAI:
		return s.jsAPITicket
	case wx.TicketTypeVerify:
		return s.verifyTicket
	}
	return ""
}
