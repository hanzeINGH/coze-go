package coze

func ptrValue[T any](s *T) T {
	if s != nil {
		return *s
	}
	var empty T
	return empty
}

func ptr[T any](s T) *T {
	return &s
}
