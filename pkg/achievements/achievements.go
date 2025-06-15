package achievements

import (
	"encoding/json"
	"fmt"
	"time"
)

type AchievementType int

const (
	AchievementFirstWin AchievementType = iota
	AchievementSpeedrun
	AchievementEfficient
	AchievementTimeAttackWin
	AchievementPerfectGame
	AchievementBridgeBuilder
	AchievementIslandHopper
	AchievementLevelCreator
	AchievementDedicated
	AchievementMaster
)

type Achievement struct {
	ID          AchievementType `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Icon        string          `json:"icon"`
	Unlocked    bool            `json:"unlocked"`
	UnlockedAt  *time.Time      `json:"unlocked_at,omitempty"`
	Progress    int             `json:"progress"`
	Target      int             `json:"target"`
	Hidden      bool            `json:"hidden"`
}

type AchievementSystem struct {
	achievements map[AchievementType]*Achievement
	statistics   *GameStatistics
	listeners    []func(*Achievement)
}

type GameStatistics struct {
	GamesPlayed       int           `json:"games_played"`
	GamesWon          int           `json:"games_won"`
	TotalMoves        int           `json:"total_moves"`
	TotalTime         time.Duration `json:"total_time"`
	BestTime          time.Duration `json:"best_time"`
	FewestMoves       int           `json:"fewest_moves"`
	BridgesBuilt      int           `json:"bridges_built"`
	TimeAttackWins    int           `json:"time_attack_wins"`
	PerfectGames      int           `json:"perfect_games"`    // Games with minimum moves
	LevelsCreated     int           `json:"levels_created"`
	PlayStreak        int           `json:"play_streak"`
	LastPlayDate      *time.Time    `json:"last_play_date,omitempty"`
}

func NewAchievementSystem() *AchievementSystem {
	system := &AchievementSystem{
		achievements: make(map[AchievementType]*Achievement),
		statistics:   &GameStatistics{FewestMoves: 999}, // Initialize with high value
		listeners:    make([]func(*Achievement), 0),
	}
	
	system.initializeAchievements()
	return system
}

func (as *AchievementSystem) initializeAchievements() {
	achievements := []*Achievement{
		{
			ID:          AchievementFirstWin,
			Name:        "First Victory",
			Description: "Win your first game",
			Icon:        "🏆",
			Target:      1,
		},
		{
			ID:          AchievementSpeedrun,
			Name:        "Speed Demon",
			Description: "Complete a level in under 30 seconds",
			Icon:        "⚡",
			Target:      1,
		},
		{
			ID:          AchievementEfficient,
			Name:        "Efficiency Expert",
			Description: "Complete a level with minimum moves",
			Icon:        "🎯",
			Target:      1,
		},
		{
			ID:          AchievementTimeAttackWin,
			Name:        "Time Master",
			Description: "Win 5 Time Attack games",
			Icon:        "⏰",
			Target:      5,
		},
		{
			ID:          AchievementPerfectGame,
			Name:        "Perfectionist",
			Description: "Achieve 10 perfect games",
			Icon:        "💎",
			Target:      10,
		},
		{
			ID:          AchievementBridgeBuilder,
			Name:        "Bridge Builder",
			Description: "Build 100 bridges",
			Icon:        "🌉",
			Target:      100,
		},
		{
			ID:          AchievementIslandHopper,
			Name:        "Island Hopper",
			Description: "Win 25 games",
			Icon:        "🏝️",
			Target:      25,
		},
		{
			ID:          AchievementLevelCreator,
			Name:        "Level Designer",
			Description: "Create 5 levels in the editor",
			Icon:        "🎨",
			Target:      5,
		},
		{
			ID:          AchievementDedicated,
			Name:        "Dedicated Player",
			Description: "Play for 7 consecutive days",
			Icon:        "🔥",
			Target:      7,
		},
		{
			ID:          AchievementMaster,
			Name:        "Island Master",
			Description: "Unlock all other achievements",
			Icon:        "👑",
			Target:      9,
			Hidden:      true,
		},
	}
	
	for _, achievement := range achievements {
		as.achievements[achievement.ID] = achievement
	}
}

func (as *AchievementSystem) OnAchievementUnlocked(callback func(*Achievement)) {
	as.listeners = append(as.listeners, callback)
}

func (as *AchievementSystem) notifyListeners(achievement *Achievement) {
	for _, callback := range as.listeners {
		callback(achievement)
	}
}

func (as *AchievementSystem) checkAchievement(id AchievementType) {
	achievement := as.achievements[id]
	if achievement == nil || achievement.Unlocked {
		return
	}
	
	if achievement.Progress >= achievement.Target {
		achievement.Unlocked = true
		now := time.Now()
		achievement.UnlockedAt = &now
		as.notifyListeners(achievement)
		
		// Check master achievement
		as.checkMasterAchievement()
	}
}

func (as *AchievementSystem) checkMasterAchievement() {
	unlockedCount := 0
	for id, achievement := range as.achievements {
		if id != AchievementMaster && achievement.Unlocked {
			unlockedCount++
		}
	}
	
	master := as.achievements[AchievementMaster]
	if master != nil && !master.Unlocked {
		master.Progress = unlockedCount
		as.checkAchievement(AchievementMaster)
	}
}

// Game event handlers
func (as *AchievementSystem) OnGameStart() {
	as.statistics.GamesPlayed++
	
	// Update play streak
	now := time.Now()
	if as.statistics.LastPlayDate != nil {
		daysSince := int(now.Sub(*as.statistics.LastPlayDate).Hours() / 24)
		if daysSince == 1 {
			as.statistics.PlayStreak++
		} else if daysSince > 1 {
			as.statistics.PlayStreak = 1
		}
	} else {
		as.statistics.PlayStreak = 1
	}
	as.statistics.LastPlayDate = &now
	
	// Check dedicated player achievement
	as.achievements[AchievementDedicated].Progress = as.statistics.PlayStreak
	as.checkAchievement(AchievementDedicated)
}

func (as *AchievementSystem) OnGameWin(moves int, gameTime time.Duration, isTimeAttack bool, isPerfect bool) {
	as.statistics.GamesWon++
	as.statistics.TotalMoves += moves
	as.statistics.TotalTime += gameTime
	
	// Update best records
	if as.statistics.BestTime == 0 || gameTime < as.statistics.BestTime {
		as.statistics.BestTime = gameTime
	}
	
	if moves < as.statistics.FewestMoves {
		as.statistics.FewestMoves = moves
	}
	
	// Time Attack specific
	if isTimeAttack {
		as.statistics.TimeAttackWins++
		as.achievements[AchievementTimeAttackWin].Progress = as.statistics.TimeAttackWins
		as.checkAchievement(AchievementTimeAttackWin)
	}
	
	// Perfect game
	if isPerfect {
		as.statistics.PerfectGames++
		as.achievements[AchievementPerfectGame].Progress = as.statistics.PerfectGames
		as.checkAchievement(AchievementPerfectGame)
		as.checkAchievement(AchievementEfficient)
	}
	
	// Check achievements
	as.achievements[AchievementFirstWin].Progress = min(1, as.statistics.GamesWon)
	as.checkAchievement(AchievementFirstWin)
	
	as.achievements[AchievementIslandHopper].Progress = as.statistics.GamesWon
	as.checkAchievement(AchievementIslandHopper)
	
	// Speed achievement (under 30 seconds)
	if gameTime < 30*time.Second {
		as.achievements[AchievementSpeedrun].Progress = 1
		as.checkAchievement(AchievementSpeedrun)
	}
}

func (as *AchievementSystem) OnBridgeBuilt() {
	as.statistics.BridgesBuilt++
	as.achievements[AchievementBridgeBuilder].Progress = as.statistics.BridgesBuilt
	as.checkAchievement(AchievementBridgeBuilder)
}

func (as *AchievementSystem) OnLevelCreated() {
	as.statistics.LevelsCreated++
	as.achievements[AchievementLevelCreator].Progress = as.statistics.LevelsCreated
	as.checkAchievement(AchievementLevelCreator)
}

func (as *AchievementSystem) GetAchievements() []*Achievement {
	result := make([]*Achievement, 0)
	for _, achievement := range as.achievements {
		if !achievement.Hidden || achievement.Unlocked {
			result = append(result, achievement)
		}
	}
	return result
}

func (as *AchievementSystem) GetStatistics() *GameStatistics {
	return as.statistics
}

func (as *AchievementSystem) GetUnlockedCount() int {
	count := 0
	for _, achievement := range as.achievements {
		if achievement.Unlocked {
			count++
		}
	}
	return count
}

func (as *AchievementSystem) GetTotalCount() int {
	return len(as.achievements)
}

// Save/Load functionality
func (as *AchievementSystem) SaveToJSON() (string, error) {
	data := struct {
		Achievements map[AchievementType]*Achievement `json:"achievements"`
		Statistics   *GameStatistics                  `json:"statistics"`
	}{
		Achievements: as.achievements,
		Statistics:   as.statistics,
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	
	return string(jsonData), nil
}

func (as *AchievementSystem) LoadFromJSON(jsonStr string) error {
	var data struct {
		Achievements map[AchievementType]*Achievement `json:"achievements"`
		Statistics   *GameStatistics                  `json:"statistics"`
	}
	
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return err
	}
	
	if data.Achievements != nil {
		as.achievements = data.Achievements
	}
	
	if data.Statistics != nil {
		as.statistics = data.Statistics
	}
	
	return nil
}

func (as *AchievementSystem) GetProgressSummary() string {
	unlocked := as.GetUnlockedCount()
	total := as.GetTotalCount()
	return fmt.Sprintf("Achievements: %d/%d unlocked", unlocked, total)
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}