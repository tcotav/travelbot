package location

import (
	"math"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v *Point) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Point) Add(other Point) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Point) Multiply(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Point) Normalize() {
	magnitude := math.Sqrt(v.X*v.X + v.Y*v.Y)
	v.X /= magnitude
	v.Y /= magnitude
}

func (v *Point) DistanceTo(v2 Point) float64 {
	return math.Sqrt(math.Pow(v2.X-v.X, 2) +
		math.Pow(v2.Y-v.Y, 2))
}
