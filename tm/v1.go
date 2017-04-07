package tm

import (
	"fmt"

	"github.com/nytopop/gohtm/cells"
	"github.com/nytopop/gohtm/vec"
)

/* TODO

- create access methods for deltaInf, deltaAcc, and deltaAnm

*/

// V1Params contains parameters for initialization of a V1
// TemporalMemory region.
type V1Params struct {
	NumColumns  int // input space dimensions
	CellsPerCol int // cells per column
	SegsPerCell int // max segments per cell
	SynsPerSeg  int // max synapses per segment

	InitPerm         float32 // initial permanence of new synapses
	SynPermConnected float32 // threshold for synapse to be deemed connected
	SynPermLearnMod  float32 // permanence mod during learning
	SynPermPunishMod float32 // permanence mod for incorrect prediction

	MaxNewSyns      int // max new synapses per segment per iteration
	ActiveThreshold int // threshold for segment to turn active
	MatchThreshold  int // threshold for segment to turn matching
}

// NewV1Params returns a default V1Params.
func NewV1Params() V1Params {
	return V1Params{
		NumColumns:       2048,
		CellsPerCol:      32,
		SegsPerCell:      32,
		SynsPerSeg:       32,
		InitPerm:         0.21,
		SynPermConnected: 0.5,
		SynPermLearnMod:  0.1,
		SynPermPunishMod: 0.01,
		MaxNewSyns:       16,
		ActiveThreshold:  12,
		MatchThreshold:   10,
	}
}

// V1 is a basic implementation of TemporalMemory.
type V1 struct {
	// Params
	P V1Params

	// State
	Cons            cells.Cells
	PrevActiveCells []bool
	PrevWinnerCells []bool
	ActiveCells     []bool
	WinnerCells     []bool
	prediction      []bool
	iteration       int

	// Metrics
	deltaInf     float64 // depolarized : active
	deltaAcc     float64 // correct : incorrect
	deltaAnm     float64 // (active-correct) : active
	nSegs, nSyns int
}

// NewV1 initializes a new TemporalMemory region with the provided V1Params.
func NewV1(p V1Params) *V1 {
	cp := cells.V1Params{
		NumColumns:  p.NumColumns,
		CellsPerCol: p.CellsPerCol,
		SegsPerCell: p.SegsPerCell,
		SynsPerSeg:  p.SynsPerSeg,
	}

	return &V1{
		P:               p,
		Cons:            cells.NewV1(cp),
		PrevActiveCells: make([]bool, 0, p.NumColumns*p.CellsPerCol),
		PrevWinnerCells: make([]bool, 0, p.NumColumns*p.CellsPerCol),
		ActiveCells:     make([]bool, 0, p.NumColumns*p.CellsPerCol),
		WinnerCells:     make([]bool, 0, p.NumColumns*p.CellsPerCol),
		prediction:      make([]bool, p.NumColumns),
		iteration:       0,
	}
}

// Compute iterates the TemporalMemory algorithm with the
// provided vector of active columns from a SpatialPooler.
func (e *V1) Compute(active []bool, learn bool) {
	if len(active) != e.P.NumColumns {
		panic("TM: Mismatched input dimensions")
	}
	// We compute metrics by taking prediction from last step
	// and comparing to currently active columns
	e.computeMetrics(active)

	// Compute active / depolarized cells
	e.activateCells(active, learn)

	// Compute active & matching dendrite segments
	e.Cons.Clear()
	e.Cons.ComputeActivity(e.ActiveCells, e.P.SynPermConnected,
		e.P.ActiveThreshold, e.P.MatchThreshold)

	// Cleanup neural net
	e.Cons.Cleanup()

	// Compute prediction, stats
	e.prediction = e.Cons.ComputePredictedCols()
	e.nSegs, e.nSyns = e.Cons.ComputeStats()

	if learn {
		e.iteration++
		e.Cons.StartNewIteration()
	}

	fmt.Println("deltaInf:", e.deltaInf)
	fmt.Println("deltaAcc:", e.deltaAcc)
	fmt.Println("deltaAnm:", e.deltaAnm)
}

// Calculate the active cells using active columns and dendrite segments.
// Grow and reinforce synapses.
func (e *V1) activateCells(active []bool, learn bool) {
	/*
		for each column
		  if col active and has active dendrite segments
		    call activatePredictedColumn
		  if col active and doesn't have active dendrite segments
		    call burstColumn
		  if col inactive and has matching dendrite segments
		    call punishpredictedcolumn
	*/

	// Allocate memory for active / winner cells
	e.PrevActiveCells = e.ActiveCells
	e.PrevWinnerCells = e.WinnerCells
	e.ActiveCells = make([]bool, e.P.NumColumns*e.P.CellsPerCol)
	e.WinnerCells = make([]bool, e.P.NumColumns*e.P.CellsPerCol)

	for i := range active {
		if active[i] {
			syns := e.Cons.ActiveSegsForCol(i)
			switch {
			case syns > 0:
				//fmt.Println("Activating predicted col", i)
				cellsToAdd := e.activatePredictedColumn(i, learn)
				for _, c := range cellsToAdd {
					e.ActiveCells[c] = true
					e.WinnerCells[c] = true
				}
			case syns == 0:
				//fmt.Println("Bursting col", i)
				cellsToAdd, winnerCell := e.burstColumn(i, learn)
				for _, c := range cellsToAdd {
					e.ActiveCells[c] = true
				}
				e.WinnerCells[winnerCell] = true
			}
		} else {
			if learn {
				if e.Cons.MatchingSegsForCol(i) > 0 {
					//fmt.Println("Punishing predicted col", i)
					e.punishPredictedColumn(i)
				}
			}
		}
	}
}

func (e *V1) activatePredictedColumn(col int, learn bool) []int {
	/*
		for each cell in col that has an active distal dendrite segment
		  mark cell as active cell
		  mark cell as winner cell
		  if learning enabled
		    strengthen active synapses
			weaken inactive synapses
			grow synapses to previous winner cells
	*/
	cellsToAdd := make([]int, 0, e.P.CellsPerCol) // TODO sizing ???

	for _, i := range e.Cons.CellsForCol(col) {
		act := e.Cons.ActiveSegsForCell(i)
		if len(act) > 0 {
			cellsToAdd = append(cellsToAdd, i)

			if learn {
				for j := range act {
					e.Cons.AdaptSegment(i, act[j], e.PrevActiveCells,
						e.P.SynPermLearnMod, e.P.SynPermLearnMod)
					e.Cons.GrowSynapses(i, act[j], e.PrevWinnerCells,
						e.P.InitPerm, e.P.MaxNewSyns)
				}
			}
		}
	}

	return cellsToAdd
}

func (e *V1) burstColumn(col int, learn bool) ([]int, int) {
	/*
		mark all cells as active
		if any matching segments
			find most active matching segment, mark its cell as winner
			if learn
				grow & reinforce synapse to prevWinnerCells
		if no matching segments
			find cell with least # segments, mark it as winner
			if learn
				if any prevWinnerCells
					add segment to this winner cell
					grow synapses to prevWinnerCells
	*/

	cellsToAdd := e.Cons.CellsForCol(col)
	var winnerCell, winnerSeg int

	segs := e.Cons.MatchingSegsForCol(col)
	switch {
	case segs > 0:
		winnerCell, winnerSeg = e.Cons.BestMatchingSegForCol(col)
		if learn {
			e.Cons.AdaptSegment(winnerCell, winnerSeg, e.PrevActiveCells,
				e.P.SynPermLearnMod, e.P.SynPermLearnMod)
			e.Cons.GrowSynapses(winnerCell, winnerSeg, e.PrevWinnerCells,
				e.P.InitPerm, e.P.MaxNewSyns)
		}
	case segs == 0:
		winnerCell = e.Cons.LeastSegsForCol(col)
		if learn {
			winnerSeg = e.Cons.CreateSegment(winnerCell)
			e.Cons.GrowSynapses(winnerCell, winnerSeg, e.PrevWinnerCells,
				e.P.InitPerm, e.P.MaxNewSyns)
		}
	}

	return cellsToAdd, winnerCell
}

func (e *V1) punishPredictedColumn(col int) {
	/*
		for each matching segment in the column
			weaken active synapses
	*/

	cells := e.Cons.CellsForCol(col)
	for i := range cells {
		segs := e.Cons.MatchingSegsForCell(cells[i])
		for j := range segs {
			e.Cons.AdaptSegment(cells[i], segs[j], e.PrevActiveCells,
				-e.P.SynPermPunishMod, 0.0)
		}
	}
}

// Reset clears temporary data so sequences are not learned between
// the current and next time step.
func (e *V1) Reset() {
	e.Cons.Clear()
	e.PrevActiveCells = make([]bool, e.P.NumColumns*e.P.CellsPerCol)
	e.PrevWinnerCells = make([]bool, e.P.NumColumns*e.P.CellsPerCol)
	e.ActiveCells = make([]bool, e.P.NumColumns*e.P.CellsPerCol)
	e.WinnerCells = make([]bool, e.P.NumColumns*e.P.CellsPerCol)
}

func (e *V1) computeMetrics(active []bool) {
	var nActive, nGood, nBad int
	for i := range active {
		if active[i] {
			nActive++
		}

		if e.prediction[i] {
			switch active[i] {
			case true:
				nGood++
			case false:
				nBad++
			}
		}
	}

	e.deltaInf = float64(nGood+nBad) / float64(nActive)
	e.deltaAcc = float64(nGood) / float64(nBad+nGood)
	e.deltaAnm = float64(nActive-nGood) / float64(nActive)
}

// GetActiveCells returns the currently active cells, in []int
// format.
func (e *V1) GetActiveCells() []int {
	return vec.ToInt(e.ActiveCells)
}

// GetAnomalyScore returns the current normalized anomaly score.
func (e *V1) GetAnomalyScore() float64 {
	return e.deltaAnm
}

// GetPrediction returns the current set of depolarized columns.
func (e *V1) GetPrediction() []bool {
	return e.prediction
}

// GetStats returns the current number of segments and synapses.
func (e *V1) GetStats() (int, int) {
	return e.nSegs, e.nSyns
}
