package systems

import (
	"fmt"
	"image/color"

	"math"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ponyo877/island-merge/pkg/island"
)

const (
	TileSize = 64
	GridOffsetX = 160
	GridOffsetY = 120
)

type RenderSystem struct {
	// Cache for tile images
	tileImages map[island.TileType]*ebiten.Image
}

func NewRenderSystem() *RenderSystem {
	rs := &RenderSystem{
		tileImages: make(map[island.TileType]*ebiten.Image),
	}
	rs.initTileImages()
	return rs
}

func (rs *RenderSystem) initTileImages() {
	// Create simple colored tiles for MVP
	colors := map[island.TileType]color.Color{
		island.TileSea:    color.RGBA{64, 164, 223, 255},   // Blue
		island.TileLand:   color.RGBA{139, 195, 74, 255},   // Green
		island.TileBridge: color.RGBA{121, 85, 72, 255},    // Brown
	}
	
	for tileType, col := range colors {
		img := ebiten.NewImage(TileSize, TileSize)
		img.Fill(col)
		rs.tileImages[tileType] = img
	}
}

func (rs *RenderSystem) Draw(screen *ebiten.Image, board *island.Board, moves int, gameWon bool) {
	// Clear screen
	screen.Fill(color.RGBA{240, 240, 240, 255})
	
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
	// Convert mouse to grid coordinates
	gridX := (mouseX - GridOffsetX) / TileSize
	gridY := (mouseY - GridOffsetY) / TileSize
	
	// Check if hover is valid
	if board.CanBuildBridge(gridX, gridY) {
		x := GridOffsetX + gridX*TileSize
		y := GridOffsetY + gridY*TileSize
		
		// Draw hover highlight
		highlight := ebiten.NewImage(TileSize, TileSize)
		highlight.Fill(color.RGBA{255, 255, 255, 64})
		
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(highlight, opt)
		
		// Draw border
		vector.StrokeRect(
			screen,
			float32(x), float32(y),
			float32(TileSize), float32(TileSize),
			2,
			color.RGBA{255, 255, 255, 128},
			false,
		)
	}
}

func (rs *RenderSystem) drawBoard(screen *ebiten.Image, board *island.Board) {
	for y := 0; y < board.Height; y++ {
		for x := 0; x < board.Width; x++ {
			tile := board.GetTile(x, y)
			if tile == nil {
				continue
			}
			
			// Draw tile
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(float64(GridOffsetX+x*TileSize), float64(GridOffsetY+y*TileSize))
			
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
		float32(GridOffsetX+x*TileSize),
		float32(GridOffsetY+y*TileSize),
		float32(GridOffsetX+(x+1)*TileSize),
		float32(GridOffsetY+y*TileSize),
		lineWidth,
		gridColor,
		false,
	)
	
	// Vertical line
	vector.StrokeLine(
		screen,
		float32(GridOffsetX+x*TileSize),
		float32(GridOffsetY+y*TileSize),
		float32(GridOffsetX+x*TileSize),
		float32(GridOffsetY+(y+1)*TileSize),
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
	x := float64(GridOffsetX + anim.X*TileSize + TileSize/2)
	y := float64(GridOffsetY + anim.Y*TileSize + TileSize/2)
	
	// Easing animation
	progress := EaseOutCubic(anim.Progress)
	
	// Expanding circle effect
	radius := float32(progress * float64(TileSize) * 0.8)
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