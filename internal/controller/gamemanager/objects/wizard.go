package objects

// Wizard is a unit object.
// State description: 1 bit for moved state
type Wizard struct {
	Owner int
	HP    int
	State int
}

// NewWizard returns a new Wizard object
func NewWizard(owner, hp, state int) *Wizard {
	return &Wizard{
		Owner: owner,
		HP:    hp,
		State: state,
	}
}

// UnitType see function from Unit
func (u *Wizard) UnitType() int {
	return UnitTypeWizard
}

// UnitMaxHP see function from Unit
func (u *Wizard) UnitMaxHP() int {
	return UnitMaxHPWizard
}

// UnitWeight see function from Unit
func (u *Wizard) UnitWeight() int {
	return UnitWeightWizard
}

// UnitMoveType see function from Unit
func (u *Wizard) UnitMoveType() int {
	return MoveTypeBlink
}

// UnitMoveRange see function from Unit
func (u *Wizard) UnitMoveRange() int {
	return UnitMoveRangeWizardMax
}

// UnitMoveRangeMin see function from Unit
func (u *Wizard) UnitMoveRangeMin() int {
	return UnitMoveRangeWizardMin
}

// UnitAttackType see function from Unit
func (u *Wizard) UnitAttackType() int {
	return AttackTypeGround
}

// UnitAttackRange see function from Unit
func (u *Wizard) UnitAttackRange() int {
	return UnitAttackRangeWizard
}

// UnitAttackRangeMin see function from Unit
func (u *Wizard) UnitAttackRangeMin() int { return 0 }

// UnitAttackPower see function from Unit
func (u *Wizard) UnitAttackPower() int {
	return UnitAttackPowerWizard
}

// UnitCost see function from Unit
func (u *Wizard) UnitCost() int {
	return UnitCostWizard
}

// GetOwner see function from Unit
func (u *Wizard) GetOwner() int {
	return u.Owner
}

// GetState see function from Unit
func (u *Wizard) GetState() int {
	return u.State
}

// GetStateBit see function from Unit
func (u *Wizard) GetStateBit(bit int) bool {
	return (u.State & bit) != 0
}

// ToggleStateBit see function from Unit
func (u *Wizard) ToggleStateBit(bit int) {
	u.State ^= bit
}

// GetHP see function from Unit
func (u *Wizard) GetHP() int {
	return u.HP
}

// SetHP see function from Unit
func (u *Wizard) SetHP(hp int) {
	u.HP = hp
}

// StartTurn see function from Unit
func (u *Wizard) StartTurn() {}

// EndTurn see function from Unit
func (u *Wizard) EndTurn() {
	// turn off `moved` bit
	u.State &= ^UnitStateBitMoved
}
