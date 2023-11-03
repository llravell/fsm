package fsm

type set[T comparable] struct {
	m map[T]struct{}
}

func newSet[T comparable]() set[T] {
	return set[T]{
		make(map[T]struct{}),
	}
}

func (s *set[T]) Add(key T) {
	s.m[key] = struct{}{}
}

func (s *set[T]) Delete(key T) {
	delete(s.m, key)
}

func (s set[T]) Has(key T) bool {
	_, has := s.m[key]

	return has
}

func (s set[T]) Size() int {
	return len(s.m)
}

func (s set[T]) Keys() []T {
	keys := make([]T, 0, s.Size())

	for key := range s.m {
		keys = append(keys, key)
	}

	return keys
}
