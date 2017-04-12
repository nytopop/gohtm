/*
Package cla provides implementation agnostic classifiers
for gohtm. A classifier is neccessary for conversions from
the TM algorithm's internal state to a real valued output.

There are currently two classifiers (both in progress).

V1 directly associates column activity patterns with bucket
values from an encoder; on calling infer with a pattern of
depolarized cells, V1 will do a simple search through all
stored patterns and return the top 4 indices by overlap.

V2 is somewhat more advanced, it uses a feedforward ANN to
classify column patterns more reliably than V2 can.
*/
package cla

// Classifier interface for classifiers to implement.
type Classifier interface {
	Compute(
		sdr []int, bidx int, actValue float64,
		learn, infer bool) Result

	// TODO
	// compute::encoderBucket, activeCells
	// infer::depolarizedCells -> Results
}

// Result should be returned by all classifiers. P represents
// a probability value.
type Result struct {
	P float64
}

// Results type for sorting by probability.
type Results []Result

func (r Results) Len() int           { return len(r) }
func (r Results) Less(i, j int) bool { return r[i].P > r[j].P }
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
