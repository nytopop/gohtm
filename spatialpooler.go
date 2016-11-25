package gohtm

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
)

/* Spatial Pooler */

type SpatialParams struct {
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

func NewSpatialParams() SpatialParams {
	return SpatialParams{
		NumColumns:          2048, // size of output vector
		NumInputs:           400,  // size of input vector
		PotentialRadius:     8,    // # of potential synapses
		PotentialPct:        0.5,  // % sample of potentials
		InitConnPct:         0.3,  // % synapses to connect on init
		SynPermConnected:    0.3,  // synapse Connected threshold
		GlobalInhibition:    true, // enable global inhibition
		LocalAreaDensity:    0.02, // sparsity of output vector
		StimulusThreshold:   0,    // used for variable sparsity inputs
		SynPermActiveMod:    0.05, // Permanence increment
		SynPermInactiveMod:  0.03, // Permanence decrement
		DutyCyclePeriod:     8,    // duty cycle period, in cycles
		MinOverlapDutyCycle: 0.04, // used to bump weak columns
		MinActiveDutyCycle:  0.04, // used to boost weak columns
		MaxBoost:            8.0,  // maximum boost value
	}
}

type ProximalSynapse struct {
	Idx       int     // input index
	Perm      float64 // Permanence value
	Connected bool    // Connected ?
}

type SPColumn struct {
	PSyns            []ProximalSynapse
	Overlap          int
	BoostedOverlap   int
	BoostFactor      float64
	OverlapDutyCycle float64
	ActiveDutyCycle  float64
	Active           bool
}

type SpatialPooler struct {
	// state
	Cols             []SPColumn
	input            SparseBinaryVector
	InhibitionRadius int
	Iteration        int
	LearnIteration   int

	// params
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

/* Initialize a new SpatialPooler with supplied SpatialParams. */
func NewSpatialPooler(p SpatialParams) SpatialPooler {
	sp := SpatialPooler{
		NumColumns:          p.NumColumns,
		NumInputs:           p.NumInputs,
		PotentialRadius:     p.PotentialRadius,
		PotentialPct:        p.PotentialPct,
		InitConnPct:         p.InitConnPct,
		SynPermConnected:    p.SynPermConnected,
		GlobalInhibition:    p.GlobalInhibition,
		LocalAreaDensity:    p.LocalAreaDensity,
		StimulusThreshold:   p.StimulusThreshold,
		SynPermActiveMod:    p.SynPermActiveMod,
		SynPermInactiveMod:  p.SynPermInactiveMod,
		DutyCyclePeriod:     p.DutyCyclePeriod,
		MinOverlapDutyCycle: p.MinOverlapDutyCycle,
		MinActiveDutyCycle:  p.MinActiveDutyCycle,
		MaxBoost:            p.MaxBoost,
		Iteration:           1,
		LearnIteration:      1,
	}

	sp.Cols = make([]SPColumn, p.NumColumns)

	for i, _ := range sp.Cols {
		sp.mapPotential(i)
		sp.initPermanence(i)
		sp.updateConnected(i)
		sp.Cols[i].BoostFactor = 1
	}

	return sp
}

func (sp SpatialPooler) Save(filename string) {
	js, _ := json.Marshal(sp)
	ioutil.WriteFile(filename, js, 0644)
}

/* Updates the Connected value on specified column's synapses. This should be called every time a synapse is modified. */
func (sp *SpatialPooler) updateConnected(col int) {
	for i, j := range sp.Cols[col].PSyns {
		if j.Perm >= sp.SynPermConnected {
			sp.Cols[col].PSyns[i].Connected = true
		} else {
			sp.Cols[col].PSyns[i].Connected = false
		}
	}
}

// TODO : Clipping min/max values
/* Initializes Permanence of synapses on specified column. This method uses a normal distribution centering around the SynPermConnected parameter, and intializes columns as Connected based on the InitConnPct parameter. */
func (sp *SpatialPooler) initPermanence(col int) {
	sd := 0.05
	var p float64
	for i, _ := range sp.Cols[col].PSyns {
		chance := rand.Float64()
		switch {
		case chance <= sp.InitConnPct:
			p = rand.NormFloat64()*sd + sp.SynPermConnected
			for p < sp.SynPermConnected {
				p = rand.NormFloat64()*sd + sp.SynPermConnected
			}
			sp.Cols[col].PSyns[i].Perm = p
		case chance > sp.InitConnPct:
			p = rand.NormFloat64()*sd + sp.SynPermConnected
			for p >= sp.SynPermConnected {
				p = rand.NormFloat64()*sd + sp.SynPermConnected
			}
			sp.Cols[col].PSyns[i].Perm = p
		}
	}
}

/* Creates potential synapses on specified column. This method will randomly sample the receptive field of a column, and sets the potential synapses to the sampled indices. */
func (sp *SpatialPooler) mapPotential(col int) {
	ratio := float64(col) / float64(sp.NumColumns)
	center := int(float64(sp.NumInputs) * ratio)

	nbs := sp.getInputNeighbors(center)
	n := int(float64(len(nbs)) * sp.PotentialPct)
	sample := uniqueRandInts(n, len(nbs))

	sp.Cols[col].PSyns = make([]ProximalSynapse, len(sample))
	for i, j := range sample {
		sp.Cols[col].PSyns[i].Idx = nbs[j]
	}
}

/* Returns neighborhood of specified input index. Uses wraparound by default. */
func (sp *SpatialPooler) getInputNeighbors(center int) (nbs []int) {
	r := sp.PotentialRadius
	for i := center - r; i <= center+r; i++ {
		switch {
		case i >= 0 && i <= sp.NumInputs-1:
			nbs = append(nbs, i)
		case i < 0:
			nbs = append(nbs, i+sp.NumInputs)
		case i > sp.NumInputs-1:
			nbs = append(nbs, i-sp.NumInputs)
		}
	}
	return
}

/* Compute active columns for a given input vector. */
func (sp *SpatialPooler) Compute(input SparseBinaryVector, learn bool) SparseBinaryVector {
	if input.x != sp.NumInputs {
		panic("Mismatched input dimensions!")
	}
	sp.input = input
	sp.Iteration += 1
	if learn {
		sp.LearnIteration += 1
	}

	sp.updateOverlaps()

	if sp.GlobalInhibition || sp.InhibitionRadius > sp.NumColumns {
		sp.inhibitColumnsGlobal(learn)
	} else {
		sp.inhibitColumnsLocal(learn)
	}

	if learn {
		sp.adaptSynapses()
		sp.updateOverlapDutyCycles()
		sp.updateActiveDutyCycles()
		sp.bumpWeakColumns()
		sp.updateBoostFactors()
	}

	// return active columns
	active := NewSparseBinaryVector(sp.NumColumns)
	for i, col := range sp.Cols {
		active.Set(i, col.Active)
	}
	return active
}

/* Update boost factors for all columns. The boost factors are based on the activation duty cycle of each column; columns that activate infrequently are boosted higher, columns that are active enough of the time are left at 1.0 boost. */
func (sp *SpatialPooler) updateBoostFactors() {
	for i, col := range sp.Cols {
		if col.ActiveDutyCycle < sp.MinActiveDutyCycle {
			boost := ((1 - sp.MaxBoost) / sp.MinActiveDutyCycle *
				col.ActiveDutyCycle) + sp.MaxBoost
			sp.Cols[i].BoostFactor = boost
		}
	}
}

/* Increase Permanence values for all synapses on weak columns. */
func (sp *SpatialPooler) bumpWeakColumns() {
	for i, col := range sp.Cols {
		if col.OverlapDutyCycle < sp.MinOverlapDutyCycle {
			for j, _ := range col.PSyns {
				sp.Cols[i].PSyns[j].Perm += sp.SynPermActiveMod
			}
			sp.updateConnected(i)
		}
	}
}

/**/
func (sp *SpatialPooler) updateOverlapDutyCycles() {
	period := sp.DutyCyclePeriod
	if period > sp.Iteration {
		period = sp.Iteration
	}

	var o float64
	for i, col := range sp.Cols {
		if col.Overlap > 0 {
			o = 1.0
		} else {
			o = 0.0
		}
		cycle := (col.OverlapDutyCycle*float64(period-1) + o) /
			float64(period)
		sp.Cols[i].OverlapDutyCycle = cycle
	}
}

/**/
func (sp *SpatialPooler) updateActiveDutyCycles() {
	period := sp.DutyCyclePeriod
	if period > sp.Iteration {
		period = sp.Iteration
	}

	var a float64
	for i, col := range sp.Cols {
		if col.Active {
			a = 1.0
		} else {
			a = 0.0
		}
		cycle := (col.ActiveDutyCycle*float64(period-1) + a) /
			float64(period)
		sp.Cols[i].ActiveDutyCycle = cycle
	}
}

/* Adapt Permanence values of synapses based on the input vector and currently active columns post-inhibition. Permanences for synapses Connected to active inputs are increased, and those Connected to inactive inputs are decreased. */
func (sp *SpatialPooler) adaptSynapses() {
	for i, col := range sp.Cols {
		if col.Active {
			for j, syn := range sp.Cols[i].PSyns {
				if sp.input.Get(syn.Idx) {
					sp.Cols[i].PSyns[j].Perm += sp.SynPermActiveMod
				} else {
					sp.Cols[i].PSyns[j].Perm -= sp.SynPermInactiveMod
				}
			}
			sp.updateConnected(i)
		}
	}
}

/* Inhibit columns globally. This method sets the active state on each column. */
func (sp *SpatialPooler) inhibitColumnsGlobal(learn bool) {
	overlaps := make([]int, sp.NumColumns)
	if learn {
		for i, col := range sp.Cols {
			overlaps[i] = col.BoostedOverlap
		}
	} else {
		for i, col := range sp.Cols {
			overlaps[i] = col.Overlap
		}
	}
	winners := sortIndices(overlaps)

	n := int(sp.LocalAreaDensity * float64(sp.NumColumns))
	start := len(winners) - n

	// Enforce Stimulus Threshold : useful for varying sparsity input
	for start < len(winners) {
		i := winners[start]
		if overlaps[i] >= sp.StimulusThreshold {
			break
		} else {
			start++
		}
	}

	for col, _ := range sp.Cols {
		sp.Cols[col].Active = false
	}

	winners = reverse(winners[start:]) // [start:]
	for _, col := range winners {
		sp.Cols[col].Active = true
	}
}

/* Inhibit columns locally. This method sets the active state on each column. */
func (sp *SpatialPooler) inhibitColumnsLocal(learn bool) {
}

/* Update the Overlap score on all columns. The Overlap is the number of Connected synapses terminating in an active input bit. */
func (sp *SpatialPooler) updateOverlaps() {
	for i, _ := range sp.Cols {
		sp.Cols[i].Overlap = 0
		for _, syn := range sp.Cols[i].PSyns {
			if syn.Connected && sp.input.Get(syn.Idx) {
				sp.Cols[i].Overlap += 1
			}
		}

		sp.Cols[i].BoostedOverlap = int(float64(sp.Cols[i].Overlap) *
			sp.Cols[i].BoostFactor)
	}
}
