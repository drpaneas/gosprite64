package main

import "testing"

func TestBuildProjectedTriangleChangesWithAngle(t *testing.T) {
	a := buildProjectedTriangle(0)
	b := buildProjectedTriangle(45)

	same := true
	for i := range a {
		if a[i].X != b[i].X || a[i].Y != b[i].Y {
			same = false
			break
		}
	}
	if same {
		t.Fatal("projected triangle should change when angle changes")
	}
}

func TestBuildProjectedTriangleUsesPerspective(t *testing.T) {
	verts := buildProjectedTriangle(0)
	for i, v := range verts {
		if v.InvW <= 0 {
			t.Fatalf("vertex %d InvW = %f, want > 0", i, v.InvW)
		}
		if v.X < 0 || v.X > screenW {
			t.Fatalf("vertex %d X = %f out of screen bounds", i, v.X)
		}
		if v.Y < 0 || v.Y > screenH {
			t.Fatalf("vertex %d Y = %f out of screen bounds", i, v.Y)
		}
	}
}
