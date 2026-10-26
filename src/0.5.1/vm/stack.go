package vm

import (
	"fmt"
	"os"
)

const StackSize = 256 //1字节

type Stack struct {
	data [StackSize]byte
	top int
}

func (s *Stack) Push(v byte) error {
	if s.top >= StackSize {
		return panicCode(0x01)
	}
	s.data[s.top] = v
	s.top++
	return nil
}

func (s *Stack) Pop() (byte,error) {
	if s.top <= 0 {
		return 0, panicCode(0x02)
	}
	s.top--
	v := s.data[s.top]
	return v,nil
}

func (s *Stack) Add() error {
	if s.top <=1 {
		return panicCode(0x02)
	}
	s.top--; s.data[s.top - 1] += s.data[s.top]
	return nil
}

func (s *Stack) Sub() error {
	if s.top <=1 {
		return panicCode(0x02)
	}
	s.top--; s.data[s.top - 1] -= s.data[s.top]
	return nil
}

func (s *Stack) Mul() error {
	if s.top <=1 {
		return panicCode(0x02)
	}
	s.top--; s.data[s.top - 1] *= s.data[s.top]
	return nil
}

func (s *Stack) Div() error {
	if s.top <=1 {
		return panicCode(0x02)
	}else if s.data[s.top] == 0 {
		return panicCode(0x06)
	}
	s.top--; s.data[s.top - 1] /= s.data[s.top]
	return nil
}

func (s *Stack) Size() int {
	return s.top
}

func (s *Stack) Look() error {
	if s.top <=0 {
		return panicCode(0x02)
	}
	fmt.Println(s.data[s.top - 1])
	return nil
}

func (s *Stack) Looc() error {
	if s.top <= 0 {
		return panicCode(0x02)
	}
	fmt.Printf("%c\n",s.data[s.top - 1])
	return nil
}

func (s *Stack) Put() error {
	if s.top <= 0 {
		return panicCode(0x02)
	}
	fmt.Print(s.data[s.top - 1])
	return nil
}

func (s *Stack) Putc() error {
	if s.top <= 0 {
		return panicCode(0x02)
	}
	fmt.Printf("%c",s.data[s.top - 1])
	return nil
}

func (s *Stack) Endl() error {
	fmt.Println()
	return nil
}

func (s *Stack) Listen() (byte, error) {
    fmt.Print("Listening> ")
    var i int
    _, err := fmt.Fscanln(os.Stdin, &i)
    if err != nil {
        return 0, err
    }
    return byte(i % 256), nil
}

func NewStack() *Stack {
	return &Stack{
		top: 0,
	}
}
