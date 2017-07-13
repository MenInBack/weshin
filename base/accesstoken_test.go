package base

import (
	"log"
	"testing"
)

// 可使用公众平台接口测试号
// https://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo?action=showinfo&t=sandbox/index
const (
	appID  = ""
	secret = ""
	openID = ""
)

func TestAccessToken(t *testing.T) {
	mp := New(appID, secret, nil)
	token, err := mp.GrantAccessToken(0)
	if err != nil {
		t.Error("grant access token failed: ", err)
	}
	log.Print("got access token: ", token)
}

func TestGetUserInfo(t *testing.T) {
	accessToken := ""
	mp := New(appID, secret, nil)
	mp.token.Set(accessToken, 0)
	info, err := mp.GetUserInfo(openID, "", 0)
	if err != nil {

		t.Error("get userinfo failed: ", err)
	}
	log.Printf("got user info: %+v", info)
}
