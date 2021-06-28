package combat

import (
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
	"gitlab.com/beewar/beewar-be/internal/utils"
)

// Combat does a combat depending on the attacker's attack type
func Combat(attacker, defender objects.Unit, dist int, underDome bool) {
	switch attacker.UnitAttackType() {
	case objects.AttackTypeGround:
		GroundCombat(attacker, defender, dist)
	case objects.AttackTypeAerial:
		AerialCombat(attacker, defender, underDome)
	}
}

// GroundCombat does a normal combat (1 attack, 1 counter-attack) and modifies the given attacker and defender units' hp
func GroundCombat(attacker, defender objects.Unit, dist int) {
	dmg := GroundAttack(attacker, defender, dist)
	defender.SetHP(defender.GetHP() - dmg)
	dmg = GroundAttack(defender, attacker, dist)
	attacker.SetHP(attacker.GetHP() - dmg)
}

// SimulateGroundCombat simulates normal combat: doesn't modify unit hp but return damage dealt instead.
// Damage dealt is returned in this order (damage to attacker, damage to defender)
func SimulateGroundCombat(attacker, defender objects.Unit, dist int) (int, int) {
	dmgDef := GroundAttack(attacker, defender, dist)
	defHP := defender.GetHP()
	defender.SetHP(defHP - dmgDef)
	dmgAtk := GroundAttack(defender, attacker, dist)
	defender.SetHP(defHP) // undo defender dmg
	return dmgAtk, dmgDef
}

// GroundAttack does a single attack only. Returns damage dealt to the defender. Damage dealt cannot exceed remaining hp
func GroundAttack(attacker, defender objects.Unit, dist int) int {
	atkDamage := 0
	switch attacker.UnitAttackType() {
	case objects.AttackTypeNone:
		// do nothing
	case objects.AttackTypeGround:
		if dist <= attacker.UnitAttackRange() {
			atkDamage = utils.CeilDivInt(attacker.GetHP()*attacker.UnitAttackPower(), 10)
		}
	case objects.AttackTypeAerial:
		// aerial attackers cannot respond to ground attacks
	default:
		panic("panic attack: unknown attack type")
	}

	return utils.MinInt(defender.GetHP(), atkDamage)
}

// AerialCombat does a normal aerial attack that cannot be countered and modifies defender hp.
// The combat is under the assumption that the defender is in range of the attacker.
func AerialCombat(attacker, defender objects.Unit, underDome bool) {
	dmgDef := AerialAttack(attacker, defender, underDome)
	defender.SetHP(defender.GetHP() - dmgDef)
}

// AerialAttack does a single aerial attack.
// The attack is under the assumption that the defender is in range of the attacker.
func AerialAttack(attacker, defender objects.Unit, underDome bool) int {
	atkDamage := 0
	if attacker.UnitAttackType() != objects.AttackTypeAerial {
		panic("panic attack: attack type is not aerial")
	}

	if underDome {
		atkDamage = utils.CeilDivInt(attacker.GetHP()*attacker.UnitAttackPower(), 20)
	} else {
		atkDamage = utils.CeilDivInt(attacker.GetHP()*attacker.UnitAttackPower(), 10)
	}

	return utils.MinInt(defender.GetHP(), atkDamage)
}
