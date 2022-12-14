package message

import "gitlab.com/beewar/beewar-be/internal/access/model"

// GameDataMessageData is message data for CmdGameData
type GameDataMessageData struct {
	Game    *model.Game `json:"game"`
	Players []*Player   `json:"players"`
}

// Player is similar to model.GameUser but with some fields removed and user data added
type Player struct {
	UserID      uint64      `json:"user_id"`
	PlayerOrder uint8       `json:"player_order"`
	FinalRank   uint8       `json:"final_rank"`
	FinalTurns  int32       `json:"final_turns"`
	User        *model.User `json:"user"`
}
