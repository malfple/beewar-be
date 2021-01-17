package message

// UnitAttackMessageData is message data for CmdUnitAttack
type UnitAttackMessageData struct {
	Y1 int `json:"y_1"`
	X1 int `json:"x_1"`
	YT int `json:"y_t"`
	XT int `json:"x_t"`
}

// UnitAttackMessageDataExt is message data for CmdUnitAttack with additional fields
// for use when BE is sending message to FE
type UnitAttackMessageDataExt struct {
	Y1    int `json:"y_1"`
	X1    int `json:"x_1"`
	YT    int `json:"y_t"`
	XT    int `json:"x_t"`
	HPAtk int `json:"hp_atk"`
	HPDef int `json:"hp_def"`
}
