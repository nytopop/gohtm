package cells

import (
	"fmt"
	"math/rand"
)

// V1Params contains parameters for
// initialization of V1 cellular state.
type V1Params struct {
	NumColumns  int `json:"numcolumns"`
	CellsPerCol int `json:"cellspercol"`
	SegsPerCell int `json:"segspercell"`
	SynsPerSeg  int `json:"synsperseg"`
}

// V1 cellular state interface.
type V1 struct {
	P         V1Params `json:"params"`
	Cells     []V1Cell `json:"cells"`
	iteration int
}

// NewV1 returns a new V1 instance, initialized
// according to provided params.
func NewV1(p V1Params) *V1 {
	return &V1{
		P:         p,
		Cells:     make([]V1Cell, p.NumColumns*p.CellsPerCol),
		iteration: 0,
	}
}

// V1Cell represents a single cell.
type V1Cell struct {
	Segments         []V1Segment `json:"segments"`
	active, matching int
}

// V1Segment represents a single dendritic
// segment attached to a cell.
type V1Segment struct {
	Synapses         []V1Synapse `json:"synapses"`
	active, matching bool
	live, dead       int
	lastIter         int
}

// V1Synapse represents the connection from the
// dendritic segment of a cell to another cell,
// called the presynaptic cell.
type V1Synapse struct {
	Cell int     `json:"cell"`
	Perm float32 `json:"perm"`
}

// CreateSegment spawns a new segment on the specified
// cell. If the number of segments that already exist
// on the cell is greater than SegsPerCell, the oldest
// or least used segment is removed and the new segment
// appended.
func (v *V1) CreateSegment(cell int) int {
	if len(v.Cells[cell].Segments) < v.P.SegsPerCell {
		v.Cells[cell].Segments = append(
			v.Cells[cell].Segments,
			V1Segment{
				active:   false,
				matching: false,
				Synapses: make([]V1Synapse, 0)})
	} else {
		// TODO :: how do we handle this?
		//         remove the oldest / least recent segment
		// sort by lastIter, then yadda yadda
		fmt.Println("too many segments on cell", cell)
	}

	return len(v.Cells[cell].Segments) - 1
}

// CreateSynapse spawns a new synapse on the specified
// segment / cell. If the number of synapses that already
// exist on the segment exceeds SynsPerSeg, the synapse
// with the lowest permanence value is removed and the
// new synapse is appended.
func (v *V1) CreateSynapse(cell, seg, target int, perm float32) {
	for i := range v.Cells[cell].Segments[seg].Synapses {
		if v.Cells[cell].Segments[seg].Synapses[i].Cell == target {
			return
		}
	}

	if len(v.Cells[cell].Segments[seg].Synapses) >= v.P.SynsPerSeg {
		var idx int
		var min float32

		min = 1.2 // over 1.0
		for i := range v.Cells[cell].Segments[seg].Synapses {
			if v.Cells[cell].Segments[seg].Synapses[i].Perm < min {
				min = v.Cells[cell].Segments[seg].Synapses[i].Perm
				idx = i
			}
		}

		// now we remove idx from the synapse slice
		v.Cells[cell].Segments[seg].Synapses = append(
			v.Cells[cell].Segments[seg].Synapses[:idx],
			v.Cells[cell].Segments[seg].Synapses[idx+1:]...)

	}

	// append the new synapse
	v.Cells[cell].Segments[seg].Synapses = append(
		v.Cells[cell].Segments[seg].Synapses,
		V1Synapse{
			Cell: target,
			Perm: perm})
}

// AdaptSegment adapts synapses on a segment to the provided
// slice of cells active in the previous time step.
func (v *V1) AdaptSegment(cell, seg int, prevActive []bool,
	inc, dec float32) {

	var perm float32
	for i := range v.Cells[cell].Segments[seg].Synapses {
		perm = v.Cells[cell].Segments[seg].Synapses[i].Perm

		// check if synapse.Cell is in prevActive
		switch prevActive[v.Cells[cell].Segments[seg].Synapses[i].Cell] {
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

		v.Cells[cell].Segments[seg].Synapses[i].Perm = perm
	}
}

// GrowSynapses grows new synapses on a segment to a randomly sampled
// set of cells selected as winners in the previous time step.
func (v *V1) GrowSynapses(cell, seg int, prevWinners []bool,
	perm float32, newSyns int) {

	// disable any candidates that are synapsed on this segment
	for _, syn := range v.Cells[cell].Segments[seg].Synapses {
		prevWinners[syn.Cell] = false
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
	act := make([]int, 0, v.Cells[cell].active)
	for i := range v.Cells[cell].Segments {
		if v.Cells[cell].Segments[i].active {
			act = append(act, i)
		}
	}
	return act
}

// ActiveSegsForCol returns the number of active segments
// attached to cells in a column.
func (v *V1) ActiveSegsForCol(col int) int {
	cells := v.CellsForCol(col)
	var segs int
	for _, i := range cells {
		segs += v.Cells[i].active
	}
	return segs
}

// MatchingSegsForCell returns the number of matching segments
// attached to a cell.
func (v *V1) MatchingSegsForCell(cell int) []int {
	mat := make([]int, 0, v.Cells[cell].matching)
	for i := range v.Cells[cell].Segments {
		if v.Cells[cell].Segments[i].matching {
			mat = append(mat, i)
		}
	}
	return mat
}

// MatchingSegsForCol returns the number of matching segments
// attached to cells in a column.
func (v *V1) MatchingSegsForCol(col int) int {
	cells := v.CellsForCol(col)
	var segs int
	for _, i := range cells {
		segs += v.Cells[i].matching
	}
	return segs
}

// LeastSegsForCol returns the cell index with the least number
// of matching segments. If there is a tie, a random selection
// is made from the tie candidates.
func (v *V1) LeastSegsForCol(col int) int {
	cells := v.CellsForCol(col)

	// find min value
	min := 9999999
	for i := range cells {
		if len(v.Cells[cells[i]].Segments) < min {
			min = len(v.Cells[cells[i]].Segments)
		}
	}

	// tie breaker
	//	minCells := make([]int, 0)
	var minCells []int
	for i := range cells {
		if len(v.Cells[cells[i]].Segments) == min {
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
		for j := range v.Cells[cells[i]].Segments {
			if v.Cells[cells[i]].Segments[j].live > max {
				max = v.Cells[cells[i]].Segments[j].live
			}
		}
	}

	// tie breaker
	maxPairs := make([][2]int, 0, 4)
	for i := range cells {
		for j := range v.Cells[cells[i]].Segments {
			if v.Cells[cells[i]].Segments[j].live == max {
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

	for i := range v.Cells {
		for j := range v.Cells[i].Segments {
			// count live synapses on each segment
			var live int
			var dead int
			for _, syn := range v.Cells[i].Segments[j].Synapses {
				// if synapse corresponds to
				// a cell in an active column
				if active[syn.Cell] {
					// if synapse is connected
					if syn.Perm >= connected {
						live++
					} else {
						dead++
					}
				}
			}

			// set active / matching
			switch {
			case live >= activeThreshold:
				v.Cells[i].Segments[j].active = true
				v.Cells[i].active++
				fallthrough
			case dead >= matchThreshold:
				v.Cells[i].Segments[j].matching = true
				v.Cells[i].matching++
			}
			v.Cells[i].Segments[j].live = live
			v.Cells[i].Segments[j].dead = dead
		}
	}
}

// Cleanup traverses all cells, segments, and synapses, performing
// maintenance. Segments with 0 synapses and synapses with a
// permanence value of < 0.001 are destroyed.
func (v *V1) Cleanup() {
	for i := range v.Cells {
	restartSegs:
		for j := range v.Cells[i].Segments {
		restartSyns:
			for k := range v.Cells[i].Segments[j].Synapses {
				if v.Cells[i].Segments[j].Synapses[k].Perm < 0.001 {
					v.Cells[i].Segments[j].Synapses = append(
						v.Cells[i].Segments[j].Synapses[:k],
						v.Cells[i].Segments[j].Synapses[k+1:]...)
					goto restartSyns
				}
			}

			if len(v.Cells[i].Segments[j].Synapses) == 0 {
				v.Cells[i].Segments = append(
					v.Cells[i].Segments[:j],
					v.Cells[i].Segments[j+1:]...)
				goto restartSegs
			}
		}
	}
}

// Clear clears temporary data from all cells and segments.
func (v *V1) Clear() {
	for i := range v.Cells {
		v.Cells[i].active = 0
		v.Cells[i].matching = 0
		for j := range v.Cells[i].Segments {
			v.Cells[i].Segments[j].active = false
			v.Cells[i].Segments[j].matching = false
			v.Cells[i].Segments[j].live = 0
			v.Cells[i].Segments[j].dead = 0
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
			if v.Cells[j].active > 0 {
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
	for i := range v.Cells {
		nSegs += len(v.Cells[i].Segments)
		for j := range v.Cells[i].Segments {
			nSyns += len(v.Cells[i].Segments[j].Synapses)
		}
	}
	return nSegs, nSyns
}
