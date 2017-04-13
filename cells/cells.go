/*
Package cells provides an implementation agnostic
interface for the storage and manipulation of synaptic
state.
*/
package cells

// Interface for cellular connectivity structures.
type Interface interface {
	CreateSegment(cell int, targets []bool, perm float32) int
	ComputeActivity(active []bool, perm float32) []int
}

// Cells is an interface for TemporalMemory
// compatible cellular and synaptic state.
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
	StartNewIteration()

	ComputePredictedCols() []bool
	ComputeStats() (int, int)
}
