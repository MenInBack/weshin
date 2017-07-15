package base

import (
	"github.com/MenInBack/weshin/wx"
)

// MPAccount for wechat official account
type MPAccount struct {
	AppID  string
	secret string
	wx.TokenStorage
}

// New MPAccount instance
func New(appID, secret string, tokenStorage wx.TokenStorage) *MPAccount {
	if tokenStorage == nil {
		tokenStorage = newDefaultTokenStorage()
	}
	return &MPAccount{
		AppID:        appID,
		secret:       secret,
		TokenStorage: tokenStorage,
	}
}

type MPAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// implements TokenStorage, without refreshing.
type defaultStorage struct {
	token string
}

func newDefaultTokenStorage() *defaultStorage {
	return new(defaultStorage)
}

func (s *defaultStorage) SetAccessToken(token string, expriresIn int64) {
	s.token = token
}

func (s *defaultStorage) GetAccessToken() string {
	return s.token
}
