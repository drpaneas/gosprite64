//go:build n64

package rspq

import (
	"github.com/clktmr/n64/rcp/cpu"
	"github.com/clktmr/n64/rcp/rsp"
)

// SubmitTask loads an OSTask onto the RSP and starts execution.
// This replicates the osSpTaskLoad + osSpTaskStartGo protocol:
//  1. Marshal the task to a 64-byte buffer
//  2. Clear signal0/1/2, set interrupt-on-break
//  3. DMA task descriptor to DMEM at 0xFC0
//  4. DMA boot microcode to IMEM
//  5. Clear halt to start execution
//
// bootCode is the rspboot microcode bytes.
// All address fields in task must be physical RDRAM addresses.
func SubmitTask(task *OSTask, bootCode []byte) {
	taskBuf := task.Marshal()
	taskSlice := cpu.MakePaddedSlice[byte](OSTaskSize)
	copy(taskSlice, taskBuf[:])
	cpu.WritebackSlice(taskSlice)

	rsp.ClearSignals(0x07)
	rsp.SetInterrupt(true)

	rsp.DMEM.WriteAt(taskSlice, int64(OSTaskDMEMOff))

	boot := cpu.CopyPaddedSlice(bootCode)
	cpu.WritebackSlice(boot)
	rsp.IMEM.WriteAt(boot, 0)

	rsp.Resume()
}

// WaitTaskDone blocks until the RSP breaks (task complete).
func WaitTaskDone() {
	rsp.IntBreak.Wait(0)
}
