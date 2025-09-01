package common

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

const (
	RandomCharset = "abcdefghijklmnopqrstuvwxyz" + // lowercase
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" + // uppercase
		"0123456789" // digits
	MinRndStrLen = 16
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GetRandomString(length int) string {
	return StringWithCharset(length, RandomCharset)
}

func GetSHA256Sum(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
