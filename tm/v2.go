package tm

import (
	"fmt"

	"github.com/nytopop/gohtm/cells"
	"github.com/pkg/errors"
)

// V2 ... extended TM. basal, apical
type V2Params struct {
	NumColumns       int     `json:"numcolumns"`
	CellsPerCol      int     `json:"cellspercol"`
	SegsPerCell      int     `json:"segspercell"`
	SynsPerSeg       int     `json:"synsperseg"`
	MatchThreshold   int     `json:"matchthreshold"`
	ActiveThreshold  int     `json:"activethreshold"`
	SynPermConnected float32 `json:"synpermconnected"`

	NumBasalCells  int `json:"numbasalcells"`
	NumApicalCells int `json:"numapicalcells"`
}

func NewV2Params() V2Params {
	return V2Params{
		NumColumns:       2048,
		CellsPerCol:      16,
		SegsPerCell:      32,
		SynsPerSeg:       16,
		MatchThreshold:   6,
		ActiveThreshold:  12,
		SynPermConnected: 0.5,
	}
}

// V2 temporal memory. Implements support for apical disambiguation.
type V2 struct {
	// Params
	P V2Params `json:"params"`

	// State
	Basal           cells.Interface `json:"basal"`
	Apical          cells.Interface `json:"apical"`
	prevActiveCells []bool
	prevWinnerCells []bool
	activeCells     []bool
	winnerCells     []bool

	// Metrics
}

func NewV2(p V2Params) Interface {
	// basal connections params, use local
	bpar := cells.V2Params{
		NumColumns:       p.NumColumns,
		CellsPerCol:      p.CellsPerCol,
		SegsPerCell:      p.SegsPerCell,
		SynsPerSeg:       p.SynsPerSeg,
		MatchThreshold:   p.MatchThreshold,
		ActiveThreshold:  p.ActiveThreshold,
		SynPermConnected: p.SynPermConnected,
	}
	// for now this works,
	// but really the apical input size can be something
	// totally different from cols*cells
	p.NumBasalCells = p.NumColumns * p.CellsPerCol
	p.NumApicalCells = p.NumColumns * p.CellsPerCol

	return &V2{
		P:      p,
		Basal:  cells.NewV2(bpar),
		Apical: cells.NewV2(bpar),
	}
}

// feedforward + feedback
func (v *V2) Compute(learn bool, cols, basal, apical []bool) error {
	switch {
	case len(cols) != v.P.NumColumns:
		return errors.WithStack(errors.New("column count mismatch"))
	case len(basal) != v.P.NumBasalCells:
		return errors.WithStack(errors.New("basal cell count mismatch"))
	case len(apical) != v.P.NumApicalCells:
		return errors.WithStack(errors.New("apical cell count mismatch"))
	case v.P.MatchThreshold >= v.P.ActiveThreshold:
		return errors.WithStack(errors.New("match >= active"))
	}

	// calculate predictions
	bActSegs, _ := v.Basal.ComputeActivity(basal)
	aActSegs, _ := v.Apical.ComputeActivity(apical)
	pCells := v.computePrediction(bActSegs, aActSegs) // export pCells
	fmt.Printf("%d cells\n", len(pCells))

	// learn
	if learn {
	}

	return nil
}

func (v *V2) computePrediction(b, a [][]int) []bool {
	// depolarize cells if active on basal && apical
	full := make([]bool, v.P.NumBasalCells) // basal+apical
	part := make([]bool, v.P.NumBasalCells) // basal
	for i := range full {
		switch {
		case len(b[i]) > 0 && len(a[i]) > 0:
			full[i] = true
		case len(b[i]) > 0:
			part[i] = true
		}
	}

	// apical connections inhibit other cells in the same col
	cols := make([]bool, v.P.NumColumns)
	for i := range full {
		col := i / v.P.CellsPerCol
		cols[col] = cols[col] || full[i]
	}

	// depolarize cells that survive inhibition
	for i := range part {
		col := i / v.P.CellsPerCol
		if part[i] && cols[col] == false {
			full[i] = true
		}
	}

	return full
}

func (v *V2) Reset() {}

func (v *V2) ActiveCells() []bool {
	return make([]bool, v.P.NumBasalCells)
}
func (v *V2) WinnerCells() []bool {
	return make([]bool, v.P.NumBasalCells)
}
