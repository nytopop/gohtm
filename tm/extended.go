package tm

import "github.com/nytopop/gohtm/cells"

// ExtendedParams contains parameters for initialization of an Extended
// TemporalMemory region.
type ExtendedParams struct {
	NumColumns          int // input space dimensions
	NumCells            int // cells per column
	SegPerCell          int
	SynPerSeg           int
	ActivationThreshold int     // # of active synapses for 'active'
	MinThreshold        int     // # of potential synapses for learning
	InitPermanence      float64 // initial permanence of new synapses
	SynPermConnected    float64 // connection threshold
	SynPermInactiveMod  float64 // perm decrement
	SynPermActiveMod    float64 // perm increment
}

// NewExtendedParams returns a default ExtendedParams.
func NewExtendedParams() ExtendedParams {
	return ExtendedParams{
		NumColumns:          2048,
		CellsPerColumn:      32,
		SegPerCell:          16,
		SynPerSeg:           16,
		ActivationThreshold: 12,
		MinThreshold:        10,
		InitPermanence:      0.21,
		SynPermConnected:    0.5,
		SynPermActiveMod:    0.05,
		SynPermInactiveMod:  0.03,
	}
}

// Extended is a TemporalMemory implementation with support
// for basal and apical dendrites connected to other regions.
type Extended struct {
	P ExtendedParams
	//Cons *cells.Connections
	Cons *cells.Cells
}

// NewExtended initializes a new TemporalMemory region with
// the provided ExtendedParams.
func NewExtended(p ExtendedParams) *Extended {
	cp := cells.NewV1Params()
	cp.Columns = p.NumColumns
	cp.CellsPerColumn = p.CellsPerColumn
	cp.SegmentsPerCell = p.SegPerCell
	cp.SynapsesPerSegment = p.SynPerSeg

	return &Extended{
		P: p,
		/*Cons: cells.NewConnections(
			p.NumColumns,
			p.NumCells,
			p.SegPerCell,
			p.SynPerSeg,
		),*/
	}
}

// Compute iterates the TemporalMemory algorithm with the provided
// vector of active columns from a SpatialPooler.
func (e *Extended) Compute(active []bool, learn bool) []bool {
	return []bool{}
}

// Reset clears temporary data so sequences are not learned between
// the current and next time step.
func (e *Extended) Reset() {
}

func (e *Extended) Save() {
}

func (e *Extended) Load() {
}
