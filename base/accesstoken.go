package base

/**
 * https://mp.weixin.qq.com/wiki/ 开始开发/获取access_token
 */

import (
	"log"

	"github.com/MenInBack/weshin/wx"
)

const (
	accessTokenPath = "https://api.weixin.qq.com/cgi-bin/token"
)

var WXConfig wx.Config

// GrantAccessToken for wechat mp
// https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
func GrantAccessToken(timeout int) (token *wx.MPAccessToken, err error) {
	if len(WXConfig.AppID) <= 0 {
		return nil, wx.ConfigError{InvalidConfig: "appID"}
	}
	if len(WXConfig.Secret) <= 0 {
		return nil, wx.ConfigError{InvalidConfig: "secret"}
	}

	req := wx.HttpClient{
		Path:    accessTokenPath,
		Timeout: timeout,
		Parameters: []wx.QueryParameter{
			{"grant_type", wx.GrantTypeCredential},
			{"appid", WXConfig.AppID},
			{"secret", WXConfig.Secret},
		},
	}

	token = new(wx.MPAccessToken)
	err = req.Get(token)
	if err != nil {
		log.Print("access token failed: ", err)
		return nil, err
	}

	return token, nil
}
