package collection

type MapKey interface {
	comparable
}

func Values[TKey MapKey, TValue any](m map[TKey]TValue) []TValue {
	vals := make([]TValue, len(m))
	i := 0
	for _, v := range m {
		vals[i] = v
		i++
	}
	return vals
}
