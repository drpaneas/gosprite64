package rspq

// MicrocodeType identifies which RSP microcode to load.
type MicrocodeType int

const (
	// Fast3D is the original Fast3D microcode used by SM64.
	Fast3D MicrocodeType = iota
	// F3DEX is the extended Fast3D with more vertex buffer slots.
	F3DEX
	// F3DEX2 is Fast3D EX version 2, the most common N64 microcode.
	F3DEX2
	// AspMain is the standard audio microcode.
	AspMain
)

func (m MicrocodeType) String() string {
	switch m {
	case Fast3D:
		return "Fast3D"
	case F3DEX:
		return "F3DEX"
	case F3DEX2:
		return "F3DEX2"
	case AspMain:
		return "aspMain"
	default:
		return "unknown"
	}
}

// Microcode holds the IMEM (code) and DMEM (data) blobs for an RSP microcode.
type Microcode struct {
	Type MicrocodeType
	Code []byte // IMEM content (up to 4KB)
	Data []byte // DMEM content (up to 4KB)
}

// MaxIMEMSize is the RSP instruction memory size.
const MaxIMEMSize = 4096

// MaxDMEMSize is the RSP data memory size.
const MaxDMEMSize = 4096
