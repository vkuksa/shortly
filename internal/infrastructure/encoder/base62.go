package encoder

import (
	"crypto/sha256"
	"math/big"
)

const base62chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const length = 10

type Base62 struct {
	length int
}

// NewBase62 returns a new Base62 encoder with the specified output length.
func NewBase62() Base62 {
	if length <= 10 {
		panic("length must be greater than 10")
	}
	return Base62{length: length}
}

// Encode generates a securely random Base62 encoded string.
func (s Base62) Encode(original string) (string, error) {
	hash := sha256.Sum256([]byte(original))
	num := new(big.Int).SetBytes(hash[:])
	base := big.NewInt(62)
	encoded := ""

	for num.Sign() > 0 {
		mod := new(big.Int)
		num.DivMod(num, base, mod)
		encoded = string(base62chars[mod.Int64()]) + encoded
	}
	return encoded[:s.length], nil
}
