package game

import (
	"image/color"
	"math"
	"math/rand"
	"strconv"
	"strings"
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
	showFPS           bool
	showDebug         bool
	speedControl      *SpeedControl
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
	circles = append(circles, NewCircle(float64(width)/2, float64(height)/2, 200.0, color.White))

	return &Game{
		width:        width,
		height:       height,
		showFPS:      true,
		showDebug:    true,
		speedControl: NewSpeedControl(),
		engine:       NewEngine(circles),
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
		g.engine.selectNearestPostion(mxf, myf)
	}
	// Handle pulling selected circle
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.engine.applyForceToSelected(mxf, myf, g.speedControl.multiplier())
	}
	// Clear selection
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.engine.deselect()
	}

	// Dynamic input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.engine.dynamicNearestPosition(mxf, myf)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		g.engine.dynamicRelease(mxf, myf)
	}

	// Toggle display of FPS and debug text/lines
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.showDebug = !g.showDebug
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.showFPS = !g.showFPS
	}

	// Handle speed control keyboard inputs
	g.speedControl.update()

	if !g.speedControl.paused() {
		// larger
		max := 100
		for i := 0; len(g.engine.circles) < max && i < 1; i++ {
			xbuffer := float64(g.width / 4)
			ybuffer := float64(g.height / 4)
			xpos := randFloat(xbuffer, float64(g.width)-xbuffer)
			ypos := randFloat(ybuffer, float64(g.height)-ybuffer)
			radius := randRadius(10, 70)
			circle := NewCircle(xpos, ypos, radius, color.White)
			g.engine.circles = append(g.engine.circles, circle)
		}
		// smaller
		for i := 0; len(g.engine.circles) < max && i < 3; i++ {
			xbuffer := float64(g.width / 4)
			ybuffer := float64(g.height / 4)
			xpos := randFloat(xbuffer, float64(g.width)-xbuffer)
			ypos := randFloat(ybuffer, float64(g.height)-ybuffer)
			radius := randRadius(5, 35)
			circle := NewCircle(xpos, ypos, radius, color.White)
			g.engine.circles = append(g.engine.circles, circle)
		}
	}

	// TODO: get proper elapsed time
	elapsedTime := 1.0
	g.engine.update(g.width, g.height, g.speedControl.multiplier(), elapsedTime)

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

	// Draw dynamic input line
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		mx, my := ebiten.CursorPosition()
		mxf := float64(mx)
		myf := float64(my)
		x2, y2, found := g.engine.getDynamicPosition(mxf, myf)
		if found {
			ebitenutil.DrawLine(
				screen,
				mxf, myf, x2, y2,
				color.RGBA{0, 255, 0, 255},
			)
		}
	}
	// Draw selected pull line
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		mxf := float64(mx)
		myf := float64(my)
		x2, y2, found := g.engine.getSelectedPosition(mxf, myf)
		if found {
			ebitenutil.DrawLine(
				screen,
				mxf, myf, x2, y2,
				color.RGBA{0, 255, 0, 255},
			)
		}
	}

	// Debug text and lines
	if g.showFPS || g.showDebug {
		var msg strings.Builder
		if g.showFPS {
			msg.WriteString("FPS: ")
			msg.WriteString(strconv.FormatFloat(ebiten.CurrentFPS(), 'f', 2, 64))
			msg.WriteString("\nTPS: ")
			msg.WriteString(strconv.FormatFloat(ebiten.CurrentTPS(), 'f', 2, 64))
			msg.WriteString("\n")
		}
		if g.showDebug {
			msg.WriteString("Game speed: ")
			msg.WriteString(strconv.Itoa(g.speedControl.control))
			msg.WriteString("\nCircle count: ")
			msg.WriteString(strconv.Itoa(len(g.engine.circles)))
			msg.WriteString("\nChecks: ")
			msg.WriteString(strconv.Itoa(g.engine.checks))
			// msg.WriteString("\nUpdate Elapsed: ")
			// msg.WriteString(strconv.FormatFloat(g.updateElapsedTime.Seconds(), 'f', 4, 64))
			// msg.WriteString("\nDraw Elapsed: ")
			// msg.WriteString(strconv.FormatFloat(g.drawElapsedTime.Seconds(), 'f', 4, 64))

			// Draw red lines between colliding circles
			for _, p := range g.engine.collidingPairs {
				ebitenutil.DrawLine(
					screen,
					g.engine.circles[p.a].posX,
					g.engine.circles[p.a].posY,
					g.engine.circles[p.b].posX,
					g.engine.circles[p.b].posY,
					color.RGBA{255, 0, 0, 30},
				)
			}
		}
		ebitenutil.DebugPrint(screen, msg.String())
	}

	// g.drawShapeFunction(screen)

	g.drawElapsedTime = time.Now().Sub(start)
}

func (g *Game) drawShapeFunction(screen *ebiten.Image) {
	// For visualizing shaping functions (see utils/shape)
	for ix := 0; ix < g.height; ix++ {
		// map x 0..1
		x := float64(ix) / float64(g.height)

		// y in terms of x
		y := shape(x)

		// translate x and y from 0..1 to 0..screen-height
		x = x * float64(g.height)
		y = y * float64(g.height)
		ebitenutil.DrawRect(screen, x, float64(g.height)-y, 2, 2, color.White)
	}
	ebitenutil.DrawRect(screen, float64(g.height), 0, 2, float64(g.height), color.White)
}

// Layout accepts the window size on desktop as the outside size, and return's
// the game's internal or pixel screen size, which is then scaled up to fit in
// the outside size. This does more for resizeable windows.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}
