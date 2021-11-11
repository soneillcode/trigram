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
	dialogueState          *state.Ngrams
}

func NewService() *Service {
	return &Service{
		defaultNumberSentences: 12,
		state:                  state.NewState(),
		dialogueState:          state.NewState(),
	}
}

func (s *Service) Learn(text string) error {
	if text == "" {
		return nil
	}

	// todo performance -> consider adding tokens to state as they are created to save storing them all.
	tokens := tokenize(text)

	processTokens(tokens, s.state, s.dialogueState)

	return nil
}

func processTokens(tokens []string, ngrams *state.Ngrams, dialogTokens *state.Ngrams) int {
	length := len(tokens)
	for i := 0; i < length; i = i + 1 {
		t := tokens[i]
		if i == 0 {
			ngrams.StoreBigram(state.MagicStartToken, t)
		}
		if i == 1 {
			ngrams.StoreTrigram(state.MagicStartToken, tokens[i-1], t)
		}
		if i > 1 {
			word1 := tokens[i-2]
			word2 := tokens[i-1]

			if word2 == state.MagicSentenceToken || word2 == state.MagicDialogToken {
				// we don't store the end of the sentence in relation to the start of one
				ngrams.StoreBigram(word2, t)
			} else {
				ngrams.StoreTrigram(word1, word2, t)
			}

			if t == state.MagicDialogToken {
				newIndex := processDialogTokens(tokens[i+1:], dialogTokens)
				i = i + newIndex + 1 // 2 for the dialog end token
			}
		}
	}
	return len(tokens) - 1
}

func processDialogTokens(tokens []string, ngrams *state.Ngrams) int {
	length := len(tokens)
	for i := 0; i < length; i = i + 1 {
		t := tokens[i]
		if i == 0 {
			ngrams.StoreBigram(state.MagicDialogToken, t)
		}
		if i == 1 {
			ngrams.StoreTrigram(state.MagicDialogToken, tokens[i-1], t)
		}
		if i > 1 {
			word1 := tokens[i-2]
			word2 := tokens[i-1]

			if word2 == state.MagicSentenceToken {
				// we don't store the end of the sentence in relation to the start of one
				ngrams.StoreBigram(word2, t)
			} else {
				ngrams.StoreTrigram(word1, word2, t)
			}
			if t == state.MagicDialogToken {
				return i
			}
		}
	}
	return len(tokens) - 1 // todo an error in the data in this case to be honest
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
		if word2 == state.MagicSentenceToken || word2 == state.MagicDialogToken {
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

		if newWord == state.MagicDialogToken {
			dialog := s.generateDialog()
			builder.WriteString(dialog)
			continue
		}

		builder.WriteRune(' ')
		if newWord == state.MagicStartToken {
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

func (s *Service) generateDialog() string {
	var builder strings.Builder
	builder.Grow(100)
	builder.WriteString("\"")

	var newWord string
	var word1 = s.dialogueState.GetBigram(state.MagicDialogToken)
	var word2 = s.dialogueState.GetTrigram(state.MagicDialogToken, word1)

	var numTokens = 0
	for numTokens < 24 {
		numTokens = numTokens + 1

		newWord = s.dialogueState.GetTrigram(word1, word2)

		if newWord == state.MagicSentenceToken {
			builder.WriteString(".")
			continue
		}

		if newWord == state.MagicDialogToken {
			break
		}
		builder.WriteString(" ")
		builder.WriteString(newWord)

		word1 = word2
		word2 = newWord

	}

	builder.WriteString("\"")
	return builder.String()
}
