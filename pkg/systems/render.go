package systems

import (
	"fmt"
	"image/color"

	"math"
	"time"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ponyo877/island-merge/pkg/island"
)

const (
	MaxTileSize = 64
	MinTileSize = 16
	GridOffsetX = 160
	GridOffsetY = 120
	MaxGridWidth = 400  // Maximum grid display width
	MaxGridHeight = 300 // Maximum grid display height
)

type RenderSystem struct {
	// Cache for tile images
	tileImages map[island.TileType]*ebiten.Image
	currentTileSize int
	viewportX, viewportY float64
	zoom float64
}

func NewRenderSystem() *RenderSystem {
	rs := &RenderSystem{
		tileImages:      make(map[island.TileType]*ebiten.Image),
		currentTileSize: MaxTileSize,
		zoom:           1.0,
	}
	rs.initTileImages()
	return rs
}

func (rs *RenderSystem) initTileImages() {
	// Initialize with max tile size, will be dynamically resized
	rs.createTileImages(MaxTileSize)
}

func (rs *RenderSystem) createTileImages(size int) {
	// Clear existing images
	rs.tileImages = make(map[island.TileType]*ebiten.Image)
	
	// Create simple colored tiles
	colors := map[island.TileType]color.Color{
		island.TileSea:    color.RGBA{64, 164, 223, 255},   // Blue
		island.TileLand:   color.RGBA{139, 195, 74, 255},   // Green
		island.TileBridge: color.RGBA{121, 85, 72, 255},    // Brown
	}
	
	for tileType, col := range colors {
		img := ebiten.NewImage(size, size)
		img.Fill(col)
		rs.tileImages[tileType] = img
	}
}

func (rs *RenderSystem) calculateTileSize(boardWidth, boardHeight int) int {
	// Calculate optimal tile size to fit the board in the available space
	maxWidthTileSize := MaxGridWidth / boardWidth
	maxHeightTileSize := MaxGridHeight / boardHeight
	
	optimalSize := min(maxWidthTileSize, maxHeightTileSize)
	
	// Clamp to min/max tile sizes
	if optimalSize < MinTileSize {
		return MinTileSize
	}
	if optimalSize > MaxTileSize {
		return MaxTileSize
	}
	
	return optimalSize
}

func (rs *RenderSystem) updateTileSize(boardWidth, boardHeight int) {
	newSize := rs.calculateTileSize(boardWidth, boardHeight)
	if newSize != rs.currentTileSize {
		rs.currentTileSize = newSize
		rs.createTileImages(newSize)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (rs *RenderSystem) Draw(screen *ebiten.Image, board *island.Board, moves int, gameWon bool) {
	// Clear screen
	screen.Fill(color.RGBA{240, 240, 240, 255})
	
	// Update tile size based on board dimensions
	if board != nil {
		rs.updateTileSize(board.Width, board.Height)
	}
	
	// Draw board
	rs.drawBoard(screen, board)
	
	// Draw UI
	rs.drawUI(screen, moves)
	
	// Draw victory message if won
	if gameWon {
		rs.drawVictory(screen)
	}
}

func (rs *RenderSystem) DrawHover(screen *ebiten.Image, board *island.Board, mouseX, mouseY int) {
	if board == nil {
		return
	}
	
	// Convert mouse to grid coordinates
	gridX := (mouseX - GridOffsetX) / rs.currentTileSize
	gridY := (mouseY - GridOffsetY) / rs.currentTileSize
	
	// Check if hover is valid
	if board.CanBuildBridge(gridX, gridY) {
		x := GridOffsetX + gridX*rs.currentTileSize
		y := GridOffsetY + gridY*rs.currentTileSize
		
		// Draw hover highlight
		highlight := ebiten.NewImage(rs.currentTileSize, rs.currentTileSize)
		highlight.Fill(color.RGBA{255, 255, 255, 64})
		
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(highlight, opt)
		
		// Draw border
		vector.StrokeRect(
			screen,
			float32(x), float32(y),
			float32(rs.currentTileSize), float32(rs.currentTileSize),
			2,
			color.RGBA{255, 255, 255, 128},
			false,
		)
	}
}

func (rs *RenderSystem) drawBoard(screen *ebiten.Image, board *island.Board) {
	if board == nil {
		return
	}
	
	for y := 0; y < board.Height; y++ {
		for x := 0; x < board.Width; x++ {
			tile := board.GetTile(x, y)
			if tile == nil {
				continue
			}
			
			// Draw tile
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(float64(GridOffsetX+x*rs.currentTileSize), float64(GridOffsetY+y*rs.currentTileSize))
			
			if img, ok := rs.tileImages[tile.Type]; ok {
				screen.DrawImage(img, opt)
			}
			
			// Draw grid lines
			rs.drawGridLines(screen, x, y)
		}
	}
}

func (rs *RenderSystem) drawGridLines(screen *ebiten.Image, x, y int) {
	gridColor := color.RGBA{200, 200, 200, 255}
	lineWidth := float32(1)
	
	// Horizontal line
	vector.StrokeLine(
		screen,
		float32(GridOffsetX+x*rs.currentTileSize),
		float32(GridOffsetY+y*rs.currentTileSize),
		float32(GridOffsetX+(x+1)*rs.currentTileSize),
		float32(GridOffsetY+y*rs.currentTileSize),
		lineWidth,
		gridColor,
		false,
	)
	
	// Vertical line
	vector.StrokeLine(
		screen,
		float32(GridOffsetX+x*rs.currentTileSize),
		float32(GridOffsetY+y*rs.currentTileSize),
		float32(GridOffsetX+x*rs.currentTileSize),
		float32(GridOffsetY+(y+1)*rs.currentTileSize),
		lineWidth,
		gridColor,
		false,
	)
}

func (rs *RenderSystem) drawUI(screen *ebiten.Image, moves int) {
	// Draw title
	ebitenutil.DebugPrintAt(screen, "Island Merge", 10, 10)
	
	// Draw moves counter
	movesText := fmt.Sprintf("Moves: %d", moves)
	ebitenutil.DebugPrintAt(screen, movesText, 10, 30)
	
	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "Click on sea tiles to build bridges", 10, 50)
	ebitenutil.DebugPrintAt(screen, "Connect all islands to win!", 10, 70)
}

func (rs *RenderSystem) drawVictory(screen *ebiten.Image) {
	// Draw semi-transparent overlay
	overlay := ebiten.NewImage(640, 480)
	overlay.Fill(color.RGBA{0, 0, 0, 128})
	screen.DrawImage(overlay, nil)
	
	// Draw victory message
	msg := "Victory! All islands connected!"
	bounds := screen.Bounds()
	x := bounds.Dx()/2 - len(msg)*3
	y := bounds.Dy()/2
	
	ebitenutil.DebugPrintAt(screen, msg, x, y)
}

func (rs *RenderSystem) DrawAnimations(screen *ebiten.Image, animations []*Animation) {
	for _, anim := range animations {
		switch anim.Type {
		case AnimationBridgeBuild:
			rs.drawBridgeBuildAnimation(screen, anim)
		case AnimationVictory:
			rs.drawVictoryAnimation(screen, anim)
		}
	}
}

func (rs *RenderSystem) drawBridgeBuildAnimation(screen *ebiten.Image, anim *Animation) {
	// Calculate position
	x := float64(GridOffsetX + anim.X*rs.currentTileSize + rs.currentTileSize/2)
	y := float64(GridOffsetY + anim.Y*rs.currentTileSize + rs.currentTileSize/2)
	
	// Easing animation
	progress := EaseOutCubic(anim.Progress)
	
	// Expanding circle effect
	radius := float32(progress * float64(rs.currentTileSize) * 0.8)
	alpha := uint8((1.0 - progress) * 200)
	
	// Draw expanding circle
	vector.DrawFilledCircle(
		screen,
		float32(x), float32(y),
		radius,
		color.RGBA{121, 85, 72, alpha},
		false,
	)
}

func (rs *RenderSystem) drawVictoryAnimation(screen *ebiten.Image, anim *Animation) {
	// Pulsing victory effect
	progress := anim.Progress
	pulse := math.Sin(progress * math.Pi * 4) * 0.1 + 1.0
	
	// Draw pulsing overlay
	overlay := ebiten.NewImage(640, 480)
	alpha := uint8(100 + 50*math.Sin(progress*math.Pi*2))
	overlay.Fill(color.RGBA{255, 215, 0, alpha}) // Gold color
	
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(pulse, pulse)
	opt.GeoM.Translate((1-pulse)*320, (1-pulse)*240)
	
	screen.DrawImage(overlay, opt)
}

func (rs *RenderSystem) DrawGameMode(screen *ebiten.Image, world interface{}) {
	// Type assertion to avoid circular import
	type gameWorld interface {
		GetMode() int
		GetScore() interface {
			GetMoves() int
			GetTime() time.Duration
		}
		GetTimeLimit() time.Duration
		GetState() int
	}
	
	if w, ok := world.(gameWorld); ok {
		mode := w.GetMode()
		score := w.GetScore()
		
		// Draw mode-specific UI
		var modeText string
		switch mode {
		case 0: // ModeClassic
			modeText = "Classic Mode"
		case 1: // ModeTimeAttack
			modeText = "Time Attack"
			// Draw timer
			remaining := w.GetTimeLimit() - score.GetTime()
			if remaining < 0 {
				remaining = 0
			}
			timerText := fmt.Sprintf("Time: %02d:%02d", 
				int(remaining.Minutes()), int(remaining.Seconds())%60)
			ebitenutil.DebugPrintAt(screen, timerText, 450, 10)
		case 2: // ModePuzzle
			modeText = "Puzzle Mode"
		}
		
		ebitenutil.DebugPrintAt(screen, modeText, 450, 30)
		
		// Draw score
		scoreText := fmt.Sprintf("Moves: %d", score.GetMoves())
		ebitenutil.DebugPrintAt(screen, scoreText, 450, 50)
		
		timeText := fmt.Sprintf("Time: %02d:%02d", 
			int(score.GetTime().Minutes()), int(score.GetTime().Seconds())%60)
		ebitenutil.DebugPrintAt(screen, timeText, 450, 70)
	}
}