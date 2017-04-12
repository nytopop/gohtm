/*
Package net provides an implementation agnostic
interface for building, linking, and managing regions.
*/
package net

import (
	"fmt"
	"log"
	"math"

	"github.com/nytopop/gohtm/region"
)

// Network interface.
type Network interface {
	Compute(data [][]bool) [][]bool
	Serialize() []byte
}

func sizeLayer(w int) int {
	return w
}

// NewNetwork uses the provided column count to generate
// a suitable network topology.
func NewNetwork(cols int) Network {
	var n int
	switch cols%2048 == 0 {
	case true:
		n = cols / 2048
	case false:
		r := float64(cols) / 2048.0
		n = int(math.Ceil(r))
	}

	s := 0
	switch {
	case n%3 == 0:
		s = 3
	case n%2 == 0:
		s = 2
	case n == 1:
		//return NewUnary(n, cols)
		s = 1
	default:
		log.Fatalln(cols, n)
		panic("Unable to find a suitable topology.")
	}
	fmt.Printf("%04d %02d %02d\n", cols, n, s)

	return &Unary{}
}

// Ternary network
// [ ] [ ] [ ] [ ] [ ] [ ]
// [         ] [         ]
// [                     ]
type Ternary struct {
}

// Binary Network
// [ ] [ ] [ ] [ ]
// [     ] [     ]
// [             ]
type Binary struct {
}

// Unary Network
// [ ]
// [ ]
// [ ]
type Unary struct {
	nCols int
	w     int
	grph  []region.Region
}

func NewUnary(w, cols int) Network {
	for i := 1; i <= w; i++ {
		fmt.Printf("%d\n", i)
	}

	return &Unary{
		nCols: cols,
		w:     w,
	}
}

func (u *Unary) Compute(data [][]bool) [][]bool {
	return make([][]bool, 0)
}

func (u *Unary) Serialize() []byte {
	return make([]byte, 256)
}
