# Memory Pools

The `dma.Pool` is a linear memory allocator for RDRAM, modeled after the allocation pattern used in Super Mario 64 (`main_pool_alloc` / `alloc_display_list`). It allocates from both ends of a contiguous memory region: the "head" grows upward for permanent allocations, and the "tail" grows downward for temporary per-frame allocations that are reset each frame.

## Why a pool allocator?

The N64 has 4MB (or 8MB with the Expansion Pak) of RDRAM and no virtual memory. Go's garbage collector is not available on bare-metal N64 builds. A pool allocator gives you:

- **Deterministic allocation** with no fragmentation (within the linear model)
- **Fast per-frame reset** for temporary data like display lists and vertex buffers
- **Simple lifetime management**: permanent data at the head, temporary data at the tail

## Pool

```go
type Pool struct {
    base     uint32
    size     uint32
    headUsed uint32
    tailUsed uint32
}
```

The pool manages a single contiguous RDRAM region from `base` to `base + size`.

```
 ┌──────────────────────────────────────────────┐
 │  HEAD (permanent) ──>          <── TAIL (temp)│
 │  [used]           [available]        [used]   │
 └──────────────────────────────────────────────┘
 base                                    base+size
```

## NewPool

Creates a memory pool at the given RDRAM address with the given size:

```go
pool := dma.NewPool(0x80200000, 512*1024) // 512KB pool at 0x80200000
```

The pool starts empty with zero bytes used on both ends.

## AllocHead

Allocates bytes from the head (permanent, grows up). Returns the RDRAM address of the allocation:

```go
addr, err := pool.AllocHead(1024)
if err != nil {
    // pool exhausted
}
```

All allocations are aligned to 16 bytes. If the requested size plus existing usage exceeds the pool, `ErrPoolExhausted` is returned.

Head allocations are permanent for the lifetime of the pool. Use these for data that persists across frames: level geometry, textures, loaded assets.

## AllocTail

Allocates bytes from the tail (temporary, grows down). Returns the RDRAM address of the allocation:

```go
addr, err := pool.AllocTail(2048)
if err != nil {
    // pool exhausted
}
```

Tail allocations grow downward from `base + size`. Like head allocations, they are 16-byte aligned.

Use tail allocations for per-frame data: display lists, transformed vertex buffers, scratch matrices.

## ResetTail

Frees all tail allocations at once. Call this at the start of each frame:

```go
pool.ResetTail()
```

After reset, all tail memory is available again. Head allocations are unaffected.

## Querying pool state

```go
pool.Available()  // remaining bytes between head and tail
pool.HeadUsed()   // bytes allocated from the head
pool.TailUsed()   // bytes allocated from the tail
pool.Base()       // pool's base RDRAM address
```

## Error handling

```go
var (
    ErrPoolExhausted = errors.New("dma: memory pool exhausted")
    ErrInvalidFree   = errors.New("dma: invalid free address")
)
```

Both `AllocHead` and `AllocTail` return `ErrPoolExhausted` when the remaining space between head and tail cannot satisfy the request.

## Alignment

All allocations are rounded up to the nearest 16-byte boundary. This ensures proper alignment for DMA transfers and RSP data structures, which require 8-byte or 16-byte alignment on the N64.

For example, `AllocHead(100)` actually consumes 112 bytes (next multiple of 16).

## Usage patterns

### Level lifecycle

```go
// At level load
pool := dma.NewPool(poolBase, poolSize)

// Permanent allocations (persist for entire level)
geoAddr, _ := pool.AllocHead(geoSize)   // level geometry
texAddr, _ := pool.AllocHead(texSize)   // textures

// Per-frame loop
for {
    pool.ResetTail()

    // Temporary allocations (freed every frame)
    dlAddr, _ := pool.AllocTail(dlSize)     // display list
    vtxAddr, _ := pool.AllocTail(vtxSize)   // transformed vertices
    mtxAddr, _ := pool.AllocTail(mtxSize)   // matrix stack

    // Build and submit frame...
}
```

### Monitoring usage

```go
fmt.Printf("Pool: %d/%d bytes used (head=%d, tail=%d, free=%d)\n",
    pool.HeadUsed()+pool.TailUsed(),
    pool.HeadUsed()+pool.TailUsed()+pool.Available(),
    pool.HeadUsed(),
    pool.TailUsed(),
    pool.Available(),
)
```

### Multiple pools

You can create separate pools for different purposes:

```go
mainPool := dma.NewPool(0x80200000, 384*1024)  // general purpose
audioPool := dma.NewPool(0x80260000, 64*1024)  // audio buffers
```

This prevents audio allocations from competing with graphics allocations.
