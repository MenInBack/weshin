package useroauth

type OAuth struct {
	appID       string
	secret      string
	accessToken AccessToken
	// userToken   UserAccessTokenStorage
}

func New(appid, secret string, accessToken AccessToken) *OAuth {
	// if userToken == nil {
	// 	userToken = newDefaultUserAccessToken()
	// }

	return &OAuth{
		appID:       appid,
		secret:      secret,
		accessToken: accessToken,
		// userToken:   userToken,
	}
}

// AccessToken of official account or thirdparty component
type AccessToken interface {
	GetAccessToken() string
}

// UserAccessToken
type UserAccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}

// type UserAccessTokenStorage interface {
// 	Set(token *UserAccessToken)
// 	Get(openID string) *UserAccessToken
// }

// type defaultUserAccessToken struct{}

// func newDefaultUserAccessToken() *defaultUserAccessToken             { return nil }
// func (t *defaultUserAccessToken) Set(token *UserAccessToken)         {}
// func (t *defaultUserAccessToken) Get(openID string) *UserAccessToken { return nil }

type UserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
}
