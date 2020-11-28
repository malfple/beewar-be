package units

// You is a unit object. It is most often the win or lose condition in a game.
type You struct {
	P int8
}

// GetUnitType see function from Unit
func (you *You) GetUnitType() int8 {
	return UnitTypeYou
}
