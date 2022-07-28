package momo

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func microTime() (float64, int64) {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	now := time.Now().In(loc)
	micSeconds := float64(now.Nanosecond()) / 1000000000
	return micSeconds, now.Unix()
}

func newSHA256(plaintext string) string {
	h := sha256.New()
	h.Write([]byte(plaintext))
	sha := h.Sum(nil)
	shaStr := hex.EncodeToString(sha)
	return shaStr
}
