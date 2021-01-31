package game

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jlafayette/2d-circle-collisions/resources/shader"
	"github.com/lucasb-eyer/go-colorful"
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
	time              int
	showFPS           bool
	showDebug         bool
	speedControl      *SpeedControl
	engine            *Engine
	circleShader      *ebiten.Shader
	updateElapsedTime time.Duration
	drawElapsedTime   time.Duration
}

// NewGame creates a new Game
func NewGame(width, height int) *Game {

	seed := time.Now().UnixNano()
	rand.Seed(seed)

	sh, err := ebiten.NewShader(shader.Circle)
	if err != nil {
		log.Fatal("Circle shader failed: ", err)
	}

	var circles []*Circle
	// circles = append(circles, NewCircle(float64(width)/2, float64(height)/2, 200.0, color.White, sh))

	var capsules []*Capsule
	w := float64(width) - 5
	h := float64(height) - 5
	capsules = append(capsules, NewCapsule(5, 5, w, 5, 10, sh))
	capsules = append(capsules, NewCapsule(5, 5, 5, h, 10, sh))
	capsules = append(capsules, NewCapsule(w, h, w, 5, 10, sh))
	capsules = append(capsules, NewCapsule(w, h, 5, h, 10, sh))
	capsules = append(capsules, NewCapsule(100, 100, 500, 500, 10, sh))

	return &Game{
		width:        width,
		height:       height,
		showFPS:      true,
		showDebug:    true,
		speedControl: NewSpeedControl(),
		engine:       NewEngine(circles, capsules),
		circleShader: sh,
	}
}

func cursorPosition() Vec2 {
	x, y := ebiten.CursorPosition()
	return Vec2{
		X: float64(x),
		Y: float64(y),
	}
}

// Update function is called every tick and updates the game's logical state.
func (g *Game) Update() error {
	g.time++

	start := time.Now()

	cursorPos := cursorPosition()

	// If cursor is over capsule end, then drag it around
	// Otherwise pull the nearest circle
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		found := g.engine.selectCapsuleAtPostion(cursorPos)
		if !found {
			g.engine.selectNearestPostion(cursorPos)
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		found := g.engine.moveSelectedCapsuleTo(cursorPos)
		if !found {
			g.engine.applyForceToSelected(cursorPos, g.speedControl.multiplier())
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.engine.deselectCapsule()
		g.engine.deselect()
	}

	// Dynamic input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.engine.dynamicNearestPosition(cursorPos)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		g.engine.dynamicRelease(cursorPos)
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
		max := 10
		for i := 0; len(g.engine.circles) < max && i < 1; i++ {
			xbuffer := float64(g.width / 4)
			ybuffer := float64(g.height / 4)
			xpos := randFloat(xbuffer, float64(g.width)-xbuffer)
			ypos := randFloat(ybuffer, float64(g.height)-ybuffer)
			radius := randRadius(10, 70)
			circle := NewCircle(xpos, ypos, radius, g.circleShader)
			g.engine.addCircle(circle)
		}
		// smaller
		for i := 0; len(g.engine.circles) < max && i < 3; i++ {
			xbuffer := float64(g.width / 4)
			ybuffer := float64(g.height / 4)
			xpos := randFloat(xbuffer, float64(g.width)-xbuffer)
			ypos := randFloat(ybuffer, float64(g.height)-ybuffer)
			radius := randRadius(5, 35)
			circle := NewCircle(xpos, ypos, radius, g.circleShader)
			g.engine.addCircle(circle)
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

	cursorPos := cursorPosition()

	screen.Fill(color.Black)
	for i := range g.engine.circles {
		g.engine.circles[i].Draw(screen)
	}
	for i := range g.engine.capsules {
		g.engine.capsules[i].Draw(screen)
	}

	// Draw dynamic input line
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		circle := g.engine.getDynamic()
		if circle != nil {
			// opposite hue
			h, _, _ := circle.color.Hcl()
			h += 180
			if h > 360 {
				h -= 360
			}
			clr := colorful.Hcl(h, 1.0, 0.75)
			drawLine(cursorPos, circle.pos, 2, screen, clr)
		}
	}

	// Draw selected pull line
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		circle := g.engine.getSelected()
		if circle != nil {
			// opposite hue
			h, _, _ := circle.color.Hcl()
			h += 180
			if h > 360 {
				h -= 360
			}
			clr := colorful.Hcl(h, 1.0, 0.75)
			drawLine(cursorPos, circle.pos, 2, screen, clr)
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
			msg.WriteString("\nMax Speed: ")
			msg.WriteString(strconv.FormatFloat(g.engine.maxSpeed, 'f', 2, 64))
			// msg.WriteString("\nUpdate Elapsed: ")
			// msg.WriteString(strconv.FormatFloat(g.updateElapsedTime.Seconds(), 'f', 4, 64))
			// msg.WriteString("\nDraw Elapsed: ")
			// msg.WriteString(strconv.FormatFloat(g.drawElapsedTime.Seconds(), 'f', 4, 64))

			// Draw red lines between colliding circles
			for _, p := range g.engine.collidingPairs {
				ebitenutil.DrawLine(
					screen,
					g.engine.circles[p.a].pos.X,
					g.engine.circles[p.a].pos.Y,
					g.engine.circles[p.b].pos.X,
					g.engine.circles[p.b].pos.Y,
					color.RGBA{255, 0, 0, 30},
				)
			}
		}
		ebitenutil.DebugPrint(screen, msg.String())
	}
	g.drawElapsedTime = time.Now().Sub(start)
}

// Draw2 for testing shader
func (g *Game) Draw2(screen *ebiten.Image) {
	screen.Fill(color.RGBA{128, 128, 128, 255})

	g.drawShader(g.width, 0, 0, screen)
	if g.showDebug {
		ebitenutil.DrawRect(screen, 300-1, 200-1, 200+2, 200+2, color.Black)
	}
	g.drawShader(200, 300, 200, screen)
}

func (g *Game) drawShader(size, x, y int, screen *ebiten.Image) {

	w, h := screen.Size()
	cx, cy := ebiten.CursorPosition()

	op := &ebiten.DrawRectShaderOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.Uniforms = map[string]interface{}{
		"Time":       float32(g.time) / 60,
		"Cursor":     []float32{float32(cx), float32(cy)},
		"ScreenSize": []float32{float32(w), float32(h)},
		"Translate":  []float32{float32(x), float32(y)},
		"Size":       []float32{float32(size), float32(size)},
	}
	// op.Images = [4]*ebiten.Image{
	// 	g.ShapeMap,
	// 	g.Texture,
	// },

	screen.DrawRectShader(size, size, g.circleShader, op)
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
