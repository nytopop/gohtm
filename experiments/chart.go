package main

import (
	"log"
	"math"
	"os"

	chart "github.com/wcharczuk/go-chart"
)

func main() {
	x, y := sineGen(512, 1.0)

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

func sineGen(n int, amp float64) ([]float64, []float64) {
	dx := (math.Pi * 2) / 64.0
	amp /= 2.0
	theta := 0.0

	x, y := make([]float64, n), make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i)
		y[i] = (math.Sin(theta) * amp) + amp
		theta += dx
	}
	return x, y
}
