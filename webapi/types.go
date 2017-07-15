package webapi

import (
	"github.com/MenInBack/weshin/wx"
)

// WebAPI for user authorization
type WebAPI struct {
	Mode        int32
	AppID       string
	ComponentID string
	secret      string
	token       wx.TokenStorage
}

// New WebAPI
func New(appID, secret, componentID string) *WebAPI {
	o := &WebAPI{
		AppID:       appID,
		secret:      secret,
		ComponentID: componentID,
	}
	if componentID == "" {
		o.Mode = wx.ModeMP
	} else {
		o.Mode = wx.ModeComponent
	}

	return o
}

// base.MPAccount or component.Component
type MPServer interface {
	// GetAccessToken() string
	SetAPITicket(*wx.APITicket)
}

// UserAccessToken holds access token for user authorization
type UserAccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}

// UserInfo for authorized users
type UserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
}
