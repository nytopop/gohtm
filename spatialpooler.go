package main

import (
	"fmt"
	"math/rand"
)

/* Spatial Pooling
1. Start with input of fixed n bits.
2. Assign fixed number of columns to region receiving the input. Each column has associated dendrite segment. Each dendrite segment has a set of potential synapses representing subset of input bits. Each potential synapse has a permanence value; based on that value, some potential synapses will be connected to the dendrite segment.
3. For any input, determine how many connected synapses on a column are connected to active input bits.
4. Number of active synapses is multiplied by a boosting factor; dynamically determined by how often a column is active relative to neighbors.
5. A fixed n of columns within the inhibition radius with highest activations after boosting become active and disable rest of columns within the radius. The inhibition radius is dynamically determined by the spread of input bits. We should now have a sparse set of active columns.
6. For each active column, we adjust the permanence values of all potential synapses. The permanence of synapses aligned with active input bits is increased. The changes may change some synapses from connected<->disconnected.
*/

type SpatialParams struct {
	numColumns       int
	numInputs        int
	columnHeight     int
	initConnPct      float64
	potentialPct     float64
	sparsity         float64
	potentialRadius  int
	synPermConnected float64
	globalInhibition bool
}

// Return default SpatialParams
func NewSpatialParams() SpatialParams {
	return SpatialParams{
		numColumns:       48,
		numInputs:        32,
		columnHeight:     16,
		initConnPct:      0.5,
		potentialPct:     0.3,
		sparsity:         0.02,
		potentialRadius:  5,
		synPermConnected: 0.3,
		globalInhibition: true,
	}
}

type SpatialPooler struct {
	// Params
	numColumns       int     // number of columns
	numInputs        int     // number of inputs
	columnHeight     int     // column height
	initConnPct      float64 // % of potential to connect on init
	potentialPct     float64 // % of receptive field to potentially con
	sparsity         float64 // sparsity of output vector
	potentialRadius  int     // radius of column's receptive field
	synPermConnected float64 // connected permanence threshhold
	globalInhibition bool    // global inhibition enabled?

	// State
	potential            SparseBinaryMatrix // map of columns to input space
	connected            SparseBinaryMatrix // connected synapses
	numConnected         []int              // number of connected syns per col
	permanences          SparseFloatMatrix  // permanences
	tieBreaker           []float64          // random tie breaker for ties
	overlapDutyCycles    []int
	activeDutyCycles     []int
	minOverlapDutyCycles []int
	minActiveDutyCycles  []int
	boostFactors         []int

	inhibitionRadius int
	iteration        int // iteration counter
}

// Initialize a new SpatialPooler
func NewSpatialPooler(p SpatialParams) SpatialPooler {
	sp := SpatialPooler{
		numColumns:       p.numColumns,
		numInputs:        p.numInputs,
		columnHeight:     p.columnHeight,
		initConnPct:      p.initConnPct,
		potentialPct:     p.potentialPct,
		sparsity:         p.sparsity,
		potentialRadius:  p.potentialRadius,
		synPermConnected: p.synPermConnected,
		globalInhibition: p.globalInhibition,
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
		sp.updateColumnPermanence(i, perm)
	}

	fmt.Println(sp.potential)
	fmt.Println(sp.permanences)
	fmt.Println(sp.connected)
	fmt.Println(sp.numConnected)

	sp.overlapDutyCycles = make([]int, sp.numColumns)
	sp.activeDutyCycles = make([]int, sp.numColumns)
	sp.minOverlapDutyCycles = make([]int, sp.numColumns)
	sp.minActiveDutyCycles = make([]int, sp.numColumns)
	sp.boostFactors = make([]int, sp.numColumns)
	for i, _ := range sp.boostFactors {
		sp.boostFactors[i] = 1
	}

	sp.inhibitionRadius = 0
	sp.updateInhibitionRadius()

	return sp
}

func (sp SpatialPooler) updateInhibitionRadius() {
	// no-op, for now TODO
}

// Update the sp.permanence values for a column
// TODO : clipping of permanence values to min, max
func (sp SpatialPooler) updateColumnPermanence(i int, perm []float64) {
	var c int
	for k, v := range perm {
		sp.permanences.Set(i, k, v)
		if v >= sp.synPermConnected {
			sp.connected.Set(i, k, true)
			c++
		} else {
			sp.connected.Set(i, k, false)
		}
	}
	sp.numConnected[i] = c
}

// Map potential connections to initial permanence values
func (sp SpatialPooler) initPermanence(potential []int) []float64 {
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

// Map a column to its input bits.
// Returns [input indices] for a column index.
func (sp SpatialPooler) mapPotential(i int) []int {
	// Find center of column's receptive field
	ratio := float64(i) / float64(sp.numColumns)
	center := int(float64(sp.numInputs) * ratio)

	// Return random sample of column's receptive field
	neigh := sp.getInputNeighborhood(center)
	n := int(float64(len(neigh)) * sp.potentialPct)
	sample := UniqueRandInts(n, len(neigh))

	field := make([]int, len(sample))
	for j, k := range sample {
		field[j] = neigh[k]
	}

	return field
}

// Return the input neighborhood of an input index
func (sp SpatialPooler) getInputNeighborhood(i int) (sv []int) {
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

func (sp SpatialPooler) Compute(sv SparseBinaryVector) SparseBinaryVector {
	return SparseBinaryVector{}
}
