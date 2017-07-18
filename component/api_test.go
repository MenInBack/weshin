package component

import (
	"github.com/MenInBack/weshin/wx"
)

// sampleStorage implements Storage using local variables
type sampleStorage struct {
	verifyTicket      string
	jsAPITicket       string
	token             string
	tokenExpireAt     int64
	authorizationCode map[string]*AuthorizationCode
	authorizerToken   map[string]*AuthorizerToken
}

func newsampleStorage() *sampleStorage {
	return &sampleStorage{
		authorizationCode: make(map[string]*AuthorizationCode),
		authorizerToken:   make(map[string]*AuthorizerToken),
	}
}

func (s *sampleStorage) SetAccessToken(token string, expiresIn int64) {
	s.token = token
	s.tokenExpireAt = expiresIn
}

func (s *sampleStorage) GetAccessToken() string {
	return s.token
}

func (s *sampleStorage) SetAuthorizationCode(code *AuthorizationCode) {
	s.authorizationCode[code.AppID] = code
}

func (s *sampleStorage) GetAuthorizerToken(authorizerAppID string) string {
	return s.authorizerToken[authorizerAppID].AccessToken
}

func (s *sampleStorage) SetAuthorizerToken(token *AuthorizerToken) {
	s.authorizerToken[token.AppID] = token
}

func (s *sampleStorage) ClearAuthorizerToken(authorizerAppID string) {
	delete(s.authorizerToken, authorizerAppID)
}

func (s *sampleStorage) SetAPITicket(ticket *wx.APITicket) {
	switch ticket.Typ {
	case wx.TicketTypeJSAPI:
		s.jsAPITicket = ticket.Ticket
	case wx.TicketTypeVerify:
		s.verifyTicket = ticket.Ticket
	}
}

func (s *sampleStorage) GetAPITicket(typ string) string {
	switch typ {
	case wx.TicketTypeJSAPI:
		return s.jsAPITicket
	case wx.TicketTypeVerify:
		return s.verifyTicket
	}
	return ""
}
