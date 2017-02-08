// Package region provides an implementation agnostic interface for linking
// and interacting with computational units of the htm algorithm.
package region

import (
	"github.com/nytopop/gohtm/sp"
	"github.com/nytopop/gohtm/tm"
)

// Region is an interface providing access to a segment of cortex.
type Region interface {
	Compute(input []bool)
	Serialize() []byte
}

/*

Spatial
  vv
Temporal
  vv
Temporal
  vv
Sensori-Motor
  vv
Motor
  vv
Feedback

*/

type V1 struct {
	Pooler sp.SpatialPooler
	L1     int //feedback
	L2, L3 tm.TemporalMemory
	L6     int // feedback pathway
}

func (r *V1) Compute(input []bool) {
}

func (r *V1) LoadFeedback() {
}

func (r *V1) Serialize() []byte {
	return []byte{}
}
