package encoder

import (
	"encoding/base64"
)

type Base64 struct {
}

func NewBase64() Base64 {
	return Base64{}
}

func (s Base64) Encode(original string) (string, error) {
	return base64.URLEncoding.EncodeToString([]byte(original)), nil
}
