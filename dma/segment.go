package dma

// SegmentTable manages N64 segment base addresses.
// The N64 uses 16 segment registers (0x00-0x0F) to translate
// segment addresses in display lists to physical RDRAM addresses.
type SegmentTable struct {
	bases [16]uint32
}

// NewSegmentTable creates a segment table with all segments set to 0.
func NewSegmentTable() *SegmentTable {
	return &SegmentTable{}
}

// Set sets the base RDRAM address for a segment (0-15).
func (st *SegmentTable) Set(segment uint8, base uint32) {
	if segment < 16 {
		st.bases[segment] = base
	}
}

// Get returns the base address for a segment.
func (st *SegmentTable) Get(segment uint8) uint32 {
	if segment < 16 {
		return st.bases[segment]
	}
	return 0
}

// Resolve translates a segment address to a physical RDRAM address.
// A segment address has the segment number in bits 28-31 and the
// offset in bits 0-27.
func (st *SegmentTable) Resolve(addr uint32) uint32 {
	seg := (addr >> 24) & 0x0F
	offset := addr & 0x00FFFFFF
	return st.bases[seg] + offset
}

// MakeSegAddr creates a segment address from a segment number and offset.
func MakeSegAddr(segment uint8, offset uint32) uint32 {
	return uint32(segment)<<24 | (offset & 0x00FFFFFF)
}
