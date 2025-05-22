package gosprite64

import (
	"io/ioutil"

	"github.com/clktmr/n64/drivers/cartfs"
)

// cartFS is the global cartridge filesystem instance
var cartFS cartfs.FS

// loadFromCartridge reads a file from the cartridge filesystem
// filename: name of the file to read from the cartridge
// Returns the file contents as a byte slice or an error if the operation fails
func loadFromCartridge(filename string) ([]byte, error) {
	f, err := cartFS.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}
