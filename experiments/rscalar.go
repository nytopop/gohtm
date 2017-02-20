package main

import (
	"fmt"

	"github.com/nytopop/gohtm/enc"
	"github.com/nytopop/gohtm/vec"
)

func main() {
	r := enc.NewRDScalar(64, 16, 0.001)

	var b []bool
	for i := 0.0; i < 2048; i += 0.1 {
		b = r.Encode(i)
		fmt.Println(vec.Pretty(b))
	}

	fmt.Println(r.Buckets(), "buckets")
}
