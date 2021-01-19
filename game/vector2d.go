package game

import "math"

// Vec2 is a 2D vector
type Vec2 struct {
	X float64
	Y float64
}

// Add two vectors
func (u Vec2) Add(v Vec2) Vec2 {
	return Vec2{
		X: u.X + v.X,
		Y: u.Y + v.Y,
	}
}

// Sub subtracts vector v from u
func (u Vec2) Sub(v Vec2) Vec2 {
	return Vec2{
		X: u.X - v.X,
		Y: u.Y - v.Y,
	}
}

// To returns vector from u to v (same as v.Sub(u))
func (u Vec2) To(v Vec2) Vec2 {
	return Vec2{
		X: v.X - u.X,
		Y: v.Y - u.Y,
	}
}

// Scaled returns vector u multiplied by s
func (u Vec2) Scaled(s float64) Vec2 {
	return Vec2{
		X: u.X * s,
		Y: u.Y * s,
	}
}

// Unit returns a unit vector with len 1 in the direction of u
func (u Vec2) Unit() Vec2 {
	return u.Scaled(1.0 / u.Len())
}

// Dot returns the dot product of vectors u and v
func (u Vec2) Dot(v Vec2) float64 {
	return u.X*v.X + u.Y*v.Y
}

// Len returns the length of vector u
func (u Vec2) Len() float64 {
	return math.Hypot(u.X, u.Y)
}

// Cross returns the cross product of vectors u and v
func (u Vec2) Cross(v Vec2) float64 {
	return u.X*v.Y - v.X*u.Y
}

// Normal returns a vector normal to u
func (u Vec2) Normal() Vec2 {
	return Vec2{-u.Y, u.X}
}
