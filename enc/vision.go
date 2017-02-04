package enc

/* RetinaEncoder
brightness, color, contrast
capture frame --> edge detect --> vectorize --> abstract
*/

/* Encoder design

Takes a single input frame, and outputs a [][]bool with the right overlap properties to be used.

Edge detection is performed on the input frame before sectorizing.

*/

type Frame struct {
	red   []uint8
	green []uint8
	blue  []uint8
}

// RetinaEncoder encodes images as SDRs, inspired by the encoding
// properties of a biological retina. Images are converted to black
// and white.
type RetinaEncoder struct {
	X, Y int // input dimensions
}

func NewRetinaEncoder(x, y int) *RetinaEncoder {
	return &RetinaEncoder{
		X: x,
		Y: y,
	}
}

func (r *RetinaEncoder) Encode(f interface{}) []bool {
	return []bool{}
}

func (r *RetinaEncoder) Decode(sv []bool) interface{} {
	return 0
}
