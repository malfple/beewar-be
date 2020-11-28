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

// NewUnitFromType creates a new unit based on the type specified
func NewUnitFromType(unitType, owner int8) Unit {
	switch unitType {
	case UnitTypeYou:
		return &You{P: owner}
	case UnitTypeInfantry:
		return &Infantry{P: owner}
	default:
		return nil
	}
}
