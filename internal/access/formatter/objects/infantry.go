package objects

// Infantry is a unit object.
type Infantry struct {
	Owner uint8
	HP    uint8
	State uint8
}

// NewInfantry returns a new Infantry object
func NewInfantry(owner, hp, state uint8) *Infantry {
	return &Infantry{
		Owner: owner,
		HP:    hp,
		State: state,
	}
}

// GetUnitType see function from Unit
func (inf *Infantry) GetUnitType() uint8 {
	return UnitTypeInfantry
}

// GetUnitState see function from Unit
func (inf *Infantry) GetUnitState() uint8 {
	return inf.State
}
