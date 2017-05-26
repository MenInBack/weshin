package wx

import ()

var State = ""

var WXConfig struct {
	APPID            string
	Secret           string
	OAuthRedirectURI string
	State            string
}

type GrantTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
}

type OAuthResp struct {
	Code  string
	State string
}

type AuthorizeCodeReq struct {
	APPID     string
	Secret    string
	Code      string
	grantType string
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}

type RefreshTokenReq struct {
	APPID        string
	RefreshToken string
	grantType    string
}

type UserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        string   `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
}
