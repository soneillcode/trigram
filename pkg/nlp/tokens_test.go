package nlp

import (
	"reflect"
	"testing"
)

func Test_getTokens(t *testing.T) {

	tests := []struct {
		desc string
		text string
		want []string
	}{
		{
			desc: "empty case",
			text: "",
			want: []string{MagicStartToken},
		},
		{
			desc: "single case",
			text: "one",
			want: []string{MagicStartToken, "one"},
		},
		{
			desc: "two word case",
			text: "To be",
			want: []string{MagicStartToken, "To", "be"},
		},
		{
			desc: "six word case",
			text: "To be or not to be",
			want: []string{MagicStartToken, "To", "be", "or", "not", "to", "be"},
		},
		{
			desc: "full stop case",
			text: "To be. Or not",
			want: []string{MagicStartToken, "To", "be", ".", MagicStartToken, "Or", "not"},
		},
		{
			desc: "multiple full stop case",
			text: "To be. Or not. To be",
			want: []string{MagicStartToken, "To", "be", ".", MagicStartToken, "Or", "not", ".", MagicStartToken, "To", "be"},
		},
		{
			desc: "treat newlines as spaces case",
			text: "To be\nor not",
			want: []string{MagicStartToken, "To", "be", "or", "not"},
		},
		{
			desc: "ignore comma case",
			text: "To be, or not",
			want: []string{MagicStartToken, "To", "be", "or", "not"},
		},
		{
			desc: "ignore semi-colon case",
			text: "To be; or not",
			want: []string{MagicStartToken, "To", "be", "or", "not"},
		},
		{
			desc: "ignore colon case",
			text: "To be: or not",
			want: []string{MagicStartToken, "To", "be", "or", "not"},
		},
		{
			desc: "ignore underscore case",
			text: "To _be_ or not",
			want: []string{MagicStartToken, "To", "be", "or", "not"},
		},
		{
			desc: "don't treat title as full stop case",
			text: "Hello Mr. Bingley. To.",
			want: []string{MagicStartToken, "Hello", "Mr.", "Bingley", ".", MagicStartToken, "To", ".", MagicStartToken},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := getTokens(tt.text)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTokens(): got: '%v', want: '%v'", got, tt.want)
			}
		})
	}

}
