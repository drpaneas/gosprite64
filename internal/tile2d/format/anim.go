package format

import (
	"encoding/binary"
	"fmt"
)

type ParsedAnim struct {
	Clips []ParsedClip
}

type ParsedClip struct {
	Name   string
	FPS    uint16
	Frames []uint16
}

func ParseAnim(raw []byte) (ParsedAnim, error) {
	var anim ParsedAnim

	h, err := ParseHeader(raw, "ANM2")
	if err != nil {
		return anim, err
	}

	payload := raw[h.HeaderBytes : h.HeaderBytes+h.PayloadBytes]
	if len(payload) < 2 {
		return anim, fmt.Errorf("format: anim payload too short: got %d bytes", len(payload))
	}

	clipCount := int(binary.LittleEndian.Uint16(payload[:2]))
	offset := 2
	anim.Clips = make([]ParsedClip, 0, clipCount)

	for range clipCount {
		if offset >= len(payload) {
			return ParsedAnim{}, fmt.Errorf("format: anim clip header truncated")
		}

		nameLen := int(payload[offset])
		offset++
		if offset+nameLen+4 > len(payload) {
			return ParsedAnim{}, fmt.Errorf("format: anim clip payload truncated")
		}

		clip := ParsedClip{
			Name: string(payload[offset : offset+nameLen]),
		}
		offset += nameLen

		frameCount := int(binary.LittleEndian.Uint16(payload[offset : offset+2]))
		offset += 2
		clip.FPS = binary.LittleEndian.Uint16(payload[offset : offset+2])
		offset += 2

		if offset+frameCount*2 > len(payload) {
			return ParsedAnim{}, fmt.Errorf("format: anim frames truncated")
		}

		clip.Frames = make([]uint16, frameCount)
		for i := range frameCount {
			clip.Frames[i] = binary.LittleEndian.Uint16(payload[offset : offset+2])
			offset += 2
		}

		anim.Clips = append(anim.Clips, clip)
	}

	return anim, nil
}
