package utils

// Utility functions related to the grid engine

// HexDistance calculates distance between two positions as if there are no obstacles
func HexDistance(y1, x1, y2, x2 int) int {
	x1 -= (y1 + 1) / 2
	x2 -= (y2 + 1) / 2
	z1 := -x1 - y1
	z2 := -x2 - y2
	return (AbsInt(x1-x2) + AbsInt(y1-y2) + AbsInt(z1-z2)) / 2
}
