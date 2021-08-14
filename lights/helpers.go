package lights

import (
	"fmt"
	"math/rand"
	"time"
)

// generateID returns a 10 character id number
// It is prefixed with the unix timestamp and a random string
func generateID() string {
	timeNow := time.Now().Unix()
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return fmt.Sprintf("%d%s", timeNow, string(b))
}
