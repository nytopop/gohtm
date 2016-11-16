package main

import (
	"math/rand"
	"sort"
)

/* Utility functions */

/* Check if []int s contains int s. Returns true if yes. */
func ContainsInt(n int, s []int) bool {
	for _, v := range s {
		if n == v {
			return true
		}
	}
	return false
}

/* Return slice of n random integers, up to max size. All integers returned will be unique, that is, there should be no duplicates in the returned slice. */
func UniqueRandInts(n, max int) (rnd []int) {
	for len(rnd) < n {
		r := rand.Intn(max)
		if !ContainsInt(r, rnd) {
			rnd = append(rnd, r)
		}
	}
	return
}

// indice sortable slice
type Slice struct {
	sort.Interface
	idx []int
}

func (s Slice) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func NewSlice(ints []int) *Slice {
	s := &Slice{
		Interface: sort.IntSlice(ints),
		idx:       make([]int, len(ints)),
	}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

/* Sorts the input slice in ascending order, and returns the indices of the original slice in sorted order. */
func SortIndices(data []int) []int {
	sorted := NewSlice(data)
	sort.Sort(sorted)
	return sorted.idx
}

/* Returns a reversed version of the input slice. */
func Reverse(data []int) []int {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}
