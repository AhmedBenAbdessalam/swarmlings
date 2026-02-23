package sim

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestBound(t *testing.T) {
	testCases := []struct {
		desc string
		Ling Ling
	}{
		{
			desc: "ling moving left past boundary gets clamped",
			Ling: Ling{X: 10, Y: 50, VX: -20, VY: 0},
		},
		{
			desc: "ling moving right past boundary gets clamped",
			Ling: Ling{X: 100, Y: 50, VX: 10, VY: 0},
		},
		{
			desc: "ling moving up past boundary gets clamped",
			Ling: Ling{X: 50, Y: 10, VX: 0, VY: -20},
		},
		{
			desc: "ling moving down past boundary gets clamped",
			Ling: Ling{X: 50, Y: 100, VX: 0, VY: 10},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			world := World{Width: 100, Height: 100, WallMargin: 25, WallForce: 1.5, MaxSpeed: 3}
			world.Lings = []Ling{tC.Ling}
			world.Update()
			ling := world.Lings[0]

			if ling.X < 0 || ling.X > float64(world.Width) || ling.Y < 0 || ling.Y > float64(world.Height) {
				t.Errorf("expected ling to be in bounds, but it was out of bounds: %v", ling)
			}
		})
	}
}

// updateBruteForce is the original O(n^2) implementation for equivalence testing.
func updateBruteForce(w *World) {
	others := make([]Ling, len(w.Lings)-1)
	for i := range w.Lings {
		copy(others[:i], w.Lings[:i])
		copy(others[i:], w.Lings[i+1:])

		vx, vy := w.Lings[i].Avoid(others, w.AvoidanceFactor, w.AvoidanceRadius)
		vx2, vy2 := w.Lings[i].Align(others, w.AlignmentFactor, w.DetectionRadius)
		vx += vx2
		vy += vy2
		vx2, vy2 = w.Lings[i].Gather(others, w.GatheringFactor, w.DetectionRadius)
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

func TestGridEquivalence(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	n := 200
	lings1 := make([]Ling, n)
	lings2 := make([]Ling, n)
	for i := range lings1 {
		l := Ling{
			X:  rng.Float64() * 1000,
			Y:  rng.Float64() * 1000,
			VX: rng.Float64()*4 - 2,
			VY: rng.Float64()*4 - 2,
		}
		lings1[i] = l
		lings2[i] = l
	}

	grid := New(lings1, 1000, 1000)
	brute := New(lings2, 1000, 1000)

	grid.Update()
	updateBruteForce(&brute)

	for i := range lings1 {
		g := grid.Lings[i]
		b := brute.Lings[i]
		if math.Abs(g.X-b.X) > 1e-9 || math.Abs(g.Y-b.Y) > 1e-9 ||
			math.Abs(g.VX-b.VX) > 1e-9 || math.Abs(g.VY-b.VY) > 1e-9 {
			t.Errorf("ling %d diverged:\n  grid:  {X:%.6f Y:%.6f VX:%.6f VY:%.6f}\n  brute: {X:%.6f Y:%.6f VX:%.6f VY:%.6f}",
				i, g.X, g.Y, g.VX, g.VY, b.X, b.Y, b.VX, b.VY)
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	for _, n := range []int{500, 1000, 5000, 10000, 20000} {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			rng := rand.New(rand.NewSource(42))
			lings := make([]Ling, n)
			for j := range lings {
				lings[j] = Ling{
					X:  rng.Float64() * 1000,
					Y:  rng.Float64() * 1000,
					VX: rng.Float64()*2 - 1,
					VY: rng.Float64()*2 - 1,
				}
			}
			world := New(lings, 1000, 1000)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				world.Update()
			}
		})
	}
}
