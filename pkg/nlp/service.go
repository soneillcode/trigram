package nlp

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"example.com/todo/pkg/nlp/state"
)

/*
	The nlp (Natural Language Processing) service.
*/

type Service struct {
	defaultLength string
	state         *state.State
}

func NewService() *Service {
	return &Service{
		defaultLength: "120",
		state:         state.NewState(),
	}
}

func (s *Service) Learn(text string) error {
	if text == "" {
		return nil
	}

	// consider adding tokens to state as they are created to save storing them all.
	tokens := tokenize(text)

	for i, t := range tokens {
		if i == 1 {
			s.state.Store("", tokens[i-1], t)
		}
		if i > 1 {
			s.state.Store(tokens[i-2], tokens[i-1], t)
		}
	}

	return nil
}

func tokenize(text string) []string {
	var tokens []string
	var isWord = false
	var word []rune
	for _, r := range text {
		if r != ' ' {
			// todo filter tokens list and common handling
			// filter end of line
			if r == '\n' {
				continue
			}
			// filter comma
			if r == ',' {
				continue
			}
			// filter semi colon
			if r == ';' {
				continue
			}
			// filter underscore
			if r == '_' {
				continue
			}
			// add full stops as a distinct token
			if r == '.' {
				if len(word) > 0 {
					tokens = append(tokens, string(word))
					word = word[:0]
				}
				tokens = append(tokens, ".")
				continue
			}
			// add quotes as a distinct token
			// todo distinct tokens list and common handling
			if r == '"' {
				if len(word) > 0 {
					tokens = append(tokens, string(word))
					word = word[:0]
				}
				tokens = append(tokens, "\"")
				continue
			}
			if r == '“' {
				if len(word) > 0 {
					tokens = append(tokens, string(word))
					word = word[:0]
				}
				tokens = append(tokens, "“")
				continue
			}
			if r == '”' {
				if len(word) > 0 {
					tokens = append(tokens, string(word))
					word = word[:0]
				}
				tokens = append(tokens, "”")
				continue
			}
			if r == '?' {
				if len(word) > 0 {
					tokens = append(tokens, string(word))
					word = word[:0]
				}
				tokens = append(tokens, "”")
				continue
			}
			// handle standard alphanumeric character
			isWord = true
			word = append(word, r)
		}
		if r == ' ' {
			if isWord {
				tokens = append(tokens, string(word))
				word = word[:0]
			}
			isWord = false
		}

	}
	return tokens
}

func (s *Service) Generate() (*string, error) {
	lengthVal, err := strconv.Atoi(s.defaultLength)
	if err != nil {
		return nil, fmt.Errorf("failed to convert length to int: %w", err)
	}
	log.Printf("length: %v", lengthVal)

	var tokens []string
	var word1 = ""
	var word2 = ""
	var word3 = s.state.Get(word1, word2)

	for numSentences := 0; numSentences < 12; {
		word3 = s.state.Get(word1, word2)
		word1 = word2
		word2 = word3
		tokens = append(tokens, word3)
		if word3 == "." {
			numSentences = numSentences + 1
		}
	}
	text := strings.Join(tokens, " ")
	return &text, nil
}
