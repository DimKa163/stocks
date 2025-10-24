package collection

func All[T comparable](slice []T, val T) bool {
	for _, v := range slice {
		if v != val {
			return false
		}
	}
	return true
}
