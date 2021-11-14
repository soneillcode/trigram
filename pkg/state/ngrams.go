package state

import (
	"math/rand"
	"strings"
	"sync"
)

// Ngrams stores ngrams such as bigrams (to, be) and trigrams(to, be, or) with the last word
// in the list being stored as frequency data. This frequency data is used to randomly generate words based on
// the preceding words.
// Example:
// if the the bigram (to, be) is stored twice and the bigram (to, go) is stored once, the func Get(to) will
// randomly return (be) twice as many times as (go)
type Ngrams interface {

	// Store takes a list of at least 2 words. All the words except the last are considered the key, and hte last word
	// is the word stored. The word's frequency is increased by one for ever addition of that word for that key.
	Store(words ...string)

	// Get returns a randomly selected word based on the frequency of words for a given ordered list of words.
	// The frequency of a word is the number of times a word has been stored for the given list of words.
	Get(words ...string) string
}

type wordFreq struct {
	mutex sync.RWMutex
	total int
	words map[string]int
}

func newWordFreq() *wordFreq {
	return &wordFreq{
		words: map[string]int{},
		mutex: sync.RWMutex{},
	}
}

func (wf *wordFreq) add(word string) {
	wf.mutex.Lock()
	defer wf.mutex.Unlock()

	freq, ok := wf.words[word]
	if !ok {
		wf.words[word] = 1
	} else {
		wf.words[word] = freq + 1
	}
	wf.total = wf.total + 1
}

func (wf *wordFreq) get(random *rand.Rand) string {
	wf.mutex.RLock()
	defer wf.mutex.RUnlock()

	var cdf = 0
	randomInt := random.Intn(wf.total) + 1

	for word, freq := range wf.words {
		cdf = cdf + freq
		if randomInt <= cdf {
			return word
		}
	}
	return ""
}

const keySeparator = "-"

func getKey(words ...string) string {
	return strings.Join(words, keySeparator)
}

func getKeyAndWord(words ...string) (string, string) {
	length := len(words)
	if length == 0 {
		return "", ""
	}
	if length == 1 {
		return words[0], ""
	}
	if length == 2 {
		return words[0], words[1]
	}
	return getKey(words[:length-1]...), words[length-1]
}
