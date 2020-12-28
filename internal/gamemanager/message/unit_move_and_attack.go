package message

// UnitMoveAndAttackMessageData is message data for CmdUnitMoveAndAttack
type UnitMoveAndAttackMessageData struct {
	Y1 int `json:"y_1"`
	X1 int `json:"x_1"`
	Y2 int `json:"y_2"`
	X2 int `json:"x_2"`
	YT int `json:"y_t"`
	XT int `json:"x_t"`
}

// UnitMoveAndAttackMessageDataExt is message data for CmdUnitMoveAndAttack with additional fields
// for use when BE is sending message to FE
type UnitMoveAndAttackMessageDataExt struct {
	Y1    int `json:"y_1"`
	X1    int `json:"x_1"`
	Y2    int `json:"y_2"`
	X2    int `json:"x_2"`
	YT    int `json:"y_t"`
	XT    int `json:"x_t"`
	HPAtk int `json:"hp_atk"`
	HPDef int `json:"hp_def"`
}
