package webapi

import (
	"testing"

	"github.com/MenInBack/weshin/wx"
)

const (
	appID       = "wx58b8557718c3b79e"
	secret      = "bfdde51505feb9b4d1de282c37ccd258"
	redirectURI = "REDIRECT_URI"
	state       = "STATE"
)

func TestJumpURL(t *testing.T) {
	webAPI := New(appID, secret, "", nil, nil)
	jumpURI, err := webAPI.JumpToAuth(wx.OAUthScopeUserInfo, redirectURI, state)
	if err != nil {
		t.Error(err)
	}
	t.Log(jumpURI)
}
