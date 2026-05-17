package save

// EEPROM implements Storage for N64 EEPROM (4Kbit or 16Kbit).
// EEPROM is accessed in 8-byte blocks via the SI/PIF interface.
//
// On real hardware, this will use osEepromRead/osEepromWrite equivalents
// from the clktmr/n64 package. This implementation provides the structure
// and validation; hardware-specific read/write is pluggable via ReadFunc
// and WriteFunc.
type EEPROM struct {
	kind      StorageType
	ReadFunc  func(addr int, buf []byte) error
	WriteFunc func(addr int, data []byte) error
}

const eepromBlockSize = 8

// NewEEPROM4K creates an EEPROM storage for 4Kbit (512 bytes).
func NewEEPROM4K() *EEPROM {
	return &EEPROM{kind: StorageEEPROM4K}
}

// NewEEPROM16K creates an EEPROM storage for 16Kbit (2048 bytes).
func NewEEPROM16K() *EEPROM {
	return &EEPROM{kind: StorageEEPROM16K}
}

func (e *EEPROM) Type() StorageType { return e.kind }
func (e *EEPROM) Size() int         { return e.kind.Size() }

func (e *EEPROM) Read(addr int, buf []byte) error {
	if addr < 0 || addr+len(buf) > e.Size() {
		return ErrOutOfRange
	}
	if e.ReadFunc == nil {
		return ErrNotAvailable
	}
	return e.ReadFunc(addr, buf)
}

func (e *EEPROM) Write(addr int, data []byte) error {
	if addr < 0 || addr+len(data) > e.Size() {
		return ErrOutOfRange
	}
	if e.WriteFunc == nil {
		return ErrNotAvailable
	}
	return e.WriteFunc(addr, data)
}
