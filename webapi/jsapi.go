package webapi

// jsapi_tikcet
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421141115
// https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1421823488&token=&lang=zh_CN

import (
	"github.com/MenInBack/weshin/wx"
)

const (
	jsAPITicketURI = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
)

// https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=jsapi
func (s *WebAPI) GetJSAPITicket(token string, timeout int) (*wx.APITicket, error) {
	req := wx.HttpClient{
		Path: jsAPITicketURI,
		Parameters: []wx.QueryParameter{
			{"access_token", token},
			{"type", wx.TicketTypeJSPAI},
		},
		Timeout: timeout,
	}

	ticket := new(wx.APITicket)
	err := req.Get(ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}
