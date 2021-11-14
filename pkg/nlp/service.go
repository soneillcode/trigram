package nlp

import (
	"strings"

	"example.com/todo/pkg/state"
)

// Service provides a Learn and Generate feature which stores data as ngrams and uses the ngrams to randomly generate
// text based on stored word frequency.
//
// Tokens are handled as arrays of string for simplicity. If performance needs improving,  could be converted to a
// data structure that processes the tokens one by one to reduce copying the arrays.
type Service struct {
	ngrams               state.Ngrams
	defaultNumberOfWords int
}

// NewService returns a new instance of a Service.
func NewService() *Service {
	return &Service{
		ngrams:               state.NewHashNgrams(),
		defaultNumberOfWords: 200,
	}
}

// Learn takes a body of text, tokenizes it and stores the tokens as ngrams with their frequency.
func (s *Service) Learn(text string) {
	if text == "" {
		return
	}

	tokens := getTokens(text)
	storeTokens(tokens, s.ngrams)
}

// Generate uses trigram word frequency data to randomly generate a body of text.
func (s *Service) Generate() (*string, error) {
	tokens := generateTokens(s.ngrams, s.defaultNumberOfWords)
	tokens = filterTokens(tokens)
	tokens = addSpaceTokens(tokens)
	text := toString(tokens)
	return &text, nil
}

func generateTokens(ngrams state.Ngrams, maxTokens int) []string {
	var tokens []string
	word1 := ngrams.Get(MagicStartToken)
	word2 := ngrams.Get(MagicStartToken, word1)
	tokens = append(tokens, word1)
	tokens = append(tokens, word2)

	var newWord string
	for numTokens := 0; numTokens < maxTokens; numTokens = numTokens + 1 {
		if word2 == MagicStartToken {
			// ignore the first word in the context of a new sentence.
			newWord = ngrams.Get(MagicStartToken)
		} else {
			newWord = ngrams.Get(word1, word2)
		}
		word1 = word2
		word2 = newWord
		tokens = append(tokens, newWord)
	}

	tokens = append(tokens, newLineWord)
	return tokens
}

// Consider providing the filter function as an argument.
func filterTokens(tokens []string) []string {
	var filtered []string
	for _, token := range tokens {
		if token == MagicStartToken {
			continue
		}
		filtered = append(filtered, token)
	}
	return filtered
}

func addSpaceTokens(tokens []string) []string {
	var newTokens []string
	for index, token := range tokens {
		if index != 0 && token != fullStopWord {
			newTokens = append(newTokens, spaceWord)
		}
		newTokens = append(newTokens, token)
	}
	return newTokens
}

func toString(tokens []string) string {
	var builder strings.Builder
	builder.Grow(len(tokens))
	for _, token := range tokens {
		builder.WriteString(token)
	}
	return builder.String()
}
