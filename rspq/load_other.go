//go:build !n64

package rspq

func Load(m *Microcode) {}
func Start()            {}
func IsStopped() bool   { return true }
func WaitDone()         {}
