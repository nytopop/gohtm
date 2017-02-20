package vec

import "sort"

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
