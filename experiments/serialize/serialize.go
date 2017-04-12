package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/nytopop/gohtm/enc"
	"github.com/nytopop/gohtm/sp"
	"github.com/nytopop/gohtm/tm"
)

/* Serialization
all types should be JSON marshalable
exported fields
  OR
MarshalJSON/UnmarshalJSON methods

encoder
sp
tm + cells
cla

*/

func main() {
	/*
		Hmmm, we can probably do networks with no feedback
	*/

	e := enc.NewRDScalar(1024, 100, 0, 1)

	spar := sp.NewV2Params()
	spar.NumColumns = 16384
	s := sp.NewV2(spar)

	tpar := tm.NewV1Params()
	tpar.NumColumns = spar.NumColumns
	t := tm.NewV1(tpar)

	for i := 0; i < 64; i++ {
		sdr, _ := e.Encode(float64(rand.Intn(16)))
		act := s.Compute(sdr, true)
		t.Compute(act, true)
	}

	dst, _ := json.Marshal(s)
	fmt.Printf("%s\n", string(dst))
	dst, _ = json.Marshal(t)
	fmt.Printf("%s\n", string(dst))
}
