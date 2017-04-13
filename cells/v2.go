package cells

import (
	"fmt"
	"math"
	"math/rand"
)

// V2Params ...
type V2Params struct {
	NumColumns  int `json:"numcolumns"`
	CellsPerCol int `json:"cellspercol"`
	SegsPerCell int `json:"segspercell"`
	SynsPerSeg  int `json:"synsperseg"`
}

// V2 ...
type V2 struct {
	P         V2Params `json:"params"`
	Cells     []V2Cell `json:"cells"`
	iteration int
}

// NewV2 ...
func NewV2(p V2Params) Interface {
	return &V2{
		P:     p,
		Cells: make([]V2Cell, p.NumColumns*p.CellsPerCol),
	}
}

// V2Cell ...
type V2Cell struct {
	Segments []V2Segment `json:"segments"`
}

// V2Segment ...
type V2Segment struct {
	Synapses []V2Synapse `json:"synapses"`
	lastIter int
}

// V2Synapse
type V2Synapse struct {
	Idx  int     `json:"idx"`
	Perm float32 `json:"perm"`
}

// CreateSegment creates a segment on the specified cell and
// returns the index. The segment will be synapsed onto up to
// c.P.SynsPerSeg randomly sampled cells from those provided.
func (c *V2) CreateSegment(cell int, targets []bool, perm float32) int {
	// allocate
	c.Cells[cell].Segments = append(
		c.Cells[cell].Segments,
		V2Segment{
			Synapses: make([]V2Synapse, c.P.SynsPerSeg),
			lastIter: c.iteration,
		})

	// TODO bounds check
	var mi, mv int
	mi, mv = math.MaxInt64, math.MaxInt64
	if len(c.Cells[cell].Segments) > c.P.SegsPerCell {
		for i := range c.Cells[cell].Segments {
			if c.Cells[cell].Segments[i].lastIter < mv {
				mv, mi = c.Cells[cell].Segments[i].lastIter, i
			}
		}

		// got it, now we remove mi
		fmt.Printf("lowest iter %d\n", mi)
	}

	// gen sample
	var cells, sample []int
	for i := range targets {
		if targets[i] {
			cells = append(cells, i)
		}
	}

	switch {
	case len(cells) >= c.P.SynsPerSeg:
		sample = rand.Perm(len(cells))[:c.P.SynsPerSeg]
	default:
		sample = rand.Perm(len(cells))
	}

	// synapse
	idx := len(c.Cells[cell].Segments) - 1
	for i := range sample {
		c.Cells[cell].Segments[idx].Synapses = append(
			c.Cells[cell].Segments[idx].Synapses,
			V2Synapse{
				Idx:  sample[i],
				Perm: perm,
			})
	}

	return idx
}
