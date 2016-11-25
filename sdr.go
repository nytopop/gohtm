package main

// SDR Utility Functions
// ********************************

/* Compute the union of input vectors. Returns a SparseBinaryVector comprising all active bits in inputs. */
func VectorUnion(input ...SparseBinaryVector) SparseBinaryVector {
	for _, root := range input {
		for _, cmp := range input {
			if root.x != cmp.x {
				panic("Mismatched vector lengths in VectorUnion()!")
			}
		}
	}

	out := NewSparseBinaryVector(input[0].x)

	var bit bool
	for i := 0; i < input[0].x; i++ {
		bit = true
		for _, sbv := range input {
			bit = sbv.Get(i) && bit
		}
		out.Set(i, bit)
	}

	return out
}

// ********************************

// Sparse Float Matrix : map backing
// ********************************
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

// ********************************

// Sparse Binary Matrix : map backing
// ********************************
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

// ********************************

// Sparse Binary Vector : map backing
// ********************************
type SparseBinaryVector struct {
	x int
	d map[int]bool
}

func NewSparseBinaryVector(x int) SparseBinaryVector {
	return SparseBinaryVector{
		x: x,
		d: map[int]bool{},
	}
}

func (sbv *SparseBinaryVector) Set(x int, v bool) {
	if v {
		sbv.d[x] = true
	} else {
		sbv.Del(x)
	}
}

func (sbv *SparseBinaryVector) Get(x int) bool {
	return sbv.d[x]
}

func (sbv *SparseBinaryVector) Del(x int) {
	delete(sbv.d, x)
}

func (sbv *SparseBinaryVector) Sparsity() float64 {
	return float64(len(sbv.d)) / float64(sbv.x)
}

func (sbv *SparseBinaryVector) Dense() []bool {
	dv := make([]bool, sbv.x)
	for i, _ := range dv {
		dv[i] = sbv.Get(i)
	}
	return dv
}

func (sbv *SparseBinaryVector) Pretty() string {
	out := ""
	for i := 0; i < sbv.x; i++ {
		if sbv.Get(i) {
			out += "1"
		} else {
			out += "0"
		}
	}
	return out
}

// ********************************
