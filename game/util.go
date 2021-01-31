package game

import (
	"math"
	"math/rand"
)

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func remap(in, inMin, inMax, outMin, outMax float64) float64 {
	return (in-inMin)/(inMax-inMin)*(outMax-outMin) + outMin
}

func clamp(in, min, max float64) float64 {
	return math.Min(max, math.Max(min, in))
}

func randRadius(min, max float64) float64 {
	x := rand.Float64()
	y := shape(x)
	return min + y*(max-min)
}

func shape(x float64) float64 {
	// experiment with shaping functions
	// https://www.iquilezles.org/www/articles/functions/functions.htm

	// a := 0.1
	// b := 2.0
	// k := math.Pow(a+b, a+b) / (math.Pow(a, a) * math.Pow(b, b))
	// y := k * math.Pow(x, a) * math.Pow(1.0-x, b)
	// return y

	return x
}
