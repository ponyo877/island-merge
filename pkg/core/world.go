package core

import (
	"time"
	
	"github.com/ponyo877/island-merge/pkg/island"
)

type World struct {
	State     GameState
	Mode      GameMode
	Board     *island.Board
	Score     Score
	GameWon   bool
	StartTime time.Time
	TimeLimit time.Duration // For Time Attack mode
}

type Score struct {
	Moves       int
	Time        time.Duration
	IslandsLeft int
	BestTime    time.Duration
	BestMoves   int
}

// Methods for interface compliance
func (w *World) GetMode() int {
	return int(w.Mode)
}

func (w *World) GetScore() interface {
	GetMoves() int
	GetTime() time.Duration
} {
	return w.Score
}

func (w *World) GetTimeLimit() time.Duration {
	return w.TimeLimit
}

func (w *World) GetState() int {
	return int(w.State)
}

func (s Score) GetMoves() int {
	return s.Moves
}

func (s Score) GetTime() time.Duration {
	return s.Time
}