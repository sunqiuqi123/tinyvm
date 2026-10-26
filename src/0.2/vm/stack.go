package vm

import "fmt"

type Stack struct {
	data [256]byte
	top int
}

func (s *Stack) Push(v byte) error {
	if s.top >= 256 {
		return panicCode(0x01)
	}
	s.data[s.top] = v
	s.top++
	return nil
}

func (s *Stack) Pop() (byte,error) {
	if s.top <= 0 {
		return 0,panicCode(0x02)
	}
	s.top--
	v := s.data[s.top]
	return v,nil
}

func (s *Stack) Add() error {
	if s.top <=1 {
		return panicCode(0x02)
	}
	s.top--; v := s.data[s.top]; s.data[s.top] = 0
	s.top--; s.data[s.top] += v; s.top++
	return nil
}

func (s *Stack) Sub() error {
	if s.top <=1 {
		return panicCode(0x02)
	}
	s.top--; v := s.data[s.top]; s.data[s.top] = 0
	s.top--; s.data[s.top] -= v; s.top++
	return nil
}

func (s *Stack) Mul() error {
	if s.top <=1 {
		return panicCode(0x02)
	}
	s.top--; v := s.data[s.top]; s.data[s.top] = 0
	s.top--; s.data[s.top] *= v; s.top++
	return nil
}

func (s *Stack) Div() error {
	if s.top <=1 {
		return panicCode(0x02)
	}
	s.top--; v := s.data[s.top]; s.data[s.top] = 0
	s.top--; s.data[s.top] /= v; s.top++
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

func NewStack() *Stack {
	return &Stack{
		top: 0,
	}
}
