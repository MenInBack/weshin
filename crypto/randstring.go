package crypto

import (
	crand "crypto/rand"
	"math/rand"
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

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func RandString(n int) []byte {
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
	return s
}
