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

/* Requirements
We should never deal with internal datastructures from the calling code.
All input/output should use generic types such as int / float64 / etc.
*/
