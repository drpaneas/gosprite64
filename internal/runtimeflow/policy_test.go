package runtimeflow

import (
	"slices"
	"testing"
)

func TestBootstrapOrder(t *testing.T) {
	want := []Stage{
		StageConsole,
		StageVideo,
		StageGameInit,
		StageAudio,
		StageLoop,
	}

	if got := BootstrapOrder(); !slices.Equal(got, want) {
		t.Fatalf("BootstrapOrder() = %v, want %v", got, want)
	}
}

func TestStatusGuards(t *testing.T) {
	tests := []struct {
		name          string
		status        Status
		wantCanDraw   bool
		wantCanQueue  bool
	}{
		{
			name:         "zero status blocks draw and audio",
			status:       Status{},
			wantCanDraw:  false,
			wantCanQueue: false,
		},
		{
			name:         "video ready allows draw only",
			status:       Status{VideoReady: true},
			wantCanDraw:  true,
			wantCanQueue: false,
		},
		{
			name:         "audio ready allows queue only",
			status:       Status{AudioReady: true},
			wantCanDraw:  false,
			wantCanQueue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.CanDraw(); got != tt.wantCanDraw {
				t.Fatalf("Status.CanDraw() = %v, want %v", got, tt.wantCanDraw)
			}

			if got := tt.status.CanQueueAudio(); got != tt.wantCanQueue {
				t.Fatalf("Status.CanQueueAudio() = %v, want %v", got, tt.wantCanQueue)
			}
		})
	}
}
