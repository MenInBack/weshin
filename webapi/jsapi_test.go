package webapi

import (
	"log"
	"testing"

	"github.com/MenInBack/weshin/base"
	"github.com/MenInBack/weshin/wx"
)

func TestJSAPI(t *testing.T) {
	mp := base.MP{
		AppID:   appID,
		Secret:  secret,
		Storage: new(sampleStorage),
	}
	token, err := mp.GrantAccessToken(0)
	if err != nil {
		t.Error(err)
	}
	log.Printf("got token: %+v\n", token)

	api := WebAPI{
		Mode:     wx.ModeMP,
		WechatMP: mp,
	}
	ticket, err := api.GetJSAPITicket(0)
	if err != nil {
		t.Error(err)
	}
	log.Printf("got ticket: %+v\n", ticket)
}
