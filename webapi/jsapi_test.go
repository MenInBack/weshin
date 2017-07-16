package webapi

import (
	"github.com/MenInBack/weshin/base"
	"log"
	"testing"
)

func TestJSAPI(t *testing.T) {
	mp := base.New(appID, secret, nil)
	token, err := mp.GrantAccessToken(0)
	if err != nil {
		t.Error(err)
	}
	log.Printf("got token: %+v\n", token)

	api := New(appID, secret, "", nil, mp.Storage)
	ticket, err := api.GetJSAPITicket(mp.AppID, mp.GetAccessToken(), 0)
	if err != nil {
		t.Error(err)
	}
	log.Printf("got ticket: %+v\n", ticket)
}
