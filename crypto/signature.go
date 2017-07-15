package crypto

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"sort"
	"strings"
)

// KeyedSignatured signature url query parameters with keys, used for jsapi ticket.
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421141115
func KeyedSignatured(values map[string]string) []byte {
	var vals []string
	var b bytes.Buffer
	for k, v := range values {
		b.Reset()
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(v)
		vals = append(vals, b.String())
	}
	sort.Strings(vals)
	sum := sha1.Sum([]byte(strings.Join(vals, "")))
	data := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(data, sum[:])
	return data
}

// Signature signature generator for wechat message
func Signature(message []string) []byte {
	sort.Strings(message)
	sum := sha1.Sum([]byte(strings.Join(message, "")))
	data := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(data, sum[:])
	return data
}
