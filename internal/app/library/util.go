package library

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func RandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func GenerateTrxID() string {
	now := time.Now()

	return now.Format("060102150405") +
		strconv.Itoa(RandInt(10, 99)) +
		strconv.Itoa(RandInt(10, 99))
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
