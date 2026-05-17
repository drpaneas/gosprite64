//go:build n64

package dma

import (
	"github.com/clktmr/n64/rcp/cpu"
	"github.com/clktmr/n64/rcp/periph"
)

func CartToRDRAM(romOffset uint32, dst []byte) error {
	dev := periph.NewDevice(cpu.Addr(0x10000000+romOffset), uint32(len(dst)))
	_, err := dev.ReadAt(dst, 0)
	if err != nil {
		return err
	}
	cpu.InvalidateSlice(dst)
	return nil
}

func SRAMRead(offset uint32, dst []byte) error {
	dev := periph.NewDevice(cpu.Addr(0x08000000), 32768)
	_, err := dev.ReadAt(dst, int64(offset))
	return err
}

func SRAMWrite(offset uint32, src []byte) error {
	dev := periph.NewDevice(cpu.Addr(0x08000000), 32768)
	_, err := dev.WriteAt(src, int64(offset))
	return err
}
