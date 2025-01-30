package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	chars         = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func GenerateChallenge(n int) string {
	token := ""
	for i := 0; i < n; i++ {
		token += string(chars[cryptoRandSecure(int64(len(chars)))])
	}

	return token
}

func cryptoRandSecure(n int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		return 0
	}

	return nBig.Int64()
}
