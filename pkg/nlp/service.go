package nlp

import (
	"fmt"
	"strings"

	"example.com/todo/pkg/state"
)

/*
	The nlp (Natural Language Processing) service.
*/
type Service struct {
	defaultNumberSentences int
	standardNgrams         *state.Ngrams
}

func NewService() *Service {
	return &Service{
		defaultNumberSentences: 12,
		standardNgrams:         state.NewNgrams(),
	}
}

func (s *Service) Learn(text string) error {
	if text == "" {
		return fmt.Errorf("missing data to learn")
	}

	// consider processing tokens as they are created, to save storing them all.
	tokens := tokenize(text)
	processTokens(tokens, s.standardNgrams)
	return nil
}

func processTokens(tokens []string, ngrams *state.Ngrams) {
	length := len(tokens)
	var currentNgrams = ngrams

	// consider using a stream of tokens instead of manual index handling
	for i := 0; i < length-2; i = i + 1 {
		current := tokens[i]
		next := tokens[i+1]
		nextAgain := tokens[i+2]

		if current == state.MagicStartToken {
			currentNgrams.StoreBigram(current, next)
		}

		// we don't store data that crosses magic tokens, this allows flexibility when generating
		if next == state.MagicStartToken {
			continue
		}
		currentNgrams.StoreTrigram(current, next, nextAgain)
	}
}

func (s *Service) Generate() (*string, error) {

	var builder strings.Builder
	builder.Grow(500)

	var newWord string
	var word1, word2 = getStartingWords(s.standardNgrams, state.MagicStartToken)
	builder.WriteString(word1)
	builder.WriteString(spaceWord)
	builder.WriteString(word2)

	maxNumTokens := 1000
	numTokens := 0
	for numSentences := 0; numSentences < s.defaultNumberSentences; {

		if word2 == state.MagicStartToken {
			newWord = s.standardNgrams.GetBigram(word2)
		} else {
			newWord = s.standardNgrams.GetTrigram(word1, word2)
		}

		word1 = word2
		word2 = newWord

		if newWord == state.MagicStartToken {
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

func getStartingWords(ngram *state.Ngrams, startToken string) (string, string) {
	// prevent infinite loops, consider getting a 'safe' Bigram that's guaranteed to be a word
	var saneLimit = 1000
	var currentIterations = 0

	var word1 = ngram.GetBigram(startToken)
	for word1 == state.MagicStartToken {
		currentIterations = currentIterations + 1
		if currentIterations > saneLimit {
			break
		}
		word1 = ngram.GetBigram(startToken)
	}

	currentIterations = 0
	var word2 = ngram.GetTrigram(startToken, word1)
	for word2 == state.MagicStartToken {
		currentIterations = currentIterations + 1
		if currentIterations > saneLimit {
			break
		}
		word2 = ngram.GetTrigram(startToken, word1)
	}

	return word1, word2
}
