package save

import "errors"

var (
	ErrNotAvailable = errors.New("save: storage type not available")
	ErrOutOfRange   = errors.New("save: address out of range")
	ErrReadFailed   = errors.New("save: read failed")
	ErrWriteFailed  = errors.New("save: write failed")
)

// StorageType identifies the save backend.
type StorageType int

const (
	StorageNone    StorageType = iota
	StorageEEPROM4K            // 512 bytes, used by SM64, Mario Kart 64
	StorageEEPROM16K           // 2048 bytes, used by Yoshi's Story
	StorageSRAM                // 32768 bytes (256Kbit), used by many games
	StorageFlashRAM            // 131072 bytes (1Mbit), used by Paper Mario, Pokemon Stadium
)

// Size returns the total byte capacity for the storage type.
func (s StorageType) Size() int {
	switch s {
	case StorageEEPROM4K:
		return 512
	case StorageEEPROM16K:
		return 2048
	case StorageSRAM:
		return 32768
	case StorageFlashRAM:
		return 131072
	default:
		return 0
	}
}

func (s StorageType) String() string {
	switch s {
	case StorageEEPROM4K:
		return "EEPROM 4Kbit"
	case StorageEEPROM16K:
		return "EEPROM 16Kbit"
	case StorageSRAM:
		return "SRAM 256Kbit"
	case StorageFlashRAM:
		return "FlashRAM 1Mbit"
	default:
		return "None"
	}
}

// Storage is the interface for N64 save backends. Implementations provide
// access to EEPROM, SRAM, or FlashRAM through a uniform API.
type Storage interface {
	Type() StorageType
	Read(addr int, buf []byte) error
	Write(addr int, data []byte) error
	Size() int
}

// ReadAll reads the entire save storage into a new byte slice.
func ReadAll(s Storage) ([]byte, error) {
	buf := make([]byte, s.Size())
	if err := s.Read(0, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// WriteAll writes data to the beginning of storage. Data must not exceed
// the storage capacity.
func WriteAll(s Storage, data []byte) error {
	if len(data) > s.Size() {
		return ErrOutOfRange
	}
	return s.Write(0, data)
}

// Checksum computes a simple additive checksum over a byte slice,
// suitable for save file integrity checks.
func Checksum(data []byte) uint32 {
	var sum uint32
	for _, b := range data {
		sum += uint32(b)
	}
	return sum
}
