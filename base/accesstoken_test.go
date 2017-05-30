package base

import (
	"github.com/MenInBack/weshin/wx"
	"log"
	"testing"
)

func init() {
	WXConfig.APPID = "YourAPPID"
	WXConfig.Secret = "YourSecret"

	grantedToken := "GrantedToken"
	someOpenID := "SomeOpenID"
}

func TestAccessToken(t *testing.T) {
	token, err := GrantAccessToken(0)
	if err != nil {
		t.Error("grant access token failed: ", err)
	}
	log.Print("got access token: ", token)
}

// CcQ13AaEriT-ZgK4idPiVThNILT9IspW0uc4S4C331SK3PNRvDVOMc35cSvOh-5YnVM8Plk4Tj1hC8Fi0NxIo9H5AT5w2Dhoa3Lc-7fIx69or1aDiIo87gRcg0nwauDDWFMfAAACLS
func TestTokenStorage(t *testing.T) {
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
	err := tokenStorage.Set(&wx.MPAccessToken{
		AccessToken: grantedToken,
		ExpiresIn:   7200,
	})
	if err != nil {
		t.Error("set failed: ", err)
	}

	info, err := GetUserInfo(someOpenID, wx.LangCN, 0)
	if err != nil {
		t.Error("get userinfo failed: ", err)
	}
	log.Printf("got user info: %+v", info)
}
