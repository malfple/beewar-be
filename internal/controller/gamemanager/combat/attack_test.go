package combat

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
	"testing"
)

func TestCombat(t *testing.T) {
	inf1 := objects.NewInfantry(1, 10, 0)
	inf2 := objects.NewInfantry(2, 10, 0)
	Combat(inf1, inf2, 1, false)
	assert.Equal(t, 7, inf1.HP)
	assert.Equal(t, 5, inf2.HP)
}

// also tests simulate normal combat
func TestGroundCombat(t *testing.T) {
	// infantry attacks infantry
	inf1 := objects.NewInfantry(1, 10, 0)
	inf2 := objects.NewInfantry(2, 10, 0)
	dmgA, dmgD := SimulateGroundCombat(inf1, inf2, 1)
	assert.Equal(t, 3, dmgA)
	assert.Equal(t, 5, dmgD)
	GroundCombat(inf1, inf2, 1)
	assert.Equal(t, 7, inf1.HP)
	assert.Equal(t, 5, inf2.HP)

	inf1.HP = 10
	inf2.HP = 2

	GroundCombat(inf1, inf2, 1)
	assert.Equal(t, 10, inf1.HP)
	assert.Equal(t, 0, inf2.HP)
}

func TestGroundCombat2(t *testing.T) {
	// wizard attack infantry
	wiz := objects.NewWizard(1, 10, 0)
	inf := objects.NewInfantry(2, 10, 0)
	dmgA, dmgD := SimulateGroundCombat(wiz, inf, 1)
	assert.Equal(t, 3, dmgA)
	assert.Equal(t, 5, dmgD)
	dmgA, dmgD = SimulateGroundCombat(wiz, inf, 2)
	assert.Equal(t, 0, dmgA)
	assert.Equal(t, 5, dmgD)
	GroundCombat(wiz, inf, 2)
	assert.Equal(t, 10, wiz.GetHP())
	assert.Equal(t, 5, inf.GetHP())

	// wizard attack jetcrew
	jet := objects.NewJetCrew(2, 8, 0)
	GroundCombat(wiz, jet, 2)
	assert.Equal(t, 10, wiz.GetHP())
	assert.Equal(t, 3, jet.GetHP())
}

func TestAerialCombat(t *testing.T) {
	// mortar attack infantry
	mortar := objects.NewMortar(1, 4, 0)
	inf := objects.NewInfantry(2, 10, 0)
	AerialCombat(mortar, inf, false)
	assert.Equal(t, 6, inf.GetHP())
	inf.SetHP(10)
	AerialCombat(mortar, inf, true)
	assert.Equal(t, 8, inf.GetHP())
}
