package gohtm

/* Connections for Temporal Memory */

// Synapse contains information about a unique synapse.
type Synapse struct {
	segment         *Segment
	preSynapticCell int
	ordinal         int
	permanence      float64
}

// Segment contains information about a unique segment.
type Segment struct {
	cell     int
	flatIdx  int
	lastIter int
	ordinal  int
	synapses map[Synapse]bool
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
}

func NewCell(seg int) *Cell {
	return &Cell{
		segments: []*Segment{},
	}
}

// Connections stores the connectivity of a TemporalMemory region.
type Connections struct {
	numCells   int
	segPerCell int
	synPerSeg  int

	numSynapses int
	iteration   int
	cells       []*Cell
}

func NewConnections(cell, seg, syn int) *Connections {
	c := Connections{
		numCells:    cell,
		segPerCell:  seg,
		synPerSeg:   syn,
		numSynapses: 0,
		iteration:   0,
	}

	c.cells = make([]*Cell, c.numCells)
	for i := 0; i < c.numCells; i++ {
		c.cells[i] = NewCell(c.segPerCell)
	}

	return &c
}

// SegmentsForCell returns segments that belong to the provided cell idx.
func (c *Connections) SegmentsForCell(cell int) []*Segment {
	return c.cells[cell].segments
}

// synsForSeg
func (c *Connections) SynapsesForSegment(seg *Segment) map[Synapse]bool {
	return seg.synapses
}

// getSeg
func (c *Connections) GetSegment(cell, idx int) *Segment {
	return c.cells[cell].segments[idx]
}

// leastRecentSeg
// minPermSynapse
// segForFlatIdx
// lenFlatList
// synsForPresynapticCell
// createSegment
// deleteSegment
// createSynapse
// deleteSynapse
