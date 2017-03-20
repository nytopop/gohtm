package main

import (
	"fmt"

	"github.com/nytopop/gohtm/enc"
	"github.com/nytopop/gohtm/sp"
	"github.com/nytopop/gohtm/vec"
)

func main() {
	e := enc.NewRDScalar(64, 7, 0, 1)
	sr := sp.NewV2(sp.NewV2Params())

	inputs := []float64{}
	for i := 0; i < 16; i++ {
		inputs = append(inputs, float64(i))
	}

	for i := range inputs {
		//fmt.Println(inputs[i])

		vector, _ := e.Encode(inputs[i])
		cols := sr.Compute(vector, true)

		//fmt.Println(vec.Pretty(vector))
		fmt.Println(vec.Pretty(cols))
	}
}
