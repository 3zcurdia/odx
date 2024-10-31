package vectors

import (
	"testing"
)

func TestCross(t *testing.T) {
	a := Init([]float64{1, 0, 0})
	b := Init([]float64{0, 1, 0})
	result := a.Cross(b)
	expected := Init([]float64{0, 0, 1})
	if !result.Equals(expected) {
		t.Errorf("Cross product of %v and %v is not %v", a, b, Init([]float64{0, 0, 1}))
	}
}

func TestDot(t *testing.T) {
	a := Init([]float64{1, 0, 0})
	b := Init([]float64{0, 1, 0})
	if a.Dot(b) != 0 {
		t.Errorf("Dot product of %v and %v is not 0", a, b)
	}
}

func TestLength(t *testing.T) {
	a := Init([]float64{1, 0, 0})
	if a.Length() != 1 {
		t.Errorf("Length of %v is not 1", a)
	}
}

func TestNormalize(t *testing.T) {
	a := Init([]float64{9, 0, 0})
	result := a.Normalize()
	expected := Init([]float64{1, 0, 0})
	if !result.Equals(expected) {
		t.Errorf("Normalized vector of %v is not %v (got %v)", a, expected, result)
	}
}
