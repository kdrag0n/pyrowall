package util

// Min returns the minimum of the given integers.
func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// Max returns the minimum of the given integers.
func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
