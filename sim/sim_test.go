package sim

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestBound(t *testing.T) {
	testCases := []struct {
		desc  string
		Boid  Boid
		World World
	}{
		{
			desc:  "boid should wrap around the world when it goes out of bounds",
			Boid:  Boid{X: 10, Y: 50, VX: -20, VY: 0},
			World: World{Width: 100, Height: 100},
		},
		{
			desc:  "boid should wrap around the world when it goes out of bounds",
			Boid:  Boid{X: 100, Y: 50, VX: 10, VY: 0},
			World: World{Width: 100, Height: 100},
		},
		{
			desc:  "boid should wrap around the world when it goes out of bounds",
			Boid:  Boid{X: 50, Y: 10, VX: 0, VY: -20},
			World: World{Width: 100, Height: 100},
		},
		{
			desc:  "boid should wrap around the world when it goes out of bounds",
			Boid:  Boid{X: 50, Y: 100, VX: 0, VY: 10},
			World: World{Width: 100, Height: 100},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			world := tC.World
			world.Boids = []Boid{tC.Boid}
			world.Update()
			boid := world.Boids[0]

			if boid.X < 0 || boid.X > float64(world.Width) || boid.Y < 0 || boid.Y > float64(world.Height) {
				t.Errorf("expected boid to be in bounds, but it was out of bounds: %v", boid)
			}

		})
	}
}

// updateBruteForce is the original O(n²) implementation for equivalence testing.
func updateBruteForce(w *World) {
	others := make([]Boid, len(w.Boids)-1)
	for i := range w.Boids {
		copy(others[:i], w.Boids[:i])
		copy(others[i:], w.Boids[i+1:])

		vx, vy := w.Boids[i].Avoid(others, w.AvoidanceFactor, w.AvoidanceRadius)
		vx2, vy2 := w.Boids[i].Align(others, w.AlignmentFactor, w.DetectionRadius)
		vx += vx2
		vy += vy2
		vx2, vy2 = w.Boids[i].Gather(others, w.GatheringFactor, w.DetectionRadius)
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

func TestGridEquivalence(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	n := 200
	boids1 := make([]Boid, n)
	boids2 := make([]Boid, n)
	for i := range boids1 {
		b := Boid{
			X:  rng.Float64() * 1000,
			Y:  rng.Float64() * 1000,
			VX: rng.Float64()*4 - 2,
			VY: rng.Float64()*4 - 2,
		}
		boids1[i] = b
		boids2[i] = b
	}

	grid := New(boids1, 1000, 1000)
	brute := New(boids2, 1000, 1000)

	grid.Update()
	updateBruteForce(&brute)

	for i := range boids1 {
		g := grid.Boids[i]
		b := brute.Boids[i]
		if math.Abs(g.X-b.X) > 1e-9 || math.Abs(g.Y-b.Y) > 1e-9 ||
			math.Abs(g.VX-b.VX) > 1e-9 || math.Abs(g.VY-b.VY) > 1e-9 {
			t.Errorf("boid %d diverged:\n  grid:  {X:%.6f Y:%.6f VX:%.6f VY:%.6f}\n  brute: {X:%.6f Y:%.6f VX:%.6f VY:%.6f}",
				i, g.X, g.Y, g.VX, g.VY, b.X, b.Y, b.VX, b.VY)
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	for _, n := range []int{500, 1000, 5000, 10000, 20000} {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			rng := rand.New(rand.NewSource(42))
			boids := make([]Boid, n)
			for j := range boids {
				boids[j] = Boid{
					X:  rng.Float64() * 1000,
					Y:  rng.Float64() * 1000,
					VX: rng.Float64()*2 - 1,
					VY: rng.Float64()*2 - 1,
				}
			}
			world := New(boids, 1000, 1000)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				world.Update()
			}
		})
	}
}
