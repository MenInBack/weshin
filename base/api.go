package base

/**
 * https://mp.weixin.qq.com/wiki/ 开始开发/获取access_token
 */

import (
	"github.com/MenInBack/weshin/wx"
)

const (
	accessTokenPath = "https://api.weixin.qq.com/cgi-bin/token"
	userinfoPath    = "https://api.weixin.qq.com/cgi-bin/user/info"
)

// GrantAccessToken for wechat mp
// https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
func (mp *MP) GrantAccessToken(timeout int) (token *MPAccessToken, err error) {
	req := wx.HttpClient{
		Path:    accessTokenPath,
		Timeout: timeout,
		Parameters: []wx.QueryParameter{
			{"grant_type", wx.GrantTypeCredential},
			{"appid", mp.AppID},
			{"secret", mp.Secret},
		},
	}

	token = new(MPAccessToken)
	err = req.Get(token)
	if err != nil {
		return nil, err
	}

	mp.SetAccessToken(token.AccessToken, token.ExpiresIn)
	return token, nil
}

// GetUserInfo with known openID
// https://mp.weixin.qq.com/wiki/ 用户管理/获取用户基本信息(UnionID机制)
// https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
func (mp *MP) GetUserInfo(openID, lang string, timeout int) (userinfo *wx.UserInfo, err error) {
	if len(openID) <= 0 {
		return nil, wx.ParameterError{InvalidParameter: "openID"}
	}
	if len(lang) <= 0 {
		lang = wx.LangCN
	} else if lang != wx.LangCN && lang != wx.LangEN && lang != wx.LangTW {
		return nil, wx.ParameterError{InvalidParameter: "lang"}
	}

	req := wx.HttpClient{
		Path:    userinfoPath,
		Timeout: timeout,
		Parameters: []wx.QueryParameter{
			{"access_token", mp.GetAccessToken()},
			{"openid", openID},
			{"lang", lang},
		},
	}

	userinfo = new(wx.UserInfo)
	err = req.Get(userinfo)
	if err != nil {
		return nil, err
	}

	return userinfo, nil
}
