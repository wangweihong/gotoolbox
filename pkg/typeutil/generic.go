package typeutil

func GenericIndirectValue[T any](p *T) T {
	var zero T
	if p != nil {
		return *p
	}
	return zero
}
