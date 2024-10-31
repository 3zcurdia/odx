package vectors

import "math"

type Vector struct {
	x float64
	y float64
	z float64
}

func Init(v []float64) Vector {
	return Vector{
		x: v[0],
		y: v[1],
		z: v[2],
	}
}

func (a *Vector) Cross(b Vector) Vector {
	return Vector{
		x: (a.y * b.z) - (a.z * b.y),
		y: (a.z * b.x) - (a.x * b.z),
		z: (a.x * b.y) - (a.y * b.x),
	}
}

func (a *Vector) Dot(b Vector) float64 {
	return (a.x * b.x) + (a.y * b.y) + (a.z * b.z)
}

func (a *Vector) Equals(b Vector) bool {
	return a.x == b.x && a.y == b.y && a.z == b.z
}

func (a *Vector) Length() float64 {
	return math.Sqrt(a.x*a.x + a.y*a.y + a.z*a.z)
}

func (a *Vector) Normalize() Vector {
	length := a.Length()
	return Vector{
		x: a.x / length,
		y: a.y / length,
		z: a.z / length,
	}
}
