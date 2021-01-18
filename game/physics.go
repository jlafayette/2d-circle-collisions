package game

import (
	"math"
)

// NewEngine initializes a new physics engine
func NewEngine(circles []*Circle, capsules []*Capsule) *Engine {
	return &Engine{
		selectedIndex: -1,
		circles:       circles,
		capsules:      capsules,
	}
}

// Engine handles collisions
type Engine struct {
	selectedIndex     int
	dynamicIndex      int
	checks            int
	circles           []*Circle
	capsules          []*Capsule
	collidingPairs    []collidingPair
	collidingCapsules []collidingCapsule
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

func (e *Engine) selectNearestPostion(x, y float64) {
	e.selectedIndex = e.circleNearestPosition(x, y)
	if e.selectedIndex >= 0 {
		e.circles[e.selectedIndex].selected = true
	}
}

func (e *Engine) dynamicNearestPosition(x, y float64) {
	e.dynamicIndex = e.circleNearestPosition(x, y)
	if e.dynamicIndex >= 0 {
		e.circles[e.dynamicIndex].selected = true
	}
}

func (e *Engine) circleNearestPosition(x, y float64) int {
	minDistance := math.MaxFloat64
	closest := -1
	for i := range e.circles {
		cx := e.circles[i].posX
		cy := e.circles[i].posY
		cr := e.circles[i].radius
		d := math.Abs((cx-x)*(cx-x) + (cy-y)*(cy-y))
		if d < (cr * cr) {
			return i
		}
		if d < minDistance {
			minDistance = d
			closest = i
		}
	}
	return closest
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

func (e *Engine) applyForceToSelected(x, y, speed float64) {
	if e.selectedIndex >= 0 {
		e.circles[e.selectedIndex].accX = (0.03 * (x - e.circles[e.selectedIndex].posX)) * speed
		e.circles[e.selectedIndex].accY = (0.03 * (y - e.circles[e.selectedIndex].posY)) * speed
	}
}

func (e *Engine) deselect() {
	if e.selectedIndex >= 0 {
		e.circles[e.selectedIndex].selected = false
	}
	e.selectedIndex = -1
}

func (e *Engine) dynamicRelease(x, y float64) {
	if e.dynamicIndex >= 0 {
		e.circles[e.dynamicIndex].selected = false
		e.circles[e.dynamicIndex].accX = 0.2 * (e.circles[e.dynamicIndex].posX - x)
		e.circles[e.dynamicIndex].accY = 0.2 * (e.circles[e.dynamicIndex].posY - y)
	}
	e.dynamicIndex = -1
}

func (e *Engine) getSelectedPosition(x, y float64) (float64, float64, bool) {
	if e.selectedIndex >= 0 {
		return e.circles[e.selectedIndex].posX, e.circles[e.selectedIndex].posY, true
	}
	return 0, 0, false
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

type collidingCapsule struct {
	i int
	x float64
	y float64
	r float64
	d float64
}

func (e *Engine) update(width, height int, speed, elapsedTime float64) {
	e.checks = 0
	steps := 5
	stepSpeed := speed / float64(steps)
	for step := steps; step > 0; step-- {
		e.updateCirclePositions(width, height, stepSpeed, elapsedTime)
		e.resolveStaticCollisions()
		e.resolveDynamicCollisions()
	}

	// // apply acceleration from static collision displacement
	// for i := range e.circles {
	// 	// should be proportional to area
	// 	multiplier := 10.0
	// 	amountX := ((e.circles[i].posX - e.circles[i].prevPosX) / e.circles[i].radius) * multiplier
	// 	amountY := ((e.circles[i].posY - e.circles[i].prevPosY) / e.circles[i].radius) * multiplier
	// 	e.circles[i].accX = amountX
	// 	e.circles[i].accY = amountY
	// }
}

func (e *Engine) updateCirclePositions(width, height int, speed, elapsedTime float64) {
	// Update ball positions
	for i := range e.circles {

		// apply friction
		accX := e.circles[i].accX - (e.circles[i].velX * 0.02 * speed)
		accY := e.circles[i].accY - (e.circles[i].velY * 0.02 * speed)

		// update velocity and position
		e.circles[i].velX += accX
		e.circles[i].velY += accY
		e.circles[i].posX += e.circles[i].velX * elapsedTime * speed
		e.circles[i].posY += e.circles[i].velY * elapsedTime * speed

		e.circles[i].accX = 0.0
		e.circles[i].accY = 0.0

		// // wrap around the screen
		// w := float64(width) + 200
		// if e.circles[i].posX < -100.0 {
		// 	e.circles[i].posX += w
		// }
		// if e.circles[i].posX > w-100.0 {
		// 	e.circles[i].posX -= w
		// }
		// h := float64(height) + 200
		// if e.circles[i].posY < -100.0 {
		// 	e.circles[i].posY += h
		// }
		// if e.circles[i].posY > h-100.0 {
		// 	e.circles[i].posY -= h
		// }

		// clamp low velocity values

		// set previous position
		e.circles[i].prevPosX = e.circles[i].posX
		e.circles[i].prevPosY = e.circles[i].posY
	}
}

func (e *Engine) resolveStaticCollisions() {
	// Resolve static collisions
	e.collidingPairs = e.collidingPairs[:0]       // clear slice but keep capacity
	e.collidingCapsules = e.collidingCapsules[:0] // clear slice but keep capacity

	for i := range e.circles {
		for j := range e.circles {
			if i == j {
				continue
			}
			e.checks++
			if e.overlap(i, j) {
				e.collidingPairs = append(e.collidingPairs, collidingPair{i, j})
				// distance between ball centers
				x1 := e.circles[i].posX
				y1 := e.circles[i].posY
				r1 := e.circles[i].radius
				x2 := e.circles[j].posX
				y2 := e.circles[j].posY
				r2 := e.circles[j].radius
				distance := math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
				distanceM := (1.0 / distance)
				if i == e.selectedIndex {
					// displace target circle away from collision
					amount := distance - r1 - r2
					e.circles[j].posX += amount * (x1 - x2) * distanceM
					e.circles[j].posY += amount * (y1 - y2) * distanceM
				} else {
					// Make displace amount depend on area
					totalAmount := distance - r1 - r2
					a1 := e.circles[i].area
					a2 := e.circles[j].area
					areaSumM := 1.0 / (a1 + a2)
					amount1 := totalAmount * a2 * areaSumM
					amount2 := totalAmount * a1 * areaSumM
					// displace current circle away from the collision
					e.circles[i].posX -= amount1 * (x1 - x2) * distanceM
					e.circles[i].posY -= amount1 * (y1 - y2) * distanceM
					// displace target circle away from collision
					e.circles[j].posX += amount2 * (x1 - x2) * distanceM
					e.circles[j].posY += amount2 * (y1 - y2) * distanceM
				}
			}
		}
		// line collisions
		for j := range e.capsules {
			lx1 := e.capsules[j].x1
			ly1 := e.capsules[j].y1
			lx2 := e.capsules[j].x2
			ly2 := e.capsules[j].y2
			lr := e.capsules[j].radius
			cx := e.circles[i].posX
			cy := e.circles[i].posY
			cr := e.circles[i].radius
			// Line vector
			lineX1 := lx2 - lx1
			lineY1 := ly2 - ly1
			// Vector from circle to start of the line
			lineX2 := cx - lx1
			lineY2 := cy - ly1

			lineLen := lineX1*lineX1 + lineY1*lineY1

			// t represents the closest point on the line segment, normalized between 0 and 1
			// where zero is the start, and one is end of the line.
			t := math.Max(0, math.Min(lineLen, (lineX1*lineX2+lineY1*lineY2))) / lineLen

			// Closest point
			closestPointX := lx1 + t*lineX1
			closestPointY := ly1 + t*lineY1

			// Distance betwen closest point and circle center
			dist := math.Sqrt((cx-closestPointX)*(cx-closestPointX) + (cy-closestPointY)*(cy-closestPointY))

			// Check for collision
			if dist <= (cr + lr) {
				e.collidingCapsules = append(
					e.collidingCapsules,
					collidingCapsule{i, closestPointX, closestPointY, lr, dist},
				)

				// Calculate displacement required
				amount := dist - cr - lr

				// displace circle away from collision
				distanceM := 1.0 / dist // Can be used to multiply instead of divide by dist
				e.circles[i].posX -= amount * (cx - closestPointX) * distanceM
				e.circles[i].posY -= amount * (cy - closestPointY) * distanceM

				// TODO: Add ball and line pair to dynamic collisions
			}
		}
	}
}

func (e *Engine) resolveDynamicCollisions() {
	// dynamic collisions
	for _, cap := range e.collidingCapsules {
		px1 := e.circles[cap.i].posX
		py1 := e.circles[cap.i].posY
		vx1 := e.circles[cap.i].velX
		vy1 := e.circles[cap.i].velY
		a1 := e.circles[cap.i].area
		px2 := cap.x
		py2 := cap.y
		vx2 := -vx1
		vy2 := -vy1
		a2 := a1

		// Distance between balls
		distance := math.Sqrt((px1-px2)*(px1-px2) + (py1-py2)*(py1-py2))

		// Normal
		nx := (px2 - px1) / distance
		ny := (py2 - py1) / distance

		// Calculate new velocities from elastic collision
		// https://en.wikipedia.org/wiki/Elastic_collision
		kx := vx1 - vx2
		ky := vy1 - vy2
		p := 2.0 * (nx*kx + ny*ky) / (a1 + a2)
		e.circles[cap.i].velX = vx1 - p*a2*nx
		e.circles[cap.i].velY = vy1 - p*a2*ny
		// e.circles[pair.b].velX = vx2 + p*a1*nx
		// e.circles[pair.b].velY = vy2 + p*a1*ny

		// // r=d−2(d⋅n)n
		// // where d⋅n is the dot product, and n must be normalized.
	}

	for _, pair := range e.collidingPairs {
		px1 := e.circles[pair.a].posX
		py1 := e.circles[pair.a].posY
		vx1 := e.circles[pair.a].velX
		vy1 := e.circles[pair.a].velY
		a1 := e.circles[pair.a].area
		px2 := e.circles[pair.b].posX
		py2 := e.circles[pair.b].posY
		vx2 := e.circles[pair.b].velX
		vy2 := e.circles[pair.b].velY
		a2 := e.circles[pair.b].area

		// Distance between balls
		distance := math.Sqrt((px1-px2)*(px1-px2) + (py1-py2)*(py1-py2))

		// Normal
		nx := (px2 - px1) / distance
		ny := (py2 - py1) / distance

		// Calculate new velocities from elastic collision
		// https://en.wikipedia.org/wiki/Elastic_collision
		kx := vx1 - vx2
		ky := vy1 - vy2
		p := 2.0 * (nx*kx + ny*ky) / (a1 + a2)
		e.circles[pair.a].velX = vx1 - p*a2*nx
		e.circles[pair.a].velY = vy1 - p*a2*ny
		e.circles[pair.b].velX = vx2 + p*a1*nx
		e.circles[pair.b].velY = vy2 + p*a1*ny
	}
}
