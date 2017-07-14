package useroauth

import (
	"github.com/MenInBack/weshin/base"
	"github.com/MenInBack/weshin/component"
	"github.com/MenInBack/weshin/wx"
)

// OAuth for user authorization
type OAuth struct {
	Mode        int32
	AppID       string
	ComponentID string
	secret      string
	server      MPServer
}

// New OAuth
func New(appID, secret string, server MPServer) *OAuth {
	o := &OAuth{
		AppID:  appID,
		secret: secret,
	}
	switch server.(type) {
	case base.MPAccount:
		o.Mode = wx.ModeMP
	case component.Component:
		o.Mode = wx.ModeComponent
		o.ComponentID = server.(component.Component).AppID
	}
	return o
}

// base.MPAccount or component.Component
type MPServer interface {
	GetAccessToken() string
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
