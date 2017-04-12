package region

import (
	"github.com/nytopop/gohtm/enc"
	"github.com/nytopop/gohtm/sp"
)

/* V2 region
spatial abstraction
spatial scoping
temporal scoping



*/

// Bot. no scoping, bUpDown
type Bot struct {
	inputEncoder enc.Encoder      // all have this
	inputPooler  sp.SpatialPooler // all have this

	// Compute(learn bool)
}

// Mid. ts scoping, bUpDown, tDownUp
type Mid struct {
}

// Top. ts scoping, tDownUp
type Top struct {
}
