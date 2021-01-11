package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// NewSpeedControl initializes a new control at full speed.
func NewSpeedControl() *SpeedControl {
	return &SpeedControl{
		control: 3,
		prev:    3,
	}
}

// SpeedControl handles time dilation and pausing.
type SpeedControl struct {
	control int
	prev    int
}

func (s *SpeedControl) multiplier() float64 {
	switch s.control {
	case 0:
		return 0.0
	case 1:
		return 0.1
	case 2:
		return 0.5
	case 3:
		return 1.0
	}
	return 1.0
}

func (s *SpeedControl) paused() bool {
	return s.control == 0
}

func (s *SpeedControl) update() {
	// Adjust game speed
	if inpututil.IsKeyJustPressed(ebiten.KeyComma) {
		if s.control > 0 {
			s.control--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPeriod) {
		if s.control < 3 {
			s.control++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if s.paused() {
			s.control = s.prev
		} else {
			s.prev = s.control
			s.control = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		s.control = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		s.control = 2
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		s.control = 3
	}
}
