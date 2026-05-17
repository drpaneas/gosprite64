//go:build !n64

package gfx

// Execute is a no-op on non-N64 targets.
func Execute(dl *DisplayList) {}

// Flush is a no-op on non-N64 targets.
func Flush() {}
