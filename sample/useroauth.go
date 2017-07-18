package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/MenInBack/weshin/base"
	"github.com/MenInBack/weshin/webapi"
	"github.com/MenInBack/weshin/wx"
)

const (
	defaultState = "STATE"
)

var config struct {
	AppID       string `json:"appID,omitempty"`
	Secret      string `json:"secret,omitempty"`
	Address     string `json:"address,omitempty"`
	HelloURI    string `json:"helloURI,omitempty"`
	CallbackURI string `json:"callbackURI,omitempty"`
}

var mp base.MP
var api webapi.WebAPI

func init() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("os.Open error: ", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("ioutil.ReadAll error: ", err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("json.Unmarshal error: ", err)
	}

	mp = base.MP{
		AppID:   config.AppID,
		Secret:  config.Secret,
		Storage: new(sampleStorage),
	}

	api = webapi.WebAPI{
		Mode:     wx.ModeMP,
		WechatMP: mp,
	}
}

func StartOAuthServer() {
	http.HandleFunc("", Hello)
	http.HandleFunc("", OAuthCallback)
	err := http.ListenAndServe(config.Address, nil)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func Hello(w http.ResponseWriter, req *http.Request) {
	log.Print("got hello req: ", req)
	name := req.URL.Query().Get("name")

	if name == "" {
		log.Print("unknown user")
		jumpURI, err := api.JumpToAuth(wx.OAUthScopeUserInfo, config.CallbackURI, defaultState)
		if err != nil {
			log.Print("jumpURI error: ", err)
			return
		}
		http.Redirect(w, req, jumpURI, http.StatusSeeOther)
		return
	}

	w.Write([]byte("hello " + name))

}

func OAuthCallback(w http.ResponseWriter, req *http.Request) {
	log.Print("got OAuth callback")

	q := req.URL.Query()
	code := q.Get("code")
	state := q.Get("state")

	if len(code) <= 0 {
		log.Print("invlaid code")
		return
	}
	if state != defaultState {
		log.Print("unmatched state")
		return
	}

	token, err := api.GrantAuthorizeToken(code, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("got access token %+v", token)

	userinfo, err := api.GetUserInfo(token.OpenID, "", 0)
	if err != nil {
		log.Print("GetUserInfo error: ", err)
		return
	}

	log.Print("got user info: ", userinfo)

	http.Redirect(w, req, config.HelloURI+"?name="+userinfo.Nickname, http.StatusSeeOther)
}

// implements TokenStorage, without refreshing.
type sampleStorage struct {
	token       string
	jsAPITicket string
}

func newsampleStorage() *sampleStorage {
	return new(sampleStorage)
}

func (s *sampleStorage) SetAccessToken(token string, expriresIn int64) {
	s.token = token
}

func (s *sampleStorage) GetAccessToken() string {
	return s.token
}

func (s *sampleStorage) SetAPITicket(ticket *wx.APITicket) {
	if ticket.Typ == wx.TicketTypeJSPAI {
		s.jsAPITicket = ticket.Ticket
	}
}

func (s *sampleStorage) GetAPITicket(typ string) string {
	if typ == wx.TicketTypeJSPAI {
		return s.jsAPITicket
	}
	return ""
}
