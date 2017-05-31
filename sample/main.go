package main

import (
	"github.com/MenInBack/weshin/userauthorize"
	"github.com/MenInBack/weshin/wx"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	helloURI     = "YourServerAddress/hello"
	callbackURI  = "YourServerAddress/oauthcallback"
	defaultState = "STATE"
)

func init() {
	userauthorize.WXConfig = wx.Config{
		AppID:  "YourAppID",
		Secret: "YourSecret",
	}
}

func main() {
	http.HandleFunc("/oauthcallback", OAuthCallback)
	http.HandleFunc("/hello", Hello)
	err := http.ListenAndServe("YourServerAddress", nil)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

func Hello(w http.ResponseWriter, req *http.Request) {
	log.Print("got hello req: ", req)
	path := strings.Trim(req.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) <= 1 {
		log.Print("unknown user")
		jumpURI, err := userauthorize.JumpToAuth(wx.OAuthScopeBase, callbackURI, defaultState)
		if err != nil {
			log.Print("jumpURI error: ", err)
			return
		}
		http.Redirect(w, req, jumpURI, http.StatusSeeOther)
		return
	}

	io.WriteString(w, "hello "+parts[1])

}

func OAuthCallback(w http.ResponseWriter, req *http.Request) {
	log.Print("got OAuth callback")

	q := req.URL.Query()
	code := q["code"][0]
	state := q["state"][0]

	if len(code) <= 0 {
		log.Print("invlaid code")
		return
	}
	if state != defaultState {
		log.Print("unmatched state")
		return
	}

	token, err := userauthorize.GrantAuthorizeToken(code, 0)
	if err != nil {
		log.Print("GrantAuthorizeToken error: ", err)
		return
	}

	userinfo, err := userauthorize.GetUserInfo(token.AccessToken, token.OpenID, "", 0)
	if err != nil {
		log.Print("GetUserInfo error: ", err)
		return
	}

	log.Print("got user info: ", userinfo)

	http.Redirect(w, req, helloURI+userinfo.Nickname, http.StatusSeeOther)

}
