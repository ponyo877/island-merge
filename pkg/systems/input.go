package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ActionType int

const (
	ActionClick ActionType = iota
)

type Action struct {
	Type ActionType
	X, Y int
}

type InputSystem struct {
	lastMouseX, lastMouseY int
}

func NewInputSystem() *InputSystem {
	return &InputSystem{}
}

func (is *InputSystem) Update() *Action {
	// Handle mouse clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return &Action{
			Type: ActionClick,
			X:    x,
			Y:    y,
		}
	}
	
	// Update mouse position for potential hover effects
	is.lastMouseX, is.lastMouseY = ebiten.CursorPosition()
	
	return nil
}