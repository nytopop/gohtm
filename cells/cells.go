// Package cells provides an implementation agnostic interface for storage
// and manipulation of cellular and synaptic state.
package cells

// Cells is an interface for TemporalMemory compatible cellular and synaptic
// state.
type Cells interface {
	CreateSegment(cell int)              // TODO
	DestroySegment(cell, seg int)        // TODO
	CreateSynapse(cell, seg, target int) // TODO
	DestroySynapse(cell, seg, syn int)   // TODO

	AdaptSynapses(cell int, prevActive []bool)
	GrowSynapses(cell int, prevWinners []bool) // TODO

	CellsForCol(col int) []int
	ActiveSegsForCell(cell int) int
	ActiveSegsForCol(col int) int
	MatchingSegsForCell(cell int) int
	MatchingSegsForCol(col int) int

	ComputeActivity(active []bool)
	Clear()
	StartNewIteration()
}
