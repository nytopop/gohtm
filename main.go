package main

func main() {
	r := NewRegion()
	r.Compute(40.0)

	// test matrix
	//sm := NewDenseMatrix(64, 32).Sparse()
	//fmt.Println(sm.Pretty())

	// test encoding
	/*
		for i := 0; i < 320; i++ {
			s := e.Encode(float64(i))
			d := e.Decode(s)
			if float64(i) != d {
				fmt.Println(i, d)
			}
		}
	*/

	// test overlap
	/*
		sdr := e.Encode(float64(12))
		for i := 0; i < 4096; i++ {
			s := e.Encode(float64(i))
			o := Overlap(sdr.Vector(), s.Vector())
			if o > 3 {
				fmt.Println(i, o)
			}
		}
	*/
}
