package state

import (
	"fmt"
	"math/rand"
	"time"
)

type Ngrams struct {
	trigrams map[string]*WordFreq
	bigrams  map[string]*WordFreq
	random   *rand.Rand
}

const keySeparator = "-"
const MagicStartToken = "MAGIC_START_TOKEN"
const MagicSentenceToken = "MAGIC_SENTENCE_TOKEN"
const MagicDialogToken = "MAGIC_DIALOG_TOKEN"

func NewState() *Ngrams {
	return &Ngrams{
		trigrams: map[string]*WordFreq{},
		bigrams:  map[string]*WordFreq{},
		random:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Ngrams) StoreTrigram(word1, word2, word3 string) {
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

func (s *Ngrams) StoreBigram(word1, word2 string) {
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

func (s *Ngrams) GetTrigram(word1, word2 string) string {
	key := getKey(word1, word2)
	wordFreq, ok := s.trigrams[key]
	if !ok {
		// return random word or stop ?
		return ""
	}
	return wordFreq.get(s.random)
}

func (s *Ngrams) GetBigram(word1 string) string {
	key := word1
	wordFreq, ok := s.bigrams[key]
	if !ok {
		// return random word or stop ?
		return ""
	}
	return wordFreq.get(s.random)
}

func getKey(word1, word2 string) string {
	return fmt.Sprintf("%s%s%s", word1, keySeparator, word2)
}

type WordFreq struct {
	total int
	words map[string]int
}

func (t *WordFreq) add(word string) {
	// todo lock
	freq, ok := t.words[word]
	if !ok {
		t.words[word] = 1
	} else {
		t.words[word] = freq + 1
	}
	t.total = t.total + 1
}

func (t *WordFreq) get(random *rand.Rand) string {
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
