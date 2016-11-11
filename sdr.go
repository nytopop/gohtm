package main

/* Sparse Distributed Representations
Matrix : 2d
Vector : 1d
Dense  : All bits stored
Sparse : Active bits stored to conserve memory
*/

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

func (sfm SparseFloatMatrix) Key(x, y int) int {
	return (x % sfm.x) + (y % sfm.y) + (x * sfm.x)
}

func (sfm SparseFloatMatrix) Set(x, y int, v float64) {
	sfm.d[sfm.Key(x, y)] = v
}

func (sfm SparseFloatMatrix) Get(x, y int) float64 {
	return sfm.d[sfm.Key(x, y)]
}

func (sfm SparseFloatMatrix) Del(x, y int) {
	delete(sfm.d, sfm.Key(x, y))
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

func (sbm SparseBinaryMatrix) Key(x, y int) int {
	return (x % sbm.x) + (y % sbm.y) + (x * sbm.x)
}

func (sbm SparseBinaryMatrix) Set(x, y int, v bool) {
	if v {
		sbm.d[sbm.Key(x, y)] = true
	} else {
		sbm.Del(x, y)
	}
}

func (sbm SparseBinaryMatrix) Get(x, y int) bool {
	return sbm.d[sbm.Key(x, y)]
}

func (sbm SparseBinaryMatrix) Del(x, y int) {
	delete(sbm.d, sbm.Key(x, y))
}

// ********************************

// Sparse Binary Vector : slice backing
// TODO : convert to map backing
// ********************************
type SparseBinaryVector struct {
	x int
	d []int
}

func NewSparseBinaryVector(x int) SparseBinaryVector {
	sv := SparseBinaryVector{
		x: x,
	}
	// Allocates for 100% sparsity, probably overkill
	//sv.d = make([]int, x)
	return sv
}

func (sv SparseBinaryVector) Pretty() string {
	out := ""
	dv := sv.Dense()
	for x := 0; x < sv.x; x++ {
		if dv[x] {
			out += "1"
		} else {
			out += "0"
		}
	}
	return out
}

// This function is broken, returns '1' in pos 0 of empty vector
func (sv SparseBinaryVector) Dense() []bool {
	dv := make([]bool, sv.x)
	for _, i := range sv.d {
		dv[i] = true
	}
	return dv
}

// ********************************
