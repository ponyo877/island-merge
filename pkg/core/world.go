package core

import (
	"time"
	
	"github.com/ponyo877/island-merge/pkg/island"
)

type World struct {
	Board   *island.Board
	Score   Score
	GameWon bool
}

type Score struct {
	Moves       int
	Time        time.Duration
	IslandsLeft int
}