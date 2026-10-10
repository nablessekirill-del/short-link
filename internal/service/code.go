package service

import (
	"crypto/rand"
	"math/big"
)

var reserved = map[string]struct{}{
	"api":     {},
	"health":  {},
	"metrics": {},
}

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewCode(n int) (string, error) {
	b := make([]byte, n)
	max := big.NewInt(int64(len(alphabet)))

	for i := range b {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = alphabet[idx.Int64()]
	}
	return string(b), nil

}

func IsReserved(code string) bool {
	_, ok := reserved[code]
	return ok
}
