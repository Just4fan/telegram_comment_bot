package utils

import (
	"encoding/base64"
)

func Base64Encode(src []byte) (dst string, err error) {
	encoding := base64.RawURLEncoding
	dst = encoding.EncodeToString(src)
	return
}

func Base64Decode(src string) (dst []byte, err error) {
	encoding := base64.RawURLEncoding
	dst, err = encoding.DecodeString(src)
	return
}
