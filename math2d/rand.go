package math2d

// Rand is a fast, seedable pseudo-random number generator using xoshiro128**.
// Suitable for deterministic gameplay - same seed always produces same sequence.
type Rand struct {
	s [4]uint32
}

// NewRand creates a new PRNG seeded with the given value.
func NewRand(seed uint64) *Rand {
	r := &Rand{}
	r.Seed(seed)
	return r
}

// Seed resets the generator state. Uses SplitMix64 to expand the seed.
func (r *Rand) Seed(seed uint64) {
	r.s[0] = uint32(splitmix64(&seed))
	r.s[1] = uint32(splitmix64(&seed))
	r.s[2] = uint32(splitmix64(&seed))
	r.s[3] = uint32(splitmix64(&seed))
	if r.s[0]|r.s[1]|r.s[2]|r.s[3] == 0 {
		r.s[0] = 1
	}
}

// Uint32 returns a pseudo-random uint32.
func (r *Rand) Uint32() uint32 {
	result := rotl32(r.s[1]*5, 7) * 9
	t := r.s[1] << 9
	r.s[2] ^= r.s[0]
	r.s[3] ^= r.s[1]
	r.s[1] ^= r.s[2]
	r.s[0] ^= r.s[3]
	r.s[2] ^= t
	r.s[3] = rotl32(r.s[3], 11)
	return result
}

// Intn returns a pseudo-random int in [0, n). Panics if n <= 0.
func (r *Rand) Intn(n int) int {
	if n <= 0 {
		panic("math2d: Intn argument must be positive")
	}
	return int(r.Uint32() % uint32(n))
}

// Float32 returns a pseudo-random float32 in [0.0, 1.0).
func (r *Rand) Float32() float32 {
	return float32(r.Uint32()>>8) / (1 << 24)
}

// RangeInt returns a pseudo-random int in [min, max).
func (r *Rand) RangeInt(min, max int) int {
	if max <= min {
		return min
	}
	return min + r.Intn(max-min)
}

// RangeFloat32 returns a pseudo-random float32 in [min, max).
func (r *Rand) RangeFloat32(min, max float32) float32 {
	return min + r.Float32()*(max-min)
}

// Bool returns true or false with roughly equal probability.
func (r *Rand) Bool() bool {
	return r.Uint32()&1 == 0
}

func rotl32(x uint32, k int) uint32 {
	return (x << k) | (x >> (32 - k))
}

func splitmix64(state *uint64) uint64 {
	*state += 0x9e3779b97f4a7c15
	z := *state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}
