package combat

import (
	"gitlab.com/beewar/beewar-be/internal/access/formatter/objects"
	"gitlab.com/beewar/beewar-be/internal/utils"
)

// NormalCombat does a normal combat (1 attack, 1 counter-attack) and modifies the given attacker and defender units' hp
func NormalCombat(attacker, defender objects.Unit, dist int) {
	Attack(attacker, defender, dist)
	Attack(defender, attacker, dist)
}

// Attack does a single attack only. Modifies the defender hp
func Attack(attacker, defender objects.Unit, dist int) {
	atkPower := 0
	switch attacker.GetUnitType() {
	case objects.UnitTypeInfantry:
		if dist <= objects.UnitAttackRangeInfantry {
			atkPower = (attacker.GetUnitHP() + 1) / 2
		}
	default:
		panic("panic attack: unknown attacker unit type")
	}

	hpDef := utils.MaxInt(defender.GetUnitHP()-atkPower, 0)
	defender.SetUnitHP(hpDef)
}
