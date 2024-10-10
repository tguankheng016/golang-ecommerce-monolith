package helpers

func MapContains[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}
