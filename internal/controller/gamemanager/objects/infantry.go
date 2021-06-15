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
func (inf *Infantry) UnitType() int {
	return UnitTypeInfantry
}

// UnitMaxHP see function from Unit
func (inf *Infantry) UnitMaxHP() int {
	return UnitMaxHPInfantry
}

// UnitWeight see function from Unit
func (inf *Infantry) UnitWeight() int {
	return UnitWeightInfantry
}

// UnitMoveType see function from Unit
func (inf *Infantry) UnitMoveType() int {
	return MoveTypeGround
}

// UnitMoveRange see function from Unit
func (inf *Infantry) UnitMoveRange() int {
	return UnitMoveRangeInfantry
}

// UnitAttackType see function from Unit
func (inf *Infantry) UnitAttackType() int {
	return AttackTypeGround
}

// UnitAttackRange see function frmo Unit
func (inf *Infantry) UnitAttackRange() int {
	return UnitAttackRangeInfantry
}

// UnitAttackPower see function from Unit
func (inf *Infantry) UnitAttackPower() int {
	return UnitAttackPowerInfantry
}

// UnitCost see function from Unit
func (inf *Infantry) UnitCost() int {
	return UnitCostInfantry
}

// GetOwner see function from Unit
func (inf *Infantry) GetOwner() int {
	return inf.Owner
}

// GetState see function from Unit
func (inf *Infantry) GetState() int {
	return inf.State
}

// GetStateBit see function from Unit
func (inf *Infantry) GetStateBit(bit int) bool {
	return (inf.State & bit) != 0
}

// ToggleStateBit see function from Unit
func (inf *Infantry) ToggleStateBit(bit int) {
	inf.State ^= bit
}

// GetHP see function from Unit
func (inf *Infantry) GetHP() int {
	return inf.HP
}

// SetHP see function from Unit
func (inf *Infantry) SetHP(hp int) {
	inf.HP = hp
}

// StartTurn see function from Unit
func (inf *Infantry) StartTurn() {}

// EndTurn see function from Unit
func (inf *Infantry) EndTurn() {
	// turn off `moved` bit
	inf.State &= ^UnitStateBitMoved
}
