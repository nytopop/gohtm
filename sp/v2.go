package sp

import (
	"log"
	"math/rand"
	"sort"
)

// V2Params contains parameters for initialization of a V2 SpatialPooler.
type V2Params struct {
	NumColumns       int
	NumInputs        int
	PotentialRadius  int
	PotentialPct     float64 //
	InitConnPct      float64 //
	SynPermConnected float32 //
	SynPermMod       float32 //
	Sparsity         float64
	DutyCyclePeriod  int
	MinDutyCycle     float64 //
	MaxBoost         float64 //
}

// NewV2Params returns a default set of V2Params.
func NewV2Params() V2Params {
	return V2Params{
		NumColumns:       64,
		NumInputs:        64,
		PotentialRadius:  0,
		PotentialPct:     0.10,
		InitConnPct:      0.3,
		SynPermConnected: 0.5,
		SynPermMod:       0.05,
		Sparsity:         0.02,
		DutyCyclePeriod:  32,
		MinDutyCycle:     0.2,
		MaxBoost:         8.0,
	}
}

// V2 SpatialPooler.
type V2 struct {
	P         V2Params
	cells     []V2Cell
	iteration int
}

// NewV2 initializes and returns a new V2 SpatialPooler with the
// provided parameters.
func NewV2(p V2Params) *V2 {
	s := &V2{
		P:         p,
		cells:     make([]V2Cell, p.NumColumns),
		iteration: 0,
	}

	// Initialize potential synapses
	for i := range s.cells {
		s.mapPotential(i)
		s.cells[i].oPeriod = make([]int, 0)
		s.cells[i].aPeriod = make([]bool, 0)
		s.cells[i].boostFactor = 1.0
	}

	return s
}

// V2Cell ...
type V2Cell struct {
	boostFactor      float64 // boost multiplier
	oPeriod          []int   // overlap period
	overlapDutyCycle float64 // overlap duty cycle
	aPeriod          []bool  // active period
	activeDutyCycle  float64 // active duty cycle
	synapses         []V2Synapse
}

// V2Synapse ...
type V2Synapse struct {
	idx  int     // input index
	perm float32 // permanence
}

// Compute ...
func (s *V2) Compute(input []bool, learn bool) []bool {
	if len(input) != s.P.NumInputs {
		log.Fatalln("SP: Mismatched input dims", s.P.NumInputs, len(input))
	}

	// Calculate overlaps and inhibit cells
	overlaps := s.calcOverlaps(input)
	activeCells := s.inhibitCells(overlaps, learn)

	// Perform learning
	if learn {
		s.adaptSynapses(input, activeCells)
		s.updateOverlapDutyCycles(overlaps)
		s.updateActiveDutyCycles(activeCells)
		s.bumpWeakCells()
		s.updateBoostFactors()
		s.iteration++
	}

	return activeCells
}

// calcOverlaps ...
func (s *V2) calcOverlaps(input []bool) []int {
	overlaps := make([]int, s.P.NumColumns)
	for i := range s.cells {
		for j := range s.cells[i].synapses {
			if s.cells[i].synapses[j].perm >= s.P.SynPermConnected {
				if input[s.cells[i].synapses[j].idx] {
					overlaps[i]++
				}
			}
		}
	}
	return overlaps
}

// V2InhCell ...
type V2InhCell struct {
	idx     int
	overlap float64
}

// V2InhNet ...
type V2InhNet []V2InhCell

func (v V2InhNet) Len() int           { return len(v) }
func (v V2InhNet) Less(i, j int) bool { return v[i].overlap > v[j].overlap }
func (v V2InhNet) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }

// inhibitCells ...
func (s *V2) inhibitCells(overlapScores []int, learn bool) []bool {
	// create sortable slice of cells w/ overlap scores
	overlaps := make(V2InhNet, s.P.NumColumns)

	// apply boosting if learning is enabled
	switch learn {
	case true:
		for i := range overlaps {
			overlaps[i] = V2InhCell{
				idx:     i,
				overlap: float64(overlapScores[i]) * s.cells[i].boostFactor,
			}
		}
	case false:
		for i := range overlaps {
			overlaps[i] = V2InhCell{
				idx:     i,
				overlap: float64(overlapScores[i]),
			}
		}
	}

	// sort cells by overlap score, descending order
	sort.Sort(overlaps)

	// calc number of active columns
	n := int(s.P.Sparsity * float64(s.P.NumColumns))

	// return surviving cells
	winners := overlaps[:n]
	activeCells := make([]bool, s.P.NumColumns)
	for i := range winners {
		activeCells[winners[i].idx] = true
	}
	return activeCells
}

// adaptSynapses ...
func (s *V2) adaptSynapses(input []bool, activeCells []bool) {
	var perm float32
	for i := range activeCells {
		if activeCells[i] {
			for j := range s.cells[i].synapses {
				// decide whether to bump up or down
				switch input[s.cells[i].synapses[j].idx] {
				case true:
					// bump up contributing synapse
					perm = s.cells[i].synapses[j].perm
					perm += s.P.SynPermMod

					// clamp [0.0 : 1.0]
					switch {
					case perm > 1.0:
						perm = 1.0
					case perm < 0.0:
						perm = 0.0
					}

					s.cells[i].synapses[j].perm = perm
				case false:
					// bump down non-contributing synapse
					perm = s.cells[i].synapses[j].perm
					perm -= s.P.SynPermMod

					// clamp [0.0 : 1.0]
					switch {
					case perm > 1.0:
						perm = 1.0
					case perm < 0.0:
						perm = 0.0
					}

					s.cells[i].synapses[j].perm = perm
				}
			}
		}
	}
}

// updateOverlapDutyCycles
func (s *V2) updateOverlapDutyCycles(overlaps []int) {
	// OVERLAP duty cycle is a moving average of the number of
	// inputs which overlapped with the each column
	for i := range s.cells {
		s.cells[i].oPeriod = append(s.cells[i].oPeriod, overlaps[i])
		if len(s.cells[i].oPeriod) > s.P.DutyCyclePeriod {
			s.cells[i].oPeriod = s.cells[i].oPeriod[1:]
		}

		var sum int
		for j := range s.cells[i].oPeriod {
			sum += s.cells[i].oPeriod[j]
		}

		// calc moving average
		s.cells[i].overlapDutyCycle =
			float64(sum) / float64(len(s.cells[i].oPeriod))
	}
}

// updateActiveDutyCycles
func (s *V2) updateActiveDutyCycles(activeCells []bool) {
	// ACTIVITY duty cycles is a moving average of
	// the frequency of activation for each column.
	for i := range s.cells {
		s.cells[i].aPeriod = append(s.cells[i].aPeriod, activeCells[i])
		if len(s.cells[i].aPeriod) > s.P.DutyCyclePeriod {
			s.cells[i].aPeriod = s.cells[i].aPeriod[1:]
		}

		var sum int
		for j := range s.cells[i].aPeriod {
			if s.cells[i].aPeriod[j] {
				sum++
			}
		}

		// calc moving average
		s.cells[i].activeDutyCycle =
			float64(sum) / float64(len(s.cells[i].aPeriod))
	}
}

// bumpWeakCells
func (s *V2) bumpWeakCells() {
	// increase permanence on all synapses
	// belonging to weak cells
	var perm float32
	for i := range s.cells {
		if s.cells[i].overlapDutyCycle < s.P.MinDutyCycle {
			for j := range s.cells[i].synapses {
				perm = s.cells[i].synapses[j].perm
				perm += s.P.SynPermMod
				switch {
				case perm > 1.0:
					perm = 1.0
				case perm < 0.0:
					perm = 0.0
				}
				s.cells[i].synapses[j].perm = perm
			}
		}
	}
}

// updateBoostFactors
func (s *V2) updateBoostFactors() {
	// boost factors are caomputed as an inverse linear relationship
	// between the active duty cycle and maximum boost setting
	// only cells with an active duty cycle below s.P.Sparsity are
	// boosted, all others remain at 1.0
	var r float64
	for i := range s.cells {
		switch {
		case s.cells[i].activeDutyCycle < s.P.Sparsity:
			r = 1 - (s.cells[i].activeDutyCycle / s.P.Sparsity)
			s.cells[i].boostFactor = r*s.P.MaxBoost + 1
		case s.cells[i].activeDutyCycle > s.P.Sparsity:
			s.cells[i].boostFactor = 1.0
		}
	}
}

// mapPotential creates potential synapses on the specified cell. This will
// grow synapses to a random sample of its receptive field.
func (s *V2) mapPotential(cell int) {
	// Find centerpoint in input space for this cell
	ratio := float64(cell) / float64(s.P.NumColumns)
	center := int(float64(s.P.NumInputs) * ratio)

	// Take random sample of inputs in receptive field of cell
	nbs := s.getInputNeighbors(center)
	n := int(float64(len(nbs)) * s.P.PotentialPct)
	sample := rand.Perm(len(nbs))[:n]

	// Grow synapses
	s.cells[cell].synapses = make([]V2Synapse, len(sample))
	for syn, idx := range sample {
		s.cells[cell].synapses[syn].idx = nbs[idx]
		s.cells[cell].synapses[syn].perm = s.getInitPerm()
	}
}

// getInputNeighbors returns all input indices within PotentialRadius of
// the provided input index. If PotentialRadius is <= 0, then all input
// indices are returned.
func (s *V2) getInputNeighbors(input int) []int {
	var nbs []int
	switch {
	case s.P.PotentialRadius <= 0:
		nbs = make([]int, s.P.NumInputs)
		for i := range nbs {
			nbs[i] = i
		}
	default:
		r := s.P.PotentialRadius
		nbs = make([]int, 0, r*2+1)
		for i := input - r; i <= input+r; i++ {
			// switch to test for overflow / underflow
			switch {
			case i >= 0 && i < s.P.NumInputs:
				// we are within normal range
				nbs = append(nbs, i)
			case i < 0:
				// we are underflowing
				nbs = append(nbs, i+s.P.NumInputs)
			case i >= s.P.NumInputs:
				// we are overflowing
				nbs = append(nbs, i-s.P.NumInputs)
			}
		}
	}
	return nbs
}

// getInitPerm returns an initial permanence value for a synapse. The
// returned permanence will be centered on a normal distribution peaking
// at SynPermConnected with a standard deviation of 0.05.
func (s *V2) getInitPerm() float32 {
	sd := 0.1

	// Determine if the perm should be connected.
	var p float64
	chance := rand.Float64()
	switch {
	case chance <= s.P.InitConnPct:
		// Generate a connected permanence.
		p = rand.NormFloat64()*sd + float64(s.P.SynPermConnected)
		for p < float64(s.P.SynPermConnected) {
			p = rand.NormFloat64()*sd + float64(s.P.SynPermConnected)
		}
	case chance > s.P.InitConnPct:
		// Generate a non-connected permanence.
		p = rand.NormFloat64()*sd + float64(s.P.SynPermConnected)
		for p >= float64(s.P.SynPermConnected) {
			p = rand.NormFloat64()*sd + float64(s.P.SynPermConnected)
		}
	}

	// Clamp 0.0 : 1.0
	switch {
	case p > 1.0:
		p = 1.0
	case p < 0.0:
		p = 0.0
	}

	return float32(p)
}
