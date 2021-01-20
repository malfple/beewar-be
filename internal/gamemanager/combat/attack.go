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
	atkDamage := 0
	switch attacker.GetAttackType() {
	case objects.AttackTypeNone:
		// do nothing
	case objects.AttackTypeGround:
		if dist <= attacker.GetAttackRange() {
			atkDamage = utils.CeilDivInt(attacker.GetUnitHP() * attacker.GetAttackPower(), 10)
		}
	default:
		panic("panic attack: unknown attack type")
	}

	hpDef := utils.MaxInt(defender.GetUnitHP()-atkDamage, 0)
	defender.SetUnitHP(hpDef)
}
