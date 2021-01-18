package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// NewCapsule creates a new line from (x1, y1) to (x2, y2)
func NewCapsule(x1, y1, x2, y2, r float64) *Capsule {
	return &Capsule{
		x1:     x1,
		y1:     y1,
		x2:     x2,
		y2:     y2,
		radius: r,
	}
}

// Capsule represents a line that can collide with circles
type Capsule struct {
	x1     float64
	y1     float64
	x2     float64
	y2     float64
	radius float64
}

// Draw the line to the screen.
func (c *Capsule) Draw(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, c.x1, c.y1, c.x2, c.y2, color.White)

	// normal vec
	nx := -(c.y2 - c.y1)
	ny := c.x2 - c.x1

	// len of normal vec
	d := math.Sqrt(nx*nx + ny*ny)

	// normalized normal vec
	nx = nx / d
	ny = nx / d

	x1 := c.x1 + nx*c.radius
	y1 := c.y1 + ny*c.radius
	x2 := c.x2 + nx*c.radius
	y2 := c.y2 + ny*c.radius
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.White)

	x1 = c.x1 - nx*c.radius
	y1 = c.y1 - ny*c.radius
	x2 = c.x2 - nx*c.radius
	y2 = c.y2 - ny*c.radius
	ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.White)
}
