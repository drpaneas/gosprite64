package gosprite64

import tilerender "github.com/drpaneas/gosprite64/internal/tile2d/render"

type runtimeState struct {
	video *videoState
	audio *audioState
	tile  *tileRuntime
}

type tileRuntime struct {
	renderer *tilerender.Renderer
	textured tilerender.TexturedSetupState
}

var activeRuntime *runtimeState

func newTileRuntime() *tileRuntime {
	return &tileRuntime{
		renderer: tilerender.NewRenderer(tilerender.RenderHooks{}),
	}
}

func newRuntimeState() *runtimeState {
	return &runtimeState{
		tile: newTileRuntime(),
	}
}

func activateRuntime(rt *runtimeState) {
	activeRuntime = rt
}

func currentRuntime() *runtimeState {
	return activeRuntime
}

func (rt *runtimeState) currentVideo() *videoState {
	if rt == nil {
		return nil
	}
	return rt.video
}

func currentVideo() *videoState {
	return currentRuntime().currentVideo()
}

func (rt *runtimeState) currentTile() *tileRuntime {
	if rt == nil {
		return nil
	}
	return rt.tile
}

func currentTile() *tileRuntime {
	return currentRuntime().currentTile()
}

func (t *tileRuntime) resetTexturedState() {
	if t == nil {
		return
	}
	t.textured = tilerender.TexturedSetupState{}
}

func (rt *runtimeState) currentAudio() *audioState {
	if rt == nil {
		return nil
	}
	return rt.audio
}

func currentAudio() *audioState {
	return currentRuntime().currentAudio()
}
