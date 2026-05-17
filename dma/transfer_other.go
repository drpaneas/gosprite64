//go:build !n64

package dma

import "errors"

var errNoHardware = errors.New("dma: not available on this platform")

func CartToRDRAM(romOffset uint32, dst []byte) error { return errNoHardware }
func SRAMRead(offset uint32, dst []byte) error       { return errNoHardware }
func SRAMWrite(offset uint32, src []byte) error       { return errNoHardware }
