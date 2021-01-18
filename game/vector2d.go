package game

// Vec2 is a 2D vector
type Vec2 struct {
	X float64
	Y float64
}

// Add two vectors
func (a Vec2) Add(b Vec2) Vec2 {
	return Vec2{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

// Sub subtracts vec b from a
func (a Vec2) Sub(b Vec2) Vec2 {
	return Vec2{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}
