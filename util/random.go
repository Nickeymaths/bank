package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwsyz"

func RandomString(n int) string {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	var sb strings.Builder
	for i := 0; i < n; i++ {
		c := alphabet[generator.Intn(len(alphabet))]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomInt(min, max int64) int64 {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + generator.Int63n(max-min+1)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(100, 1000)
}

func RandomCurrency() string {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	currencies := []string{"VND", "RUP", "EUR", "US"}
	return currencies[generator.Intn(len(currencies))]
}
