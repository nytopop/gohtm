/*
Package sp provides an implementation agnostic interface
for spatial poolers.

The spatial pooler learns and extracts stable feature
vectors from a binary input space. Each cell in a
spatial pooler acts as a feature learning network itself,
and the pooler output is the result of a competitive
activation / inhibition process between cells.

There are two spatial pooler implementations working right
now. Use V2, which fixes some architectural mistakes from V1.
*/
package sp

// SpatialPooler ...
type SpatialPooler interface {
	Compute(input []bool, learn bool) []bool
}
