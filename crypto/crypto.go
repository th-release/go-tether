package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
)

func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func Md5Base64Digest(data string) string {
	hash := md5.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func HmacSignature(secret string, canonicalString string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(canonicalString))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
