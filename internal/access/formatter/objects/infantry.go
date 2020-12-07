package objects

// Infantry is a unit object.
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

// GetWeight see function from Unit
func (inf *Infantry) GetWeight() int {
	return UnitWeightInfantry
}

// GetUnitOwner see function from Unit
func (inf *Infantry) GetUnitOwner() int {
	return inf.Owner
}

// GetUnitState see function from Unit
func (inf *Infantry) GetUnitState() int {
	return inf.State
}
