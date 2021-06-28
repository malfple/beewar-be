package objects

// Mortar is a unit object.
// State description: 1 bit for moved state
type Mortar struct {
	Owner int
	HP    int
	State int
}

// NewMortar returns a new Wizard object
func NewMortar(owner, hp, state int) *Mortar {
	return &Mortar{
		Owner: owner,
		HP:    hp,
		State: state,
	}
}

// UnitType see function from Unit
func (u *Mortar) UnitType() int {
	return UnitTypeMortar
}

// UnitMaxHP see function from Unit
func (u *Mortar) UnitMaxHP() int {
	return UnitMaxHPMortar
}

// UnitWeight see function from Unit
func (u *Mortar) UnitWeight() int {
	return UnitWeightMortar
}

// UnitMoveType see function from Unit
func (u *Mortar) UnitMoveType() int {
	return MoveTypeGround
}

// UnitMoveRange see function from Unit
func (u *Mortar) UnitMoveRange() int {
	return UnitMoveRangeMortar
}

// UnitMoveRangeMin see function from Unit
func (u *Mortar) UnitMoveRangeMin() int { return 0 }

// UnitAttackType see function from Unit
func (u *Mortar) UnitAttackType() int {
	return AttackTypeAerial
}

// UnitAttackRange see function from Unit
func (u *Mortar) UnitAttackRange() int {
	return UnitAttackRangeMortarMax
}

// UnitAttackRangeMin see function from Unit
func (u *Mortar) UnitAttackRangeMin() int {
	return UnitAttackRangeMortarMin
}

// UnitAttackPower see function from Unit
func (u *Mortar) UnitAttackPower() int {
	return UnitAttackPowerMortar
}

// UnitCost see function from Unit
func (u *Mortar) UnitCost() int {
	return UnitCostMortar
}

// GetOwner see function from Unit
func (u *Mortar) GetOwner() int {
	return u.Owner
}

// GetState see function from Unit
func (u *Mortar) GetState() int {
	return u.State
}

// GetStateBit see function from Unit
func (u *Mortar) GetStateBit(bit int) bool {
	return (u.State & bit) != 0
}

// ToggleStateBit see function from Unit
func (u *Mortar) ToggleStateBit(bit int) {
	u.State ^= bit
}

// GetHP see function from Unit
func (u *Mortar) GetHP() int {
	return u.HP
}

// SetHP see function from Unit
func (u *Mortar) SetHP(hp int) {
	u.HP = hp
}

// StartTurn see function from Unit
func (u *Mortar) StartTurn() {}

// EndTurn see function from Unit
func (u *Mortar) EndTurn() {
	// turn off `moved` bit
	u.State &= ^UnitStateBitMoved
}
