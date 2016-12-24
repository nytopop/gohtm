package gohtm

import "fmt"

/* Temporal Memory */

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

// TemporalMemory is a sequence learning and prediction algorithm.
type TemporalMemory struct {
	// state
	cons         *Connections
	activeCells  []int
	winnerCells  []int
	activeSegs   []int
	matchingSegs []int

	// params
	numColumns          int
	numCells            int
	activationThreshold int
	minThreshold        int
	initPermanence      float64
	synPermConnected    float64
	synPermActiveMod    float64
	synPermInactiveMod  float64
	segPerCell          int
	synPerSeg           int
}

// NewTemporalMemory initializes a new TemporalMemory region with
// the provided TemporalParams.
func NewTemporalMemory(p TemporalParams) TemporalMemory {
	tm := TemporalMemory{
		numColumns:          p.NumColumns,
		numCells:            p.NumCells,
		activationThreshold: p.ActivationThreshold,
		minThreshold:        p.MinThreshold,
		initPermanence:      p.InitPermanence,
		synPermConnected:    p.SynPermConnected,
		synPermActiveMod:    p.SynPermActiveMod,
		synPermInactiveMod:  p.SynPermInactiveMod,
		segPerCell:          p.SegPerCell,
		synPerSeg:           p.SynPerSeg,
	}

	totalCells := tm.numColumns * tm.numCells
	//	tm.cons = NewConnections(totalCells, tm.segPerCell, tm.synPerSeg)
	tm.cons = NewConnections(tm.numColumns, tm.numCells, tm.segPerCell, tm.synPerSeg)
	tm.activeCells = make([]int, totalCells)
	tm.winnerCells = make([]int, totalCells)

	return tm
}

// Compute iterates the TemporalMemory algorithm with the provided
// vector of active columns from a SpatialPooler.
func (tm *TemporalMemory) Compute(active SparseBinaryVector, learn bool) SparseBinaryVector {
	// tm.cons.Snapshot
	tm.activateCells(active.Dense(), learn)
	tm.activateDendrites(learn)

	// TODO : switch to slices for synpatic modifications

	return SparseBinaryVector{}
}

// Reset clears temporary data so sequences are not learned between
// the current and next time step.
func (tm *TemporalMemory) Reset() {
}

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
		} /*else {
			// terrible prediction
			if len(predicted) > 0 {
				fmt.Println("punishing", i)
				tm.punishColumn(i)
			}
		}*/
	}
}

// punishColumn does something??
func (tm *TemporalMemory) punishColumn(col int) {
}

/*
// activateCells does something
func (tm *TemporalMemory) activateCells(active []bool, learn bool) {
	prevActiveCells := tm.activeCells
	prevWinnerCells := tm.winnerCells
	tm.activeCells = make([]int, tm.numColumns*tm.numCells)
	tm.winnerCells = make([]int, tm.numColumns*tm.numCells)

	for i, col := range active {
		if col {
			for _, cell := range tm.cons.CellsForCol(i) {
				if cell.state == 1 {
					tm.activatePredictedColumn(i, prevActiveCells, prevWinnerCells, learn)
					goto predicted
				}
			}
			tm.burstColumn(i)
		predicted:
		} else {
			if learn {
				for _, cell := range tm.cons.CellsForCol(i) {
					if cell.state == 1 {
						tm.punishPredictedColumn(i)
						break
					}
				}
			}
		}

		// yay we've finished... or not?
	}

}
*/

/*
// column, prevactivecells, prevwinnercells
// learn
func (tm *TemporalMemory) activatePredictedColumn(col int, prevActive, prevWinner []int, learn bool) {
	fmt.Println("Activating column", col)
	for _, cell := range tm.cons.CellsForCol(col) {
		if cell.state == 1 {
			cell.state = 2
			break
		}
	}

	// TODO
	if learn {
		// strengthen active synapses
		// weaken inactive synapses
		// grow synapses to previous winner cells
	}
}
*/

/*
func (tm *TemporalMemory) burstColumn(col int) {
	fmt.Println("Bursting column", col)
	cells := tm.cons.CellsForCol(col)
	for _, cell := range cells {
		cell.state = 2
	}

	for _, cell := range cells {
		// hoow to find active distal dendrites?
		fmt.Println(cell)
	}
}

func (tm *TemporalMemory) punishPredictedColumn(col int) {
}

func (tm *TemporalMemory) activateDendrites(learn bool) {
}
*/
