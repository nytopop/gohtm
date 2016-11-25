package gohtm

import "fmt"

/* Temporal Memory */

type TemporalParams struct {
	numColumns          int
	numCells            int
	activationThreshold int
	minThreshold        int
	initPermanence      float64
	synPermConnected    float64
	segPerCell          int
	synPerSeg           int
}

func NewTemporalParams() TemporalParams {
	return TemporalParams{
		numColumns:          128,  // input space dimensions
		numCells:            4,    // cells per column
		activationThreshold: 12,   // # of active synapses for 'active'
		minThreshold:        10,   // # of potential synapses for learning
		initPermanence:      0.21, // initial permanence of new synapses
		synPermConnected:    0.5,  // connection threshold
		segPerCell:          64,
		synPerSeg:           64,
	}
}

/*// Each synapse contains information about the link
type DistalSynapse struct {
	seg       *list.Element
	connected bool
	perm      float64
}

// Each segment has a variable linked list of synapses
type Segment struct {
	dsyns list.List // { DistalSynapse }
}

// Each cell has a variable linked list of Segments
type Cell struct {
	state int // 0,1,2,3 : inactive, predicted, active, bursting
	//seg   Segment
	seg list.List // { Segment }
}

// Each column is a static sized slice of Cells
type TMColumn []Cell

//type TMColumn list.List

// TODO : synapses should be separated from cols/cells/segments
// synapse should be two ptrs
*/

type TemporalMemory struct {
	// state
	connections Connections

	// params
	numColumns          int
	numCells            int
	activationThreshold int
	minThreshold        int
	initPermanence      float64
	synPermConnected    float64
}

func NewTemporalMemory(p TemporalParams) TemporalMemory {
	tm := TemporalMemory{
		numColumns:          p.numColumns,
		numCells:            p.numCells,
		activationThreshold: p.activationThreshold,
		minThreshold:        p.minThreshold,
		initPermanence:      p.initPermanence,
		synPermConnected:    p.synPermConnected,
	}

	totalCells := p.numColumns * p.numCells
	fmt.Println(totalCells, "total cells")
	tm.connections = NewConnections(totalCells, p.segPerCell, p.synPerSeg)
	fmt.Println(tm.connections)

	/*
		tm.cols = make([]TMColumn, p.numColumns)
		for i, _ := range tm.cols {
			// i is column index
			tm.cols[i] = make(TMColumn, p.numCells)

			for j, _ := range tm.cols[i] {
				// j is cell index

				tm.cols[i][j].seg = Segment{
					dsyns: make([]DistalSynapse, 128),
				}
			}
		}
	*/

	return tm
}

func (tm *TemporalMemory) Compute(active SparseBinaryVector) SparseBinaryVector {
	return SparseBinaryVector{}
}
