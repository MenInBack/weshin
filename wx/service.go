package wx

import (
	"bytes"
	"log"
)

func JumpToAuth(scope string) (jumpURL string, err error) {
	if len(WXConfig.APPID) <= 0 {
		return "", ConfigError{InvalidConfig: "APPID"}
	}
	if len(WXConfig.OAuthRedirectURI) <= 0 {
		return "", ConfigError{InvalidConfig: "OAuthRedirectURI"}
	}
	if len(WXConfig.State) <= 0 {
		return "", ConfigError{InvalidConfig: "state"}
	}
	if scope != OAUthScopeUserInfo || scope != OAuthScopeBase {
		return "", ParameterError{InvalidParameter: "scope"}
	}

	u := bytes.NewBufferString(oAuthURI)
	u.WriteString("?appid=")
	u.WriteString(WXConfig.APPID)
	u.WriteString("&redirect_uri=")
	u.WriteString(WXConfig.OAuthRedirectURI)
	u.WriteString("&response_type=code")
	u.WriteString("&scope=")
	u.WriteString(scope)
	u.WriteString("&state=")
	u.WriteString(WXConfig.State)
	u.WriteString("#wechat_redirect")

	log.Print("jump uri for authorization: ", u.String())
	return u.String(), nil
}

func AuthorizeCode(code string, timeout int) (token *AccessToken, err error) {
	log.Print("authorizing code: ", code)

	if len(WXConfig.APPID) <= 0 {
		return nil, ConfigError{InvalidConfig: "appID"}
	}
	if len(WXConfig.Secret) <= 0 {
		return nil, ConfigError{InvalidConfig: "secret"}
	}
	if len(code) <= 0 {
		return nil, ParameterError{InvalidParameter: "code"}
	}

	req := httpClient{
		Path:    authorizeCodeReqURI,
		Timeout: timeout,
		Parameters: []queryParameter{
			{"appid", WXConfig.APPID},
			{"secret", WXConfig.Secret},
			{"code", code},
			{"grant_type", GrantTypeAuthorize},
		},
	}

	token = new(AccessToken)
	err = req.Get(token)
	if err != nil {
		log.Print("authorize code failed: ", err)
		return nil, err
	}
	return token, nil
}

func RefreshToken(refreshToken string, timeout int) (token *AccessToken, err error) {
	if len(WXConfig.APPID) <= 0 {
		return nil, ConfigError{InvalidConfig: "appID"}
	}
	if len(refreshToken) <= 0 {
		return nil, ParameterError{InvalidParameter: "refreshToken"}
	}

	req := httpClient{
		Path:    refreshTokenURI,
		Timeout: timeout,
		Parameters: []queryParameter{
			{"appid", WXConfig.APPID},
			{"grant_type", GrantTypeRefresh},
			{"refresh_token", refreshToken},
		},
	}

	token = new(AccessToken)
	err = req.Get(token)
	if err != nil {
		log.Print("refresh token failed: ", err)
		return nil, err
	}

	return token, err
}