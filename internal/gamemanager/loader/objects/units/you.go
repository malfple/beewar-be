package units

// You is a unit object. It is most often the win or lose condition in a game.
type You struct {
	Owner uint8
	HP    uint8
	State uint8
}

// NewYou returns a new You object
func NewYou(owner, hp, state uint8) *You {
	return &You{
		Owner: owner,
		HP:    hp,
		State: state,
	}
}

// GetUnitType see function from Unit
func (you *You) GetUnitType() uint8 {
	return UnitTypeYou
}

// GetUnitState see function from Unit
func (you *You) GetUnitState() uint8 {
	return you.State
}
