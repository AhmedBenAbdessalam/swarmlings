package sim

import "math"

type World struct {
	Lings           []Ling
	Width           int
	Height          int
	AvoidanceFactor float64
	AlignmentFactor float64
	GatheringFactor float64
	DetectionRadius float64
	AvoidanceRadius float64
	MaxSpeed        float64
	WallMargin      float64
	WallForce       float64
	grid            *Grid
	neighbors       []Ling
}

func New(lings []Ling, w, h int) World {
	return World{
		Lings:           lings,
		Width:           w,
		Height:          h,
		AvoidanceFactor: 1.0,
		AlignmentFactor: 0.003,
		GatheringFactor: 0.0005,
		DetectionRadius: 100,
		AvoidanceRadius: 20,
		MaxSpeed:        3,
		WallMargin:      75,
		WallForce:       1.5,
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
	w.grid.Populate(w.Lings)

	for i := range w.Lings {
		w.neighbors = w.grid.Neighbors(w.Lings[i].X, w.Lings[i].Y, i, w.Lings, w.neighbors)

		vx, vy := w.Lings[i].Avoid(w.neighbors, w.AvoidanceFactor, w.AvoidanceRadius)
		vx2, vy2 := w.Lings[i].Align(w.neighbors, w.AlignmentFactor, w.DetectionRadius)
		vx += vx2
		vy += vy2
		vx2, vy2 = w.Lings[i].Gather(w.neighbors, w.GatheringFactor, w.DetectionRadius)
		vx += vx2
		vy += vy2
		vx2, vy2 = w.Lings[i].WallAvoid(float64(w.Width), float64(w.Height), w.WallMargin, w.WallForce)
		vx += vx2
		vy += vy2
		w.Lings[i].VX += vx
		w.Lings[i].VY += vy
		speed := math.Hypot(w.Lings[i].VX, w.Lings[i].VY)
		if speed > w.MaxSpeed {
			w.Lings[i].VX = w.Lings[i].VX / speed * w.MaxSpeed
			w.Lings[i].VY = w.Lings[i].VY / speed * w.MaxSpeed
		}

		w.Lings[i].Move()
		w.Lings[i].Clamp(float64(w.Width), float64(w.Height))
	}
}

func (w *World) UpdatePositions(ratioX, ratioY float64) {
	for i := range w.Lings {
		w.Lings[i].X *= ratioX
		w.Lings[i].Y *= ratioY
	}
}
