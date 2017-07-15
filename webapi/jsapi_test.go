package webapi

import (
	"github.com/MenInBack/weshin/base"
	"testing"
)

func TestJSAPI(t *testing.T) {
	mp := base.New(appID, secret, nil)
	token, err := mp.GrantAccessToken(0)
	if err != nil {
		t.Error(err)
	}
	t.Logf("got token: %+v", token)

	api := New(appID, secret, "", nil, mp.Storage)
	ticket, err := api.GetJSAPITicket(mp.GetAccessToken(), 0)
	if err != nil {
		t.Error(err)
	}
	t.Logf("got ticket: %+v", ticket)
}
