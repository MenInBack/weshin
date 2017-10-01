// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
package pay

import (
	"bytes"
	"crypto/aes"
	"crypto/hmac"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/MenInBack/weshin/wx"
)

const (
	NonceLength = 32
)

// initialize random source with crypto/rand
func init() {
	buf := make([]byte, 8)
	n, e := crand.Read(buf)
	if n != 8 || e != nil {
		panic("init random source failed")
	}
	s := uint64(buf[0])<<56 | uint64(buf[1])<<48 | uint64(buf[2])<<40 | uint64(buf[3])<<32 |
		uint64(buf[4])<<24 | uint64(buf[5])<<16 | uint64(buf[6])<<8 | uint64(buf[7])
	rand.Seed(int64(s))
}

func sign(fields []field, key string, typ SignType) (string, error) {
	if len(fields) <= 0 {
		return "", wx.WeshinError{Detail: "empty query parameter"}
	}
	sort.Slice(fields, func(i, j int) bool { return strings.Compare(fields[i].name, fields[j].name) < 0 })
	fields = append(fields, field{"key", key})

	var buf bytes.Buffer
	for _, f := range fields {
		buf.WriteString(f.name)
		buf.WriteByte('=')
		buf.WriteString(f.value)
		buf.WriteByte('&')
	}
	str := bytes.TrimSuffix(buf.Bytes(), []byte{'&'})

	switch typ {
	case MD5:
		return fmt.Sprintf("%X", md5.Sum(str)), nil
	case HMAC:
		mac := hmac.New(sha256.New, []byte(key))
		mac.Write(str)
		return fmt.Sprintf("%X", mac.Sum(nil)), nil
	}
	return "", wx.WeshinError{Detail: "unknown sign type"}
}

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func randomString(n int) string {
	s := make([]byte, 0, n)
	cache := rand.Uint64()
	remain := 64
	for i := 0; i < n; remain -= 6 {
		if remain < 6 {
			cache, remain = rand.Uint64(), 64
		}
		b := cache & (1<<6 - 1)
		if b < 62 {
			s = append(s, chars[b])
			i++
		}
		cache >>= 6
	}
	return string(s)
}

type field struct {
	name  string
	value string
}

//解密步骤如下：
//（1）对加密串A做base64解码，得到加密串B
//（2）对商户key做md5，得到32位小写key* ( key设置路径：微信商户平台(pay.weixin.qq.com)-->账户设置-->API安全-->密钥设置 )
//（3）用key*对加密串B做AES-256-ECB解密
func decodeNoticeMessage(info, key string) ([]byte, error) {
	cipher := make([]byte, base64.RawStdEncoding.DecodedLen(len(info)))
	base64.RawStdEncoding.Decode(cipher, []byte(info))

	hashKey := md5.Sum([]byte(key))
	hexKey := make([]byte, hex.EncodedLen(len(hashKey)))
	hex.Encode(hexKey, hashKey[:])

	block, e := aes.NewCipher(hexKey)
	if e != nil {
		return nil, e
	}

	// padding cipher
	if len(cipher)%block.BlockSize() != 0 {
		pad := make([]byte, block.BlockSize()-len(cipher)%block.BlockSize())
		cipher = append(cipher, pad...)
	}
	buf := cipher

	// aes-256-ecb decrypt
	for len(cipher) > 0 {
		block.Decrypt(buf, cipher)
		cipher = cipher[block.BlockSize():]
	}

	return buf, nil
}
