package gosprite64

import "testing"

func TestMenuNavDown(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "Start"},
		{Label: "Options"},
		{Label: "Quit"},
	})
	if m.Cursor() != 0 {
		t.Fatalf("initial cursor should be 0, got %d", m.Cursor())
	}
	m.MoveDown()
	if m.Cursor() != 1 {
		t.Fatalf("after MoveDown, cursor should be 1, got %d", m.Cursor())
	}
	m.MoveDown()
	if m.Cursor() != 2 {
		t.Fatalf("after second MoveDown, cursor should be 2, got %d", m.Cursor())
	}
}

func TestMenuNavUp(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "Start"},
		{Label: "Options"},
		{Label: "Quit"},
	})
	m.MoveDown()
	m.MoveDown()
	m.MoveUp()
	if m.Cursor() != 1 {
		t.Fatalf("after MoveUp, cursor should be 1, got %d", m.Cursor())
	}
}

func TestMenuWrapDown(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "A"},
		{Label: "B"},
	})
	m.Wrap = true
	m.MoveDown()
	m.MoveDown()
	if m.Cursor() != 0 {
		t.Fatalf("wrap-around should go to 0, got %d", m.Cursor())
	}
}

func TestMenuWrapUp(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "A"},
		{Label: "B"},
		{Label: "C"},
	})
	m.Wrap = true
	m.MoveUp()
	if m.Cursor() != 2 {
		t.Fatalf("wrap-up should go to last, got %d", m.Cursor())
	}
}

func TestMenuNoWrapDown(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "A"},
		{Label: "B"},
	})
	m.MoveDown()
	m.MoveDown()
	if m.Cursor() != 1 {
		t.Fatalf("no-wrap should clamp at last, got %d", m.Cursor())
	}
}

func TestMenuNoWrapUp(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "A"},
		{Label: "B"},
	})
	m.MoveUp()
	if m.Cursor() != 0 {
		t.Fatalf("no-wrap should clamp at 0, got %d", m.Cursor())
	}
}

func TestMenuSelected(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "Start"},
		{Label: "Options"},
	})
	m.MoveDown()
	sel := m.Selected()
	if sel.Label != "Options" {
		t.Fatalf("expected 'Options', got %q", sel.Label)
	}
}

func TestMenuConfirm(t *testing.T) {
	confirmed := ""
	m := NewMenu([]MenuItem{
		{Label: "Start", OnConfirm: func() { confirmed = "start" }},
		{Label: "Options", OnConfirm: func() { confirmed = "options" }},
	})
	m.MoveDown()
	m.Confirm()
	if confirmed != "options" {
		t.Fatalf("expected 'options', got %q", confirmed)
	}
}

func TestMenuConfirmNilCallback(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "Start"},
	})
	m.Confirm()
}

func TestMenuEmpty(t *testing.T) {
	m := NewMenu(nil)
	m.MoveDown()
	m.MoveUp()
	m.Confirm()
	if m.Cursor() != 0 {
		t.Fatal("empty menu cursor should stay at 0")
	}
}

func TestMenuSetCursor(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "A"},
		{Label: "B"},
		{Label: "C"},
	})
	m.SetCursor(2)
	if m.Cursor() != 2 {
		t.Fatalf("expected 2, got %d", m.Cursor())
	}
	m.SetCursor(99)
	if m.Cursor() != 2 {
		t.Fatalf("out-of-range should clamp to last, got %d", m.Cursor())
	}
	m.SetCursor(-1)
	if m.Cursor() != 0 {
		t.Fatalf("negative should clamp to 0, got %d", m.Cursor())
	}
}

func TestMenuCount(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "A"},
		{Label: "B"},
	})
	if m.Count() != 2 {
		t.Fatalf("expected 2, got %d", m.Count())
	}
}

func TestMenuDisabledItemSkipped(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "A"},
		{Label: "B", Disabled: true},
		{Label: "C"},
	})
	m.MoveDown()
	if m.Cursor() != 2 {
		t.Fatalf("disabled item should be skipped, cursor should be 2, got %d", m.Cursor())
	}
}

func TestMenuDisabledItemSkippedUp(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "A"},
		{Label: "B", Disabled: true},
		{Label: "C"},
	})
	m.SetCursor(2)
	m.MoveUp()
	if m.Cursor() != 0 {
		t.Fatalf("disabled item should be skipped going up, cursor should be 0, got %d", m.Cursor())
	}
}

func TestMenuAllDisabledNoWrap(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "A", Disabled: true},
		{Label: "B", Disabled: true},
		{Label: "C", Disabled: true},
	})
	m.MoveDown()
	if m.Cursor() != 0 {
		t.Fatalf("all disabled: MoveDown should stay at 0, got %d", m.Cursor())
	}
	m.MoveUp()
	if m.Cursor() != 0 {
		t.Fatalf("all disabled: MoveUp should stay at 0, got %d", m.Cursor())
	}
	m.Confirm()
}

func TestMenuAllDisabledWrap(t *testing.T) {
	m := NewMenu([]MenuItem{
		{Label: "A", Disabled: true},
		{Label: "B", Disabled: true},
	})
	m.Wrap = true
	m.MoveDown()
	if m.Cursor() != 0 {
		t.Fatalf("all disabled wrap: MoveDown should stay at 0, got %d", m.Cursor())
	}
	m.MoveUp()
	if m.Cursor() != 0 {
		t.Fatalf("all disabled wrap: MoveUp should stay at 0, got %d", m.Cursor())
	}
}
