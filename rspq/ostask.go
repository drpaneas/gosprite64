package rspq

import "encoding/binary"

// OSTask is the 64-byte task descriptor DMA'd to DMEM at offset 0xFC0.
// All addresses must be physical (KSEG0/KSEG1 stripped).
type OSTask struct {
	Type          uint32
	Flags         uint32
	UcodeBoot     uint32
	UcodeBootSize uint32
	Ucode         uint32
	UcodeSize     uint32
	UcodeData     uint32
	UcodeDataSize uint32
	DRAMStack     uint32
	DRAMStackSize uint32
	OutputBuff    uint32
	OutputBuffEnd uint32
	DataPtr       uint32
	DataSize      uint32
	YieldDataPtr  uint32
	YieldDataSize uint32
}

const (
	TaskGfx   = 1
	TaskAudio = 2

	FlagDPWait  = 0x0002
	FlagYielded = 0x0001

	OSTaskSize    = 64
	OSTaskDMEMOff = 0xFC0
)

// Marshal serializes the OSTask to a 64-byte big-endian buffer.
func (t *OSTask) Marshal() [OSTaskSize]byte {
	var buf [OSTaskSize]byte
	binary.BigEndian.PutUint32(buf[0x00:], t.Type)
	binary.BigEndian.PutUint32(buf[0x04:], t.Flags)
	binary.BigEndian.PutUint32(buf[0x08:], t.UcodeBoot)
	binary.BigEndian.PutUint32(buf[0x0C:], t.UcodeBootSize)
	binary.BigEndian.PutUint32(buf[0x10:], t.Ucode)
	binary.BigEndian.PutUint32(buf[0x14:], t.UcodeSize)
	binary.BigEndian.PutUint32(buf[0x18:], t.UcodeData)
	binary.BigEndian.PutUint32(buf[0x1C:], t.UcodeDataSize)
	binary.BigEndian.PutUint32(buf[0x20:], t.DRAMStack)
	binary.BigEndian.PutUint32(buf[0x24:], t.DRAMStackSize)
	binary.BigEndian.PutUint32(buf[0x28:], t.OutputBuff)
	binary.BigEndian.PutUint32(buf[0x2C:], t.OutputBuffEnd)
	binary.BigEndian.PutUint32(buf[0x30:], t.DataPtr)
	binary.BigEndian.PutUint32(buf[0x34:], t.DataSize)
	binary.BigEndian.PutUint32(buf[0x38:], t.YieldDataPtr)
	binary.BigEndian.PutUint32(buf[0x3C:], t.YieldDataSize)
	return buf
}
