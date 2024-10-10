package helpers

func SliceToMap[T comparable](s []T) map[T]struct{} {
	m := make(map[T]struct{})
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}

func SliceContains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
