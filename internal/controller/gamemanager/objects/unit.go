package objects

/*
Unit is an object that describes units on the map.
To implement this interface, the basic fields that a unit struct should have are:
 - Owner (owner of this unit)
 - HP (health point)
 - State (state/flags), which can split into multiple fields, or renamed
*/
type Unit interface {
	// UnitType gets the unit type of the current unit object
	UnitType() int
	// UnitMaxHP returns the maximum hp of the unit type
	UnitMaxHP() int
	// UnitWeight returns the weight characteristic of the unit type
	UnitWeight() int
	// UnitMoveType returns movement type of the unit type
	UnitMoveType() int
	// UnitMoveRange returns movement range of the unit type
	UnitMoveRange() int
	// UnitMoveRangeMin returns min movement range of the unit type (if exists)
	UnitMoveRangeMin() int
	// UnitAttackType returns attack type of the unit type
	UnitAttackType() int
	// UnitAttackRange returns attack range of the unit type
	UnitAttackRange() int
	// UnitAttackRangeMin returns min attack range of the unit type (if exists)
	UnitAttackRangeMin() int
	// UnitAttackPower returns attack power of the unit type, multiplied by 10
	UnitAttackPower() int
	// UnitCost returns value of the unit type.
	UnitCost() int
	// GetOwner returns the owner of the unit
	GetOwner() int
	// GetState returns unit states
	GetState() int
	// GetStateBit returns one specified unit state bit as a bool (0 or 1)
	GetStateBit(bit int) bool
	// ToggleStateBit toggles state using the bit given. bit has to be one of the const below
	ToggleStateBit(bit int)
	// GetHP returns hp of the unit
	GetHP() int
	// SetHP sets the hp of the unit
	SetHP(hp int)
	// StartTurn triggers start-of-turn effects
	StartTurn()
	// EndTurn ends the turn for the unit, reset states and trigger any end-of-turn effects
	EndTurn()
}

// UnitType defines the unit type number for each unit
const (
	UnitTypeQueen    = 1
	UnitTypeInfantry = 3
	UnitTypeJetCrew  = 4
	UnitTypeWizard   = 5
	UnitTypeTank     = 6
	UnitTypeMortar   = 9
)

// UnitMaxHP defines the maximum hp for each unit
const (
	UnitMaxHPQueen    = 10
	UnitMaxHPInfantry = 10
	UnitMaxHPJetCrew  = 8
	UnitMaxHPWizard   = 10
	UnitMaxHPTank     = 14
	UnitMaxHPMortar   = 4
)

// move types
// move types of a specific unit type defined in each unit file.
const (
	// MoveTypeNone means the unit cannot move
	MoveTypeNone = 0
	// MoveTypeGround is a normal ground move. Shortest path can be used to check this
	MoveTypeGround = 1
	// MoveTypeBlink instantly teleports to any cell inside range without any movement penalty
	MoveTypeBlink = 2
)

// unit weights
// 0 = light. 1 = heavy. 2 = unpassable
// weight is used to determine whether a unit can pass another unit.
// 2 units can pass through each other if the sum of their weight <= 1 AND they have the same owner

// UnitWeight defines the weight stat for each unit
const (
	UnitWeightQueen    = 0
	UnitWeightInfantry = 0
	UnitWeightJetCrew  = 0
	UnitWeightWizard   = 0
	UnitWeightTank     = 1
	UnitWeightMortar   = 0
)

// UnitMoveRange defines the movement range for each unit.
// When appended with Min or Max, it becomes the lower or upper limit of the movement range.
const (
	UnitMoveRangeQueen     = 1
	UnitMoveRangeInfantry  = 3
	UnitMoveRangeJetCrew   = 6
	UnitMoveRangeWizardMin = 2
	UnitMoveRangeWizardMax = 3
	UnitMoveRangeTank      = 4
	UnitMoveRangeMortar    = 4
)

// unit attack types
// attack types of a specific unit type defined in each unit file
const (
	// AttackTypeNone means the unit cannot attack
	AttackTypeNone = 0
	// AttackTypeGround is a normal melee attack
	AttackTypeGround = 1
	// AttackTypeAerial is a ranged aerial attack and can be mitigated with domes
	AttackTypeAerial = 2
)

// UnitAttackRange defines the attack range stat of each unit.
// When appended with Min or Max, it becomes the lower or upper limit of the attack range.
const (
	UnitAttackRangeInfantry  = 1
	UnitAttackRangeJetCrew   = 1
	UnitAttackRangeWizard    = 2
	UnitAttackRangeTank      = 1
	UnitAttackRangeMortarMin = 2
	UnitAttackRangeMortarMax = 3
)

// UnitAttackPower defines the attack power stat of each unit.
// this is multiplied by 10 to avoid floating point. So 5 is actually 0.5
const (
	UnitAttackPowerInfantry = 5
	UnitAttackPowerJetCrew  = 5
	UnitAttackPowerWizard   = 5
	UnitAttackPowerTank     = 5
	UnitAttackPowerMortar   = 10
)

// state bit constants. always in the form of 2^n
const (
	// UnitStateBitMoved defines the bit in unit states that specifies moved state
	UnitStateBitMoved = 1
)

// UnitCost defines the cost of each unit
const (
	UnitCostQueen    = 10000
	UnitCostInfantry = 300
	UnitCostJetCrew  = 320
	UnitCostWizard   = 500
	UnitCostTank     = 700
	UnitCostMortar   = 800
)
