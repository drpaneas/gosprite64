package gosprite64

type runtimeState struct{}

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
