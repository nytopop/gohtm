package enc

import (
	"math"
	"math/rand"

	"github.com/nytopop/gohtm/vec"
)

// RDScalar implements a random distributed scalar encoder.
type RDScalar struct {
	n          uint32  // size of output vector
	w          int     // number of active bits
	r          float64 // resolution
	maxOverlap int
	series     []uint32 // pseudorandom, non-repeating series
}

// NewRDScalar initializes an RDScalar encoder.
func NewRDScalar(n uint32, w, o int, r float64) *RDScalar {
	return &RDScalar{
		n:          n,
		w:          w,
		r:          r,
		maxOverlap: o,
		series:     make([]uint32, 0),
	}
}

// Encode encodes a float value to a bit vector.
func (r *RDScalar) Encode(s interface{}) ([]bool, int) {
	// ensure we get a float64
	if _, ok := s.(float64); !ok {
		panic("RDScalar did not receive float64 input value!")
	}

	// calculate bucket value
	rb := s.(float64) / r.r
	dif := rb - float64(int(rb))

	// round up / down for bucket switchover at .5
	var b int
	switch {
	case dif >= 0.5:
		b = int(math.Ceil(rb))
	case dif < 0.5:
		b = int(math.Floor(rb))
	}

	// create bucket if it doesn't exist
	if len(r.series) < b+r.w {
		r.extendSeries(b)
	}

	// return the bucket
	return vec.ToBool32(r.series[b:b+r.w], r.n), b
}

// extendSeries will recursively create all buckets up to b.
func (r *RDScalar) extendSeries(b int) {
	if b+r.w > len(r.series) {
		r.extendSeries(b - 1)
	}

	idx := len(r.series) - r.w - 2
	if idx < 0 {
		idx = 0
	}

	var newVal uint32
	for {
		newVal = uint32(rand.Intn(int(r.n)))
		if r.isValidBucket(idx, newVal) {
			break
		}
	}

	r.series = append(r.series, newVal)
}

// isValidBucket performs several tests to ensure the validity of a new value
// to use in extending r.series. Values chosen to extend the series are such
// that there is no duplicates with the last w-1 elements, and that new w
// length sequences formed do not overlap with any previous sequences by more
// than r.maxOverlap.
func (r *RDScalar) isValidBucket(b int, newVal uint32) bool {
	// invalidate if newVal is within r.series[b:]
	if vec.Contains32(r.series[b:], newVal) {
		return false
	}

	// validate if no previous sequences or r.maxOverlap == 0
	if len(r.series) < r.w || r.maxOverlap == 0 {
		return true
	}

	bucket := r.series[b : b+r.w-1]
	bucket = append(bucket, newVal)

	start := r.series[0:r.w] // starting overlap calc
	overlap := vec.Overlap32(bucket, start)

	// compute overlap on a sliding window...
	for i := 0; i < b-r.w-1; i++ {
		// decrement overlap if removed idx is in bucket
		if vec.Contains32(bucket, r.series[i]) {
			overlap--
		}

		// increment overlap if next idx is in bucket
		if vec.Contains32(bucket, r.series[i+r.w-1]) {
			overlap++
		}

		// invalidate if overlap is over threshold
		if overlap > r.maxOverlap {
			return false
		}
	}

	// valid if we got this far
	return true
}

func (r *RDScalar) Buckets() int {
	return len(r.series) - r.w + 1
}

func (r *RDScalar) Decode(s []bool) interface{} {
	return 0
}
