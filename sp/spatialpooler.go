// Package sp provides an implementation agnostic interface for
// spatial pooling.
package sp

// SpatialPooler asdf
type SpatialPooler interface {
	Compute(input []bool, learn bool) []bool
}
