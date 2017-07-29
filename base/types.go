package base

import (
	"github.com/MenInBack/weshin/wx"
)

// MP for wechat official account
type MP struct {
	AppID          string
	Secret         string
	EncodingAESKey string
	Token          string
	Storage
}

// implements wx.WechatMP
func (mp MP) GetAppID() string {
	return mp.AppID
}
func (mp MP) GetSecret() string {
	return mp.Secret
}
func (mp MP) GetEncodingAESKey() string {
	return mp.EncodingAESKey
}

type MPAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Storage holds access token and JS-SDK ticket
type Storage interface {
	wx.AccessTokenStorage
	wx.JSTicketStorage
}
