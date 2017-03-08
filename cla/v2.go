package cla

type V2Params struct {
	Steps         []int
	Alpha         float64
	ActValueAlpha float64
}

func NewV2Params() V2Params {
	return V2Params{
		Steps:         []int{1},
		Alpha:         0.001,
		ActValueAlpha: 0.3,
	}
}

func max(a []int) int {
	var max int
	for i := range a {
		if a[i] > max {
			max = a[i]
		}
	}
	return max
}

// V2 Classifier.
type V2 struct {
	P V2Params

	patternHistory [][]int
	maxSteps       int
	maxInputIdx    int
	maxBucketIdx   int
	weightMatrix   [][]float64
}

// NewV2 returns a new V2 Classifier initialized with the provided V2Params.
func NewV2(p V2Params) *V2 {
	weights := make([][]float64, 1)
	weights[0] = make([]float64, 1)

	return &V2{
		P:              p,
		patternHistory: make([][]int, 0),
		maxSteps:       max(p.Steps) + 1,
		maxInputIdx:    0,
		maxBucketIdx:   0,
		weightMatrix:   weights,
	}
}

func (c *V2) Compute(
	sdr []int, bidx int, actValue float64,
	learn, infer bool) Result {

	// store pattern in history, pop off first when we exceed limit
	c.patternHistory = append(c.patternHistory, sdr)
	if len(c.patternHistory) > c.maxSteps {
		c.patternHistory = c.patternHistory[1:]
	}

	return Result{}
}

func (c *V2) infer() {
}

func (c *V2) inferSingleStep() {
}
