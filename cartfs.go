package gosprite64

import (
	"io"

	"github.com/clktmr/n64/drivers/cartfs"
)

// cartFS is the global cartridge filesystem instance
var _cartFS cartfs.FS // Prefixed with underscore to indicate intended usage in the future

// loadFromCartridge reads a file from the cartridge filesystem
// filename: name of the file to read from the cartridge
// Returns the file contents as a byte slice or an error if the operation fails
// LoadFromCartridge reads a file from the cartridge filesystem
// filename: name of the file to read from the cartridge
// Returns the file contents as a byte slice or an error if the operation fails
func LoadFromCartridge(filename string) ([]byte, error) {
	f, err := _cartFS.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := f.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	return io.ReadAll(f)
}
