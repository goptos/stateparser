package stacks

type Stack[T any] struct {
	stack   []T
	pointer int
}

func New[T any]() Stack[T] {
	return Stack[T]{
		stack:   []T{},
		pointer: -1,
	}
}

func (_self *Stack[T]) Depth() int {
	return _self.pointer
}

func (_self *Stack[T]) Push(t T) {
	_self.stack = append(_self.stack, t)
	_self.pointer = len(_self.stack) - 1
}

func (_self *Stack[T]) Peak() T {
	return _self.stack[_self.pointer]
}

func (_self *Stack[T]) Pop() T {
	var tmp = _self.stack[_self.pointer]
	_self.stack = _self.stack[0:_self.pointer]
	_self.pointer = len(_self.stack) - 1
	return tmp
}
