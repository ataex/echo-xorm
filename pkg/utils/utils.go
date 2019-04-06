package utils

import (
	"encoding/base64"

	"golang.org/x/crypto/sha3"
)

// GetSHA3Hash ...
func GetSHA3Hash(data string) string {
	h := make([]byte, 64)
	sha3.ShakeSum256(h, []byte(data))
	return base64.StdEncoding.EncodeToString(h)
}
