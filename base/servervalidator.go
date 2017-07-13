package base

// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421135319

import (
	"net/http"

	"github.com/MenInBack/weshin/crypto"
)

func StartServerValidator(address, path string) {
	http.HandleFunc(path, serverValidator)
	go http.ListenAndServe(address, nil)
}

func serverValidator(w http.ResponseWriter, req *http.Request) {
	queries := req.URL.Query()
	signature := queries.Get("signature")
	timestamp := queries.Get("timestamp")
	nonce := queries.Get("nonce")
	echostr := queries.Get("echostr")

	sig := string(crypto.Signature([]string{timestamp, nonce, echostr}))
	if sig == signature {
		w.Write([]byte(echostr))
	}
}
