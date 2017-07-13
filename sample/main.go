package main

import (
	"log"
	"net/http"

	"github.com/MenInBack/weshin/useroauth"
	"github.com/MenInBack/weshin/wx"
)

const (
	appID  = ""
	secret = ""
	token  = "-itCg_1p7WX8MnjHVgHLA2MEywstb2I0JUUGaAHAFFR"

	address      = ""
	helloURI     = ""
	callbackURI  = ""
	defaultState = "STATE"
)

func main() {
	http.HandleFunc("", Hello)
	http.HandleFunc("", OAuthCallback)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func Hello(w http.ResponseWriter, req *http.Request) {
	log.Print("got hello req: ", req)
	name := req.URL.Query().Get("name")

	if name == "" {
		log.Print("unknown user")
		oAuth := useroauth.New(appID, secret)
		jumpURI, err := oAuth.JumpToAuth(wx.OAUthScopeUserInfo, callbackURI, defaultState)
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

	oAuth := useroauth.New(appID, secret)
	token, err := oAuth.GrantAuthorizeToken(code, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("got access token %+v", token)

	userinfo, err := useroauth.GetUserInfo(token.OpenID, token.AccessToken, "", 0)
	if err != nil {
		log.Print("GetUserInfo error: ", err)
		return
	}

	log.Print("got user info: ", userinfo)

	http.Redirect(w, req, helloURI+"?name="+userinfo.Nickname, http.StatusSeeOther)
}
