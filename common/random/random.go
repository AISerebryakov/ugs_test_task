package random

import (
	"math/rand"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes = "0123456789"
	symbolBytes = letterBytes + numberBytes

	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 53 / letterIdxBits

	numberIdxBits = 6
	numberIdxMask = 1<<numberIdxBits - 1
	numberIdxMax  = 11 / numberIdxBits

	symbolIdxBits = 6
	symbolIdxMask = 1<<symbolIdxBits - 1
	symbolIdxMax  = 63 / symbolIdxBits
)

func String(n int) string {
	if n == 0 {
		return ""
	}
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), symbolIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), symbolIdxMax
		}
		if idx := int(cache & symbolIdxMask); idx < len(symbolBytes) {
			b[i] = symbolBytes[idx]
			i--
		}
		cache >>= symbolIdxBits
		remain--
	}
	return string(b)
}

func Numbers(n int) string {
	if n == 0 {
		return ""
	}
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), numberIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), numberIdxMax
		}
		if idx := int(cache & numberIdxMask); idx < len(numberBytes) {
			b[i] = numberBytes[idx]
			i--
		}
		cache >>= numberIdxBits
		remain--
	}
	return string(b)
}

func Letters(n int) string {
	if n == 0 {
		return ""
	}
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func GenerateRequestId() string {
	return String(32)
}
