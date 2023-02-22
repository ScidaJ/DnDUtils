package commands

import (
	"math/rand"
)

func randomColor() *int {
	min := 0
	max := 16777215
	random := rand.Intn(max-min+1) + min
	return &random
}
