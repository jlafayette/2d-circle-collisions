// Code generated by file2byteslice. DO NOT EDIT.
// (gofmt is fine after generating)

package shader

var Circle = []byte("// +build ignore\r\n\r\npackage shader\r\n\r\nvar Translate vec2\r\nvar Size vec2\r\n\r\nfunc Fragment(position vec4, texCoord vec2, color vec4) vec4 {\r\n\r\n\t// vec2 st = gl_FragCoord.xy/u_resolution;\r\n\tst := (position.xy / Size) - (Translate / Size)\r\n\r\n\t// The DISTANCE from the pixel to the center\r\n\tdist := distance(st, vec2(0.5))\r\n\r\n\t// Make circle edge 1 pixel away from the edges of the image to avoid\r\n\t// clipping\r\n\tpixel := 1.0 / Size.x\r\n\tradius := 0.5 - pixel\r\n\r\n\t// falloff of 1 pixel\r\n\tfalloff := pixel\r\n\r\n\tclr := smoothstep(radius, radius-falloff, dist)\r\n\r\n\t// Uncomment this line to create an outline instead of a filled circle\r\n\t// clr = clr * smoothstep(radius-(falloff*2.0), radius-falloff, dist)\r\n\r\n\treturn vec4(clr)\r\n}\r\n")
