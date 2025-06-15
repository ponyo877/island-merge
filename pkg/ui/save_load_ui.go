package ui

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ponyo877/island-merge/pkg/storage"
)

type SaveLoadUI struct {
	saveSystem    *storage.SaveSystem
	showPanel     bool
	selectedTab   int // 0: Save/Load, 1: Settings, 2: Import/Export
	settings      *storage.GameSettings
	statusMessage string
	statusTime    time.Time
	OnSaveGame    func()
	OnLoadGame    func()
}

func NewSaveLoadUI(saveSystem *storage.SaveSystem) *SaveLoadUI {
	settings, _ := saveSystem.LoadSettings()
	return &SaveLoadUI{
		saveSystem:  saveSystem,
		showPanel:   false,
		selectedTab: 0,
		settings:    settings,
	}
}

func (slui *SaveLoadUI) TogglePanel() {
	slui.showPanel = !slui.showPanel
	if slui.showPanel {
		// Refresh settings when opening
		settings, _ := slui.saveSystem.LoadSettings()
		slui.settings = settings
	}
}

func (slui *SaveLoadUI) IsOpen() bool {
	return slui.showPanel
}

func (slui *SaveLoadUI) Update() {
	// Clear status message after 3 seconds
	if !slui.statusTime.IsZero() && time.Since(slui.statusTime) > 3*time.Second {
		slui.statusMessage = ""
		slui.statusTime = time.Time{}
	}
}

func (slui *SaveLoadUI) HandleClick(x, y int) bool {
	if !slui.showPanel {
		return false
	}
	
	// Panel bounds
	panelX, panelY := 120, 60
	panelWidth, panelHeight := 400, 360
	
	// Check if clicking outside panel
	if x < panelX || x > panelX+panelWidth || y < panelY || y > panelY+panelHeight {
		slui.showPanel = false
		return true
	}
	
	// Close button
	if x >= panelX+panelWidth-30 && x <= panelX+panelWidth-10 && y >= panelY+10 && y <= panelY+30 {
		slui.showPanel = false
		return true
	}
	
	// Tab buttons
	tabWidth := 120
	tabY := panelY + 40
	for i := 0; i < 3; i++ {
		tabX := panelX + 20 + i*tabWidth
		if x >= tabX && x <= tabX+tabWidth-10 && y >= tabY && y <= tabY+30 {
			slui.selectedTab = i
			return true
		}
	}
	
	// Tab-specific clicks
	switch slui.selectedTab {
	case 0:
		return slui.handleSaveLoadClick(x, y, panelX, panelY)
	case 1:
		return slui.handleSettingsClick(x, y, panelX, panelY)
	case 2:
		return slui.handleImportExportClick(x, y, panelX, panelY)
	}
	
	return true
}

func (slui *SaveLoadUI) handleSaveLoadClick(x, y, panelX, panelY int) bool {
	buttonY := panelY + 120
	buttonWidth, buttonHeight := 160, 40
	spacing := 20
	
	// Save Game button
	saveX := panelX + 30
	if x >= saveX && x <= saveX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight {
		slui.saveGame()
		return true
	}
	
	// Load Game button
	loadX := saveX + buttonWidth + spacing
	if x >= loadX && x <= loadX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight {
		slui.loadGame()
		return true
	}
	
	// Delete Save button
	deleteY := buttonY + buttonHeight + 20
	if x >= saveX && x <= saveX+buttonWidth && y >= deleteY && y <= deleteY+buttonHeight {
		slui.deleteSave()
		return true
	}
	
	// Auto-save toggle
	autoSaveY := deleteY + buttonHeight + 20
	if x >= saveX && x <= saveX+20 && y >= autoSaveY && y <= autoSaveY+20 {
		slui.settings.AutoSave = !slui.settings.AutoSave
		slui.saveSystem.SaveSettings(slui.settings)
		return true
	}
	
	return true
}

func (slui *SaveLoadUI) handleSettingsClick(x, y, panelX, panelY int) bool {
	startY := panelY + 100
	checkboxSize := 20
	spacing := 30
	
	checkboxes := []struct {
		setting *bool
		y       int
	}{
		{&slui.settings.SoundEnabled, startY},
		{&slui.settings.MusicEnabled, startY + spacing},
		{&slui.settings.ShowTutorial, startY + spacing*2},
		{&slui.settings.AutoSave, startY + spacing*3},
	}
	
	checkboxX := panelX + 30
	for _, checkbox := range checkboxes {
		if x >= checkboxX && x <= checkboxX+checkboxSize && 
		   y >= checkbox.y && y <= checkbox.y+checkboxSize {
			*checkbox.setting = !*checkbox.setting
			slui.saveSystem.SaveSettings(slui.settings)
			slui.showStatus("Settings saved!")
			return true
		}
	}
	
	// Animation speed slider (simplified - just buttons)
	sliderY := startY + spacing*4
	slowButtonX := checkboxX
	fastButtonX := checkboxX + 100
	
	if y >= sliderY && y <= sliderY+20 {
		if x >= slowButtonX && x <= slowButtonX+40 {
			slui.settings.AnimationSpeed = 0.5
			slui.saveSystem.SaveSettings(slui.settings)
			slui.showStatus("Animation speed: Slow")
			return true
		}
		if x >= fastButtonX && x <= fastButtonX+40 {
			slui.settings.AnimationSpeed = 2.0
			slui.saveSystem.SaveSettings(slui.settings)
			slui.showStatus("Animation speed: Fast")
			return true
		}
	}
	
	return true
}

func (slui *SaveLoadUI) handleImportExportClick(x, y, panelX, panelY int) bool {
	buttonY := panelY + 120
	buttonWidth, buttonHeight := 160, 40
	spacing := 20
	
	// Export button
	exportX := panelX + 30
	if x >= exportX && x <= exportX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight {
		slui.exportData()
		return true
	}
	
	// Clear Data button
	clearY := buttonY + buttonHeight + spacing
	if x >= exportX && x <= exportX+buttonWidth && y >= clearY && y <= clearY+buttonHeight {
		slui.clearAllData()
		return true
	}
	
	return true
}

func (slui *SaveLoadUI) saveGame() {
	// Signal to main game to save
	if slui.OnSaveGame != nil {
		slui.OnSaveGame()
	}
	slui.showStatus("Game saved!")
}

func (slui *SaveLoadUI) loadGame() {
	if slui.saveSystem.HasSavedGame() {
		// Signal to main game to load
		if slui.OnLoadGame != nil {
			slui.OnLoadGame()
		}
		slui.showStatus("Game loaded!")
	} else {
		slui.showStatus("No saved game found!")
	}
}

func (slui *SaveLoadUI) deleteSave() {
	slui.saveSystem.DeleteSavedGame()
	slui.showStatus("Save deleted!")
}

func (slui *SaveLoadUI) exportData() {
	// In a real implementation, this would create a download or copy to clipboard
	slui.showStatus("Data exported to console!")
	fmt.Println("Exporting save data...")
	// This is where we'd implement actual export functionality
}

func (slui *SaveLoadUI) clearAllData() {
	slui.saveSystem.ClearAllData()
	slui.showStatus("All data cleared!")
}

func (slui *SaveLoadUI) showStatus(message string) {
	slui.statusMessage = message
	slui.statusTime = time.Now()
}

func (slui *SaveLoadUI) Draw(screen *ebiten.Image) {
	if !slui.showPanel {
		return
	}
	
	// Dark overlay
	overlay := ebiten.NewImage(640, 480)
	overlay.Fill(color.RGBA{0, 0, 0, 128})
	screen.DrawImage(overlay, nil)
	
	// Panel background
	panelX, panelY := 120, 60
	panelWidth, panelHeight := 400, 360
	
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
	ebitenutil.DebugPrintAt(screen, "Game Settings", panelX+20, panelY+15)
	
	// Close button
	vector.DrawFilledRect(screen, float32(panelX+panelWidth-30), float32(panelY+10), 20, 20, color.RGBA{200, 100, 100, 255}, false)
	ebitenutil.DebugPrintAt(screen, "X", panelX+panelWidth-25, panelY+15)
	
	// Draw tabs
	slui.drawTabs(screen, panelX, panelY)
	
	// Draw tab content
	switch slui.selectedTab {
	case 0:
		slui.drawSaveLoadTab(screen, panelX, panelY)
	case 1:
		slui.drawSettingsTab(screen, panelX, panelY)
	case 2:
		slui.drawImportExportTab(screen, panelX, panelY)
	}
	
	// Status message
	if slui.statusMessage != "" {
		statusY := panelY + panelHeight - 30
		ebitenutil.DebugPrintAt(screen, slui.statusMessage, panelX+20, statusY)
	}
}

func (slui *SaveLoadUI) drawTabs(screen *ebiten.Image, panelX, panelY int) {
	tabs := []string{"Save/Load", "Settings", "Data"}
	tabWidth := 120
	tabHeight := 30
	tabY := panelY + 40
	
	for i, tabName := range tabs {
		tabX := panelX + 20 + i*tabWidth
		
		// Tab background
		bgColor := color.RGBA{200, 200, 200, 255}
		if i == slui.selectedTab {
			bgColor = color.RGBA{150, 150, 250, 255}
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
		textX := tabX + (tabWidth-len(tabName)*6)/2
		textY := tabY + tabHeight/2 - 4
		ebitenutil.DebugPrintAt(screen, tabName, textX, textY)
	}
}

func (slui *SaveLoadUI) drawSaveLoadTab(screen *ebiten.Image, panelX, panelY int) {
	startY := panelY + 90
	
	// Info text
	ebitenutil.DebugPrintAt(screen, "Game Save Management", panelX+20, startY)
	
	hasSave := slui.saveSystem.HasSavedGame()
	saveStatus := "No saved game"
	if hasSave {
		saveStatus = "Saved game available"
	}
	ebitenutil.DebugPrintAt(screen, saveStatus, panelX+20, startY+20)
	
	// Buttons
	buttonY := panelY + 120
	buttonWidth, buttonHeight := 160, 40
	spacing := 20
	
	// Save Game button
	slui.drawButton(screen, panelX+30, buttonY, buttonWidth, buttonHeight, "Save Game", color.RGBA{100, 200, 100, 255})
	
	// Load Game button
	loadColor := color.RGBA{100, 100, 200, 255}
	if !hasSave {
		loadColor = color.RGBA{150, 150, 150, 255} // Disabled
	}
	slui.drawButton(screen, panelX+30+buttonWidth+spacing, buttonY, buttonWidth, buttonHeight, "Load Game", loadColor)
	
	// Delete Save button
	deleteY := buttonY + buttonHeight + 20
	deleteColor := color.RGBA{200, 100, 100, 255}
	if !hasSave {
		deleteColor = color.RGBA{150, 150, 150, 255} // Disabled
	}
	slui.drawButton(screen, panelX+30, deleteY, buttonWidth, buttonHeight, "Delete Save", deleteColor)
	
	// Auto-save checkbox
	autoSaveY := deleteY + buttonHeight + 20
	slui.drawCheckbox(screen, panelX+30, autoSaveY, slui.settings.AutoSave, "Auto-save enabled")
}

func (slui *SaveLoadUI) drawSettingsTab(screen *ebiten.Image, panelX, panelY int) {
	startY := panelY + 90
	
	ebitenutil.DebugPrintAt(screen, "Game Settings", panelX+20, startY)
	
	checkboxY := startY + 30
	spacing := 30
	
	// Sound settings
	slui.drawCheckbox(screen, panelX+30, checkboxY, slui.settings.SoundEnabled, "Sound Effects")
	slui.drawCheckbox(screen, panelX+30, checkboxY+spacing, slui.settings.MusicEnabled, "Background Music")
	slui.drawCheckbox(screen, panelX+30, checkboxY+spacing*2, slui.settings.ShowTutorial, "Show Tutorial")
	slui.drawCheckbox(screen, panelX+30, checkboxY+spacing*3, slui.settings.AutoSave, "Auto-save")
	
	// Animation speed
	speedY := checkboxY + spacing*4
	ebitenutil.DebugPrintAt(screen, "Animation Speed:", panelX+30, speedY)
	
	// Speed buttons
	slowColor := color.RGBA{150, 150, 150, 255}
	if slui.settings.AnimationSpeed == 0.5 {
		slowColor = color.RGBA{100, 200, 100, 255}
	}
	slui.drawButton(screen, panelX+30, speedY+20, 40, 20, "Slow", slowColor)
	
	normalColor := color.RGBA{150, 150, 150, 255}
	if slui.settings.AnimationSpeed == 1.0 {
		normalColor = color.RGBA{100, 200, 100, 255}
	}
	slui.drawButton(screen, panelX+80, speedY+20, 50, 20, "Normal", normalColor)
	
	fastColor := color.RGBA{150, 150, 150, 255}
	if slui.settings.AnimationSpeed == 2.0 {
		fastColor = color.RGBA{100, 200, 100, 255}
	}
	slui.drawButton(screen, panelX+140, speedY+20, 40, 20, "Fast", fastColor)
}

func (slui *SaveLoadUI) drawImportExportTab(screen *ebiten.Image, panelX, panelY int) {
	startY := panelY + 90
	
	ebitenutil.DebugPrintAt(screen, "Data Management", panelX+20, startY)
	
	// Storage usage
	usage := slui.saveSystem.GetStorageUsage()
	infoY := startY + 30
	for key, exists := range usage {
		status := "❌"
		if exists {
			status = "✅"
		}
		text := fmt.Sprintf("%s %s", status, key)
		ebitenutil.DebugPrintAt(screen, text, panelX+30, infoY)
		infoY += 15
	}
	
	// Buttons
	buttonY := panelY + 120
	buttonWidth, buttonHeight := 160, 40
	spacing := 20
	
	slui.drawButton(screen, panelX+30, buttonY, buttonWidth, buttonHeight, "Export Data", color.RGBA{100, 200, 200, 255})
	
	clearY := buttonY + buttonHeight + spacing
	slui.drawButton(screen, panelX+30, clearY, buttonWidth, buttonHeight, "Clear All Data", color.RGBA{200, 100, 100, 255})
}

func (slui *SaveLoadUI) drawButton(screen *ebiten.Image, x, y, width, height int, text string, bgColor color.Color) {
	vector.DrawFilledRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		bgColor,
		false,
	)
	
	vector.StrokeRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		2,
		color.RGBA{100, 100, 100, 255},
		false,
	)
	
	textX := x + (width-len(text)*6)/2
	textY := y + height/2 - 4
	ebitenutil.DebugPrintAt(screen, text, textX, textY)
}

func (slui *SaveLoadUI) drawCheckbox(screen *ebiten.Image, x, y int, checked bool, label string) {
	size := 20
	
	// Checkbox background
	bgColor := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(size), float32(size), bgColor, false)
	
	// Checkbox border
	vector.StrokeRect(screen, float32(x), float32(y), float32(size), float32(size), 2, color.RGBA{100, 100, 100, 255}, false)
	
	// Check mark
	if checked {
		ebitenutil.DebugPrintAt(screen, "✓", x+4, y+4)
	}
	
	// Label
	ebitenutil.DebugPrintAt(screen, label, x+size+10, y+6)
}

func (slui *SaveLoadUI) DrawSettingsButton(screen *ebiten.Image, x, y float64) {
	width, height := 100.0, 30.0
	
	vector.DrawFilledRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		color.RGBA{200, 200, 200, 255},
		false,
	)
	
	vector.StrokeRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		2,
		color.RGBA{150, 150, 150, 255},
		false,
	)
	
	ebitenutil.DebugPrintAt(screen, "⚙️ Settings", int(x+10), int(y+10))
}

func (slui *SaveLoadUI) IsSettingsButtonClicked(x, y int) bool {
	return x >= 10 && x <= 110 && y >= 10 && y <= 40
}