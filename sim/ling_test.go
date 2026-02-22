package sim

import (
	"testing"
)

func TestAvoid(t *testing.T) {
	testCases := []struct {
		desc       string
		Ling       Ling
		Lings      []Ling
		expectedVX float64
		expectedVY float64
	}{
		{
			desc: "ling should avoid other lings that are too close",
			Ling: Ling{X: 50, Y: 50, VX: 0, VY: 0},
			Lings: []Ling{
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
			desc: "ling should not avoid other lings that are far enough",
			Ling: Ling{X: 50, Y: 50, VX: 0, VY: 0},
			Lings: []Ling{
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
			vx, vy := tC.Ling.Avoid(tC.Lings[1:], 1.0, 20)
			if vx != tC.expectedVX || vy != tC.expectedVY {
				t.Errorf("expected ling velocity to be vx=%v, vy=%v, but got vx=%v, vy=%v", tC.expectedVX, tC.expectedVY, vx, vy)
			}
		})
	}
}

func TestAlign(t *testing.T) {
	testCases := []struct {
		desc       string
		Ling       Ling
		Lings      []Ling
		expectedVX float64
		expectedVY float64
	}{
		{
			desc: "ling should align with other lings that are close enough",
			Ling: Ling{X: 50, Y: 50, VX: 0, VY: 0},
			Lings: []Ling{
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
			desc: "ling should not align with other lings that are far enough",
			Ling: Ling{X: 50, Y: 50, VX: 0, VY: 0},
			Lings: []Ling{
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
			vx, vy := tC.Ling.Align(tC.Lings[1:], 1.0, 20)
			if vx != tC.expectedVX || vy != tC.expectedVY {
				t.Errorf("expected ling velocity to be vx=%v, vy=%v, but got vx=%v, vy=%v", tC.expectedVX, tC.expectedVY, vx, vy)
			}
		})
	}
}

func TestGather(t *testing.T) {
	testCases := []struct {
		desc  string
		Ling  Ling
		Lings []Ling
		expectedVX float64
		expectedVY float64
	}{
		{
			desc: "ling should gather with other lings that are close enough",
			Ling: Ling{X: 50, Y: 50, VX: 0, VY: 0},
			Lings: []Ling{
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
			desc: "ling should not gather with other lings that are far enough",
			Ling: Ling{X: 50, Y: 50, VX: 0, VY: 0},
			Lings: []Ling{
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
			vx, vy := tC.Ling.Gather(tC.Lings[1:], 1.0, 20)
			if vx != tC.expectedVX || vy != tC.expectedVY {
				t.Errorf("expected ling velocity to be vx=%v, vy=%v, but got vx=%v, vy=%v", tC.expectedVX, tC.expectedVY, vx, vy)
			}
		})
	}
}
