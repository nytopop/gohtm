package region

import (
	"fmt"

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

// NewV1 returns a new V1 Region initialized with the provided V1Params.
func NewV1(p V1Params) *V1 {
	epar := enc.NewScalarParams()
	spar := sp.NewV1Params()
	tpar := tm.NewV1Params()
	cpar := cla.NewV2Params()

	e := enc.NewScalar(epar)

	spar.NumInputs = e.Bits
	fmt.Println(e.Bits)
	tpar.NumColumns = spar.NumColumns

	s := sp.NewV1(spar)
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
	Datapoint    int
	AnomalyScore float64
	Prediction   cla.V2Results
}

func (r *V1) Compute(datapoint int, learn bool) V1Result {
	// Encode to vector
	inputvector := r.e.Encode(datapoint)

	// Compute SP and TM
	activecolumns := r.s.Compute(inputvector, learn)
	r.t.Compute(activecolumns, learn)

	// Get active cells and run classifier
	activeCells := r.t.GetActiveCells()
	prediction := r.c.Compute(activeCells, datapoint, true, true)
	anomaly := r.t.GetAnomalyScore()

	return V1Result{
		Datapoint:    datapoint,
		AnomalyScore: anomaly,
		Prediction:   prediction,
	}
}
