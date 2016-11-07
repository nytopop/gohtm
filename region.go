package main

/* Region
 */

type Region struct {
	enc Encoder
	sp  SpatialPooler
	tm  TemporalMemory
}

func NewRegion() Region {
	return Region{
		enc: NewRDScalarEncoder(64, 4, 1),
	}
}

func (r Region) Compute() {
}
