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
	numColumns      int
	numInputs       int
	columnHeight    int
	initConnPercent float64
	sparsity        float64
	potentialRadius int
}

// Return default SpatialParams
func NewSpatialParams() SpatialParams {
	return SpatialParams{
		numColumns:      48,
		numInputs:       32,
		columnHeight:    16,
		initConnPercent: 0.3,
		sparsity:        0.02,
		potentialRadius: 5,
	}
}

type SpatialPooler struct {
	// Params
	numColumns      int     // number of columns
	numInputs       int     // number of inputs
	columnHeight    int     // column height
	initConnPercent float64 // percent to connect on init
	sparsity        float64 // sparsity of output vector
	potentialRadius int     // radius of column's receptive field

	// State
	columns      SparseBinaryVector // columns
	potential    SparseBinaryMatrix // map of columns to input space
	connected    SparseBinaryMatrix // connected synapses
	numConnected []int              // number of connected syns per col
	permanences  SparseFloatMatrix  // permanences
	tieBreaker   []float64          // random tie breaker for ties

	iteration int // iteration counter
}

// Initialize a new SpatialPooler
func NewSpatialPooler(p SpatialParams) SpatialPooler {
	sp := SpatialPooler{
		numColumns:      p.numColumns,
		numInputs:       p.numInputs,
		columnHeight:    p.columnHeight,
		initConnPercent: p.initConnPercent,
		sparsity:        p.sparsity,
		potentialRadius: p.potentialRadius,
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

	// Initialize permanence values
	for i := 0; i < sp.numColumns; i++ {
		potential := sp.mapPotential(i)
		fmt.Println(potential)
		// perm := initPermanance(potential, connPercent)
		// potential.d[i] = potential ?
	}

	return sp
}

// Map a column to its input bits.
// Returns [input indices] for a column index.
func (sp SpatialPooler) mapPotential(i int) []int {
	// Find center of column's receptive field
	ratio := float64(i) / float64(sp.numColumns)
	center := int(float64(sp.numInputs) * ratio)

	// Return random sample of column's receptive field
	neigh := sp.getInputNeighborhood(center)
	n := int(float64(len(neigh)) * sp.initConnPercent)
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
