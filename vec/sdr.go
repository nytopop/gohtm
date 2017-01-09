package vec

// SDR Utility Functions

// Compute the union of input vectors.
// Returns a SparseBinaryVector comprising all active bits in inputs.
func VectorUnion(input ...SparseBinaryVector) SparseBinaryVector {
	for _, root := range input {
		for _, cmp := range input {
			if root.X != cmp.X {
				panic("Mismatched vector lengths in VectorUnion()!")
			}
		}
	}

	out := NewSparseBinaryVector(input[0].X)

	var bit bool
	for i := 0; i < input[0].X; i++ {
		bit = true
		for _, sbv := range input {
			bit = sbv.Get(i) && bit
		}
		out.Set(i, bit)
	}

	return out
}

/*
// Sparse Float Matrix : map backing
type SparseFloatMatrix struct {
	x, y int
	d    map[int]float64
}

func NewSparseFloatMatrix(x, y int) SparseFloatMatrix {
	return SparseFloatMatrix{
		x: x,
		y: y,
		d: map[int]float64{},
	}
}

func (sfm *SparseFloatMatrix) Key(x, y int) int {
	if x >= sfm.x || y >= sfm.y {
		panic("matrix out of bounds!")
	}
	return (x % sfm.x) + (y % sfm.y) + (x * sfm.x)
}

func (sfm *SparseFloatMatrix) Set(x, y int, v float64) {
	sfm.d[sfm.Key(x, y)] = v
}

func (sfm *SparseFloatMatrix) Get(x, y int) float64 {
	return sfm.d[sfm.Key(x, y)]
}

func (sfm *SparseFloatMatrix) Del(x, y int) {
	delete(sfm.d, sfm.Key(x, y))
}

func (sfm *SparseFloatMatrix) Exists(x, y int) bool {
	if _, ok := sfm.d[sfm.Key(x, y)]; ok {
		return true
	}
	return false
}

// Sparse Binary Matrix : map backing
type SparseBinaryMatrix struct {
	x, y int
	d    map[int]bool
}

func NewSparseBinaryMatrix(x, y int) SparseBinaryMatrix {
	return SparseBinaryMatrix{
		x: x,
		y: y,
		d: map[int]bool{},
	}
}

func (sbm *SparseBinaryMatrix) Key(x, y int) int {
	if x >= sbm.x || y >= sbm.y {
		panic("matrix out of bounds!")
	}
	return (x % sbm.x) + (y % sbm.y) + (x * sbm.x)
}

func (sbm *SparseBinaryMatrix) Set(x, y int, v bool) {
	if v {
		sbm.d[sbm.Key(x, y)] = true
	} else {
		sbm.Del(x, y)
	}
}

func (sbm *SparseBinaryMatrix) Get(x, y int) bool {
	return sbm.d[sbm.Key(x, y)]
}

func (sbm *SparseBinaryMatrix) Del(x, y int) {
	delete(sbm.d, sbm.Key(x, y))
}
*/

// SparseBinaryVector is a sparsely allocated binary vector type,
// storing only active bits.
type SparseBinaryVector struct {
	X int
	d map[int]bool
}

// NewSparseBinaryVector allocates returns new SparseBinaryVector,
// allocated to a maximum length of x.
func NewSparseBinaryVector(x int) SparseBinaryVector {
	return SparseBinaryVector{
		X: x,
		d: map[int]bool{},
	}
}

// Set sets the value of position x in the vector to v. If v is true,
// the value is encoded, if v is false the value will be deleted to
// maintain sparsity.
func (sbv *SparseBinaryVector) Set(x int, v bool) {
	if v {
		sbv.d[x] = true
	} else {
		sbv.Del(x)
	}
}

// Get returns the value of position x in the vector.
func (sbv *SparseBinaryVector) Get(x int) bool {
	return sbv.d[x]
}

// Del deletes the entry at x from the vector.
func (sbv *SparseBinaryVector) Del(x int) {
	delete(sbv.d, x)
}

// Sparsity returns the sparsity of the vector, as computed by
// (activebits / length).
func (sbv *SparseBinaryVector) Sparsity() float64 {
	return float64(len(sbv.d)) / float64(sbv.X)
}

// Subsample returns w random bits from the SBV.
func (sbv *SparseBinaryVector) Subsample(w int) SparseBinaryVector {
	sl := make([]int, len(sbv.d))
	var i int
	for k := range sbv.d {
		sl[i] = k
		i++
	}

	vec := NewSparseBinaryVector(sbv.X)
	sub := UniqueRandInts(w, len(sl))
	for _, v := range sub {
		vec.Set(sl[v], true)
	}

	return vec
}

// Dense converts and returns the vector as a dense []bool representation.
func (sbv *SparseBinaryVector) Dense() []bool {
	dv := make([]bool, sbv.X)
	for i := range dv {
		dv[i] = sbv.Get(i)
	}
	return dv
}

// Pretty returns a stringified representation of the vector.
func (sbv *SparseBinaryVector) Pretty() string {
	out := ""
	for i := 0; i < sbv.X; i++ {
		if sbv.Get(i) {
			out += "1"
		} else {
			out += "0"
		}
	}
	return out
}
