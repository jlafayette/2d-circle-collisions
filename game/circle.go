package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lucasb-eyer/go-colorful"
)

// Bresenham algorithm for rasterizing a circle
// Draw a circle that fills given image
// Before drawing, ensure that the image has odd width and height for best
// results:
//   if width%2 == 0 {
//     width++
//     height = width
//   }
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

// Use shader to draw a circle
func drawCircleToImage(img *ebiten.Image, shader *ebiten.Shader) {
	w, h := img.Size()
	x := 0.0
	y := 0.0
	op := &ebiten.DrawRectShaderOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.Uniforms = map[string]interface{}{
		"Translate": []float32{float32(x), float32(y)},
		"Size":      []float32{float32(w), float32(h)},
	}
	img.DrawRectShader(w, h, shader, op)
}

// NewCircle creates a new circle at position x,y with radius r
func NewCircle(x, y, r float64, shader *ebiten.Shader) *Circle {

	var width = int(r)*2 + 3
	var height = width

	img := ebiten.NewImage(width, height)

	drawCircleToImage(img, shader)

	// mod controls the accumulation of activity based on speed
	maxMod := remap(r, 5, 70, 5, 2)

	// dim rate controls how fast activity fades
	dimRate := remap(r, 5, 70, 0.07, 0.01)

	// max change is the max that activity can reach
	maxCharge := 1.5

	return &Circle{
		selected:  false,
		pos:       Vec2{x, y},
		radius:    r,
		area:      math.Pi * r * r,
		maxMod:    maxMod,
		dimRate:   dimRate,
		maxCharge: maxCharge,
		color:     randomCircleColor(),
		image:     img,
	}
}

func randomCircleColor() colorful.Color {
	hue := randFloat(0, 360)
	if hue > 360 {
		hue -= 360
	}
	chroma := 0.45    // -1 .. 1
	lightness := 0.45 // 0 .. 1
	return colorful.Hcl(hue, chroma, lightness)
}

// Circle represents a circle
type Circle struct {
	selected bool
	id       int
	pos      Vec2
	prevPos  Vec2
	vel      Vec2
	acc      Vec2
	radius   float64
	area     float64
	speed    float64

	activity  float64
	maxMod    float64
	dimRate   float64
	maxCharge float64

	color colorful.Color
	image *ebiten.Image
}

func (c *Circle) postUpdate() {
	c.speed = c.vel.Len()
	mod := remap(c.speed, 0, 100, 0, c.maxMod)
	c.activity += mod
	c.activity -= c.dimRate
	c.activity = math.Min(math.Max(c.activity, 0), c.maxCharge)
}

func (c *Circle) addCollisionEnergy(energy float64) {
	mod := remap(energy, 0, 500, 0, c.maxMod)
	c.activity += mod
}

// Draw the circle to the screen.
func (c Circle) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// set chroma and lightness based on speed
	if c.selected {
		c.activity = math.Max(c.activity, 1.0)
	}
	hue, chroma, lightness := c.color.Hcl()
	chroma = remap(math.Min(c.activity, 1), 0, 1, 0, 1)
	lightness = remap(math.Min(c.activity, 1), 0, 1, 0.45, 0.9)
	c.color = colorful.Hcl(hue, chroma, lightness)
	r := c.color.R
	g := c.color.G
	b := c.color.B
	if c.selected {
		h, s, v := c.color.Hsv()
		col := colorful.Hsv(h, s, math.Min(v+0.25, 1))
		r = col.R
		g = col.G
		b = col.B
	}

	// Draw motion blur effect that fades as the circle slows
	if c.speed > 10 {
		a := remap(clamp(c.speed, 10, 75), 10, 75, 0, 0.95)
		op.GeoM.Translate(c.prevPos.X-c.radius, c.prevPos.Y-c.radius)
		op.ColorM.Scale(r, g, b, a)
		screen.DrawImage(c.image, op)
		drawLine(c.pos, c.prevPos, c.radius*1.9, screen, c.color, a)
	}

	// Draw the circle
	op.GeoM.Reset()
	op.ColorM.Reset()
	op.ColorM.Scale(r, g, b, 1)
	op.GeoM.Translate(c.pos.X-c.radius, c.pos.Y-c.radius)
	screen.DrawImage(c.image, op)

}
