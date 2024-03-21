package driver

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nfwGytautas/appy"
)

type InitializeArgs struct {
	ConnectionString string
	Version          string
	Migration        MigrationFn
}

var gDatabaseConnection *pgxpool.Pool
var gLogger appy.Logger

func Initialize(logger appy.Logger, args InitializeArgs) error {
	gLogger = logger

	gLogger.Info("Initializing driver, version: '%s'", args.Version)

	// Open connection
	err := openConnection(args.ConnectionString)
	if err != nil {
		gLogger.Error("Failed to open connection")
		return err
	}

	// Get version
	currentVersion := getDatabaseVersion()

	// Check for migration
	tx, err := StartTransaction()
	if err != nil {
		gLogger.Error("Failed to start transaction")
		return err
	}
	defer tx.Rollback()

	gLogger.Info("Migrating database to '%s'", args.Version)
	err = args.Migration(tx, currentVersion)
	if err != nil {
		gLogger.Error("Failed to migrate to the correct datamodel version")
		return err
	}

	err = tx.Commit()
	if err != nil {
		gLogger.Error("Failed to commit migrations")
		return err
	}

	return nil
}
