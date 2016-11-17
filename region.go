package main

import "fmt"

/* Region. A type comprising an encoder, spatial pooler, and temporal memory segment. */
type Region struct {
	enc Encoder
	sp  SpatialPooler
	tm  TemporalMemory
}

/* Create a new region with default parameters. */
func NewRegion() Region {
	return Region{
		enc: NewRDScalarEncoder(400, 21, 1),
		sp:  NewSpatialPooler(NewSpatialParams()),
		tm:  NewTemporalMemory(NewTemporalParams()),
	}
}

/* Encode a datapoint, call compute on temporal memory and spatial pooler. */
func (r *Region) Compute(data interface{}) {
	// encode input and prettyprint
	fmt.Println("Encoding:", data)
	sv := r.enc.Encode(data)

	rv := r.sp.Compute(sv, true)
	fmt.Println("Sparsity:", rv.Sparsity())
	fmt.Println(rv.Pretty())
}
