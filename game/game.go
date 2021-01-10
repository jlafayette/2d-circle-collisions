package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	width             int
	height            int
	engine            *Engine
	updateElapsedTime time.Duration
	drawElapsedTime   time.Duration
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
		width:  width,
		height: height,
		engine: NewEngine(circles),
	}
}

// Update function is called every tick and updates the game's logical state.
func (g *Game) Update() error {

	start := time.Now()

	mx, my := ebiten.CursorPosition()
	mxf := float64(mx)
	myf := float64(my)

	// If a circle is under the mouse curser, then select it
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.engine.selectAtPostion(mxf, myf)
	}

	// Handle dragging selected circle
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.engine.moveSelectedTo(mxf, myf)
	}

	// Clear selection
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.engine.deselect()
	}

	for i := 0; len(g.engine.circles) < 600 && i < 2; i++ {
		xbuffer := float64(g.width / 4)
		ybuffer := float64(g.height / 4)
		xpos := randFloat(xbuffer, float64(g.width)-xbuffer)
		ypos := randFloat(ybuffer, float64(g.height)-ybuffer)
		radius := randFloat(5, 50)
		circle := NewCircle(xpos, ypos, radius, color.White)
		g.engine.circles = append(g.engine.circles, circle)
	}

	// TODO: get proper elapsed time
	g.engine.update(g.width, g.height, 1.0)

	g.updateElapsedTime = time.Now().Sub(start)

	return nil
}

// Draw is called every frame. The frame frequency depends on the display's
// refresh rate, so if the display is 60 Hz, Draw will be called 60 times per
// second.
func (g *Game) Draw(screen *ebiten.Image) {
	start := time.Now()

	screen.Fill(color.Black)
	for i := range g.engine.circles {
		g.engine.circles[i].Draw(screen)
	}

	// Debug text
	msg := fmt.Sprintf(
		"FPS: %0.2f\nTPS: %0.2f\nUpdate Elapsed: %0.4f\nDraw Elapsed: %0.4f",
		ebiten.CurrentFPS(),
		ebiten.CurrentTPS(),
		g.updateElapsedTime.Seconds(),
		g.drawElapsedTime.Seconds(),
	)
	ebitenutil.DebugPrint(screen, msg)

	g.drawElapsedTime = time.Now().Sub(start)
}

// Layout accepts the window size on desktop as the outside size, and return's
// the game's internal or pixel screen size, which is then scaled up to fit in
// the outside size. This does more for resizeable windows.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}
