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
		maxSteps:       max(p.Steps),
		maxInputIdx:    0,
		maxBucketIdx:   0,
		weightMatrix:   weights,
	}
}

type V2Result struct {
	P float64
}

type V2Results []V2Result

func (v V2Results) Len() int           { return len(v) }
func (v V2Results) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v V2Results) Less(i, j int) bool { return v[i].P < v[j].P }

func (c *V2) Compute(sdr []int, actValue int, learn, infer bool) V2Results {

	return V2Results{}
}

func (c *V2) infer() {
}

func (c *V2) inferSingleStep() {
}
