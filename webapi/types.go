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
	wx.TokenStorage
	wx.TicketStorage
}

// New WebAPI
func New(appID, secret, componentID string, token wx.TokenStorage, ticket wx.TicketStorage) *WebAPI {
	o := &WebAPI{
		AppID:         appID,
		secret:        secret,
		ComponentID:   componentID,
		TokenStorage:  token,
		TicketStorage: ticket,
	}
	if componentID == "" {
		o.Mode = wx.ModeMP
	} else {
		o.Mode = wx.ModeComponent
	}

	return o
}

// UserAccessToken holds access token for user authorization
type UserAccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}
