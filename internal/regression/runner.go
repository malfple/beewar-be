package regression

import "gitlab.com/otqee/otqee-be/internal/regression/accesstester"

// RunRegressionTests runs the regression tests for this whole project.
// this function does not handle initializations, which should be handled by the caller of this function.
func RunRegressionTests() bool {
	if !accesstester.RunAccessRegressionTester() {
		return false
	}
	return true
}
