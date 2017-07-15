package base

import (
	"github.com/MenInBack/weshin/wx"
)

// MPAccount for wechat official account
type MPAccount struct {
	AppID  string
	secret string
	Storage
}

// New MPAccount instance
func New(appID, secret string, storage Storage) *MPAccount {
	if storage == nil {
		storage = newDefaultStorage()
	}
	return &MPAccount{
		AppID:   appID,
		secret:  secret,
		Storage: storage,
	}
}

type MPAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type Storage interface {
	wx.TokenStorage
	wx.TicketStorage
}

// implements TokenStorage, without refreshing.
type defaultStorage struct {
	token       string
	jsAPITicket string
}

func newDefaultStorage() *defaultStorage {
	return new(defaultStorage)
}

func (s *defaultStorage) SetAccessToken(token string, expriresIn int64) {
	s.token = token
}

func (s *defaultStorage) GetAccessToken() string {
	return s.token
}

func (s *defaultStorage) SetAPITicket(ticket *wx.APITicket) {
	if ticket.Typ == wx.TicketTypeJSPAI {
		s.jsAPITicket = ticket.Ticket
	}
}

func (s *defaultStorage) GetAPITicket(typ string) string {
	if typ == wx.TicketTypeJSPAI {
		return s.jsAPITicket
	}
	return ""
}
