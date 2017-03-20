// Package sp provides an implementation agnostic interface for
// spatial pooling.
package sp

/* Spatial Pooler

The spatial pooler learns and extracts stable feature vectors from a binary
input space. Each cell in a spatial pooler acts as a feature learning network
itself, and the pooler output is the result of a competitive activation / inhibitory process between cells.
*/

// SpatialPooler ...
type SpatialPooler interface {
	Compute(input []bool, learn bool) []bool
}
