package gosprite64

import (
	"io"

	"github.com/clktmr/n64/drivers/cartfs"
	tileloader "github.com/drpaneas/gosprite64/internal/tile2d/loader"
)

var _cartFS cartfs.FS

func RegisterAssetFS(f cartfs.FS) {
	_cartFS = f
}

// LoadFromCartridge reads a file from the cartridge filesystem.
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
