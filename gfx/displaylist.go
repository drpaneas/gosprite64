package gfx

// DisplayList builds a sequence of N64 Gfx commands.
type DisplayList struct {
	cmds []Gfx
}

// NewDisplayList creates a display list with the given initial capacity.
func NewDisplayList(capacity int) *DisplayList {
	return &DisplayList{cmds: make([]Gfx, 0, capacity)}
}

// Len returns the number of commands in the display list.
func (dl *DisplayList) Len() int { return len(dl.cmds) }

// Commands returns the raw command slice.
func (dl *DisplayList) Commands() []Gfx { return dl.cmds }

// Reset clears the display list for reuse without deallocating.
func (dl *DisplayList) Reset() { dl.cmds = dl.cmds[:0] }

func (dl *DisplayList) append(w0, w1 uint32) {
	dl.cmds = append(dl.cmds, Gfx{W0: w0, W1: w1})
}

func (dl *DisplayList) appendRaw(words ...uint64) {
	if len(words) == 0 {
		return
	}
	raw := make([]uint64, len(words))
	copy(raw, words)
	dl.cmds = append(dl.cmds, Gfx{Raw: raw})
}

func shiftL(v, shift, width uint32) uint32 {
	return (v & ((1 << width) - 1)) << shift
}

// SPMatrix appends a matrix load/multiply command.
// addr is the RDRAM address (or segment address) of an N64Mtx.
// flags is a combination of MtxProjection/MtxModelView, MtxLoad/MtxMul, MtxPush/MtxNoPush.
func (dl *DisplayList) SPMatrix(addr uint32, flags uint8) {
	w0 := shiftL(OpSPMatrix, 24, 8) | shiftL(uint32(flags), 16, 8) | shiftL(64, 0, 16) // sizeof(Mtx) = 64
	dl.append(w0, addr)
}

// SPVertex loads vertices into the RSP vertex buffer.
// addr is the address of the vertex array, n is the count (1-16), v0 is the start index.
func (dl *DisplayList) SPVertex(addr uint32, n, v0 uint8) {
	w0 := shiftL(OpSPVertex, 24, 8) | shiftL(uint32((n-1)<<4|v0), 16, 8) | shiftL(uint32(n)*16, 0, 16)
	dl.append(w0, addr)
}

// SP1Triangle draws a single triangle. v0, v1, v2 are vertex buffer indices.
func (dl *DisplayList) SP1Triangle(v0, v1, v2, flag uint8) {
	w0 := shiftL(OpSP1Triangle, 24, 8)
	w1 := shiftL(uint32(flag), 24, 8) | shiftL(uint32(v0)*10, 16, 8) |
		shiftL(uint32(v1)*10, 8, 8) | shiftL(uint32(v2)*10, 0, 8)
	dl.append(w0, w1)
}

// SP2Triangles draws two triangles. In F3D, this emits two SP1Triangle commands.
func (dl *DisplayList) SP2Triangles(v00, v01, v02, flag0, v10, v11, v12, flag1 uint8) {
	dl.SP1Triangle(v00, v01, v02, flag0)
	dl.SP1Triangle(v10, v11, v12, flag1)
}

// SPDisplayList calls a child display list (with return).
func (dl *DisplayList) SPDisplayList(addr uint32) {
	w0 := shiftL(OpSPDisplayList, 24, 8) | shiftL(DLPush, 16, 8)
	dl.append(w0, addr)
}

// SPBranchList jumps to a display list (no return).
func (dl *DisplayList) SPBranchList(addr uint32) {
	w0 := shiftL(OpSPDisplayList, 24, 8) | shiftL(DLNoPush, 16, 8)
	dl.append(w0, addr)
}

// SPEndDisplayList terminates display list processing.
func (dl *DisplayList) SPEndDisplayList() {
	dl.append(shiftL(OpSPEndDisplayList, 24, 8), 0)
}

// SPSegment sets a segment register for address translation.
func (dl *DisplayList) SPSegment(segment uint8, base uint32) {
	w0 := shiftL(OpSPMoveWord, 24, 8) | shiftL(MwSegment, 8, 8) | shiftL(uint32(segment)*4, 0, 8)
	dl.append(w0, base)
}

// SPViewport sets the viewport parameters.
func (dl *DisplayList) SPViewport(addr uint32) {
	w0 := shiftL(OpSPMoveMem, 24, 8) | shiftL(MvViewport, 16, 8) | shiftL(16, 0, 16) // sizeof(Vp) = 16
	dl.append(w0, addr)
}

// SPPerspNormalize sets the perspective normalization value (F3D path via RDPHALF_1).
func (dl *DisplayList) SPPerspNormalize(s uint16) {
	dl.append(shiftL(OpSPRDPHalf1, 24, 8), uint32(s))
}

// SPSetGeometryMode enables geometry mode flags.
func (dl *DisplayList) SPSetGeometryMode(flags uint32) {
	dl.append(shiftL(OpSPSetGeomMode, 24, 8), flags)
}

// SPClearGeometryMode disables geometry mode flags.
func (dl *DisplayList) SPClearGeometryMode(flags uint32) {
	dl.append(shiftL(OpSPClearGeomMode, 24, 8), flags)
}

// DPSetColorImage sets the RDP color framebuffer target.
func (dl *DisplayList) DPSetColorImage(fmt, siz uint8, width uint16, addr uint32) {
	w0 := shiftL(OpDPSetColorImage, 24, 8) | shiftL(uint32(fmt), 21, 3) |
		shiftL(uint32(siz), 19, 2) | shiftL(uint32(width-1), 0, 12)
	dl.append(w0, addr)
}

// DPSetTextureImage sets the texture source for subsequent loads.
func (dl *DisplayList) DPSetTextureImage(fmt, siz uint8, width uint16, addr uint32) {
	w0 := shiftL(OpDPSetTextureImage, 24, 8) | shiftL(uint32(fmt), 21, 3) |
		shiftL(uint32(siz), 19, 2) | shiftL(uint32(width-1), 0, 12)
	dl.append(w0, addr)
}

// DPSetZImage sets the Z-buffer address.
func (dl *DisplayList) DPSetZImage(addr uint32) {
	dl.append(shiftL(OpDPSetZImage, 24, 8), addr)
}

// DPSetFillColor sets the fill color (packed 16-bit RGBA doubled, or 32-bit).
func (dl *DisplayList) DPSetFillColor(color uint32) {
	dl.append(shiftL(OpDPSetFillColor, 24, 8), color)
}

// DPFillRect fills a rectangle with the current fill color.
// Coordinates are in 10.2 fixed-point format.
func (dl *DisplayList) DPFillRect(ulx, uly, lrx, lry uint16) {
	w0 := shiftL(OpDPFillRect, 24, 8) | shiftL(uint32(lrx), 12, 12) | shiftL(uint32(lry), 0, 12)
	w1 := shiftL(uint32(ulx), 12, 12) | shiftL(uint32(uly), 0, 12)
	dl.append(w0, w1)
}

// DPSetTile configures a tile descriptor in TMEM.
func (dl *DisplayList) DPSetTile(fmt, siz uint8, line uint16, tmem uint16,
	tile, palette uint8, cmt, maskt, shiftt, cms, masks, shifts uint8) {
	w0 := shiftL(OpDPSetTile, 24, 8) | shiftL(uint32(fmt), 21, 3) |
		shiftL(uint32(siz), 19, 2) | shiftL(uint32(line), 9, 9) | shiftL(uint32(tmem), 0, 9)
	w1 := shiftL(uint32(tile), 24, 3) | shiftL(uint32(palette), 20, 4) |
		shiftL(uint32(cmt), 18, 2) | shiftL(uint32(maskt), 14, 4) |
		shiftL(uint32(shiftt), 10, 4) | shiftL(uint32(cms), 8, 2) |
		shiftL(uint32(masks), 4, 4) | shiftL(uint32(shifts), 0, 4)
	dl.append(w0, w1)
}

// DPLoadBlock loads a contiguous block of texture data into TMEM.
func (dl *DisplayList) DPLoadBlock(tile uint8, uls, ult, lrs, dxt uint16) {
	maxTxl := uint16(2047)
	if lrs > maxTxl {
		lrs = maxTxl
	}
	w0 := shiftL(OpDPLoadBlock, 24, 8) | shiftL(uint32(uls), 12, 12) | shiftL(uint32(ult), 0, 12)
	w1 := shiftL(uint32(tile), 24, 3) | shiftL(uint32(lrs), 12, 12) | shiftL(uint32(dxt), 0, 12)
	dl.append(w0, w1)
}

// DPSetTileSize sets the texture tile size.
func (dl *DisplayList) DPSetTileSize(tile uint8, uls, ult, lrs, lrt uint16) {
	w0 := shiftL(OpDPSetTileSize, 24, 8) | shiftL(uint32(uls), 12, 12) | shiftL(uint32(ult), 0, 12)
	w1 := shiftL(uint32(tile), 24, 3) | shiftL(uint32(lrs), 12, 12) | shiftL(uint32(lrt), 0, 12)
	dl.append(w0, w1)
}

// DPSetCombineMode sets the RDP color combiner (raw 64-bit encoding).
func (dl *DisplayList) DPSetCombineMode(w0hi uint32, w1 uint32) {
	dl.append(shiftL(OpDPSetCombine, 24, 8)|w0hi, w1)
}

// DPSetRenderMode sets the RDP render mode via SetOtherModeL.
func (dl *DisplayList) DPSetRenderMode(c0, c1 uint32) {
	w0 := shiftL(OpDPSetOtherModeL, 24, 8) | shiftL(3, 8, 8) | shiftL(29, 0, 8)
	dl.append(w0, c0|c1)
}

// DPSetScissor sets the scissor rectangle. Coords in screen pixels.
func (dl *DisplayList) DPSetScissor(mode uint8, ulx, uly, lrx, lry uint16) {
	w0 := shiftL(OpDPSetScissor, 24, 8) |
		shiftL(uint32(ulx)<<2, 12, 12) | shiftL(uint32(uly)<<2, 0, 12)
	w1 := shiftL(uint32(mode), 24, 2) |
		shiftL(uint32(lrx)<<2, 12, 12) | shiftL(uint32(lry)<<2, 0, 12)
	dl.append(w0, w1)
}

// DPSetEnvColor sets the environment color.
func (dl *DisplayList) DPSetEnvColor(r, g, b, a uint8) {
	dl.append(shiftL(OpDPSetEnvColor, 24, 8),
		shiftL(uint32(r), 24, 8)|shiftL(uint32(g), 16, 8)|shiftL(uint32(b), 8, 8)|shiftL(uint32(a), 0, 8))
}

// DPSetPrimColor sets the primitive color.
func (dl *DisplayList) DPSetPrimColor(minLevel, fracLevel, r, g, b, a uint8) {
	w0 := shiftL(OpDPSetPrimColor, 24, 8) | shiftL(uint32(minLevel), 8, 8) | shiftL(uint32(fracLevel), 0, 8)
	w1 := shiftL(uint32(r), 24, 8) | shiftL(uint32(g), 16, 8) | shiftL(uint32(b), 8, 8) | shiftL(uint32(a), 0, 8)
	dl.append(w0, w1)
}

// DPSetFogColor sets the fog color.
func (dl *DisplayList) DPSetFogColor(r, g, b, a uint8) {
	dl.append(shiftL(OpDPSetFogColor, 24, 8),
		shiftL(uint32(r), 24, 8)|shiftL(uint32(g), 16, 8)|shiftL(uint32(b), 8, 8)|shiftL(uint32(a), 0, 8))
}

// DPPipeSync inserts a pipeline sync.
func (dl *DisplayList) DPPipeSync() {
	dl.append(shiftL(OpDPPipeSync, 24, 8), 0)
}

// DPTileSync inserts a tile sync.
func (dl *DisplayList) DPTileSync() {
	dl.append(shiftL(OpDPTileSync, 24, 8), 0)
}

// DPLoadSync inserts a load sync.
func (dl *DisplayList) DPLoadSync() {
	dl.append(shiftL(OpDPLoadSync, 24, 8), 0)
}

// DPFullSync signals RDP completion.
func (dl *DisplayList) DPFullSync() {
	dl.append(shiftL(OpDPFullSync, 24, 8), 0)
}

// DPRaw appends one or more fully packed raw RDP command words.
// This is intended for advanced paths, such as CPU-built triangle commands.
func (dl *DisplayList) DPRaw(words ...uint64) {
	dl.appendRaw(words...)
}
