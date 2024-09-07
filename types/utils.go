package types

// ContainsString checks if the target is presented in arr.
// This non-generic version of contains is faster than `slices.Contains`.
func ContainsString(target string, arr []string) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}

	return false
}

// ContainsUint16 checks if the target is presented in arr.
// This non-generic version of contains is faster than `slices.Contains`.
func ContainsUint16(target uint16, arr []uint16) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}

	return false
}
