package format

import (
	"encoding/binary"
	"fmt"
)

type ParsedBundle struct {
	Entries []BundleEntry
}

const (
	BundleKindSheet uint8 = 1
	BundleKindMap   uint8 = 2
	BundleKindAnim  uint8 = 3
)

type BundleEntry struct {
	Kind uint8
	Name string
	Path string
}

func ParseBundle(raw []byte) (ParsedBundle, error) {
	var bundle ParsedBundle

	h, err := ParseHeader(raw, "BND2")
	if err != nil {
		return bundle, err
	}

	payload := raw[h.HeaderBytes : h.HeaderBytes+h.PayloadBytes]
	if len(payload) < 2 {
		return bundle, fmt.Errorf("format: bundle payload too short: got %d bytes", len(payload))
	}

	entryCount := int(binary.LittleEndian.Uint16(payload[:2]))
	offset := 2
	bundle.Entries = make([]BundleEntry, 0, entryCount)

	for range entryCount {
		if offset+3 > len(payload) {
			return ParsedBundle{}, fmt.Errorf("format: bundle entry header truncated")
		}

		entry := BundleEntry{
			Kind: payload[offset],
		}
		offset++

		nameLen := int(payload[offset])
		offset++
		pathLen := int(payload[offset])
		offset++

		if offset+nameLen+pathLen > len(payload) {
			return ParsedBundle{}, fmt.Errorf("format: bundle entry payload truncated")
		}

		entry.Name = string(payload[offset : offset+nameLen])
		offset += nameLen
		entry.Path = string(payload[offset : offset+pathLen])
		offset += pathLen

		bundle.Entries = append(bundle.Entries, entry)
	}

	return bundle, nil
}
