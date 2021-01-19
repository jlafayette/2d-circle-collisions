package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// NewCapsule creates a new line from (x1, y1) to (x2, y2)
func NewCapsule(x1, y1, x2, y2, r float64, shader *ebiten.Shader) *Capsule {
	width := int(r)*2 + 3
	height := width

	img := ebiten.NewImage(width, height)

	drawCircleToImage(img, shader)
	return &Capsule{
		start:  Vec2{x1, y1},
		end:    Vec2{x2, y2},
		radius: r,
		image:  img,
	}
}

// Capsule represents a line that can collide with circles
type Capsule struct {
	start  Vec2
	end    Vec2
	radius float64
	image  *ebiten.Image
}

// Draw the line to the screen.
func (c *Capsule) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(0.5, 0.5, 0.5, 1)
	op.GeoM.Translate(c.start.X-c.radius, c.start.Y-c.radius)
	screen.DrawImage(c.image, op)
	op.GeoM.Reset()
	op.GeoM.Translate(c.end.X-c.radius, c.end.Y-c.radius)
	screen.DrawImage(c.image, op)

	// ebitenutil.DrawLine(screen, c.start.X, c.start.Y, c.end.X, c.end.Y, color.RGBA{255, 0, 0, 255})

	lineV := c.start.To(c.end)
	offset := lineV.Normal().Unit().Scaled(c.radius)

	start := c.start.Add(offset)
	end := c.end.Add(offset)
	ebitenutil.DrawLine(screen, start.X, start.Y, end.X, end.Y, color.White)

	start = c.start.Sub(offset)
	end = c.end.Sub(offset)
	ebitenutil.DrawLine(screen, start.X, start.Y, end.X, end.Y, color.White)
}
