package residency

import "testing"

func TestChunkCachePreloadsVisibleRing(t *testing.T) {
	c := NewChunkCache(ChunkPolicy{Radius: 1, Capacity: 9})
	c.EnsureVisible(ChunkCoord{X: 3, Y: 2})
	if got := c.Len(); got != 9 {
		t.Fatalf("Len() = %d, want 9", got)
	}
}

func TestResidencyStoreTracksLifetimeClass(t *testing.T) {
	s := NewStore()
	s.Pin("hud.sheet", ResidencyPermanent)
	if got := s.ClassOf("hud.sheet"); got != ResidencyPermanent {
		t.Fatalf("ClassOf() = %v, want %v", got, ResidencyPermanent)
	}
}
