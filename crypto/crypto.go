package crypto

// wechat message crypto api
// https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1419318482&token=&lang=zh_CN

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
)

const (
	encodingKeySize = 43
	randStringLen   = 16
	aesKeySize      = 32
	msgSizeLength   = 4
)

type MessageCrypto struct {
	Token     string
	AppID     string
	aesKey    []byte
	nonce     string
	timeStamp string
	userName  string
}

// New wechat message crypto
func New(encodingAESKey, token, appID string) (*MessageCrypto, error) {
	if len(encodingAESKey) != encodingKeySize {
		return nil, errors.New("invalid encodingAESKey")
	}
	key, err := base64.RawStdEncoding.DecodeString(encodingAESKey)
	if err != nil {
		return nil, err
	}
	return &MessageCrypto{
		Token:  token,
		AppID:  appID,
		aesKey: key,
	}, nil
}

// Encrypt a message object into xml marshaled encrypted message
func (mc *MessageCrypto) Encrypt(msg []byte, nonce, timestamp string) (data []byte, err error) {
	//1.add rand str ,len, appid
	pad, err := mc.messagePadding(msg)
	if err != nil {
		return nil, CryptoError{"padding failed", err}
	}

	//2. AES Encrypt
	block, err := aes.NewCipher(mc.aesKey)
	if err != nil {
		return nil, CryptoError{"aes.NewCipher failed", err}
	}
	mode := cipher.NewCBCEncrypter(block, mc.aesKey[:aes.BlockSize])
	// pad = pad[aes.BlockSize:]
	mode.CryptBlocks(pad, pad)

	//3. base64Encode
	buf := make([]byte, base64.RawStdEncoding.EncodedLen(len(pad)))
	base64.RawStdEncoding.Encode(buf, pad)

	//4. compute signature
	if len(nonce) > 0 {
		mc.nonce = nonce
	}
	if len(timestamp) > 0 {
		mc.timeStamp = timestamp
	}

	sign := mc.signature(buf)

	//5. Gen xml
	rply := encryptedMessage{
		ToUserName:   cdata{mc.userName},
		Encrypt:      cdata{string(buf)},
		MsgSignature: cdata{string(sign)},
		TimeStamp:    mc.timeStamp,
		Nonce:        cdata{mc.nonce},
	}
	data, err = xml.Marshal(rply)
	if err != nil {
		return nil, CryptoError{"xml.Marshal failed", err}
	}

	return data, nil
}

// Decrypt and validate wechat message
func (mc *MessageCrypto) Decrypt(src []byte, signature, nonce, timestamp string) (msg []byte, err error) {
	//1.validate xml format
	encrypted := new(encryptedMessage)
	err = xml.Unmarshal(src, encrypted)
	if err != nil {
		return nil, CryptoError{"xml.Unmarshal failed", err}
	}
	if len(encrypted.Encrypt.Data) == 0 {
		return nil, CryptoError{"invalid message", errors.New("got nothing")}
	}
	if len(nonce) > 0 {
		mc.nonce = nonce
	} else {
		mc.nonce = encrypted.Nonce.Data
	}
	if len(timestamp) > 0 {
		mc.timeStamp = timestamp
	} else {
		mc.timeStamp = encrypted.TimeStamp
	}
	mc.userName = encrypted.FromUserName.Data

	//2.validate signature
	sign := mc.signature([]byte(encrypted.Encrypt.Data))
	if len(signature) <= 0 {
		signature = encrypted.MsgSignature.Data
	}
	if !bytes.Equal(sign, []byte(signature)) {
		return nil, CryptoError{"invalid message", errors.New("signature mismatch")}
	}

	//3.decode base64
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(encrypted.Encrypt.Data)))
	_, err = base64.StdEncoding.Decode(buf, []byte(encrypted.Encrypt.Data))
	if err != nil {
		return nil, CryptoError{"base64.Decode failed ", err}
	}
	l := len(buf)
	l = l / aes.BlockSize * aes.BlockSize
	buf = buf[:l] // drop base64 trailling zeros, resize buf to multiplex of aes.BlockSize

	//4.decode aes
	block, err := aes.NewCipher(mc.aesKey)
	if err != nil {
		return nil, CryptoError{"aes.NewCipher", err}
	}
	mode := cipher.NewCBCDecrypter(block, mc.aesKey[:aes.BlockSize])
	mode.CryptBlocks(buf, buf)

	// 5. remove rand str, appid and trailling padding
	data, err := mc.messageUnpadding(buf)
	if err != nil {
		return nil, CryptoError{"unpadding message failed", err}
	}

	return data, nil
}

// random(16B) + msg_len(4B) + msg + appid
func (mc *MessageCrypto) messagePadding(msg []byte) ([]byte, error) {
	size := new(bytes.Buffer)
	err := binary.Write(size, binary.BigEndian, int32(len(msg)))
	if err != nil {
		return nil, err
	}
	rStr, err := randString(randStringLen)
	if err != nil {
		return nil, err
	}

	data := bytes.NewBuffer(rStr)
	data.Write(size.Bytes())
	data.Write(msg)
	data.Write([]byte(mc.AppID))

	l := data.Len()
	padding := aesKeySize - l%aesKeySize
	for i := 0; i < padding; i++ {
		data.WriteByte(byte(padding))
	}
	return data.Bytes(), nil
}

func (mc *MessageCrypto) messageUnpadding(src []byte) ([]byte, error) {
	if len(src) <= randStringLen+msgSizeLength {
		return nil, errors.New("decrypted message too short")
	}
	src = src[randStringLen:] // drop random string
	sizeBuf := bytes.NewBuffer(src[:msgSizeLength])
	var size int32
	err := binary.Read(sizeBuf, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}

	src = src[msgSizeLength:] // drop message size
	if len(src) <= int(size) {
		return nil, errors.New("unpadded message too short")
	}

	appid := src[int(size):]
	if !bytes.Equal(appid, []byte(mc.AppID)) {
		return nil, errors.New("appid mismatch")
	}

	return src[:int(size)], nil
}

// msg_signature=sha1(sort(Token、timestamp、nonce, msg_encrypt))
func (mc *MessageCrypto) signature(encrypt []byte) []byte {
	return Signature([]string{mc.Token, mc.timeStamp, mc.nonce, string(encrypt)})
}

type encryptedMessage struct {
	XMLName      xml.Name `xml:"xml"`
	FromUserName cdata    `xml:"FromUserName,omitempty"`
	ToUserName   cdata    `xml:"ToUserName,omitempty"`
	Encrypt      cdata    `xml:"Encrypt"`
	MsgSignature cdata    `xml:"MsgSignature"`
	TimeStamp    string   `xml:"TimeStamp"`
	Nonce        cdata    `xml:"Nonce"`
}

type cdata struct {
	Data string `xml:",cdata"`
}

func randString(len int) ([]byte, error) {
	byteLen := base64.RawStdEncoding.DecodedLen(len) // no padding
	buf := make([]byte, byteLen)
	n, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}
	if n != byteLen {
		return nil, errors.New("rand.Read exception")
	}
	data := make([]byte, len)
	base64.RawStdEncoding.Encode(data, buf)
	return data[:len], nil
}

type CryptoError struct {
	Detail string
	Err    error
}

func (e CryptoError) Error() string {
	return fmt.Sprintf("crypto error - %s: %s", e.Detail, e.Err.Error())
}
