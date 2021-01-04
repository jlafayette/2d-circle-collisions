package game

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	twoPi  = math.Pi * 2
	halfPi = math.Pi / 2
)

// Game implements ebiten.Game interface and stores the game state.
//
// The methods run in the following order (each one is run once in this order
// if fps is 60 and display is 60 Hz):
//	Update
//	Draw
//	Layout
type Game struct {
	width   int
	height  int
	circles []*Circle
}

// NewGame creates a new Game
func NewGame(width, height int) *Game {

	seed := time.Now().UnixNano()
	rand.Seed(seed)

	var circles []*Circle

	// reference circles
	circles = append(circles, NewCircle(0.0, 0.0, 25.0, color.White))
	circles = append(circles, NewCircle(float64(width), float64(height), 25.0, color.White))

	return &Game{
		width:   width,
		height:  height,
		circles: circles,
	}
}

// Update function is called every tick and updates the game's logical state.
func (g *Game) Update() error {

	for i := 0; len(g.circles) < 500 && i < 5; i++ {
		xpos := randFloat(100, float64(g.width)-100)
		ypos := randFloat(100, float64(g.height)-100)
		radius := randFloat(100, 100)
		circle := NewCircle(xpos, ypos, radius, color.White)
		g.circles = append(g.circles, circle)
	}

	return nil
}

// Draw is called every frame. The frame frequency depends on the display's
// refresh rate, so if the display is 60 Hz, Draw will be called 60 times per
// second.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for i := range g.circles {
		g.circles[i].Draw(screen)
	}
}

// Layout accepts the window size on desktop as the outside size, and return's
// the game's internal or pixel screen size, which is then scaled up to fit in
// the outside size. This does more for resizeable windows.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}
