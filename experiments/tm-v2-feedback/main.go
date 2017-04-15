package main

import (
	"fmt"
	"log"

	"github.com/nytopop/gohtm/tm"
)

/* Temporal Memory w/ Feedback
we need multiple tm regions for testing...

actually do proper error handling
this panic business is preposterous!

-- TM
tm.TemporalMemory -> tm.Interface
Compute(learn bool, args ...[]bool)
  call compute with variadic args
  args[0] is always column activity
  if std   : args[1] is ???
  if apical: args[1] is self.activeCells

New(cols int,

tm.V1 -> tm.TemporalMemory
tm.V2 -> tm.FBTemporalMemory

-- SP
sp.SpatialPooler -> sp.Interface

-- CLA
cla.Classifier -> cla.Interface

-- REGION
region.Region -> region.Interface

// NET
net.Network -> net.Interface

*/

func main() {
	tpar := tm.NewV2Params()
	t := tm.NewV2(tpar)

	for i := 0; i < 16; i++ {
		if err := t.Compute(true,
			make([]bool, 2048),
			make([]bool, 2048*16),
			make([]bool, 2048*16)); err != nil {
			log.Fatalf("%+v", err)
		}

		fmt.Printf("bang %d\n", i)
	}
}
