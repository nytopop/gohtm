package main

import "math/rand"

/* Utility functions
 */

// Check if s contains n
func ContainsInt(n int, s []int) bool {
	for _, v := range s {
		if n == v {
			return true
		}
	}
	return false
}

// Return []int of n unique pseudorandom integers up to max
func UniqueRandInts(n, max int) (rnd []int) {
	for len(rnd) < n {
		r := rand.Intn(max)
		if !ContainsInt(r, rnd) {
			rnd = append(rnd, r)
		}
	}
	return
}
