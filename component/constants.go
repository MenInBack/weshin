package component

const (
	accessTokenURI         = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
	preAuthCodeURI         = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode"
	authorizeURI           = "https://mp.weixin.qq.com/cgi-bin/componentloginpage"
	authorizationInfoURI   = "https://api.weixin.qq.com/cgi-bin/component/api_query_auth"
	authorizerTokenURI     = "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token"
	authorizerInfoURI      = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info"
	getAuthorizerOptionURI = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_option"
	setAuthorizerOptionURI = "https://api.weixin.qq.com/cgi-bin/component/api_set_authorizer_option"
)

const (
	NotifyTypeVerifyTicket     = "component_verify_ticket"
	NotifyTypeUnauthorized     = "unauthorized"
	NotifyTypeAuthorized       = "authorized"
	NotifyTypeUpdateAuthorized = "updateauthorized"
)
