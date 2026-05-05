package runtimeflow

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
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
		name         string
		status       Status
		wantCanDraw  bool
		wantCanQueue bool
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

func TestRuntimeBootstrapIsExplicit(t *testing.T) {
	logsSource := mustReadRepoFile(t, "logs.go")
	if strings.Contains(logsSource, "func init()") {
		t.Fatal("logs.go must not hide bootstrap in package init")
	}
	if !strings.Contains(logsSource, "func setupConsole()") {
		t.Fatal("logs.go must expose setupConsole")
	}

	runtimeSource := mustReadRepoFile(t, "runtime.go")
	for _, snippet := range []string{
		"type runtimeState struct{}",
		"var activeRuntime *runtimeState",
		"func newRuntimeState() *runtimeState",
		"func activateRuntime(rt *runtimeState)",
		"func currentRuntime() *runtimeState",
	} {
		if !strings.Contains(runtimeSource, snippet) {
			t.Fatalf("runtime.go must contain %q", snippet)
		}
	}

	gameLoopSource := mustReadRepoFile(t, "gameloop.go")
	assertOrderedSubstrings(t, gameLoopSource,
		"setupConsole()",
		"rt := newRuntimeState()",
		"activateRuntime(rt)",
		"videoInit()",
		"rdp.RDP.SetScissor(",
		"g.Init()",
		"initAudioV1()",
	)
}

func mustReadRepoFile(t *testing.T, name string) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	path := filepath.Join(filepath.Dir(file), "..", "..", name)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}

	return string(content)
}

func assertOrderedSubstrings(t *testing.T, content string, snippets ...string) {
	t.Helper()

	last := -1
	for _, snippet := range snippets {
		index := strings.Index(content, snippet)
		if index == -1 {
			t.Fatalf("missing %q", snippet)
		}
		if index <= last {
			t.Fatalf("%q must appear after previous bootstrap step", snippet)
		}
		last = index
	}
}
