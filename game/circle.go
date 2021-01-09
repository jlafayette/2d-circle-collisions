package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// var (
// 	selectedColor = color.RGBA{0, 100, 255, 255}
// 	collidedColor = color.RGBA{255, 0, 0, 255}
// 	defaultColor  = color.RGBA{200, 200, 200, 255}
// )

// Bresenham algorithm for rasterizing a circle
// Draw a circle that fills given image
func bresenham(color color.Color, img *ebiten.Image) {
	width, height := img.Size()
	x := width / 2
	y := height / 2
	r := width / 2

	if r < 0 {
		return
	}

	x1, y1, err := -r, 0, 2-2*r
	for {
		img.Set(x-x1, y+y1, color)
		img.Set(x-y1, y-x1, color)
		img.Set(x+x1, y-y1, color)
		img.Set(x+y1, y+x1, color)
		r = err
		if r > x1 {
			x1++
			err += x1*2 + 1
		}
		if r <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}
}

// NewCircle creates a new circle at position x,y with radius r
func NewCircle(x, y, r float64, color color.Color) *Circle {

	var width = int(r) * 2
	var height = width
	if width%2 == 0 {
		width++
		height = width
	}
	img := ebiten.NewImage(width, height)
	bresenham(color, img)

	return &Circle{
		selected: false,
		posX:     x,
		posY:     y,
		radius:   r,
		image:    img,
	}
}

// Circle represents a circle
type Circle struct {
	selected bool
	posX     float64
	posY     float64
	velX     float64
	velY     float64
	accX     float64
	accY     float64
	radius   float64
	image    *ebiten.Image
}

// Draw the circle to the screen.
func (c Circle) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.posX-c.radius, c.posY-c.radius)

	if c.selected {
		op.ColorM.Scale(0, 0.5, 1, 1)
	}

	screen.DrawImage(c.image, op)
}
