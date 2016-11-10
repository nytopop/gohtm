package main

import "fmt"

/* Region
 */

type Region struct {
	enc Encoder
	sp  SpatialPooler
	tm  TemporalMemory
}

func NewRegion() Region {
	sparams := NewSpatialParams()
	return Region{
		enc: NewRDScalarEncoder(400, 21, 1),
		sp:  NewSpatialPooler(sparams),
	}
}

func (r Region) Compute(data interface{}) {
	// encode input and prettyprint
	fmt.Println("Encoding... ", data)
	sv := r.enc.Encode(data)
	fmt.Println(sv.Pretty())

	sv = r.sp.Compute(sv)
}
