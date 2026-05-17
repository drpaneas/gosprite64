package dma

import "errors"

var (
	ErrPoolExhausted = errors.New("dma: memory pool exhausted")
	ErrInvalidFree   = errors.New("dma: invalid free address")
)

// Pool is a simple linear memory allocator for RDRAM, matching the N64's
// main_pool_alloc / alloc_display_list pattern used in SM64.
// It allocates from both ends: the "head" grows upward for permanent
// allocations, the "tail" grows downward for per-frame allocations
// that are reset each frame.
type Pool struct {
	base     uint32
	size     uint32
	headUsed uint32
	tailUsed uint32
}

// NewPool creates a memory pool starting at the given RDRAM address with the given size.
func NewPool(base, size uint32) *Pool {
	return &Pool{base: base, size: size}
}

// AllocHead allocates bytes from the head (permanent, grows up).
// Returns the RDRAM address of the allocation.
func (p *Pool) AllocHead(size uint32) (uint32, error) {
	aligned := (size + 15) & ^uint32(15) // 16-byte alignment
	if p.headUsed+aligned+p.tailUsed > p.size {
		return 0, ErrPoolExhausted
	}
	addr := p.base + p.headUsed
	p.headUsed += aligned
	return addr, nil
}

// AllocTail allocates bytes from the tail (temporary, grows down).
// Returns the RDRAM address of the allocation.
func (p *Pool) AllocTail(size uint32) (uint32, error) {
	aligned := (size + 15) & ^uint32(15)
	if p.headUsed+p.tailUsed+aligned > p.size {
		return 0, ErrPoolExhausted
	}
	p.tailUsed += aligned
	addr := p.base + p.size - p.tailUsed
	return addr, nil
}

// ResetTail frees all tail allocations (called each frame).
func (p *Pool) ResetTail() {
	p.tailUsed = 0
}

// HeadUsed returns bytes allocated from the head.
func (p *Pool) HeadUsed() uint32 { return p.headUsed }

// TailUsed returns bytes allocated from the tail.
func (p *Pool) TailUsed() uint32 { return p.tailUsed }

// Available returns remaining bytes between head and tail.
func (p *Pool) Available() uint32 {
	return p.size - p.headUsed - p.tailUsed
}

// Base returns the pool's base address.
func (p *Pool) Base() uint32 { return p.base }
