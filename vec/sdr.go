// Package vec provides vector manipulation functions for the htm algorithm.
package vec

// Active returns a slice indices to active vector positions.
func Active(s []bool) (out []int) {
	for i := range s {
		if s[i] {
			out = append(out, i)
		}
	}
	return
}

// Pretty returns a prettified string representation of a vector.
func Pretty(s []bool) (out string) {
	for i := range s {
		if s[i] {
			out += "1"
		} else {
			out += "0"
		}
	}
	return
}

// Sparsity returns a sparsity measure of a vector.
func Sparsity(s []bool) float64 {
	var sum int
	for i := range s {
		if s[i] {
			sum++
		}
	}
	return float64(sum / len(s))
}

// Overlap returns the overlap in bits of two bit vectors.
func Overlap(a, b []bool) int {
	var overlap int
	for i := range a {
		if a[i] == b[i] {
			overlap++
		}
	}
	return overlap
}

// Equal compares two bit vectors and returns true if they are equal,
// false if not. Best case performance is faster than Overlap because
// the function returns as soon as an inequality is detected.
func Equal(a, b []bool) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// VectorUnion returns the union of all provided SparseBinaryVector.
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
