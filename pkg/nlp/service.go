package nlp

import (
	"fmt"
	"strings"

	"example.com/todo/pkg/state"
)

// Service provides a Learn and Generate feature which stores data as ngrams and uses the ngrams to randomly generate
// text based on stored word frequency.
type Service struct {
	defaultNumberSentences int
	standardNgrams         state.Ngrams
}

func NewService() *Service {
	return &Service{
		defaultNumberSentences: 12,
		standardNgrams:         state.NewHashNgrams(),
	}
}

// Learn takes a body of text, tokenizes it and stores the tokens as ngrams with their frequency.
func (s *Service) Learn(text string) error {
	if text == "" {
		return fmt.Errorf("missing data to learn")
	}

	// consider processing tokens as they are created, to save storing them all.
	tokens := getTokens(text)
	storeTokens(tokens, s.standardNgrams)
	return nil
}

// Generate uses word frequency data to randomly generate a body of text.
func (s *Service) Generate() (*string, error) {

	var builder strings.Builder
	builder.Grow(500)

	var newWord string
	var word1, word2 = getStartingWords(s.standardNgrams, MagicStartToken)
	builder.WriteString(word1)
	builder.WriteString(spaceWord)
	builder.WriteString(word2)

	maxNumTokens := 1000
	numTokens := 0
	for numSentences := 0; numSentences < s.defaultNumberSentences; {

		if word2 == MagicStartToken {
			newWord = s.standardNgrams.GetBigram(word2)
		} else {
			newWord = s.standardNgrams.GetTrigram(word1, word2)
		}

		word1 = word2
		word2 = newWord

		if newWord == MagicStartToken {
			numSentences = numSentences + 1
			continue
		}

		if newWord != fullStopWord {
			builder.WriteString(spaceWord)
		}

		builder.WriteString(newWord)

		numTokens = numTokens + 1
		if numTokens > maxNumTokens {
			break
		}
	}
	builder.WriteString("\n")
	text := builder.String()
	return &text, nil
}

// Minor issue, sometimes we return Magic start tokens. The normal loop deals with this condition but
// we need special handling for the initial words.
func getStartingWords(ngram state.Ngrams, startToken string) (string, string) {
	var word1 = ngram.GetBigram(startToken)
	var word2 = ngram.GetTrigram(startToken, word1)
	if word2 == MagicStartToken {
		word2 = ngram.GetBigram(startToken)
	}
	return word1, word2
}
