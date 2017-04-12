/*
Package tm provides an implementation agnostic
interface for temporal memory in gohtm.
*/
package tm

// TemporalMemory is an interface defining the functionality of
// a temporal memory region with no feedback.

// TemporalMemory is an interface for a
// temporal memory region with no feedback.
type TemporalMemory interface {
	Compute(active []bool, learn bool)
	Reset()
	GetActiveCells() []int
	GetAnomalyScore() (anomaly float64)
	GetPrediction() (prediction []bool)
	GetStats() (segments, synapses int)
}

// TemporalMemoryFB is an interface for a
// temporal memory region capable of utilizing
// feedforward input as well as a
type TemporalMemoryFB interface {
	// hmmm, is ETM only for layer 4?

	// as in,

	// l3 learns distal::sensor
	// l4 learns basal::l3cells, apical::motor/efference
}

/* package layer

type Temporal interface // l2, l3
type SensoriMotor interface

*/

// Feedforward + feedback (sequence memory)
//   0.7 internal distal segments (feedforward)
//   0.3 external distal segments (feedback)

// we need to make predictions based on
// (lower state, higher predicted)

// normally, we learn from active/prevActive local cells

// activate feedforward from active cols
// tx from Cells::active-0 to Cells::active-1
// learn transitions from t-1 -> t-0 segment by segment

// ? two learn passes perhaps
// for distal,
//   activeCells[-1] -> activeCells[0]
// for apical,
//   predictCells[-1] -> activeCells[0]

// activate feedforward from active cols
// + last timesteps higher depolarized
// + this timesteps local active
// + last timesteps local depolarized
// tx from [la, ld, ha, hd]
// feedback[
// compare (cells::active-0, Ext::predicted-0 : Cells::active-1 )
