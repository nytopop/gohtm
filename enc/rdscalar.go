package enc

import (
	"math"
	"math/rand"

	"github.com/nytopop/gohtm/vec"
)

// TODO : cleanup + Decode()

// RDScalar implements a random distributed scalar encoder.
type RDScalar struct {
	N          uint32   `json:"n"`
	W          int      `json:"w"`
	R          float64  `json:"r"`
	MaxOverlap int      `json:"maxOverlap"`
	Series     []uint32 `json:"series"`
}

// NewRDScalar initializes an RDScalar encoder.
func NewRDScalar(n uint32, w, o int, r float64) *RDScalar {
	return &RDScalar{
		N:          n,
		W:          w,
		R:          r,
		MaxOverlap: o,
		Series:     make([]uint32, 0),
	}
}

// Encode encodes a float value to a bit vector.
func (r *RDScalar) Encode(s interface{}) ([]bool, int) {
	// ensure we get a float64
	if _, ok := s.(float64); !ok {
		panic("RDScalar did not receive float64 input value!")
	}

	// calculate bucket value
	rb := s.(float64) / r.R
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
	if len(r.Series) < b+r.W {
		r.extendSeries(b)
	}

	// return the bucket
	return vec.ToBool32(r.Series[b:b+r.W], r.N), b
}

// extendSeries will recursively create all buckets up to b.
func (r *RDScalar) extendSeries(b int) {
	if b+r.W > len(r.Series) {
		r.extendSeries(b - 1)
	}

	idx := len(r.Series) - r.W - 2
	if idx < 0 {
		idx = 0
	}

	var newVal uint32
	for {
		newVal = uint32(rand.Intn(int(r.N)))
		if r.isValidBucket(idx, newVal) {
			break
		}
	}

	r.Series = append(r.Series, newVal)
}

// isValidBucket performs several tests to ensure the
// validity of a new value to use in extending r.Series.
// Values chosen to extend the series are such that
// there are no duplicates with the last w-1 elements,
// and that new w length sequences formed do not overlap
// with any previous sequences by more than r.MaxOverlap.
func (r *RDScalar) isValidBucket(b int, newVal uint32) bool {
	// invalidate if newVal is within r.Series[b:]
	if vec.Contains32(r.Series[b:], newVal) {
		return false
	}

	// validate if no previous sequences or r.MaxOverlap == 0
	if len(r.Series) < r.W || r.MaxOverlap == 0 {
		return true
	}

	bucket := r.Series[b : b+r.W-1]
	bucket = append(bucket, newVal)

	start := r.Series[0:r.W] // starting overlap calc
	overlap := vec.Overlap32(bucket, start)

	// compute overlap on a sliding window...
	for i := 0; i < b-r.W-1; i++ {
		// decrement overlap if removed idx is in bucket
		if vec.Contains32(bucket, r.Series[i]) {
			overlap--
		}

		// increment overlap if next idx is in bucket
		if vec.Contains32(bucket, r.Series[i+r.W-1]) {
			overlap++
		}

		// invalidate if overlap is over threshold
		if overlap > r.MaxOverlap {
			return false
		}
	}

	// valid if we got this far
	return true
}

func (r *RDScalar) Buckets() int {
	return len(r.Series) - r.W + 1
}

func (r *RDScalar) Decode(s []bool) interface{} {
	return 0
}
