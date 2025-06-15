package editor

import (
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ponyo877/island-merge/pkg/island"
)

type EditorMode int

const (
	ModePaint EditorMode = iota
	ModeErase
	ModeTest
)

type Tool int

const (
	ToolLand Tool = iota
	ToolSea
	ToolEmpty
)

type LevelEditor struct {
	Board          *island.Board
	Mode           EditorMode
	Tool           Tool
	IsPlaying      bool
	TestBoard      *island.Board // For testing the level
	UIButtons      []*UIButton
	OnLevelCreated func()        // Callback for achievement tracking
}

type UIButton struct {
	Text     string
	X, Y     float64
	Width    float64
	Height   float64
	Action   func()
	Color    color.Color
	Hovered  bool
}

const (
	EditorTileSize   = 32
	EditorGridX      = 50
	EditorGridY      = 100
	EditorGridWidth  = 16
	EditorGridHeight = 12
)

func NewLevelEditor() *LevelEditor {
	board := island.NewBoard(EditorGridWidth, EditorGridHeight)
	
	editor := &LevelEditor{
		Board:     board,
		Mode:      ModePaint,
		Tool:      ToolLand,
		IsPlaying: false,
		UIButtons: make([]*UIButton, 0),
	}
	
	editor.setupUI()
	return editor
}

func (le *LevelEditor) setupUI() {
	buttonY := 20.0
	buttonWidth := 80.0
	buttonHeight := 30.0
	spacing := 10.0
	
	buttons := []struct {
		text   string
		color  color.Color
		action func()
	}{
		{"Land", color.RGBA{139, 195, 74, 255}, func() { le.Tool = ToolLand }},
		{"Sea", color.RGBA{64, 164, 223, 255}, func() { le.Tool = ToolSea }},
		{"Empty", color.RGBA{200, 200, 200, 255}, func() { le.Tool = ToolEmpty }},
		{"Clear", color.RGBA{255, 100, 100, 255}, func() { le.clearBoard() }},
		{"Test", color.RGBA{100, 255, 100, 255}, func() { le.testLevel() }},
		{"Export", color.RGBA{255, 255, 100, 255}, func() { le.exportLevel() }},
		{"Back", color.RGBA{150, 150, 150, 255}, nil}, // Will be handled by parent
	}
	
	for i, btn := range buttons {
		button := &UIButton{
			Text:   btn.text,
			X:      50 + float64(i)*(buttonWidth+spacing),
			Y:      buttonY,
			Width:  buttonWidth,
			Height: buttonHeight,
			Action: btn.action,
			Color:  btn.color,
		}
		le.UIButtons = append(le.UIButtons, button)
	}
}

func (le *LevelEditor) Update(mouseX, mouseY int, clicked bool) bool {
	// Update UI buttons
	backClicked := false
	for i, btn := range le.UIButtons {
		btn.Hovered = float64(mouseX) >= btn.X && float64(mouseX) <= btn.X+btn.Width &&
			float64(mouseY) >= btn.Y && float64(mouseY) <= btn.Y+btn.Height
		
		if btn.Hovered && clicked {
			if btn.Action != nil {
				btn.Action()
			} else if i == len(le.UIButtons)-1 { // Back button
				backClicked = true
			}
		}
	}
	
	if backClicked {
		return true // Signal to return to menu
	}
	
	// Handle grid clicks
	if clicked {
		gridX := (mouseX - EditorGridX) / EditorTileSize
		gridY := (mouseY - EditorGridY) / EditorTileSize
		
		if gridX >= 0 && gridX < EditorGridWidth && gridY >= 0 && gridY < EditorGridHeight {
			if le.IsPlaying {
				le.handleTestClick(gridX, gridY)
			} else {
				le.paintTile(gridX, gridY)
			}
		}
	}
	
	return false
}

func (le *LevelEditor) handleTestClick(x, y int) {
	if le.TestBoard == nil {
		return
	}
	
	// Convert to game coordinates (test board uses smaller tiles)
	if le.TestBoard.CanBuildBridge(x, y) {
		le.TestBoard.BuildBridge(x, y)
	}
}

func (le *LevelEditor) paintTile(x, y int) {
	switch le.Tool {
	case ToolLand:
		le.Board.SetTile(x, y, island.TileLand)
	case ToolSea:
		le.Board.SetTile(x, y, island.TileSea)
	case ToolEmpty:
		le.Board.SetTile(x, y, island.TileEmpty)
	}
}

func (le *LevelEditor) clearBoard() {
	for y := 0; y < le.Board.Height; y++ {
		for x := 0; x < le.Board.Width; x++ {
			le.Board.SetTile(x, y, island.TileEmpty)
		}
	}
}

func (le *LevelEditor) testLevel() {
	if le.IsPlaying {
		le.IsPlaying = false
		le.TestBoard = nil
	} else {
		// Create test board copy
		le.TestBoard = island.NewBoard(le.Board.Width, le.Board.Height)
		for y := 0; y < le.Board.Height; y++ {
			for x := 0; x < le.Board.Width; x++ {
				tile := le.Board.GetTile(x, y)
				if tile != nil {
					le.TestBoard.SetTile(x, y, tile.Type)
				}
			}
		}
		le.IsPlaying = true
	}
}

func (le *LevelEditor) exportLevel() {
	levelData := le.createLevelData()
	jsonData, err := json.MarshalIndent(levelData, "", "  ")
	if err != nil {
		fmt.Println("Export error:", err)
		return
	}
	
	// In a real implementation, this would save to file or clipboard
	fmt.Println("Level exported:")
	fmt.Println(string(jsonData))
	
	// Notify achievement system (this will be called from the game)
	if le.OnLevelCreated != nil {
		le.OnLevelCreated()
	}
}

func (le *LevelEditor) createLevelData() map[string]interface{} {
	tiles := make([][]int, le.Board.Height)
	for y := 0; y < le.Board.Height; y++ {
		tiles[y] = make([]int, le.Board.Width)
		for x := 0; x < le.Board.Width; x++ {
			tile := le.Board.GetTile(x, y)
			if tile != nil {
				tiles[y][x] = int(tile.Type)
			}
		}
	}
	
	return map[string]interface{}{
		"name":   "Custom Level",
		"width":  le.Board.Width,
		"height": le.Board.Height,
		"tiles":  tiles,
	}
}

func (le *LevelEditor) Draw(screen *ebiten.Image) {
	// Clear background
	screen.Fill(color.RGBA{240, 240, 240, 255})
	
	// Draw title
	ebitenutil.DebugPrintAt(screen, "Level Editor", 300, 20)
	
	// Draw UI buttons
	le.drawUI(screen)
	
	// Draw grid
	le.drawGrid(screen)
	
	// Draw instructions
	le.drawInstructions(screen)
}

func (le *LevelEditor) drawUI(screen *ebiten.Image) {
	for _, btn := range le.UIButtons {
		// Button background
		btnColor := btn.Color
		if btn.Hovered {
			// Brighten on hover
			r, g, b, a := btn.Color.RGBA()
			btnColor = color.RGBA{
				uint8(min(255, int((r>>8)+30))),
				uint8(min(255, int((g>>8)+30))),
				uint8(min(255, int((b>>8)+30))),
				uint8(a >> 8),
			}
		}
		
		vector.DrawFilledRect(
			screen,
			float32(btn.X), float32(btn.Y),
			float32(btn.Width), float32(btn.Height),
			btnColor,
			false,
		)
		
		// Button border
		vector.StrokeRect(
			screen,
			float32(btn.X), float32(btn.Y),
			float32(btn.Width), float32(btn.Height),
			2,
			color.RGBA{100, 100, 100, 255},
			false,
		)
		
		// Button text
		textX := int(btn.X + btn.Width/2 - float64(len(btn.Text)*3))
		textY := int(btn.Y + btn.Height/2 - 4)
		ebitenutil.DebugPrintAt(screen, btn.Text, textX, textY)
	}
	
	// Draw current tool indicator
	toolText := fmt.Sprintf("Current Tool: %s", le.getToolName())
	ebitenutil.DebugPrintAt(screen, toolText, 50, 70)
}

func (le *LevelEditor) getToolName() string {
	switch le.Tool {
	case ToolLand:
		return "Land"
	case ToolSea:
		return "Sea"
	case ToolEmpty:
		return "Empty"
	default:
		return "Unknown"
	}
}

func (le *LevelEditor) drawGrid(screen *ebiten.Image) {
	board := le.Board
	if le.IsPlaying && le.TestBoard != nil {
		board = le.TestBoard
	}
	
	for y := 0; y < board.Height; y++ {
		for x := 0; x < board.Width; x++ {
			drawX := EditorGridX + x*EditorTileSize
			drawY := EditorGridY + y*EditorTileSize
			
			tile := board.GetTile(x, y)
			tileColor := color.RGBA{200, 200, 200, 255} // Empty
			
			if tile != nil {
				switch tile.Type {
				case island.TileLand:
					tileColor = color.RGBA{139, 195, 74, 255} // Green
				case island.TileSea:
					tileColor = color.RGBA{64, 164, 223, 255} // Blue
				case island.TileBridge:
					tileColor = color.RGBA{121, 85, 72, 255} // Brown
				}
			}
			
			// Draw tile
			vector.DrawFilledRect(
				screen,
				float32(drawX), float32(drawY),
				float32(EditorTileSize), float32(EditorTileSize),
				tileColor,
				false,
			)
			
			// Draw grid lines
			vector.StrokeRect(
				screen,
				float32(drawX), float32(drawY),
				float32(EditorTileSize), float32(EditorTileSize),
				1,
				color.RGBA{150, 150, 150, 255},
				false,
			)
		}
	}
}

func (le *LevelEditor) drawInstructions(screen *ebiten.Image) {
	instructions := []string{
		"Click tiles to paint with selected tool",
		"Use Test button to play your level",
		"Export saves level data to console",
	}
	
	for i, instruction := range instructions {
		ebitenutil.DebugPrintAt(screen, instruction, 50, 450+i*15)
	}
	
	if le.IsPlaying {
		ebitenutil.DebugPrintAt(screen, "TEST MODE - Click Test again to return to editing", 50, 400)
	}
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}