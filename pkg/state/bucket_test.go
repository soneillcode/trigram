package state

import (
	"fmt"
	"math/rand"
	"testing"
)

func Test_BucketNgrams(t *testing.T) {
	random := rand.New(rand.NewSource(16))
	implFunc := func() Ngrams {
		return NewBucketNgrams(random, 12, 12)
	}
	name := "bucketNgrams"
	testEmptyStoreAndGet(name, implFunc, t)
	testBasicBigramStoreAndGet(name, implFunc, t)
	testBasicTrigramStoreAndGet(name, implFunc, t)
	testWordFrequency(name, implFunc, t)
	testConcurrentAccess(name, implFunc, t)
}

func Test_getBucketId(t *testing.T) {
	tests := []struct {
		name string
		key  string
		size int
		want int
	}{
		{
			name: "empty",
			key:  "",
			size: 1,
			want: 0,
		},
		{
			name: "single character, one bucket",
			key:  "a",
			size: 1,
			want: 0,
		},
		{
			name: "single character, two buckets",
			key:  "b",
			size: 2,
			want: 0,
		},
		{
			name: "single character, higher value, more buckets",
			key:  "dfg876sdfgkdgiuhxg",
			size: 12,
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBucketId(tt.key, tt.size); got != tt.want {
				t.Errorf("getBucketId() = %v, want %v", got, tt.want)
			}
		})
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Test_getBucketIdSafety(t *testing.T) {
	var values []string

	for i := 0; i < 1024; i = i + 1 {
		values = append(values, randSeq(16))
	}

	var bucketIds []int
	for _, v := range values {
		bucketIds = append(bucketIds, getBucketId(v, 12))
	}

	for index, id := range bucketIds {
		if id < 0 || id > 11 {
			t.Errorf("getBucketId(): index:%v id:%v out of bucket range ( 0 - 11 )", index, id)
		}
	}

}

func Test_getBucketIdDistribution(t *testing.T) {
	var values []string

	for i := 0; i < 1024; i = i + 1 {
		length := rand.Intn(32)
		values = append(values, randSeq(length))
	}

	bucketIds := map[int]int{}
	for _, v := range values {
		id := getBucketId(v, 12)
		freq, ok := bucketIds[id]
		if !ok {
			bucketIds[id] = 1
		} else {
			bucketIds[id] = freq + 1
		}
	}

	for k, val := range bucketIds {
		fmt.Printf("key: %v value: %v\n", k, val)
		if val < 1 || val > 2*(1024/12) {
			t.Errorf("getBucketId(): id:%v val:%v not correctly distributed", k, val)
		}
	}
	//t.Errorf("getBucketId(): index:%v id:%v out of bucket range ( 0 - 11 )", index, id)

}
