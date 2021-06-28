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
func (u *Queen) UnitType() int {
	return UnitTypeQueen
}

// UnitMaxHP see function from Unit
func (u *Queen) UnitMaxHP() int {
	return UnitMaxHPQueen
}

// UnitWeight see function from Unit
func (u *Queen) UnitWeight() int {
	return UnitWeightQueen
}

// UnitMoveType see function from Unit
func (u *Queen) UnitMoveType() int {
	return MoveTypeGround
}

// UnitMoveRange see function from Unit
func (u *Queen) UnitMoveRange() int {
	return UnitMoveRangeQueen
}

// UnitMoveRangeMin see function from Unit
func (u *Queen) UnitMoveRangeMin() int { return 0 }

// UnitAttackType see function from Unit
func (u *Queen) UnitAttackType() int {
	return AttackTypeNone
}

// UnitAttackRange see function from Unit
func (u *Queen) UnitAttackRange() int { return 0 }

// UnitAttackRangeMin see function from Unit
func (u *Queen) UnitAttackRangeMin() int { return 0 }

// UnitAttackPower see function from Unit
func (u *Queen) UnitAttackPower() int { return 0 }

// UnitCost see function from Unit
func (u *Queen) UnitCost() int {
	return UnitCostQueen
}

// GetOwner see function from Unit
func (u *Queen) GetOwner() int {
	return u.Owner
}

// GetState see function from Unit
func (u *Queen) GetState() int {
	return u.State
}

// GetStateBit see function from Unit
func (u *Queen) GetStateBit(bit int) bool {
	return (u.State & bit) != 0
}

// ToggleStateBit see function from Unit
func (u *Queen) ToggleStateBit(bit int) {
	u.State ^= bit
}

// GetHP see function from Unit
func (u *Queen) GetHP() int {
	return u.HP
}

// SetHP see function from Unit
func (u *Queen) SetHP(hp int) {
	u.HP = hp
}

// StartTurn see function from Unit
func (u *Queen) StartTurn() {}

// EndTurn see function from Unit
func (u *Queen) EndTurn() {
	// turn off `moved` bit
	u.State &= ^UnitStateBitMoved
}
