package deepswap

func toPointer[T any](t T) *T {
	return &t
}

func toInterface[T any](t T) any {
	return t
}
