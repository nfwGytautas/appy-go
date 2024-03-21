package driver

var IgnoreColumn ignoreColumn

type ignoreColumn struct{}

func (ignoreColumn) Scan(value interface{}) error {
	return nil
}

type ScanFn[T any] func(row Scannable, entry *T) error

func ParseRows[T any](rows RowsResult, scanFn ScanFn[T]) ([]T, error) {
	defer rows.Close()

	entries := []T{}

	for rows.Next() {
		var entry T
		err := scanFn(rows, &entry)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func ParseRow[T any](row RowResult, scanFn ScanFn[T]) (*T, error) {
	var entry T
	err := scanFn(row, &entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

type ExecFn func(*Tx, ...any) (ExecResult, error)

type QueryFn[Rt Scannable] func(*Tx, ...any) (Rt, error)

func QueryRows[T any, Rt Scannable](tx *Tx, scanFn ScanFn[T], query QueryFn[Rt], args ...any) ([]T, error) {
	res, err := query(tx, args...)
	if err != nil {
		return nil, err
	}

	return parseScannable(res, scanFn)
}

func parseScannable[T any](rows Scannable, scanFn ScanFn[T]) ([]T, error) {
	defer rows.Close()

	entries := []T{}

	for rows.Next() {
		var entry T
		err := scanFn(rows, &entry)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return entries, nil
}
