package format

import (
	"encoding/binary"
	"fmt"
)

const mapPayloadSize = 12

type ParsedMap struct {
	Width       uint16
	Height      uint16
	LayerCount  uint16
	CellBits    uint8
	ChunkWidth  uint16
	ChunkHeight uint16
	Layers      []ParsedMapLayer
}

type ParsedMapLayer struct {
	SheetID uint16
	Cells   []uint16
}

func ParseMap(raw []byte) (ParsedMap, error) {
	var m ParsedMap

	h, err := ParseHeader(raw, "MAP2")
	if err != nil {
		return m, err
	}

	payload := raw[h.HeaderBytes : h.HeaderBytes+h.PayloadBytes]
	if len(payload) < mapPayloadSize {
		return m, fmt.Errorf("format: map payload too short: got %d bytes", len(payload))
	}

	m.Width = binary.LittleEndian.Uint16(payload[0:2])
	m.Height = binary.LittleEndian.Uint16(payload[2:4])
	m.LayerCount = binary.LittleEndian.Uint16(payload[4:6])
	m.CellBits = payload[6]
	m.ChunkWidth = binary.LittleEndian.Uint16(payload[8:10])
	m.ChunkHeight = binary.LittleEndian.Uint16(payload[10:12])

	if m.CellBits != 8 && m.CellBits != 16 {
		return m, fmt.Errorf("format: unsupported cell width %d", m.CellBits)
	}

	cellCount := int(m.Width) * int(m.Height)
	if cellCount == 0 {
		return m, nil
	}

	cellBytes := 1
	if m.CellBits == 16 {
		cellBytes = 2
	}

	remaining := payload[mapPayloadSize:]
	if len(remaining) == 0 {
		m.Layers = make([]ParsedMapLayer, int(m.LayerCount))
		for i := range m.Layers {
			m.Layers[i].SheetID = 1
			m.Layers[i].Cells = make([]uint16, cellCount)
		}
		return m, nil
	}

	layerHeaderBytes := int(m.LayerCount) * 2
	wantBytes := int(m.LayerCount) * cellCount * cellBytes
	withLayerHeaders := layerHeaderBytes + wantBytes
	hasLayerHeaders := false
	switch len(remaining) {
	case wantBytes:
		hasLayerHeaders = false
	case withLayerHeaders:
		hasLayerHeaders = true
	default:
		if len(remaining) < wantBytes {
			return m, fmt.Errorf("format: map cell payload too short: got %d bytes, want at least %d", len(remaining), wantBytes)
		}
		return m, fmt.Errorf("format: unexpected map payload size %d", len(remaining))
	}

	m.Layers = make([]ParsedMapLayer, int(m.LayerCount))
	offset := 0
	sheetIDs := make([]uint16, int(m.LayerCount))
	if hasLayerHeaders {
		for i := range sheetIDs {
			sheetIDs[i] = binary.LittleEndian.Uint16(remaining[offset : offset+2])
			offset += 2
		}
	}
	for layerIdx := range m.Layers {
		layer := ParsedMapLayer{SheetID: 1, Cells: make([]uint16, cellCount)}
		if hasLayerHeaders {
			layer.SheetID = sheetIDs[layerIdx]
		}
		for i := range layer.Cells {
			if m.CellBits == 8 {
				layer.Cells[i] = uint16(remaining[offset])
				offset++
				continue
			}
			layer.Cells[i] = binary.LittleEndian.Uint16(remaining[offset : offset+2])
			offset += 2
		}
		m.Layers[layerIdx] = layer
	}

	return m, nil
}
