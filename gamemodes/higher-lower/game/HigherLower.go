package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"drinkers.beer/waterfall/gamemodes/higher-lower/models"
	"drinkers.beer/waterfall/gamemodes/pirate/errors"
	"drinkers.beer/waterfall/websockets"
)

type HigherLower struct {
	socket   *websockets.SocketHub
	Settings models.GameSettings `json:"settings"`
	GamePlay models.Gameplay     `json:"gameplay"`
	Theme    models.Theme        `json:"theme"`
}

func (game *HigherLower) Init(player models.HigherLowerPlayer) {
	game.GamePlay.CurrentNumber = game.newNumber()
	game.GamePlay.Players.Players = make([]models.HigherLowerPlayer, 0)
	game.socket = websockets.GetNewHub()

	game.AddPlayer(player)
	game.GamePlay.Players.CurrentPlayer = player.UUID
	game.GamePlay.Players.NextPlayer = player.UUID
}

func (game *HigherLower) start() {
	game.Settings.HasStarted = true
	game.socket.BroadcastAll(models.GameStateChange(true))
}

func (game *HigherLower) stop() {
	game.Settings.HasStarted = false
	game.socket.BroadcastAll(models.GameStateChange(false))
}

func (game *HigherLower) promote(uuid string) {
	if !game.isPlayer(uuid) {
		return
	}

	game.Settings.OwnerID = uuid

	game.socket.BroadcastAll(models.PlayerPromoted(uuid))
}

func (game *HigherLower) nextTurn(action string) {

	nextNum := game.newNumber()
	shouldNextPlayer := true

	switch action {
	case "higher":
		if nextNum >= game.GamePlay.CurrentNumber {
			shouldNextPlayer = false
			break
		}
	case "lower":
		if nextNum <= game.GamePlay.CurrentNumber {
			shouldNextPlayer = false
			break
		}
	default:
		break
	}
	game.GamePlay.CurrentNumber = nextNum
	game.GamePlay.Suite = game.newSuite()

	if !shouldNextPlayer {
		game.GamePlay.CurrentStreak = game.GamePlay.CurrentStreak + 1
		go game.sendPrompt()
	} else {
		game.nextPlayer()
		game.GamePlay.Players.NextPlayer = game.GamePlay.Players.Players[game.getNextPlayerIndex()].UUID
		game.GamePlay.CurrentStreak = 0
	}

	if game.GamePlay.MaxRounds != 0 && game.GamePlay.MaxRounds == game.GamePlay.CurrentRounds {
		game.socket.BroadcastMessage("", []byte("Finished"))
		return
	}

	game.socket.BroadcastAll(models.NextTurn(
		game.GamePlay.Players.CurrentPlayer,
		game.GamePlay.Players.NextPlayer,
		game.GamePlay.CurrentNumber,
		game.GamePlay.Suite,
		!shouldNextPlayer,
		game.GamePlay.CurrentStreak,
	))

}

func (game *HigherLower) newNumber() int {
	min, max := game.getMaxNumber()
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func (game *HigherLower) newSuite() int {
	return rand.Intn(4)
}

func (game HigherLower) sendPrompt() {
	details := game.GamePlay.Pack.Prompts[rand.Intn(len(game.GamePlay.Pack.Prompts))]

	var prompt models.Prompt
	jsonData, _ := json.Marshal(details.Details)
	err := json.Unmarshal(jsonData, &prompt)

	if err != nil {
		fmt.Println("Failed to send Prompt: " + details.UUID)
		return
	}
	prompt.Owner = game.GamePlay.Players.CurrentPlayer
	game.socket.BroadcastAll(models.SendPrompt(prompt))
}

func (game *HigherLower) nextPlayer() {
	index := game.getNextPlayerIndex()
	if game.GamePlay.MaxRounds != 0 && (index == 0) {
		game.GamePlay.CurrentRounds++
	}
	game.GamePlay.Players.CurrentPlayer = game.GamePlay.Players.Players[index].UUID
}

func (game *HigherLower) AddPlayer(player models.HigherLowerPlayer) {
	if game.isPlayer(player.UUID) {
		return
	}

	game.GamePlay.Players.Players = append(game.GamePlay.Players.Players, player)

	if len(game.GamePlay.Players.Players) > 0 {
		game.socket.BroadcastAll(models.NewPlayerAppeared(
			player,
			game.GamePlay.Players.CurrentPlayer,
			game.GamePlay.Players.CurrentPlayer,
		))
	}
}

func (game *HigherLower) KickPlayer(uuid string) {
	if !game.isPlayer(uuid) {
		return
	}

	players := make([]models.HigherLowerPlayer, 0)

	for i := 0; i < len(game.GamePlay.Players.Players); i++ {
		if game.GamePlay.Players.Players[i].UUID == uuid {
			continue
		}
		players = append(players, game.GamePlay.Players.Players[i])
	}

	game.GamePlay.Players.Players = players
	game.socket.BroadcastAll(models.PlayerBooted(uuid))
}

func (game HigherLower) isPlayer(uuid string) bool {
	_, err := game.getPlayer(uuid)

	return err == nil
}

func (game HigherLower) getPlayer(uuid string) (player models.HigherLowerPlayer, err error) {
	for i := 0; i < len(game.GamePlay.Players.Players); i++ {
		temp := game.GamePlay.Players.Players[i]
		if temp.UUID == uuid {
			player = temp
			return
		}
	}
	err = errors.PlayerNotFound{
		UUID: uuid,
	}
	return
}

func (game HigherLower) getPlayerIndex(uuid string) (pos int, err error) {
	for i := 0; i < len(game.GamePlay.Players.Players); i++ {
		temp := game.GamePlay.Players.Players[i]
		if temp.UUID == uuid {
			pos = i
			return
		}
	}
	err = errors.PlayerNotFound{
		UUID: uuid,
	}

	return
}

func (game HigherLower) getNextPlayerIndex() int {
	currentIndex, err := game.getPlayerIndex(game.GamePlay.Players.CurrentPlayer)
	nextIndex := currentIndex + 1
	if err != nil {
		return 0
	}

	if (nextIndex > len(game.GamePlay.Players.Players)-1) || (len(game.GamePlay.Players.Players) == 1) {
		return 0
	} else {
		return currentIndex + 1
	}
}

func (game HigherLower) getMaxNumber() (min int, max int) {
	min, max = 0, 0
	switch game.Settings.GameType {
	case models.CARD:
		max = 13
	case models.DICE:
		max = 6
	}

	return
}

func (game *HigherLower) HandleSocket() {
	h := game.socket
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case userMessage := <-h.Broadcast:
			var data map[string][]byte
			json.Unmarshal(userMessage, &data)

			game.handleMessage(string(data["id"]), data["message"])

		}
	}
}

func (game HigherLower) GetSocket() *websockets.SocketHub {
	return game.socket
}

func (game HigherLower) HasStarted() bool {
	return game.Settings.HasStarted
}
