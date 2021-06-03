package combat

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
	"testing"
)

// also tests simulate normal combat
func TestNormalCombat(t *testing.T) {
	// infantry attacks infantry
	inf1 := &objects.Infantry{
		Owner: 1,
		HP:    10,
		State: 0,
	}
	inf2 := &objects.Infantry{
		Owner: 2,
		HP:    10,
		State: 0,
	}
	dmgA, dmgD := SimulateNormalCombat(inf1, inf2, 1)
	assert.Equal(t, 3, dmgA)
	assert.Equal(t, 5, dmgD)
	NormalCombat(inf1, inf2, 1)
	assert.Equal(t, 7, inf1.HP)
	assert.Equal(t, 5, inf2.HP)

	inf1.HP = 10
	inf2.HP = 2

	NormalCombat(inf1, inf2, 1)
	assert.Equal(t, 10, inf1.HP)
	assert.Equal(t, 0, inf2.HP)
}
