//go:build n64

package gfx

import (
	"encoding/binary"
	"image"
	"image/color"
	_ "unsafe"

	"github.com/clktmr/n64/rcp/cpu"
	"github.com/clktmr/n64/rcp/rdp"
	"github.com/drpaneas/gosprite64/rspq"
)

type rawRDPCommand uint64

//go:linkname rdpPush github.com/clktmr/n64/rcp/rdp.(*DisplayList).Push
func rdpPush(dl *rdp.DisplayList, cmds ...rawRDPCommand)

func pushRawRDP(words []uint64) {
	if len(words) == 0 {
		return
	}
	cmds := make([]rawRDPCommand, len(words))
	for i, w := range words {
		cmds[i] = rawRDPCommand(w)
	}
	rdpPush(&rdp.RDP, cmds...)
}

// PushRaw submits one or more fully packed raw RDP command words.
func PushRaw(words ...uint64) {
	pushRawRDP(words)
}

// Execute walks the display list and translates DP opcodes into rdp.RDP calls.
// SP commands (vertex, triangle, matrix) are skipped - this is an HLE path
// that only handles the RDP-direct subset.
func Execute(dl *DisplayList) {
	for _, cmd := range dl.Commands() {
		if len(cmd.Raw) > 0 {
			pushRawRDP(cmd.Raw)
			continue
		}

		opcode := cmd.W0 >> 24

		switch opcode {
		case OpDPPipeSync:
			rdp.RDP.Push(rdp.SyncPipe)
		case OpDPTileSync:
			rdp.RDP.Push(rdp.SyncTile)
		case OpDPLoadSync:
			rdp.RDP.Push(rdp.SyncLoad)
		case OpDPFullSync:
			rdp.RDP.Push(rdp.SyncFull)

		case OpDPSetFillColor:
			// W1 contains two packed RGBA5551 pixels; decode the first one.
			rgba16 := uint16(cmd.W1 >> 16)
			r := uint8((rgba16 >> 11) & 0x1F)
			g := uint8((rgba16 >> 6) & 0x1F)
			b := uint8((rgba16 >> 1) & 0x1F)
			a := uint8(rgba16 & 1)
			rdp.RDP.SetFillColor(color.RGBA{
				R: (r << 3) | (r >> 2),
				G: (g << 3) | (g >> 2),
				B: (b << 3) | (b >> 2),
				A: a * 0xFF,
			})

		case OpDPFillRect:
			// W0: [31:24] opcode, [23:12] lrx, [11:0] lry (10.2 fixed)
			// W1: [23:12] ulx, [11:0] uly (10.2 fixed)
			lrx := int((cmd.W0 >> 12) & 0xFFF)
			lry := int(cmd.W0 & 0xFFF)
			ulx := int((cmd.W1 >> 12) & 0xFFF)
			uly := int(cmd.W1 & 0xFFF)
			rdp.RDP.FillRectangle(image.Rect(ulx>>2, uly>>2, lrx>>2, lry>>2))

		case OpDPSetScissor:
			// W0: [23:12] ulx<<2, [11:0] uly<<2
			// W1: [25:24] mode, [23:12] lrx<<2, [11:0] lry<<2
			ulx := int((cmd.W0 >> 12) & 0xFFF)
			uly := int(cmd.W0 & 0xFFF)
			lrx := int((cmd.W1 >> 12) & 0xFFF)
			lry := int(cmd.W1 & 0xFFF)
			rdp.RDP.SetScissor(image.Rect(ulx>>2, uly>>2, lrx>>2, lry>>2), rdp.InterlaceNone)

		case OpDPSetEnvColor:
			r := uint8(cmd.W1 >> 24)
			g := uint8(cmd.W1 >> 16)
			b := uint8(cmd.W1 >> 8)
			a := uint8(cmd.W1)
			rdp.RDP.SetEnvironmentColor(color.RGBA{R: r, G: g, B: b, A: a})

		case OpDPSetPrimColor:
			r := uint8(cmd.W1 >> 24)
			g := uint8(cmd.W1 >> 16)
			b := uint8(cmd.W1 >> 8)
			a := uint8(cmd.W1)
			rdp.RDP.SetPrimitiveColor(color.RGBA{R: r, G: g, B: b, A: a})

		case OpDPSetFogColor:
			r := uint8(cmd.W1 >> 24)
			g := uint8(cmd.W1 >> 16)
			b := uint8(cmd.W1 >> 8)
			a := uint8(cmd.W1)
			rdp.RDP.SetBlendColor(color.RGBA{R: r, G: g, B: b, A: a})

		case OpSPEndDisplayList:
			return
		}
	}
}

// Flush submits any buffered RDP commands to the hardware.
func Flush() {
	rdp.RDP.Flush()
}

// ExecuteViaRSP submits a display list to the RSP for processing by
// F3DEX2 microcode. This is the real 3D path.
func ExecuteViaRSP(dl *DisplayList, bootCode, ucodeText, ucodeData []byte) {
	if dl == nil || dl.Len() == 0 {
		return
	}

	cmds := dl.Commands()

	gfxBytes := cpu.MakePaddedSlice[byte](len(cmds) * 8)
	for i, cmd := range cmds {
		binary.BigEndian.PutUint32(gfxBytes[i*8:], cmd.W0)
		binary.BigEndian.PutUint32(gfxBytes[i*8+4:], cmd.W1)
	}
	cpu.WritebackSlice(gfxBytes)

	dramStack := cpu.MakePaddedSlice[byte](1024)
	outputBuf := cpu.MakePaddedSlice[byte](4096)
	yieldBuf := cpu.MakePaddedSlice[byte](0xC00)

	ucText := cpu.CopyPaddedSlice(ucodeText)
	ucData := cpu.CopyPaddedSlice(ucodeData)
	cpu.WritebackSlice(ucText)
	cpu.WritebackSlice(ucData)

	task := rspq.OSTask{
		Type:          rspq.TaskGfx,
		Flags:         0,
		UcodeBoot:     0,
		UcodeBootSize: uint32(len(bootCode)),
		Ucode:         uint32(cpu.PhysicalAddressSlice(ucText)),
		UcodeSize:     uint32(len(ucText)),
		UcodeData:     uint32(cpu.PhysicalAddressSlice(ucData)),
		UcodeDataSize: uint32(len(ucData)),
		DRAMStack:     uint32(cpu.PhysicalAddressSlice(dramStack)),
		DRAMStackSize: uint32(len(dramStack)),
		OutputBuff:    uint32(cpu.PhysicalAddressSlice(outputBuf)),
		OutputBuffEnd: uint32(cpu.PhysicalAddressSlice(outputBuf)) + uint32(len(outputBuf)),
		DataPtr:       uint32(cpu.PhysicalAddressSlice(gfxBytes)),
		DataSize:      uint32(len(gfxBytes)),
		YieldDataPtr:  uint32(cpu.PhysicalAddressSlice(yieldBuf)),
		YieldDataSize: uint32(len(yieldBuf)),
	}

	rspq.SubmitTask(&task, bootCode)
	rspq.WaitTaskDone()
}
