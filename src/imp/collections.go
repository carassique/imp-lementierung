package imp

type Stack[T any] interface {
	Pop() T
	Peek() T
	Push(val T)
	IsEmpty() bool
}

type SliceStack[T any] struct {
	slice *[]T
}

func (s *SliceStack[T]) Pop() T {
	deref := *s.slice
	value := deref[0]
	*s.slice = deref[1:]
	return value
}

func (s *SliceStack[T]) Peek() T {
	deref := *s.slice
	value := deref[0]
	return value
}

func (s *SliceStack[T]) Push(val T) {
	*s.slice = append(*s.slice, val)
}

func (s *SliceStack[T]) IsEmpty() bool {
	return len(*s.slice) == 0
}

func MakeStack[T any](val ...T) Stack[T] {
	slice := val
	return &SliceStack[T]{slice: &slice}
}
