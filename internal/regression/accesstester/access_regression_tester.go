package accesstester

import "gitlab.com/otqee/otqee-be/internal/logger"

// RunAccessRegressionTester runs the regression tests for access layer.
// the db does not need to be seeded
// may result in dirty data if the test fails mid-way
func RunAccessRegressionTester() bool {
	logger.GetLogger().Info("access regression test")
	if !TestUserAccess() {
		return false
	}
	return true
}
