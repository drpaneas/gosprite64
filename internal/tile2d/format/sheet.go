package format

import (
	"encoding/binary"
	"fmt"
)

const sheetPayloadSize = 16
const legacySheetPayloadSize = 12

type ParsedSheet struct {
	TileWidth      uint16
	TileHeight     uint16
	TileCount      uint16
	PaletteEntries uint16
	AtlasWidth     uint16
	AtlasHeight    uint16
	DataOffset     uint32
	Pixels         []byte
}

func ParseSheet(raw []byte) (ParsedSheet, error) {
	var sheet ParsedSheet

	h, err := ParseHeader(raw, "SHT2")
	if err != nil {
		return sheet, err
	}

	payload := raw[h.HeaderBytes : h.HeaderBytes+h.PayloadBytes]
	if len(payload) < legacySheetPayloadSize {
		return sheet, fmt.Errorf("format: sheet payload too short: got %d bytes", len(payload))
	}

	sheet.TileWidth = binary.LittleEndian.Uint16(payload[0:2])
	sheet.TileHeight = binary.LittleEndian.Uint16(payload[2:4])
	sheet.TileCount = binary.LittleEndian.Uint16(payload[4:6])
	sheet.PaletteEntries = binary.LittleEndian.Uint16(payload[6:8])
	if len(payload) < sheetPayloadSize {
		sheet.DataOffset = binary.LittleEndian.Uint32(payload[8:12])
		return sheet, nil
	}
	sheet.AtlasWidth = binary.LittleEndian.Uint16(payload[8:10])
	sheet.AtlasHeight = binary.LittleEndian.Uint16(payload[10:12])
	sheet.DataOffset = binary.LittleEndian.Uint32(payload[12:16])

	if sheet.AtlasWidth == 0 || sheet.AtlasHeight == 0 {
		return sheet, fmt.Errorf("format: invalid sheet atlas size %dx%d", sheet.AtlasWidth, sheet.AtlasHeight)
	}

	pixelCount := int(sheet.AtlasWidth) * int(sheet.AtlasHeight) * 4
	if pixelCount == 0 {
		return sheet, nil
	}
	if int(sheet.DataOffset) < int(h.HeaderBytes)+sheetPayloadSize {
		return sheet, fmt.Errorf("format: invalid sheet data offset %d", sheet.DataOffset)
	}
	if int(sheet.DataOffset)+pixelCount > len(raw) {
		return sheet, fmt.Errorf("format: sheet pixel payload too short: got %d bytes, want %d", len(raw)-int(sheet.DataOffset), pixelCount)
	}

	sheet.Pixels = append([]byte(nil), raw[sheet.DataOffset:int(sheet.DataOffset)+pixelCount]...)

	return sheet, nil
}
