package math2d

import "testing"

func TestLayerMaskSameLayer(t *testing.T) {
	const LayerPlayer Layer = 1 << 0
	if !LayerPlayer.Matches(LayerPlayer) {
		t.Fatal("same layer should match")
	}
}

func TestLayerMaskDifferentLayers(t *testing.T) {
	const LayerPlayer Layer = 1 << 0
	const LayerEnemy Layer = 1 << 1
	if LayerPlayer.Matches(LayerEnemy) {
		t.Fatal("different layers should not match")
	}
}

func TestLayerMaskCombined(t *testing.T) {
	const LayerPlayer Layer = 1 << 0
	const LayerEnemy Layer = 1 << 1
	const LayerBullet Layer = 1 << 2

	mask := LayerPlayer | LayerEnemy
	if !mask.Matches(LayerPlayer) {
		t.Fatal("combined mask should match player")
	}
	if !mask.Matches(LayerEnemy) {
		t.Fatal("combined mask should match enemy")
	}
	if mask.Matches(LayerBullet) {
		t.Fatal("combined mask should not match bullet")
	}
}

func TestLayerNone(t *testing.T) {
	if LayerNone.Matches(Layer(1)) {
		t.Fatal("LayerNone should match nothing")
	}
}

func TestLayerAll(t *testing.T) {
	if !LayerAll.Matches(Layer(1)) {
		t.Fatal("LayerAll should match everything")
	}
	if !LayerAll.Matches(Layer(1 << 15)) {
		t.Fatal("LayerAll should match everything")
	}
}

func TestColliderCheckLayers(t *testing.T) {
	const LayerPlayer Layer = 1 << 0
	const LayerEnemy Layer = 1 << 1
	const LayerWall Layer = 1 << 2

	player := Collider{
		Bounds: Rect{X: 0, Y: 0, W: 10, H: 10},
		Layer:  LayerPlayer,
		Mask:   LayerEnemy | LayerWall,
	}
	enemy := Collider{
		Bounds: Rect{X: 5, Y: 5, W: 10, H: 10},
		Layer:  LayerEnemy,
		Mask:   LayerPlayer,
	}
	wall := Collider{
		Bounds: Rect{X: 5, Y: 5, W: 10, H: 10},
		Layer:  LayerWall,
		Mask:   LayerNone,
	}

	if !ColliderOverlap(player, enemy) {
		t.Fatal("player (mask includes enemy) should collide with enemy")
	}
	if !ColliderOverlap(enemy, player) {
		t.Fatal("enemy (mask includes player) should collide with player")
	}
	if !ColliderOverlap(player, wall) {
		t.Fatal("player (mask includes wall) should collide with wall")
	}
	if ColliderOverlap(enemy, wall) {
		t.Fatal("enemy (mask excludes wall) should NOT collide with wall")
	}
}

func TestColliderNoSpatialOverlap(t *testing.T) {
	const LayerA Layer = 1 << 0
	a := Collider{Bounds: Rect{X: 0, Y: 0, W: 5, H: 5}, Layer: LayerA, Mask: LayerA}
	b := Collider{Bounds: Rect{X: 20, Y: 20, W: 5, H: 5}, Layer: LayerA, Mask: LayerA}
	if ColliderOverlap(a, b) {
		t.Fatal("spatially separated colliders should not overlap")
	}
}
