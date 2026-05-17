package gosprite64

// FrameInput captures the controller state for one player in one frame.
type FrameInput struct {
	Buttons ButtonMask
	StickX  int8
	StickY  int8
}

// ReplayData holds a complete recorded input sequence for all players.
type ReplayData struct {
	PlayerCount int
	FrameCount  int
	frames      [][]FrameInput // [player][frame]
}

// InputRecorder captures per-frame controller state during gameplay.
type InputRecorder struct {
	playerCount int
	frames      [][]FrameInput
	frameCount  int
}

// NewInputRecorder creates a recorder for the given number of players.
func NewInputRecorder(playerCount int) *InputRecorder {
	if playerCount <= 0 {
		playerCount = 1
	}
	return &InputRecorder{
		playerCount: playerCount,
		frames:      make([][]FrameInput, playerCount),
	}
}

// CaptureFrame records one frame of input for the given player.
func (r *InputRecorder) CaptureFrame(player int, input FrameInput) {
	if r == nil || player < 0 || player >= r.playerCount {
		return
	}
	r.frames[player] = append(r.frames[player], input)
	maxLen := 0
	for _, pf := range r.frames {
		if len(pf) > maxLen {
			maxLen = len(pf)
		}
	}
	r.frameCount = maxLen
}

// Finish finalizes the recording and returns the replay data.
func (r *InputRecorder) Finish() *ReplayData {
	if r == nil {
		return &ReplayData{}
	}
	copied := make([][]FrameInput, r.playerCount)
	for i, pf := range r.frames {
		copied[i] = make([]FrameInput, len(pf))
		copy(copied[i], pf)
	}
	return &ReplayData{
		PlayerCount: r.playerCount,
		FrameCount:  r.frameCount,
		frames:      copied,
	}
}

// InputPlayer replays recorded input frame by frame.
type InputPlayer struct {
	data    *ReplayData
	cursors []int
}

// NewInputPlayer creates a player for the given replay data.
func NewInputPlayer(data *ReplayData) *InputPlayer {
	if data == nil {
		data = &ReplayData{}
	}
	cursors := make([]int, data.PlayerCount)
	return &InputPlayer{
		data:    data,
		cursors: cursors,
	}
}

// NextFrame returns the next frame of input for the given player.
// Returns false when all frames have been consumed for that player.
func (p *InputPlayer) NextFrame(player int) (FrameInput, bool) {
	if p == nil || p.data == nil || player < 0 || player >= p.data.PlayerCount {
		return FrameInput{}, false
	}
	cursor := p.cursors[player]
	if cursor >= len(p.data.frames[player]) {
		return FrameInput{}, false
	}
	input := p.data.frames[player][cursor]
	p.cursors[player]++
	return input, true
}

// Done returns true when all players have consumed all their frames.
func (p *InputPlayer) Done() bool {
	if p == nil || p.data == nil || p.data.FrameCount == 0 {
		return true
	}
	for i := 0; i < p.data.PlayerCount; i++ {
		if p.cursors[i] < len(p.data.frames[i]) {
			return false
		}
	}
	return true
}

// Reset restarts playback from the beginning.
func (p *InputPlayer) Reset() {
	if p == nil {
		return
	}
	for i := range p.cursors {
		p.cursors[i] = 0
	}
}

// CurrentFrame returns the current playback position (frame index of player 0).
func (p *InputPlayer) CurrentFrame() int {
	if p == nil || len(p.cursors) == 0 {
		return 0
	}
	return p.cursors[0]
}
