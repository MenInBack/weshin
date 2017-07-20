package base

import (
	"log"
	"testing"

	"github.com/MenInBack/weshin/wx"
)

// 可使用公众平台接口测试号
// https://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo?action=showinfo&t=sandbox/index
const (
	appID  = "wx58b8557718c3b79e"
	secret = "bfdde51505feb9b4d1de282c37ccd258"
	openID = "oicAdwjVCJlTfem2YzNULtPm1-2g"
)

func TestAccessToken(t *testing.T) {
	mp := MP{
		AppID:   appID,
		Secret:  secret,
		Storage: new(sampleStorage),
	}
	token, err := mp.GrantAccessToken(0)
	if err != nil {
		t.Error("grant access token failed: ", err)
	}
	log.Printf("got access token: %+v", token)
}

func TestGetUserInfo(t *testing.T) {
	mp := MP{
		AppID:   appID,
		Secret:  secret,
		Storage: new(sampleStorage),
	}
	token, err := mp.GrantAccessToken(0)
	if err != nil {
		t.Error("grant access token failed: ", err)
	}
	log.Print("got access token: ", token)

	info, err := mp.GetUserInfo(openID, "", 0)
	if err != nil {
		t.Error("get userinfo failed: ", err)
	}
	log.Printf("got user info: %+v", info)
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
	s.token = token
}

func (s *sampleStorage) GetAccessToken() string {
	return s.token
}

func (s *sampleStorage) SetAPITicket(ticket *wx.APITicket) {
	if ticket.Typ == wx.TicketTypeJSAPI {
		s.jsAPITicket = ticket.Ticket
	}
}

func (s *sampleStorage) GetJSAPITicket() *wx.APITicket {
	return &wx.APITicket{
		Typ:    wx.TicketTypeJSAPI,
		Ticket: s.jsAPITicket,
	}
}

// func (s *sampleStorage) GetAPITicket(typ string) string {
// 	if typ == wx.TicketTypeJSAPI {
// 		return s.jsAPITicket
// 	}
// 	return ""
// }
