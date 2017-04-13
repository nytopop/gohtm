/*
Package tm provides an implementation agnostic
interface for temporal memory in gohtm.
*/
package tm

// Interface for temporal memory. Updated API to allow the
// option of feedback input. For no feedback, simply pass
// a 0 length apical input.
type Interface interface {
	Compute(learn bool, cols, basal, apical []bool) error
	Reset()
	ActiveCells() []bool
	WinnerCells() []bool
}

// TemporalMemory is an interface for a temporal
// memory region with no feedback.
type TemporalMemory interface {
	Compute(active []bool, learn bool)
	Reset()
	GetActiveCells() []int
	GetAnomalyScore() float64
	GetPrediction() []bool
	GetStats() (segments, synapses int)
}
