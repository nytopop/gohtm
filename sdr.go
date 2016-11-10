package main

/* Sparse Distributed Representations
Matrix : 2d
Vector : 1d
Dense  : All bits stored
Sparse : Active bits stored to conserve memory
*/

// Sparse Representations
// ********************************
type FloatEntry struct {
	x, y int
	d    float64
}

type SparseFloatMatrix struct {
	x, y int
	d    []FloatEntry
}

func NewSparseFloatMatrix(x, y int) SparseFloatMatrix {
	sm := SparseFloatMatrix{}
	sm.x = x
	sm.y = y
	// Allocates for 100% sparsity, probably overkill
	sm.d = make([]FloatEntry, x*y)
	return sm
}

type BinaryEntry struct {
	x, y int
}

type SparseBinaryMatrix struct {
	x, y int
	d    []BinaryEntry
}

func NewSparseBinaryMatrix(x, y int) SparseBinaryMatrix {
	sm := SparseBinaryMatrix{}
	sm.x = x
	sm.y = y
	// Allocates for 100% sparsity, probably overkill
	sm.d = make([]BinaryEntry, x*y)
	return sm
}

type SparseBinaryVector struct {
	x int
	d []int
}

func NewSparseBinaryVector(x int) SparseBinaryVector {
	sv := SparseBinaryVector{}
	sv.x = x
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
