package cells

import (
	"fmt"
	"math/rand"
)

/* TODO

 */

// V1Params contains parameters for initialization of V1 cellular state.
type V1Params struct {
	NumColumns  int
	CellsPerCol int
	SegsPerCell int
	SynsPerSeg  int
}

// V1 cellular state interface.
type V1 struct {
	P         V1Params
	cells     []V1Cell
	iteration int
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
	active, matching int
	segments         []V1Segment
}

// V1Segment represents a single dendritic segment attached to a cell.
type V1Segment struct {
	active, matching bool
	live, dead       int
	lastIter         int
	synapses         []V1Synapse
}

// V1Synapse represents the connection from the dendritic segment of a cell
// to another cell, called the presynaptic cell.
type V1Synapse struct {
	cell int
	perm float32
}

// CreateSegment spawns a new segment on the specified cell. If the number
// of segments that already exist on the cell is greater than SegsPerCell,
// the oldest or least used segment is removed and the new segment appended.
func (v *V1) CreateSegment(cell int) int {
	if len(v.cells[cell].segments) < v.P.SegsPerCell {
		v.cells[cell].segments = append(
			v.cells[cell].segments,
			V1Segment{
				active:   false,
				matching: false,
				synapses: make([]V1Synapse, 0)})
	} else {
		// TODO :: how do we handle this?
		//         remove the oldest / least recent segment
		fmt.Println("too many segments on cell", cell)
	}

	return len(v.cells[cell].segments) - 1
}

// CreateSynapse spawns a new synapse on the specified segment / cell. If the
// number of synapses that already exist on the segment exceeds SynsPerSeg,
// the synapse with the lowest permanence value is removed and the new
// synapse is appended.
func (v *V1) CreateSynapse(cell, seg, target int, perm float32) {
	for i := range v.cells[cell].segments[seg].synapses {
		if v.cells[cell].segments[seg].synapses[i].cell == target {
			return
		}
	}

	if len(v.cells[cell].segments[seg].synapses) >= v.P.SynsPerSeg {
		var idx int
		var min float32

		min = 1.2 // over 1.0
		for i := range v.cells[cell].segments[seg].synapses {
			if v.cells[cell].segments[seg].synapses[i].perm < min {
				min = v.cells[cell].segments[seg].synapses[i].perm
				idx = i
			}
		}

		// now we remove idx from the synapse slice
		v.cells[cell].segments[seg].synapses = append(
			v.cells[cell].segments[seg].synapses[:idx],
			v.cells[cell].segments[seg].synapses[idx+1:]...)

	}

	// append the new synapse
	v.cells[cell].segments[seg].synapses = append(
		v.cells[cell].segments[seg].synapses,
		V1Synapse{
			cell: target,
			perm: perm})
}

// AdaptSegment adapts synapses on a segment to the provided slice
// of cells active in the previous time step.
func (v *V1) AdaptSegment(cell, seg int, prevActive []bool,
	inc, dec float32) {

	var perm float32
	for i := range v.cells[cell].segments[seg].synapses {
		perm = v.cells[cell].segments[seg].synapses[i].perm

		// check if synapse.cell is in prevActive
		switch prevActive[v.cells[cell].segments[seg].synapses[i].cell] {
		case true:
			perm += inc
		case false:
			perm -= dec
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
func (v *V1) GrowSynapses(cell, seg int, prevWinners []bool,
	perm float32, newSyns int) {

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
	case len(candidates) <= newSyns:
		// we use all
		sample = rand.Perm(len(candidates))
	case len(candidates) > newSyns:
		// we permutate and slice how much we need
		sample = rand.Perm(len(candidates))[:newSyns]
	}

	// grow new synapses
	for i := range sample {
		v.CreateSynapse(cell, seg, candidates[sample[i]], perm)
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

// LeastSegsForCol returns the cell index with the least number
// of matching segments. If there is a tie, a random selection
// is made from the tie candidates.
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
	//	minCells := make([]int, 0)
	var minCells []int
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
func (v *V1) ComputeActivity(active []bool, connected float32,
	activeThreshold, matchThreshold int) {
	/*
		for each segment with activity >= ActiveThreshold
		  mark segment active
		for each segment with unconnected activity >= MatchingThreshold
		  mark segment matching
	*/

	for i := range v.cells {
		for j := range v.cells[i].segments {
			// count live synapses on each segment
			var live int
			var dead int
			for _, syn := range v.cells[i].segments[j].synapses {
				// if the synapse corresponds to a cell in an active column
				if active[syn.cell] {
					// if synapse is connected
					if syn.perm >= connected {
						live++
					} else {
						dead++
					}
				}
			}

			// set active / matching
			switch {
			case live >= activeThreshold:
				v.cells[i].segments[j].active = true
				v.cells[i].active++
				fallthrough
			case dead >= matchThreshold:
				v.cells[i].segments[j].matching = true
				v.cells[i].matching++
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
			v.cells[i].segments[j].live = 0
			v.cells[i].segments[j].dead = 0
		}
	}
}

// StartNewIteration increments the iteration counter.
func (v *V1) StartNewIteration() {
	v.iteration++
}

// ComputePredictedCols computes which columns contain
// depolarized cells and returns them.
func (v *V1) ComputePredictedCols() []bool {
	prediction := make([]bool, v.P.NumColumns)
	for i := range prediction {
		for _, j := range v.CellsForCol(i) {
			if v.cells[j].active > 0 {
				prediction[i] = true
				break
			}
		}
	}
	return prediction
}

// ComputeStats returns the total number of segments and synapses.
func (v *V1) ComputeStats() (int, int) {
	var nSegs, nSyns int
	for i := range v.cells {
		nSegs += len(v.cells[i].segments)
		for j := range v.cells[i].segments {
			nSyns += len(v.cells[i].segments[j].synapses)
		}
	}
	return nSegs, nSyns
}
