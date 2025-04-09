package appy_driver

func IDScanner(rr Scannable, entry *uint64) error {
	return rr.Scan(entry)
}

func BoolScanner(rr Scannable, entry *bool) error {
	return rr.Scan(entry)
}

func IntScanner(rr Scannable, entry *int) error {
	return rr.Scan(entry)
}

func StringScanner(rr Scannable, entry *string) error {
	return rr.Scan(entry)
}
