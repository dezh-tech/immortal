package websocket

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/rand"
)

const (
	success      = "success"
	databaseFail = "db_fail"
	parseFail    = "parse_fail"
	authFail     = "parse_fail"
	limitsFail   = "limits_fail"
	serverFail   = "server_fail"
	invalidFail  = "invalid_fail"

	chars         = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func measureLatency(ht prometheus.Histogram) func() {
	start := time.Now()

	return func() {
		ht.Observe(time.Since(start).Seconds())
	}
}

func generateChallenge(n int) string {
	src := rand.NewSource(uint64(time.Now().UnixNano()))
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Uint64(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Uint64(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(chars) {
			b[i] = chars[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
