package combat

import (
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
	"gitlab.com/beewar/beewar-be/internal/utils"
)

// NormalCombat does a normal combat (1 attack, 1 counter-attack) and modifies the given attacker and defender units' hp
func NormalCombat(attacker, defender objects.Unit, dist int) {
	dmg := Attack(attacker, defender, dist)
	defender.SetUnitHP(defender.GetUnitHP() - dmg)
	dmg = Attack(defender, attacker, dist)
	attacker.SetUnitHP(attacker.GetUnitHP() - dmg)
}

// SimulateNormalCombat simulates normal combat: doesn't modify unit hp but return damage dealt instead.
// Damage dealt is returned in this order (damage to attacker, damage to defender)
func SimulateNormalCombat(attacker, defender objects.Unit, dist int) (int, int) {
	dmgDef := Attack(attacker, defender, dist)
	defHP := defender.GetUnitHP()
	defender.SetUnitHP(defHP - dmgDef)
	dmgAtk := Attack(defender, attacker, dist)
	defender.SetUnitHP(defHP) // undo defender dmg
	return dmgAtk, dmgDef
}

// Attack does a single attack only. Returns damage dealt to the defender. Damage dealt cannot exceed remaining hp
func Attack(attacker, defender objects.Unit, dist int) int {
	atkDamage := 0
	switch attacker.GetAttackType() {
	case objects.AttackTypeNone:
		// do nothing
	case objects.AttackTypeGround:
		if dist <= attacker.GetAttackRange() {
			atkDamage = utils.CeilDivInt(attacker.GetUnitHP()*attacker.GetAttackPower(), 10)
		}
	default:
		panic("panic attack: unknown attack type")
	}

	return utils.MinInt(defender.GetUnitHP(), atkDamage)
}
