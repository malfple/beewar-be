package objects

// JetCrew is a unit object.
// State description: 1 bit for moved state
type JetCrew struct {
	Owner int
	HP    int
	State int
}

// NewJetCrew returns a new JetCrew object
func NewJetCrew(owner, hp, state int) *JetCrew {
	return &JetCrew{
		Owner: owner,
		HP:    hp,
		State: state,
	}
}

// UnitType see function from Unit
func (u *JetCrew) UnitType() int {
	return UnitTypeJetCrew
}

// UnitMaxHP see function from Unit
func (u *JetCrew) UnitMaxHP() int {
	return UnitMaxHPJetCrew
}

// UnitWeight see function from Unit
func (u *JetCrew) UnitWeight() int {
	return UnitWeightJetCrew
}

// UnitMoveType see function from Unit
func (u *JetCrew) UnitMoveType() int {
	return MoveTypeGround
}

// UnitMoveRange see function from Unit
func (u *JetCrew) UnitMoveRange() int {
	return UnitMoveRangeJetCrew
}

// UnitMoveRangeMin see function from Unit
func (u *JetCrew) UnitMoveRangeMin() int { return 0 }

// UnitAttackType see function from Unit
func (u *JetCrew) UnitAttackType() int {
	return AttackTypeGround
}

// UnitAttackRange see function from Unit
func (u *JetCrew) UnitAttackRange() int {
	return UnitAttackRangeJetCrew
}

// UnitAttackRangeMin see function from Unit
func (u *JetCrew) UnitAttackRangeMin() int { return 0 }

// UnitAttackPower see function from Unit
func (u *JetCrew) UnitAttackPower() int {
	return UnitAttackPowerJetCrew
}

// UnitCost see function from Unit
func (u *JetCrew) UnitCost() int {
	return UnitCostJetCrew
}

// GetOwner see function from Unit
func (u *JetCrew) GetOwner() int {
	return u.Owner
}

// GetState see function from Unit
func (u *JetCrew) GetState() int {
	return u.State
}

// GetStateBit see function from Unit
func (u *JetCrew) GetStateBit(bit int) bool {
	return (u.State & bit) != 0
}

// ToggleStateBit see function from Unit
func (u *JetCrew) ToggleStateBit(bit int) {
	u.State ^= bit
}

// GetHP see function from Unit
func (u *JetCrew) GetHP() int {
	return u.HP
}

// SetHP see function from Unit
func (u *JetCrew) SetHP(hp int) {
	u.HP = hp
}

// StartTurn see function from Unit
func (u *JetCrew) StartTurn() {}

// EndTurn see function from Unit
func (u *JetCrew) EndTurn() {
	// turn off `moved` bit
	u.State &= ^UnitStateBitMoved
}
