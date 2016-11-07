package main

/* Sparse Distributed Representations
Tensor : 3d
Matrix : 2d
Vector : 1d
Dense  : All bits stored
Sparse : Active bits stored to conserve memory
*/

// Dense Representations
// ********************************
type DenseTensor []DenseMatrix

func NewDenseTensor(x, y, z int) DenseTensor {
	t := make(DenseTensor, x)
	for i, _ := range t {
		t[i] = make(DenseMatrix, y)
		for ii, _ := range t[i] {
			t[i][ii] = make(DenseVector, z)
		}
	}
	return t
}

// TODO : func DenseTensor Sparse

type DenseMatrix []DenseVector

func NewDenseMatrix(x, y int) DenseMatrix {
	m := make(DenseMatrix, x)
	for i, _ := range m {
		m[i] = make(DenseVector, y)
	}

	return m
}

func (dm DenseMatrix) Sparse() SparseMatrix {
	sm := SparseMatrix{}
	sm.x = len(dm)
	sm.y = len(dm[0])
	for _, v := range dm {
		sm.d = append(sm.d, v.Sparse())
	}
	return sm
}

type DenseVector []bool

func NewDenseVector(x int) DenseVector {
	v := make(DenseVector, x)
	return v
}

func (dv DenseVector) Sparse() SparseVector {
	sv := SparseVector{}
	sv.x = len(dv)
	for i, v := range dv {
		if v {
			sv.d = append(sv.d, i)
		}
	}
	return sv
}

// ********************************

// Sparse Representations
// ********************************
type SparseTensor struct {
	x, y, z int
	d       []SparseMatrix
}

// TODO : func SparseTensor Dense()
// TODO : func SparseTensor Pretty()

type SparseMatrix struct {
	x, y int
	d    []SparseVector
}

func (sm SparseMatrix) Dense() DenseMatrix {
	dm := NewDenseMatrix(sm.x, sm.y)
	for i, sv := range sm.d {
		dm[i] = sv.Dense()
	}
	return dm
}

func (sm SparseMatrix) Pretty() string {
	out := ""
	dm := sm.Dense()
	for y := 0; y < sm.y; y++ {
		for x := 0; x < sm.x; x++ {
			if dm[x][y] {
				out += "1"
			} else {
				out += "0"
			}
		}
		out += "\n"
	}
	return out
}

type SparseVector struct {
	x int
	d []int
}

func (sv SparseVector) Dense() DenseVector {
	dv := NewDenseVector(sv.x)
	for _, i := range sv.d {
		dv[i] = true
	}
	return dv
}

func (sv SparseVector) Pretty() string {
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

// ********************************

// Annex Code : TODO remove
// ********************************

type Vector []bool

// Create Vector of len x, full of empty values

// Convert Vector to SDR
func (vec Vector) SDR() SDR {
	s := SDR{}
	s.n = len(vec)

	for i, v := range vec {
		if v == true {
			s.w = append(s.w, i)
		}
	}
	return s
}

// Return overlap of two vectors
func Overlap(x, y Vector) int {
	if len(x) != len(y) {
		panic("Mismatched Vectors!")
	}

	o := 0
	for i := 0; i < len(x); i++ {
		if x[i] && y[i] {
			o += 1
		}
	}

	return o
}

// TODO : Return union of input vectors
/*
func Union(x ...Vector) Vector {
	return Vector{}
}*/

type SDR struct {
	n int   // width
	w []int // indices of active bits
}

// Return sparsity of SDR
func (s SDR) Sparsity() float32 {
	return float32(len(s.w)) / float32(s.n)
}

// Convert SDR to Vector
func (s SDR) Vector() Vector {
	v := Vector{}
	for i := 0; i < s.n; i++ {
		v = append(v, false)
	}
	for _, a := range s.w {
		v[a] = true
	}
	return v
}

func (s SDR) Pretty() string {
	out := ""
	vec := s.Vector()
	for _, v := range vec {
		if v {
			out += "1"
		} else {
			out += "0"
		}
	}
	return out
}

// ********************************
