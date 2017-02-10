// Package tm provides an implementation agnostic interface for temporal
// memory in the htm algorithm.
package tm

// TemporalMemory is an interface defining the functionality of
// a temporal memory region.
type TemporalMemory interface {
	Compute(active []bool, learn bool)
	Reset()
	GetAnomaly()
	GetPrediction() []bool
	GetState()
}
