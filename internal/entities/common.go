package entities

func equalSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	unique := make(map[T]bool)
	for _, s := range a {
		unique[s] = true
	}
	for _, s := range b {
		if !unique[s] {
			return false
		}
	}
	return true
}
