package ui

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ponyo877/island-merge/pkg/achievements"
)

type AchievementNotification struct {
	Achievement *achievements.Achievement
	StartTime   time.Time
	Duration    time.Duration
	Y           float64
}

type AchievementsUI struct {
	achievementSystem *achievements.AchievementSystem
	notifications     []*AchievementNotification
	showPanel         bool
	panelScroll       float64
}

func NewAchievementsUI(system *achievements.AchievementSystem) *AchievementsUI {
	ui := &AchievementsUI{
		achievementSystem: system,
		notifications:     make([]*AchievementNotification, 0),
		showPanel:         false,
	}
	
	// Listen for new achievements
	system.OnAchievementUnlocked(ui.onAchievementUnlocked)
	
	return ui
}

func (aui *AchievementsUI) onAchievementUnlocked(achievement *achievements.Achievement) {
	notification := &AchievementNotification{
		Achievement: achievement,
		StartTime:   time.Now(),
		Duration:    time.Second * 4,
		Y:           -100, // Start off-screen
	}
	
	aui.notifications = append(aui.notifications, notification)
}

func (aui *AchievementsUI) Update() {
	now := time.Now()
	
	// Update notifications
	activeNotifications := make([]*AchievementNotification, 0)
	for _, notification := range aui.notifications {
		elapsed := now.Sub(notification.StartTime)
		
		if elapsed < notification.Duration {
			// Animate notification sliding in and out
			progress := float64(elapsed) / float64(notification.Duration)
			
			if progress < 0.2 {
				// Slide in
				slideProgress := progress / 0.2
				notification.Y = -100 + slideProgress*120 // Slide to Y=20
			} else if progress < 0.8 {
				// Stay visible
				notification.Y = 20
			} else {
				// Slide out
				slideProgress := (progress - 0.8) / 0.2
				notification.Y = 20 - slideProgress*120 // Slide to Y=-100
			}
			
			activeNotifications = append(activeNotifications, notification)
		}
	}
	
	aui.notifications = activeNotifications
}

func (aui *AchievementsUI) TogglePanel() {
	aui.showPanel = !aui.showPanel
	aui.panelScroll = 0
}

func (aui *AchievementsUI) HandleScroll(deltaY float64) {
	if aui.showPanel {
		aui.panelScroll += deltaY * 20
		aui.panelScroll = math.Max(0, aui.panelScroll)
	}
}

func (aui *AchievementsUI) HandleClick(x, y int) bool {
	if !aui.showPanel {
		return false
	}
	
	// Check if clicking close button
	if x >= 580 && x <= 620 && y >= 20 && y <= 60 {
		aui.showPanel = false
		return true
	}
	
	return true // Consume click when panel is open
}

func (aui *AchievementsUI) Draw(screen *ebiten.Image) {
	// Draw notifications
	aui.drawNotifications(screen)
	
	// Draw achievements panel if open
	if aui.showPanel {
		aui.drawAchievementsPanel(screen)
	}
}

func (aui *AchievementsUI) drawNotifications(screen *ebiten.Image) {
	for _, notification := range aui.notifications {
		aui.drawNotification(screen, notification)
	}
}

func (aui *AchievementsUI) drawNotification(screen *ebiten.Image, notification *AchievementNotification) {
	x := 50.0
	y := notification.Y
	width := 300.0
	height := 60.0
	
	// Background with glow effect
	glowColor := color.RGBA{255, 215, 0, 100} // Gold glow
	for i := 0; i < 3; i++ {
		vector.DrawFilledRect(
			screen,
			float32(x-float64(i)*2), float32(y-float64(i)*2),
			float32(width+float64(i)*4), float32(height+float64(i)*4),
			glowColor,
			false,
		)
	}
	
	// Main background
	vector.DrawFilledRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		color.RGBA{0, 0, 0, 200},
		false,
	)
	
	// Border
	vector.StrokeRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		2,
		color.RGBA{255, 215, 0, 255},
		false,
	)
	
	// Achievement unlocked text
	ebitenutil.DebugPrintAt(screen, "Achievement Unlocked!", int(x+10), int(y+10))
	
	// Achievement name and icon
	nameText := fmt.Sprintf("%s %s", notification.Achievement.Icon, notification.Achievement.Name)
	ebitenutil.DebugPrintAt(screen, nameText, int(x+10), int(y+25))
	
	// Description
	ebitenutil.DebugPrintAt(screen, notification.Achievement.Description, int(x+10), int(y+40))
}

func (aui *AchievementsUI) drawAchievementsPanel(screen *ebiten.Image) {
	// Panel background
	panelX := 100.0
	panelY := 50.0
	panelWidth := 440.0
	panelHeight := 380.0
	
	// Dark background overlay
	overlay := ebiten.NewImage(640, 480)
	overlay.Fill(color.RGBA{0, 0, 0, 128})
	screen.DrawImage(overlay, nil)
	
	// Panel background
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
	ebitenutil.DebugPrintAt(screen, "Achievements", int(panelX+20), int(panelY+20))
	
	// Close button
	vector.DrawFilledRect(screen, 580, 20, 40, 40, color.RGBA{200, 100, 100, 255}, false)
	ebitenutil.DebugPrintAt(screen, "X", 595, 35)
	
	// Progress summary
	summary := aui.achievementSystem.GetProgressSummary()
	ebitenutil.DebugPrintAt(screen, summary, int(panelX+20), int(panelY+40))
	
	// Achievement list
	achievements := aui.achievementSystem.GetAchievements()
	startY := panelY + 70 - aui.panelScroll
	
	for i, achievement := range achievements {
		itemY := startY + float64(i*70)
		
		// Skip if outside visible area
		if itemY < panelY+60 || itemY > panelY+panelHeight-10 {
			continue
		}
		
		aui.drawAchievementItem(screen, achievement, panelX+10, itemY, panelWidth-20)
	}
}

func (aui *AchievementsUI) drawAchievementItem(screen *ebiten.Image, achievement *achievements.Achievement, x, y, width float64) {
	height := 60.0
	
	// Background color based on unlock status
	bgColor := color.RGBA{200, 200, 200, 255}
	if achievement.Unlocked {
		bgColor = color.RGBA{144, 238, 144, 255} // Light green
	}
	
	vector.DrawFilledRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		bgColor,
		false,
	)
	
	// Border
	borderColor := color.RGBA{150, 150, 150, 255}
	if achievement.Unlocked {
		borderColor = color.RGBA{255, 215, 0, 255} // Gold
	}
	
	vector.StrokeRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		2,
		borderColor,
		false,
	)
	
	// Icon and name
	nameText := fmt.Sprintf("%s %s", achievement.Icon, achievement.Name)
	ebitenutil.DebugPrintAt(screen, nameText, int(x+10), int(y+10))
	
	// Description
	ebitenutil.DebugPrintAt(screen, achievement.Description, int(x+10), int(y+25))
	
	// Progress bar
	if !achievement.Unlocked && achievement.Target > 1 {
		progressText := fmt.Sprintf("Progress: %d/%d", achievement.Progress, achievement.Target)
		ebitenutil.DebugPrintAt(screen, progressText, int(x+10), int(y+40))
		
		// Progress bar
		barWidth := 200.0
		barHeight := 8.0
		barX := x + width - barWidth - 10
		barY := y + 45
		
		// Background
		vector.DrawFilledRect(
			screen,
			float32(barX), float32(barY),
			float32(barWidth), float32(barHeight),
			color.RGBA{100, 100, 100, 255},
			false,
		)
		
		// Progress
		progress := float64(achievement.Progress) / float64(achievement.Target)
		progressWidth := barWidth * math.Min(1.0, progress)
		
		vector.DrawFilledRect(
			screen,
			float32(barX), float32(barY),
			float32(progressWidth), float32(barHeight),
			color.RGBA{0, 200, 0, 255},
			false,
		)
	} else if achievement.Unlocked {
		ebitenutil.DebugPrintAt(screen, "UNLOCKED", int(x+width-80), int(y+40))
	}
}

func (aui *AchievementsUI) DrawAchievementButton(screen *ebiten.Image, x, y float64) {
	width := 120.0
	height := 30.0
	
	// Button background
	vector.DrawFilledRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		color.RGBA{255, 215, 0, 255},
		false,
	)
	
	// Button border
	vector.StrokeRect(
		screen,
		float32(x), float32(y),
		float32(width), float32(height),
		2,
		color.RGBA{200, 170, 0, 255},
		false,
	)
	
	// Button text
	unlocked := aui.achievementSystem.GetUnlockedCount()
	total := aui.achievementSystem.GetTotalCount()
	buttonText := fmt.Sprintf("🏆 %d/%d", unlocked, total)
	ebitenutil.DebugPrintAt(screen, buttonText, int(x+10), int(y+10))
}

func (aui *AchievementsUI) IsAchievementButtonClicked(x, y int) bool {
	// Check if clicking achievement button (positioned at top right)
	return x >= 500 && x <= 620 && y >= 10 && y <= 40
}