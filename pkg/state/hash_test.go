package state

import (
	"math/rand"
	"testing"
)

func Test_HashNgrams(t *testing.T) {
	random := rand.New(rand.NewSource(16))
	implFunc := func() Ngrams {
		return NewHashNgrams(random)
	}
	name := "hashNgrams"
	testEmptyStoreAndGet(name, implFunc, t)
	testBasicBigramStoreAndGet(name, implFunc, t)
	testBasicTrigramStoreAndGet(name, implFunc, t)
	testWordFrequency(name, implFunc, t)
	testConcurrentAccess(name, implFunc, t)
}
