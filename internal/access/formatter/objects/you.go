package objects

// You is a unit object. It is most often the win or lose condition in a game.
// State description: 1 bit for moved state
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

// GetWeight see function from Unit
func (you *You) GetWeight() int {
	return UnitWeightYou
}

// GetUnitOwner see function from Unit
func (you *You) GetUnitOwner() int {
	return you.Owner
}

// GetUnitState see function from Unit
func (you *You) GetUnitState() int {
	return you.State
}

// ToggleUnitStateBit see function from Unit
func (you *You) ToggleUnitStateBit(bit int) {
	you.State ^= bit
}

// GetUnitHP see function from Unit
func (you *You) GetUnitHP() int {
	return you.HP
}

// SetUnitHP see function from Unit
func (you *You) SetUnitHP(hp int) {
	you.HP = hp
}

// GetUnitStateBit see function from Unit
func (you *You) GetUnitStateBit(bit int) bool {
	return (you.State & bit) != 0
}

// StartTurn see function from Unit
func (you *You) StartTurn() {}

// EndTurn see function from Unit
func (you *You) EndTurn() {
	// turn off `moved` bit
	you.State &= ^UnitStateBitMoved
}
