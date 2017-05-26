package wx

import (
	"testing"
)

func init() {
	WXConfig.APPID = "APPID"
	WXConfig.Secret = "SECRET"
	WXConfig.State = "STATE"
}

func TestAuthorizeCode(t *testing.T) {
	code := "CODE"

	token, err := AuthorizeCode(code, 0)
	if err != nil {
		t.Error(err)
	}

	t.Log(token)
}
