package gosprite64

// GameState represents a single game screen or mode (title, menu, gameplay, pause, etc.).
// Implement this interface for each distinct screen in your game.
type GameState interface {
	Enter()
	Update()
	Draw()
	Exit()
}

// StateMachine manages a stack of GameStates. The top state receives
// Update and Draw calls. Push overlays a new state (e.g. pause menu),
// Pop removes it, Switch replaces the top state entirely.
type StateMachine struct {
	stack []GameState
}

// NewStateMachine creates a state machine with the given initial state.
// Call Init() to trigger the initial state's Enter().
func NewStateMachine(initial GameState) *StateMachine {
	sm := &StateMachine{}
	if initial != nil {
		sm.stack = []GameState{initial}
	}
	return sm
}

// Init triggers Enter() on the initial state. Call this once,
// typically inside your Game.Init() implementation.
func (sm *StateMachine) Init() {
	if len(sm.stack) > 0 {
		sm.stack[len(sm.stack)-1].Enter()
	}
}

// Update delegates to the top state's Update.
func (sm *StateMachine) Update() {
	if len(sm.stack) == 0 {
		return
	}
	sm.stack[len(sm.stack)-1].Update()
}

// Draw delegates to the top state's Draw.
func (sm *StateMachine) Draw() {
	if len(sm.stack) == 0 {
		return
	}
	sm.stack[len(sm.stack)-1].Draw()
}

// Switch replaces the top state. Calls Exit on the old top and Enter on the new one.
func (sm *StateMachine) Switch(state GameState) {
	if state == nil {
		return
	}
	if len(sm.stack) > 0 {
		sm.stack[len(sm.stack)-1].Exit()
		sm.stack[len(sm.stack)-1] = state
	} else {
		sm.stack = append(sm.stack, state)
	}
	state.Enter()
}

// Push overlays a new state on top of the current one.
// The current state is NOT exited - it remains in the stack.
// Use this for pause menus, dialog overlays, etc.
func (sm *StateMachine) Push(state GameState) {
	if state == nil {
		return
	}
	sm.stack = append(sm.stack, state)
	state.Enter()
}

// Pop removes the top state and returns to the one below it.
// If only one state remains, Pop is a no-op (the game must always have at least one state).
func (sm *StateMachine) Pop() {
	if len(sm.stack) <= 1 {
		return
	}
	top := sm.stack[len(sm.stack)-1]
	sm.stack = sm.stack[:len(sm.stack)-1]
	top.Exit()
}

// Current returns the top state, or nil if the stack is empty.
func (sm *StateMachine) Current() GameState {
	if len(sm.stack) == 0 {
		return nil
	}
	return sm.stack[len(sm.stack)-1]
}

// Depth returns the number of states on the stack.
func (sm *StateMachine) Depth() int {
	return len(sm.stack)
}
