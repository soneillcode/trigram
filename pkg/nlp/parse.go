package nlp

import "example.com/todo/pkg/state"

var Titles = []string{"Mr", "Mrs", "Dr", "Ms"}

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
					w := string(word)
					if isTitle(w) {
						tokens = append(tokens, w+".")
					} else {
						tokens = append(tokens, w)
						tokens = append(tokens, state.MagicSentenceToken)
					}
					word = word[:0]
					continue
				}
				tokens = append(tokens, state.MagicSentenceToken)
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
			if isWord && len(word) > 0 {
				tokens = append(tokens, string(word))
				word = word[:0]
			}
			isWord = false
		}

	}
	return tokens
}

func isTitle(word string) bool {
	for _, t := range Titles {
		if t == word {
			word = word[:0]
			return true
		}
	}
	return false
}
