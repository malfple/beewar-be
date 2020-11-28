package units

// Infantry is a unit object.
type Infantry struct {
	P int8
}

// GetUnitType see function from Unit
func (inf *Infantry) GetUnitType() int8 {
	return UnitTypeInfantry
}
