package math2d

import "testing"

func TestGridNew(t *testing.T) {
	g := NewGrid[int](8, 17)
	if g.Cols() != 8 || g.Rows() != 17 {
		t.Fatalf("expected 8x17, got %dx%d", g.Cols(), g.Rows())
	}
}

func TestGridSetGet(t *testing.T) {
	g := NewGrid[int](8, 17)
	g.Set(3, 5, 42)
	if g.Get(3, 5) != 42 {
		t.Fatalf("expected 42, got %d", g.Get(3, 5))
	}
}

func TestGridGetOutOfBounds(t *testing.T) {
	g := NewGrid[int](8, 17)
	v := g.Get(100, 100)
	if v != 0 {
		t.Fatalf("out-of-bounds Get should return zero value, got %d", v)
	}
}

func TestGridSetOutOfBounds(t *testing.T) {
	g := NewGrid[int](8, 17)
	g.Set(-1, 0, 99)
	g.Set(0, -1, 99)
	g.Set(8, 0, 99)
	g.Set(0, 17, 99)
}

func TestGridClear(t *testing.T) {
	g := NewGrid[int](4, 4)
	g.Set(1, 1, 5)
	g.Clear()
	if g.Get(1, 1) != 0 {
		t.Fatal("Clear should zero all cells")
	}
}

func TestGridFill(t *testing.T) {
	g := NewGrid[int](3, 3)
	g.Fill(7)
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if g.Get(c, r) != 7 {
				t.Fatalf("cell (%d,%d) = %d, want 7", c, r, g.Get(c, r))
			}
		}
	}
}

func TestGridInBounds(t *testing.T) {
	g := NewGrid[int](8, 17)
	if !g.InBounds(0, 0) {
		t.Fatal("(0,0) should be in bounds")
	}
	if !g.InBounds(7, 16) {
		t.Fatal("(7,16) should be in bounds")
	}
	if g.InBounds(8, 0) {
		t.Fatal("(8,0) should be out of bounds")
	}
	if g.InBounds(-1, 0) {
		t.Fatal("(-1,0) should be out of bounds")
	}
}

func TestGridScanRow(t *testing.T) {
	g := NewGrid[int](8, 4)
	g.Set(2, 1, 3)
	g.Set(3, 1, 3)
	g.Set(4, 1, 3)
	g.Set(5, 1, 3)

	runs := g.ScanRow(1, func(v int) int { return v })
	found := false
	for _, run := range runs {
		if run.Value == 3 && run.Length >= 4 {
			found = true
			if run.Start != 2 {
				t.Fatalf("expected run start at col 2, got %d", run.Start)
			}
		}
	}
	if !found {
		t.Fatal("should find a run of 4 threes in row 1")
	}
}

func TestGridScanCol(t *testing.T) {
	g := NewGrid[int](4, 8)
	g.Set(1, 0, 5)
	g.Set(1, 1, 5)
	g.Set(1, 2, 5)
	g.Set(1, 3, 5)
	g.Set(1, 4, 5)

	runs := g.ScanCol(1, func(v int) int { return v })
	found := false
	for _, run := range runs {
		if run.Value == 5 && run.Length >= 4 {
			found = true
		}
	}
	if !found {
		t.Fatal("should find a run of 5 fives in column 1")
	}
}

func TestGridScanRowSkipsZero(t *testing.T) {
	g := NewGrid[int](8, 4)
	g.Set(0, 0, 0)
	g.Set(1, 0, 0)
	g.Set(2, 0, 0)
	g.Set(3, 0, 0)

	runs := g.ScanRow(0, func(v int) int { return v })
	for _, run := range runs {
		if run.Value == 0 && run.Length >= 4 {
			t.Fatal("zero-value runs should be excluded (empty cells)")
		}
	}
}

func TestGridCountValue(t *testing.T) {
	g := NewGrid[int](4, 4)
	g.Set(0, 0, 3)
	g.Set(1, 1, 3)
	g.Set(2, 2, 3)
	g.Set(3, 3, 5)

	count := g.CountValue(3)
	if count != 3 {
		t.Fatalf("expected 3 cells with value 3, got %d", count)
	}
}

func TestGridFindAll(t *testing.T) {
	g := NewGrid[int](4, 4)
	g.Set(0, 0, 3)
	g.Set(2, 1, 3)
	g.Set(3, 3, 3)

	cells := g.FindAll(func(v int) bool { return v == 3 })
	if len(cells) != 3 {
		t.Fatalf("expected 3 matching cells, got %d", len(cells))
	}
}

func TestGridNeighbors(t *testing.T) {
	g := NewGrid[int](4, 4)
	g.Fill(1)

	n := g.Neighbors4(0, 0)
	if len(n) != 2 {
		t.Fatalf("corner should have 2 neighbors, got %d", len(n))
	}

	n = g.Neighbors4(1, 1)
	if len(n) != 4 {
		t.Fatalf("interior should have 4 neighbors, got %d", len(n))
	}
}
