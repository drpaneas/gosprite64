package gfx

// F3D command opcodes (Fast3D original, matching SM64 default build).
const (
	// SP commands (RSP geometry)
	OpSPNoop           = 0x00
	OpSPMatrix         = 0x01
	OpSPVertex         = 0x04
	OpSPDisplayList    = 0x06
	OpSPMoveMem        = 0x03
	OpSP1Triangle      = 0xBF
	OpSPPopMatrix      = 0xBD
	OpSPSetGeomMode    = 0xB7
	OpSPClearGeomMode  = 0xB6
	OpSPEndDisplayList = 0xB8
	OpSPMoveWord       = 0xBC
	OpSPTexture        = 0xBB
	OpSPRDPHalf1       = 0xB4

	// DP commands (RDP rasterizer)
	OpDPSetColorImage   = 0xFF
	OpDPSetTextureImage = 0xFD
	OpDPSetCombine      = 0xFC
	OpDPSetTile         = 0xF5
	OpDPLoadBlock       = 0xF3
	OpDPSetTileSize     = 0xF2
	OpDPSetOtherModeH   = 0xBA
	OpDPSetOtherModeL   = 0xB9
	OpDPSetEnvColor     = 0xFB
	OpDPSetPrimColor    = 0xFA
	OpDPSetFogColor     = 0xF8
	OpDPSetFillColor    = 0xF7
	OpDPFillRect        = 0xF6
	OpDPPipeSync        = 0xE7
	OpDPTileSync        = 0xE8
	OpDPLoadSync        = 0xE6
	OpDPFullSync        = 0xE9
	OpDPSetScissor      = 0xED
	OpDPSetZImage       = 0xFE
)

// G_MTX flags (F3D-style, not F3DEX2).
const (
	MtxModelView  = 0x00
	MtxProjection = 0x01
	MtxMul        = 0x00
	MtxLoad       = 0x02
	MtxNoPush     = 0x00
	MtxPush       = 0x04
)

// G_MV indices for gSPMoveMem (F3D).
const (
	MvViewport = 0x80
)

// G_MW indices for gSPMoveWord.
const (
	MwSegment   = 0x06
	MwPerspNorm = 0x0E
)

// Geometry mode flags.
const (
	GZBuffer       = 0x00000001
	GShade         = 0x00000004
	GCullFront     = 0x00000200
	GCullBack      = 0x00000400
	GCullBoth      = GCullFront | GCullBack
	GFog           = 0x00010000
	GLighting      = 0x00020000
	GTexGen        = 0x00040000
	GTexGenLinear  = 0x00080000
	GShadingSmooth = 0x00200000
)

// Image format constants.
const (
	FmtRGBA = 0
	FmtYUV  = 1
	FmtCI   = 2
	FmtIA   = 3
	FmtI    = 4

	Siz4b  = 0
	Siz8b  = 1
	Siz16b = 2
	Siz32b = 3
)

// DL push/nopush.
const (
	DLPush   = 0x00
	DLNoPush = 0x01
)

// Render mode cycle 1 presets.
const (
	RMAAZBOpa1   = 0x00442078
	RMAAZBOpa2   = 0x00112078
	RMAAZBXlu1   = 0x00442478
	RMAAZBXlu2   = 0x00112478
	RMOpaOpaque1 = 0x0F0A4000
	RMOpaOpaque2 = 0x00102000
	RMXluSurf1   = 0x00404240
	RMXluSurf2   = 0x00104240
)
