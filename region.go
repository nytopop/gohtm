package main

import "fmt"

/* Region Parameters */
type RegionParams struct {
	//ep EncoderParams
	sp SpatialParams
	tp TemporalParams
}

/* Result struct for a region's output. */
type RegionResult struct {
	data    interface{}
	encoded SparseBinaryVector
	spatial SparseBinaryVector
}

/* Region. A type comprising an encoder, spatial pooler, and temporal memory segment. */
type Region struct {
	enc Encoder
	sp  SpatialPooler
	tm  TemporalMemory

	iteration int
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
func (r *Region) Compute(data interface{}, learn bool) RegionResult {
	r.iteration++

	// encode input and prettyprint
	sv := r.enc.Encode(data)
	rv := r.sp.Compute(sv, learn)
	fmt.Println("Data:", data)
	fmt.Println("Sparsity:", rv.Sparsity())
	fmt.Println(rv.Pretty())

	//r.sp.Save("sp.json")

	return RegionResult{
		data:    data,
		encoded: sv,
		spatial: rv,
	}
}

func (r *Region) PredictK(k int) {
	for i := 0; i < k; i++ {
		// call compute recursively
	}
}
