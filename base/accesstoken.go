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

// GrantAccessToken for wechat mp
// https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
func (s *MPAccount) GrantAccessToken(timeout int) (token *MPAccessToken, err error) {
	req := wx.HttpClient{
		Path:    accessTokenPath,
		Timeout: timeout,
		Parameters: []wx.QueryParameter{
			{"grant_type", wx.GrantTypeCredential},
			{"appid", s.appID},
			{"secret", s.secret},
		},
	}

	token = new(MPAccessToken)
	err = req.Get(token)
	if err != nil {
		log.Print("access token failed: ", err)
		return nil, err
	}

	return token, nil
}
