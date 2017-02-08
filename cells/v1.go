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
	NumColumns         int     // number of columns
	CellsPerCol        int     // cells per column
	SegsPerCell        int     // max segments per cell
	SynsPerSeg         int     // max synapses per segment
	SynPermConnected   float32 // synapse permanence connection threshold
	SynPermActiveMod   float32 // synapse permanence increment
	SynPermInactiveMod float32 // synapse permanence decrement
	InitPerm           float32 // initial permanence of new synapses
	ActiveThreshold    int     // number live synapses to mark segment active
	MatchingThreshold  int     // # of potential synapses for learning ????
}

// NewV1Params returns a default parameter set.
func NewV1Params() V1Params {
	return V1Params{
		NumColumns:         64,
		CellsPerCol:        8,
		SegsPerCell:        16,
		SynsPerSeg:         16,
		SynPermConnected:   0.5,
		SynPermActiveMod:   0.05,
		SynPermInactiveMod: 0.03,
		InitPerm:           0.21,
		ActiveThreshold:    12,
		MatchingThreshold:  8,
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
	v := &V1{
		P:         p,
		iteration: 0,
	}

	// allocate cells
	v.cells = make([]V1Cell, p.NumColumns*p.CellsPerCol)
	return v
}

// V1Cell represents a single cell.
type V1Cell struct {
	active, matching int // number of active & matching segments
	segments         []V1Segment
}

// V1Segment represents a single dendritic segment attached to a cell.
type V1Segment struct {
	active, matching bool
	synapses         []V1Synapse
}

// V1Synapse represents the connection from the dendritic segment of a cell
// to another cell, called the presynaptic cell.
type V1Synapse struct {
	cell int     // presynaptic cell
	perm float32 // permanence of connection
}

func (v *V1) CreateSegment(cell int) {
}

func (v *V1) DestroySegment(cell, seg int) {
}

func (v *V1) CreateSynapse(cell, seg, target int) {
}

func (v *V1) DestroySynapse(cell, seg, syn int) {
}

// AdaptSynapses adapts synapses on a cell to the provided slice
// of cells that were active in the previous time step.
func (v *V1) AdaptSynapses(cell int, prevActive []bool) {
	// delete queue for synapses
	synQueue := make([][2]int, 0)

	var perm float32
	for i := range v.cells[cell].segments {
		for j := range v.cells[cell].segments[i].synapses {
			perm = v.cells[cell].segments[i].synapses[j].perm

			// check if synapse.cell is in prevActiveCells
			switch prevActive[v.cells[cell].segments[i].synapses[j].cell] {
			case true:
				perm += v.P.SynPermActiveMod
			case false:
				perm -= v.P.SynPermInactiveMod
			}

			// constrain permanence to [0.0 : 1.0]
			switch {
			case perm < 0.0:
				perm = 0.0
			case perm > 1.0:
				perm = 1.0
			}

			// If permanence < 0.001, add synapse to destruction queue
			// If not, update permanence
			switch {
			case perm < 0.001:
				var syn [2]int
				syn[0], syn[1] = i, j
				synQueue = append(synQueue, syn)
			default:
				v.cells[cell].segments[i].synapses[j].perm = perm
			}
		}
	}

	// Destroy synapses in queue
	for i := range synQueue {
		v.DestroySynapse(cell, synQueue[i][0], synQueue[i][1])
	}

	// Add empty segments to destruction queue
	segQueue := make([]int, 0)
	for i := range v.cells[cell].segments {
		if len(v.cells[cell].segments[i].synapses) == 0 {
			segQueue = append(segQueue, i)
		}
	}

	// Destroy segments in queue
	for i := range segQueue {
		v.DestroySegment(cell, segQueue[i])
	}
}

// GrowSynapses grows new synapses on a cell to the provided slice
// of cells chosen as winners in the previous time step.
func (v *V1) GrowSynapses(cell int, prevWinners []bool) {
	// for each active seg on cell
	// grow new synapses to prev winner cells
}

// CellsForCol returns a slice of all cell indices within the
// provided column index.
func (v *V1) CellsForCol(col int) []int {
	out := make([]int, 0, v.P.CellsPerCol)
	for i := col * v.P.CellsPerCol; i < (col+1)*v.P.CellsPerCol; i++ {
		out = append(out, i)
	}
	return out
}

// ActiveSegsForCell returns the number of active segments
// attached to a cell.
func (v *V1) ActiveSegsForCell(cell int) int {
	return v.cells[cell].active
}

// ActiveSegsForCol returns the number of active segments
// attached to cells in a column.
func (v *V1) ActiveSegsForCol(col int) int {
	cells := v.CellsForCol(col)
	var syns int
	for _, i := range cells {
		syns += v.cells[i].active
	}
	return syns
}

// MatchingSegsForCell returns the number of matching segments
// attached to a cell.
func (v *V1) MatchingSegsForCell(cell int) int {
	return v.cells[cell].matching
}

// MatchingSegsForCol returns the number of matching segments
// attached to cells in a column.
func (v *V1) MatchingSegsForCol(col int) int {
	cells := v.CellsForCol(col)
	var syns int
	for _, i := range cells {
		syns += v.cells[i].matching
	}
	return syns
}

// ComputeActivity computes cell, segment, and synapse activity
// in regards to currently active columns.
func (v *V1) ComputeActivity(active []bool) {
	for i := range v.cells {
		for j := range v.cells[i].segments {
			// count live synapses on each segment
			var live int
			for _, syn := range v.cells[i].segments[j].synapses {
				// if the synapse corresponds to a cell in an active column
				if active[syn.cell] {
					// if synapse is connected
					if syn.perm >= v.P.SynPermConnected {
						live += 1
					}
				}
			}
			// if over threshold, activate dendrite segment
			switch {
			case live >= v.P.ActiveThreshold:
				v.cells[i].segments[j].active = true
				v.cells[i].active += 1
			case live >= v.P.MatchingThreshold:
				v.cells[i].segments[j].matching = true
				v.cells[i].matching += 1
			}
		}

	}
}

// Clear clears temporary data from all cells and segments.
func (v *V1) Clear() {
	for i := range v.cells {
		v.cells[i].active = 0
		v.cells[i].matching = 0
		for j := range v.cells[i].segments {
			v.cells[i].segments[j].active = false
			v.cells[i].segments[j].matching = false
		}
	}
}

// StartNewIteration increments the iteration counter.
func (v *V1) StartNewIteration() {
	v.iteration += 1
}
