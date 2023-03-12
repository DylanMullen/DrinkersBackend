package models

type RequestID int

const (
	STATE_CHANGE_REQUEST RequestID = 0
	JOIN_REQUEST         RequestID = 1
	LEAVE_REQUEST        RequestID = 2
	NEXT_TURN            RequestID = 3
	KICK_PLAYER          RequestID = 4
	PROMOTE_PLAYER       RequestID = 5
	PROMPT               RequestID = 6
	THEME                RequestID = 7
)

type Request struct {
	ID      RequestID              `json:"id"`
	Sender  string                 `json:"sender"`
	Content map[string]interface{} `json:"content"`
}

type CreateReq struct {
	Owner     HigherLowerPlayer `json:"owner"`
	GameType  GameType          `json:"type"`
	MaxRounds int               `json:"maxRounds"`
	Debug     bool              `json:"debug"`
}

type JoinReq struct {
	JoinCode string            `json:"joinCode"`
	Player   HigherLowerPlayer `json:"player"`
}

type NextTurnReq struct {
	Action string `json:"action"`
}

type ThemeReq struct {
	Type  string                 `json:"type"`
	Table string                 `jons:"color"`
	Red   map[string]interface{} `json:"red"`
	Black map[string]interface{} `json:"black"`
}
