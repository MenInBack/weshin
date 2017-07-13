package base

// MPService for wechat official account
type MPService struct {
	appID  string
	secret string
	token  TokenStorage
}

// New MPService instance
func New(appID, secret string, tokenStorage TokenStorage) *MPService {
	if tokenStorage == nil {
		tokenStorage = newDefaultTokenStorage()
	}
	return &MPService{
		appID: appID,
		token: tokenStorage,
	}
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
