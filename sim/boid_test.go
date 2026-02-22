package sim

import (
	"testing"
)

func TestAvoid(t *testing.T) {
	testCases := []struct {
		desc       string
		Boid       Boid
		Boids      []Boid
		expectedVX float64
		expectedVY float64
	}{
		{
			desc: "boid should avoid other boids that are too close",
			Boid: Boid{X: 50, Y: 50, VX: 0, VY: 0},
			Boids: []Boid{
				{X: 50, Y: 50, VX: 0, VY: 0},
				{X: 60, Y: 50, VX: 0, VY: 0},
				{X: 60, Y: 50, VX: 0, VY: 0},
				{X: 50, Y: 60, VX: 0, VY: 0},
				{X: 50, Y: 60, VX: 0, VY: 0},
			},
			expectedVX: -0.2,
			expectedVY: -0.2,
		},
		{
			desc: "boid should not avoid other boids that are far enough",
			Boid: Boid{X: 50, Y: 50, VX: 0, VY: 0},
			Boids: []Boid{
				{X: 50, Y: 50, VX: 0, VY: 0},
				{X: 80, Y: 50, VX: 0, VY: 0},
				{X: 80, Y: 50, VX: 0, VY: 0},
				{X: 50, Y: 80, VX: 0, VY: 0},
				{X: 50, Y: 80, VX: 0, VY: 0},
			},
			expectedVX: 0.0,
			expectedVY: 0.0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vx, vy := tC.Boid.Avoid(tC.Boids[1:], 1.0, 20)
			if vx != tC.expectedVX || vy != tC.expectedVY {
				t.Errorf("expected boid velocity to be vx=%v, vy=%v, but got vx=%v, vy=%v", tC.expectedVX, tC.expectedVY, vx, vy)
			}
		})
	}
}

func TestAlign(t *testing.T) {
	testCases := []struct {
		desc       string
		Boid       Boid
		Boids      []Boid
		expectedVX float64
		expectedVY float64
	}{
		{
			desc: "boid should align with other boids that are close enough",
			Boid: Boid{X: 50, Y: 50, VX: 0, VY: 0},
			Boids: []Boid{
				{X: 50, Y: 50, VX: 0, VY: 0},
				{X: 55, Y: 50, VX: 1, VY: 0},
				{X: 45, Y: 50, VX: 1, VY: 0},
				{X: 50, Y: 55, VX: 1, VY: 0},
				{X: 50, Y: 45, VX: 1, VY: 0},
			},
			expectedVX: 1.0,
			expectedVY: 0.0,
		},
		{
			desc: "boid should not align with other boids that are far enough",
			Boid: Boid{X: 50, Y: 50, VX: 0, VY: 0},
			Boids: []Boid{
				{X: 50, Y: 50, VX: 0, VY: 0},
				{X: 80, Y: 50, VX: 1, VY: 0},
				{X: 20, Y: 50, VX: 1, VY: 0},
				{X: 50, Y: 80, VX: 1, VY: 0},
				{X: 50, Y: 20, VX: 1, VY: 0},
			},
			expectedVX: 0.0,
			expectedVY: 0.0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vx, vy := tC.Boid.Align(tC.Boids[1:], 1.0, 20)
			if vx != tC.expectedVX || vy != tC.expectedVY {
				t.Errorf("expected boid velocity to be vx=%v, vy=%v, but got vx=%v, vy=%v", tC.expectedVX, tC.expectedVY, vx, vy)
			}
		})
	}
}

func TestGather(t *testing.T) {
	testCases := []struct {
		desc  string
		Boid  Boid
		Boids []Boid
		expectedVX float64
		expectedVY float64
	}{
		{
			desc: "boid should gather with other boids that are close enough",
			Boid: Boid{X: 50, Y: 50, VX: 0, VY: 0},
			Boids: []Boid{
				{X: 50, Y: 50, VX: 0, VY: 0},
				{X: 55, Y: 50, VX: 0, VY: 0},
				{X: 55, Y: 50, VX: 0, VY: 0},
				{X: 50, Y: 55, VX: 0, VY: 0},
				{X: 50, Y: 55, VX: 0, VY: 0},
			},
			expectedVX: 2.5,
			expectedVY: 2.5,
		},
		{
			desc: "boid should not gather with other boids that are far enough",
			Boid: Boid{X: 50, Y: 50, VX: 0, VY: 0},
			Boids: []Boid{
				{X: 50, Y: 50, VX: 0, VY: 0},
				{X: 80, Y: 50, VX: 0, VY: 0},
				{X: 20, Y: 50, VX: 0, VY: 0},
				{X: 50, Y: 80, VX: 0, VY: 0},
				{X: 50, Y: 20, VX: 0, VY: 0},
			},
			expectedVX: 0.0,
			expectedVY: 0.0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vx, vy := tC.Boid.Gather(tC.Boids[1:], 1.0, 20)
			if vx != tC.expectedVX || vy != tC.expectedVY {
				t.Errorf("expected boid velocity to be vx=%v, vy=%v, but got vx=%v, vy=%v", tC.expectedVX, tC.expectedVY, vx, vy)
			}
		})
	}
}
