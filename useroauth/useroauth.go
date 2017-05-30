package useroauth

/**
 * https://mp.weixin.qq.com/wiki/ 微信网页开发/微信网页授权
 */

import (
	"bytes"
	"log"
	"net/url"
``
	"github.com/MenInBack/weshin/wx"
)

const (
	defaultSate = "mib_test"

	oAuthPath        = "https://open.weixin.qq.com/connect/oauth2/authorize"
	accessTokenPath  = "https://api.weixin.qq.com/sns/oauth2/access_token"
	refreshTokenPath = "https://api.weixin.qq.com/sns/oauth2/refresh_token"
	userinfoPath     = "https://api.weixin.qq.com/sns/userinfo"
)

var WXConfig wx.WXConfig

// JumpToAuth compose jump uri for user authorization
// https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect
func JumpToAuth(scope, redirectURI, state string) (jumpURL string, err error) {
	if len(WXConfig.APPID) <= 0 {
		return "", wx.ConfigError{InvalidConfig: "APPID"}
	}
	if scope == "" {
		scope = wx.OAuthScopeBase
	} else if scope != wx.OAUthScopeUserInfo && scope != wx.OAuthScopeBase {
		return "", wx.ParameterError{InvalidParameter: "scope"}
	}
	if redirectURI == "" {
		return "", wx.ParameterError{InvalidParameter: "redirectURI"}
	}
	if state == "" {
		state = defaultSate
	}

	u := bytes.NewBufferString(oAuthPath)
	u.WriteString("?appid=")
	u.WriteString(WXConfig.APPID)
	u.WriteString("&redirect_uri=")
	u.WriteString(url.QueryEscape(redirectURI))
	u.WriteString("&response_type=code")
	u.WriteString("&scope=")
	u.WriteString(scope)
	u.WriteString("&state=")
	u.WriteString(state)
	u.WriteString("#wechat_redirect")

	log.Print("jump uri for authorization: ", u.String())
	return u.String(), nil
}

// GrantAuthorizeToken grant access token for user authorization
// https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code
func GrantAuthorizeToken(code string, timeout int) (token *wx.UserAccessToken, err error) {
	log.Print("authorizing code: ", code)

	if len(WXConfig.APPID) <= 0 {
		return nil, wx.ConfigError{InvalidConfig: "appID"}
	}
	if len(WXConfig.Secret) <= 0 {
		return nil, wx.ConfigError{InvalidConfig: "secret"}
	}
	if len(code) <= 0 {
		return nil, wx.ParameterError{InvalidParameter: "code"}
	}

	req := wx.HttpClient{
		Path:    accessTokenPath,
		Timeout: timeout,
		Parameters: []wx.QueryParameter{
			{"appid", WXConfig.APPID},
			{"secret", WXConfig.Secret},
			{"code", code},
			{"grant_type", wx.GrantTypeAuthorize},
		},
	}

	token = new(wx.UserAccessToken)
	err = req.Get(token)
	if err != nil {
		log.Print("authorize code failed: ", err)
		return nil, err
	}
	return token, nil
}

// RefreshAuthorizeToken refresh access token for user authorization
// https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=APPID&grant_type=refresh_token&refresh_token=REFRESH_TOKEN
func RefreshAuthorizeToken(refreshToken string, timeout int) (token *wx.UserAccessToken, err error) {
	if len(WXConfig.APPID) <= 0 {
		return nil, wx.ConfigError{InvalidConfig: "appID"}
	}
	if len(refreshToken) <= 0 {
		return nil, wx.ParameterError{InvalidParameter: "refresh token"}
	}

	req := wx.HttpClient{
		Path:    refreshTokenPath,
		Timeout: timeout,
		Parameters: []wx.QueryParameter{
			{"appid", WXConfig.APPID},
			{"grant_type", wx.GrantTypeRefresh},
			{"refresh_token", refreshToken},
		},
	}

	token = new(wx.UserAccessToken)
	err = req.Get(token)
	if err != nil {
		log.Print("refresh token failed: ", err)
		return nil, err
	}

	return token, err
}

// GetUserInfo get authorized user info
// https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
func GetUserInfo(token, openID, lang string, timeout int) (info *wx.UserInfo, err error) {
	if len(token) <= 0 {
		return nil, wx.ParameterError{InvalidParameter: "access token"}
	}
	if len(openID) <= 0 {
		return nil, wx.ParameterError{InvalidParameter: "openID"}
	}
	if lang == "" {
		lang = wx.LangCN
	} else if lang != wx.LangCN && lang != wx.LangTW && lang != wx.LangEN {
		return nil, wx.ParameterError{InvalidParameter: "lang"}
	}

	req := wx.HttpClient{
		Path:    userinfoPath,
		Timeout: timeout,
		Parameters: []wx.QueryParameter{
			{"access_token", token},
			{"openid", openID},
			{"lang", lang},
		},
	}

	info = new(wx.UserInfo)
	err = req.Get(info)
	if err != nil {
		log.Print("query user info failed: ", err)
		return nil, err
	}

	return info, err
}
