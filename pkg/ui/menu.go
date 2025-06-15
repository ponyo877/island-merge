package ui

import (
	"image/color"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MenuItem struct {
	Text     string
	Action   func()
	X, Y     float64
	Width    float64
	Height   float64
	Hovered  bool
	Selected bool
}

type Menu struct {
	Title      string
	Items      []*MenuItem
	Background color.Color
}

func NewMainMenu(onModeSelect func(int)) *Menu {
	menu := &Menu{
		Title:      "Island Merge",
		Background: color.RGBA{240, 240, 240, 255},
		Items:      make([]*MenuItem, 0),
	}
	
	// Menu items
	items := []struct {
		text   string
		action func()
	}{
		{"Select Level", func() { onModeSelect(0) }}, // Level Select
		{"Time Attack", func() { onModeSelect(1) }}, // ModeTimeAttack
		{"Puzzle Mode", func() { onModeSelect(2) }}, // ModePuzzle
		{"Level Editor", func() { onModeSelect(3) }}, // Level Editor
	}
	
	startY := 200.0
	for i, item := range items {
		menuItem := &MenuItem{
			Text:   item.text,
			Action: item.action,
			X:      320 - 100, // Center
			Y:      startY + float64(i*60),
			Width:  200,
			Height: 40,
		}
		menu.Items = append(menu.Items, menuItem)
	}
	
	return menu
}

func (m *Menu) Update(mouseX, mouseY int, clicked bool) {
	for _, item := range m.Items {
		// Check hover
		item.Hovered = float64(mouseX) >= item.X && float64(mouseX) <= item.X+item.Width &&
			float64(mouseY) >= item.Y && float64(mouseY) <= item.Y+item.Height
		
		// Check click
		if item.Hovered && clicked && item.Action != nil {
			item.Action()
		}
	}
}

func (m *Menu) Draw(screen *ebiten.Image) {
	// Clear background
	screen.Fill(m.Background)
	
	// Draw title
	titleX := 320 - len(m.Title)*6 // Rough centering
	ebitenutil.DebugPrintAt(screen, m.Title, titleX, 100)
	
	// Draw menu items
	for _, item := range m.Items {
		// Background
		bgColor := color.RGBA{200, 200, 200, 255}
		if item.Hovered {
			bgColor = color.RGBA{150, 150, 250, 255}
		}
		
		vector.DrawFilledRect(
			screen,
			float32(item.X), float32(item.Y),
			float32(item.Width), float32(item.Height),
			bgColor,
			false,
		)
		
		// Border
		vector.StrokeRect(
			screen,
			float32(item.X), float32(item.Y),
			float32(item.Width), float32(item.Height),
			2,
			color.RGBA{100, 100, 100, 255},
			false,
		)
		
		// Text
		textX := int(item.X + item.Width/2 - float64(len(item.Text)*3))
		textY := int(item.Y + item.Height/2 - 4)
		ebitenutil.DebugPrintAt(screen, item.Text, textX, textY)
	}
}