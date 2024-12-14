package appy_driver

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type internalInfo struct {
	DBVersion string
}

// Initialise database
func openConnection(connectionString string) error {
	var err error

	gDatabaseConnection, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return err
	}

	return nil
}

func readInternals() (internalInfo, error) {
	internals := internalInfo{}

	row := gDatabaseConnection.QueryRow(context.Background(), "SELECT * FROM \"driver-internal\" LIMIT 1;")
	err := row.Scan(&internals.DBVersion)

	return internals, err
}

func getDatabaseVersion() string {
	internals, err := readInternals()

	if err != nil {
		return "nil"
	}

	return internals.DBVersion
}
