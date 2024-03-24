package driver

import (
	"github.com/jackc/pgx/v5/pgxpool"
	appy_logger "github.com/nfwGytautas/appy/logger"
)

type InitializeArgs struct {
	ConnectionString string
	Version          string
	Migration        MigrationFn
}

var gDatabaseConnection *pgxpool.Pool

func Initialize(args InitializeArgs) error {
	appy_logger.Get().Info("Initializing driver, version: '%s'", args.Version)

	// Open connection
	err := openConnection(args.ConnectionString)
	if err != nil {
		appy_logger.Get().Error("Failed to open connection")
		return err
	}

	// Get version
	currentVersion := getDatabaseVersion()

	// Check for migration
	tx, err := StartTransaction()
	if err != nil {
		appy_logger.Get().Error("Failed to start transaction")
		return err
	}
	defer tx.Rollback()

	appy_logger.Get().Info("Migrating database to '%s'", args.Version)
	err = args.Migration(tx, currentVersion)
	if err != nil {
		appy_logger.Get().Error("Failed to migrate to the correct datamodel version")
		return err
	}

	err = tx.Commit()
	if err != nil {
		appy_logger.Get().Error("Failed to commit migrations")
		return err
	}

	return nil
}
