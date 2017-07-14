package base

// MPAccount for wechat official account
type MPAccount struct {
	appID  string
	secret string
	token  TokenStorage
}

// New MPAccount instance
func New(appID, secret string, tokenStorage TokenStorage) *MPAccount {
	if tokenStorage == nil {
		tokenStorage = newDefaultTokenStorage()
	}
	return &MPAccount{
		appID:  appID,
		secret: secret,
		token:  tokenStorage,
	}
}

// implements useroauth.MPServer interface
func (s MPAccount) GetAccessToken() string {
	return s.token.Get()
}

// TokenStorage holds official account's access token, and is responsible for token refreshing.
type TokenStorage interface {
	Set(token string, expriresIn int64)
	Get() string
}

type MPAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// implements TokenStorage, without refreshing.
type defaultStorage struct {
	token string
}

func newDefaultTokenStorage() *defaultStorage {
	return new(defaultStorage)
}

func (s *defaultStorage) Set(token string, expriresIn int64) {
	s.token = token
}

func (s *defaultStorage) Get() string {
	return s.token
}
