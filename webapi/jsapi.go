package webapi

// jsapi_tikcet
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421141115
// https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1421823488&token=&lang=zh_CN

import (
	"time"

	"github.com/MenInBack/weshin/base"
	"github.com/MenInBack/weshin/component"
	"github.com/MenInBack/weshin/wx"
)

const (
	jsAPITicketURI = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
)

// GetJSAPITicket for js_api config
// https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=jsapi
func (s *WebAPI) GetJSAPITicket(timeout int) (*wx.APITicket, error) {
	var token string
	switch s.Mode {
	case wx.ModeMP:
		token = s.WechatMP.(base.MP).GetAccessToken()
	case wx.ModeComponent:
		// for component mode, token is authorizer access token, not component access token.
		token = s.WechatMP.(component.Component).GetAuthorizerToken(s.AppID)
	}

	req := wx.HttpClient{
		Path: jsAPITicketURI,
		Parameters: []wx.QueryParameter{
			{"access_token", token},
			{"type", wx.TicketTypeJSAPI},
		},
		Timeout: timeout,
	}

	ticket := new(wx.APITicket)
	err := req.Get(ticket)
	if err != nil {
		return nil, err
	}
	ticket.AppID = appID
	ticket.Typ = wx.TicketTypeJSAPI
	ticket.CreateAt = time.Now().Unix()

	go s.SetAPITicket(ticket)

	return ticket, nil
}
