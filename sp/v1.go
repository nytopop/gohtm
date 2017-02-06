package sp

import (
	"math/rand"

	"github.com/nytopop/gohtm/vec"
)

// V1Params contains parameters for initialization of a V1 SpatialPooler.
type V1Params struct {
	NumColumns          int
	NumInputs           int
	PotentialRadius     int
	PotentialPct        float64
	InitConnPct         float64
	SynPermConnected    float64
	GlobalInhibition    bool
	LocalAreaDensity    float64
	StimulusThreshold   int
	SynPermActiveMod    float64
	SynPermInactiveMod  float64
	DutyCyclePeriod     int
	MinOverlapDutyCycle float64
	MinActiveDutyCycle  float64
	MaxBoost            float64
}

// NewV1Params returns a default V1Params set.
func NewV1Params() V1Params {
	return V1Params{
		NumColumns:          2048, // size of output vector
		NumInputs:           400,  // size of input vector
		PotentialRadius:     8,    // # of potential synapses
		PotentialPct:        0.5,  // % sample of potentials
		InitConnPct:         0.3,  // % synapses to connect on init
		SynPermConnected:    0.3,  // synapse connected threshold
		GlobalInhibition:    true, // enable global inhibition
		LocalAreaDensity:    0.02, // sparsity of output vector
		StimulusThreshold:   0,    // used for variable sparsity inputs
		SynPermActiveMod:    0.05, // permanence increment
		SynPermInactiveMod:  0.03, // permanence decrement
		DutyCyclePeriod:     8,    // duty cycle period, in cycles
		MinOverlapDutyCycle: 0.04, // used to bump weak columns
		MinActiveDutyCycle:  0.04, // used to boost weak columns
		MaxBoost:            8.0,  // maximum boost value
	}
}

type proximalSynapse struct {
	idx       int     // input index
	perm      float64 // permanence value
	connected bool    // connected ?
}

type spColumn struct {
	psyns            []proximalSynapse
	overlap          int
	boostedOverlap   int
	boostFactor      float64
	overlapDutyCycle float64
	activeDutyCycle  float64
	active           bool
}

// V1 implementation of spatial pooling.
type V1 struct {
	// state
	cols             []spColumn
	input            []bool
	inhibitionRadius int
	iteration        int
	learnIteration   int

	// params
	P V1Params
}

// NewV1 initializes a new V1 SpatialPooler with the provided V1Params.
func NewV1(p V1Params) *V1 {
	sp := &V1{
		P:              p,
		iteration:      1,
		learnIteration: 1,
	}

	sp.cols = make([]spColumn, p.NumColumns)

	for i := range sp.cols {
		sp.mapPotential(i)
		sp.initpermanence(i)
		sp.updateconnected(i)
		sp.cols[i].boostFactor = 1
	}

	return sp
}

// Updates the connected value on specified column's synapses.
// This should be called every time a synapse is modified.
func (sp *V1) updateconnected(col int) {
	for i, j := range sp.cols[col].psyns {
		if j.perm >= sp.P.SynPermConnected {
			sp.cols[col].psyns[i].connected = true
		} else {
			sp.cols[col].psyns[i].connected = false
		}
	}
}

// TODO : Clipping min/max values
// Initializes permanence of synapses on specified column.
// This method uses a normal distribution centering around the
// synPermConnected parameter, and intializes columns as connected
// based on the initConnPct parameter.
func (sp *V1) initpermanence(col int) {
	sd := 0.05
	var p float64
	for i := range sp.cols[col].psyns {
		chance := rand.Float64()
		switch {
		case chance <= sp.P.InitConnPct:
			p = rand.NormFloat64()*sd + sp.P.SynPermConnected
			for p < sp.P.SynPermConnected {
				p = rand.NormFloat64()*sd + sp.P.SynPermConnected
			}
			sp.cols[col].psyns[i].perm = p
		case chance > sp.P.InitConnPct:
			p = rand.NormFloat64()*sd + sp.P.SynPermConnected
			for p >= sp.P.SynPermConnected {
				p = rand.NormFloat64()*sd + sp.P.SynPermConnected
			}
			sp.cols[col].psyns[i].perm = p
		}
	}
}

// Creates potential synapses on specified column. This method will
// randomly sample the receptive field of a column, and sets the
// potential synapses to the sampled indices.
func (sp *V1) mapPotential(col int) {
	ratio := float64(col) / float64(sp.P.NumColumns)
	center := int(float64(sp.P.NumInputs) * ratio)

	nbs := sp.getInputNeighbors(center)
	n := int(float64(len(nbs)) * sp.P.PotentialPct)
	sample := vec.UniqueRandInts(n, len(nbs))

	sp.cols[col].psyns = make([]proximalSynapse, len(sample))
	for i, j := range sample {
		sp.cols[col].psyns[i].idx = nbs[j]
	}
}

// Returns neighborhood of specified input index. Uses wraparound
// by default.
func (sp *V1) getInputNeighbors(center int) (nbs []int) {
	r := sp.P.PotentialRadius
	for i := center - r; i <= center+r; i++ {
		switch {
		case i >= 0 && i <= sp.P.NumInputs-1:
			nbs = append(nbs, i)
		case i < 0:
			nbs = append(nbs, i+sp.P.NumInputs)
		case i > sp.P.NumInputs-1:
			nbs = append(nbs, i-sp.P.NumInputs)
		}
	}
	return
}

// Compute runs an input vector through the SpatialPooler algorithm,
// and returns a vector containing the active columns. The learn
// parameter specifies whether learning should be performed.
func (sp *V1) Compute(input []bool, learn bool) []bool {
	if len(input) != sp.P.NumInputs {
		panic("Mismatched input dimensions!")
	}
	sp.input = input
	sp.iteration++
	if learn {
		sp.learnIteration++
	}

	sp.updateoverlaps()

	if sp.P.GlobalInhibition || sp.inhibitionRadius > sp.P.NumColumns {
		sp.inhibitColumnsGlobal(learn)
	} else {
		sp.inhibitColumnsLocal(learn)
	}

	if learn {
		sp.adaptSynapses()
		sp.updateoverlapDutyCycles()
		sp.updateactiveDutyCycles()
		sp.bumpWeakColumns()
		sp.updateboostFactors()
	}

	// return active columns
	active := make([]bool, len(sp.cols))
	for i := range sp.cols {
		active[i] = sp.cols[i].active
	}
	return active
}

// Update boost factors for all columns. The boost factors are based
// on the activation duty cycle of each column; columns that activate
// infrequently are boosted higher, columns that are active enough of
// the time are left at 1.0 boost.
func (sp *V1) updateboostFactors() {
	for i, col := range sp.cols {
		if col.activeDutyCycle < sp.P.MinActiveDutyCycle {
			boost := ((1 - sp.P.MaxBoost) / sp.P.MinActiveDutyCycle *
				col.activeDutyCycle) + sp.P.MaxBoost
			sp.cols[i].boostFactor = boost
		}
	}
}

// Increase permanence values for all synapses on weak columns.
func (sp *V1) bumpWeakColumns() {
	for i, col := range sp.cols {
		if col.overlapDutyCycle < sp.P.MinOverlapDutyCycle {
			for j := range col.psyns {
				sp.cols[i].psyns[j].perm += sp.P.SynPermActiveMod
			}
			sp.updateconnected(i)
		}
	}
}

// Update the operlap duty cycles of each column.
func (sp *V1) updateoverlapDutyCycles() {
	period := sp.P.DutyCyclePeriod
	if period > sp.iteration {
		period = sp.iteration
	}

	var o float64
	for i, col := range sp.cols {
		if col.overlap > 0 {
			o = 1.0
		} else {
			o = 0.0
		}
		cycle := (col.overlapDutyCycle*float64(period-1) + o) /
			float64(period)
		sp.cols[i].overlapDutyCycle = cycle
	}
}

// Update the active duty cycles of each column.
func (sp *V1) updateactiveDutyCycles() {
	period := sp.P.DutyCyclePeriod
	if period > sp.iteration {
		period = sp.iteration
	}

	var a float64
	for i, col := range sp.cols {
		if col.active {
			a = 1.0
		} else {
			a = 0.0
		}
		cycle := (col.activeDutyCycle*float64(period-1) + a) /
			float64(period)
		sp.cols[i].activeDutyCycle = cycle
	}
}

// Adapt permanence values of synapses based on the input vector and
// currently active columns post-inhibition. Permanences for synapses
// connected to active inputs are increased, and those connected to
// inactive inputs are decreased.
func (sp *V1) adaptSynapses() {
	for i, col := range sp.cols {
		if col.active {
			for j, syn := range sp.cols[i].psyns {
				if sp.input[syn.idx] {
					sp.cols[i].psyns[j].perm += sp.P.SynPermActiveMod
				} else {
					sp.cols[i].psyns[j].perm -= sp.P.SynPermInactiveMod
				}
			}
			sp.updateconnected(i)
		}
	}
}

// Inhibit columns globally. This method sets the active state on
// each column.
func (sp *V1) inhibitColumnsGlobal(learn bool) {
	overlaps := make([]int, sp.P.NumColumns)
	if learn {
		for i, col := range sp.cols {
			overlaps[i] = col.boostedOverlap
		}
	} else {
		for i, col := range sp.cols {
			overlaps[i] = col.overlap
		}
	}

	// HOTSPOT : ~26% of execution time
	winners := vec.SortIndices(overlaps)

	n := int(sp.P.LocalAreaDensity * float64(sp.P.NumColumns))
	start := len(winners) - n

	// Enforce Stimulus Threshold : useful for varying sparsity input
	for start < len(winners) {
		i := winners[start]
		if overlaps[i] >= sp.P.StimulusThreshold {
			break
		} else {
			start++
		}
	}

	for col := range sp.cols {
		sp.cols[col].active = false
	}

	winners = vec.Reverse(winners[start:]) // [start:]
	for _, col := range winners {
		sp.cols[col].active = true
	}
}

// Inhibit columns locally. This method sets the active state on
// each column.
func (sp *V1) inhibitColumnsLocal(learn bool) {
}

// Update the overlap score on all columns. The overlap is the
// number of connected synapses terminating in an active input bit.
func (sp *V1) updateoverlaps() {
	for i := range sp.cols {
		sp.cols[i].overlap = 0
		for _, syn := range sp.cols[i].psyns {
			if syn.connected && sp.input[syn.idx] {
				sp.cols[i].overlap++
			}
		}

		sp.cols[i].boostedOverlap = int(float64(sp.cols[i].overlap) *
			sp.cols[i].boostFactor)
	}
}
