package gohtm

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

// Encoder is an interface for all sparse encoders. A valid Encoder
// should implement a Decode and an Encode method, persistent
// state is not required.
type Encoder interface {
	Encode(interface{}) SparseBinaryVector
	Decode(SparseBinaryVector) interface{}
}

// RDScalarEncoder is an implementation of a random distributed scalar
// encoder. Scalar values are mapped to buckets, which are randomly
// assigned to groups of bits in the output space. Requires persistent
// state.
type RDScalarEncoder struct {
	n       int
	w       int
	r       float64
	buckets map[int][]int
}

// NewRDScalarEncoder returns a new RDScalarEncoder with the supplied
// parameters. n denotes the size of the output vector, w denotes how
// many bits should be active for an encoded scalar value. r denotes
// the resolution of the encoder; setting r to 1 will produce a bucket
// for every interval of 1 in the input space.
func NewRDScalarEncoder(n, w int, r float64) *RDScalarEncoder {
	return &RDScalarEncoder{
		n:       n,
		w:       w,
		r:       r,
		buckets: map[int][]int{},
	}
}

// Encode encodes a scalar value and returns a SparseBinaryVector.
// The provided scalar value must be float64.
func (rdse *RDScalarEncoder) Encode(s interface{}) SparseBinaryVector {
	b := int(s.(float64) / rdse.r)
	if _, ok := rdse.buckets[b]; !ok {
		rdse.newBucket(b)
	}

	sv := NewSparseBinaryVector(rdse.n)
	for _, v := range rdse.buckets[b] {
		sv.Set(v, true)
	}

	return sv
}

func (rdse *RDScalarEncoder) newBucket(b int) {
	for i := 0; i <= b; i++ {
		if _, ok := rdse.buckets[i]; !ok {
			rdse.buckets[i] = jank(i, rdse.n, rdse.w)
			sort.Ints(rdse.buckets[i])
		}
	}
}

// Decode decodes the provided SparseBinaryVector back into a scalar
// value. If decoding fails, 0.0 is returned.
func (rdse *RDScalarEncoder) Decode(s SparseBinaryVector) interface{} {
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
			if !containsInt(t, out) {
				out = append(out, t)
			}
		}
	}

	return out
}

func fakeHash(n, max, seed int) int {
	rand.Seed(int64(seed))
	for i := 0; i < n; i++ {
		rand.Intn(max)
	}
	return rand.Intn(max)
}
