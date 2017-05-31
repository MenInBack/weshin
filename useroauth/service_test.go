package useroauth

import (
	"testing"

	"github.com/MenInBack/weshin/wx"
)

/**
 * https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect
 */

func TestJumpURL(t *testing.T) {

	jumpURI, err := JumpToAuth(wx.OAUthScopeUserInfo, "redirectURI", "")
	if err != nil {
		t.Error(err)
	}

	t.Log(jumpURI)
}

func TestGrantAuthorizeToken(t *testing.T) {}
