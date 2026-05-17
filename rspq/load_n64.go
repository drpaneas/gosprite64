//go:build n64

package rspq

import (
	"github.com/clktmr/n64/rcp/cpu"
	"github.com/clktmr/n64/rcp/rsp"
	"github.com/clktmr/n64/rcp/rsp/ucode"
)

func Load(m *Microcode) {
	if m == nil {
		return
	}
	u := ucode.NewUCode(m.Type.String(), cpu.Addr(0), m.Code, m.Data)
	rsp.Load(u)
}

func Start() {
	rsp.Resume()
}

func IsStopped() bool {
	return rsp.Stopped()
}

func WaitDone() {
	rsp.IntBreak.Sleep()
}
