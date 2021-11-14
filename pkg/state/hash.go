package state

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type HashNgrams struct {
	trigrams map[string]*WordFreq
	triMutex sync.RWMutex
	bigrams  map[string]*WordFreq
	biMutex  sync.RWMutex
	random   *rand.Rand
}

func NewHashNgrams() *HashNgrams {
	return &HashNgrams{
		trigrams: map[string]*WordFreq{},
		bigrams:  map[string]*WordFreq{},
		random:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *HashNgrams) StoreTrigram(word1, word2, word3 string) {
	s.triMutex.Lock()
	defer s.triMutex.Unlock()

	key := getKey(word1, word2)
	wordFreq, ok := s.trigrams[key]
	if !ok {
		s.trigrams[key] = &WordFreq{
			words: map[string]int{},
		}
		wordFreq = s.trigrams[key]
	}
	wordFreq.add(word3)
}

func (s *HashNgrams) StoreBigram(word1, word2 string) {
	s.biMutex.Lock()
	defer s.biMutex.Unlock()

	key := word1
	wordFreq, ok := s.bigrams[key]
	if !ok {
		s.bigrams[key] = &WordFreq{
			words: map[string]int{},
		}
		wordFreq = s.bigrams[key]
	}
	wordFreq.add(word2)
}

func (s *HashNgrams) GetTrigram(word1, word2 string) string {
	s.triMutex.RLock()
	defer s.triMutex.RUnlock()

	key := getKey(word1, word2)
	wordFreq, ok := s.trigrams[key]
	if !ok {
		// return random word or stop ?
		return ""
	}
	return wordFreq.get(s.random)
}

func (s *HashNgrams) GetBigram(word1 string) string {
	s.biMutex.RLock()
	defer s.biMutex.RUnlock()

	key := word1
	wordFreq, ok := s.bigrams[key]
	if !ok {
		// return random word or stop ?
		return ""
	}
	return wordFreq.get(s.random)
}

const keySeparator = "-"

func getKey(word1, word2 string) string {
	return fmt.Sprintf("%s%s%s", word1, keySeparator, word2)
}
