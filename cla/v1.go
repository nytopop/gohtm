package cla

import (
	"sort"

	"github.com/nytopop/gohtm/vec"
)

// V1Params contains parameters for initialization of a V1 Classifier.
type V1Params struct {
}

// NewV1Params returns a default set of parameters for a V1 Classifier.
func NewV1Params() V1Params {
	return V1Params{}
}

// V1 Classifier.
type V1 struct {
	P       V1Params
	entries V1Sortable
}

// NewV1 returns a new V1 Classifier initialized with the provided V1Params.
func NewV1(p V1Params) *V1 {
	return &V1{
		P: p,
	}
}

// V1Entry represents a single entry in the classifier. The Overlap value
// is only populated on return from the (v *V1) Classify() method.
type V1Entry struct {
	Overlap     int // This is only populated on return values
	InputVector []int
	SDR         []int
}
type V1Sortable []V1Entry

func (v V1Sortable) Len() int           { return len(v) }
func (v V1Sortable) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v V1Sortable) Less(i, j int) bool { return v[i].Overlap < v[j].Overlap }

// Associate takes an input of the currently active columns from a
// SpatialPooler and the output vector from an Encoder, then stores
// them in the internal database.
func (c *V1) Associate(active, vector []bool) {
	/* Duplicate detection
	Each SDR should be unique.
	Each InputVector may be associated with multiple SDRs.

	This is legal:
	IV   : SDR
	0100 : 00010010
	0100 : 00010001

	This is _not_ legal:
	IV   : SDR
	0100 : 00010010
	1000 : 00010010

	This means that learning in the spatial pooler will not disrupt
	classifier accuracy. If multiple InputVectors were represented
	by the same SDR, there would be no way to reliably differentiate
	what a particular SDR received from temporal memory was
	representing.
	*/

	// convert []bool to []int
	vectorInts, activeInts := vec.ToInt(vector), vec.ToInt(active)

	// return without storing if active is a duplicate
	for i := range c.entries {
		if vec.Equal(c.entries[i].SDR, activeInts) {
			return
		}
	}

	// append a new entry
	c.entries = append(c.entries, V1Entry{
		InputVector: vectorInts,
		SDR:         activeInts})
}

// Classify searches for any SDRs that overlap with prediction, looks
// up the associated input vectors, and outputs them sorted by overlap
// amount. Only SDRs that overlap are returned.
func (c *V1) Classify(prediction []bool) V1Sortable {
	// convert []bool to []int
	predictionInts := vec.ToInt(prediction)

	// search for overlapping SDRs
	var output V1Sortable
	for i := range c.entries {
		overlap := vec.Overlap(predictionInts, c.entries[i].SDR)
		if overlap > 0 {
			entry := V1Entry{
				Overlap:     overlap,
				InputVector: c.entries[i].InputVector,
				SDR:         c.entries[i].SDR,
			}
			output = append(output, entry)
		}
	}

	// sort and return output
	sort.Sort(output)
	return output
}
