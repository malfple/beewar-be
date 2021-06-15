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
	// UnitAttackType returns attack type of the unit type
	UnitAttackType() int
	// UnitAttackRange returns attack range of the unit type
	UnitAttackRange() int
	// UnitAttackPower returns attack power of the unit type, multiplied by 10
	UnitAttackPower() int
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
	// UnitCost returns value of the unit type.
	UnitCost() int
	// StartTurn triggers start-of-turn effects
	StartTurn()
	// EndTurn ends the turn for the unit, reset states and trigger any end-of-turn effects
	EndTurn()
}

// these are the unit types
const (
	// UnitTypeQueen defines the unit type number of Queen
	UnitTypeQueen = 1
	// UnitTypeInfantry defines the unit type number of Infantry
	UnitTypeInfantry = 3
)

// max hp
const (
	// UnitMaxHPQueen defines the maximum hp for Queen
	UnitMaxHPQueen = 10
	// UnitMaxHPInfantry defines the maximum hp for Infantry
	UnitMaxHPInfantry = 10
)

// move types
// move types of a specific unit type defined in each unit file
const (
	// MoveTypeNone means the unit cannot move
	MoveTypeNone = 0
	// MoveTypeGround is a normal ground move. Shortest path can be used to check this
	MoveTypeGround = 1
)

// unit weights
// 0 = light. 1 = heavy. 2 = unpassable
// weight is used to determine whether a unit can pass another unit.
// 2 units can pass through each other if the sum of their weight <= 1 AND they have the same owner
const (
	// UnitWeightQueen defines weight stat of Queen
	UnitWeightQueen = 0
	// UnitWeightInfantry defines weight stat of Infantry
	UnitWeightInfantry = 0
)

// unit move range
const (
	// UnitMoveRangeQueen defines movement range stat of Queen
	UnitMoveRangeQueen = 1
	// UnitMoveRangeInfantry defines movement range stat of Infantry
	UnitMoveRangeInfantry = 3
)

// unit attack types
// attack types of a specific unit type defined in each unit file
const (
	// AttackTypeNone means the unit cannot attack
	AttackTypeNone = 0
	// AttackTypeGround is a normal melee attack
	AttackTypeGround = 1
)

// unit attack range
const (
	// UnitAttackRangeQueen defines attack range stat of Queen
	UnitAttackRangeQueen = 0
	// UnitAttackRangeInfantry defines attack range stat of Infantry
	UnitAttackRangeInfantry = 1
)

// unit attack power
// this is multiplied by 10 to avoid floating point. So 5 is actually 0.5
const (
	// UnitAttackPowerQueen doesn't mean anything because Queen cannot attack
	UnitAttackPowerQueen = 0
	// UnitAttackPowerInfantry defines attack power stat of Infantry
	UnitAttackPowerInfantry = 5
)

// state bit constants. always in the form of 2^n
const (
	// UnitStateBitMoved defines the bit in unit states that specifies moved state
	UnitStateBitMoved = 1
)

// unit cost
const (
	// UnitCostQueen defines the cost of Queen
	UnitCostQueen = 10000
	// UnitCostInfantry defines the cost of Infantry
	UnitCostInfantry = 300
)
