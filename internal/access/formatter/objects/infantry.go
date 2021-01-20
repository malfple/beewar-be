package objects

const (
	// UnitWeightInfantry defines weight stat of Infantry
	UnitWeightInfantry = 0
	// UnitMoveRangeInfantry defines movement range stat of Infantry
	UnitMoveRangeInfantry = 3
	// UnitAttackRangeInfantry defines attack range stat of Infantry
	UnitAttackRangeInfantry = 1
	// UnitAttackPowerInfantry defines attack power stat of Infantry
	UnitAttackPowerInfantry = 5
)

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

// GetUnitType see function from Unit
func (inf *Infantry) GetUnitType() int {
	return UnitTypeInfantry
}

// GetUnitOwner see function from Unit
func (inf *Infantry) GetUnitOwner() int {
	return inf.Owner
}

// GetUnitState see function from Unit
func (inf *Infantry) GetUnitState() int {
	return inf.State
}

// GetUnitStateBit see function from Unit
func (inf *Infantry) GetUnitStateBit(bit int) bool {
	return (inf.State & bit) != 0
}

// ToggleUnitStateBit see function from Unit
func (inf *Infantry) ToggleUnitStateBit(bit int) {
	inf.State ^= bit
}

// GetUnitHP see function from Unit
func (inf *Infantry) GetUnitHP() int {
	return inf.HP
}

// SetUnitHP see function from Unit
func (inf *Infantry) SetUnitHP(hp int) {
	inf.HP = hp
}

// GetWeight see function from Unit
func (inf *Infantry) GetWeight() int {
	return UnitWeightInfantry
}

// GetMoveType see function from Unit
func (inf *Infantry) GetMoveType() int {
	return MoveTypeGround
}

// GetMoveRange see function from Unit
func (inf *Infantry) GetMoveRange() int {
	return UnitMoveRangeInfantry
}

// GetAttackType see function from Unit
func (inf *Infantry) GetAttackType() int {
	return AttackTypeGround
}

// GetAttackRange see function frmo Unit
func (inf *Infantry) GetAttackRange() int {
	return UnitAttackRangeInfantry
}

// GetAttackPower see function from Unit
func (inf *Infantry) GetAttackPower() int {
	return UnitAttackPowerInfantry
}

// StartTurn see function from Unit
func (inf *Infantry) StartTurn() {}

// EndTurn see function from Unit
func (inf *Infantry) EndTurn() {
	// turn off `moved` bit
	inf.State &= ^UnitStateBitMoved
}
