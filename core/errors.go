package core

// Check is a convenience method for panicking on errors.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
