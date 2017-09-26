package wx

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// WechatError response from weixin
type WechatError struct {
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (e WechatError) Error() string {
	return fmt.Sprintf("wechat error: [%d] %s", e.ErrCode, e.ErrMsg)
}

func handleRespError(data []byte) error {
	err := new(WechatError)
	e := json.Unmarshal(data, &err)
	if e != nil {
		log.Print("unmarshal response data to WechatError error: ", err)
		return nil
	}
	if err.ErrCode != 0 {
		return err
	}
	return nil
}

type HttpError struct {
	State int
}

func (e HttpError) Error() string {
	return fmt.Sprintf("http error: [%d] %s", e.State, http.StatusText(e.State))
}

// NotifyError for async notifies
type NotifyError struct {
	Err error
}

func (e NotifyError) Error() string {
	return fmt.Sprintf("error when handling message notify %s", e.Err.Error())
}

// handler names for NotifyError
const (
	TicketHander     = "verify ticket"
	AuthorizerHander = "authorizer"
	UserAuthHandler  = "user authorization"
)

// ConfigError for invalid configuration
type ConfigError struct {
	InvalidConfig string
}

func (e ConfigError) Error() string {
	return fmt.Sprint("config error: invalid ", e.InvalidConfig)
}

// ParameterError for invalid parameter
type ParameterError struct {
	InvalidParameter string
}

func (e ParameterError) Error() string {
	return fmt.Sprint("parameter error: invalid ", e.InvalidParameter)
}

type WeshinError struct {
	Detail string
}

func (e WeshinError) Error() string {
	return fmt.Sprint("weshin error: ", e.Detail)
}
