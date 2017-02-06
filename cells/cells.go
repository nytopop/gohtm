// Package cells provides an implementation agnostic interface for storage
// and manipulation of cellular and synaptic state.
package cells

// Cells is an interface for TemporalMemory compatible cellular and synaptic
// state.
type Cells interface {
	CreateSegment(cell int)
	DestroySegment(cell, seg int)
	CreateSynapse(cell, seg, target int)
	DestroySynapse(cell, seg, syn int)
	DepolarizedForCol(col int) []int
	ComputeStatistics(active []bool)
}
