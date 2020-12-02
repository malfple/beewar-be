package objects

// Unit is an object that describes units on the map
// to implement this interface, the basic fields that a unit struct should have are:
// Owner (owner of this unit)
// HP (health point)
// State (state/flags), which can split into multiple fields, or renamed
type Unit interface {
	// GetUnitType gets the unit type of the current unit object
	GetUnitType() uint8
	// GetUnitState combines unit states and return it in one byte
	GetUnitState() uint8
}

// these are the unit types
const (
	// UnitTypeYou defines a unit type
	UnitTypeYou = 1
	// UnitTypeInfantry defines a unit type
	UnitTypeInfantry = 3
)
