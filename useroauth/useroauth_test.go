package useroauth

import (
	"testing"

	"github.com/MenInBack/weshin/wx"
)

const (
	appID       = "APPID"
	secret      = ""
	redirectURI = "REDIRECT_URI"
	state       = "STATE"
)

func TestJumpURL(t *testing.T) {
	uri := `https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect`

	oAuth := New(appID, secret)
	jumpURI, err := oAuth.JumpToAuth(wx.OAUthScopeUserInfo, redirectURI, state)
	if err != nil {
		t.Error(err)
	}
	if jumpURI != uri {
		t.Error("incorrect uri")
	}
	t.Log(jumpURI)
}
