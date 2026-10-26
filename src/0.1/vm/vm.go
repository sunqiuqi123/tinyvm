package vm

import (
	"fmt"
)

type Stack struct {
	data [256]uint64
	top int
}

func panic(code byte) error {
	switch code {
	case 0x00:
		return fmt.Errorf("Panic 0x00:Ok\n")
	case 0x01:
		return fmt.Errorf("Panic 0x01: Stack overflow.\n")
	case 0x02:
		return fmt.Errorf("Panic 0x02: Stack downflow.\n")
	default:
		return fmt.Errorf("Unknown panic.")
	}
}

func (s *Stack) Push(v uint64) error {
	if s.top >= 256 {
		return panic(0x01)
	}
	s.data[s.top] = v
	s.top++
	return nil
}

func (s *Stack) Pop() (uint64,error) {
	if s.top <= 0 {
		return 0,panic(0x02)
	}
	s.top--
	v := s.data[s.top]
	return v,nil
}

func (s *Stack) Size() int {
	return s.top
}

func (s *Stack) Look() uint64 {
	return s.data[s.top - 1]
}

func NewStack() *Stack {
	return &Stack{
		top: 0,
	}
}

