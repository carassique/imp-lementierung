package imp

type Stack[T any] interface {
	pop() T
	peek() T
	push(val T)
	isEmpty() bool
}

type SliceStack[T any] struct {
	slice *[]T
}

func (s *SliceStack[T]) pop() T {
	deref := *s.slice
	value := deref[0]
	*s.slice = deref[1:]
	return value
}

func (s *SliceStack[T]) peek() T {
	deref := *s.slice
	value := deref[0]
	return value
}

func (s *SliceStack[T]) push(val T) {
	*s.slice = append(*s.slice, val)
}

func (s *SliceStack[T]) isEmpty() bool {
	return len(*s.slice) == 0
}

func makeStack[T any](val ...T) Stack[T] {
	slice := val
	return &SliceStack[T]{slice: &slice}
}
