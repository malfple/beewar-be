package objects

// You is a unit object. It is most often the win or lose condition in a game.
type You struct {
	Owner int
	HP    int
	State int
}

// NewYou returns a new You object
func NewYou(owner, hp, state int) *You {
	return &You{
		Owner: owner,
		HP:    hp,
		State: state,
	}
}

// GetUnitType see function from Unit
func (you *You) GetUnitType() int {
	return UnitTypeYou
}

// GetUnitState see function from Unit
func (you *You) GetUnitState() int {
	return you.State
}
