package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijkmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().Unix())
}

func RandomInt(min, max uint) uint {
	minInt := int64(min)
	maxInt := int64(max)
	return uint(minInt + rand.Int63n(maxInt-minInt+1))
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn((k))]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomUserStrID() string {
	return RandomString(6)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomURL() string {
	return fmt.Sprintf("https://%s.jpg", RandomString(6))
}
