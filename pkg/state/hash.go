package state

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

// HashNgrams implements Ngrams using a hash map. It uses hash for fast reads. A hash map is not safe for concurrent
// writes, so we lock a mutex to prevent concurrent writes.
type HashNgrams struct {
	mutex  sync.RWMutex
	ngrams map[string]*WordFreq
	random *rand.Rand
}

func NewHashNgrams() *HashNgrams {
	return &HashNgrams{
		ngrams: map[string]*WordFreq{},
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *HashNgrams) Store(words ...string) {
	key, word := getKeyAndWord(words...)
	if key == "" || word == "" {
		// consider handling this edge better
		return
	}

	wordFreq := s.getWordFreq(key)
	wordFreq.add(word)
}

func (s *HashNgrams) getWordFreq(key string) *WordFreq {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	wordFreq, ok := s.ngrams[key]
	if !ok {
		s.ngrams[key] = &WordFreq{
			words: map[string]int{},
		}
		wordFreq = s.ngrams[key]
	}
	return wordFreq
}

func (s *HashNgrams) Get(words ...string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	key := getKey(words...)

	wordFreq, ok := s.ngrams[key]
	if !ok {
		// return random word or stop ?
		return ""
	}
	return wordFreq.get(s.random)
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
