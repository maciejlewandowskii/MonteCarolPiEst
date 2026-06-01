package montecarlo

import "math/rand"

type Point struct {
	X, Y   float64
	Inside bool
}

type Result struct {
	Points      []Point
	Pi          float64
	InsideCount int
}

type PartialResult struct {
	Points      []Point
	InsideCount int
}

func EstimatePi(n int) PartialResult {
	points := make([]Point, n)
	inside := 0
	for i := range points {
		x, y := rand.Float64(), rand.Float64()
		in := x*x+y*y <= 1.0
		if in {
			inside++
		}
		points[i] = Point{X: x, Y: y, Inside: in}
	}
	return PartialResult{
		Points:      points,
		InsideCount: inside,
	}
}
