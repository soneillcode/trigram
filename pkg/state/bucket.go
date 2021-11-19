package state

import (
	"math/rand"
)

type bucketNgrams struct {
	numberOfBuckets int
	buckets         []*bucket
}

type bucket struct {
	// consider - since only one worker can access a wordFreq, a non locking version could be used instead
	frequencies map[string]*wordFreq
	random      *rand.Rand
	jobChan     chan Job
}

type Job struct {
	key   string
	value string
}

// NewBucketNgrams divides the data into buckets, each with a worker goroutine. As there is a single worker per
// bucket there are no locks. Storage operations are queued and processed one at a time to allow concurrent storage.
// This means that data becomes eventually consistent, as the queue is processed.
// The buckets are selected using a bucket id generation algorithm, the bucket id must be deterministic and as evenly
// distributed across all the buckets as possible.
func NewBucketNgrams(random *rand.Rand, numberOfBuckets int, queueSize int) Ngrams {
	b := &bucketNgrams{
		numberOfBuckets: numberOfBuckets,
		buckets:         make([]*bucket, numberOfBuckets),
	}
	for i := 0; i < numberOfBuckets; i = i + 1 {
		b.buckets[i] = &bucket{
			frequencies: map[string]*wordFreq{},
			random:      random,
			jobChan:     make(chan Job, queueSize),
		}
		go b.buckets[i].run()
	}
	return b
}

func (b *bucketNgrams) Store(words ...string) {
	key, word := getKeyAndWord(words...)
	if key == "" || word == "" {
		return
	}
	bid := getBucketId(key, b.numberOfBuckets)
	b.buckets[bid].store(key, word)
}

func (b *bucketNgrams) Get(words ...string) string {
	key := getKey(words...)
	bid := getBucketId(key, b.numberOfBuckets)
	return b.buckets[bid].get(key)
}

func (b *bucket) store(key string, word string) {
	b.jobChan <- Job{key: key, value: word}
}

func (b *bucket) run() {
	for job := range b.jobChan {
		f := b.frequencies[job.key]
		if f == nil {
			f = newWordFreq()
			b.frequencies[job.key] = f
		}
		f.add(job.value)
	}
}

func (b *bucket) get(key string) string {
	f, has := b.frequencies[key]
	if !has {
		return ""
	}
	return f.get(b.random)
}

// getBucketId returns a bucket id, deterministically based on the key and evenly distributed across all possible buckets.
// consider not using a random function but a process of converting the key into a number.
func getBucketId(key string, numBuckets int) int {
	if numBuckets == 0 {
		return 0
	}
	var keyVal int64
	for _, r := range key {
		keyVal = keyVal + int64(r)
	}
	rand.Seed(keyVal)
	return rand.Intn(numBuckets)
}
