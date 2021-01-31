package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lucasb-eyer/go-colorful"
)

var rect1x1 *ebiten.Image

func init() {
	rect1x1 = ebiten.NewImage(1, 1)
	rect1x1.Fill(color.White)
}

func drawLine(start, end Vec2, thickness float64, target *ebiten.Image, color colorful.Color) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(color.R, color.G, color.B, 1.0)

	lineV := start.To(end)
	unitOffset := lineV.Normal().Unit()
	offset := unitOffset.Scaled(thickness / 2.0)

	// Use rectangle image to draw line by rotating, scaling, and moving
	op.GeoM.Scale(lineV.Len(), thickness)
	op.GeoM.Rotate(lineV.Angle())
	op.GeoM.Translate(start.X-offset.X, start.Y-offset.Y)
	target.DrawImage(rect1x1, op)
}
