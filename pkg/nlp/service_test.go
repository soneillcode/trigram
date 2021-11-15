package nlp

import (
	"math/rand"
	"sync"
	"testing"
)

func TestServiceEmpty(t *testing.T) {
	// use a static seed to generate 'random' numbers deterministically
	random := rand.New(rand.NewSource(16))
	service := NewService(random, 6)

	result := service.Generate()
	expected := "       "
	if result != expected {
		t.Errorf("generate did not generate the expected result: got: '%s' expected: '%s'", result, expected)
	}
}

func TestServiceSingleLine(t *testing.T) {
	// use a static seed to generate 'random' numbers deterministically
	random := rand.New(rand.NewSource(16))
	service := NewService(random, 6)
	service.Learn("A specific piece of unique text.")

	result := service.Generate()
	expected := "A specific piece of unique text."
	if result != expected {
		t.Errorf("generate did not generate the expected result: got: '%s' expected: '%s'", result, expected)
	}
}

func TestServiceMultipleSentences(t *testing.T) {
	// use a static seed to generate 'random' numbers deterministically
	random := rand.New(rand.NewSource(16))
	service := NewService(random, 21)
	service.Learn("A specific piece of unique text.")

	result := service.Generate()
	expected := "A specific piece of unique text. A specific piece of unique text. A specific piece of unique text."
	if result != expected {
		t.Errorf("generate did not generate the expected result: got: '%s' expected: '%s'", result, expected)
	}
}

func TestServiceConcurrentAccess(t *testing.T) {
	// use a static seed to generate 'random' numbers deterministically
	random := rand.New(rand.NewSource(16))
	service := NewService(random, 6)
	service.Learn("A specific piece of unique text.")

	mutex := sync.Mutex{}
	waitGroup := sync.WaitGroup{}

	results := map[int]string{}
	for i := 0; i < 1000; i = i + 1 {
		waitGroup.Add(2)
		go func() {
			service.Learn("A specific piece of unique text.")
			waitGroup.Done()
		}()
		go func(index int) {
			mutex.Lock()
			results[index] = service.Generate()
			mutex.Unlock()
			waitGroup.Done()
		}(i)
	}

	waitGroup.Wait()

	expected := "A specific piece of unique text."
	for _, result := range results {
		if result != expected {
			t.Errorf("generate did not generate the expected result: got: '%s' expected: '%s'", result, expected)
		}
	}
}
