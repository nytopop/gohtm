package gohtm

import "fmt"

// RegionParams is a meta-struct containing parameters for a SpatialPooler
// and a TemporalMemory instance.
type RegionParams struct {
	sp SpatialParams
	tp TemporalParams
}

// RegionResult contains the output of calling Compute on a Region.
type RegionResult struct {
	data     interface{}
	encoded  SparseBinaryVector
	spatial  SparseBinaryVector
	temporal SparseBinaryVector
}

// Region wraps an Encoder, SpatialPooler, and TemporalMemory instance
// into one object for ease of use.
type Region struct {
	enc Encoder
	sp  SpatialPooler
	tm  TemporalMemory

	iteration int
}

// NewRegion returns a new region.
func NewRegion() Region {
	return Region{
		enc: NewRDScalarEncoder(400, 21, 1),
		sp:  NewSpatialPooler(NewSpatialParams()),
		tm:  NewTemporalMemory(NewTemporalParams()),
	}
}

// Compute encodes a provided datapoint, calls Compute on the
// SpatialPooler and TemporalMemory, and returns the result.
func (r *Region) Compute(data interface{}, learn bool) RegionResult {
	r.iteration++

	// encode input and prettyprint
	sv := r.enc.Encode(data)
	rv := r.sp.Compute(sv, learn)
	tv := r.tm.Compute(rv, learn)

	// Encoder
	fmt.Println("Data:", data)

	// Spatial Pooling
	fmt.Println("SP Sparsity:", rv.Sparsity())
	fmt.Println(rv.Pretty())

	// Temporal Memory
	fmt.Println("TM Sparsity:", tv.Sparsity())
	fmt.Println(tv.Pretty())

	//r.sp.Save("sp.json")

	return RegionResult{
		data:    data,
		encoded: sv,
		spatial: rv,
		//	temporal: tv,
	}
}

// PredictK recursively predicts k steps into the future. No learning
// is performed.
func (r *Region) PredictK(k int) {
	for i := 0; i < k; i++ {
		// call compute recursively
	}
}
