package ui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ponyo877/island-merge/pkg/levels"
)

type LevelSelectUI struct {
	levelManager     *levels.LevelManager
	selectedDifficulty levels.Difficulty
	scrollOffset     float64
	showPanel        bool
	OnLevelSelected  func(*levels.LevelData)
	OnBack          func()
}

func NewLevelSelectUI(levelManager *levels.LevelManager) *LevelSelectUI {
	return &LevelSelectUI{
		levelManager:       levelManager,
		selectedDifficulty: levels.DifficultyBeginner,
		scrollOffset:       0,
		showPanel:          false,
	}
}

func (lsui *LevelSelectUI) Show() {
	lsui.showPanel = true
	lsui.scrollOffset = 0
}

func (lsui *LevelSelectUI) Hide() {
	lsui.showPanel = false
}

func (lsui *LevelSelectUI) IsShown() bool {
	return lsui.showPanel
}

func (lsui *LevelSelectUI) HandleClick(x, y int) bool {
	if !lsui.showPanel {
		return false
	}
	
	panelX, panelY := 50, 30
	panelWidth, panelHeight := 540, 420
	
	// Check if clicking outside panel
	if x < panelX || x > panelX+panelWidth || y < panelY || y > panelY+panelHeight {
		lsui.Hide()
		if lsui.OnBack != nil {
			lsui.OnBack()
		}
		return true
	}
	
	// Back button
	if x >= panelX+panelWidth-40 && x <= panelX+panelWidth-10 && y >= panelY+10 && y <= panelY+40 {
		lsui.Hide()
		if lsui.OnBack != nil {
			lsui.OnBack()
		}
		return true
	}
	
	// Difficulty tabs
	tabWidth := 120
	tabY := panelY + 50
	for i := 0; i < 4; i++ {
		tabX := panelX + 20 + i*tabWidth
		if x >= tabX && x <= tabX+tabWidth-10 && y >= tabY && y <= tabY+30 {
			lsui.selectedDifficulty = levels.Difficulty(i)
			lsui.scrollOffset = 0
			return true
		}
	}
	
	// Level selection
	lsui.handleLevelClick(x, y, panelX, panelY)
	
	return true
}

func (lsui *LevelSelectUI) handleLevelClick(x, y, panelX, panelY int) {
	levelSet := lsui.getCurrentLevelSet()
	if levelSet == nil {
		return
	}
	
	levelsStartY := panelY + 120
	levelWidth := 100
	levelHeight := 80
	levelsPerRow := 5
	spacing := 10
	
	for i, level := range levelSet.Levels {
		row := i / levelsPerRow
		col := i % levelsPerRow
		
		levelX := panelX + 20 + col*(levelWidth+spacing)
		levelY := int(float64(levelsStartY + row*(levelHeight+spacing)) - lsui.scrollOffset)
		
		// Skip if not visible
		if levelY < levelsStartY-levelHeight || levelY > panelY+400 {
			continue
		}
		
		if x >= levelX && x <= levelX+levelWidth && y >= levelY && y <= levelY+levelHeight {
			if level.Unlocked && lsui.OnLevelSelected != nil {
				lsui.OnLevelSelected(level)
				lsui.Hide()
			}
			return
		}
	}
}

func (lsui *LevelSelectUI) HandleScroll(deltaY float64) {
	if !lsui.showPanel {
		return
	}
	
	lsui.scrollOffset += deltaY * 20
	lsui.scrollOffset = math.Max(0, lsui.scrollOffset)
}

func (lsui *LevelSelectUI) getCurrentLevelSet() *levels.LevelSet {
	for _, levelSet := range lsui.levelManager.LevelSets {
		if levelSet.Difficulty == lsui.selectedDifficulty {
			return levelSet
		}
	}
	return nil
}

func (lsui *LevelSelectUI) Draw(screen *ebiten.Image) {
	if !lsui.showPanel {
		return
	}
	
	// Dark overlay
	overlay := ebiten.NewImage(640, 480)
	overlay.Fill(color.RGBA{0, 0, 0, 128})
	screen.DrawImage(overlay, nil)
	
	// Panel background
	panelX, panelY := 50, 30
	panelWidth, panelHeight := 540, 420
	
	vector.DrawFilledRect(
		screen,
		float32(panelX), float32(panelY),
		float32(panelWidth), float32(panelHeight),
		color.RGBA{240, 240, 240, 255},
		false,
	)
	
	// Panel border
	vector.StrokeRect(
		screen,
		float32(panelX), float32(panelY),
		float32(panelWidth), float32(panelHeight),
		3,
		color.RGBA{100, 100, 100, 255},
		false,
	)
	
	// Title
	ebitenutil.DebugPrintAt(screen, "Select Level", panelX+20, panelY+15)
	
	// Back button
	vector.DrawFilledRect(screen, float32(panelX+panelWidth-40), float32(panelY+10), 30, 30, color.RGBA{200, 100, 100, 255}, false)
	ebitenutil.DebugPrintAt(screen, "←", panelX+panelWidth-30, panelY+20)
	
	// Draw difficulty tabs
	lsui.drawDifficultyTabs(screen, panelX, panelY)
	
	// Draw current level set
	levelSet := lsui.getCurrentLevelSet()
	if levelSet != nil {
		lsui.drawLevelSet(screen, levelSet, panelX, panelY)
	}
}

func (lsui *LevelSelectUI) drawDifficultyTabs(screen *ebiten.Image, panelX, panelY int) {
	difficulties := []struct {
		name string
		diff levels.Difficulty
	}{
		{"Beginner", levels.DifficultyBeginner},
		{"Intermediate", levels.DifficultyIntermediate},
		{"Expert", levels.DifficultyExpert},
		{"Master", levels.DifficultyMaster},
	}
	
	tabWidth := 120
	tabHeight := 30
	tabY := panelY + 50
	
	for i, difficulty := range difficulties {
		tabX := panelX + 20 + i*tabWidth
		
		// Tab background
		bgColor := color.RGBA{200, 200, 200, 255}
		if difficulty.diff == lsui.selectedDifficulty {
			bgColor = color.RGBA{150, 150, 250, 255}
		}
		
		// Check if difficulty is unlocked
		levelSet := lsui.getLevelSetByDifficulty(difficulty.diff)
		isUnlocked := lsui.isDifficultyUnlocked(levelSet)
		if !isUnlocked {
			bgColor = color.RGBA{150, 150, 150, 128}
		}
		
		vector.DrawFilledRect(
			screen,
			float32(tabX), float32(tabY),
			float32(tabWidth-10), float32(tabHeight),
			bgColor,
			false,
		)
		
		// Tab border
		vector.StrokeRect(
			screen,
			float32(tabX), float32(tabY),
			float32(tabWidth-10), float32(tabHeight),
			1,
			color.RGBA{100, 100, 100, 255},
			false,
		)
		
		// Tab text
		textX := tabX + (tabWidth-len(difficulty.name)*6)/2
		textY := tabY + tabHeight/2 - 4
		ebitenutil.DebugPrintAt(screen, difficulty.name, textX, textY)
		
		if !isUnlocked {
			ebitenutil.DebugPrintAt(screen, "🔒", tabX+tabWidth-25, textY)
		}
	}
}

func (lsui *LevelSelectUI) drawLevelSet(screen *ebiten.Image, levelSet *levels.LevelSet, panelX, panelY int) {
	// Level set description
	descY := panelY + 90
	ebitenutil.DebugPrintAt(screen, levelSet.Description, panelX+20, descY)
	
	// Level grid
	levelsStartY := panelY + 120
	levelWidth := 100
	levelHeight := 80
	levelsPerRow := 5
	spacing := 10
	
	for i, level := range levelSet.Levels {
		row := i / levelsPerRow
		col := i % levelsPerRow
		
		levelX := panelX + 20 + col*(levelWidth+spacing)
		levelY := int(float64(levelsStartY + row*(levelHeight+spacing)) - lsui.scrollOffset)
		
		// Skip if not visible
		if levelY < levelsStartY-levelHeight || levelY > panelY+400 {
			continue
		}
		
		lsui.drawLevelButton(screen, level, levelX, levelY, levelWidth, levelHeight)
	}
}

func (lsui *LevelSelectUI) drawLevelButton(screen *ebiten.Image, level *levels.LevelData, x, y, width, height int) {
	// Background color based on status
	var bgColor color.Color
	if !level.Unlocked {
		bgColor = color.RGBA{150, 150, 150, 255} // Locked
	} else if level.Completed {
		bgColor = color.RGBA{144, 238, 144, 255} // Completed (light green)
	} else {
		bgColor = color.RGBA{255, 248, 220, 255} // Available (light yellow)
	}
	
	vector.DrawFilledRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		bgColor,
		false,
	)
	
	// Border
	borderColor := color.RGBA{100, 100, 100, 255}
	if level.Completed {
		borderColor = color.RGBA{255, 215, 0, 255} // Gold border for completed
	}
	
	vector.StrokeRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		2,
		borderColor,
		false,
	)
	
	// Level name (shortened for display)
	nameLines := lsui.splitLevelName(level.Name, width-10)
	for i, line := range nameLines {
		textX := x + (width-len(line)*6)/2
		textY := y + 10 + i*12
		ebitenutil.DebugPrintAt(screen, line, textX, textY)
	}
	
	// Size indicator
	sizeText := fmt.Sprintf("%dx%d", level.Width, level.Height)
	sizeX := x + (width-len(sizeText)*6)/2
	sizeY := y + height - 30
	ebitenutil.DebugPrintAt(screen, sizeText, sizeX, sizeY)
	
	// Stars (if completed)
	if level.Completed && level.BestScore != nil {
		lsui.drawStars(screen, level.BestScore.Stars, x+width-25, y+5)
	}
	
	// Lock icon (if locked)
	if !level.Unlocked {
		ebitenutil.DebugPrintAt(screen, "🔒", x+width/2-6, y+height/2-6)
	}
	
	// Difficulty indicator
	diffColor := lsui.getDifficultyColor(level.Difficulty)
	vector.DrawFilledRect(
		screen,
		float32(x+5), float32(y+5),
		10, 10,
		diffColor,
		false,
	)
}

func (lsui *LevelSelectUI) splitLevelName(name string, maxWidth int) []string {
	maxChars := maxWidth / 6 // Approximate character width
	if len(name) <= maxChars {
		return []string{name}
	}
	
	// Simple word wrap
	words := []string{}
	current := ""
	for _, char := range name {
		if char == ' ' && len(current) > maxChars/2 {
			words = append(words, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		words = append(words, current)
	}
	
	if len(words) == 0 {
		return []string{name[:maxChars]}
	}
	
	return words
}

func (lsui *LevelSelectUI) drawStars(screen *ebiten.Image, stars, x, y int) {
	for i := 0; i < 3; i++ {
		starChar := "☆"
		if i < stars {
			starChar = "★"
		}
		ebitenutil.DebugPrintAt(screen, starChar, x, y+i*8)
	}
}

func (lsui *LevelSelectUI) getDifficultyColor(difficulty levels.Difficulty) color.Color {
	switch difficulty {
	case levels.DifficultyBeginner:
		return color.RGBA{0, 255, 0, 255} // Green
	case levels.DifficultyIntermediate:
		return color.RGBA{255, 255, 0, 255} // Yellow
	case levels.DifficultyExpert:
		return color.RGBA{255, 165, 0, 255} // Orange
	case levels.DifficultyMaster:
		return color.RGBA{255, 0, 0, 255} // Red
	default:
		return color.RGBA{128, 128, 128, 255} // Gray
	}
}

func (lsui *LevelSelectUI) getLevelSetByDifficulty(difficulty levels.Difficulty) *levels.LevelSet {
	for _, levelSet := range lsui.levelManager.LevelSets {
		if levelSet.Difficulty == difficulty {
			return levelSet
		}
	}
	return nil
}

func (lsui *LevelSelectUI) isDifficultyUnlocked(levelSet *levels.LevelSet) bool {
	if levelSet == nil {
		return false
	}
	
	// Count completed levels
	completedCount := 0
	for _, set := range lsui.levelManager.LevelSets {
		for _, level := range set.Levels {
			if level.Completed {
				completedCount++
			}
		}
	}
	
	return completedCount >= levelSet.UnlockLevel
}