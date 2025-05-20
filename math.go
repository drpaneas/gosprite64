package gosprite64

import (
	"math"
	"math/rand/v2"
)

// Number is a constraint that permits any numeric type.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Flr returns the floor of a generic number.
// It works with any numeric type.
//
// Example:
//
//	Flr(3.14) // returns 3
//	Flr(-2.71) // returns -3
func Flr[T Number](a T) int {
	floatVal := float64(a)
	floorVal := math.Floor(floatVal)
	return int(floorVal)
}

func Rnd[T Number](a T) int {
	limit := float64(a)

	if limit <= 0 {
		return 0
	}

	// rand.Float64() returns a float64 in [0.0, 1.0)
	// Multiplying by limit gives a float64 in [0.0, limit)
	// Applying Floor and converting to int gives an integer in [0, floor(limit))
	return int(math.Floor(rand.Float64() * limit))
}
