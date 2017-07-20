package crypto

import (
	"log"
	"testing"
)

func TestKeyedSignature(t *testing.T) {
	sig := KeyedSignatured(map[string]string{
		"jsapi_ticket": "sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90u5h9nbSlYy3-Sl-HhTdfl2fzFy1AOcHKP7qg",
		"timestamp":    "1414587457",
		"url":          "http://mp.weixin.qq.com?params=value",
		"noncestr":     "Wm3WZYTPz0wzccnW",
	})

	signature := "0f9de62fce790f9a083d5c99e95740ceb90c27ed"
	log.Println(string(sig))
	if string(sig) != signature {
		t.Error("incorrect signature")
	}
}
