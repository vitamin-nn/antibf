package integration

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateRandStr(n int) string { // nolint: deadcode,unused
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateRandIP() string { // nolint: deadcode,unused
	var s [4]string
	for i := 0; i < 4; i++ {
		s[i] = strconv.Itoa(rand.Intn(255) + 1)
	}
	return strings.Join(s[:], ".")
}
