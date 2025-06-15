package systems

import (
	"time"
)

type AnimationType int

const (
	AnimationBridgeBuild AnimationType = iota
	AnimationTileHover
	AnimationVictory
)

type Animation struct {
	Type       AnimationType
	X, Y       int
	StartTime  time.Time
	Duration   time.Duration
	Progress   float64
	Data       interface{}
}

type AnimationSystem struct {
	animations []*Animation
}

func NewAnimationSystem() *AnimationSystem {
	return &AnimationSystem{
		animations: make([]*Animation, 0),
	}
}

func (as *AnimationSystem) AddAnimation(animType AnimationType, x, y int, duration time.Duration) {
	anim := &Animation{
		Type:      animType,
		X:         x,
		Y:         y,
		StartTime: time.Now(),
		Duration:  duration,
		Progress:  0,
	}
	as.animations = append(as.animations, anim)
}

func (as *AnimationSystem) Update() {
	now := time.Now()
	
	// Update animations and remove completed ones
	activeAnimations := make([]*Animation, 0)
	for _, anim := range as.animations {
		elapsed := now.Sub(anim.StartTime)
		anim.Progress = float64(elapsed) / float64(anim.Duration)
		
		if anim.Progress < 1.0 {
			activeAnimations = append(activeAnimations, anim)
		}
	}
	
	as.animations = activeAnimations
}

func (as *AnimationSystem) GetAnimations() []*Animation {
	return as.animations
}

// Easing functions for smooth animations
func EaseOutCubic(t float64) float64 {
	t = t - 1
	return t*t*t + 1
}

func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	t = 2*t - 2
	return 1 + t*t*t/2
}