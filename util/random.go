package util

import (
	"math/rand"
	"time"
)

var randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

func ExecuteRandomFunc(funcs []func()) {
	r := randomGenerator.Intn(len(funcs))

	funcs[r]()
}
