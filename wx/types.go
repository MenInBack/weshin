package wx

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

type APITicket struct {
	Typ       string `json:"-,omitempty"`
	Ticket    string `json:"ticket,omitempty"`
	CreateAt  int64  `json:"create_at,omitempty"`
	ExpiresIn int64  `json:"expires_in,omitempty"`
}

// TokenStorage holds access token
type TokenStorage interface {
	// SetAccessToken is called inside GrantAccessToken to update access token,
	// access token refreshing should be arranged whenever setted.
	SetAccessToken(token string, expiresIn int64)
	// GetAccessToken availale currently.
	GetAccessToken() string
}
