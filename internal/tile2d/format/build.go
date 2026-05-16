package format

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
)

type MapConfig struct {
	Width       uint16 `json:"width"`
	Height      uint16 `json:"height"`
	LayerCount  uint16 `json:"layer_count"`
	CellBits    uint8  `json:"cell_bits"`
	ChunkWidth  uint16 `json:"chunk_width"`
	ChunkHeight uint16 `json:"chunk_height"`
	Layers      []MapLayerConfig `json:"layers"`
}

type MapLayerConfig struct {
	SheetID uint16   `json:"sheet_id"`
	Cells   []uint16 `json:"cells"`
}

type AnimConfig struct {
	Clips []AnimClipConfig `json:"clips"`
}

type AnimClipConfig struct {
	Name   string   `json:"name"`
	FPS    uint16   `json:"fps"`
	Frames []uint16 `json:"frames"`
}

func BuildSheet(img image.Image, tileWidth, tileHeight int) ([]byte, error) {
	if img == nil {
		return nil, fmt.Errorf("format: nil image")
	}
	if tileWidth <= 0 || tileHeight <= 0 {
		return nil, fmt.Errorf("format: tile size must be positive")
	}

	bounds := img.Bounds()
	if bounds.Dx()%tileWidth != 0 || bounds.Dy()%tileHeight != 0 {
		return nil, fmt.Errorf("format: image size %dx%d not divisible by tile size %dx%d", bounds.Dx(), bounds.Dy(), tileWidth, tileHeight)
	}

	tileCount := (bounds.Dx() / tileWidth) * (bounds.Dy() / tileHeight)
	if tileCount > 0xFFFF {
		return nil, fmt.Errorf("format: tile count %d exceeds uint16", tileCount)
	}

	paletteEntries := uniqueColorCount(img)
	if paletteEntries > 0xFFFF {
		return nil, fmt.Errorf("format: palette entries %d exceeds uint16", paletteEntries)
	}

	payload := make([]byte, sheetPayloadSize)
	binary.LittleEndian.PutUint16(payload[0:2], uint16(tileWidth))
	binary.LittleEndian.PutUint16(payload[2:4], uint16(tileHeight))
	binary.LittleEndian.PutUint16(payload[4:6], uint16(tileCount))
	binary.LittleEndian.PutUint16(payload[6:8], uint16(paletteEntries))
	binary.LittleEndian.PutUint16(payload[8:10], uint16(bounds.Dx()))
	binary.LittleEndian.PutUint16(payload[10:12], uint16(bounds.Dy()))
	binary.LittleEndian.PutUint32(payload[12:16], headerSize+sheetPayloadSize)

	nrgba := toNRGBA(img)
	return encodeAsset("SHT2", append(payload, nrgba.Pix...)), nil
}

func BuildMap(cfg MapConfig) ([]byte, error) {
	if cfg.Width == 0 || cfg.Height == 0 {
		return nil, fmt.Errorf("format: map dimensions must be non-zero")
	}
	if cfg.LayerCount == 0 {
		return nil, fmt.Errorf("format: layer_count must be non-zero")
	}
	if cfg.CellBits != 8 && cfg.CellBits != 16 {
		return nil, fmt.Errorf("format: unsupported cell width %d", cfg.CellBits)
	}
	if cfg.ChunkWidth == 0 || cfg.ChunkHeight == 0 {
		return nil, fmt.Errorf("format: chunk dimensions must be non-zero")
	}
	if len(cfg.Layers) > 0 && len(cfg.Layers) != int(cfg.LayerCount) {
		return nil, fmt.Errorf("format: got %d layers, want %d", len(cfg.Layers), cfg.LayerCount)
	}

	payload := make([]byte, mapPayloadSize)
	binary.LittleEndian.PutUint16(payload[0:2], cfg.Width)
	binary.LittleEndian.PutUint16(payload[2:4], cfg.Height)
	binary.LittleEndian.PutUint16(payload[4:6], cfg.LayerCount)
	payload[6] = cfg.CellBits
	binary.LittleEndian.PutUint16(payload[8:10], cfg.ChunkWidth)
	binary.LittleEndian.PutUint16(payload[10:12], cfg.ChunkHeight)

	cellCount := int(cfg.Width) * int(cfg.Height)
	var cellPayload bytes.Buffer
	var layerPayload bytes.Buffer
	for layerIdx := 0; layerIdx < int(cfg.LayerCount); layerIdx++ {
		var cells []uint16
		sheetID := uint16(1)
		if layerIdx < len(cfg.Layers) {
			if cfg.Layers[layerIdx].SheetID != 0 {
				sheetID = cfg.Layers[layerIdx].SheetID
			}
			cells = cfg.Layers[layerIdx].Cells
		}
		if err := binary.Write(&layerPayload, binary.LittleEndian, sheetID); err != nil {
			return nil, err
		}
		if len(cells) == 0 {
			cells = make([]uint16, cellCount)
		}
		if len(cells) != cellCount {
			return nil, fmt.Errorf("format: layer %d has %d cells, want %d", layerIdx, len(cells), cellCount)
		}
		for _, cell := range cells {
			if cfg.CellBits == 8 {
				if cell > 0xFF {
					return nil, fmt.Errorf("format: cell value %d exceeds 8-bit storage", cell)
				}
				cellPayload.WriteByte(byte(cell))
				continue
			}
			if err := binary.Write(&cellPayload, binary.LittleEndian, cell); err != nil {
				return nil, err
			}
		}
	}

	return encodeAsset("MAP2", append(append(payload, layerPayload.Bytes()...), cellPayload.Bytes()...)), nil
}

func BuildAnim(cfg AnimConfig) ([]byte, error) {
	var payload bytes.Buffer
	if len(cfg.Clips) > 0xFFFF {
		return nil, fmt.Errorf("format: clip count %d exceeds uint16", len(cfg.Clips))
	}
	if err := binary.Write(&payload, binary.LittleEndian, uint16(len(cfg.Clips))); err != nil {
		return nil, err
	}

	for _, clip := range cfg.Clips {
		if len(clip.Name) > 0xFF {
			return nil, fmt.Errorf("format: clip name %q too long", clip.Name)
		}
		if len(clip.Frames) > 0xFFFF {
			return nil, fmt.Errorf("format: clip %q has too many frames", clip.Name)
		}

		payload.WriteByte(byte(len(clip.Name)))
		payload.WriteString(clip.Name)
		if err := binary.Write(&payload, binary.LittleEndian, uint16(len(clip.Frames))); err != nil {
			return nil, err
		}
		if err := binary.Write(&payload, binary.LittleEndian, clip.FPS); err != nil {
			return nil, err
		}
		for _, frame := range clip.Frames {
			if err := binary.Write(&payload, binary.LittleEndian, frame); err != nil {
				return nil, err
			}
		}
	}

	return encodeAsset("ANM2", payload.Bytes()), nil
}

func BuildBundle(entries []BundleEntry) ([]byte, error) {
	var payload bytes.Buffer
	if len(entries) > 0xFFFF {
		return nil, fmt.Errorf("format: entry count %d exceeds uint16", len(entries))
	}
	if err := binary.Write(&payload, binary.LittleEndian, uint16(len(entries))); err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if len(entry.Name) > 0xFF || len(entry.Path) > 0xFF {
			return nil, fmt.Errorf("format: bundle entry %q too long", entry.Name)
		}
		payload.WriteByte(entry.Kind)
		payload.WriteByte(byte(len(entry.Name)))
		payload.WriteByte(byte(len(entry.Path)))
		payload.WriteString(entry.Name)
		payload.WriteString(entry.Path)
	}

	return encodeAsset("BND2", payload.Bytes()), nil
}

func encodeAsset(magic string, payload []byte) []byte {
	raw := make([]byte, headerSize+len(payload))
	copy(raw[:4], []byte(magic))
	binary.LittleEndian.PutUint16(raw[4:6], 1)
	binary.LittleEndian.PutUint32(raw[8:12], headerSize)
	binary.LittleEndian.PutUint32(raw[12:16], uint32(len(payload)))
	copy(raw[headerSize:], payload)
	return raw
}

func uniqueColorCount(img image.Image) int {
	colors := make(map[color.NRGBA]struct{})
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			colors[color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)] = struct{}{}
		}
	}
	return len(colors)
}

func toNRGBA(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	dst := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x-bounds.Min.X, y-bounds.Min.Y, img.At(x, y))
		}
	}
	return dst
}
