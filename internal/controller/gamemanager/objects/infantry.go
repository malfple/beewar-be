package objects

// Infantry is a unit object.
// State description: 1 bit for moved state
type Infantry struct {
	Owner int
	HP    int
	State int
}

// NewInfantry returns a new Infantry object
func NewInfantry(owner, hp, state int) *Infantry {
	return &Infantry{
		Owner: owner,
		HP:    hp,
		State: state,
	}
}

// UnitType see function from Unit
func (u *Infantry) UnitType() int {
	return UnitTypeInfantry
}

// UnitMaxHP see function from Unit
func (u *Infantry) UnitMaxHP() int {
	return UnitMaxHPInfantry
}

// UnitWeight see function from Unit
func (u *Infantry) UnitWeight() int {
	return UnitWeightInfantry
}

// UnitMoveType see function from Unit
func (u *Infantry) UnitMoveType() int {
	return MoveTypeGround
}

// UnitMoveRange see function from Unit
func (u *Infantry) UnitMoveRange() int {
	return UnitMoveRangeInfantry
}

// UnitMoveRangeMin see function from Unit
func (u *Infantry) UnitMoveRangeMin() int { return 0 }

// UnitAttackType see function from Unit
func (u *Infantry) UnitAttackType() int {
	return AttackTypeGround
}

// UnitAttackRange see function from Unit
func (u *Infantry) UnitAttackRange() int {
	return UnitAttackRangeInfantry
}

// UnitAttackRangeMin see function from Unit
func (u *Infantry) UnitAttackRangeMin() int { return 0 }

// UnitAttackPower see function from Unit
func (u *Infantry) UnitAttackPower() int {
	return UnitAttackPowerInfantry
}

// UnitCost see function from Unit
func (u *Infantry) UnitCost() int {
	return UnitCostInfantry
}

// GetOwner see function from Unit
func (u *Infantry) GetOwner() int {
	return u.Owner
}

// GetState see function from Unit
func (u *Infantry) GetState() int {
	return u.State
}

// GetStateBit see function from Unit
func (u *Infantry) GetStateBit(bit int) bool {
	return (u.State & bit) != 0
}

// ToggleStateBit see function from Unit
func (u *Infantry) ToggleStateBit(bit int) {
	u.State ^= bit
}

// GetHP see function from Unit
func (u *Infantry) GetHP() int {
	return u.HP
}

// SetHP see function from Unit
func (u *Infantry) SetHP(hp int) {
	u.HP = hp
}

// StartTurn see function from Unit
func (u *Infantry) StartTurn() {}

// EndTurn see function from Unit
func (u *Infantry) EndTurn() {
	// turn off `moved` bit
	u.State &= ^UnitStateBitMoved
}
