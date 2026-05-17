package save

// SRAM implements Storage for N64 SRAM (256Kbit / 32KB).
// SRAM is accessed via PI DMA at address 0x08000000.
type SRAM struct {
	ReadFunc  func(addr int, buf []byte) error
	WriteFunc func(addr int, data []byte) error
}

func NewSRAM() *SRAM {
	return &SRAM{}
}

func (s *SRAM) Type() StorageType { return StorageSRAM }
func (s *SRAM) Size() int         { return StorageSRAM.Size() }

func (s *SRAM) Read(addr int, buf []byte) error {
	if addr < 0 || addr+len(buf) > s.Size() {
		return ErrOutOfRange
	}
	if s.ReadFunc == nil {
		return ErrNotAvailable
	}
	return s.ReadFunc(addr, buf)
}

func (s *SRAM) Write(addr int, data []byte) error {
	if addr < 0 || addr+len(data) > s.Size() {
		return ErrOutOfRange
	}
	if s.WriteFunc == nil {
		return ErrNotAvailable
	}
	return s.WriteFunc(addr, data)
}
