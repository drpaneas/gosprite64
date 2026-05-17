package math2d

// Grid is a generic 2D grid of cells. Useful for tile-based collision,
// puzzle games (Dr. Mario style), and spatial queries.
// Column (X) is first, Row (Y) is second: Get(col, row).
type Grid[T comparable] struct {
	cols  int
	rows  int
	cells []T
}

// NewGrid creates a grid with the given dimensions, all cells zero-valued.
func NewGrid[T comparable](cols, rows int) *Grid[T] {
	return &Grid[T]{
		cols:  cols,
		rows:  rows,
		cells: make([]T, cols*rows),
	}
}

func (g *Grid[T]) Cols() int { return g.cols }
func (g *Grid[T]) Rows() int { return g.rows }

func (g *Grid[T]) InBounds(col, row int) bool {
	return col >= 0 && col < g.cols && row >= 0 && row < g.rows
}

func (g *Grid[T]) index(col, row int) int {
	return row*g.cols + col
}

// Get returns the cell value. Returns zero value if out of bounds.
func (g *Grid[T]) Get(col, row int) T {
	if !g.InBounds(col, row) {
		var zero T
		return zero
	}
	return g.cells[g.index(col, row)]
}

// Set writes a cell value. No-op if out of bounds.
func (g *Grid[T]) Set(col, row int, value T) {
	if !g.InBounds(col, row) {
		return
	}
	g.cells[g.index(col, row)] = value
}

// Clear sets all cells to their zero value.
func (g *Grid[T]) Clear() {
	var zero T
	for i := range g.cells {
		g.cells[i] = zero
	}
}

// Fill sets all cells to the given value.
func (g *Grid[T]) Fill(value T) {
	for i := range g.cells {
		g.cells[i] = value
	}
}

// Run describes a consecutive sequence of cells with the same group value.
type Run struct {
	Start  int
	Length int
	Value  int
}

// ScanRow scans a row for consecutive runs of cells with the same non-zero
// group value. The groupFn maps cell values to group identifiers; cells
// returning 0 are treated as empty and break runs.
func (g *Grid[T]) ScanRow(row int, groupFn func(T) int) []Run {
	if row < 0 || row >= g.rows {
		return nil
	}
	var runs []Run
	currentGroup := 0
	start := 0
	length := 0

	for col := 0; col < g.cols; col++ {
		group := groupFn(g.Get(col, row))
		if group != 0 && group == currentGroup {
			length++
		} else {
			if currentGroup != 0 && length > 0 {
				runs = append(runs, Run{Start: start, Length: length, Value: currentGroup})
			}
			currentGroup = group
			start = col
			length = 1
		}
	}
	if currentGroup != 0 && length > 0 {
		runs = append(runs, Run{Start: start, Length: length, Value: currentGroup})
	}
	return runs
}

// ScanCol scans a column for consecutive runs (vertical).
func (g *Grid[T]) ScanCol(col int, groupFn func(T) int) []Run {
	if col < 0 || col >= g.cols {
		return nil
	}
	var runs []Run
	currentGroup := 0
	start := 0
	length := 0

	for row := 0; row < g.rows; row++ {
		group := groupFn(g.Get(col, row))
		if group != 0 && group == currentGroup {
			length++
		} else {
			if currentGroup != 0 && length > 0 {
				runs = append(runs, Run{Start: start, Length: length, Value: currentGroup})
			}
			currentGroup = group
			start = row
			length = 1
		}
	}
	if currentGroup != 0 && length > 0 {
		runs = append(runs, Run{Start: start, Length: length, Value: currentGroup})
	}
	return runs
}

// GridCell identifies a cell position in the grid.
type GridCell struct {
	Col, Row int
}

// CountValue returns how many cells equal the given value.
func (g *Grid[T]) CountValue(value T) int {
	count := 0
	for _, c := range g.cells {
		if c == value {
			count++
		}
	}
	return count
}

// FindAll returns positions of all cells matching the predicate.
func (g *Grid[T]) FindAll(pred func(T) bool) []GridCell {
	var result []GridCell
	for row := 0; row < g.rows; row++ {
		for col := 0; col < g.cols; col++ {
			if pred(g.Get(col, row)) {
				result = append(result, GridCell{Col: col, Row: row})
			}
		}
	}
	return result
}

// Neighbors4 returns the values of the 4 orthogonal neighbors (up, down, left, right)
// that are in bounds.
func (g *Grid[T]) Neighbors4(col, row int) []T {
	offsets := [4][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	var result []T
	for _, off := range offsets {
		nc, nr := col+off[0], row+off[1]
		if g.InBounds(nc, nr) {
			result = append(result, g.Get(nc, nr))
		}
	}
	return result
}
