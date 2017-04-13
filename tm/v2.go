package tm

import (
	"fmt"

	"github.com/nytopop/gohtm/cells"
	"github.com/pkg/errors"
)

// V2 ... extended TM. basal, apical
type V2Params struct {
	Apical bool `json:"apical"`

	NumColumns      int `json:"numcolumns"`
	CellsPerCol     int `json:"cellspercol"`
	SegsPerCell     int `json:"segspercell"`
	SynsPerSeg      int `json:"synsperseg"`
	MatchThreshold  int `json:"matchthreshold"`
	ActiveThreshold int `json:"activethreshold"`

	NumBasalCells  int `json:"numbasalcells"`
	NumApicalCells int `json:"numapicalcells"`
}

func NewV2Params() V2Params {
	return V2Params{
		Apical:          false,
		NumColumns:      2048,
		CellsPerCol:     16,
		SegsPerCell:     32,
		SynsPerSeg:      16,
		MatchThreshold:  6,
		ActiveThreshold: 12,
		NumBasalCells:   65535,
		NumApicalCells:  65535,
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
		NumColumns:  p.NumColumns,
		CellsPerCol: p.CellsPerCol,
		SegsPerCell: p.SegsPerCell,
		SynsPerSeg:  p.SynsPerSeg,
	}

	apar := cells.V2Params{
		NumColumns:  p.NumColumns,
		CellsPerCol: p.CellsPerCol / 2,
		SegsPerCell: p.SegsPerCell / 2,
		SynsPerSeg:  p.SynsPerSeg,
	}

	return &V2{
		P:      p,
		Basal:  cells.NewV2(bpar),
		Apical: cells.NewV2(apar),
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

	fmt.Printf("cols %d\nbasal %d\napical %d\n",
		len(cols), len(basal), len(apical))

	return nil
}

func (v *V2) Reset() {}

func (v *V2) ActiveCells() []bool {
	return make([]bool, v.P.NumBasalCells)
}
func (v *V2) WinnerCells() []bool {
	return make([]bool, v.P.NumBasalCells)
}
