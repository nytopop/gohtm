package main

import (
	"fmt"

	"github.com/nytopop/gohtm/region"
	"github.com/nytopop/gohtm/vec"
)

func main() {
	r := region.NewV1(region.NewV1Params())

	/* Measurements

	benchmark slowdown per iteration in ms
	100, 105, 110, 115 : 5 ms/iter
	100, 101, 102.01, 103.0301 : 1 %/iter

	*/

	var res [32768]float64
	_, sine := vec.SineGen(32768, 1, 0.05)

	//p := sp.HyperSearchV2(sine, 4)
	//fmt.Println(p)

	for i := 0; i < len(sine); i++ {
		fmt.Println(i)
		fmt.Println("data:", sine[i])

		result := r.Compute(sine[i], true)
		res[i] = result.AnomalyScore

		fmt.Println()
		// get avg of slice 256 entries

		var n = 512
		var avg, total = 0.0, 0.0
		if i >= n {
			for j := range res[i-n : i] {
				total += res[i-n : i][j]
			}
			avg = total / float64(len(res[i-n:i]))
		}

		if avg != 0 {
		}
		//fmt.Println("anom:", result.AnomalyScore)
		//fmt.Println(" avg:", avg)

	}
}
