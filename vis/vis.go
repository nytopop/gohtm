// Package vis implements visualization primitives for gohtm.
package vis

import (
	"fmt"
	"log"

	"github.com/nytopop/gohtm/sp"
	"github.com/nytopop/gohtm/tm"
	"github.com/veandco/go-sdl2/sdl"
)

type Visualizer struct {
	s sp.SpatialPooler
	t tm.TemporalMemory
}

func Draw() {
	var window *sdl.Window
	var r *sdl.Renderer
	var points []sdl.Point
	//var rect sdl.Rect
	//var rects []sdl.Rect

	window, err := sdl.CreateWindow("gothm",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		300, 300, 0x00)

	if err != nil {
		log.Fatalln(err)
	}
	defer window.Destroy()

	window.Maximize()
	x, y := window.GetMaximumSize()
	fmt.Println(x, y)

	r, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Destroy()

	r.Clear()

	r.SetDrawColor(255, 255, 0, 255)
	points = []sdl.Point{{0, 0}, {100, 300}, {100, 300}, {200, 0}}
	r.DrawLines(points)

	r.Present()

	sdl.Delay(5000)
}
