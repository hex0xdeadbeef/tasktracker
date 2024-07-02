package app

import (
	"fmt"
	"tasktracker/configs"
	"tasktracker/internal/database/postgres"
	"tasktracker/internal/logger"
)

// Run maintains the application's execution and returns an error if any
func Run() (err error) {

	logger, err := logger.Init(&configs.Cfg)
	if err != nil {
		return fmt.Errorf("initializing logger: %w", err)
	}
	defer func() {
		loggerSyncErr := logger.L.Sync()

		if loggerSyncErr == nil {
			return
		}

		if err != nil {
			err = fmt.Errorf("during logger synchronization: %w; %w", loggerSyncErr, err)
		}

		err = loggerSyncErr
	}()

	if err := configs.Load(); err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	db, err := postgres.OpenDB(logger)
	if err != nil {
		return fmt.Errorf("during db creation: %w", err)
	}
	defer func() {
		dbClosingErr := db.Close()

		if dbClosingErr == nil {
			return
		}

		if err != nil {
			err = fmt.Errorf("%w; %w", dbClosingErr, err)
		}

		err = dbClosingErr
	}()

	return nil
}
