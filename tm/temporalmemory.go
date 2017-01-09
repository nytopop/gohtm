// Package tm provides an implementation agnostic interface for temporal
// memory in the htm algorithm.
package tm

import "github.com/nytopop/gohtm/vec"

// TemporalMemory is an interface defining the functionality of
// a temporal memory region.
type TemporalMemory interface {
	Compute(vec.SparseBinaryVector, bool) vec.SparseBinaryVector
	Reset()
	//Save()
	//Load()
}
