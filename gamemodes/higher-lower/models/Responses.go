package models

import (
	"encoding/json"
)

type Response struct {
	ID      int `json:"id"`
	Content any `json:"body"`
}

type Turn struct {
	Player        string `json:"current"`
	DidWin        bool   `json:"winner"`
	Next          string `json:"next"`
	CurrentNumber int    `json:"number"`
	Suite         int    `jons:"suite"`
	Streak        int    `json:"streak"`
}
type Join struct {
	Current    string            `json:"current"`
	NextPlayer string            `json:"next"`
	Player     HigherLowerPlayer `json:"player"`
}

type Kick struct {
	Player string `json:"uuid"`
}

type GameState struct {
	Started bool `json:"started"`
}

type Promotion struct {
	UUID string `json:"uuid"`
}

type Prompt struct {
	Owner       string `json:"owner"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        int    `json:"time,omitempty"`
}

func Marshal(id int, content any) []byte {
	response, _ := json.Marshal(Response{
		ID:      id,
		Content: content,
	})
	return response
}

func GameStateChange(started bool) []byte {
	return Marshal(int(STATE_CHANGE_REQUEST), GameState{
		Started: started,
	})
}

func NewPlayerAppeared(player HigherLowerPlayer, current string, next string) []byte {

	return Marshal(int(JOIN_REQUEST), Join{
		Player:     player,
		Current:    current,
		NextPlayer: next,
	})
}

func PlayerBooted(player string) []byte {
	return Marshal(int(LEAVE_REQUEST), Kick{
		Player: player,
	})
}

func NextTurn(player string, next string, currentNumber int, suite int, winner bool, streak int) []byte {

	return Marshal(int(NEXT_TURN), Turn{
		Player:        player,
		DidWin:        winner,
		Next:          next,
		CurrentNumber: currentNumber,
		Suite:         suite,
		Streak:        streak,
	})
}

func PlayerPromoted(uuid string) []byte {
	return Marshal(int(PROMOTE_PLAYER), Promotion{
		UUID: uuid,
	})
}

func SendPrompt(prompt Prompt) []byte {
	return Marshal(int(PROMPT), prompt)
}
