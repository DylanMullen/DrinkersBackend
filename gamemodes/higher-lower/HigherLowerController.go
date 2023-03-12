package higherlower

import (
	"net/http"

	"drinkers.beer/waterfall/gamemodes/higher-lower/game"
	"drinkers.beer/waterfall/gamemodes/higher-lower/models"
	"drinkers.beer/waterfall/gamemodes/pirate/errors"
	"drinkers.beer/waterfall/shared/packs"
	"drinkers.beer/waterfall/websockets"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	games  map[string]*game.HigherLower = make(map[string]*game.HigherLower, 0)
	loader packs.PackLoader             = packs.PackLoader{Location: "./gamemodes/higher-lower/packs"}
)

func HandleRequests(router *gin.RouterGroup) {
	router.POST("/create", handleCreate)
	router.POST("/join", handleJoin)
	router.GET("/ws/:gameID", handleWebsocketRequest)

	LoadPacks()
}

func LoadPacks() {
	loader.LoadAllPacks()
}

func handleCreate(c *gin.Context) {
	var req models.CreateReq
	err := c.BindJSON(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			ID:      0,
			Content: "Could not create game",
		})
		return
	}

	game := createGame(req)

	c.JSON(http.StatusOK, models.Response{
		ID:      0,
		Content: game,
	})
}

func handleJoin(c *gin.Context) {
	var req models.JoinReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			ID:      1,
			Content: "Could not join game",
		})
		return
	}

	game, err := getHigherLowerGameByCode(req.JoinCode)

	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			ID:      1,
			Content: "Could not find game",
		})
		return
	}

	game.AddPlayer(req.Player)

	c.JSON(http.StatusOK, models.Response{
		ID:      1,
		Content: game,
	})
}

func createGame(req models.CreateReq) game.HigherLower {

	id := uuid.NewString()
	if req.Debug {
		id = "debug"
	}

	prompts := make([]packs.Prompt, 0)

	for _, v := range loader.Packs {
		prompts = append(prompts, v.Prompts...)
	}

	pack := packs.Pack{
		Settings: loader.Packs["3956ceb7-fcdf-4691-bc3b-80db4085c2be"].Settings,
		Prompts:  prompts,
	}

	game := game.HigherLower{
		Settings: models.GameSettings{
			UUID:     id,
			JoinCode: id[:6],
			OwnerID:  req.Owner.UUID,
			GameType: req.GameType,
		},
		GamePlay: models.Gameplay{
			MaxRounds: req.MaxRounds,
			Pack:      pack,
		},
	}

	game.Init(req.Owner)
	go game.HandleSocket()

	games[game.Settings.UUID] = &game

	return game

}

func getHigherLowerGame(uuid string) (game *game.HigherLower, err error) {
	for _, v := range games {
		if v.Settings.UUID == uuid {
			game = v
			return
		}
	}
	err = errors.GameNotFound{
		Message: uuid,
	}
	return
}

func getHigherLowerGameByCode(code string) (game *game.HigherLower, err error) {
	for _, v := range games {
		if v.Settings.JoinCode == code {
			game = v
			return
		}
	}
	err = errors.GameNotFound{
		Message: code,
	}
	return
}

func handleWebsocketRequest(c *gin.Context) {
	gameId := c.Param("gameID")

	game, err := getHigherLowerGame(gameId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": errors.GameNotFound{Message: "UUID: " + gameId}})
		return
	}
	websockets.HandleWs(game.GetSocket(), c)
}
