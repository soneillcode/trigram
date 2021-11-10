package state

import (
	"fmt"
	"math/rand"
	"time"
)

type State struct {
	trigrams map[string]*Trigram
	random   *rand.Rand
}

const keySeparator = " "
const MagicStartToken = "aaaaaa"
const MagicWildCard = "EYSDHXCHYREccvdf363"

func NewState() *State {
	return &State{
		trigrams: map[string]*Trigram{},
		random:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *State) Store(word1, word2, word3 string) {
	key := getKey(word1, word2)
	tri, ok := s.trigrams[key]
	if !ok {
		tri = s.add(key)
	}
	tri.add(word3)
}

func (s *State) add(key string) *Trigram {
	// todo lock
	s.trigrams[key] = &Trigram{
		words: map[string]int{},
	}
	return s.trigrams[key]
}

func (s *State) Get(word1, word2 string) string {
	key := getKey(word1, word2)
	tri, ok := s.trigrams[key]
	if !ok {
		// return random word or stop ?
		return ""
	}
	return tri.get(s.random)
}

func getKey(word1, word2 string) string {
	if word2 == "" {
		return MagicStartToken
	}
	if word2 == "." {
		return MagicStartToken
	}
	if word1 == "" {
		word1 = MagicStartToken
	}
	return fmt.Sprintf("%s%s%s", word1, keySeparator, word2)
}

type Trigram struct {
	total int
	words map[string]int
}

func (t *Trigram) add(word string) {
	// todo lock
	freq, ok := t.words[word]
	if !ok {
		t.words[word] = 1
	} else {
		t.words[word] = freq + 1
	}
	t.total = t.total + 1
}

func (t *Trigram) get(random *rand.Rand) string {
	// todo lock
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
