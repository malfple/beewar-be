package loader

// validate functions return empty string if no errors

// validates of unit is owned by the userID given.
// also validates position inside map, and if a unit exists in the given position
func (gl *GameLoader) validateUnitOwned(userID uint64, y, x int) string {
	if y < 0 || y > gl.Height || x < 0 || x > gl.Width {
		return errMsgInvalidPos
	}
	if gl.Units[y][x] == nil {
		return errMsgInvalidPos
	}
	// player doesn't own the unit
	if gl.UserIDToPlayerMap[userID] != gl.Units[y][x].GetUnitOwner() {
		return errMsgUnitNotOwned
	}

	return ""
}
