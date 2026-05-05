package runtimeflow

type Stage uint8

const (
	StageConsole Stage = iota
	StageVideo
	StageGameInit
	StageAudio
	StageLoop
)

func BootstrapOrder() []Stage {
	return []Stage{
		StageConsole,
		StageVideo,
		StageGameInit,
		StageAudio,
		StageLoop,
	}
}

type Status struct {
	VideoReady bool
	AudioReady bool
}

func (s Status) CanDraw() bool {
	return s.VideoReady
}

func (s Status) CanQueueAudio() bool {
	return s.AudioReady
}
