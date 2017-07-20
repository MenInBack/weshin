package webapi

import (
	"log"
	"testing"

	"github.com/MenInBack/weshin/base"
	"github.com/MenInBack/weshin/wx"
)

const (
	appID       = "wx58b8557718c3b79e"
	secret      = "bfdde51505feb9b4d1de282c37ccd258"
	redirectURI = "REDIRECT_URI"
	state       = "STATE"
)

func TestJumpURL(t *testing.T) {
	mp := base.MP{
		AppID:   appID,
		Secret:  secret,
		Storage: new(sampleStorage),
	}
	api := WebAPI{
		Mode:     wx.ModeMP,
		WechatMP: mp,
	}
	jumpURI := api.JumpToAuth(wx.OAUthScopeUserInfo, redirectURI, state)
	t.Log(jumpURI)
}

// implements TokenStorage, without refreshing.
type sampleStorage struct {
	token       string
	jsAPITicket string
}

func newsampleStorage() *sampleStorage {
	return new(sampleStorage)
}

func (s *sampleStorage) SetAccessToken(token string, expriresIn int64) {
	log.Print("token setted: ", token)
	s.token = token
}

func (s *sampleStorage) GetAccessToken() string {
	return s.token
}

func (s *sampleStorage) SetAPITicket(ticket *wx.APITicket) {
	if ticket.Typ == wx.TicketTypeJSAPI {
		log.Println("jsapi ticket setted: ", ticket)
		s.jsAPITicket = ticket.Ticket
	}
}

func (s *sampleStorage) GetJSAPITicket() string {
	return s.jsAPITicket
}

// func (s *sampleStorage) GetAPITicket(typ string) string {
// 	if typ == wx.TicketTypeJSAPI {
// 		return s.jsAPITicket
// 	}
// 	return ""
// }
