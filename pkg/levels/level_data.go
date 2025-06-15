package levels

import (
	"time"

	"github.com/ponyo877/island-merge/pkg/island"
)

type Difficulty int

const (
	DifficultyBeginner Difficulty = iota
	DifficultyIntermediate
	DifficultyExpert
	DifficultyMaster
)

type LevelData struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Difficulty  Difficulty            `json:"difficulty"`
	Width       int                   `json:"width"`
	Height      int                   `json:"height"`
	Grid        [][]island.TileType   `json:"grid"`
	OptimalMoves int                  `json:"optimal_moves"`
	TimeLimit   time.Duration         `json:"time_limit,omitempty"`
	Objectives  []Objective           `json:"objectives"`
	Unlocked    bool                  `json:"unlocked"`
	Completed   bool                  `json:"completed"`
	BestScore   *Score               `json:"best_score,omitempty"`
}

type Objective struct {
	Type        string `json:"type"`        // "connect_all", "min_bridges", "time_limit"
	Target      int    `json:"target"`
	Description string `json:"description"`
}

type Score struct {
	Moves     int           `json:"moves"`
	Time      time.Duration `json:"time"`
	Stars     int           `json:"stars"` // 1-3 stars based on performance
	Date      time.Time     `json:"date"`
}

type LevelSet struct {
	Name        string       `json:"name"`
	Difficulty  Difficulty   `json:"difficulty"`
	Description string       `json:"description"`
	Levels      []*LevelData `json:"levels"`
	UnlockLevel int          `json:"unlock_level"` // Level required to unlock this set
}

// Level manager handles all level data
type LevelManager struct {
	LevelSets    []*LevelSet         `json:"level_sets"`
	CurrentLevel *LevelData          `json:"current_level,omitempty"`
	Progress     map[string]*Score   `json:"progress"` // levelID -> best score
}

func NewLevelManager() *LevelManager {
	lm := &LevelManager{
		LevelSets: make([]*LevelSet, 0),
		Progress:  make(map[string]*Score),
	}
	
	lm.initializeDefaultLevels()
	return lm
}

func (lm *LevelManager) initializeDefaultLevels() {
	// Beginner levels (5x5 to 8x8)
	beginnerSet := &LevelSet{
		Name:        "Island Basics",
		Difficulty:  DifficultyBeginner,
		Description: "Learn the fundamentals of island connecting",
		UnlockLevel: 0,
		Levels:      make([]*LevelData, 0),
	}
	
	// Add beginner levels
	beginnerSet.Levels = append(beginnerSet.Levels, lm.createBeginnerLevels()...)
	
	// Intermediate levels (10x10 to 15x15)
	intermediateSet := &LevelSet{
		Name:        "Island Chains",
		Difficulty:  DifficultyIntermediate,
		Description: "More complex island arrangements",
		UnlockLevel: 3, // Unlock after completing 3 beginner levels
		Levels:      make([]*LevelData, 0),
	}
	
	intermediateSet.Levels = append(intermediateSet.Levels, lm.createIntermediateLevels()...)
	
	// Expert levels (20x20 to 25x25)
	expertSet := &LevelSet{
		Name:        "Island Archipelago",
		Difficulty:  DifficultyExpert,
		Description: "Master the art of large-scale connecting",
		UnlockLevel: 8, // Unlock after completing intermediate levels
		Levels:      make([]*LevelData, 0),
	}
	
	expertSet.Levels = append(expertSet.Levels, lm.createExpertLevels()...)
	
	// Master levels (various challenging configurations)
	masterSet := &LevelSet{
		Name:        "Island Master",
		Difficulty:  DifficultyMaster,
		Description: "Ultimate challenges for true masters",
		UnlockLevel: 15,
		Levels:      make([]*LevelData, 0),
	}
	
	masterSet.Levels = append(masterSet.Levels, lm.createMasterLevels()...)
	
	lm.LevelSets = []*LevelSet{beginnerSet, intermediateSet, expertSet, masterSet}
	
	// Mark first level as unlocked
	if len(beginnerSet.Levels) > 0 {
		beginnerSet.Levels[0].Unlocked = true
	}
}

func (lm *LevelManager) createBeginnerLevels() []*LevelData {
	levels := make([]*LevelData, 0)
	
	// Level 1: Simple 3-island triangle (5x5)
	level1 := &LevelData{
		ID:          "beginner_01",
		Name:        "First Steps",
		Description: "Connect three islands in a simple triangle",
		Difficulty:  DifficultyBeginner,
		Width:       5,
		Height:      5,
		OptimalMoves: 2,
		Objectives: []Objective{
			{Type: "connect_all", Target: 1, Description: "Connect all islands"},
		},
	}
	level1.Grid = lm.createGrid(5, 5, [][]int{
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0},
	})
	levels = append(levels, level1)
	
	// Level 2: Four corners (6x6)
	level2 := &LevelData{
		ID:          "beginner_02",
		Name:        "Four Corners",
		Description: "Islands at each corner need connecting",
		Difficulty:  DifficultyBeginner,
		Width:       6,
		Height:      6,
		OptimalMoves: 5,
		Objectives: []Objective{
			{Type: "connect_all", Target: 1, Description: "Connect all corner islands"},
		},
	}
	level2.Grid = lm.createGrid(6, 6, [][]int{
		{1, 0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 1},
	})
	levels = append(levels, level2)
	
	// Level 3: Cross pattern (7x7)
	level3 := &LevelData{
		ID:          "beginner_03",
		Name:        "Island Cross",
		Description: "Connect islands arranged in a cross pattern",
		Difficulty:  DifficultyBeginner,
		Width:       7,
		Height:      7,
		OptimalMoves: 4,
		Objectives: []Objective{
			{Type: "connect_all", Target: 1, Description: "Connect all islands"},
			{Type: "min_bridges", Target: 4, Description: "Use minimum bridges"},
		},
	}
	level3.Grid = lm.createGrid(7, 7, [][]int{
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 1, 0, 0, 1},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
	})
	levels = append(levels, level3)
	
	// Level 4: Circle formation (8x8)
	level4 := &LevelData{
		ID:          "beginner_04",
		Name:        "Island Circle",
		Description: "Islands forming a circle - find the optimal path",
		Difficulty:  DifficultyBeginner,
		Width:       8,
		Height:      8,
		OptimalMoves: 6,
		Objectives: []Objective{
			{Type: "connect_all", Target: 1, Description: "Connect all islands"},
		},
	}
	level4.Grid = lm.createGrid(8, 8, [][]int{
		{0, 0, 0, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 1, 0, 0, 0},
	})
	levels = append(levels, level4)
	
	return levels
}

func (lm *LevelManager) createIntermediateLevels() []*LevelData {
	levels := make([]*LevelData, 0)
	
	// Level 5: Scattered islands (10x10)
	level5 := &LevelData{
		ID:          "intermediate_01",
		Name:        "Scattered Isles",
		Description: "Many small islands scattered across the sea",
		Difficulty:  DifficultyIntermediate,
		Width:       10,
		Height:      10,
		OptimalMoves: 8,
		TimeLimit:   time.Minute * 3,
		Objectives: []Objective{
			{Type: "connect_all", Target: 1, Description: "Connect all islands"},
			{Type: "time_limit", Target: 180, Description: "Complete within 3 minutes"},
		},
	}
	level5.Grid = lm.createGrid(10, 10, [][]int{
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{0, 0, 0, 1, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 1, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	})
	levels = append(levels, level5)
	
	// Level 6: Maze-like (12x12)
	level6 := &LevelData{
		ID:          "intermediate_02",
		Name:        "Island Maze",
		Description: "Navigate through a maze of islands",
		Difficulty:  DifficultyIntermediate,
		Width:       12,
		Height:      12,
		OptimalMoves: 12,
		Objectives: []Objective{
			{Type: "connect_all", Target: 1, Description: "Connect all islands"},
			{Type: "min_bridges", Target: 12, Description: "Find the optimal path"},
		},
	}
	// Create a maze-like pattern
	mazePattern := make([][]int, 12)
	for i := range mazePattern {
		mazePattern[i] = make([]int, 12)
	}
	// Add strategic island placements
	mazePattern[1][1] = 1
	mazePattern[1][5] = 1
	mazePattern[1][10] = 1
	mazePattern[5][1] = 1
	mazePattern[5][5] = 1
	mazePattern[5][10] = 1
	mazePattern[10][1] = 1
	mazePattern[10][5] = 1
	mazePattern[10][10] = 1
	
	level6.Grid = lm.createGrid(12, 12, mazePattern)
	levels = append(levels, level6)
	
	// Level 7: Dense cluster (15x15)
	level7 := &LevelData{
		ID:          "intermediate_03",
		Name:        "Dense Archipelago",
		Description: "Many islands clustered together",
		Difficulty:  DifficultyIntermediate,
		Width:       15,
		Height:      15,
		OptimalMoves: 15,
		Objectives: []Objective{
			{Type: "connect_all", Target: 1, Description: "Connect all islands"},
		},
	}
	
	// Create a dense cluster pattern
	clusterPattern := make([][]int, 15)
	for i := range clusterPattern {
		clusterPattern[i] = make([]int, 15)
	}
	
	// Add islands in clusters
	clusters := [][]int{
		{2, 2}, {3, 2}, {2, 3},           // Top-left cluster
		{12, 2}, {13, 2}, {12, 3},        // Top-right cluster
		{7, 7}, {8, 7}, {7, 8}, {8, 8},   // Center cluster
		{2, 12}, {3, 12}, {2, 13},        // Bottom-left cluster
		{12, 12}, {13, 12}, {12, 13},     // Bottom-right cluster
	}
	
	for _, pos := range clusters {
		if pos[0] < 15 && pos[1] < 15 {
			clusterPattern[pos[1]][pos[0]] = 1
		}
	}
	
	level7.Grid = lm.createGrid(15, 15, clusterPattern)
	levels = append(levels, level7)
	
	return levels
}

func (lm *LevelManager) createExpertLevels() []*LevelData {
	levels := make([]*LevelData, 0)
	
	// Level 8: Large spiral (20x20)
	level8 := &LevelData{
		ID:          "expert_01",
		Name:        "Spiral Galaxy",
		Description: "Islands arranged in a vast spiral pattern",
		Difficulty:  DifficultyExpert,
		Width:       20,
		Height:      20,
		OptimalMoves: 25,
		TimeLimit:   time.Minute * 5,
		Objectives: []Objective{
			{Type: "connect_all", Target: 1, Description: "Connect all islands"},
			{Type: "time_limit", Target: 300, Description: "Complete within 5 minutes"},
		},
	}
	
	spiralPattern := lm.createSpiralPattern(20, 20)
	level8.Grid = lm.createGrid(20, 20, spiralPattern)
	levels = append(levels, level8)
	
	// Level 9: Maximum size challenge (25x25)
	level9 := &LevelData{
		ID:          "expert_02",
		Name:        "Continental Drift",
		Description: "The ultimate island connecting challenge",
		Difficulty:  DifficultyExpert,
		Width:       25,
		Height:      25,
		OptimalMoves: 35,
		TimeLimit:   time.Minute * 8,
		Objectives: []Objective{
			{Type: "connect_all", Target: 1, Description: "Connect all continents"},
			{Type: "time_limit", Target: 480, Description: "Complete within 8 minutes"},
			{Type: "min_bridges", Target: 35, Description: "Achieve optimal efficiency"},
		},
	}
	
	continentalPattern := lm.createContinentalPattern(25, 25)
	level9.Grid = lm.createGrid(25, 25, continentalPattern)
	levels = append(levels, level9)
	
	return levels
}

func (lm *LevelManager) createMasterLevels() []*LevelData {
	levels := make([]*LevelData, 0)
	
	// Master level: Symmetric beauty (20x20)
	master1 := &LevelData{
		ID:          "master_01",
		Name:        "Perfect Symmetry",
		Description: "A perfectly symmetric island arrangement",
		Difficulty:  DifficultyMaster,
		Width:       20,
		Height:      20,
		OptimalMoves: 18,
		TimeLimit:   time.Minute * 4,
		Objectives: []Objective{
			{Type: "connect_all", Target: 1, Description: "Connect all islands"},
			{Type: "min_bridges", Target: 18, Description: "Perfect efficiency required"},
		},
	}
	
	symmetricPattern := lm.createSymmetricPattern(20, 20)
	master1.Grid = lm.createGrid(20, 20, symmetricPattern)
	levels = append(levels, master1)
	
	return levels
}

// Helper functions to create patterns
func (lm *LevelManager) createGrid(width, height int, pattern [][]int) [][]island.TileType {
	grid := make([][]island.TileType, height)
	for y := range grid {
		grid[y] = make([]island.TileType, width)
		for x := range grid[y] {
			if y < len(pattern) && x < len(pattern[y]) && pattern[y][x] == 1 {
				grid[y][x] = island.TileLand
			} else {
				grid[y][x] = island.TileSea
			}
		}
	}
	return grid
}

func (lm *LevelManager) createSpiralPattern(width, height int) [][]int {
	pattern := make([][]int, height)
	for i := range pattern {
		pattern[i] = make([]int, width)
	}
	
	centerX, centerY := width/2, height/2
	radius := 2
	
	for angle := 0; angle < 720; angle += 30 { // Two full rotations
		x := centerX + int(float64(radius)*0.1*float64(angle)*0.017453292519943295) // Convert to radians
		y := centerY + int(float64(radius)*0.1*float64(angle)*0.017453292519943295)
		
		if x >= 0 && x < width && y >= 0 && y < height {
			pattern[y][x] = 1
		}
	}
	
	return pattern
}

func (lm *LevelManager) createContinentalPattern(width, height int) [][]int {
	pattern := make([][]int, height)
	for i := range pattern {
		pattern[i] = make([]int, width)
	}
	
	// Create several "continents" - large clusters of islands
	continents := []struct {
		centerX, centerY, size int
	}{
		{6, 6, 3},     // Top-left continent
		{18, 6, 4},    // Top-right continent
		{6, 18, 3},    // Bottom-left continent
		{18, 18, 4},   // Bottom-right continent
		{12, 12, 2},   // Central island
	}
	
	for _, continent := range continents {
		for dy := -continent.size; dy <= continent.size; dy++ {
			for dx := -continent.size; dx <= continent.size; dx++ {
				x := continent.centerX + dx
				y := continent.centerY + dy
				
				if x >= 0 && x < width && y >= 0 && y < height {
					// Create irregular continent shape
					if dx*dx+dy*dy <= continent.size*continent.size {
						pattern[y][x] = 1
					}
				}
			}
		}
	}
	
	return pattern
}

func (lm *LevelManager) createSymmetricPattern(width, height int) [][]int {
	pattern := make([][]int, height)
	for i := range pattern {
		pattern[i] = make([]int, width)
	}
	
	// Create symmetric points
	points := []struct{ x, y int }{
		{3, 3}, {5, 2}, {8, 4}, {10, 7}, {12, 3},
	}
	
	for _, point := range points {
		// Place original point
		if point.x < width && point.y < height {
			pattern[point.y][point.x] = 1
		}
		
		// Place symmetric points (4-way symmetry)
		symX1 := width - 1 - point.x
		symY1 := height - 1 - point.y
		symX2 := point.x
		symY2 := height - 1 - point.y
		symX3 := width - 1 - point.x
		symY3 := point.y
		
		if symX1 >= 0 && symX1 < width && symY1 >= 0 && symY1 < height {
			pattern[symY1][symX1] = 1
		}
		if symX2 >= 0 && symX2 < width && symY2 >= 0 && symY2 < height {
			pattern[symY2][symX2] = 1
		}
		if symX3 >= 0 && symX3 < width && symY3 >= 0 && symY3 < height {
			pattern[symY3][symX3] = 1
		}
	}
	
	return pattern
}

// Level management methods
func (lm *LevelManager) GetLevelByID(id string) *LevelData {
	for _, levelSet := range lm.LevelSets {
		for _, level := range levelSet.Levels {
			if level.ID == id {
				return level
			}
		}
	}
	return nil
}

func (lm *LevelManager) UnlockNextLevel(completedLevelID string) {
	// Find the completed level and unlock the next one
	for _, levelSet := range lm.LevelSets {
		for i, level := range levelSet.Levels {
			if level.ID == completedLevelID {
				level.Completed = true
				
				// Unlock next level in same set
				if i+1 < len(levelSet.Levels) {
					levelSet.Levels[i+1].Unlocked = true
				}
				
				// Check if we should unlock next difficulty set
				lm.checkUnlockNextDifficulty()
				return
			}
		}
	}
}

func (lm *LevelManager) checkUnlockNextDifficulty() {
	completedCount := 0
	
	for _, levelSet := range lm.LevelSets {
		for _, level := range levelSet.Levels {
			if level.Completed {
				completedCount++
			}
		}
		
		// Check if next level set should be unlocked
		for _, nextSet := range lm.LevelSets {
			if nextSet.UnlockLevel <= completedCount {
				for _, level := range nextSet.Levels {
					if !level.Unlocked {
						level.Unlocked = true
						return // Unlock first level of next set
					}
				}
			}
		}
	}
}

func (lm *LevelManager) CalculateStars(level *LevelData, moves int, completionTime time.Duration) int {
	stars := 1 // Base completion star
	
	// Perfect moves = 3 stars
	if moves <= level.OptimalMoves {
		stars = 3
	} else if moves <= level.OptimalMoves+2 {
		stars = 2
	}
	
	// Time bonus (if there's a time limit)
	if level.TimeLimit > 0 {
		if completionTime <= level.TimeLimit/2 {
			stars = 3
		} else if completionTime <= level.TimeLimit*3/4 {
			stars = max(stars, 2)
		}
	}
	
	return stars
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}