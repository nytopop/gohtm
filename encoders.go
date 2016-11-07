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
	Encode(interface{}) SparseVector
	Decode(SparseVector) interface{}
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

func (rdse RDScalarEncoder) Encode(s interface{}) SparseVector {
	sv := NewSparseVector(rdse.n)

	b := int(s.(float64) / rdse.r)
	if _, ok := rdse.buckets[b]; !ok {
		rdse.NewBucket(b)
	}

	for _, v := range rdse.buckets[b] {
		sv.d = append(sv.d, v)
	}

	return sv
}

func (rdse RDScalarEncoder) NewBucket(b int) {
	for i := 0; i <= b; i++ {
		if _, ok := rdse.buckets[i]; !ok {
			rdse.buckets[i] = jank(i, rdse.n, rdse.w)
			sort.Ints(rdse.buckets[i])
		}
	}
}

func (rdse RDScalarEncoder) Decode(s SparseVector) interface{} {
	for k, v := range rdse.buckets {
		if reflect.DeepEqual(v, s.x) {
			return float64(k) * rdse.r
		}
	}
	return 0.0
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

// ***********************************
