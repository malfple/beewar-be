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

// MinInt returns min of 2 integers
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CeilDivInt return Ceil(a / b), but integer.
// a and b has to be non-negative. No handling for negative yet
func CeilDivInt(a, b int) int {
	res := a / b
	a -= res * b
	if a != 0 {
		res++
	}
	return res
}
