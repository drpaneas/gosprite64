package save

import "testing"

func TestStorageTypeSize(t *testing.T) {
	tests := []struct {
		st   StorageType
		size int
	}{
		{StorageEEPROM4K, 512},
		{StorageEEPROM16K, 2048},
		{StorageSRAM, 32768},
		{StorageFlashRAM, 131072},
		{StorageNone, 0},
	}
	for _, tt := range tests {
		if got := tt.st.Size(); got != tt.size {
			t.Errorf("%s.Size() = %d, want %d", tt.st, got, tt.size)
		}
	}
}

func TestEEPROMBoundsCheck(t *testing.T) {
	e := NewEEPROM4K()
	err := e.Read(500, make([]byte, 20))
	if err != ErrOutOfRange {
		t.Fatalf("expected ErrOutOfRange, got %v", err)
	}
}

func TestEEPROMNoBackend(t *testing.T) {
	e := NewEEPROM4K()
	err := e.Read(0, make([]byte, 8))
	if err != ErrNotAvailable {
		t.Fatalf("expected ErrNotAvailable, got %v", err)
	}
}

func TestEEPROMWithBackend(t *testing.T) {
	backing := make([]byte, 512)
	for i := range backing {
		backing[i] = byte(i)
	}
	e := NewEEPROM4K()
	e.ReadFunc = func(addr int, buf []byte) error {
		copy(buf, backing[addr:addr+len(buf)])
		return nil
	}
	buf := make([]byte, 8)
	if err := e.Read(0, buf); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 8; i++ {
		if buf[i] != byte(i) {
			t.Fatalf("byte %d: expected %d, got %d", i, i, buf[i])
		}
	}
}

func TestChecksum(t *testing.T) {
	data := []byte{1, 2, 3, 4}
	if got := Checksum(data); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}
