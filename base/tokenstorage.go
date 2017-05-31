package base

import (
	"log"
	"time"

	"github.com/MenInBack/weshin/wx"
)

type TokenStorage interface {
	Set(token *wx.MPAccessToken) error
	Get() (token *wx.MPAccessToken, err error)
	ArrangeRefresh()
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

func newTokenStorage() TokenStorage {
	return new(defaultStorage)
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

// will start an infinite refresh loop
func (s *defaultStorage) ArrangeRefresh() {
	time.AfterFunc(
		time.Duration(s.token.ExpiresIn*9/10)*time.Second,
		func() {
			t, err := GrantAccessToken(0)
			if err != nil {
				log.Print("refresh access token failed: ", err)
				return
			}

			err = s.Set(t)
			if err != nil {
				log.Print("set new access token failed: ", err)
				s.ArrangeRefresh()
			}
		})
}
