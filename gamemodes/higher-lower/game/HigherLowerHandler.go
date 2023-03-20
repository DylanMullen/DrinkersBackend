package game

import (
	"encoding/json"
	"strings"

	"drinkers.beer/waterfall/gamemodes/higher-lower/models"
)

func (game *HigherLower) handleMessage(senderID string, message []byte) {
	var request models.Request
	err := json.Unmarshal(message, &request)

	if err != nil {
		return
	}

	switch request.ID {
	case models.STATE_CHANGE_REQUEST:
		handleState(game, request)
	case models.JOIN_REQUEST:
		handleJoin(game, request)
	case models.LEAVE_REQUEST:
		handleLeave(game, request, false)
	case models.NEXT_TURN:
		handleNextTurn(game, request)
	case models.KICK_PLAYER:
		handleLeave(game, request, true)
	case models.PROMOTE_PLAYER:
		handlePromotion(game, request)
	case models.THEME:
		handleTheme(game, request)
	default:
		break
	}
}

func handleJoin(game *HigherLower, request models.Request) {
	var req models.JoinReq
	jsonData, _ := json.Marshal(request.Content)
	json.Unmarshal(jsonData, &req)
	game.AddPlayer(req.Player)
}

func handleLeave(game *HigherLower, request models.Request, kick bool) {
	if kick && !isOwner(*game, request.Sender) {
		return
	}

	val, ok := request.Content["uuid"].(string)
	if ok {
		game.KickPlayer(val)
	}
}

func handleState(game *HigherLower, request models.Request) {
	if !isOwner(*game, request.Sender) {
		return
	}

	val, ok := request.Content["state"].(bool)

	if !ok {
		return
	}

	if val {
		game.start()
	} else {
		game.stop()
	}
}

func handleNextTurn(game *HigherLower, request models.Request) {
	if !game.HasStarted() {
		return
	}

	var req models.NextTurnReq
	jsonData, _ := json.Marshal(request.Content)
	json.Unmarshal(jsonData, &req)

	game.nextTurn(strings.ToLower(req.Action))
}

func handlePromotion(game *HigherLower, request models.Request) {
	if !isOwner(*game, request.Sender) {
		return
	}

	val, ok := request.Content["uuid"].(string)
	if ok {
		game.promote(val)
	}
}

func handleTheme(game *HigherLower, request models.Request) {
	var req models.ThemeReq
	jsonData, _ := json.Marshal(request.Content)
	err := json.Unmarshal(jsonData, &req)

	if err != nil {
		return
	}

	if req.Type == "table" {
		game.Theme.Table = req.Table
	} else {
		game.Theme.Red = req.Red
		game.Theme.Black = req.Black
	}

	go game.socket.BroadcastAll(models.Marshal(int(models.THEME), game.Theme))
}

func isOwner(game HigherLower, uuid string) bool {
	return game.Settings.OwnerID == uuid
}
