package residency

type ChunkPolicy struct {
	Radius   int
	Capacity int
}

type ChunkCoord struct {
	X int
	Y int
}

type ChunkCache struct {
	policy ChunkPolicy
	loaded map[ChunkCoord]struct{}
	order  []ChunkCoord
}

func NewChunkCache(policy ChunkPolicy) *ChunkCache {
	return &ChunkCache{
		policy: policy,
		loaded: make(map[ChunkCoord]struct{}),
	}
}

func (c *ChunkCache) EnsureVisible(center ChunkCoord) {
	if c == nil {
		return
	}

	radius := max(0, c.policy.Radius)
	for y := center.Y - radius; y <= center.Y+radius; y++ {
		for x := center.X - radius; x <= center.X+radius; x++ {
			c.touch(ChunkCoord{X: x, Y: y})
		}
	}
}

func (c *ChunkCache) Len() int {
	if c == nil {
		return 0
	}
	return len(c.loaded)
}

func (c *ChunkCache) touch(coord ChunkCoord) {
	if _, ok := c.loaded[coord]; ok {
		return
	}

	c.loaded[coord] = struct{}{}
	c.order = append(c.order, coord)
	c.trim()
}

func (c *ChunkCache) trim() {
	if c.policy.Capacity <= 0 {
		return
	}

	for len(c.order) > c.policy.Capacity {
		evict := c.order[0]
		c.order = c.order[1:]
		delete(c.loaded, evict)
	}
}
