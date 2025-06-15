package core

import (
	"time"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ponyo877/island-merge/pkg/achievements"
	"github.com/ponyo877/island-merge/pkg/editor"
	"github.com/ponyo877/island-merge/pkg/island"
	"github.com/ponyo877/island-merge/pkg/systems"
	"github.com/ponyo877/island-merge/pkg/ui"
)

type Game struct {
	world           *World
	input           *systems.InputSystem
	render          *systems.RenderSystem
	animation       *systems.AnimationSystem
	mainMenu        *ui.Menu
	levelEditor     *editor.LevelEditor
	achievementSys  *achievements.AchievementSystem
	achievementUI   *ui.AchievementsUI
}

func NewGame() *Game {
	achievementSys := achievements.NewAchievementSystem()
	
	levelEditor := editor.NewLevelEditor()
	
	game := &Game{
		input:          systems.NewInputSystem(),
		render:         systems.NewRenderSystem(),
		animation:      systems.NewAnimationSystem(),
		levelEditor:    levelEditor,
		achievementSys: achievementSys,
		achievementUI:  ui.NewAchievementsUI(achievementSys),
	}
	
	// Set up level editor callback
	levelEditor.OnLevelCreated = func() {
		achievementSys.OnLevelCreated()
	}
	
	game.mainMenu = ui.NewMainMenu(game.handleMenuAction)
	
	// Initialize with menu state
	game.world = &World{
		State: StateMenu,
		Mode:  ModeClassic,
	}
	
	return game
}

func (g *Game) handleMenuAction(action int) {
	if action == 3 { // Level Editor
		g.world.State = StateLevelEditor
	} else {
		g.startGameMode(action)
	}
}

func (g *Game) startGameMode(mode int) {
	board := island.NewBoard(5, 5)
	board.SetupLevel1() // Simple predefined level for MVP
	
	g.world = &World{
		State:     StatePlaying,
		Mode:      GameMode(mode),
		Board:     board,
		Score:     Score{},
		StartTime: time.Now(),
	}
	
	// Set time limit for Time Attack mode
	if mode == 1 { // ModeTimeAttack
		g.world.TimeLimit = time.Minute * 2 // 2 minutes
	}
	
	// Track game start
	g.achievementSys.OnGameStart()
}

func (g *Game) Update() error {
	// Update animations and achievements UI
	g.animation.Update()
	g.achievementUI.Update()
	
	// Handle input based on game state
	if action := g.input.Update(); action != nil {
		// Check for achievement button click first
		if action.Type == systems.ActionClick && g.achievementUI.IsAchievementButtonClicked(action.X, action.Y) {
			g.achievementUI.TogglePanel()
		} else if g.achievementUI.HandleClick(action.X, action.Y) {
			// Achievement UI handled the click
		} else {
			switch g.world.State {
			case StateMenu:
				g.mainMenu.Update(action.X, action.Y, action.Type == systems.ActionClick)
			case StatePlaying:
				g.handleGameAction(action)
			case StateLevelEditor:
				if g.levelEditor.Update(action.X, action.Y, action.Type == systems.ActionClick) {
					g.world.State = StateMenu // Return to menu
				}
			}
		}
	}
	
	// Update game logic for playing state
	if g.world.State == StatePlaying && g.world.Board != nil {
		// Update timer
		g.world.Score.Time = time.Since(g.world.StartTime)
		
		// Check time limit for Time Attack mode
		if g.world.Mode == ModeTimeAttack && g.world.TimeLimit > 0 {
			if g.world.Score.Time >= g.world.TimeLimit {
				g.world.State = StateGameOver
			}
		}
		
		// Check win condition
		if g.world.Board.IsAllConnected() && !g.world.GameWon {
			g.world.GameWon = true
			// Add victory animation
			g.animation.AddAnimation(systems.AnimationVictory, 320, 240, time.Second*2)
			
			// Track achievement progress
			gameTime := g.world.Score.Time
			moves := g.world.Score.Moves
			isTimeAttack := g.world.Mode == ModeTimeAttack
			isPerfect := moves <= 2 // For the basic level, optimal is 2 moves
			
			g.achievementSys.OnGameWin(moves, gameTime, isTimeAttack, isPerfect)
		}
	}
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.world.State {
	case StateMenu:
		g.mainMenu.Draw(screen)
	case StatePlaying, StateGameOver:
		if g.world.Board != nil {
			g.render.Draw(screen, g.world.Board, g.world.Score.Moves, g.world.GameWon)
			g.render.DrawHover(screen, g.world.Board, g.input.MouseX, g.input.MouseY)
			g.render.DrawGameMode(screen, g.world)
		}
		g.render.DrawAnimations(screen, g.animation.GetAnimations())
		// Draw achievement button and UI
		g.achievementUI.DrawAchievementButton(screen, 500, 10)
	case StateLevelEditor:
		g.levelEditor.Draw(screen)
	}
	
	// Always draw achievement notifications and panel on top
	g.achievementUI.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func (g *Game) handleGameAction(action *systems.Action) {
	if action.Type == systems.ActionClick {
		// Convert screen coordinates to grid coordinates
		// Account for grid offset (160, 120) and tile size (64)
		gridX := (action.X - 160) / 64
		gridY := (action.Y - 120) / 64
		
		// Try to build bridge
		if g.world.Board.CanBuildBridge(gridX, gridY) {
			g.world.Board.BuildBridge(gridX, gridY)
			g.world.Score.Moves++
			// Add build animation
			g.animation.AddAnimation(systems.AnimationBridgeBuild, gridX, gridY, time.Millisecond*500)
			// Track bridge building achievement
			g.achievementSys.OnBridgeBuilt()
		}
	}
}