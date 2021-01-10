package game

import (
	"math"
)

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
	if e.selectedIndex >= 0 {
		e.circles[e.selectedIndex].selected = true
	}
}

func (e *Engine) circleAtPosition(x, y float64) int {
	for i := range e.circles {
		cx := e.circles[i].posX
		cy := e.circles[i].posY
		cr := e.circles[i].radius
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
	if e.selectedIndex >= 0 {
		e.circles[e.selectedIndex].selected = false
	}
	e.selectedIndex = -1
}

func (e *Engine) overlap(i, j int) bool {
	x1 := e.circles[i].posX
	y1 := e.circles[i].posY
	r1 := e.circles[i].radius
	x2 := e.circles[j].posX
	y2 := e.circles[j].posY
	r2 := e.circles[j].radius
	return math.Abs((x1-x2)*(x1-x2)+(y1-y2)*(y1-y2)) < (r1+r2)*(r1+r2)
}

func (e *Engine) update() {
	collided := true
	maxSteps := 5
	for step := maxSteps; step > 0 && collided; step-- {
		collided = false
		for i := range e.circles {
			for j := range e.circles {
				if i == j {
					continue
				}
				if e.overlap(i, j) {
					collided = true
					// distance between ball centers
					x1 := e.circles[i].posX
					y1 := e.circles[i].posY
					r1 := e.circles[i].radius
					x2 := e.circles[j].posX
					y2 := e.circles[j].posY
					r2 := e.circles[j].radius
					distance := math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
					if i == e.selectedIndex {
						// displace target circle away from collision
						amount := distance - r1 - r2
						e.circles[j].posX += amount * (x1 - x2) / distance
						e.circles[j].posY += amount * (y1 - y2) / distance
					} else {
						amount := 0.5 * (distance - r1 - r2)
						// displace current circle away from the collision
						e.circles[i].posX -= amount * (x1 - x2) / distance
						e.circles[i].posY -= amount * (y1 - y2) / distance
						// displace target circle away from collision
						e.circles[j].posX += amount * (x1 - x2) / distance
						e.circles[j].posY += amount * (y1 - y2) / distance
					}
				}
			}
		}
	}
}
