package state

import (
	"math/rand"
	"sync"
	"time"
)

// hashNgrams implements Ngrams using a hash map. It uses hash for fast reads. A hash map is not safe for concurrent
// writes, so we lock a mutex to prevent concurrent writes.
type hashNgrams struct {
	mutex  sync.RWMutex
	ngrams map[string]*wordFreq
	random *rand.Rand
}

func NewHashNgrams() Ngrams {
	return &hashNgrams{
		mutex:  sync.RWMutex{},
		ngrams: map[string]*wordFreq{},
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *hashNgrams) Store(words ...string) {
	key, word := getKeyAndWord(words...)
	if key == "" || word == "" {
		// consider handling this case better
		return
	}

	wf := s.getWordFreq(key)
	wf.add(word)
}

func (s *hashNgrams) getWordFreq(key string) *wordFreq {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	wf, ok := s.ngrams[key]
	if !ok {
		s.ngrams[key] = newWordFreq()
		wf = s.ngrams[key]
	}
	return wf
}

func (s *hashNgrams) Get(words ...string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	key := getKey(words...)

	wf, ok := s.ngrams[key]
	if !ok {
		return ""
	}
	return wf.get(s.random)
}
