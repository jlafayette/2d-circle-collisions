package game

import "math"

// NewEngine initializes a new physics engine
func NewEngine(circles []*Circle) *Engine {
	return &Engine{
		circles: circles,
	}
}

// Engine handles collisions
type Engine struct {
	circles []*Circle
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
	// TODO: implement static collisions for circles
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
