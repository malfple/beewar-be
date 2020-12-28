package utils

// AbsInt returns absolute value
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// MaxInt returns max of 2 integers
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
