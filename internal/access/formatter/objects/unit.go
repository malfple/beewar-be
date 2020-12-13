package objects

// Unit is an object that describes units on the map
// to implement this interface, the basic fields that a unit struct should have are:
// Owner (owner of this unit)
// HP (health point)
// State (state/flags), which can split into multiple fields, or renamed
type Unit interface {
	// GetUnitType gets the unit type of the current unit object
	GetUnitType() int
	// GetWeight returns the weight characteristic of the unit type
	GetWeight() int
	// GetUnitOwner returns the owner of the unit
	GetUnitOwner() int
	// GetUnitState returns unit states
	GetUnitState() int
	// GetUnitStateBit returns one specified unit state bit as a bool (0 or 1)
	GetUnitStateBit(bit int) bool
	// ToggleUnitStateBit toggles state using the bit given. bit has to be one of the const below
	ToggleUnitStateBit(bit int)
	// StartTurn triggers start-of-turn effects
	StartTurn()
	// EndTurn ends the turn for the unit, reset states and trigger any end-of-turn effects
	EndTurn()
}

// these are the unit types
const (
	// UnitTypeYou defines the unit type number of You
	UnitTypeYou = 1
	// UnitTypeInfantry defines the unit type number of Infantry
	UnitTypeInfantry = 3
)

// unit weights
// 0 = light. 1 = heavy. 2 = unpassable
// weight is used to determine whether a unit can pass another unit.
// 2 units can pass through each other if the sum of their weight <= 1 AND they have the same owner
const (
	// UnitWeightYou defines unit weight of You
	UnitWeightYou = 0
	// UnitWeightInfantry defines unit weight of Infantry
	UnitWeightInfantry = 0
)

// unit steps
// for units that has a maximum movement range
const (
	// UnitMoveStepsYou defines unit movement range of You
	UnitMoveStepsYou = 1
	// UnitMoveStepsInfantry defines unit movement range of Infantry
	UnitMoveStepsInfantry = 3
)

// state bit constants. always in the form of 2^n
const (
	// UnitStateBitMoved defines the bit in unit states that specifies moved state
	UnitStateBitMoved = 1
)
