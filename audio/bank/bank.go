package bank

import (
	"encoding/binary"
	"errors"
	"math"
)

// Instrument represents a single instrument in a sound bank.
type Instrument struct {
	ID         uint8
	Volume     uint8
	Pan        uint8
	Priority   uint8
	SampleRate uint16
	KeyLow     uint8
	KeyHigh    uint8
	Sounds     []Sound
}

// Sound is a single sample within an instrument, with key range and tuning.
type Sound struct {
	SampleAddr uint32
	SampleLen  uint32
	LoopStart  uint32
	LoopEnd    uint32
	LoopCount  int32
	KeyBase    uint8
	Tuning     float32
}

// Bank is a collection of instruments loaded from a sound bank file.
type Bank struct {
	Instruments []Instrument
	SampleData  []byte
}

// NewBank creates an empty sound bank.
func NewBank() *Bank {
	return &Bank{}
}

// GetInstrument returns the instrument at the given index, or nil.
func (b *Bank) GetInstrument(idx uint8) *Instrument {
	if int(idx) < len(b.Instruments) {
		return &b.Instruments[idx]
	}
	return nil
}

// InstrumentCount returns the number of instruments in the bank.
func (b *Bank) InstrumentCount() int {
	return len(b.Instruments)
}

var ErrInvalidBank = errors.New("bank: invalid bank data")

// LoadBank parses a simple sound bank binary format.
// Format: [u8 instrumentCount] then per instrument:
//
//	[u8 id] [u8 volume] [u8 pan] [u8 priority]
//	[u16be sampleRate] [u8 keyLow] [u8 keyHigh]
//	[u8 soundCount] then per sound:
//	  [u32be sampleAddr] [u32be sampleLen]
//	  [u32be loopStart] [u32be loopEnd] [s32be loopCount]
//	  [u8 keyBase] [3 bytes padding]
//	  [f32be tuning]
func LoadBank(data []byte) (*Bank, error) {
	if len(data) < 1 {
		return nil, ErrInvalidBank
	}
	b := &Bank{}
	pos := 0
	instCount := int(data[pos])
	pos++

	for i := 0; i < instCount; i++ {
		if pos+9 > len(data) {
			return nil, ErrInvalidBank
		}
		inst := Instrument{
			ID:         data[pos],
			Volume:     data[pos+1],
			Pan:        data[pos+2],
			Priority:   data[pos+3],
			SampleRate: binary.BigEndian.Uint16(data[pos+4 : pos+6]),
			KeyLow:     data[pos+6],
			KeyHigh:    data[pos+7],
		}
		soundCount := int(data[pos+8])
		pos += 9

		for j := 0; j < soundCount; j++ {
			if pos+28 > len(data) {
				return nil, ErrInvalidBank
			}
			snd := Sound{
				SampleAddr: binary.BigEndian.Uint32(data[pos : pos+4]),
				SampleLen:  binary.BigEndian.Uint32(data[pos+4 : pos+8]),
				LoopStart:  binary.BigEndian.Uint32(data[pos+8 : pos+12]),
				LoopEnd:    binary.BigEndian.Uint32(data[pos+12 : pos+16]),
				LoopCount:  int32(binary.BigEndian.Uint32(data[pos+16 : pos+20])),
				KeyBase:    data[pos+20],
				Tuning:     math.Float32frombits(binary.BigEndian.Uint32(data[pos+24 : pos+28])),
			}
			pos += 28
			inst.Sounds = append(inst.Sounds, snd)
		}
		b.Instruments = append(b.Instruments, inst)
	}
	return b, nil
}
