package randutil

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())

}

// RandBool This function returns a random boolean value based on the current time
func RandBool() bool {
	return rand.Intn(2) == 1
}

func RandString(runes []rune, size int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

	if runes != nil {
		letterRunes = runes
	}

	b := make([]rune, size)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandNumSets(n int) string {
	if n == 0 {
		n = 6
	}
	//random number
	var letterRunes = []rune("1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandNumRange(min int, absDelta int) int {
	if min <= 0 {
		min = 0
	}
	if absDelta == 0 {
		absDelta = 1
	}
	max := min + int(math.Abs(float64(absDelta)))
	return rand.Intn(max-min+1) + min
}
