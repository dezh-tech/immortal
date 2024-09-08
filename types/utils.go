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

// ContainsKind checks if the target is presented in arr.
// This non-generic version of contains is faster than `slices.Contains`.
func ContainsKind(target Kind, arr []Kind) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}

	return false
}
