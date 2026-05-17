//go:build !n64

package gfx

// Execute is a no-op on non-N64 targets.
func Execute(dl *DisplayList) {}

// Flush is a no-op on non-N64 targets.
func Flush() {}

// ExecuteViaRSP is a no-op on non-N64 builds.
func ExecuteViaRSP(dl *DisplayList, bootCode, ucodeText, ucodeData []byte) {}
