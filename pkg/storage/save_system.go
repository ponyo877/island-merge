package storage

import (
	"fmt"
	"time"
)

const (
	SaveKeyGameState     = "island_merge_game_state"
	SaveKeyAchievements  = "island_merge_achievements"
	SaveKeySettings      = "island_merge_settings"
	SaveKeyCustomLevels  = "island_merge_custom_levels"
	SaveKeyProgress      = "island_merge_progress"
)

// GameSaveData represents the complete saved game state
type GameSaveData struct {
	Version       string                 `json:"version"`
	SavedAt       time.Time              `json:"saved_at"`
	CurrentGame   *CurrentGameState      `json:"current_game,omitempty"`
	Achievements  interface{}            `json:"achievements,omitempty"`
	Settings      *GameSettings          `json:"settings"`
	Progress      *GameProgress          `json:"progress"`
	CustomLevels  []CustomLevel          `json:"custom_levels,omitempty"`
}

// CurrentGameState stores the state of an ongoing game
type CurrentGameState struct {
	Mode        int           `json:"mode"`
	Board       BoardData     `json:"board"`
	Score       ScoreData     `json:"score"`
	StartTime   time.Time     `json:"start_time"`
	TimeLimit   time.Duration `json:"time_limit,omitempty"`
	GameWon     bool          `json:"game_won"`
}

// BoardData represents the game board state
type BoardData struct {
	Width   int     `json:"width"`
	Height  int     `json:"height"`
	Tiles   [][]int `json:"tiles"`
	Islands []int   `json:"islands"`
}

// ScoreData represents the current score
type ScoreData struct {
	Moves    int           `json:"moves"`
	Time     time.Duration `json:"time"`
	BestTime time.Duration `json:"best_time,omitempty"`
}

// GameSettings stores user preferences
type GameSettings struct {
	SoundEnabled     bool    `json:"sound_enabled"`
	MusicEnabled     bool    `json:"music_enabled"`
	AnimationSpeed   float64 `json:"animation_speed"`
	ShowTutorial     bool    `json:"show_tutorial"`
	AutoSave         bool    `json:"auto_save"`
	PreferredMode    int     `json:"preferred_mode"`
}

// GameProgress tracks overall game progress
type GameProgress struct {
	CompletedLevels   []string  `json:"completed_levels"`
	HighScores        []Score   `json:"high_scores"`
	TotalPlayTime     time.Duration `json:"total_play_time"`
	LastPlayed        time.Time `json:"last_played"`
	UnlockedModes     []int     `json:"unlocked_modes"`
}

// Score represents a high score entry
type Score struct {
	Level     string        `json:"level"`
	Mode      int           `json:"mode"`
	Moves     int           `json:"moves"`
	Time      time.Duration `json:"time"`
	Date      time.Time     `json:"date"`
	PlayerID  string        `json:"player_id,omitempty"`
}

// CustomLevel represents a user-created level
type CustomLevel struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	Author      string    `json:"author,omitempty"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Tiles       [][]int   `json:"tiles"`
	Difficulty  string    `json:"difficulty,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
}

// SaveSystem manages all save/load operations
type SaveSystem struct {
	storage *LocalStorage
}

func NewSaveSystem() *SaveSystem {
	return &SaveSystem{
		storage: NewLocalStorage(),
	}
}

// SaveGameState saves the current game state
func (ss *SaveSystem) SaveGameState(gameState *CurrentGameState) error {
	return ss.storage.Set(SaveKeyGameState, gameState)
}

// LoadGameState loads the saved game state
func (ss *SaveSystem) LoadGameState() (*CurrentGameState, error) {
	var gameState CurrentGameState
	err := ss.storage.Get(SaveKeyGameState, &gameState)
	if err != nil {
		return nil, err
	}
	return &gameState, nil
}

// HasSavedGame checks if there's a saved game
func (ss *SaveSystem) HasSavedGame() bool {
	return ss.storage.Exists(SaveKeyGameState)
}

// DeleteSavedGame removes the saved game state
func (ss *SaveSystem) DeleteSavedGame() {
	ss.storage.Remove(SaveKeyGameState)
}

// SaveAchievements saves achievement data
func (ss *SaveSystem) SaveAchievements(achievements interface{}) error {
	return ss.storage.Set(SaveKeyAchievements, achievements)
}

// LoadAchievements loads achievement data
func (ss *SaveSystem) LoadAchievements(target interface{}) error {
	return ss.storage.Get(SaveKeyAchievements, target)
}

// SaveSettings saves game settings
func (ss *SaveSystem) SaveSettings(settings *GameSettings) error {
	return ss.storage.Set(SaveKeySettings, settings)
}

// LoadSettings loads game settings
func (ss *SaveSystem) LoadSettings() (*GameSettings, error) {
	var settings GameSettings
	err := ss.storage.Get(SaveKeySettings, &settings)
	if err != nil {
		// Return default settings if none found
		return ss.GetDefaultSettings(), nil
	}
	return &settings, nil
}

// GetDefaultSettings returns default game settings
func (ss *SaveSystem) GetDefaultSettings() *GameSettings {
	return &GameSettings{
		SoundEnabled:   true,
		MusicEnabled:   true,
		AnimationSpeed: 1.0,
		ShowTutorial:   true,
		AutoSave:       true,
		PreferredMode:  0, // Classic mode
	}
}

// SaveProgress saves game progress
func (ss *SaveSystem) SaveProgress(progress *GameProgress) error {
	return ss.storage.Set(SaveKeyProgress, progress)
}

// LoadProgress loads game progress
func (ss *SaveSystem) LoadProgress() (*GameProgress, error) {
	var progress GameProgress
	err := ss.storage.Get(SaveKeyProgress, &progress)
	if err != nil {
		// Return default progress if none found
		return &GameProgress{
			CompletedLevels: []string{},
			HighScores:      []Score{},
			UnlockedModes:   []int{0}, // Start with Classic mode unlocked
			LastPlayed:      time.Now(),
		}, nil
	}
	return &progress, nil
}

// SaveCustomLevel saves a custom level
func (ss *SaveSystem) SaveCustomLevel(level *CustomLevel) error {
	levels, err := ss.LoadCustomLevels()
	if err != nil {
		levels = []CustomLevel{}
	}
	
	// Check if level already exists and update, or add new
	found := false
	for i, existingLevel := range levels {
		if existingLevel.ID == level.ID {
			levels[i] = *level
			found = true
			break
		}
	}
	
	if !found {
		levels = append(levels, *level)
	}
	
	return ss.storage.Set(SaveKeyCustomLevels, levels)
}

// LoadCustomLevels loads all custom levels
func (ss *SaveSystem) LoadCustomLevels() ([]CustomLevel, error) {
	var levels []CustomLevel
	err := ss.storage.Get(SaveKeyCustomLevels, &levels)
	if err != nil {
		return []CustomLevel{}, nil
	}
	return levels, nil
}

// DeleteCustomLevel deletes a custom level
func (ss *SaveSystem) DeleteCustomLevel(levelID string) error {
	levels, err := ss.LoadCustomLevels()
	if err != nil {
		return err
	}
	
	// Filter out the level to delete
	var newLevels []CustomLevel
	for _, level := range levels {
		if level.ID != levelID {
			newLevels = append(newLevels, level)
		}
	}
	
	return ss.storage.Set(SaveKeyCustomLevels, newLevels)
}

// ExportSaveData exports all save data as JSON
func (ss *SaveSystem) ExportSaveData() (*GameSaveData, error) {
	saveData := &GameSaveData{
		Version: "1.0",
		SavedAt: time.Now(),
	}
	
	// Load current game state
	if gameState, err := ss.LoadGameState(); err == nil {
		saveData.CurrentGame = gameState
	}
	
	// Load achievements
	var achievements interface{}
	if err := ss.LoadAchievements(&achievements); err == nil {
		saveData.Achievements = achievements
	}
	
	// Load settings
	if settings, err := ss.LoadSettings(); err == nil {
		saveData.Settings = settings
	}
	
	// Load progress
	if progress, err := ss.LoadProgress(); err == nil {
		saveData.Progress = progress
	}
	
	// Load custom levels
	if levels, err := ss.LoadCustomLevels(); err == nil {
		saveData.CustomLevels = levels
	}
	
	return saveData, nil
}

// ImportSaveData imports save data from JSON
func (ss *SaveSystem) ImportSaveData(saveData *GameSaveData) error {
	if saveData.CurrentGame != nil {
		if err := ss.SaveGameState(saveData.CurrentGame); err != nil {
			return fmt.Errorf("failed to import game state: %w", err)
		}
	}
	
	if saveData.Achievements != nil {
		if err := ss.SaveAchievements(saveData.Achievements); err != nil {
			return fmt.Errorf("failed to import achievements: %w", err)
		}
	}
	
	if saveData.Settings != nil {
		if err := ss.SaveSettings(saveData.Settings); err != nil {
			return fmt.Errorf("failed to import settings: %w", err)
		}
	}
	
	if saveData.Progress != nil {
		if err := ss.SaveProgress(saveData.Progress); err != nil {
			return fmt.Errorf("failed to import progress: %w", err)
		}
	}
	
	if saveData.CustomLevels != nil {
		for _, level := range saveData.CustomLevels {
			if err := ss.SaveCustomLevel(&level); err != nil {
				return fmt.Errorf("failed to import custom level %s: %w", level.ID, err)
			}
		}
	}
	
	return nil
}

// ClearAllData removes all saved data
func (ss *SaveSystem) ClearAllData() {
	ss.storage.Remove(SaveKeyGameState)
	ss.storage.Remove(SaveKeyAchievements)
	ss.storage.Remove(SaveKeySettings)
	ss.storage.Remove(SaveKeyCustomLevels)
	ss.storage.Remove(SaveKeyProgress)
}

// GetStorageUsage returns information about storage usage
func (ss *SaveSystem) GetStorageUsage() map[string]bool {
	return map[string]bool{
		"game_state":    ss.storage.Exists(SaveKeyGameState),
		"achievements":  ss.storage.Exists(SaveKeyAchievements),
		"settings":      ss.storage.Exists(SaveKeySettings),
		"custom_levels": ss.storage.Exists(SaveKeyCustomLevels),
		"progress":      ss.storage.Exists(SaveKeyProgress),
	}
}