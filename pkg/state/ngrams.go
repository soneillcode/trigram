package state

import (
	"math/rand"
	"sync"
)

type Ngrams interface {
	StoreTrigram(word1, word2, word3 string)
	StoreBigram(word1, word2 string)
	GetTrigram(word1, word2 string) string
	GetBigram(word1 string) string
}

type WordFreq struct {
	mutex sync.RWMutex
	total int
	words map[string]int
}

func (t *WordFreq) add(word string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	freq, ok := t.words[word]
	if !ok {
		t.words[word] = 1
	} else {
		t.words[word] = freq + 1
	}
	t.total = t.total + 1
}

func (t *WordFreq) get(random *rand.Rand) string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	var cdf = 0
	for word, freq := range t.words {
		if cdf == 0 {
			cdf = freq
		}
		if random.Intn(t.total) <= cdf {
			return word
		}
		cdf = cdf + freq
	}
	return ""
}
