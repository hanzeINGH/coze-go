package internal

func Value[T any](s *T) T {
	if s != nil {
		return *s
	}
	var empty T
	return empty
}

func Ptr[T any](s T) *T {
	return &s
}
