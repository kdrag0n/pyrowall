package util

// PanicIf is a convenience method for panicking on non-nil errors.
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
