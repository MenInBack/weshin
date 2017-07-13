package base

/**
 * https://mp.weixin.qq.com/wiki/ 用户管理/获取用户基本信息(UnionID机制)
 */

import (
	"log"

	"github.com/MenInBack/weshin/wx"
)

const (
	userinfoPath = "https://api.weixin.qq.com/cgi-bin/user/info"
)

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

// GetUserInfo with known openID
// https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
func (s *MPService) GetUserInfo(openID, lang string, timeout int) (userinfo *UserInfo, err error) {
	if len(openID) <= 0 {
		return nil, wx.ParameterError{InvalidParameter: "openID"}
	}
	if len(lang) <= 0 {
		lang = wx.LangCN
	} else if lang != wx.LangCN && lang != wx.LangEN && lang != wx.LangTW {
		return nil, wx.ParameterError{InvalidParameter: "lang"}
	}

	req := wx.HttpClient{
		Path:    userinfoPath,
		Timeout: timeout,
		Parameters: []wx.QueryParameter{
			{"access_token", s.token.Get()},
			{"openid", openID},
			{"lang", lang},
		},
	}

	userinfo = new(UserInfo)
	err = req.Get(userinfo)
	if err != nil {
		log.Print("get userinfo failed: ", err)
		return nil, err
	}

	return userinfo, nil
}
