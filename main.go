package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jlafayette/2d-circle-collisions/game"
)

var (
	// Screen width and height indicates how many pixels to draw, not the
	// window dimensions.
	screenWidth  = 1880 // 1036
	screenHeight = 1040 // 640
)

func main() {

	// In this test, window size is equal to screen size, so no pixelation
	// or stretching will occur.
	ebiten.SetWindowSize(screenWidth, screenHeight)

	ebiten.SetWindowTitle("2D Collisions")
	game := game.NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
