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
func (r *RDScalar) Encode(s interface{}) []bool {
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
	return vec.ToBool32(r.series[b:b+r.w], r.n)
}

func (r *RDScalar) extendSeries(b int) {
	// calc # of new entries
	n := (b + r.w) - len(r.series)
	for i := 0; i < n; i++ {
		// get subset of last w+1 elems in series
		idx := len(r.series) - r.w - 2
		if idx < 0 {
			idx = 0
		}

		// find random value not in subset
		var rVal uint32
		for {
			// get random number from [0 : r.n]
			rVal = uint32(rand.Intn(int(r.n)))

			// ensure number isn't a duplicate from subset
			if !vec.Contains32(r.series[idx:], rVal) {
				// ensure non-adjescent buckets don't overlap
				if r.isValidBucket(idx, rVal) {
					break
				}
			}
		}
		r.series = append(r.series, rVal)
	}
}

func (r *RDScalar) isValidBucket(b int, newVal uint32) bool {
	if len(r.series) < r.w {
		return true
	}

	if r.maxOverlap == 0 {
		return true
	}

	bucket := r.series[b : b+r.w-1]
	bucket = append(bucket, newVal)

	var overlap int
	start := r.series[0:r.w] // starting overlap calc
	overlap = vec.Overlap32(bucket, start)

	// compute overlap on a sliding window...
	var cursor uint32
	for i := 0; i < b-r.w; i++ {
		// decrement overlap if the idx we are removing was in bucket
		cursor = r.series[i]
		if vec.Contains32(bucket, cursor) {
			overlap--
		}

		// increment overlap if the next idx is in bucket
		cursor = r.series[i+r.w-1]
		if vec.Contains32(bucket, cursor) {
			overlap++
		}

		if overlap > r.maxOverlap {
			return false
		}
	}

	return true
}

func (r *RDScalar) Buckets() int {
	return len(r.series) - r.w + 1
}
