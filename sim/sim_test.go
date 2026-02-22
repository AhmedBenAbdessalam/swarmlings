package sim

import "testing"

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
