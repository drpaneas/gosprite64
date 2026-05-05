package apinames

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCoreGameplayNaming(t *testing.T) {
	gameLoop := mustReadRepoFile(t, "gameloop.go")

	requireContains(t, gameLoop, "type Game interface {")
	requireContains(t, gameLoop, "func Run(g Game)")
	requireNotContains(t, gameLoop, "type Gamelooper interface")
	requireNotContains(t, gameLoop, "func Run(g Gamelooper)")
}

func TestExamplesUseQualifiedImports(t *testing.T) {
	clearscreen := mustReadRepoFile(t, "examples/clearscreen/main.go")
	requireContains(t, clearscreen, `"github.com/drpaneas/gosprite64"`)
	requireContains(t, clearscreen, "gosprite64.Run(&Game{})")
	requireNotContains(t, clearscreen, `import . "github.com/drpaneas/gosprite64"`)

	calibration := mustReadRepoFile(t, "examples/calibration/main.go")
	requireContains(t, calibration, `"github.com/drpaneas/gosprite64"`)
	requireContains(t, calibration, "gosprite64.Run(&Game{})")
	requireNotContains(t, calibration, `import . "github.com/drpaneas/gosprite64"`)
}

func TestHelloWorldGuideUsesQualifiedImports(t *testing.T) {
	guide := mustReadRepoFile(t, "docs/hello_world.md")
	requireContains(t, guide, `import "github.com/drpaneas/gosprite64"`)
	requireContains(t, guide, "gosprite64.Run(&Game{})")
	requireNotContains(t, guide, `import . "github.com/drpaneas/gosprite64"`)
	requireNotContains(t, guide, "Gamelooper")
}

func TestAudioGuideUsesQualifiedImports(t *testing.T) {
	guide := mustReadRepoFile(t, "docs/audio.md")
	requireNotContains(t, guide, `. "github.com/drpaneas/gosprite64"`)
	requireContains(t, guide, "gosprite64.PlaySoundEffect(")
}

func TestInputNaming(t *testing.T) {
	controls := mustReadRepoFile(t, "controls.go")

	for _, snippet := range []string{
		"ButtonA",
		"ButtonB",
		"ButtonZ",
		"ButtonStart",
		"ButtonDPadUp",
		"ButtonDPadDown",
		"ButtonDPadLeft",
		"ButtonDPadRight",
		"func IsButtonDown(",
		"func IsButtonJustPressed(",
		"func StickPosition(",
	} {
		requireContains(t, controls, snippet)
	}

	for _, snippet := range []string{
		"BtnA",
		"BtnB",
		"Btn(",
		"Btnp(",
		"GetStick(",
		"RIGHT =",
		"LEFT  =",
		"UP    =",
		"DOWN  =",
		"X     =",
		"O     =",
	} {
		requireNotContains(t, controls, snippet)
	}
}

func TestDrawHelperNaming(t *testing.T) {
	shapes := mustReadRepoFile(t, "shapes.go")
	requireContains(t, shapes, "func DrawRect(")
	requireContains(t, shapes, "func DrawLine(")
	requireContains(t, shapes, "func FillRect(")
	requireNotContains(t, shapes, "func Line(")
	requireNotContains(t, shapes, "func Rectfill(")

	gameLoop := mustReadRepoFile(t, "gameloop.go")
	requireNotContains(t, gameLoop, "func DrawRectFill(")

	text := mustReadRepoFile(t, "text.go")
	requireContains(t, text, "func DrawText(")
	requireNotContains(t, text, "func Print(")
}

func TestDocsTeachOnlyCanonicalDrawNames(t *testing.T) {
	for _, path := range []string{
		"README.md",
		"docs/introduction.md",
		"docs/hello_world.md",
		"docs/square_pixels.md",
	} {
		content := mustReadRepoFile(t, path)
		requireNotContains(t, content, "`Rectfill`")
		requireNotContains(t, content, "`DrawRectFill`")
		requireNotContains(t, content, "`Line`")
		requireNotContains(t, content, "`Print`")
	}
}

func TestMathShorthandHelpersAreGone(t *testing.T) {
	requireMissingRepoFile(t, "math.go")

	pong := mustReadRepoFile(t, "examples/pong/main.go")
	requireContains(t, pong, "rand.IntN(2)")
	requireNotContains(t, pong, "Flr(")
	requireNotContains(t, pong, "Rnd(")

	spaceInvaders := mustReadRepoFile(t, "examples/space_invaders/main.go")
	requireNotContains(t, spaceInvaders, "Btn(")
	requireNotContains(t, spaceInvaders, "Btnp(")
	requireNotContains(t, spaceInvaders, "Rnd(")
}

func TestAudioNaming(t *testing.T) {
	audioSource := mustReadRepoFile(t, "audio.go")
	for _, snippet := range []string{
		"type AudioBundle struct {",
		"func RegisterAudioBundle(bundle AudioBundle)",
		"func PlaySoundEffect(",
		"func PlayMusic(",
		"func StopMusic()",
		"func SetSoundEffectVolume(",
		"func SetMusicVolume(",
	} {
		requireContains(t, audioSource, snippet)
	}

	for _, snippet := range []string{
		"RegisterAudioV1",
		"RegisterSFXNameResolver",
		"PlayEffect(",
		"PlayTrack(",
		"StopTrack(",
		"SetEffectVolume(",
		"SetTrackVolume(",
	} {
		requireNotContains(t, audioSource, snippet)
	}
}

func TestLegacyNamesAreGoneFromDocsAndExamples(t *testing.T) {
	legacy := []string{
		"Gamelooper",
		"Btn(",
		"Btnp(",
		"Flr(",
		"Rnd(",
		"Rectfill",
		"DrawRectFill",
		"`Line`",
		".Line(",
		"Print(",
		"RegisterAudioV1",
		"RegisterSFXNameResolver",
		"PlayEffect(",
		"PlayTrack(",
		"StopTrack(",
		"SetEffectVolume(",
		"SetTrackVolume(",
	}

	for _, path := range []string{
		"README.md",
		"docs/introduction.md",
		"docs/hello_world.md",
		"docs/square_pixels.md",
		"docs/audio.md",
		"examples/clearscreen/main.go",
		"examples/calibration/main.go",
		"examples/pong/main.go",
		"examples/space_invaders/main.go",
	} {
		content := mustReadRepoFile(t, path)
		for _, snippet := range legacy {
			requireNotContains(t, content, snippet)
		}
	}
}

func mustReadRepoFile(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(repoFilePath(t, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(data)
}

func repoFilePath(t *testing.T, name string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Join(filepath.Dir(file), "..", "..", name)
}

func requireMissingRepoFile(t *testing.T, name string) {
	t.Helper()
	if _, err := os.Stat(repoFilePath(t, name)); err == nil {
		t.Fatalf("expected %s to be missing", name)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat %s: %v", name, err)
	}
}

func requireContains(t *testing.T, content, snippet string) {
	t.Helper()
	if !strings.Contains(content, snippet) {
		t.Fatalf("missing %q", snippet)
	}
}

func requireNotContains(t *testing.T, content, snippet string) {
	t.Helper()
	if strings.Contains(content, snippet) {
		t.Fatalf("unexpected %q", snippet)
	}
}
