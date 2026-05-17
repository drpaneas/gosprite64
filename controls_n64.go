//go:build n64

package gosprite64

import (
	"github.com/clktmr/n64/rcp/serial"
	"github.com/clktmr/n64/rcp/serial/joybus"
)

func rumbleWrite(port int, enabled bool) {
	block := serial.NewCommandBlock(0)
	for i := 0; i < port; i++ {
		joybus.ControlByte(block, joybus.CtrlSkip)
	}
	cmd, err := joybus.NewWritePakCommand(block)
	if err != nil {
		return
	}
	cmd.SetAddress(0xC000)
	data := make([]byte, 32)
	if enabled {
		for i := range data {
			data[i] = 0x01
		}
	}
	cmd.SetData(data)
	serial.Run(block)
}
