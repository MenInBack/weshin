package base

import (
	"github.com/MenInBack/weshin/wx"
)

type TokenStorage interface {
	Set(token wx.MPAccessToken) error
	Get() (token *wx.MPAccessToken, err error)
}

type defaultStorage struct {
	token *wx.MPAccessToken
}

func (s *defaultStorage) Set(token *wx.MPAccessToken) error {
	s.token = token
}

func (s *defaultStorage) Get() (token *wx.MPAccessToken, err error) {
	if s.token == nil {
		return nil, wx.ConfigError{InvalidConfig: "accessToken"}
	}
	return s.token, nil
}

var tokenStorage TokenStorage

func init() {
	tokenStorage = defaultStorage{}
}

func UseStorage(s TokenStorage) {
	tokenStorage = s
}
