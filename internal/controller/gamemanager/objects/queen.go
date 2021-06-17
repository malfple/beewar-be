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

// UnitType see function from Unit
func (queen *Queen) UnitType() int {
	return UnitTypeQueen
}

// UnitMaxHP see function from Unit
func (queen *Queen) UnitMaxHP() int {
	return UnitMaxHPQueen
}

// UnitWeight see function from Unit
func (queen *Queen) UnitWeight() int {
	return UnitWeightQueen
}

// UnitMoveType see function from Unit
func (queen *Queen) UnitMoveType() int {
	return MoveTypeGround
}

// UnitMoveRange see function from Unit
func (queen *Queen) UnitMoveRange() int {
	return UnitMoveRangeQueen
}

// UnitAttackType see function from Unit
func (queen *Queen) UnitAttackType() int {
	return AttackTypeNone
}

// UnitAttackRange see function from Unit
func (queen *Queen) UnitAttackRange() int {
	return UnitAttackRangeQueen
}

// UnitAttackPower see function from Unit
func (queen *Queen) UnitAttackPower() int {
	return UnitAttackPowerQueen
}

// UnitCost see function from Unit
func (queen *Queen) UnitCost() int {
	return UnitCostQueen
}

// GetOwner see function from Unit
func (queen *Queen) GetOwner() int {
	return queen.Owner
}

// GetState see function from Unit
func (queen *Queen) GetState() int {
	return queen.State
}

// GetStateBit see function from Unit
func (queen *Queen) GetStateBit(bit int) bool {
	return (queen.State & bit) != 0
}

// ToggleStateBit see function from Unit
func (queen *Queen) ToggleStateBit(bit int) {
	queen.State ^= bit
}

// GetHP see function from Unit
func (queen *Queen) GetHP() int {
	return queen.HP
}

// SetHP see function from Unit
func (queen *Queen) SetHP(hp int) {
	queen.HP = hp
}

// StartTurn see function from Unit
func (queen *Queen) StartTurn() {}

// EndTurn see function from Unit
func (queen *Queen) EndTurn() {
	// turn off `moved` bit
	queen.State &= ^UnitStateBitMoved
}
