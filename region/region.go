// Package region provides an implementation agnostic interface for linking
// and interacting with computational units of the htm algorithm.
package region

// Region is an interface providing access to a segment of cortex.
type Region interface {
	Compute(input []bool)
	Serialize() []byte
}
