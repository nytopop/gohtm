package vec

import (
	"math"
	"math/rand"
	"sort"
)

// SineGen generates a sine wave, with n elements, amplitude of amp,
// and noise ratio of noise.
func SineGen(n int, amp, noise float64) ([]float64, []float64) {
	dx := (math.Pi * 2) / 64.0
	scaler := 1.0 / amp
	amp /= 2.0
	theta := 0.0

	x, y := make([]float64, n), make([]float64, n)
	var dirt float64
	for i := 0; i < n; i++ {
		// compute sine
		x[i] = float64(i)
		y[i] = (math.Sin(theta) * amp) + amp

		// compute noise, if any
		dirt = rand.Float64() * scaler * noise
		switch rand.Intn(2) {
		case 0:
			y[i] += dirt
		case 1:
			y[i] -= dirt
		}

		// constrain to 0.0:amp
		switch {
		case y[i] > amp*2.0:
			y[i] = amp * 2.0
		case y[i] < 0.0:
			y[i] = 0.0
		}

		// increment theta
		theta += dx
	}
	return x, y
}

// indice sortable slice
type slice struct {
	sort.Interface
	idx []int
}

func (s slice) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func newSlice(ints []int) *slice {
	s := &slice{
		Interface: sort.IntSlice(ints),
		idx:       make([]int, len(ints)),
	}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

// SortIndices sorts data in ascending order and returns the indices of the
// original unsorted slice in now sorted order.
// SortIndices([]int{16, 8, 4, 2, 1}) returns []int{4, 3, 2, 1, 0}
func SortIndices(data []int) []int {
	sorted := newSlice(data)
	sort.Sort(sorted)
	return sorted.idx
}

// Reverse reverses and returns data.
func Reverse(data []int) []int {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}
