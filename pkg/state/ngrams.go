package state

import (
	"math/rand"
	"sync"
)

// Ngrams stores ngrams such as bigrams (to, be) and trigrams(to, be, or) with the last word
// in the list being stored as frequency data. This frequency data is used to randomly generate words based on
// the preceding words.
// Example:
// if the the bigram (to, be) is stored twice and the bigram (to, go) is stored once, the func Get(to) will
// randomly return (be) twice as many times as (go)
type Ngrams interface {
	Store(words ...string)
	Get(words ...string) string
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
