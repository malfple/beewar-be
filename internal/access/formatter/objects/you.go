package objects

const (
	// UnitWeightYou defines weight stat of You
	UnitWeightYou = 0
	// UnitMoveRangeYou defines movement range stat of You
	UnitMoveRangeYou = 1
	// UnitAttackRangeYou defines attack range stat of You
	UnitAttackRangeYou = 0
	// UnitAttackPowerYou doesn't mean anything because You cannot attack
	UnitAttackPowerYou = 0
)

// You is a unit object. It is most often the win or lose condition in a game.
// State description: 1 bit for moved state
type You struct {
	Owner int
	HP    int
	State int
}

// NewYou returns a new You object
func NewYou(owner, hp, state int) *You {
	return &You{
		Owner: owner,
		HP:    hp,
		State: state,
	}
}

// GetUnitType see function from Unit
func (you *You) GetUnitType() int {
	return UnitTypeYou
}

// GetUnitOwner see function from Unit
func (you *You) GetUnitOwner() int {
	return you.Owner
}

// GetUnitState see function from Unit
func (you *You) GetUnitState() int {
	return you.State
}

// GetUnitStateBit see function from Unit
func (you *You) GetUnitStateBit(bit int) bool {
	return (you.State & bit) != 0
}

// ToggleUnitStateBit see function from Unit
func (you *You) ToggleUnitStateBit(bit int) {
	you.State ^= bit
}

// GetUnitHP see function from Unit
func (you *You) GetUnitHP() int {
	return you.HP
}

// SetUnitHP see function from Unit
func (you *You) SetUnitHP(hp int) {
	you.HP = hp
}

// GetWeight see function from Unit
func (you *You) GetWeight() int {
	return UnitWeightYou
}

// GetMoveType see function from Unit
func (you *You) GetMoveType() int {
	return MoveTypeGround
}

// GetMoveRange see function from Unit
func (you *You) GetMoveRange() int {
	return UnitMoveRangeYou
}

// GetAttackType see function from Unit
func (you *You) GetAttackType() int {
	return AttackTypeNone
}

// GetAttackRange see function from Unit
func (you *You) GetAttackRange() int {
	return UnitAttackRangeYou
}

// GetAttackPower see function from Unit
func (you *You) GetAttackPower() int {
	return UnitAttackPowerYou
}

// StartTurn see function from Unit
func (you *You) StartTurn() {}

// EndTurn see function from Unit
func (you *You) EndTurn() {
	// turn off `moved` bit
	you.State &= ^UnitStateBitMoved
}
