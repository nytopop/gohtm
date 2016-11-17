package main

func main() {
	r := NewRegion()
	//	r.Compute(48.0)

	for i := 0.0; i <= 128; i++ {
		r.Compute(i)
	}

}
