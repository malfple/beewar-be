package objects

// Infantry is a unit object.
// State description: 1 bit for moved state
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

// GetUnitStateBit see function from Unit
func (inf *Infantry) GetUnitStateBit(bit int) bool {
	return (inf.State & bit) != 0
}

// ToggleUnitStateBit see function from Unit
func (inf *Infantry) ToggleUnitStateBit(bit int) {
	inf.State ^= bit
}

// StartTurn see function from Unit
func (inf *Infantry) StartTurn() {}

// EndTurn see function from Unit
func (inf *Infantry) EndTurn() {
	// turn off `moved` bit
	inf.State &= ^UnitStateBitMoved
}
