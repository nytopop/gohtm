package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/nytopop/gohtm/enc"
	"github.com/nytopop/gohtm/sp"
	"github.com/nytopop/gohtm/tm"
)

func main() {
	// Start CPU profiler
	f, err := os.Create("cpuprofile")
	if err != nil {
		log.Fatalln(err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatalln(err)
	}
	defer pprof.StopCPUProfile()

	// RETINA
	/*
		r := enc.NewRetina(640, 480)
		r.Size()
	*/

	// Temporal Memory Test
	e := enc.NewScalar(enc.NewScalarParams())
	spar := sp.NewV1Params()
	spar.NumInputs = e.Bits
	s := sp.NewV1(spar)
	t := tm.NewV1(tm.NewV1Params())

	seq := make([][]int, 128)
	for i := range seq {
		seq[i] = make([]int, 13)
		for j := range seq[i] {
			seq[i][j] = j * 5
		}
	}

	start := time.Now()
	var v []bool
	for i := range seq {
		//fmt.Println()

		var anom float64
		for j := range seq[i] {
			v = e.Encode(seq[i][j])
			v = s.Compute(v, true)
			t.Compute(v, true)
			anom += t.GetAnomalyScore()
		}
		anom /= float64(len(seq[i]))
		fmt.Println("Avg anomaly score:", anom)
		//t.Reset()

		if i%512 == 0 {
			elap := time.Since(start)

			per := float32(elap.Nanoseconds()) / float32(512) / 1000 / 1000
			fmt.Println(512, "in", elap, "||", per, "ms per iteration")
			fmt.Println(i, "so far")

			start = time.Now()
		}
	}

	elap := time.Since(start)
	per := float32(elap.Nanoseconds()) / float32(len(seq)) / 1000 / 1000
	fmt.Println(len(seq), "in", elap, "||", per, "ms per iteration")

	// Spatial pooler speed test
	/*
		total := time.Now()

		k := make(chan int)
		var wg sync.WaitGroup
		for j := 0; j < 4; j++ {
			wg.Add(1)
			go func(j int) {
				encoder := enc.NewScalar(enc.NewScalarParams())
				spar := sp.NewV1Params()
				spar.NumInputs = encoder.Bits
				spool := sp.NewV1(spar)
				var v []bool
				n := 10000
				for i := 1; i <= n; i++ {
					v = encoder.Encode(rand.Intn(255))
					spool.Compute(v, true)
				}

				wg.Done()

				fmt.Println("ping", n)
				k <- n
			}(j)

		}

		wg.Wait()
		var t int
		var n int
		for i := range k {
			n += i
			t += 1
			if t == 4 {
				close(k)
			}
		}
		elap := time.Since(total)
		rate := int(float64(n) / elap.Seconds())
		fmt.Println(n, "in", elap)
		fmt.Println(rate, "per second")
	*/
}
