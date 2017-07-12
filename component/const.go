package component

const (
	accessTokenURI         = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
	preAuthCodeURI         = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode"
	authorizationInfoURI   = "https://api.weixin.qq.com/cgi-bin/component/api_query_auth"
	refreshAuthURI         = "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token"
	authorizerInfoURI      = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info"
	getAuthorizerOptionURI = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_option"
	setAuthorizerOptionURI = "https://api.weixin.qq.com/cgi-bin/component/api_set_authorizer_option"
)

const (
	NotifyTypeUnauthorized     = "unauthorized"
	NotifyTypeAuthorized       = "authorized"
	NotifyTypeUpdateAuthorized = "updateauthorized"
)
