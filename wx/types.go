package wx

type UserInfo struct {
	Subscribe     int32    `json:"subscribe"`
	OpenID        string   `json:"openid"`
	Nickname      string   `json:"nickname"`
	Sex           int32    `json:"sex"`
	Language      string   `json:"language"`
	Province      string   `json:"province"`
	City          string   `json:"city"`
	Country       string   `json:"country"`
	HeadImgURL    string   `json:"headimgurl"`
	Privilege     []string `json:"privilege"`
	SubscribeTime int64    `json:"subscribe_time"`
	Remark        string   `json:"remark"`
	GroupID       int32    `json:"groupid"`
	TagIDList     []int32  `json:"tagid_list"`
	UnionID       string   `json:"unionid"`
}

type APITicket struct {
	Typ       string `json:"-"`
	AppID     string `json:"-"`
	Ticket    string `json:"ticket,omitempty"`
	CreateAt  int64  `json:"create_at,omitempty"`
	ExpiresIn int64  `json:"expires_in,omitempty"`
}

type WechatMP interface {
	GetAppID() string
	GetSecret() string
	GetEncodingAESKey() string

	TokenStorage
	TicketStorage
}

// TokenStorage holds access token
type TokenStorage interface {
	// SetAccessToken is called inside GrantAccessToken to update access token,
	// access token refreshing should be arranged whenever setted.
	SetAccessToken(token string, expiresIn int64)
	// GetAccessToken availale currently.
	GetAccessToken() string
}

// TicketStorage holds verify ticket and js_api ticket
type TicketStorage interface {
	// SetAPITicket for verify ticket or js_api ticket.
	SetAPITicket(*APITicket)
	// GetAPITicket for verify ticket or js_api ticket.
	// GetAPITicket(typ string) string
}
