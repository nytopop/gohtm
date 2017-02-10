package cells

import "math/rand"

/*
Synapses connect a segment --> cell

cell1: A.1.2.3.4
cell2: B.1.2.3.4

syn1: A.1 --> B
syn2: B.1 --> A
NOT DUPLICATES!

Do everything at the index level
cells -> columns -> segments -> synapses
*/

// V1Params contains parameters for initialization of V1 cellular state.
type V1Params struct {
	NumColumns         int     // number of columns
	CellsPerCol        int     // cells per column
	SegsPerCell        int     // max segments per cell
	SynsPerSeg         int     // max synapses per segment
	SynPermConnected   float32 // synapse permanence connection threshold
	SynPermActiveMod   float32 // synapse permanence increment
	SynPermInactiveMod float32 // synapse permanence decrement
	InitPerm           float32 // initial permanence of new synapses
	ActiveThreshold    int     // number live synapses to mark segment active
	MatchingThreshold  int     // # of potential synapses for learning ????
	MaxNewSyns         int     // # of new synapses to grow per cell
}

// NewV1Params returns a default parameter set.
func NewV1Params() V1Params {
	return V1Params{
		NumColumns:         64,
		CellsPerCol:        8,
		SegsPerCell:        16,
		SynsPerSeg:         16,
		SynPermConnected:   0.5,
		SynPermActiveMod:   0.05,
		SynPermInactiveMod: 0.03,
		InitPerm:           0.21,
		ActiveThreshold:    12,
		MatchingThreshold:  8,
		MaxNewSyns:         20,
	}
}

// V1 cellular state interface.
type V1 struct {
	P            V1Params
	cells        []V1Cell
	iteration    int
	nSegs, nSyns int
}

// NewV1 returns a new V1 instance, initialized according to provided params.
func NewV1(p V1Params) *V1 {
	v := &V1{
		P:         p,
		iteration: 0,
	}

	// allocate cells
	v.cells = make([]V1Cell, p.NumColumns*p.CellsPerCol)
	return v
}

// V1Cell represents a single cell.
type V1Cell struct {
	active, matching int // number of active & matching segments
	segments         []V1Segment
}

// V1Segment represents a single dendritic segment attached to a cell.
type V1Segment struct {
	active, matching bool
	live, dead       int
	synapses         []V1Synapse
}

// V1Synapse represents the connection from the dendritic segment of a cell
// to another cell, called the presynaptic cell.
type V1Synapse struct {
	cell int     // presynaptic cell
	perm float32 // permanence of connection
}

func (v *V1) CreateSegment(cell int) int {
	// BUG :: v.P.SegsPerCell ignored!
	v.cells[cell].segments = append(
		v.cells[cell].segments,
		V1Segment{
			active:   false,
			matching: false,
			synapses: make([]V1Synapse, 0)})

	return len(v.cells[cell].segments) - 1
}

func (v *V1) CreateSynapse(cell, seg, target int) {
	// BUG :: v.P.SynsPerSeg ignored!
	// BUG :: duplicate synapses! check before creating a new synapse
	v.cells[cell].segments[seg].synapses = append(
		v.cells[cell].segments[seg].synapses,
		V1Synapse{
			cell: target,
			perm: v.P.InitPerm})
}

// AdaptSegment adapts synapses on a segment to the provided slice
// of cells active in the previous time step.
func (v *V1) AdaptSegment(cell, seg int, prevActive []bool) {
	// delete queue for synapses
	//synQueue := make([]int, 0)

	var perm float32
	for i := range v.cells[cell].segments[seg].synapses {
		perm = v.cells[cell].segments[seg].synapses[i].perm

		// check if synapse.cell is in prevActive
		switch prevActive[v.cells[cell].segments[seg].synapses[i].cell] {
		case true:
			perm += v.P.SynPermActiveMod
		case false:
			perm -= v.P.SynPermInactiveMod
		}

		// constrain perm to [0.0 : 1.0]
		switch {
		case perm < 0.0:
			perm = 0.0
		case perm > 1.0:
			perm = 1.0
		}

		v.cells[cell].segments[seg].synapses[i].perm = perm
	}
}

// PunishSegment will weaken synapses on a segment that are connected
// to any cells active during the previous time step.
func (v *V1) PunishSegment(cell, seg int, prevActive []bool) {
	var perm float32
	for i := range v.cells[cell].segments[seg].synapses {
		perm = v.cells[cell].segments[seg].synapses[i].perm

		// check if synapse.cell is in prevActive
		if prevActive[v.cells[cell].segments[seg].synapses[i].cell] {
			perm -= v.P.SynPermInactiveMod
		}

		// constrain perm to [0.0 : 1.0]
		switch {
		case perm < 0.0:
			perm = 0.0
		case perm > 1.0:
			perm = 1.0
		}

		v.cells[cell].segments[seg].synapses[i].perm = perm
	}
}

// GrowSynapses grows new synapses on a segment to a randomly sampled
// set of cells selected as winners in the previous time step.
func (v *V1) GrowSynapses(cell, seg int, prevWinners []bool) {
	// disable any candidates that are synapsed on this segment
	for _, syn := range v.cells[cell].segments[seg].synapses {
		prevWinners[syn.cell] = false
	}

	// allocate candidates
	candidates := make([]int, 0, 48)
	for i := range prevWinners {
		if prevWinners[i] {
			candidates = append(candidates, i)
		}
	}

	// decide whether to sample or target all winner cells
	var sample []int
	switch {
	case len(candidates) <= v.P.MaxNewSyns:
		// we use all
		sample = rand.Perm(len(candidates))
	case len(candidates) > v.P.MaxNewSyns:
		// we permutate and slice how much we need
		sample = rand.Perm(len(candidates))[:v.P.MaxNewSyns]
	}

	// grow new synapses
	for i := range sample {
		v.CreateSynapse(cell, seg, candidates[sample[i]])
	}
}

// CellsForCol returns a slice of all cell indices within the
// provided column index.
func (v *V1) CellsForCol(col int) []int {
	out := make([]int, 0, v.P.CellsPerCol)
	for i := col * v.P.CellsPerCol; i < (col+1)*v.P.CellsPerCol; i++ {
		out = append(out, i)
	}
	return out
}

// ActiveSegsForCell returns a []int of all active segments
// attached to a cell.
func (v *V1) ActiveSegsForCell(cell int) []int {
	act := make([]int, 0, v.cells[cell].active)
	for i := range v.cells[cell].segments {
		if v.cells[cell].segments[i].active {
			act = append(act, i)
		}
	}
	return act
}

// ActiveSegsForCol returns the number of active segments
// attached to cells in a column.
func (v *V1) ActiveSegsForCol(col int) int {
	cells := v.CellsForCol(col)
	var syns int
	for _, i := range cells {
		syns += v.cells[i].active
	}
	return syns
}

// MatchingSegsForCell returns the number of matching segments
// attached to a cell.
func (v *V1) MatchingSegsForCell(cell int) []int {
	mat := make([]int, 0, v.cells[cell].matching)
	for i := range v.cells[cell].segments {
		if v.cells[cell].segments[i].matching {
			mat = append(mat, i)
		}
	}
	return mat
}

// MatchingSegsForCol returns the number of matching segments
// attached to cells in a column.
func (v *V1) MatchingSegsForCol(col int) int {
	cells := v.CellsForCol(col)
	var syns int
	for _, i := range cells {
		syns += v.cells[i].matching
	}
	return syns
}

// LeastMatchingSegsForCol returns the cell index with the
// least number of matching segments. If there is a tie, a
// random selection is made from the tie candidates.
func (v *V1) LeastSegsForCol(col int) int {
	cells := v.CellsForCol(col)

	// find min value
	min := 9999999
	for i := range cells {
		if len(v.cells[cells[i]].segments) < min {
			min = len(v.cells[cells[i]].segments)
		}
	}

	// tie breaker
	minCells := make([]int, 0)
	for i := range cells {
		if len(v.cells[cells[i]].segments) == min {
			minCells = append(minCells, cells[i])
		}
	}

	choice := rand.Intn(len(minCells))
	return minCells[choice]
}

// BestMatchingSegForCol returns the cell index with the
// matching segment that has the highest number of live
// synapses within the provided volumn. If there is a tie,
// a randoml selection is made from the tie candidates.
func (v *V1) BestMatchingSegForCol(col int) (int, int) {
	cells := v.CellsForCol(col)

	// find max value
	var max int
	for i := range cells {
		for j := range v.cells[cells[i]].segments {
			if v.cells[cells[i]].segments[j].live > max {
				max = v.cells[cells[i]].segments[j].live
			}
		}
	}

	// tie breaker
	maxPairs := make([][2]int, 0, 4)
	for i := range cells {
		for j := range v.cells[cells[i]].segments {
			if v.cells[cells[i]].segments[j].live == max {
				var pair [2]int
				pair[0], pair[1] = cells[i], j
				maxPairs = append(maxPairs, pair)
			}
		}
	}

	choice := rand.Intn(len(maxPairs))

	return maxPairs[choice][0], maxPairs[choice][1]
}

// ComputeActivity computes cell, segment, and synapse activity
// in regards to currently active columns.
func (v *V1) ComputeActivity(active []bool) {
	/*
		for each segment with activity >= ActiveThreshold
		  mark segment active
		for each segment with unconnected activity >= MatchingThreshold
		  mark segment matching
	*/

	act := make([]int, 0, 48)
	for i := range active {
		if active[i] {
			act = append(act, i)
		}
	}

	for i := range v.cells {
		for j := range v.cells[i].segments {
			// count live synapses on each segment
			var live int
			var dead int
			for _, syn := range v.cells[i].segments[j].synapses {
				// if the synapse corresponds to a cell in an active column
				if active[syn.cell] {
					// if synapse is connected
					if syn.perm >= v.P.SynPermConnected {
						live += 1
					} else {
						dead += 1
					}
				}
			}

			// set active / matching
			switch {
			case live >= v.P.ActiveThreshold:
				v.cells[i].segments[j].active = true
				v.cells[i].active += 1
				fallthrough
			case dead >= v.P.MatchingThreshold:
				v.cells[i].segments[j].matching = true
				v.cells[i].matching += 1
			}
			v.cells[i].segments[j].live = live
			v.cells[i].segments[j].dead = dead
		}

	}
}

// Cleanup traverses all cells, segments, and synapses, performing
// maintenance. Segments with 0 synapses and synapses with a
// permanence value of < 0.001 are destroyed.
func (v *V1) Cleanup() {
	for i := range v.cells {
	restartSegs:
		for j := range v.cells[i].segments {
		restartSyns:
			for k := range v.cells[i].segments[j].synapses {
				if v.cells[i].segments[j].synapses[k].perm < 0.001 {
					v.cells[i].segments[j].synapses = append(
						v.cells[i].segments[j].synapses[:k],
						v.cells[i].segments[j].synapses[k+1:]...)
					goto restartSyns
				}
			}

			if len(v.cells[i].segments[j].synapses) == 0 {
				v.cells[i].segments = append(
					v.cells[i].segments[:j],
					v.cells[i].segments[j+1:]...)
				goto restartSegs
			}
		}
	}
}

// Clear clears temporary data from all cells and segments.
func (v *V1) Clear() {
	for i := range v.cells {
		v.cells[i].active = 0
		v.cells[i].matching = 0
		for j := range v.cells[i].segments {
			v.cells[i].segments[j].active = false
			v.cells[i].segments[j].matching = false
		}
	}
}

func (v *V1) Counts() {
	var nSegs, nSyns int

	for i := range v.cells {
		for j := range v.cells[i].segments {
			nSegs++
			nSyns += len(v.cells[i].segments[j].synapses)
		}
	}

	v.nSegs, v.nSyns = nSegs, nSyns
}

// StartNewIteration increments the iteration counter.
func (v *V1) StartNewIteration() {
	v.iteration += 1
}
