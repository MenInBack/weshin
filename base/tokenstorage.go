package base

import (
	"github.com/MenInBack/weshin/wx"
)

type TokenStorage interface {
	Set(token *wx.MPAccessToken) error
	Get() (token *wx.MPAccessToken, err error)
}

var tokenStorage = newTokenStorage()

// UseCustomStorage sets custom storage for access token,
// should be called on begging of custom code.
func UseCustomStorage(s TokenStorage) {
	tokenStorage = s
}

type defaultStorage struct {
	token *wx.MPAccessToken
}

func (s *defaultStorage) Set(token *wx.MPAccessToken) error {
	s.token = token
	return nil
}

func (s *defaultStorage) Get() (token *wx.MPAccessToken, err error) {
	if s.token == nil {
		return nil, wx.ConfigError{InvalidConfig: "accessToken"}
	}
	return s.token, nil
}

func newTokenStorage() TokenStorage {
	return new(defaultStorage)
}
