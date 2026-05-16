package gosprite64

import (
	"fmt"

	"github.com/drpaneas/gosprite64/internal/tile2d/format"
	tileloader "github.com/drpaneas/gosprite64/internal/tile2d/loader"
)

type Bundle struct {
	manifest format.ParsedBundle
	loader   tileloader.Loader
}

func OpenBundle(path string) (*Bundle, error) {
	return OpenBundleWithLoader(path, cartLoader{})
}

func OpenBundleWithLoader(path string, l tileloader.Loader) (*Bundle, error) {
	if l == nil {
		return nil, fmt.Errorf("open bundle: nil loader")
	}

	manifest, err := tileloader.OpenBundle(path, l)
	if err != nil {
		return nil, err
	}

	return &Bundle{
		manifest: manifest,
		loader:   l,
	}, nil
}

func (b *Bundle) ManifestCount() int {
	if b == nil {
		return 0
	}
	return len(b.manifest.Entries)
}

func (b *Bundle) LoadSheet(name string) (*Sheet, error) {
	entry, err := b.entryByKindAndName(format.BundleKindSheet, name)
	if err != nil {
		return nil, err
	}
	return b.loadSheetEntry(entry)
}

func (b *Bundle) LoadMap(name string) (*Map, error) {
	entry, err := b.entryByKindAndName(format.BundleKindMap, name)
	if err != nil {
		return nil, err
	}
	return b.loadMapEntry(entry)
}

func (b *Bundle) LoadAnimation(name string) (*AnimationSet, error) {
	entry, err := b.entryByKindAndName(format.BundleKindAnim, name)
	if err != nil {
		return nil, err
	}
	return b.loadAnimEntry(entry)
}

func (b *Bundle) entryByKindAndName(kind uint8, name string) (format.BundleEntry, error) {
	if b == nil {
		return format.BundleEntry{}, fmt.Errorf("bundle is nil")
	}

	for _, entry := range b.manifest.Entries {
		if entry.Kind == kind && entry.Name == name {
			return entry, nil
		}
	}

	return format.BundleEntry{}, fmt.Errorf("bundle entry %q of kind %d not found", name, kind)
}

func (b *Bundle) loadSheetEntry(entry format.BundleEntry) (*Sheet, error) {
	parsed, err := tileloader.LoadSheet(entry.Path, b.loader)
	if err != nil {
		return nil, err
	}
	return &Sheet{parsed: parsed}, nil
}

func (b *Bundle) loadMapEntry(entry format.BundleEntry) (*Map, error) {
	parsed, err := tileloader.LoadMap(entry.Path, b.loader)
	if err != nil {
		return nil, err
	}
	return &Map{parsed: parsed}, nil
}

func (b *Bundle) loadAnimEntry(entry format.BundleEntry) (*AnimationSet, error) {
	parsed, err := tileloader.LoadAnim(entry.Path, b.loader)
	if err != nil {
		return nil, err
	}
	return &AnimationSet{name: entry.Name, parsed: parsed}, nil
}
