package objects

// Queen is a unit object. It is most often the win or lose condition in a game.
// State description: 1 bit for moved state
type Queen struct {
	Owner int
	HP    int
	State int
}

// NewQueen returns a new Queen object
func NewQueen(owner, hp, state int) *Queen {
	return &Queen{
		Owner: owner,
		HP:    hp,
		State: state,
	}
}

// GetUnitType see function from Unit
func (queen *Queen) GetUnitType() int {
	return UnitTypeQueen
}

// GetUnitOwner see function from Unit
func (queen *Queen) GetUnitOwner() int {
	return queen.Owner
}

// GetUnitState see function from Unit
func (queen *Queen) GetUnitState() int {
	return queen.State
}

// GetUnitStateBit see function from Unit
func (queen *Queen) GetUnitStateBit(bit int) bool {
	return (queen.State & bit) != 0
}

// ToggleUnitStateBit see function from Unit
func (queen *Queen) ToggleUnitStateBit(bit int) {
	queen.State ^= bit
}

// GetUnitHP see function from Unit
func (queen *Queen) GetUnitHP() int {
	return queen.HP
}

// SetUnitHP see function from Unit
func (queen *Queen) SetUnitHP(hp int) {
	queen.HP = hp
}

// GetWeight see function from Unit
func (queen *Queen) GetWeight() int {
	return UnitWeightQueen
}

// GetMoveType see function from Unit
func (queen *Queen) GetMoveType() int {
	return MoveTypeGround
}

// GetMoveRange see function from Unit
func (queen *Queen) GetMoveRange() int {
	return UnitMoveRangeQueen
}

// GetAttackType see function from Unit
func (queen *Queen) GetAttackType() int {
	return AttackTypeNone
}

// GetAttackRange see function from Unit
func (queen *Queen) GetAttackRange() int {
	return UnitAttackRangeQueen
}

// GetAttackPower see function from Unit
func (queen *Queen) GetAttackPower() int {
	return UnitAttackPowerQueen
}

// StartTurn see function from Unit
func (queen *Queen) StartTurn() {}

// EndTurn see function from Unit
func (queen *Queen) EndTurn() {
	// turn off `moved` bit
	queen.State &= ^UnitStateBitMoved
}
