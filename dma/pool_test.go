package dma

import "testing"

func TestPoolAllocHead(t *testing.T) {
	p := NewPool(0x8000_0000, 1024)
	addr, err := p.AllocHead(32)
	if err != nil {
		t.Fatalf("AllocHead(32) returned error: %v", err)
	}
	if addr != 0x8000_0000 {
		t.Errorf("expected base address 0x80000000, got 0x%08X", addr)
	}
	if p.HeadUsed() != 32 {
		t.Errorf("expected HeadUsed=32, got %d", p.HeadUsed())
	}
}

func TestPoolAllocHeadAlignment(t *testing.T) {
	p := NewPool(0x8000_0000, 1024)

	addr1, err := p.AllocHead(3)
	if err != nil {
		t.Fatalf("AllocHead(3) returned error: %v", err)
	}
	if addr1 != 0x8000_0000 {
		t.Errorf("first alloc: expected 0x80000000, got 0x%08X", addr1)
	}
	if p.HeadUsed() != 16 {
		t.Errorf("3 bytes should round to 16, got HeadUsed=%d", p.HeadUsed())
	}

	addr2, err := p.AllocHead(1)
	if err != nil {
		t.Fatalf("AllocHead(1) returned error: %v", err)
	}
	if addr2 != 0x8000_0010 {
		t.Errorf("second alloc: expected 0x80000010, got 0x%08X", addr2)
	}
}

func TestPoolAllocTail(t *testing.T) {
	p := NewPool(0x8000_0000, 256)
	addr, err := p.AllocTail(32)
	if err != nil {
		t.Fatalf("AllocTail(32) returned error: %v", err)
	}
	expected := uint32(0x8000_0000 + 256 - 32)
	if addr != expected {
		t.Errorf("expected 0x%08X, got 0x%08X", expected, addr)
	}
}

func TestPoolExhaustion(t *testing.T) {
	p := NewPool(0x8000_0000, 128)

	_, err := p.AllocHead(64)
	if err != nil {
		t.Fatalf("AllocHead(64) returned error: %v", err)
	}
	_, err = p.AllocTail(64)
	if err != nil {
		t.Fatalf("AllocTail(64) returned error: %v", err)
	}

	_, err = p.AllocHead(1)
	if err != ErrPoolExhausted {
		t.Errorf("expected ErrPoolExhausted, got %v", err)
	}

	_, err = p.AllocTail(1)
	if err != ErrPoolExhausted {
		t.Errorf("expected ErrPoolExhausted from tail, got %v", err)
	}
}

func TestPoolResetTail(t *testing.T) {
	p := NewPool(0x8000_0000, 256)
	initial := p.Available()

	_, err := p.AllocTail(64)
	if err != nil {
		t.Fatalf("AllocTail(64) returned error: %v", err)
	}
	if p.Available() >= initial {
		t.Errorf("Available() should have decreased after tail alloc")
	}

	p.ResetTail()
	if p.Available() != initial {
		t.Errorf("after ResetTail: expected Available()=%d, got %d", initial, p.Available())
	}
}
