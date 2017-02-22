package main

import (
	"log"
	"os"

	"github.com/nytopop/gohtm/vec"
	chart "github.com/wcharczuk/go-chart"
)

func main() {
	x, y := vec.SineGen(512, 1.0, 0.1)

	graph := chart.Chart{
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(2).WithAlpha(64),
					StrokeWidth: 2.0,
				},
				XValues: x,
				YValues: y,
			},
		},
	}

	if file, err := os.Create("chart.png"); err != nil {
		log.Fatalln(err)
	} else {
		defer file.Close()
		if err := graph.Render(chart.PNG, file); err != nil {
			log.Fatalln(err)
		}
	}
}
