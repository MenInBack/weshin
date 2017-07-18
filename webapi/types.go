package webapi

import (
	"github.com/MenInBack/weshin/wx"
)

type WebAPI struct {
	Mode        int32
	ComponentID string
	wx.WechatMP
}

// UserAccessToken holds access token for user authorization
type UserAccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}
