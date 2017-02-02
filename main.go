package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/nytopop/gohtm/enc"
	"github.com/nytopop/gohtm/sp"
)

func main() {
	// Start CPU profiler
	f, err := os.Create("cpuprofile")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	epar := enc.NewScalarEncoderParams()
	encoder := enc.NewScalarEncoder(epar)

	spar := sp.NewSpatialParams()
	spar.NumColumns = 2048
	spar.NumInputs = encoder.Bits
	spool := sp.NewSpatialPooler(spar)

	n := 10000

	total := time.Now()
	start := time.Now()

	var v []bool
	for i := 1; i <= n; i++ {
		// encoding
		v = encoder.Encode(rand.Intn(255))
		//encoder.Decode(v)

		// spatial pooling
		spool.Compute(v, true)

		// temporal memory

		if i%1000 == 0 {
			fmt.Println(i, "in", time.Since(start))
			start = time.Now()
		}
	}

	fmt.Println(n, "in", time.Since(total))
}
