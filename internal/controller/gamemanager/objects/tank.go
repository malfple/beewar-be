package objects

import "gitlab.com/beewar/beewar-be/internal/utils"

// Tank is a unit object.
// State description: 1 bit for moved state
type Tank struct {
	Owner int
	HP    int
	State int
}

// NewTank returns a new Tank object
func NewTank(owner, hp, state int) *Tank {
	return &Tank{
		Owner: owner,
		HP:    hp,
		State: state,
	}
}

// UnitType see function from Unit
func (u *Tank) UnitType() int {
	return UnitTypeTank
}

// UnitMaxHP see function from Unit
func (u *Tank) UnitMaxHP() int {
	return UnitMaxHPTank
}

// UnitWeight see function from Unit
func (u *Tank) UnitWeight() int {
	return UnitWeightTank
}

// UnitMoveType see function from Unit
func (u *Tank) UnitMoveType() int {
	return MoveTypeGround
}

// UnitMoveRange see function from Unit
func (u *Tank) UnitMoveRange() int {
	return UnitMoveRangeTank
}

// UnitMoveRangeMin see function from Unit
func (u *Tank) UnitMoveRangeMin() int { return 0 }

// UnitAttackType see function from Unit
func (u *Tank) UnitAttackType() int {
	return AttackTypeGround
}

// UnitAttackRange see function from Unit
func (u *Tank) UnitAttackRange() int {
	return UnitAttackRangeTank
}

// UnitAttackRangeMin see function from Unit
func (u *Tank) UnitAttackRangeMin() int { return 0 }

// UnitAttackPower see function from Unit
func (u *Tank) UnitAttackPower() int {
	return UnitAttackPowerTank
}

// UnitCost see function from Unit
func (u *Tank) UnitCost() int {
	return UnitCostTank
}

// GetOwner see function from Unit
func (u *Tank) GetOwner() int {
	return u.Owner
}

// GetState see function from Unit
func (u *Tank) GetState() int {
	return u.State
}

// GetStateBit see function from Unit
func (u *Tank) GetStateBit(bit int) bool {
	return (u.State & bit) != 0
}

// ToggleStateBit see function from Unit
func (u *Tank) ToggleStateBit(bit int) {
	u.State ^= bit
}

// GetHP see function from Unit
func (u *Tank) GetHP() int {
	return u.HP
}

// SetHP see function from Unit
func (u *Tank) SetHP(hp int) {
	u.HP = hp
}

// StartTurn see function from Unit
func (u *Tank) StartTurn() {}

// EndTurn see function from Unit
func (u *Tank) EndTurn() {
	// turn off `moved` bit
	u.State &= ^UnitStateBitMoved
	// heal
	u.HP = utils.MinInt(u.HP+2, UnitMaxHPTank)
}
