# Save Data

Persist game progress using EEPROM, SRAM, or FlashRAM cartridge storage.

```go
import "github.com/drpaneas/gosprite64/save"
```

## Storage types

The N64 cartridge supports several save storage technologies. Each has different capacity and access characteristics:

| Type | Constructor | Capacity | Notable games |
|---|---|---|---|
| EEPROM 4Kbit | `save.NewEEPROM4K()` | 512 bytes | Super Mario 64, Mario Kart 64 |
| EEPROM 16Kbit | `save.NewEEPROM16K()` | 2,048 bytes | Yoshi's Story |
| SRAM 256Kbit | `save.NewSRAM()` | 32,768 bytes | Many third-party games |
| FlashRAM 1Mbit | `save.NewFlashRAM()` | 131,072 bytes | Paper Mario, Pokemon Stadium |

Choose the smallest type that fits your save data. EEPROM 4K is enough for most games that only store settings and a few save slots. FlashRAM gives 128 KB for games with large worlds or replay data.

## The Storage interface

All storage types implement the `Storage` interface:

```go
type Storage interface {
    Type() StorageType  // which backend (StorageEEPROM4K, StorageSRAM, etc.)
    Read(addr int, buf []byte) error
    Write(addr int, data []byte) error
    Size() int          // total capacity in bytes
}
```

`Read` fills `buf` with data starting at byte address `addr`. `Write` writes `data` starting at `addr`. Both return `save.ErrOutOfRange` if the operation would exceed the storage capacity.

## Creating storage

```go
// Pick the type that matches your cartridge configuration
storage := save.NewEEPROM4K()    // 512 bytes
storage := save.NewEEPROM16K()   // 2 KB
storage := save.NewSRAM()        // 32 KB
storage := save.NewFlashRAM()    // 128 KB
```

## Reading and writing

### Low-level: Read and Write

Read and write at specific byte addresses:

```go
// Write 8 bytes at address 0
err := storage.Write(0, []byte{1, 2, 3, 4, 5, 6, 7, 8})

// Read them back
buf := make([]byte, 8)
err := storage.Read(0, buf)
```

### High-level: ReadAll and WriteAll

`ReadAll` reads the entire storage into a new byte slice. `WriteAll` writes a byte slice to the beginning of storage:

```go
// Read the full contents
data, err := save.ReadAll(storage)

// Write new data (must not exceed capacity)
err := save.WriteAll(storage, data)
```

`WriteAll` returns `save.ErrOutOfRange` if the data is larger than the storage capacity.

## Checksum

`Checksum` computes a simple additive checksum over a byte slice. Use it to detect corrupted save data:

```go
func Checksum(data []byte) uint32
```

```go
sum := save.Checksum(saveData)
```

The checksum is a straightforward sum of all bytes cast to `uint32`. It catches accidental corruption (bit flips, incomplete writes) but is not cryptographically secure.

## Errors

| Error | Meaning |
|---|---|
| `save.ErrNotAvailable` | Storage backend not initialized (no read/write function set) |
| `save.ErrOutOfRange` | Address + length exceeds storage capacity |
| `save.ErrReadFailed` | Hardware read error |
| `save.ErrWriteFailed` | Hardware write error |

## Complete save/load example

```go
type SaveData struct {
    Level    uint8
    Score    uint32
    Lives    uint8
    Checksum uint32
}

func encodeSave(s *SaveData) []byte {
    buf := make([]byte, 7)
    buf[0] = s.Level
    buf[1] = byte(s.Score >> 24)
    buf[2] = byte(s.Score >> 16)
    buf[3] = byte(s.Score >> 8)
    buf[4] = byte(s.Score)
    buf[5] = s.Lives
    // Compute checksum over the payload bytes
    sum := save.Checksum(buf[:6])
    buf[6] = byte(sum)
    return buf
}

func decodeSave(data []byte) (*SaveData, bool) {
    if len(data) < 7 {
        return nil, false
    }
    s := &SaveData{
        Level: data[0],
        Score: uint32(data[1])<<24 | uint32(data[2])<<16 |
               uint32(data[3])<<8 | uint32(data[4]),
        Lives: data[5],
    }
    expected := byte(save.Checksum(data[:6]))
    if data[6] != expected {
        return nil, false // corrupted
    }
    return s, true
}

func saveGame(storage save.Storage, s *SaveData) error {
    return save.WriteAll(storage, encodeSave(s))
}

func loadGame(storage save.Storage) (*SaveData, error) {
    data, err := save.ReadAll(storage)
    if err != nil {
        return nil, err
    }
    s, ok := decodeSave(data)
    if !ok {
        return nil, fmt.Errorf("save data corrupted")
    }
    return s, nil
}
```

## EEPROM details

EEPROM is accessed in 8-byte blocks through the N64's serial interface (SI/PIF). The `EEPROM` type handles block alignment internally. EEPROM is the most common save type for first-party Nintendo 64 games.

## SRAM details

SRAM is battery-backed static RAM accessed via PI DMA at address `0x08000000`. It offers fast random access but requires a battery in the cartridge to retain data when powered off.

## FlashRAM details

FlashRAM uses a command-based protocol via the PI interface. Write operations require erasing 16 KB sectors before writing new data. The `FlashRAM` type handles the erase-before-write protocol internally. FlashRAM does not require a battery - data persists without power.

## Reference Example

See `examples/save_demo` in the GoSprite64 repository for a minimal working save/load example using SRAM.
