package stacks

import "fmt"

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

func (_self *Stack[T]) Push(t T) {
	fmt.Printf("stack.Push(%v)\n", t)
	_self.stack = append(_self.stack, t)
	_self.pointer = len(_self.stack) - 1
}

func (_self *Stack[T]) Peak() T {
	fmt.Printf("stack.Peak() %v\n", _self.stack[_self.pointer])
	return _self.stack[_self.pointer]
}

func (_self *Stack[T]) Pop() T {
	fmt.Printf("stack.Pop() %v\n", _self.stack[_self.pointer])
	var tmp = _self.stack[_self.pointer]
	_self.stack = _self.stack[0:_self.pointer]
	_self.pointer = len(_self.stack) - 1
	return tmp
}
