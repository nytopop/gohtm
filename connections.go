package gohtm

/* Connections for Temporal Memory */

/* Synapse */
type Synapse struct {
	segment *Segment
}

/* Segment */
type Segment struct {
	cell     int
	flatIdx  int
	synapses []Synapse
}

func NewSegment() Segment {
	return Segment{}
}

/* Cell */
type Cell struct {
	segments []Segment
}

func NewCell(seg int) Cell {
	return Cell{
		segments: []Segment{},
	}
}

/* Connections */
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
