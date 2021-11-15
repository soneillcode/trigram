package state

import (
	"math/rand"
	"testing"
)

func Test_BstNgrams(t *testing.T) {
	random := rand.New(rand.NewSource(16))
	implFunc := func() Ngrams {
		return NewBstNgrams(random)
	}
	name := "bstNgrams"
	testEmptyStoreAndGet(name, implFunc, t)
	testBasicBigramStoreAndGet(name, implFunc, t)
	testBasicTrigramStoreAndGet(name, implFunc, t)
	testWordFrequency(name, implFunc, t)
	testConcurrentAccess(name, implFunc, t)
}
