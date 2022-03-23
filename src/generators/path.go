package generators

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")
var lettersLength = len(letters)
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomPath(minLen int32, maxLen int32) string {
	r := make([]rune, minLen+random.Int31n(maxLen-minLen))
	for i := range r {
		r[i] = letters[random.Intn(lettersLength)]
	}
	return "/" + string(r)
}
