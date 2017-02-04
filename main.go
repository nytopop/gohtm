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

	/* sizing
	assert(fieldsize % sectorsize == 0 || boxSize)
	numsectors == (fieldsize / boxsize) - 1

	*/
	// Split visual field into (sectorSize / 2)
	fieldSize := 512
	sectorSize := 16
	boxSize := sectorSize / 2

	// Ensure receptive fields sized approprately
	mod := fieldSize % sectorSize
	if (mod != 0) && (mod != boxSize) {
		log.Fatalln("Mismatched field / sector size")
	}

	// calc box / sector sizes
	// numBoxes := fieldSize / boxSize
	numSectors := fieldSize / sectorSize

	// generate bounds
	sectors := make([]int, numSectors)
	for i := range sectors {
		// start at 0th, end at sectorSizeth
		sectors[i] = i * sectorSize
	}

	fmt.Println(sectors)

	// Spatial pooler speed test
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

	elap := time.Since(total)
	rate := int(float64(n) / elap.Seconds())
	fmt.Println(n, "in", elap)
	fmt.Println(rate, "per second")
}
