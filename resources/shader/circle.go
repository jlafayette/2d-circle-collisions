// +build ignore

package shader

var Translate vec2
var Size vec2

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {

	// vec2 st = gl_FragCoord.xy/u_resolution;
	st := (position.xy / Size) - (Translate / Size)

	pixel := 1.0 / Size.x

	// The DISTANCE from the pixel to the center
	d := distance(st, vec2(0.5))

	r := 0.5 - pixel

	// // falloff
	// // 0.01 for large (r>=100.0)
	// // 0.1 for small (r=5.0)
	// low1 := 5.0
	// high1 := 100.0
	// low2 := 0.9
	// high2 := 0.99
	// f := clamp(Size.x, low1, high1)
	// f = low2 + (f-low1)*(high2-low2)/(high1-low1)
	// f = 1 - f
	f := pixel

	c := smoothstep(r, r-f, d)
	return vec4(c)
}
