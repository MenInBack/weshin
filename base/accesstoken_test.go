package base

import (
	"github.com/MenInBack/weshin/wx"
	"log"
	"testing"
)

func init() {
	WXConfig.AppID = "YourAppID"
	WXConfig.Secret = "YourSecret"
}

func TestAccessToken(t *testing.T) {
	token, err := GrantAccessToken(0)
	if err != nil {
		t.Error("grant access token failed: ", err)
	}
	log.Print("got access token: ", token)
}

func TestTokenStorage(t *testing.T) {
	grantedToken := "GrantedToken"

	err := tokenStorage.Set(&wx.MPAccessToken{
		AccessToken: grantedToken,
		ExpiresIn:   7200,
	})
	if err != nil {
		t.Error("set failed: ", err)
	}

	tk, err := tokenStorage.Get()
	if err != nil {
		t.Error("get failed: ", err)
	}

	log.Print("got token: ", tk)

	if tk.AccessToken != grantedToken {
		t.Error("token mismatch")
	}
	if tk.ExpiresIn != 7200 {
		t.Error("exporesIn mismatch")
	}
}

func TestGetUserInfo(t *testing.T) {
	grantedToken := "GrantedToken"

	err := tokenStorage.Set(&wx.MPAccessToken{
		AccessToken: grantedToken,
		ExpiresIn:   7200,
	})
	if err != nil {
		t.Error("set failed: ", err)
	}

	someOpenID := "someOpenID"
	info, err := GetUserInfo(someOpenID, wx.LangCN, 0)
	if err != nil {
		t.Error("get userinfo failed: ", err)
	}
	log.Printf("got user info: %+v", info)
}
