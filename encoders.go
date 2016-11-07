package main

import (
	"math/rand"
	"reflect"
	"sort"
)

/* Encoder Design Guidelines
1. Semantically similar data should result in SDRs with overlapping active bits.
2. The same input should always produce the same SDR as output.
3. The output should have the same dimensionality (total number of bits) for all inputs.
4. The output should have similar sparsity for all inputs and have enough one-bits to handle noise and subsampling.
*/

type Encoder interface {
	Encode(interface{}) SDR
	Decode(SDR) interface{}
}

// Random Distributed Scalar Encoder
// **********************************
type RDScalarEncoder struct {
	n       int
	w       int
	r       float64
	buckets map[int][]int
}

func NewRDScalarEncoder(n, w int, r float64) RDScalarEncoder {
	return RDScalarEncoder{
		n:       n,
		w:       w,
		r:       r,
		buckets: map[int][]int{},
	}
}

func (se RDScalarEncoder) Encode(s interface{}) SDR {
	vec := Vector{}
	for i := 0; i < se.n; i++ {
		vec = append(vec, false)
	}

	b := int(s.(float64) / se.r)
	if _, ok := se.buckets[b]; !ok {
		se.NewBucket(b)
	}

	for _, v := range se.buckets[b] {
		vec[v] = true
	}

	return vec.SDR()
}

func (se RDScalarEncoder) NewBucket(b int) {
	for i := 0; i <= b; i++ {
		if _, ok := se.buckets[i]; !ok {
			se.buckets[i] = jank(i, se.n, se.w)
			sort.Ints(se.buckets[i])
		}
	}
}

// Literally the worst algorithm in the known universe.
func jank(n, max, w int) []int {
	out := []int{}
	// loop until successful
	for seed := 0; len(out) < w; seed++ {
		out = []int{}

		for i := 0; i < w; i++ {
			t := fakeHash(n+i, max, seed)
			if !intInSlice(t, out) {
				out = append(out, t)
			}
		}
	}

	return out
}

func intInSlice(n int, s []int) bool {
	for _, v := range s {
		if n == v {
			return true
		}
	}
	return false
}

func fakeHash(n, max, seed int) int {
	rand.Seed(int64(seed))
	for i := 0; i < n; i++ {
		rand.Intn(max)
	}
	return rand.Intn(max)
}

func (se RDScalarEncoder) Decode(s SDR) interface{} {
	for k, v := range se.buckets {
		if reflect.DeepEqual(v, s.w) {
			return float64(k) * se.r
		}
	}

	return 0.0
}

// ***********************************
