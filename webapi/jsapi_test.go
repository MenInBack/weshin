package webapi

import (
	"github.com/MenInBack/weshin/base"
	"testing"
)

const (
	appID  = ""
	secret = ""
)

func TestJSAPI(t *testing.T) {
	mp := base.New(appID, secret, nil)
	_, err := mp.GrantAccessToken(0)
	if err != nil {
		t.Error(err)
	}

	api := New(appID, secret, mp)
	ticket, err := api.GetJSAPITicket(0)
	if err != nil(
		t.Error(err)
	)
	

}
