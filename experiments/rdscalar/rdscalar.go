package main

import (
	"fmt"

	"github.com/nytopop/gohtm/enc"
	"github.com/nytopop/gohtm/vec"
)

func main() {
	r := enc.NewRDScalar(2048, 20, 4, 0.0001)
	fmt.Println("Seeding scalar encoder")
	r.Encode(1.0)

	_, sine := vec.SineGen(4096, 1.0, 0.02)
	for i := range sine {
		b := r.Encode(sine[i])
		fmt.Println(vec.ToInt(b), sine[i], r.Buckets())
	}
	fmt.Println(r.Buckets(), "buckets")
}
