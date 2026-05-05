package gosprite64

type runtimeState struct {
	video *videoState
}

var activeRuntime *runtimeState

func newRuntimeState() *runtimeState {
	return &runtimeState{}
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
