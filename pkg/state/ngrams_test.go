package state

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func testEmptyStoreAndGet(name string, newNgramFunc func() Ngrams, t *testing.T) {
	impl := newNgramFunc()
	word := impl.Get("foo")
	if word != "" {
		t.Errorf("for type %s: failed to get value from empty ngram, expected empty string", name)
	}
}

func testBasicBigramStoreAndGet(name string, newNgramFunc func() Ngrams, t *testing.T) {
	impl := newNgramFunc()
	impl.Store("foo", "bar")
	word := impl.Get("foo")
	if word != "bar" {
		t.Errorf("for type %s: failed to get stored bigram", name)
	}
}

func testBasicTrigramStoreAndGet(name string, newNgramFunc func() Ngrams, t *testing.T) {
	impl := newNgramFunc()
	impl.Store("to", "be", "or")
	word := impl.Get("to", "be")
	if word != "or" {
		t.Errorf("for type %s: failed to get stored trigram", name)
	}
}

func testWordFrequency(name string, newNgramFunc func() Ngrams, t *testing.T) {
	impl := newNgramFunc()

	impl.Store("to", "be", "one")
	impl.Store("to", "be", "one")
	impl.Store("to", "be", "one")
	impl.Store("to", "be", "one")
	impl.Store("to", "be", "one")
	impl.Store("to", "be", "one")
	impl.Store("to", "be", "one")
	impl.Store("to", "be", "one")
	impl.Store("to", "be", "one")

	impl.Store("to", "be", "two")

	results := map[string]int{}
	// we use 100 as a rough percentage estimation
	for i := 0; i < 1000; i = i + 1 {
		w := impl.Get("to", "be")
		results[w] = results[w] + 1
	}

	tolerance := 50 // 10%

	if results["one"] < 900-tolerance || results["one"] > 900+tolerance {
		t.Errorf("for type %s: failed to get expected frequency of random words. expected 90 percent of words to be: '%s', but was actually: %v", name, "one", results["one"]/10)
	}
	if results["two"] < 100-tolerance || results["two"] > 100+tolerance {
		t.Errorf("for type %s: failed to get expected frequency of random words. expected 10 percent of words to be: '%s', but was actually: %v", name, "two", results["two"]/10)
	}
}

func testConcurrentAccess(name string, newNgramFunc func() Ngrams, t *testing.T) {
	impl := newNgramFunc()
	mutex := sync.Mutex{}
	waitGroup := sync.WaitGroup{}

	impl.Store("A", "specific", "piece") // we pre-store a value as the read can happen before the write

	results := map[int]string{}
	for i := 0; i < 1000; i = i + 1 {
		waitGroup.Add(2)
		go func() {
			impl.Store("A", "specific", "piece")
			waitGroup.Done()
		}()
		go func(index int) {
			mutex.Lock()
			results[index] = impl.Get("A", "specific")
			mutex.Unlock()
			waitGroup.Done()
		}(i)
	}

	waitGroup.Wait()

	expected := "piece"
	for _, result := range results {
		if result != expected {
			t.Errorf("for type: %s: generate did not generate the expected result: got: '%s' expected: '%s'", name, result, expected)
		}
	}
}

func TestRandomFreqOneThird(t *testing.T) {
	wf := newWordFreq()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	wf.add("one")
	wf.add("one")

	wf.add("two")

	results := map[string]int{}
	// we use 100 as a rough percentage estimation
	for i := 0; i < 100; i = i + 1 {
		word := wf.get(random)
		results[word] = results[word] + 1
	}

	tolerance := 5 // 10%

	if results["one"] < 66-tolerance || results["one"] > 66+tolerance {
		t.Errorf("failed to get expected frequency of random words. expected 66 percent of words to be: '%s', but was actually: %v", "one", results["one"])
	}

	if results["two"] < 33-tolerance || results["two"] > 33+tolerance {
		t.Errorf("failed to get expected frequency of random words. expected 33 percent of words to be: '%s', but was actually: %v", "two", results["two"])
	}
}

func TestRandomFreqOneHalf(t *testing.T) {
	wf := newWordFreq()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	wf.add("one")
	wf.add("one")

	wf.add("two")
	wf.add("two")

	results := map[string]int{}
	// we use 100 as a rough percentage estimation
	for i := 0; i < 1000; i = i + 1 {
		word := wf.get(random)
		results[word] = results[word] + 1
	}

	tolerance := 50 // 10%

	if results["two"] < 500-tolerance || results["two"] > 500+tolerance {
		t.Errorf("failed to get expected frequency of random words. expected 33 percent of words to be: '%s', but was actually: %v", "two", results["two"]/10)
	}

	if results["one"] < 500-tolerance || results["one"] > 500+tolerance {
		t.Errorf("failed to get expected frequency of random words. expected 33 percent of words to be: '%s', but was actually: %v", "one", results["one"]/10)
	}
}

func TestRandomFreqOneTenth(t *testing.T) {
	wf := newWordFreq()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	wf.add("one")
	wf.add("one")
	wf.add("one")
	wf.add("one")
	wf.add("one")
	wf.add("one")
	wf.add("one")
	wf.add("one")
	wf.add("one")

	wf.add("two")

	results := map[string]int{}
	// we use 100 as a rough percentage estimation
	for i := 0; i < 1000; i = i + 1 {
		word := wf.get(random)
		results[word] = results[word] + 1
	}

	tolerance := 50 // 10%

	if results["two"] < 100-tolerance || results["two"] > 100+tolerance {
		t.Errorf("failed to get expected frequency of random words. expected 10 percent of words to be: '%s', but was actually: %v", "two", results["two"]/10)
	}

	if results["one"] < 900-tolerance || results["one"] > 900+tolerance {
		t.Errorf("failed to get expected frequency of random words. expected 90 percent of words to be: '%s', but was actually: %v", "one", results["one"]/10)
	}
}

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
