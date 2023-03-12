package models

import (
	"drinkers.beer/waterfall/models"
	"drinkers.beer/waterfall/shared/packs"
)

type GameType string

const (
	CARD GameType = "CARD"
	DICE GameType = "DICE"
)

type GameSettings struct {
	UUID       string   `json:"gameID"`
	JoinCode   string   `json:"joinCode"`
	OwnerID    string   `json:"ownerID"`
	GameType   GameType `json:"type"`
	HasStarted bool     `json:"started"`
}

type Gameplay struct {
	CurrentNumber int     `json:"currentNumber"`
	Suite         int     `json:"suite"`
	MaxRounds     int     `json:"maxRounds"`
	CurrentStreak int     `json:"streak"`
	Players       Players `json:"players"`

	CurrentRounds int        `json:"-"`
	Pack          packs.Pack `json:"-"`
}

type Theme struct {
	Table string                 `json:"table,omitempty"`
	Red   map[string]interface{} `json:"red,omitempty"`
	Black map[string]interface{} `json:"black,omitempty"`
}

type Players struct {
	CurrentPlayer string              `json:"current"`
	NextPlayer    string              `json:"next"`
	Players       []HigherLowerPlayer `json:"players"`
}

type HigherLowerPlayer struct {
	models.User
}

// type CardStyle struct {
// 	Card CardTheme `json:"card"`
// 	Pip  PipTheme  `json:"pip"`
// }

// type CardTheme struct {
// 	BG string `json:"cardBackground"`
// 	Border string `json:""`
// }

// type PipTheme struct {
// }
