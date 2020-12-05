package message

// UnitMoveMessageData is message data for CmdUnitMove
type UnitMoveMessageData struct {
	Y1 int `json:"y_1"`
	X1 int `json:"x_1"`
	Y2 int `json:"y_2"`
	X2 int `json:"x_2"`
}
