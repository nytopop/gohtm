package cells

/*
// Synapse contains information about a unique synapse.
type Synapse struct {
	segment         *Segment
	preSynapticCell int // should be a *Cell
	ordinal         int
	perm            float64
	modPerm         float64 // temporary permanence modification
}

// NewSynapse creates a new synapse.
func NewSynapse(seg *Segment, pre int, perm float64, ord int) Synapse {
	return Synapse{
		segment:         seg,
		preSynapticCell: pre,
		ordinal:         ord,
		perm:            perm,
	}
}

// Segment contains information about a unique segment.
type Segment struct {
	cell     int
	flatIdx  int
	lastIter int
	ordinal  int
	synapses map[Synapse]bool
	state    bool // false inactive, true active
}

// NewSegment returns a new blank Segment with the supplied parameters.
func NewSegment(cell, flatIdx, lastIter, ordinal int) *Segment {
	return &Segment{
		cell:     cell,
		flatIdx:  flatIdx,
		lastIter: lastIter,
		ordinal:  ordinal,
		synapses: map[Synapse]bool{},
	}
}

// Cell contains a number of connection points to other Cells.
type Cell struct {
	segments []*Segment
	state    int // 0 inactive; 1 predicted; 2 active
}

func NewCell() *Cell {
	return &Cell{
		segments: []*Segment{},
		state:    0,
	}
}

// Connections stores the connectivity of a TemporalMemory region.
type Connections struct {
	numColumns int
	numCells   int
	segPerCell int
	synPerSeg  int

	numSegments    int
	numSynapses    int
	nextFlatIdx    int
	nextSegOrdinal int
	nextSynOrdinal int
	iteration      int
	cells          []*Cell
}

// NewConnections returns a new Connections.
func NewConnections(numColumns, numCells, segPerCell, synPerSeg int) *Connections {
	c := Connections{
		numColumns:     numColumns,
		numCells:       numCells,
		segPerCell:     segPerCell,
		synPerSeg:      synPerSeg,
		numSegments:    0,
		numSynapses:    0,
		nextFlatIdx:    0,
		nextSegOrdinal: 0,
		nextSynOrdinal: 0,
		iteration:      0,
		cells:          make([]*Cell, numColumns*numCells),
	}

	for i := range c.cells {
		c.cells[i] = NewCell()
	}

	return &c
}

// PredictedForCol returns any depolarized cells in a column.
func (c *Connections) PredictedForCol(col int) []*Cell {
	pre := []*Cell{}

	for _, cell := range c.CellsForCol(col) {
		if cell.state == 1 {
			pre = append(pre, cell)
		}
	}

	return pre
}

// CellsForCol returns []*Cell for the provided column index.
func (c *Connections) CellsForCol(col int) []*Cell {
	cells := make([]*Cell, c.numCells)

	var idx, min, max int
	min = col * c.numCells
	max = min + c.numCells
	for i := min; i < max; i++ {
		cells[idx] = c.cells[i]
		idx++
	}

	return cells
}

// SegmentsForCell returns segments that belong to the provided cell idx.
func (c *Connections) SegmentsForCell(cell int) []*Segment {
	return c.cells[cell].segments
}

// SynapsesForSegment returns a map[Synapse]bool for the
// provided *Segment.
func (c *Connections) SynapsesForSegment(seg *Segment) map[Synapse]bool {
	return seg.synapses
}

// GetSegment returns the Segment identified by a
// 2-dimensional address (cell, idx)
func (c *Connections) GetSegment(cell, idx int) *Segment {
	return c.cells[cell].segments[idx]
}

// leastRecentSeg returns the *Segment that was updated
// the longest time ago.
func (c *Connections) leastRecentSegment(cell int) *Segment {
	var min *Segment
	var minIter = math.MaxInt64
	for _, seg := range c.SegmentsForCell(cell) {
		if seg.lastIter < minIter {
			min = seg
			minIter = seg.lastIter
		}
	}
	return min
}

// minPermSynapse returns the Synapse with smallest permanence
// on provided *Segment.
func (c *Connections) minPermSynapse(seg *Segment) Synapse {
	var min Synapse
	var minPerm = 500.0
	for syn := range seg.synapses {
		if syn.perm < minPerm {
			min = syn
			minPerm = syn.perm
		}
	}
	return min
}

// segForFlatIdx
// lenFlatList
// synsForPresynapticCell

// CreateSegment creates a segment on the provided cell index.
func (c *Connections) CreateSegment(cellIdx int) *Segment {
	for len(c.SegmentsForCell(cellIdx)) >= c.segPerCell {
		c.DestroySegment(c.leastRecentSegment(cellIdx))
	}

	flatIdx := c.nextFlatIdx
	c.nextFlatIdx++

	ordinal := c.nextSegOrdinal
	c.nextSegOrdinal++

	seg := NewSegment(cellIdx, flatIdx, c.iteration, ordinal)
	c.cells[cellIdx].segments = append(c.cells[cellIdx].segments, seg)

	c.updateCounts()

	return seg
}

// DestroySegment destroys the provided *Segment, removing it from
// all datastructures.
func (c *Connections) DestroySegment(seg *Segment) {
	// TODO remove from all main datastructures

	// remove from the cells segment list
	for i, v := range c.cells[seg.cell].segments {
		if v == seg {
			// verbosity up the arse
			copy(c.cells[seg.cell].segments[i:], c.cells[seg.cell].segments[i+1:])
			c.cells[seg.cell].segments[len(c.cells[seg.cell].segments)-1] = nil
			c.cells[seg.cell].segments = c.cells[seg.cell].segments[:len(c.cells[seg.cell].segments)-1]
			c.updateCounts()
		}
	}
}

// CreateSynapse creates a new synapse on a segment.
// Params: seg is the segment on which to create the synapse, pre
// is the presynaptic cell to connect to, perm is permanence.
func (c *Connections) CreateSynapse(seg *Segment, pre int, perm float64) {
	for len(c.SynapsesForSegment(seg)) >= c.synPerSeg {
		c.DestroySynapse(c.minPermSynapse(seg))
	}

	syn := NewSynapse(seg, pre, perm, c.nextSynOrdinal)
	c.nextSynOrdinal++

	seg.synapses[syn] = true

	c.updateCounts()
}

// DestroySynapse removes the provided synapse from a segment.
func (c *Connections) DestroySynapse(syn Synapse) {
	delete(syn.segment.synapses, syn)
	c.updateCounts()
}

// updateCounts updates the numSegments and numSynapses counters
func (c *Connections) updateCounts() {
	var sumSegs int
	var sumSyns int
	for _, cell := range c.cells {
		for _, seg := range cell.segments {
			sumSegs++
			for _ = range seg.synapses {
				sumSyns++
			}
		}
	}
	c.numSegments = sumSegs
	c.numSynapses = sumSyns
}

// updateSynPerm updates the permanence value of a provided synapse
func (c *Connections) updateSynPerm(syn Synapse, perm float64) {
	new := syn
	new.perm = perm

	delete(syn.segment.synapses, syn)
	new.segment.synapses[new] = true
}

func (c *Connections) tempUpdateSynPerm(syn Synapse, mod float64) {
	new := syn
	new.modPerm = syn.modPerm + mod

	delete(syn.segment.synapses, syn)
	new.segment.synapses[new] = true
}
*/
