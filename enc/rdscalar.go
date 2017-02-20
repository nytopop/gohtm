package enc

import (
	"math"
	"math/rand"

	"github.com/nytopop/gohtm/vec"
)

// RDScalar implements a random distributed scalar encoder.
type RDScalar struct {
	n      uint32   // size of output vector
	w      int      // number of active bits
	r      float64  // resolution
	series []uint32 // pseudorandom, non-repeating series
}

// NewRDScalar initializes an RDScalar encoder.
func NewRDScalar(n uint32, w int, r float64) *RDScalar {
	return &RDScalar{
		n: n,
		w: w,
		r: r,
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
		subset := r.series[idx:]

		// find random value not in subset
		var rVal int
		for {
			rVal = rand.Intn(int(r.n))
			if !vec.Contains32(subset, uint32(rVal)) {
				break
			}
		}

		r.series = append(r.series, uint32(rVal))
	}
}

func (r *RDScalar) Buckets() int {
	return len(r.series) - r.w
}
