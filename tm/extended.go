package tm

import (
	"github.com/nytopop/gohtm/cells"
	"github.com/nytopop/gohtm/vec"
)

// TemporalParams contains parameters for initialization of a
// TemporalMemory region.
type TemporalParams struct {
	NumColumns          int     // input space dimensions
	NumCells            int     // cells per column
	ActivationThreshold int     // # of active synapses for 'active'
	MinThreshold        int     // # of potential synapses for learning
	InitPermanence      float64 // initial permanence of new synapses
	SynPermConnected    float64 // connection threshold
	SynPermInactiveMod  float64 // perm decrement
	SynPermActiveMod    float64 // perm increment
	SegPerCell          int
	SynPerSeg           int
}

// NewTemporalParams returns a default TemporalParams.
func NewTemporalParams() TemporalParams {
	return TemporalParams{
		NumColumns:          2048,
		NumCells:            32,
		ActivationThreshold: 12,
		MinThreshold:        10,
		InitPermanence:      0.21,
		SynPermConnected:    0.5,
		SynPermActiveMod:    0.05,
		SynPermInactiveMod:  0.03,
		SegPerCell:          16,
		SynPerSeg:           16,
	}
}

// ExtendedTemporalMemory is a sequence learning and prediction algorithm.
type ExtendedTemporalMemory struct {
	P    TemporalParams
	Cons *cells.Connections
}

// NewExtendedTemporalMemory initializes a new TemporalMemory region with
// the provided TemporalParams.
func NewExtendedTemporalMemory(p TemporalParams) *ExtendedTemporalMemory {
	return &ExtendedTemporalMemory{
		P: p,
		Cons: cells.NewConnections(
			p.NumColumns,
			p.NumCells,
			p.SegPerCell,
			p.SynPerSeg,
		),
	}
}

// Compute iterates the TemporalMemory algorithm with the provided
// vector of active columns from a SpatialPooler.
func (tm *ExtendedTemporalMemory) Compute(active vec.SparseBinaryVector, learn bool) vec.SparseBinaryVector {
	// tm.cons.Snapshot??
	//tm.activateCells(active.Dense(), learn)
	//tm.activateDendrites(learn)

	return vec.SparseBinaryVector{}
}

// Reset clears temporary data so sequences are not learned between
// the current and next time step.
func (tm *ExtendedTemporalMemory) Reset() {
}

/*
// activateDendrites : TM 2
func (tm *TemporalMemory) activateDendrites(learn bool) {
	// loop cells
	for _, cell := range tm.cons.cells {
		// loop dendrites
		var activeSegs = 0
		for _, seg := range cell.segments {
			// count synapses connected to active cells
			// for each synapse, check if presynaptic is active
			var activeSyns = 0
			for syn := range seg.synapses {
				c1 := syn.perm >= tm.synPermConnected
				c2 := tm.cons.cells[syn.preSynapticCell].state == 2
				if c1 && c2 {
					activeSyns++
				}
			}

			// activate segment if over threshold
			if activeSyns >= tm.activationThreshold {
				seg.state = true

				// modify synapses
				for syn := range seg.synapses {
					if tm.cons.cells[syn.preSynapticCell].state == 2 {
						tm.cons.tempUpdateSynPerm(syn, tm.synPermActiveMod)
					} else {
						tm.cons.tempUpdateSynPerm(syn, tm.synPermInactiveMod)
					}
				}

				activeSegs++
			} else {
				seg.state = false
			}
		}

		switch {
		// deactivate cell if no active dendrites and not active
		case activeSegs == 0 && cell.state != 2:
			cell.state = 0

		// depolarize cell if active dendrites and not active
		case activeSegs != 0 && cell.state != 2:
			cell.state = 1
			// we
		}
	}
}
*/
/*
// activateCells : TM 1
func (tm *TemporalMemory) activateCells(feedforward []bool, learn bool) {
	for i, col := range feedforward {
		if col {
			predicted := tm.cons.PredictedForCol(i)
			switch {
			case len(predicted) == 0:
				// burst
				fmt.Println("bursting", i)
				for _, cell := range tm.cons.CellsForCol(i) {
					cell.state = 2
				}

			case len(predicted) > 0:
				// activate predicted cells
				fmt.Println("activating", i)
				for _, cell := range predicted {
					cell.state = 2
				}
			}
		}// else {
//			// terrible prediction
//			if len(predicted) > 0 {
//				fmt.Println("punishing", i)
//				tm.punishColumn(i)
//			}
//		}
	}
}

// punishColumn does something??
func (tm *TemporalMemory) punishColumn(col int) {
}
*/
