// Package cells provides an implementation agnostic interface for storage
// and manipulation of cellular and synaptic state.
package cells

// Cells is an interface for TemporalMemory compatible cellular and synaptic
// state.
type Cells interface {
	CreateSegment(cell int) int                        // TODO
	CreateSynapse(cell, seg, target int, perm float32) // TODO

	AdaptSegment(cell, seg int, prevActive []bool,
		inc, dec float32)
	GrowSynapses(cell, seg int, prevWinners []bool,
		perm float32, newSyns int)

	CellsForCol(col int) []int
	ActiveSegsForCell(cell int) []int
	ActiveSegsForCol(col int) int
	MatchingSegsForCell(cell int) []int
	MatchingSegsForCol(col int) int
	LeastSegsForCol(col int) int
	BestMatchingSegForCol(col int) (int, int)

	ComputeActivity(active []bool, connected float32,
		activeThreshold, matchThreshold int)
	Cleanup()
	Clear()
	Counts()
	StartNewIteration()
}
