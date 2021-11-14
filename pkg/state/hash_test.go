package state

import "testing"

func Test_HashNgrams(t *testing.T) {
	testNgramImpl("hash", NewHashNgrams, t)
}
