package sim

import "math"

func DistanceSquared(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}

func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(DistanceSquared(x1, y1, x2, y2))
}
