package cells

/*
Synapses connect a segment --> cell

cell1: A.1.2.3.4
cell2: B.1.2.3.4

syn1: A.1 --> B
syn2: B.1 --> A
NOT DUPLICATES!

Do everything at the index level
cells -> columns -> segments -> synapses
*/

// V1Params contains parameters for initialization of V1 cellular state.
type V1Params struct {
	Columns            int
	CellsPerColumn     int
	SegmentsPerCell    int
	SynapsesPerSegment int
}

// NewV1Params returns a default parameter set.
func NewV1Params() V1Params {
	return V1Params{
		Columns:            64,
		CellsPerColumn:     8,
		SegmentsPerCell:    16,
		SynapsesPerSegment: 16,
	}
}

// V1 cellular state interface.
type V1 struct {
	P         V1Params
	cells     []V1Cell
	iteration int
}

// NewV1 returns a new V1 instance, initialized according to provided params.
func NewV1(p V1Params) *V1 {
	v := V1{
		P:         p,
		iteration: 0,
	}

	// allocate cells
	v.cells = make([]V1Cell, p.Columns*p.CellsPerColumn)
	return &v
}

// V1Cell represents a neuron.
type V1Cell struct {
	state    int
	segments []V1Segment
}

// V1Segment represents a single dendritic segment attached to a cell.
type V1Segment struct {
	active   bool
	live     int
	synapses []V1Synapse
}

// V1Synapse represents the connection from the dendritic segment of a cell
// to another cell, called the presynaptic cell.
type V1Synapse struct {
	cell      int // presynaptic cell
	perm      int // permanence is integer, faster
	active    bool
	connected bool
}

func (v *V1) CreateSegment(cell int) {
}

func (v *V1) DestroySegment(cell, seg int) {
}

func (v *V1) CreateSynapse(cell, seg, target int) {
}

func (v *V1) DestroySynapse(cell, seg, syn int) {
}

func (v *V1) DepolarizedForCol(col int) []int {
	return []int{}
}

func (v *V1) ComputeStatistics(active []bool) {
	v.Clear()
	for i := range v.cells {
		for j := range v.cells[i].segments {
			for k := range v.cells[i].segments[j].synapses {
				if v.cells[i].segments[j].synapses[k].connected &&
					active[v.cells[i].segments[j].synapses[k].cell] {
					// increment counter, activate synapse
				}
				// i, j, k !
				// check if this synapses
			}
		}
	}
	/*
		for each segment, count # connected synapses to active cells : live

	*/
}

func (v *V1) Clear() {
	for i := range v.cells {
		for j := range v.cells[i].segments {
			for k := range v.cells[i].segments[j].synapses {
				v.cells[i].segments[j].synapses[k].active = false
				v.cells[i].segments[j].synapses[k].connected = false
			}
		}
	}
}
