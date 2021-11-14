package state

import "testing"

func Test_BstNgrams(t *testing.T) {
	testNgramImpl("bst", NewBstNgrams, t)
}
