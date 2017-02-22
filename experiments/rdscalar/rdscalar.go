package main

import (
	"fmt"

	"github.com/nytopop/gohtm/enc"
	"github.com/nytopop/gohtm/vec"
)

func main() {
	r := enc.NewRDScalar(1024, 40, 0, 0.0001)

	_, sine := vec.SineGen(4096, 1.0, 0.02)
	for i := range sine {
		b := r.Encode(sine[i])
		fmt.Println(vec.Pretty(b), sine[i])
	}
	fmt.Println(r.Buckets(), "buckets")
	/*
		var b []bool
		for i := 0.0; i < 1.0; i += 0.1 {
			b = r.Encode(i)
			fmt.Println(vec.Pretty(b))
		}

		time.Sleep(2 * time.Second)

		for i := 0.0; i < 1.0; i += 0.0001 {
			b = r.Encode(i)
			fmt.Println(vec.Pretty(b))
		}
	*/
}
