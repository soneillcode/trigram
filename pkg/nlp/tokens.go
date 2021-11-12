package nlp

import "example.com/todo/pkg/state"

var titles = []string{"Mr", "Mrs", "Dr", "Ms"}
var ignoreCharacters = []rune{',', ';', ':', '_', '"', '“', '”', '“'}

const periodCharacter = '.'
const spaceCharacter = ' '
const newLineCharacter = '\n'

const fullStopWord = "."
const spaceWord = " "

// getTokens iterates through a body of text and creates tokens delineated by space and newline characters.
// tokens are defined as strings of non-space characters. It also adds magic tokens to mark the start and end of sentences.
// Sentences are handled as a special case to improve the quality of text generation.
func getTokens(text string) []string {
	var tokens []string
	var word []rune

	tokens = append(tokens, state.MagicStartToken)
	for _, character := range text {

		if character == spaceCharacter || character == newLineCharacter {
			if len(word) > 0 {
				tokens = append(tokens, string(word))
				word = word[:0]
			}
			continue
		}

		// filter out some characters
		if shouldIgnore(character) {
			continue
		}

		// add full stops as distinct tokens
		if character == periodCharacter {
			if len(word) > 0 {
				w := string(word)
				if isTitle(w) {
					word = append(word, character)
					if len(word) > 0 {
						tokens = append(tokens, string(word))
						word = word[:0]
					}
				} else {
					if len(word) > 0 {
						tokens = append(tokens, string(word))
						word = word[:0]
					}
					tokens = append(tokens, fullStopWord)
					word = word[:0]
					tokens = append(tokens, state.MagicStartToken)
				}
				continue
			}
			tokens = append(tokens, fullStopWord)
			continue
		}

		// handle standard character
		word = append(word, character)
	}
	return tokens
}

func shouldIgnore(character rune) bool {
	for _, ignore := range ignoreCharacters {
		if character == ignore {
			return true
		}
	}
	return false
}

func isTitle(word string) bool {
	for _, title := range titles {
		if title == word {
			return true
		}
	}
	return false
}

// storeTokens puts the given tokens in the ngram data structure.
func storeTokens(tokens []string, ngrams *state.Ngrams) {
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
