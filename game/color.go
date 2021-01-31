package game

import "github.com/lucasb-eyer/go-colorful"

func contrastColor(clr colorful.Color) colorful.Color {
	h, _, _ := clr.Hcl()
	h += 180
	if h > 360 {
		h -= 360
	}
	return colorful.Hcl(h, 1.0, 0.75)
}
