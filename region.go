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
	fmt.Println("Initializing Region...")
	sparams := NewSpatialParams()
	return Region{
		enc: NewRDScalarEncoder(400, 21, 1),
		sp:  NewSpatialPooler(sparams),
	}
}

/* Encode a datapoint, call compute on temporal memory and spatial pooler. */
func (r *Region) Compute(data interface{}) {
	// encode input and prettyprint
	fmt.Println("Encoding...", data)
	sv := r.enc.Encode(data)
	//fmt.Println(sv.Pretty())

	fmt.Println("Pooling...")
	rv := r.sp.Compute(sv, true)
	fmt.Println("Sparsity:", rv.Sparsity())
	fmt.Println(rv.Pretty())

	fmt.Println()
}
