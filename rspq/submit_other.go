//go:build !n64

package rspq

// SubmitTask is a no-op on non-N64 builds.
func SubmitTask(task *OSTask, bootCode []byte) {}

// WaitTaskDone is a no-op on non-N64 builds.
func WaitTaskDone() {}
