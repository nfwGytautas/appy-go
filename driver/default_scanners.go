package appy_driver

func IDScanner(rr Scannable, entry *uint64) error {
	return rr.Scan(entry)
}
