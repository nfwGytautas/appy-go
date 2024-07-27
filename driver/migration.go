package appy_driver

import appy_logger "github.com/nfwGytautas/appy-go/logger"

type MigrationFn func(*Tx, string) error

func MigrateToVersion(tx *Tx, versionName, query string) error {
	_, err := tx.Exec(query)
	if err != nil {
		appy_logger.Logger().Error("%v", err)
		return err
	}

	_, err = tx.Exec("UPDATE \"driver-internal\" SET \"db_version\" = $1 WHERE \"db_version\" > ''", versionName)
	if err != nil {
		appy_logger.Logger().Error("%v", err)
		return err
	}

	return nil
}
