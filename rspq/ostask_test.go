package rspq

import (
	"encoding/binary"
	"testing"
)

func TestOSTaskMarshalSize(t *testing.T) {
	task := OSTask{Type: TaskGfx}
	buf := task.Marshal()
	if len(buf) != OSTaskSize {
		t.Fatalf("expected %d bytes, got %d", OSTaskSize, len(buf))
	}
}

func TestOSTaskMarshalFields(t *testing.T) {
	task := OSTask{
		Type:          TaskGfx,
		Flags:         FlagDPWait,
		DataPtr:       0x80200000,
		DataSize:      1024,
		UcodeBootSize: 208,
	}
	buf := task.Marshal()
	if binary.BigEndian.Uint32(buf[0x00:]) != TaskGfx {
		t.Fatal("type mismatch")
	}
	if binary.BigEndian.Uint32(buf[0x04:]) != FlagDPWait {
		t.Fatal("flags mismatch")
	}
	if binary.BigEndian.Uint32(buf[0x30:]) != 0x80200000 {
		t.Fatal("data_ptr mismatch")
	}
	if binary.BigEndian.Uint32(buf[0x34:]) != 1024 {
		t.Fatal("data_size mismatch")
	}
}

func TestOSTaskDMEMOffset(t *testing.T) {
	if OSTaskDMEMOff != 0x1000-OSTaskSize {
		t.Fatalf("DMEM offset should be IMEM_START - sizeof(OSTask), got 0x%X", OSTaskDMEMOff)
	}
}
