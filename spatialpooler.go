package main

import "math/rand"

/* Spatial Pooling
1. Start with input of fixed n bits.
2. Assign fixed number of columns to region receiving the input. Each column has associated dendrite segment. Each dendrite segment has a set of potential synapses representing subset of input bits. Each potential synapse has a permanence value; based on that value, some potential synapses will be connected to the dendrite segment.
3. For any input, determine how many connected synapses on a column are connected to active input bits.
4. Number of active synapses is multiplied by a boosting factor; dynamically determined by how often a column is active relative to neighbors.
5. A fixed n of columns within the inhibition radius with highest activations after boosting become active and disable rest of columns within the radius. The inhibition radius is dynamically determined by the spread of input bits. We should now have a sparse set of active columns.
6. For each active column, we adjust the permanence values of all potential synapses. The permanence of synapses aligned with active input bits is increased. The changes may change some synapses from connected<->disconnected.
*/

type SpatialParams struct {
	numColumns        int
	numInputs         int
	columnHeight      int
	initConnPct       float64
	potentialPct      float64
	sparsity          float64
	potentialRadius   int
	synPermConnected  float64
	globalInhibition  bool
	localAreaDensity  float64
	stimulusThreshold int
	synPermActiveMod  float64
	dutyCyclePeriod   int
}

// Return default SpatialParams
func NewSpatialParams() SpatialParams {
	return SpatialParams{
		numColumns:        1024,
		numInputs:         400,
		columnHeight:      16,
		initConnPct:       0.5,
		potentialPct:      0.3,
		sparsity:          0.02,
		potentialRadius:   5,
		synPermConnected:  0.3,
		globalInhibition:  true,
		localAreaDensity:  0.5,
		stimulusThreshold: 2,
		synPermActiveMod:  0.05,
		dutyCyclePeriod:   5,
	}
}

type SpatialPooler struct {
	// Params
	numColumns        int     // number of columns
	numInputs         int     // number of inputs
	columnHeight      int     // column height
	initConnPct       float64 // % of potential to connect on init
	potentialPct      float64 // % of receptive field to potentially con
	sparsity          float64 // sparsity of output vector
	potentialRadius   int     // radius of column's receptive field
	synPermConnected  float64 // connected permanence threshhold
	globalInhibition  bool    // global inhibition enabled?
	localAreaDensity  float64 // local area density?
	stimulusThreshold int     // stimulus threshold ?
	synPermActiveMod  float64 // permanence increment/decrement val

	// State
	potential            SparseBinaryMatrix // map of columns to input space
	connected            SparseBinaryMatrix // connected synapses
	numConnected         []int              // number of connected syns per col
	permanences          SparseFloatMatrix  // permanences
	tieBreaker           []float64          // random tie breaker for ties
	overlapDutyCycles    []float64
	activeDutyCycles     []float64
	minOverlapDutyCycles []int
	minActiveDutyCycles  []int
	boostFactors         []int
	boostedOverlaps      []int

	inhibitionRadius int
	dutyCyclePeriod  int
	iteration        int // iteration counter
	learnIteration   int
}

/* Initialize a new SpatialPooler. */
func NewSpatialPooler(p SpatialParams) SpatialPooler {
	sp := SpatialPooler{
		numColumns:        p.numColumns,
		numInputs:         p.numInputs,
		columnHeight:      p.columnHeight,
		initConnPct:       p.initConnPct,
		potentialPct:      p.potentialPct,
		sparsity:          p.sparsity,
		potentialRadius:   p.potentialRadius,
		synPermConnected:  p.synPermConnected,
		globalInhibition:  p.globalInhibition,
		localAreaDensity:  p.localAreaDensity,
		stimulusThreshold: p.stimulusThreshold,
		synPermActiveMod:  p.synPermActiveMod,
		dutyCyclePeriod:   p.dutyCyclePeriod,
	}

	// Initialize data structures
	sp.potential = NewSparseBinaryMatrix(sp.numColumns, sp.numInputs)
	sp.permanences = NewSparseFloatMatrix(sp.numColumns, sp.numInputs)
	sp.connected = NewSparseBinaryMatrix(sp.numColumns, sp.numInputs)
	sp.numConnected = make([]int, sp.numColumns)

	// initialize random tie breaker
	sp.tieBreaker = make([]float64, sp.numColumns)
	for i, _ := range sp.tieBreaker {
		sp.tieBreaker[i] = 0.01 + rand.Float64()
	}
	// Initialize sp.potential, sp.permanences, sp.connected, sp.numConnected
	for i := 0; i < sp.numColumns; i++ {
		potential := sp.mapPotential(i)
		perm := sp.initPermanence(potential)
		for _, v := range potential {
			sp.potential.Set(i, v, true)
		}
		sp.updateColumnPermanence(i, perm, potential)
	}

	/*
		fmt.Println(sp.potential)
		fmt.Println(sp.permanences)
		fmt.Println(sp.connected)
		fmt.Println(sp.numConnected)
	*/

	sp.overlapDutyCycles = make([]float64, sp.numColumns)
	sp.activeDutyCycles = make([]float64, sp.numColumns)
	sp.minOverlapDutyCycles = make([]int, sp.numColumns)
	sp.minActiveDutyCycles = make([]int, sp.numColumns)
	sp.boostFactors = make([]int, sp.numColumns)
	for i, _ := range sp.boostFactors {
		sp.boostFactors[i] = 1
	}

	sp.inhibitionRadius = 0
	sp.updateInhibitionRadius()
	sp.updateBookkeepingVars(false)
	return sp
}

func (sp SpatialPooler) updateInhibitionRadius() {
	// no-op, for now TODO
}

// TODO : clipping of permanence values to min, max
/* Update the permanence values for synapses 'potential', with values from 'perm', on column index 'i'. */
func (sp *SpatialPooler) updateColumnPermanence(i int, perm []float64, potential []int) {
	var c int
	for k, v := range perm {
		sp.permanences.Set(i, potential[k], v)
		if v >= sp.synPermConnected {
			sp.connected.Set(i, potential[k], true)
			c++
		} else {
			sp.connected.Set(i, potential[k], false)
		}
	}
	sp.numConnected[i] = c
}

/* Initialize the permanence values of synapses in the potential pool. This method will initialize synapses with a permanence value within a normal distribution centering at sp.synPermConnected, with a P(sp.initConnPct) chance of being initialized connected. Standard deviation of 0.05. */
func (sp *SpatialPooler) initPermanence(potential []int) []float64 {
	perm := make([]float64, len(potential))
	sd := 0.05
	for i := 0; i < len(perm); i++ {
		var p float64
		chance := rand.Float64()
		switch {
		case chance <= sp.initConnPct:
			// set to right side of normal distribution
			p = rand.NormFloat64()*sd + sp.synPermConnected
			for p < sp.synPermConnected {
				p = rand.NormFloat64()*sd + sp.synPermConnected
			}
			perm[i] = p
		case chance > sp.initConnPct:
			// set to left side of normal distribution
			p = rand.NormFloat64()*sd + sp.synPermConnected
			for p > sp.synPermConnected {
				p = rand.NormFloat64()*sd + sp.synPermConnected
			}
			perm[i] = p
		}
	}
	return perm
}

/* Map a column to its input bits. This method will randomly sample the space of a column's receptive field, and returns the indices of the sampled inputs. */
func (sp *SpatialPooler) mapPotential(i int) []int {
	// Find center of column's receptive field
	ratio := float64(i) / float64(sp.numColumns)
	center := int(float64(sp.numInputs) * ratio)

	// Return random sample of columns's receptive field
	neigh := sp.getInputNeighborhood(center)
	n := int(float64(len(neigh)) * sp.potentialPct)
	sample := UniqueRandInts(n, len(neigh))

	field := make([]int, len(sample))
	for j, k := range sample {
		field[j] = neigh[k]
	}

	return field
}

/* Return the input neighborhood of an input index. */
func (sp *SpatialPooler) getInputNeighborhood(i int) (sv []int) {
	r := sp.potentialRadius
	for ii := i - r; ii <= i+r; ii++ {
		switch {
		case ii >= 0 && ii <= sp.numInputs-1:
			sv = append(sv, ii)
		case ii < 0:
			sv = append(sv, ii+sp.numInputs)
		case ii > sp.numInputs-1:
			sv = append(sv, ii-sp.numInputs)
		}
	}
	return
}

/* Primary method to be called on a SpatialPooler struct. This method returns the active columns for a provided input vector, and learns if the learn flag is set. */
func (sp *SpatialPooler) Compute(in SparseBinaryVector, learn bool) SparseBinaryVector {
	if in.x != sp.numInputs {
		panic("Mismatched input dimensions!")
	}

	sp.updateBookkeepingVars(learn)
	overlaps := sp.calcOverlap(in)

	// Apply boost if learn
	sp.boostedOverlaps = make([]int, sp.numColumns)
	if learn {
		for i := 0; i < sp.numColumns; i++ {
			sp.boostedOverlaps[i] = overlaps[i] * sp.boostFactors[i]
		}
	} else {
		sp.boostedOverlaps = overlaps
	}

	// Apply inhibition to find winning columns
	var activeColumns []int
	if sp.globalInhibition || sp.inhibitionRadius > sp.columnHeight {
		activeColumns = sp.inhibitColumnsGlobal()
	} else {
		activeColumns = sp.inhibitColumnsLocal()
	}

	if learn {
		sp.adaptSynapses(in, activeColumns)
		sp.updateDutyCycles(overlaps, activeColumns)
		sp.bumpWeakColumns()
		sp.updateBoostFactors()
	}

	return SparseBinaryVector{}
}

/* Update boost factors for all columns. */
func (sp *SpatialPooler) updateBoostFactors() {

}

/* Increase the permanence values of synapses whose columns are too dormant. Dormance is determined by the overlap duty cycle of the column. */
func (sp *SpatialPooler) bumpWeakColumns() {
	for i, v := range sp.overlapDutyCycles {
		if v < float64(sp.minOverlapDutyCycles[i]) {
			for j := 0; j < sp.numInputs; j++ {
				if sp.permanences.Exists(i, j) {
					perm := sp.permanences.Get(i, j) + sp.synPermActiveMod
					sp.permanences.Set(i, j, perm)
				}
			}
		}
	}
}

/* Updates the duty cycles for each column. Overlap duty cycle is a moving average of the number of inputs overlapped with each column; active duty cycle is a moving average of the activation frequency for each column. */
func (sp *SpatialPooler) updateDutyCycles(overlaps, active []int) {
	period := sp.dutyCyclePeriod
	if period > sp.iteration {
		period = sp.iteration
	}

	// update overlap duty cycles
	for i, _ := range sp.overlapDutyCycles {
		sp.overlapDutyCycles[i] = (sp.overlapDutyCycles[i]*float64((period-1)) + float64(overlaps[i])) / float64(period)
	}

	// update active duty cycles
	act := NewSparseBinaryVector(sp.numColumns)
	for _, c := range active {
		act.Set(c, true)
	}

	for i, _ := range sp.activeDutyCycles {
		var a float64
		if act.Get(i) {
			a = 1
		} else {
			a = 0
		}
		sp.activeDutyCycles[i] = (sp.activeDutyCycles[i]*float64((period-1)) + a) / float64(period)
	}
}

/* Adapt the permanence values of synapses based on the input vector and currently active columns post-inhibition. Permanences for synapses connected to active inputs are increased, and those connected to inactive inputs are decreased. */
func (sp *SpatialPooler) adaptSynapses(in SparseBinaryVector, active []int) {
	// loop through active columns
	for _, i := range active {
		// loop through synapses on the column
		for j := 0; j < sp.numInputs; j++ {
			if sp.permanences.Exists(i, j) {
				if in.Get(j) {
					perm := sp.permanences.Get(i, j) + sp.synPermActiveMod
					sp.permanences.Set(i, j, perm)
				} else {
					perm := sp.permanences.Get(i, j) - sp.synPermActiveMod
					sp.permanences.Set(i, j, perm)
				}
			}
		}
	}
}

/* Inhibit columns globally. This method returns the indices of columns that remain as winners after the inhibition process. */
func (sp *SpatialPooler) inhibitColumnsGlobal() []int {
	// active per inh area
	n := int(sp.localAreaDensity * float64(sp.numColumns))

	// get winners by sorting sp.boostedOverlaps
	indices := SortIndices(sp.boostedOverlaps)

	// enforce stimulus threshold
	start := len(indices) - n
	for start < len(indices) {
		i := indices[start]
		if sp.boostedOverlaps[i] >= sp.stimulusThreshold {
			break
		} else {
			start++
		}
	}

	return Reverse(indices[start:])
}

// Inhibit columns locally
func (sp *SpatialPooler) inhibitColumnsLocal() []int {
	return []int{}
}

/* Calculate the overlap of columns. The overlap is the number of connected synapses that are connected to an active input bit. Returns a slice of integers with indices correlating to columns. */
func (sp *SpatialPooler) calcOverlap(in SparseBinaryVector) []int {
	out := make([]int, sp.numColumns)

	// iterate through columns
	for i := 0; i < sp.numColumns; i++ {
		// for every column, iterate through connected synapses
		for j := 0; j < sp.numInputs; j++ {
			// for every connected synapse, check if input is active
			if sp.connected.Get(i, j) && in.Get(j) {
				out[i]++
			}
		}
	}

	return out
}

/* Update the sp.iteration and sp.learnIteration counters. */
func (sp *SpatialPooler) updateBookkeepingVars(learn bool) {
	sp.iteration += 1
	if learn {
		sp.learnIteration += 1
	}
}
