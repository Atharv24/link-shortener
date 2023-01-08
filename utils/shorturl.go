package utils

import (
	"math/rand"
	"time"
)

const (
	alphanumeric = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateShortURL() string {
	// Generate random alphanumeric string for short URL
	var b [7]byte
	for i := range b {
		b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(b[:])
}
