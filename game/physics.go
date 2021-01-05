package game

import "math"

// NewEngine initializes a new physics engine
func NewEngine(circles []*Circle) *Engine {
	return &Engine{
		selectedIndex: -1,
		circles:       circles,
	}
}

// Engine handles collisions
type Engine struct {
	selectedIndex int
	circles       []*Circle
}

func (e *Engine) selectAtPostion(x, y float64) {
	e.selectedIndex = e.circleAtPosition(x, y)
}

func (e *Engine) circleAtPosition(x, y float64) int {
	for i := range e.circles {
		var cx = e.circles[i].posX
		var cy = e.circles[i].posY
		var cr = e.circles[i].radius
		if math.Abs((cx-x)*(cx-x)+(cy-y)*(cy-y)) < (cr * cr) {
			return i
		}
	}
	return -1
}

func (e *Engine) moveSelectedTo(x, y float64) {
	e.moveCircleTo(e.selectedIndex, x, y)
}

func (e *Engine) moveCircleTo(index int, x, y float64) {
	if index >= 0 && index < len(e.circles) {
		e.circles[index].posX = x
		e.circles[index].posY = y
	}
}

func (e *Engine) deselect() {
	e.selectedIndex = -1
}

func (e *Engine) overlap(i, j int) bool {
	var x1 = e.circles[i].posX
	var y1 = e.circles[i].posY
	var r1 = e.circles[i].radius
	var x2 = e.circles[j].posX
	var y2 = e.circles[j].posY
	var r2 = e.circles[j].radius
	return math.Abs((x1-x2)*(x1-x2)+(y1-y2)*(y1-y2)) < (r1+r2)*(r1+r2)
}

func (e *Engine) update() {
	collided := true
	maxSteps := 5
	for step := maxSteps; step > 0 && collided; step-- {
		collided = false
		for i := 0; i < len(e.circles); i++ {
			for j := 0; j < len(e.circles); j++ {
				if i == j {
					continue
				}
				if e.overlap(i, j) {
					collided = true
					// distance between ball centers
					var x1 = e.circles[i].posX
					var y1 = e.circles[i].posY
					var r1 = e.circles[i].radius
					var x2 = e.circles[j].posX
					var y2 = e.circles[j].posY
					var r2 = e.circles[j].radius
					var distance = math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
					var overlapAmount = 0.5 * (distance - r1 - r2)
					// displace current circle away from the collision
					e.circles[i].posX -= overlapAmount * (x1 - x2) / distance
					e.circles[i].posY -= overlapAmount * (y1 - y2) / distance

					// displace target circle away from collision
					e.circles[j].posX += overlapAmount * (x1 - x2) / distance
					e.circles[j].posY += overlapAmount * (y1 - y2) / distance
				}
			}
		}
	}
}
