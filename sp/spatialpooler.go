// Package sp provides an implementation agnostic interface for
// spatial pooling.
package sp

type SpatialPooler interface {
	Compute(input []bool, learn bool) []bool
}

/*
go with the V1, V2, V3, etc naming scheme

enc: Encoder
scalar, rdse, retina, audio

sp: SpatialPooler
v1, v2, etc

tm: TemporalMemory
v1, v2, etc

region: Region
v1, v2, etc

cells: Cells
v1, v2, etc
*/
