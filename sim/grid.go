package sim

import "math"

type Grid struct {
	cells    [][]int
	cols     int
	rows     int
	cellSize float64
}

func NewGrid(width, height int, cellSize float64) *Grid {
	if cellSize < 1 {
		cellSize = 1
	}
	cols := int(math.Ceil(float64(width) / cellSize))
	rows := int(math.Ceil(float64(height) / cellSize))
	cells := make([][]int, cols*rows)
	return &Grid{
		cells:    cells,
		cols:     cols,
		rows:     rows,
		cellSize: cellSize,
	}
}

func (g *Grid) Populate(lings []Ling) {
	// Reset length, keep capacity
	for i := range g.cells {
		g.cells[i] = g.cells[i][:0]
	}

	for i, b := range lings {
		col := int(b.X / g.cellSize)
		row := int(b.Y / g.cellSize)
		if col < 0 {
			col = 0
		} else if col >= g.cols {
			col = g.cols - 1
		}
		if row < 0 {
			row = 0
		} else if row >= g.rows {
			row = g.rows - 1
		}
		idx := row*g.cols + col
		g.cells[idx] = append(g.cells[idx], i)
	}
}

func (g *Grid) Neighbors(x, y float64, excludeIndex int, lings []Ling, buf []Ling) []Ling {
	buf = buf[:0]

	col := int(x / g.cellSize)
	row := int(y / g.cellSize)
	if col < 0 {
		col = 0
	} else if col >= g.cols {
		col = g.cols - 1
	}
	if row < 0 {
		row = 0
	} else if row >= g.rows {
		row = g.rows - 1
	}

	for dr := -1; dr <= 1; dr++ {
		nr := row + dr
		if nr < 0 || nr >= g.rows {
			continue
		}
		for dc := -1; dc <= 1; dc++ {
			nc := col + dc
			if nc < 0 || nc >= g.cols {
				continue
			}
			for _, bi := range g.cells[nr*g.cols+nc] {
				if bi != excludeIndex {
					buf = append(buf, lings[bi])
				}
			}
		}
	}
	return buf
}

func (g *Grid) NeedsRebuild(width, height int, cellSize float64) bool {
	newCols := int(math.Ceil(float64(width) / cellSize))
	newRows := int(math.Ceil(float64(height) / cellSize))
	return g.cols != newCols || g.rows != newRows || g.cellSize != cellSize
}
