package sim

import (
	"testing"
)

func TestNewGrid(t *testing.T) {
	g := NewGrid(100, 200, 50)
	if g.cols != 2 {
		t.Errorf("expected 2 cols, got %d", g.cols)
	}
	if g.rows != 4 {
		t.Errorf("expected 4 rows, got %d", g.rows)
	}
	if g.cellSize != 50 {
		t.Errorf("expected cellSize 50, got %f", g.cellSize)
	}
	if len(g.cells) != 8 {
		t.Errorf("expected 8 cells, got %d", len(g.cells))
	}
}

func TestPopulateAndClear(t *testing.T) {
	g := NewGrid(100, 100, 50)
	boids := []Boid{
		{X: 10, Y: 10},  // cell (0,0)
		{X: 60, Y: 10},  // cell (0,1)
		{X: 10, Y: 60},  // cell (1,0)
		{X: 60, Y: 60},  // cell (1,1)
	}
	g.Populate(boids)

	// Each cell should have exactly one boid
	for i, cell := range g.cells {
		if len(cell) != 1 {
			t.Errorf("cell %d: expected 1 boid, got %d", i, len(cell))
		}
	}

	// Populate again should clear first (no double-counting)
	g.Populate(boids)
	for i, cell := range g.cells {
		if len(cell) != 1 {
			t.Errorf("after repopulate, cell %d: expected 1 boid, got %d", i, len(cell))
		}
	}
}

func TestPopulateBoundaryClamp(t *testing.T) {
	g := NewGrid(100, 100, 50)
	// Boid at exact edge should not panic
	boids := []Boid{
		{X: 100, Y: 100},
		{X: 0, Y: 0},
	}
	g.Populate(boids) // should not panic

	// Boid at exact edge should be in the last cell
	lastCell := g.cells[g.rows*g.cols-1]
	if len(lastCell) != 1 {
		t.Errorf("expected boid at (100,100) in last cell, got %d boids", len(lastCell))
	}
	firstCell := g.cells[0]
	if len(firstCell) != 1 {
		t.Errorf("expected boid at (0,0) in first cell, got %d boids", len(firstCell))
	}
}

func TestNeighborsBasic(t *testing.T) {
	g := NewGrid(200, 200, 50)
	boids := []Boid{
		{X: 25, Y: 25},  // cell (0,0) — this is "self"
		{X: 75, Y: 25},  // cell (0,1) — neighbor
		{X: 25, Y: 75},  // cell (1,0) — neighbor
		{X: 125, Y: 125}, // cell (2,2) — not adjacent to (0,0) in 4x4 grid
	}
	g.Populate(boids)

	var buf []Boid
	neighbors := g.Neighbors(25, 25, 0, boids, buf)

	// Should find boids 1 and 2 (adjacent cells), but not 0 (self) or 3 (far)
	if len(neighbors) != 2 {
		t.Fatalf("expected 2 neighbors, got %d", len(neighbors))
	}
}

func TestNeighborsExcludesSelf(t *testing.T) {
	g := NewGrid(200, 200, 50)
	boids := []Boid{
		{X: 25, Y: 25},
		{X: 30, Y: 30}, // same cell
	}
	g.Populate(boids)

	var buf []Boid
	neighbors := g.Neighbors(25, 25, 0, boids, buf)
	if len(neighbors) != 1 {
		t.Fatalf("expected 1 neighbor (self excluded), got %d", len(neighbors))
	}
	if neighbors[0].X != 30 {
		t.Errorf("expected neighbor at X=30, got X=%f", neighbors[0].X)
	}
}

func TestNeighborsToroidalWrap(t *testing.T) {
	g := NewGrid(200, 200, 50) // 4x4 grid
	boids := []Boid{
		{X: 10, Y: 10},   // cell (0,0) — top-left corner
		{X: 190, Y: 10},  // cell (0,3) — top-right corner
		{X: 10, Y: 190},  // cell (3,0) — bottom-left corner
		{X: 190, Y: 190}, // cell (3,3) — bottom-right corner
	}
	g.Populate(boids)

	// Querying from top-left should wrap and find top-right and bottom-left
	var buf []Boid
	neighbors := g.Neighbors(10, 10, 0, boids, buf)

	// Should find boids 1, 2, and 3 — all are in adjacent cells via toroidal wrap
	if len(neighbors) != 3 {
		t.Fatalf("expected 3 neighbors via toroidal wrap, got %d", len(neighbors))
	}
}

func TestNeighborsSmallGrid(t *testing.T) {
	// Grid smaller than 3x3 — cells should be deduped
	g := NewGrid(100, 100, 60) // 2x2 grid (ceil(100/60)=2)
	boids := []Boid{
		{X: 10, Y: 10},
		{X: 70, Y: 70},
	}
	g.Populate(boids)

	var buf []Boid
	neighbors := g.Neighbors(10, 10, 0, boids, buf)
	// With a 2x2 grid and toroidal wrap, all cells are neighbors.
	// Should still find boid 1, not duplicate it.
	if len(neighbors) != 1 {
		t.Fatalf("expected 1 neighbor (deduped small grid), got %d", len(neighbors))
	}
}

func TestNeedsRebuild(t *testing.T) {
	g := NewGrid(100, 100, 50)
	if g.NeedsRebuild(100, 100, 50) {
		t.Error("expected no rebuild needed for same dimensions")
	}
	if !g.NeedsRebuild(200, 100, 50) {
		t.Error("expected rebuild needed for different width")
	}
	if !g.NeedsRebuild(100, 100, 25) {
		t.Error("expected rebuild needed for different cellSize")
	}
}
