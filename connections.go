package gohtm

import "math"

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

func NewCell() *Cell {
	return &Cell{
		segments: []*Segment{},
	}
}

// Connections stores the connectivity of a TemporalMemory region.
type Connections struct {
	numCells   int
	segPerCell int
	synPerSeg  int

	numSegments    int
	numSynapses    int
	nextFlatIdx    int
	nextSegOrdinal int
	iteration      int
	cells          []*Cell
}

func NewConnections(cell, seg, syn int) *Connections {
	c := Connections{
		numCells:       cell,
		segPerCell:     seg,
		synPerSeg:      syn,
		numSegments:    0,
		numSynapses:    0,
		nextFlatIdx:    0,
		nextSegOrdinal: 0,
		iteration:      0,
	}

	c.cells = make([]*Cell, c.numCells)
	for i := 0; i < c.numCells; i++ {
		c.cells[i] = NewCell()
	}

	return &c
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

// LeastRecentSeg
func (c *Connections) LeastRecentSegment(cell int) *Segment {
	var min *Segment
	var minIter int = math.MaxInt64
	for _, seg := range c.SegmentsForCell(cell) {
		if seg.lastIter < minIter {
			min = seg
			minIter = seg.lastIter
		}
	}
	return min
}

// minPermSynapse
// segForFlatIdx
// lenFlatList
// synsForPresynapticCell

// CreateSegment creates a segment on the provided cell index.
func (c *Connections) CreateSegment(cellIdx int) *Segment {
	for len(c.SegmentsForCell(cellIdx)) >= c.segPerCell {
		c.DestroySegment(c.LeastRecentSegment(cellIdx))
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
//
// seg : *Segment on which to create synapse.
//
// pre : presynaptic cell from which to connect.
//
// perm: initial permanence value
func (c *Connections) CreateSynapse(seg *Segment, pre int, perm float64) {
}

// deleteSynapse
func (c *Connections) DestroySynapse() {
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
