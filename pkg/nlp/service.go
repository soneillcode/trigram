package nlp

import (
	"strings"

	"example.com/todo/pkg/state"
)

/*
	The nlp (Natural Language Processing) service.
*/

type Service struct {
	defaultNumberSentences int
	state                  *state.Ngrams
}

func NewService() *Service {
	return &Service{
		defaultNumberSentences: 12,
		state:                  state.NewState(),
	}
}

func (s *Service) Learn(text string) error {
	if text == "" {
		return nil
	}

	// consider adding tokens to state as they are created to save storing them all.
	tokens := tokenize(text)

	for i, t := range tokens {
		if i == 0 {
			s.state.StoreBigram(state.MagicStartToken, t)
		}
		if i == 1 {
			s.state.StoreTrigram(state.MagicStartToken, tokens[i-1], t)
		}
		if i > 1 {
			word1 := tokens[i-2]
			word2 := tokens[i-1]
			if word2 == state.MagicSentenceToken {
				// we don't store the end of the sentence in relation to the start of one
				s.state.StoreBigram(word2, t)
			} else {
				s.state.StoreTrigram(word1, word2, t)
			}
		}
	}

	return nil
}

func (s *Service) Generate() (*string, error) {

	var builder strings.Builder
	builder.Grow(100)

	var newWord string
	var word1 = s.state.GetBigram(state.MagicStartToken)
	var word2 = s.state.GetTrigram(state.MagicStartToken, word1)

	maxNumTokens := 1000
	numTokens := 0
	for numSentences := 0; numSentences < s.defaultNumberSentences; {
		if word2 == state.MagicSentenceToken {
			newWord = s.state.GetBigram(word2)
		} else {
			newWord = s.state.GetTrigram(word1, word2)
		}
		word1 = word2
		word2 = newWord

		if newWord == state.MagicSentenceToken {
			numSentences = numSentences + 1
			builder.WriteString(".")
			continue
		}
		builder.WriteRune(' ')
		if newWord == state.MagicStartToken {
			builder.WriteString(" ")
			continue
		}
		builder.WriteString(newWord)
		numTokens = numTokens + 1
		if numTokens > maxNumTokens {
			break
		}
	}
	text := builder.String()
	return &text, nil
}
