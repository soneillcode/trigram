package state

import (
	"fmt"
	"testing"
)

func Test_getKey(t *testing.T) {
	separator := keySeparator
	tests := []struct {
		desc  string
		words []string
		want  string
	}{
		{
			desc:  "empty case",
			words: []string{},
			want:  "",
		},
		{
			desc:  "single case",
			words: []string{"foo"},
			want:  "foo",
		},
		{
			desc:  "two word case",
			words: []string{"foo", "bar"},
			want:  fmt.Sprintf("foo%sbar", separator),
		},
		{
			desc:  "five word case",
			words: []string{"foo", "bar", "to", "be", "not"},
			want:  fmt.Sprintf("foo%sbar%sto%sbe%snot", separator, separator, separator, separator),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got := getKey(tt.words...); got != tt.want {
				t.Errorf("getKey(): got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_getKeyAndWord(t *testing.T) {
	separator := keySeparator
	tests := []struct {
		desc  string
		words []string
		key   string
		word  string
	}{
		{
			desc:  "empty case",
			words: nil,
			key:   "",
			word:  "",
		},
		{
			desc:  "single case",
			words: []string{"foo"},
			key:   "foo",
			word:  "",
		},
		{
			desc:  "two word case",
			words: []string{"foo", "bar"},
			key:   "foo",
			word:  "bar",
		},
		{
			desc:  "three word case",
			words: []string{"foo", "bar", "to"},
			key:   fmt.Sprintf("foo%sbar", separator),
			word:  "to",
		},
		{
			desc:  "four word case",
			words: []string{"foo", "bar", "to", "be"},
			key:   fmt.Sprintf("foo%sbar%sto", separator, separator),
			word:  "be",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			key, word := getKeyAndWord(tt.words...)
			if key != tt.key || word != tt.word {
				t.Errorf("getKeyAndWord(): got: %v,%v, want: %v,%v", key, word, tt.key, tt.word)
			}
		})
	}
}
