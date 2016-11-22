package cryptoutils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5Encode(base64.URLEncoding.EncodeToString(b))

}

func Get16Md5Encode(src string) string {
	runes := []rune(Md5Encode(src))
	return string(runes[7:23])
}

func Md5Encode(src string) string {
	h := md5.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func Base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}
func Base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}
