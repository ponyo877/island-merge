package core

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
	StateLevelSelect
	StateLevelEditor
)

type GameMode int

const (
	ModeClassic GameMode = iota
	ModeTimeAttack
	ModePuzzle
)