package save

// FlashRAM implements Storage for N64 FlashRAM (1Mbit / 128KB).
// FlashRAM uses a command-based protocol via PI at 0x08000000.
// Write operations require erasing sectors before writing.
type FlashRAM struct {
	ReadFunc  func(addr int, buf []byte) error
	WriteFunc func(addr int, data []byte) error
}

const flashSectorSize = 16384 // 16KB sectors, 8 sectors total

func NewFlashRAM() *FlashRAM {
	return &FlashRAM{}
}

func (f *FlashRAM) Type() StorageType { return StorageFlashRAM }
func (f *FlashRAM) Size() int         { return StorageFlashRAM.Size() }

func (f *FlashRAM) Read(addr int, buf []byte) error {
	if addr < 0 || addr+len(buf) > f.Size() {
		return ErrOutOfRange
	}
	if f.ReadFunc == nil {
		return ErrNotAvailable
	}
	return f.ReadFunc(addr, buf)
}

func (f *FlashRAM) Write(addr int, data []byte) error {
	if addr < 0 || addr+len(data) > f.Size() {
		return ErrOutOfRange
	}
	if f.WriteFunc == nil {
		return ErrNotAvailable
	}
	return f.WriteFunc(addr, data)
}
