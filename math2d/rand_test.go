package math2d

import "testing"

func TestRandDeterministic(t *testing.T) {
	r1 := NewRand(42)
	r2 := NewRand(42)
	for i := 0; i < 100; i++ {
		a := r1.Uint32()
		b := r2.Uint32()
		if a != b {
			t.Fatalf("iteration %d: same seed should produce same output: %d != %d", i, a, b)
		}
	}
}

func TestRandDifferentSeeds(t *testing.T) {
	r1 := NewRand(1)
	r2 := NewRand(2)
	same := 0
	for i := 0; i < 100; i++ {
		if r1.Uint32() == r2.Uint32() {
			same++
		}
	}
	if same > 5 {
		t.Fatalf("different seeds produced %d/100 identical values", same)
	}
}

func TestRandIntn(t *testing.T) {
	r := NewRand(99)
	for i := 0; i < 1000; i++ {
		v := r.Intn(10)
		if v < 0 || v >= 10 {
			t.Fatalf("Intn(10) = %d, out of range", v)
		}
	}
}

func TestRandIntnOne(t *testing.T) {
	r := NewRand(1)
	for i := 0; i < 100; i++ {
		v := r.Intn(1)
		if v != 0 {
			t.Fatalf("Intn(1) should always return 0, got %d", v)
		}
	}
}

func TestRandFloat32(t *testing.T) {
	r := NewRand(77)
	for i := 0; i < 1000; i++ {
		v := r.Float32()
		if v < 0 || v >= 1 {
			t.Fatalf("Float32() = %f, out of [0,1) range", v)
		}
	}
}

func TestRandRangeInt(t *testing.T) {
	r := NewRand(55)
	for i := 0; i < 1000; i++ {
		v := r.RangeInt(5, 10)
		if v < 5 || v >= 10 {
			t.Fatalf("RangeInt(5, 10) = %d, out of range", v)
		}
	}
}

func TestRandRangeFloat32(t *testing.T) {
	r := NewRand(88)
	for i := 0; i < 1000; i++ {
		v := r.RangeFloat32(2.0, 5.0)
		if v < 2.0 || v >= 5.0 {
			t.Fatalf("RangeFloat32(2, 5) = %f, out of range", v)
		}
	}
}

func TestRandBool(t *testing.T) {
	r := NewRand(123)
	trues := 0
	for i := 0; i < 10000; i++ {
		if r.Bool() {
			trues++
		}
	}
	if trues < 4000 || trues > 6000 {
		t.Fatalf("Bool() distribution: %d/10000 trues, expected ~5000", trues)
	}
}

func TestRandSeed(t *testing.T) {
	r := NewRand(42)
	first := r.Uint32()
	r.Seed(42)
	second := r.Uint32()
	if first != second {
		t.Fatalf("re-seeding should reset state: %d != %d", first, second)
	}
}
