package wx

// official account mode
const (
	ModeComponent = 1
	ModeMP        = 2
)

// ticket type
const (
	TicketTypeJSPAI  = "jsapi_ticket"
	TicketTypeVerify = "component_verify_ticket"
)

// grant type
const (
	GrantTypeRefresh    = "refresh_token"
	GrantTypeAuthorize  = "authorization_code"
	GrantTypeCredential = "client_credential"
)

// oauth scope
const (
	OAuthScopeBase     = "snsapi_base"
	OAUthScopeUserInfo = "snsapi_userinfo"
)

// language option
const (
	LangCN = "zh_CN"
	LangTW = "zh_TW"
	LangEN = "en"
)
