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
	selectedIndex  int
	dynamicIndex   int
	circles        []*Circle
	collidingPairs []collidingPair
}

func (e *Engine) selectAtPostion(x, y float64) {
	e.selectedIndex = e.circleAtPosition(x, y)
	if e.selectedIndex >= 0 {
		e.circles[e.selectedIndex].selected = true
	}
}

func (e *Engine) dynamicAtPosition(x, y float64) {
	e.dynamicIndex = e.circleAtPosition(x, y)
	if e.dynamicIndex >= 0 {
		e.circles[e.dynamicIndex].selected = true
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
		e.circles[e.selectedIndex].velX = 0
		e.circles[e.selectedIndex].velY = 0
		e.circles[e.selectedIndex].selected = false
	}
	e.selectedIndex = -1
}

func (e *Engine) dynamicRelease(x, y float64) {
	if e.dynamicIndex >= 0 {
		e.circles[e.dynamicIndex].selected = false

		e.circles[e.dynamicIndex].velX = 0.2 * (e.circles[e.dynamicIndex].posX - x)
		e.circles[e.dynamicIndex].velY = 0.2 * (e.circles[e.dynamicIndex].posY - y)
	}
	e.dynamicIndex = -1
}

func (e *Engine) getDynamicPosition(x, y float64) (float64, float64, bool) {
	if e.dynamicIndex >= 0 {
		return e.circles[e.dynamicIndex].posX, e.circles[e.dynamicIndex].posY, true
	}
	return 0, 0, false
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

type collidingPair struct {
	a int
	b int
}

func (e *Engine) update(width, height int, elapsedTime float64) {

	// Update ball positions
	for i := range e.circles {
		// apply friction
		e.circles[i].velX = e.circles[i].velX * 0.97
		e.circles[i].velY = e.circles[i].velY * 0.97

		// update velocity and position
		e.circles[i].velX += e.circles[i].accX * elapsedTime
		e.circles[i].velY += e.circles[i].accY * elapsedTime
		e.circles[i].posX += e.circles[i].velX * elapsedTime
		e.circles[i].posY += e.circles[i].velY * elapsedTime

		// wrap around the screen
		w := float64(width) + 200
		if e.circles[i].posX < -100.0 {
			e.circles[i].posX += w
		}
		if e.circles[i].posX > w-100.0 {
			e.circles[i].posX -= w
		}
		h := float64(height) + 200
		if e.circles[i].posY < -100.0 {
			e.circles[i].posY += h
		}
		if e.circles[i].posY > h-100.0 {
			e.circles[i].posY -= h
		}

		// is this needed? seems like it's good to not perpetually accelerate
		e.circles[i].accX = 0
		e.circles[i].accY = 0

		// clamp low velocity values

		// set previous position
		e.circles[i].prevPosX = e.circles[i].posX
		e.circles[i].prevPosY = e.circles[i].posY
	}

	// Resolve static collisions
	collided := true
	maxSteps := 5
	// clear slice but keep capacity
	e.collidingPairs = e.collidingPairs[:0]
	for step := maxSteps; step > 0 && collided; step-- {
		collided = false
		for i := range e.circles {
			for j := range e.circles {
				if i == j {
					continue
				}
				if e.overlap(i, j) {
					collided = true
					e.collidingPairs = append(e.collidingPairs, collidingPair{i, j})
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

	// apply acceleration from static collision displacement
	for i := range e.circles {
		// should be proportional to area
		multiplier := 10.0
		amountX := ((e.circles[i].posX - e.circles[i].prevPosX) / e.circles[i].radius) * multiplier
		amountY := ((e.circles[i].posY - e.circles[i].prevPosY) / e.circles[i].radius) * multiplier
		e.circles[i].accX = amountX
		e.circles[i].accY = amountY
	}
}
