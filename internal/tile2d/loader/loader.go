package loader

import "fmt"

type Loader interface {
	ReadAsset(path string) ([]byte, error)
}

type MemoryLoader struct {
	assets map[string][]byte
}

func NewMemoryLoader(assets map[string][]byte) MemoryLoader {
	cloned := make(map[string][]byte, len(assets))
	for path, raw := range assets {
		cloned[path] = append([]byte(nil), raw...)
	}
	return MemoryLoader{assets: cloned}
}

func (m MemoryLoader) ReadAsset(path string) ([]byte, error) {
	raw, ok := m.assets[path]
	if !ok {
		return nil, fmt.Errorf("loader: asset %q not found", path)
	}
	return append([]byte(nil), raw...), nil
}
