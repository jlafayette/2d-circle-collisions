package game

import (
	"math"
	"sort"
)

// NewEngine initializes a new physics engine
func NewEngine(width, height int, circles []*Circle, capsules []*Capsule, rectangles []*collisionRect) *Engine {

	e := &Engine{
		minArea:         99999999,
		steps:           10,
		inverseSteps:    1 / 10,
		selectedCapsule: capsuleSelection{-1, true},
		capsules:        capsules,
		collisionRects:  rectangles,
	}
	for _, circle := range circles {
		e.addCircle(circle)
	}
	return e
}

type collisionRect struct {
	upperLeft     Vec2
	lowerRight    Vec2
	collidePoints []Vec2
}

// Engine handles collisions
type Engine struct {
	checks            int
	minArea           float64
	maxArea           float64
	maxRadius         float64
	maxSpeed          float64
	steps             int
	inverseSteps      float64
	selectedCircle    circleSelection
	selectedCapsule   capsuleSelection
	circles           []*Circle
	capsules          []*Capsule
	collisionRects    []*collisionRect
	collidingPairs    []collidingPair
	collidingCapsules []collidingCapsule
}

type capsuleSelection struct {
	index int
	start bool
}

type circleSelection struct {
	pointer   *Circle
	isDynamic bool
}

func (e *Engine) addCircle(circle *Circle) {
	for i := range e.circles {
		if e.circles[i].pos.X == circle.pos.X && e.circles[i].pos.Y == circle.pos.Y {
			circle.pos.X += 0.1
			circle.pos.Y += 0.1
		}
	}
	e.circles = append(e.circles, circle)
	e.maxRadius = math.Max(e.maxRadius, circle.radius)
	e.minArea = math.Min(e.minArea, circle.area)
	e.maxArea = math.Max(e.maxArea, circle.area)
	circle.id = len(e.circles)
}

func (e *Engine) selectAtPostion(pos Vec2) {
	circle := e.circleAtPosition(pos)
	e.selectedCircle.pointer = circle
	if circle != nil {
		circle.selected = true
		e.selectedCircle.isDynamic = false
	}
}

func (e *Engine) dynamicAtPosition(pos Vec2) {
	circle := e.circleAtPosition(pos)
	e.selectedCircle.pointer = circle
	if circle != nil {
		circle.selected = true
		e.selectedCircle.isDynamic = true
	}
}

func (e *Engine) selectNearestPostion(pos Vec2) {
	circle := e.circleNearestPosition(pos)
	e.selectedCircle.pointer = circle
	if circle != nil {
		circle.selected = true
		e.selectedCircle.isDynamic = false
	}
}

func (e *Engine) dynamicNearestPosition(pos Vec2) {
	circle := e.circleNearestPosition(pos)
	e.selectedCircle.pointer = circle
	if circle != nil {
		circle.selected = true
		e.selectedCircle.isDynamic = true
	}
}

func (e *Engine) circleNearestPosition(pos Vec2) *Circle {
	minDistance := math.MaxFloat64
	var closest *Circle
	for i := range e.circles {
		cx := e.circles[i].pos.X
		cy := e.circles[i].pos.Y
		cr := e.circles[i].radius
		d := (cx-pos.X)*(cx-pos.X) + (cy-pos.Y)*(cy-pos.Y)
		if d < (cr * cr) {
			return e.circles[i]
		}
		if d < minDistance {
			minDistance = d
			closest = e.circles[i]
		}
	}
	return closest
}

func (e *Engine) circleAtPosition(pos Vec2) *Circle {
	for i := range e.circles {
		cx := e.circles[i].pos.X
		cy := e.circles[i].pos.Y
		cr := e.circles[i].radius
		if (cx-pos.X)*(cx-pos.X)+(cy-pos.Y)*(cy-pos.Y) < (cr * cr) {
			return e.circles[i]
		}
	}
	return nil
}

func (e *Engine) moveSelectedTo(pos Vec2) {
	if e.selectedCircle.pointer != nil {
		e.selectedCircle.pointer.pos = pos
	}
}

func (e *Engine) applyForceToSelected(pos Vec2, speed float64) {
	circle := e.selectedCircle.pointer
	if circle != nil {
		force := pos.Sub(circle.pos)
		circle.acc = force.Scaled(0.03).Scaled(speed)
	}
}

func (e *Engine) deselect() {
	if e.selectedCircle.pointer != nil {
		e.selectedCircle.pointer.selected = false
		e.selectedCircle.pointer = nil
	}
}

func (e *Engine) dynamicRelease(pos Vec2) {
	circle := e.selectedCircle.pointer
	if circle != nil {
		circle.selected = false
		force := circle.pos.Sub(pos)
		s := remap(circle.area, e.minArea, e.maxArea, 0.225, 0.04)
		circle.acc = force.Scaled(s)
		circle.activity += force.Len() * 0.1
	}
	e.selectedCircle.pointer = nil
	e.selectedCircle.isDynamic = false
}

func (e *Engine) getSelected() *Circle {
	if e.selectedCircle.pointer != nil && !e.selectedCircle.isDynamic {
		return e.selectedCircle.pointer
	}
	return nil
}

func (e *Engine) getSelectedPosition() (Vec2, bool) {
	circle := e.selectedCircle.pointer
	if circle != nil && !e.selectedCircle.isDynamic {
		return circle.pos, true
	}
	return Vec2{0, 0}, false
}

func (e *Engine) getDynamic() *Circle {
	circle := e.selectedCircle.pointer
	if circle != nil && e.selectedCircle.isDynamic {
		return circle
	}
	return nil
}

func (e *Engine) selectCapsuleAtPostion(pos Vec2) bool {
	for i := range e.capsules {
		v := e.capsules[i].start
		r := e.circles[i].radius
		if (v.X-pos.X)*(v.X-pos.X)+(v.Y-pos.Y)*(v.Y-pos.Y) < (r * r) {
			e.selectedCapsule.index = i
			e.selectedCapsule.start = true
			return true
		}
		v = e.capsules[i].end
		if (v.X-pos.X)*(v.X-pos.X)+(v.Y-pos.Y)*(v.Y-pos.Y) < (r * r) {
			e.selectedCapsule.index = i
			e.selectedCapsule.start = false
			return true
		}
	}
	e.selectedCapsule.index = -1
	return false
}
func (e *Engine) moveSelectedCapsuleTo(pos Vec2) bool {
	if e.selectedCapsule.index >= 0 {
		if e.selectedCapsule.start {
			e.capsules[e.selectedCapsule.index].start = pos
		} else {
			e.capsules[e.selectedCapsule.index].end = pos
		}
		return true
	}
	return false
}
func (e *Engine) deselectCapsule() {
	e.selectedCapsule.index = -1
}

func (e *Engine) overlap(i, j int) bool {

	// This looks ugly, but here it is without all the index lookups
	// (x1-x2)*(x1-x2)+(y1-y2)*(y1-y2) < (r1+r2)*(r1+r2)

	return (e.circles[i].pos.X-e.circles[j].pos.X)*(e.circles[i].pos.X-e.circles[j].pos.X)+(e.circles[i].pos.Y-e.circles[j].pos.Y)*(e.circles[i].pos.Y-e.circles[j].pos.Y) < (e.circles[i].radius+e.circles[j].radius)*(e.circles[i].radius+e.circles[j].radius)
}

type collidingPair struct {
	a int
	b int
}

type collidingCapsule struct {
	i   int
	r   float64
	d   float64
	pos Vec2
}

func (e *Engine) update(width, height int, speed, elapsedTime float64) {
	e.checks = 0

	// set previous position
	for i := range e.circles {
		e.circles[i].prevPos = e.circles[i].pos
	}

	for i := range e.collisionRects {
		e.collisionRects[i].collidePoints = e.collisionRects[i].collidePoints[:0] // clear slice but keep capacity
	}

	stepSpeed := speed / float64(e.steps)
	for step := e.steps; step > 0; step-- {
		e.updateCirclePositions(width, height, stepSpeed, elapsedTime)
		e.sortCircles()
		e.resolveStaticCollisions()
		e.resolveDynamicCollisions()
	}

	// find max speed
	e.maxSpeed = 0
	for i := range e.circles {
		e.circles[i].postUpdate()
		e.maxSpeed = math.Max(e.maxSpeed, e.circles[i].speed)
	}
}

func (e *Engine) updateCirclePositions(width, height int, speed, elapsedTime float64) {
	// Update ball positions
	for i := range e.circles {

		// apply friction
		frictionAmount := remap(e.circles[i].area, e.minArea, e.maxArea, 0.015, 0.007)
		friction := e.circles[i].acc.Sub(e.circles[i].vel.Scaled(frictionAmount).Scaled(speed))

		// update velocity and position
		e.circles[i].vel = e.circles[i].vel.Add(friction)

		posChange := e.circles[i].vel.Scaled(elapsedTime).Scaled(speed)
		e.circles[i].pos = e.circles[i].pos.Add(posChange)

		e.circles[i].acc = Vec2{0, 0}
	}
}

func (e *Engine) sortCircles() {
	// sort by x position
	sort.Slice(e.circles, func(i, j int) bool {
		return e.circles[i].pos.X < e.circles[j].pos.X
	})
}

func (e *Engine) resolveStaticCollisions() {
	// Resolve static collisions
	e.collidingPairs = e.collidingPairs[:0]       // clear slice but keep capacity
	e.collidingCapsules = e.collidingCapsules[:0] // clear slice but keep capacity

	for i := range e.circles {
		for j := i + 1; j < len(e.circles); j++ {
			e.checks++
			if e.overlap(i, j) {
				e.collidingPairs = append(e.collidingPairs, collidingPair{i, j})
				// distance between ball centers
				r1 := e.circles[i].radius
				r2 := e.circles[j].radius
				v := e.circles[i].pos.Sub(e.circles[j].pos)
				distance := v.Len()
				unit := v.Scaled(1.0 / distance)
				if e.selectedCircle.pointer != nil && e.circles[i].id == e.selectedCircle.pointer.id {
					// displace target circle away from collision
					amount := distance - r1 - r2
					e.circles[j].pos = e.circles[j].pos.Add(unit.Scaled(amount))
				} else {
					// Make displace amount depend on area
					totalAmount := distance - r1 - r2
					a1 := e.circles[i].area
					a2 := e.circles[j].area
					areaSumM := 1.0 / (a1 + a2)
					amount1 := totalAmount * a2 * areaSumM
					amount2 := totalAmount * a1 * areaSumM
					// displace current circle away from the collision
					e.circles[i].pos = e.circles[i].pos.Sub(unit.Scaled(amount1))
					// displace target circle away from collision
					e.circles[j].pos = e.circles[j].pos.Add(unit.Scaled(amount2))

					// boose circle activity based on speed of collision
					energy := e.circles[i].speed + e.circles[j].speed
					e.circles[i].addCollisionEnergy(energy * e.inverseSteps)
					e.circles[j].addCollisionEnergy(energy * e.inverseSteps)
				}
			} else {
				if e.circles[j].pos.X > e.circles[i].pos.X+e.circles[i].radius+e.maxRadius {
					break
				}
			}
		}

		// line collisions
		for j := range e.capsules {
			lx1 := e.capsules[j].start.X
			ly1 := e.capsules[j].start.Y
			lx2 := e.capsules[j].end.X
			ly2 := e.capsules[j].end.Y
			lr := e.capsules[j].radius
			cx := e.circles[i].pos.X
			cy := e.circles[i].pos.Y
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
					collidingCapsule{i, lr, dist, Vec2{closestPointX, closestPointY}},
				)

				// Calculate displacement required
				amount := dist - cr - lr

				// displace circle away from collision
				distanceM := 1.0 / dist // Can be used to multiply instead of divide by dist
				e.circles[i].pos.X -= amount * (cx - closestPointX) * distanceM
				e.circles[i].pos.Y -= amount * (cy - closestPointY) * distanceM

				// TODO: Add ball and line pair to dynamic collisions
			}
		}

		// Rectangle collisions
		for j := range e.collisionRects {
			upperLeft := e.collisionRects[j].upperLeft
			lowerRight := e.collisionRects[j].lowerRight
			// nearest point
			x := clamp(e.circles[i].pos.X, upperLeft.X, lowerRight.X)
			y := clamp(e.circles[i].pos.Y, upperLeft.Y, lowerRight.Y)
			nearest := Vec2{x, y}
			v := e.circles[i].pos.To(nearest)
			dist := v.Len()
			if dist < e.circles[i].radius {

				// If circle is mostly inside, push nearest point out to nearest edge
				// TODO: Move this to dynamic collision resolution section
				dTp := math.Abs(upperLeft.Y - y)
				dBt := math.Abs(lowerRight.Y - y)
				dLf := math.Abs(upperLeft.X - x)
				dRt := math.Abs(lowerRight.X - x)
				if dTp <= dBt && dTp <= dLf && dTp <= dRt {
					y = upperLeft.Y
					e.circles[i].vel.Y = -e.circles[i].vel.Y
				} else if dBt <= dTp && dBt <= dLf && dBt <= dRt {
					y = lowerRight.Y
					e.circles[i].vel.Y = -e.circles[i].vel.Y
				} else if dLf <= dTp && dLf <= dBt && dLf <= dRt {
					x = upperLeft.X
					e.circles[i].vel.X = -e.circles[i].vel.X
				} else if dRt <= dTp && dRt <= dBt && dRt <= dLf {
					x = lowerRight.X
					e.circles[i].vel.X = -e.circles[i].vel.X
				} else {
					x = lowerRight.X
					e.circles[i].vel.X = -e.circles[i].vel.X
				}

				if dist > 0 {
					// Circle is mostly outside
					e.collisionRects[j].collidePoints = append(e.collisionRects[j].collidePoints, nearest)
					// Calculate displacement required
					amount := dist - e.circles[i].radius
					// displace circle away from collision
					e.circles[i].pos = e.circles[i].pos.Add(v.Unit().Scaled(amount))
				} else {

					e.collisionRects[j].collidePoints = append(e.collisionRects[j].collidePoints, Vec2{x, y})

					nearest = Vec2{x, y}
					v = e.circles[i].pos.To(nearest)
					dist = v.Len()
					// Calculate displacement required
					amount := dist + e.circles[i].radius
					// displace circle away from collision
					e.circles[i].pos = e.circles[i].pos.Add(v.Unit().Scaled(amount))
				}
			}
		}
	}
}

func (e *Engine) resolveDynamicCollisions() {
	// dynamic collisions
	for _, cap := range e.collidingCapsules {
		a1 := e.circles[cap.i].area
		v2 := e.circles[cap.i].vel.Scaled(-1.0)
		a2 := a1

		// Normalized
		nV := e.circles[cap.i].pos.To(cap.pos).Unit()

		// Calculate new velocities from elastic collision
		// https://en.wikipedia.org/wiki/Elastic_collision
		kV := e.circles[cap.i].vel.Sub(v2)
		p := 2.0 * nV.Dot(kV) / (a1 + a2)
		e.circles[cap.i].vel = e.circles[cap.i].vel.Sub(nV.Scaled(p).Scaled(a2))
	}

	for _, pair := range e.collidingPairs {
		a1 := e.circles[pair.a].area
		a2 := e.circles[pair.b].area

		// Normalized
		nV := e.circles[pair.a].pos.To(e.circles[pair.b].pos).Unit()

		// Calculate new velocities from elastic collision
		// https://en.wikipedia.org/wiki/Elastic_collision
		kV := e.circles[pair.a].vel.Sub(e.circles[pair.b].vel)
		p := 2.0 * nV.Dot(kV) / (a1 + a2)
		e.circles[pair.a].vel = e.circles[pair.a].vel.Sub(nV.Scaled(p).Scaled(a2))
		e.circles[pair.b].vel = e.circles[pair.b].vel.Add(nV.Scaled(p).Scaled(a1))
	}
}
