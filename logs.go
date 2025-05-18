package gosprite64

import (
	_ "embed"
	"embedded/rtos"
	"log"
	"os"
	"syscall"

	"github.com/drpaneas/n64/drivers/carts"
	_ "github.com/drpaneas/n64/machine"

	"github.com/embeddedgo/fs/termfs"
)

func init() {
	var err error
	var cart carts.Cart

	// Redirect stdout and stderr either to cart's logger
	if cart = carts.ProbeAll(); cart == nil {
		return
	}

	devConsole := termfs.NewLight("termfs", nil, cart)
	rtos.Mount(devConsole, "/dev/console")
	os.Stdout, err = os.OpenFile("/dev/console", syscall.O_WRONLY, 0)
	if err != nil {
		panic(err)
	}
	os.Stderr = os.Stdout

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}
