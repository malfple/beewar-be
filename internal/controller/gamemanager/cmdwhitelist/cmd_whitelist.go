package cmdwhitelist

import (
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
)

// UnitMoveMap indicates which unit types can use message.CmdUnitMove
var UnitMoveMap = map[int]bool{
	objects.UnitTypeQueen:    true,
	objects.UnitTypeInfantry: true,
}

// UnitAttackMap indicates which unit types can use message.CmdUnitAttack
var UnitAttackMap = map[int]bool{
	objects.UnitTypeInfantry: true,
}

// UnitMoveAndAttackMap indicates which unit types can use message.CmdUnitMoveAndAttack
var UnitMoveAndAttackMap = map[int]bool{
	objects.UnitTypeInfantry: true,
}
