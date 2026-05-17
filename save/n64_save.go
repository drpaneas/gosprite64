//go:build n64

package save

import (
	"github.com/clktmr/n64/rcp/cpu"
	"github.com/clktmr/n64/rcp/periph"
)

const sramPIAddr = cpu.Addr(0x08000000)

func NewN64SRAM() *SRAM {
	dev := periph.NewDevice(sramPIAddr, uint32(StorageSRAM.Size()))
	s := NewSRAM()
	s.ReadFunc = func(addr int, buf []byte) error {
		padded := cpu.MakePaddedSlice[byte](len(buf))
		_, err := dev.ReadAt(padded, int64(addr))
		if err != nil {
			return err
		}
		copy(buf, padded)
		return nil
	}
	s.WriteFunc = func(addr int, data []byte) error {
		padded := cpu.CopyPaddedSlice(data)
		cpu.WritebackSlice(padded)
		_, err := dev.WriteAt(padded, int64(addr))
		return err
	}
	return s
}

func NewN64FlashRAM() *FlashRAM {
	dev := periph.NewDevice(sramPIAddr, uint32(StorageFlashRAM.Size()))
	f := NewFlashRAM()
	f.ReadFunc = func(addr int, buf []byte) error {
		padded := cpu.MakePaddedSlice[byte](len(buf))
		_, err := dev.ReadAt(padded, int64(addr))
		if err != nil {
			return err
		}
		copy(buf, padded)
		return nil
	}
	f.WriteFunc = func(addr int, data []byte) error {
		padded := cpu.CopyPaddedSlice(data)
		cpu.WritebackSlice(padded)
		_, err := dev.WriteAt(padded, int64(addr))
		return err
	}
	return f
}
