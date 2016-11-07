package main

/* Sparse Distributed Representations
We store only active bits to conserve memory
*/

type Vector []bool

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
