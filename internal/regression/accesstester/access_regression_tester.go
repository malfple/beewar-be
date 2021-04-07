package accesstester

import (
	"database/sql"
	"errors"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
)

// RunAccessRegressionTester runs the regression tests for access layer.
// the db does not need to be seeded
// may result in dirty data if the test fails mid-way
func RunAccessRegressionTester() bool {
	logger.GetLogger().Info("access regression test")

	if !testExecWithTransaction() {
		return false
	}
	if !TestUserAccess() {
		return false
	}
	if !TestMapAccess() {
		return false
	}
	if !TestGameAccess() {
		return false
	}
	return true
}

func testExecWithTransaction() bool {
	logger.GetLogger().Info("test ExecWithTransaction")

	// test do nothing
	err := access.ExecWithTransaction(func(tx *sql.Tx) error {
		// do nothing
		return nil
	})
	if err != nil {
		logger.GetLogger().Error("error do nothing", zap.Error(err))
		return false
	}

	// test failure
	if err = access.CreateGameUser(0, 0, 0); err != nil {
		logger.GetLogger().Error("fail to create dummy gameuser")
		return false
	}
	gu := access.QueryGameUsersByGameID(0)[0]
	gu.PlayerOrder = 3
	err = access.ExecWithTransaction(func(tx *sql.Tx) error {
		// try update
		err2 := access.UpdateGameUserUsingTx(tx, gu)
		if err2 != nil {
			return err2
		}

		return errors.New("fail haha")
	})
	if err == nil {
		logger.GetLogger().Error("should have returned error")
		return false
	}
	if err.Error() != "fail haha" {
		logger.GetLogger().Error("error do something", zap.Error(err))
		return false
	}
	// check to see if it's rollback-ed
	if gu2 := access.QueryGameUsersByGameID(0)[0]; gu2.PlayerOrder != 0 {
		logger.GetLogger().Error("not rollback-ed")
		return false
	}

	return true
}
