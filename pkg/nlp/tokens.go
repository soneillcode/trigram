package nlp

import "example.com/todo/pkg/state"

var titles = []string{"Mr", "Mrs", "Dr", "Ms"}
var ignoreCharacters = []rune{',', ';', ':', '_'}
var dialogCharacters = []rune{'"', '“', '”', '“'}

const periodCharacter = '.'
const spaceCharacter = ' '
const newLineCharacter = '\n'

const fullStopWord = "."
const spaceWord = " "
const dialogQuoteWord = "\""

// tokenize splits a body of text by space and newline characters. It also adds magic tokens to mark the start and end
// of sentences and dialog.
func tokenize(text string) []string {
	var tokens []string
	var word []rune

	addCurrentWordToTokens := func() {
		if len(word) > 0 {
			tokens = append(tokens, string(word))
			// note: the intent is to use slices to empty the array instead of allocating more memory
			// not sure this is correct
			word = word[:0]
		}
	}

	tokens = append(tokens, state.MagicStartToken)
	for _, character := range text {

		if character == spaceCharacter || character == newLineCharacter {
			addCurrentWordToTokens()
			continue
		}

		// filter out some characters
		if shouldIgnore(character) {
			continue
		}

		// handle dialog characters
		if isDialog(character) {
			addCurrentWordToTokens()
			tokens = append(tokens, state.MagicDialogToken)
			continue
		}

		// add full stops as distinct tokens
		if character == periodCharacter {
			if len(word) > 0 {
				w := string(word)
				if isTitle(w) {
					word = append(word, character)
					addCurrentWordToTokens()
				} else {
					addCurrentWordToTokens()
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

func isDialog(character rune) bool {
	for _, dialog := range dialogCharacters {
		if character == dialog {
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
