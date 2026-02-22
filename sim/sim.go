package sim

import "math"

type World struct {
	Boids           []Boid
	Width           int
	Height          int
	AvoidanceFactor float64
	AlignmentFactor float64
	GatheringFactor float64
	DetectionRadius float64
	AvoidanceRadius float64
	MaxSpeed        float64
	grid            *Grid
	neighbors       []Boid
}

func New(boids []Boid, w, h int) World {
	return World{
		Boids:           boids,
		Width:           w,
		Height:          h,
		AvoidanceFactor: 1.0,
		AlignmentFactor: 0.005,
		GatheringFactor: 0.001,
		DetectionRadius: 100,
		AvoidanceRadius: 20,
		MaxSpeed:        3,
	}
}

func (w *World) Update() {
	cellSize := w.DetectionRadius
	if cellSize < 1 {
		cellSize = 1
	}

	if w.grid == nil || w.grid.NeedsRebuild(w.Width, w.Height, cellSize) {
		w.grid = NewGrid(w.Width, w.Height, cellSize)
	}
	w.grid.Populate(w.Boids)

	for i := range w.Boids {
		w.neighbors = w.grid.Neighbors(w.Boids[i].X, w.Boids[i].Y, i, w.Boids, w.neighbors)

		vx, vy := w.Boids[i].Avoid(w.neighbors, w.AvoidanceFactor, w.AvoidanceRadius)
		vx2, vy2 := w.Boids[i].Align(w.neighbors, w.AlignmentFactor, w.DetectionRadius)
		vx += vx2
		vy += vy2
		vx2, vy2 = w.Boids[i].Gather(w.neighbors, w.GatheringFactor, w.DetectionRadius)
		vx += vx2
		vy += vy2
		w.Boids[i].VX += vx
		w.Boids[i].VY += vy
		speed := math.Hypot(w.Boids[i].VX, w.Boids[i].VY)
		if speed > w.MaxSpeed {
			w.Boids[i].VX = w.Boids[i].VX / speed * w.MaxSpeed
			w.Boids[i].VY = w.Boids[i].VY / speed * w.MaxSpeed
		}

		w.Boids[i].Move()
		w.Boids[i].Wrap(w.Width, w.Height)
	}
}

func (w *World) UpdatePositions(ratioX, ratioY float64) {
	for i := range w.Boids {
		w.Boids[i].X *= ratioX
		w.Boids[i].Y *= ratioY
	}
}
