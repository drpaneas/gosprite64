package gosprite64

import (
	"io"

	"github.com/clktmr/n64/drivers/cartfs"
	tileloader "github.com/drpaneas/gosprite64/internal/tile2d/loader"
)

// cartFS is the global cartridge filesystem instance
var _cartFS cartfs.FS // Prefixed with underscore to indicate intended usage in the future

func RegisterAssetFS(f cartfs.FS) {
	_cartFS = f
}

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

var _ tileloader.Loader = cartLoader{}

type cartLoader struct{}

func (cartLoader) ReadAsset(path string) ([]byte, error) {
	return LoadFromCartridge(path)
}
