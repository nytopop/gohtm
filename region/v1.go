package region

import (
	"github.com/nytopop/gohtm/cla"
	"github.com/nytopop/gohtm/enc"
	"github.com/nytopop/gohtm/sp"
	"github.com/nytopop/gohtm/tm"
)

// V1Params contains parameters for initialization of a V1 Region.
type V1Params struct {
}

// NewV1Params returns a default set of parameters for a V1 Region.
func NewV1Params() V1Params {
	return V1Params{}
}

// V1 Region. Combines an Encoder, SpatialPooler, TemporalMemory,
// and Classifier for ease of use and composability.
type V1 struct {
	P V1Params
	e enc.Encoder
	s sp.SpatialPooler
	t tm.TemporalMemory
	c cla.Classifier
}

// BuildV1 constructs a region from a JSON encoded region specification.
func BuildV1(rspec string) *V1 {
	return &V1{}
}

// NewV1 returns a new V1 Region initialized with the provided V1Params.
func NewV1(p V1Params) *V1 {
	spar := sp.NewV2Params()
	tpar := tm.NewV1Params()
	cpar := cla.NewV2Params()

	e := enc.NewRDScalar(1024, 102, 0, 0.02)

	spar.NumInputs = 1024
	spar.NumColumns = 2048
	tpar.NumColumns = spar.NumColumns

	s := sp.NewV2(spar)
	t := tm.NewV1(tpar)
	c := cla.NewV2(cpar)

	return &V1{
		P: p,
		e: e,
		s: s,
		t: t,
		c: c,
	}
}

type V1Result struct {
	Datapoint    float64
	AnomalyScore float64
	Prediction   cla.Result
}

func (r *V1) Reset() {
	r.t.Reset()
}

func (r *V1) Compute(datapoint float64, learn bool) V1Result {
	// Encode to vector
	inputvector, bidx := r.e.Encode(datapoint)

	// Compute SP and TM
	activecolumns := r.s.Compute(inputvector, learn)
	r.t.Compute(activecolumns, learn)

	//r.t.Compute(inputvector, learn)

	// Get active cells and run classifier
	activeCells := r.t.GetActiveCells()
	prediction := r.c.Compute(activeCells, bidx, datapoint, true, true)
	anomaly := r.t.GetAnomalyScore()

	return V1Result{
		Datapoint:    datapoint,
		AnomalyScore: anomaly,
		Prediction:   prediction,
	}
}
