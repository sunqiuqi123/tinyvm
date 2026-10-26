package main

import (
	"fmt"
	"tinyvm/vm"
)

func main() {
	stack1 := vm.NewStack()
	stack1.Push(10)
	stack1.Push(20)
	stack1.Push(30)
	stack1.Pop()
	fmt.Println(stack1.Look())
}
