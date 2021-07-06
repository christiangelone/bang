package sugar

func Has(err error) bool {
	return err != nil
}

func NotHas(err error) bool {
	return err == nil
}

func Is(errA, errB error) bool {
	return errA == errB
}
