package model

// GameUser is a db model
type GameUser struct {
	ID          uint64 `json:"id"`
	GameID      uint64 `json:"game_id"`
	UserID      uint64 `json:"user_id"`
	PlayerOrder uint8  `json:"player_order"`
	FinalRank   uint8  `json:"final_rank"`
	FinalTurns  int32  `json:"final_turns"`
}
