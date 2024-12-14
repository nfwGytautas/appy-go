package appy_driver

func IDScanner(rr Scannable, entry *uint64) error {
	return rr.Scan(entry)
}

func BoolScanner(rr Scannable, entry *bool) error {
	var i int

	err := rr.Scan(&i)
	if err != nil {
		return err
	}

	if i == 0 {
		*entry = false
	} else {
		*entry = true
	}

	return nil
}
