// Package vec provides vector manipulation functions for the htm algorithm.
package vec

import "sort"

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

// Sparsity returns a sparsity measure of a []bool.
func Sparsity(s []bool) float64 {
	var sum int
	for i := range s {
		if s[i] {
			sum++
		}
	}
	return float64(sum) / float64(len(s))
}

// ToInt returns a sorted []int of the indices of any true bits in a.
func ToInt(a []bool) []int {
	var b []int
	for i := range a {
		if a[i] {
			b = append(b, i)
		}
	}
	sort.Ints(b)
	return b
}

// ToBool converts a
func ToBool(a []int, n int) []bool {
	b := make([]bool, n)
	for i := range a {
		b[a[i]] = true
	}
	return b
}
func ToBool32(a []uint32, n uint32) []bool {
	b := make([]bool, n)
	for i := range a {
		b[a[i]] = true
	}
	return b
}

// Contains returns true if a contains b, false if not.
// This is intended for randomly distributed values, so
// O(n) is the best we can do.
func Contains(a []int, b int) bool {
	for i := range a {
		if a[i] == b {
			return true
		}
	}
	return false
}
func Contains32(a []uint32, b uint32) bool {

	for i := range a {
		if a[i] == b {
			return true
		}
	}
	return false
}

// Overlap returns the number of overlapping entries in a and b. It is
// assumed that a and b are sorted in ascending order.
func Overlap(a, b []int) int {
	var overlap int
	for i := range a {
		j := sort.Search(
			len(b),
			func(j int) bool {
				return b[j] >= a[i]
			})
		if !(j >= len(b) || b[j] != a[i]) {
			overlap++
		}
	}
	return overlap
}

func Overlap32(a, b []uint32) int {
	var overlap int
	for i := range a {
		for j := range b {
			if a[i] == b[j] {
				overlap++
			}
		}
	}
	return overlap
}

// Equal returns the equality of two []int slices. Slightly faster
// than binary searching the whole
func Equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
