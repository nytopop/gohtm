package main

import "math/rand"

/* Spatial Pooler */

type SpatialParams struct {
	numColumns          int
	numInputs           int
	potentialRadius     int
	potentialPct        float64
	initConnPct         float64
	synPermConnected    float64
	globalInhibition    bool
	localAreaDensity    float64
	stimulusThreshold   int
	synPermActiveMod    float64
	synPermInactiveMod  float64
	dutyCyclePeriod     int
	minOverlapDutyCycle float64
	minActiveDutyCycle  float64
	maxBoost            float64
}

func NewSpatialParams() SpatialParams {
	return SpatialParams{
		numColumns:          2048, // size of output vector
		numInputs:           400,  // size of input vector
		potentialRadius:     8,    // # of potential synapses
		potentialPct:        0.5,  // % sample of potentials
		initConnPct:         0.3,  // % synapses to connect on init
		synPermConnected:    0.3,  // synapse connected threshold
		globalInhibition:    true, // enable global inhibition
		localAreaDensity:    0.02, // sparsity of output vector
		stimulusThreshold:   0,    // used for variable sparsity inputs
		synPermActiveMod:    0.05, // permanence increment
		synPermInactiveMod:  0.03, // permanence decrement
		dutyCyclePeriod:     8,    // duty cycle period, in cycles
		minOverlapDutyCycle: 0.04, // used to bump weak columns
		minActiveDutyCycle:  0.04, // used to boost weak columns
		maxBoost:            8.0,  // maximum boost value
	}
}

type ProximalSynapse struct {
	idx       int     // input index
	perm      float64 // permanence value
	connected bool    // connected ?
}

type SPColumn struct {
	psyns            []ProximalSynapse
	overlap          int
	boostedOverlap   int
	boostFactor      float64
	overlapDutyCycle float64
	activeDutyCycle  float64
	active           bool
}

type SpatialPooler struct {
	// state
	cols             []SPColumn
	input            SparseBinaryVector
	inhibitionRadius int
	iteration        int
	learnIteration   int

	// params
	numColumns          int
	numInputs           int
	potentialRadius     int
	potentialPct        float64
	initConnPct         float64
	synPermConnected    float64
	globalInhibition    bool
	localAreaDensity    float64
	stimulusThreshold   int
	synPermActiveMod    float64
	synPermInactiveMod  float64
	dutyCyclePeriod     int
	minOverlapDutyCycle float64
	minActiveDutyCycle  float64
	maxBoost            float64
}

/* Initialize a new SpatialPooler with supplied SpatialParams. */
func NewSpatialPooler(p SpatialParams) SpatialPooler {
	sp := SpatialPooler{
		numColumns:          p.numColumns,
		numInputs:           p.numInputs,
		potentialRadius:     p.potentialRadius,
		potentialPct:        p.potentialPct,
		initConnPct:         p.initConnPct,
		synPermConnected:    p.synPermConnected,
		globalInhibition:    p.globalInhibition,
		localAreaDensity:    p.localAreaDensity,
		stimulusThreshold:   p.stimulusThreshold,
		synPermActiveMod:    p.synPermActiveMod,
		synPermInactiveMod:  p.synPermInactiveMod,
		dutyCyclePeriod:     p.dutyCyclePeriod,
		minOverlapDutyCycle: p.minOverlapDutyCycle,
		minActiveDutyCycle:  p.minActiveDutyCycle,
		maxBoost:            p.maxBoost,
		iteration:           1,
		learnIteration:      1,
	}

	sp.cols = make([]SPColumn, p.numColumns)

	for i, _ := range sp.cols {
		sp.mapPotential(i)
		sp.initPermanence(i)
		sp.updateConnected(i)
		sp.cols[i].boostFactor = 1
	}

	return sp
}

/* Updates the connected value on specified column's synapses. This should be called every time a synapse is modified. */
func (sp *SpatialPooler) updateConnected(col int) {
	for i, j := range sp.cols[col].psyns {
		if j.perm >= sp.synPermConnected {
			sp.cols[col].psyns[i].connected = true
		} else {
			sp.cols[col].psyns[i].connected = false
		}
	}
}

// TODO : Clipping min/max values
/* Initializes permanence of synapses on specified column. This method uses a normal distribution centering around the synPermConnected parameter, and intializes columns as connected based on the initConnPct parameter. */
func (sp *SpatialPooler) initPermanence(col int) {
	sd := 0.05
	var p float64
	for i, _ := range sp.cols[col].psyns {
		chance := rand.Float64()
		switch {
		case chance <= sp.initConnPct:
			p = rand.NormFloat64()*sd + sp.synPermConnected
			for p < sp.synPermConnected {
				p = rand.NormFloat64()*sd + sp.synPermConnected
			}
			sp.cols[col].psyns[i].perm = p
		case chance > sp.initConnPct:
			p = rand.NormFloat64()*sd + sp.synPermConnected
			for p >= sp.synPermConnected {
				p = rand.NormFloat64()*sd + sp.synPermConnected
			}
			sp.cols[col].psyns[i].perm = p
		}
	}
}

/* Creates potential synapses on specified column. This method will randomly sample the receptive field of a column, and sets the potential synapses to the sampled indices. */
func (sp *SpatialPooler) mapPotential(col int) {
	ratio := float64(col) / float64(sp.numColumns)
	center := int(float64(sp.numInputs) * ratio)

	nbs := sp.getInputNeighbors(center)
	n := int(float64(len(nbs)) * sp.potentialPct)
	sample := UniqueRandInts(n, len(nbs))

	sp.cols[col].psyns = make([]ProximalSynapse, len(sample))
	for i, j := range sample {
		sp.cols[col].psyns[i].idx = nbs[j]
	}
}

/* Returns neighborhood of specified input index. Uses wraparound by default. */
func (sp *SpatialPooler) getInputNeighbors(center int) (nbs []int) {
	r := sp.potentialRadius
	for i := center - r; i <= center+r; i++ {
		switch {
		case i >= 0 && i <= sp.numInputs-1:
			nbs = append(nbs, i)
		case i < 0:
			nbs = append(nbs, i+sp.numInputs)
		case i > sp.numInputs-1:
			nbs = append(nbs, i-sp.numInputs)
		}
	}
	return
}

/* Compute active columns for a given input vector. */
func (sp *SpatialPooler) Compute(input SparseBinaryVector, learn bool) SparseBinaryVector {
	if input.x != sp.numInputs {
		panic("Mismatched input dimensions!")
	}
	sp.input = input
	sp.iteration += 1
	if learn {
		sp.learnIteration += 1
	}

	sp.updateOverlaps()

	if sp.globalInhibition || sp.inhibitionRadius > sp.numColumns {
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
	active := NewSparseBinaryVector(sp.numColumns)
	for i, col := range sp.cols {
		active.Set(i, col.active)
	}
	return active
}

/* Update boost factors for all columns. The boost factors are based on the activation duty cycle of each column; columns that activate infrequently are boosted higher, columns that are active enough of the time are left at 1.0 boost. */
func (sp *SpatialPooler) updateBoostFactors() {
	for i, col := range sp.cols {
		if col.activeDutyCycle < sp.minActiveDutyCycle {
			boost := ((1 - sp.maxBoost) / sp.minActiveDutyCycle * col.activeDutyCycle) + sp.maxBoost
			sp.cols[i].boostFactor = boost
		}
	}
}

/* Increase permanence values for all synapses on weak columns. */
func (sp *SpatialPooler) bumpWeakColumns() {
	for i, col := range sp.cols {
		if col.overlapDutyCycle < sp.minOverlapDutyCycle {
			for j, _ := range col.psyns {
				sp.cols[i].psyns[j].perm += sp.synPermActiveMod
			}
			sp.updateConnected(i)
		}
	}
}

/**/
func (sp *SpatialPooler) updateOverlapDutyCycles() {
	period := sp.dutyCyclePeriod
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
		cycle := (col.overlapDutyCycle*float64(period-1) + o) / float64(period)
		sp.cols[i].overlapDutyCycle = cycle
	}
}

/**/
func (sp *SpatialPooler) updateActiveDutyCycles() {
	period := sp.dutyCyclePeriod
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
		cycle := (col.activeDutyCycle*float64(period-1) + a) / float64(period)
		sp.cols[i].activeDutyCycle = cycle
	}
}

/* Adapt permanence values of synapses based on the input vector and currently active columns post-inhibition. Permanences for synapses connected to active inputs are increased, and those connected to inactive inputs are decreased. */
func (sp *SpatialPooler) adaptSynapses() {
	for i, col := range sp.cols {
		if col.active {
			for j, syn := range sp.cols[i].psyns {
				if sp.input.Get(syn.idx) {
					sp.cols[i].psyns[j].perm += sp.synPermActiveMod
				} else {
					sp.cols[i].psyns[j].perm -= sp.synPermInactiveMod
				}
			}
			sp.updateConnected(i)
		}
	}
}

/* Inhibit columns globally. This method sets the active state on each column. */
func (sp *SpatialPooler) inhibitColumnsGlobal(learn bool) {
	overlaps := make([]int, sp.numColumns)
	if learn {
		for i, col := range sp.cols {
			overlaps[i] = col.boostedOverlap
		}
	} else {
		for i, col := range sp.cols {
			overlaps[i] = col.overlap
		}
	}
	winners := SortIndices(overlaps)

	n := int(sp.localAreaDensity * float64(sp.numColumns))
	start := len(winners) - n

	// Enforce Stimulus Threshold : useful for varying sparsity input
	for start < len(winners) {
		i := winners[start]
		if overlaps[i] >= sp.stimulusThreshold {
			break
		} else {
			start++
		}
	}

	for col, _ := range sp.cols {
		sp.cols[col].active = false
	}

	winners = Reverse(winners[start:]) // [start:]
	for _, col := range winners {
		sp.cols[col].active = true
	}
}

/* Inhibit columns locally. This method sets the active state on each column. */
func (sp *SpatialPooler) inhibitColumnsLocal(learn bool) {
}

/* Update the overlap score on all columns. The overlap is the number of connected synapses terminating in an active input bit. */
func (sp *SpatialPooler) updateOverlaps() {
	for i, _ := range sp.cols {
		sp.cols[i].overlap = 0
		for _, syn := range sp.cols[i].psyns {
			if syn.connected && sp.input.Get(syn.idx) {
				sp.cols[i].overlap += 1
			}
		}

		sp.cols[i].boostedOverlap = int(float64(sp.cols[i].overlap) * sp.cols[i].boostFactor)
	}
}
