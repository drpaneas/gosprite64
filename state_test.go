package gosprite64

import "testing"

type testState struct {
	name    string
	log     *[]string
}

func (s *testState) Enter() {
	*s.log = append(*s.log, s.name+":enter")
}

func (s *testState) Update() {
	*s.log = append(*s.log, s.name+":update")
}

func (s *testState) Draw() {
	*s.log = append(*s.log, s.name+":draw")
}

func (s *testState) Exit() {
	*s.log = append(*s.log, s.name+":exit")
}

func TestStateMachineSwitchCallsEnterExit(t *testing.T) {
	log := make([]string, 0)
	title := &testState{name: "title", log: &log}
	game := &testState{name: "game", log: &log}

	sm := NewStateMachine(title)
	sm.Init()

	if len(log) != 1 || log[0] != "title:enter" {
		t.Fatalf("expected [title:enter], got %v", log)
	}

	sm.Switch(game)
	expected := []string{"title:enter", "title:exit", "game:enter"}
	if len(log) != 3 {
		t.Fatalf("expected %v, got %v", expected, log)
	}
	for i, e := range expected {
		if log[i] != e {
			t.Fatalf("log[%d] = %q, want %q", i, log[i], e)
		}
	}
}

func TestStateMachineUpdateDrawDelegatesToTop(t *testing.T) {
	log := make([]string, 0)
	title := &testState{name: "title", log: &log}

	sm := NewStateMachine(title)
	sm.Init()
	log = log[:0]

	sm.Update()
	sm.Draw()

	expected := []string{"title:update", "title:draw"}
	if len(log) != 2 {
		t.Fatalf("expected %v, got %v", expected, log)
	}
	for i, e := range expected {
		if log[i] != e {
			t.Fatalf("log[%d] = %q, want %q", i, log[i], e)
		}
	}
}

func TestStateMachinePushPop(t *testing.T) {
	log := make([]string, 0)
	game := &testState{name: "game", log: &log}
	pause := &testState{name: "pause", log: &log}

	sm := NewStateMachine(game)
	sm.Init()
	log = log[:0]

	sm.Push(pause)
	if len(log) != 1 || log[0] != "pause:enter" {
		t.Fatalf("push should call pause:enter, got %v", log)
	}

	sm.Update()
	sm.Draw()
	if log[1] != "pause:update" || log[2] != "pause:draw" {
		t.Fatalf("update/draw should delegate to top (pause), got %v", log)
	}

	log = log[:0]
	sm.Pop()
	if len(log) != 1 || log[0] != "pause:exit" {
		t.Fatalf("pop should call pause:exit, got %v", log)
	}

	log = log[:0]
	sm.Update()
	if log[0] != "game:update" {
		t.Fatalf("after pop, update should go to game, got %v", log)
	}
}

func TestStateMachinePopLastStateIsNoop(t *testing.T) {
	log := make([]string, 0)
	only := &testState{name: "only", log: &log}

	sm := NewStateMachine(only)
	sm.Init()
	log = log[:0]

	sm.Pop()
	if len(log) != 0 {
		t.Fatalf("popping last state should be a no-op, got %v", log)
	}

	sm.Update()
	if log[0] != "only:update" {
		t.Fatalf("should still delegate to the only state, got %v", log)
	}
}

func TestStateMachineSwitchOnEmptyAfterMultiplePops(t *testing.T) {
	log := make([]string, 0)
	a := &testState{name: "a", log: &log}
	b := &testState{name: "b", log: &log}

	sm := NewStateMachine(a)
	sm.Init()
	sm.Push(b)
	log = log[:0]

	sm.Switch(a)
	if log[0] != "b:exit" || log[1] != "a:enter" {
		t.Fatalf("switch should exit top and enter new, got %v", log)
	}
}

func TestStateMachineDepth(t *testing.T) {
	log := make([]string, 0)
	a := &testState{name: "a", log: &log}
	b := &testState{name: "b", log: &log}
	c := &testState{name: "c", log: &log}

	sm := NewStateMachine(a)
	sm.Init()

	if sm.Depth() != 1 {
		t.Fatalf("expected depth 1, got %d", sm.Depth())
	}

	sm.Push(b)
	sm.Push(c)
	if sm.Depth() != 3 {
		t.Fatalf("expected depth 3, got %d", sm.Depth())
	}

	sm.Pop()
	if sm.Depth() != 2 {
		t.Fatalf("expected depth 2, got %d", sm.Depth())
	}
}

func TestStateMachineNilStateIgnored(t *testing.T) {
	log := make([]string, 0)
	a := &testState{name: "a", log: &log}

	sm := NewStateMachine(a)
	sm.Init()
	log = log[:0]

	sm.Push(nil)
	sm.Switch(nil)

	if len(log) != 0 {
		t.Fatalf("nil state operations should be no-ops, got %v", log)
	}
}

func TestStateMachineUpdateDrawWithEmptyStackIsNoop(t *testing.T) {
	sm := &StateMachine{}
	sm.Update()
	sm.Draw()
}

func TestStateMachineCurrent(t *testing.T) {
	log := make([]string, 0)
	a := &testState{name: "a", log: &log}
	b := &testState{name: "b", log: &log}

	sm := NewStateMachine(a)
	sm.Init()

	if sm.Current() != a {
		t.Fatal("Current() should return the initial state")
	}

	sm.Push(b)
	if sm.Current() != b {
		t.Fatal("Current() should return pushed state")
	}

	sm.Pop()
	if sm.Current() != a {
		t.Fatal("Current() should return a after pop")
	}
}

func TestStateMachineNilInitialThenSwitch(t *testing.T) {
	sm := NewStateMachine(nil)
	if sm.Depth() != 0 {
		t.Fatalf("nil initial should give depth 0, got %d", sm.Depth())
	}
	sm.Init()

	log := make([]string, 0)
	a := &testState{name: "a", log: &log}
	sm.Switch(a)
	if sm.Current() != a {
		t.Fatal("switch on empty stack should add the state")
	}
	if log[0] != "a:enter" {
		t.Fatalf("expected a:enter, got %v", log)
	}
}

func TestStateMachineDoubleInit(t *testing.T) {
	log := make([]string, 0)
	a := &testState{name: "a", log: &log}
	sm := NewStateMachine(a)
	sm.Init()
	sm.Init()
	if len(log) != 2 || log[0] != "a:enter" || log[1] != "a:enter" {
		t.Fatalf("double Init calls Enter twice: %v", log)
	}
}
