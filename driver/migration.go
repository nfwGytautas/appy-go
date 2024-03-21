package driver

type MigrationFn func(*Tx, string) error

func MigrateToVersion(tx *Tx, versionName, query string) error {
	_, err := tx.Exec(query)
	if err != nil {
		gLogger.Error("%v", err)
		return err
	}

	_, err = tx.Exec("UPDATE \"driver-internal\" SET \"db_version\" = $1 WHERE \"db_version\" > ''", versionName)
	if err != nil {
		gLogger.Error("%v", err)
		return err
	}

	return nil
}
