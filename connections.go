package gohtm

/* Connections for Temporal Memory */

// Synapse connects a Segment to another Segment.
type Synapse struct {
	segment *Segment
}

// Segment contains information about a cell's connectivity.
type Segment struct {
	cell     int
	flatIdx  int
	lastIter int
	ordinal  int
	synapses []Synapse
}

// NewSegment returns a new blank Segment with the supplied parameters.
func NewSegment(cell, flatIdx, lastIter, ordinal int) Segment {
	return Segment{
		cell:     cell,
		flatIdx:  flatIdx,
		lastIter: lastIter,
		ordinal:  ordinal,
		synapses: []Synapse{},
	}
}

// Cell contains a number of connection points to other Cells.
type Cell struct {
	segments []Segment
}

func NewCell(seg int) Cell {
	return Cell{
		segments: []Segment{},
	}
}

// Connections stores the connectivity of a TemporalMemory region.
type Connections struct {
	numCells   int
	segPerCell int
	synPerSeg  int

	cells []Cell
}

func NewConnections(cell, seg, syn int) Connections {
	c := Connections{
		numCells:   cell,
		segPerCell: seg,
		synPerSeg:  syn,
	}

	c.cells = make([]Cell, c.numCells)
	for i := 0; i < c.numCells; i++ {
		c.cells[i] = NewCell(c.segPerCell)
	}

	return c
}
