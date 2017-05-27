package wx

const (
	GrantTypeRefresh   = "refresh_token"
	GrantTypeAuthorize = "authorization_code"
)

const (
	OAuthScopeBase     = "snsapi_base"
	OAUthScopeUserInfo = "snsapi_userinfo"
)

const (
	oAuthURI        = "https://open.weixin.qq.com/connect/oauth2/authorize"
	accessTokenURI  = "https://api.weixin.qq.com/sns/oauth2/access_token"
	refreshTokenURI = "https://api.weixin.qq.com/sns/oauth2/refresh_token"
	userInfoURI     = "https://api.weixin.qq.com/sns/userinfo"
)

const (
	LangCN = "zh_CN"
	LangTW = "zh_TW"
	LangEN = "en"
)
