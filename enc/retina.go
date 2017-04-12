package enc

import (
	"fmt"
	"log"
)

/* RetinaEncoder
brightness, color, contrast
capture frame --> edge detect --> vectorize --> abstract
*/

/* Encoder design

Takes a single input frame, and outputs a [][]bool with the right overlap properties to be used.

Edge detection is performed on the input frame before sectorizing.

Pixel --> 'Edginess' --> 0-32 Scalar Encoder -->

1. Edge detection
2. Vectorization
3.

*/

// Retina encodes images as SDRs, inspired by the encoding properties
// of a biological retina. Images are converted to black and white.
type Retina struct {
	X, Y int // input dimensions
}

// NewRetina asdf
func NewRetina(x, y int) *Retina {
	return &Retina{
		X: x,
		Y: y,
	}
}

// Encode asdf
func (r *Retina) Encode(f interface{}) []bool {
	return []bool{}
}

// Decode asdf
func (r *Retina) Decode(sv []bool) interface{} {
	return 0
}

// Size asdf
func (r *Retina) Size() {
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
}
