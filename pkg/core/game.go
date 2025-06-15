package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ponyo877/island-merge/pkg/island"
	"github.com/ponyo877/island-merge/pkg/systems"
)

type Game struct {
	world  *World
	input  *systems.InputSystem
	render *systems.RenderSystem
}

func NewGame() *Game {
	board := island.NewBoard(5, 5)
	board.SetupLevel1() // Simple predefined level for MVP
	
	world := &World{
		Board: board,
		Score: Score{},
	}
	
	return &Game{
		world:  world,
		input:  systems.NewInputSystem(),
		render: systems.NewRenderSystem(),
	}
}

func (g *Game) Update() error {
	// Handle input
	if action := g.input.Update(); action != nil {
		g.handleAction(action)
	}
	
	// Check win condition
	if g.world.Board.IsAllConnected() && !g.world.GameWon {
		g.world.GameWon = true
	}
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.render.Draw(screen, g.world.Board, g.world.Score.Moves, g.world.GameWon)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func (g *Game) handleAction(action *systems.Action) {
	if action.Type == systems.ActionClick {
		// Convert screen coordinates to grid coordinates
		// Account for grid offset (160, 120) and tile size (64)
		gridX := (action.X - 160) / 64
		gridY := (action.Y - 120) / 64
		
		// Try to build bridge
		if g.world.Board.CanBuildBridge(gridX, gridY) {
			g.world.Board.BuildBridge(gridX, gridY)
			g.world.Score.Moves++
		}
	}
}