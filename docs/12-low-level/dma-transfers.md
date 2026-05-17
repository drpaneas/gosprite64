# DMA Transfers

The `dma` package provides low-level data transfer functions for the N64. The N64's CPU cannot directly access cartridge ROM - all data must be copied into RDRAM via DMA (Direct Memory Access). This package also covers SRAM access for save data, MIO0 decompression, memory pools, and segment address translation.

## Cartridge ROM to RDRAM

### CartToRDRAM

Copies data from cartridge ROM into an RDRAM buffer:

```go
func CartToRDRAM(romOffset uint32, dst []byte) error
```

- `romOffset` - byte offset into the cartridge ROM (relative to 0x10000000)
- `dst` - destination slice in RDRAM, determines the transfer size

```go
buf := make([]byte, 4096)
err := dma.CartToRDRAM(0x100000, buf)
if err != nil {
    panic(err)
}
```

Internally, this creates a peripheral device at the ROM address and performs a read via the N64's PI (Parallel Interface). After the transfer, it invalidates the CPU data cache for the destination region to ensure coherency.

This is an N64-only function (build tag `n64`). On other platforms, use the standard file I/O or embedded filesystem instead.

## SRAM access

SRAM is the N64's battery-backed save memory, mapped at physical address 0x08000000. The standard SRAM size is 32KB.

### SRAMRead

Reads data from SRAM into a buffer:

```go
func SRAMRead(offset uint32, dst []byte) error
```

```go
saveData := make([]byte, 256)
err := dma.SRAMRead(0, saveData)
```

### SRAMWrite

Writes data from a buffer to SRAM:

```go
func SRAMWrite(offset uint32, src []byte) error
```

```go
err := dma.SRAMWrite(0, saveData)
```

Both functions access SRAM through the PI peripheral interface. The offset is relative to the start of SRAM (0x08000000). These are N64-only functions.

## MIO0 decompression

### DecompressMIO0

Decompresses MIO0-compressed data, the standard compression format used in many N64 games (including Super Mario 64):

```go
func DecompressMIO0(src []byte) ([]byte, error)
```

```go
compressed := make([]byte, compressedSize)
dma.CartToRDRAM(romOffset, compressed)

decompressed, err := dma.DecompressMIO0(compressed)
if err != nil {
    panic(err)
}
```

The MIO0 format structure:

| Offset | Size | Description |
|--------|------|-------------|
| 0 | 4 bytes | Magic: "MIO0" |
| 4 | 4 bytes | Decompressed size (big-endian) |
| 8 | 4 bytes | Offset to compressed data |
| 12 | 4 bytes | Offset to uncompressed data |
| 16+ | variable | Layout bits, compressed data, uncompressed data |

The algorithm uses layout bits to decide whether each output byte is a literal copy (from the uncompressed stream) or a back-reference (length + offset pair from the compressed stream). Back-references copy 3-18 bytes from previously decompressed output.

Returns `ErrInvalidMIO0` if the magic bytes are wrong or the data is truncated.

## Memory pools

### Pool

A linear memory allocator for RDRAM that allocates from both ends:

```go
type Pool struct { /* ... */ }
```

See [Memory Pools](./memory-pools.md) for the full `Pool` API including `NewPool`, `AllocHead`, `AllocTail`, `ResetTail`, and `Available`.

## Segment table

### SegmentTable

The N64 uses 16 segment registers (0x00-0x0F) to translate segmented addresses in display lists to physical RDRAM addresses. A segment address encodes the segment number in bits 28-31 and the offset in bits 0-27.

```go
type SegmentTable struct {
    bases [16]uint32
}
```

### NewSegmentTable

Creates a segment table with all segments set to 0:

```go
st := dma.NewSegmentTable()
```

### Set / Get

Set or query a segment's base address:

```go
st.Set(6, 0x80200000)       // segment 6 maps to RDRAM 0x80200000
base := st.Get(6)            // returns 0x80200000
```

Segment numbers must be 0-15. Out-of-range values are silently ignored (Set) or return 0 (Get).

### Resolve

Translates a segmented address to a physical RDRAM address:

```go
physAddr := st.Resolve(0x06001000)
// segment 6 base + offset 0x001000
```

The segment number is extracted from bits 24-27, the offset from bits 0-23. The result is `bases[segment] + offset`.

### MakeSegAddr

Creates a segmented address from a segment number and offset:

```go
addr := dma.MakeSegAddr(6, 0x1000)
// returns 0x06001000
```

## Typical usage pattern

A common N64 game initialization sequence:

```go
// Set up segment table
segments := dma.NewSegmentTable()

// Load level data from cartridge
levelData := make([]byte, levelSize)
dma.CartToRDRAM(levelROMOffset, levelData)

// If data is compressed
decompressed, err := dma.DecompressMIO0(levelData)

// Set up memory pool for the level
pool := dma.NewPool(poolBase, poolSize)

// Allocate permanent data from the head
modelAddr, _ := pool.AllocHead(modelSize)

// Set segment to point at the loaded data
segments.Set(6, modelAddr)

// Allocate per-frame display list from the tail
dlAddr, _ := pool.AllocTail(dlSize)

// At frame end, reset tail allocations
pool.ResetTail()
```
