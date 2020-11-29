package units

// Unit is an object that describes units on the map
// to implement this interface, the basic fields that a unit struct should have are:
// P (owner of this unit)
type Unit interface {
	GetUnitType() int8
}

// these are the unit types
const (
	// UnitTypeYou defines a unit type
	UnitTypeYou = 1
	// UnitTypeInfantry defines a unit type
	UnitTypeInfantry = 3
)
