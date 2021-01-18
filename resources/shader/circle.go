// +build ignore

package shader

var Translate vec2
var Size vec2

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {

	// vec2 st = gl_FragCoord.xy/u_resolution;
	st := (position.xy / Size) - (Translate / Size)

	// The DISTANCE from the pixel to the center
	dist := distance(st, vec2(0.5))

	// Make circle edge 1 pixel away from the edges of the image to avoid
	// clipping
	pixel := 1.0 / Size.x
	radius := 0.5 - pixel

	// falloff of 1 pixel
	falloff := pixel

	clr := smoothstep(radius, radius-falloff, dist)

	// Uncomment this line to create an outline instead of a filled circle
	// clr = clr * smoothstep(radius-(falloff*2.0), radius-falloff, dist)

	return vec4(clr)
}
