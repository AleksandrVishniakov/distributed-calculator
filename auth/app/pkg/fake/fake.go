package fake

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func RandPassword() string {
	rand.NewSource(time.Now().UnixNano())
	return gofakeit.Password(true, true, true, true, true, rand.Intn(56)+8)
}
