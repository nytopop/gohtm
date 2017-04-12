package tm

/* Implementation Details

Things in progress
- fix docs in all code
- enc: refactor, multiencoder, retina, sound
- sp: v2 is ready I think, needs tweaking
- tm: extended TM in progress
- region: region in progress + network
- net: network deps on region, region deps on ETM
*/

// V2 ... extended TM. basal, apical
type V2Params struct {
}

func NewV2Params() V2Params {
	return V2Params{}
}

type V2 struct {
	// Params
	P V2Params

	// State
	// Metrics
}

func NewV2(p V2Params) *V2 {
	return &V2{
		P: p,
	}
}

// feedforward + feedback
func (v *V2) Compute(ff, fb []bool) []bool {
	return []bool{}
}
