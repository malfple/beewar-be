package loader

import (
	"gitlab.com/otqee/otqee-be/internal/access/formatter/objects"
)

// CmdWhitelistUnitMove indicates which unit types can use message.CmdUnitMove
var CmdWhitelistUnitMove = map[int]bool{
	objects.UnitTypeYou:      true,
	objects.UnitTypeInfantry: true,
}
