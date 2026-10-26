package vm

import (
	"encoding/binary"
)

type Memory struct {
	data [1024]byte
}

func (m *Memory) LoadByte(addr int) (byte,error) {
	if (addr >= 1024) || (addr < 0) {
		return 0,panicCode(0x04)
	}
	return m.data[addr],nil
}

func (m *Memory) StoreByte(addr int, value byte) error {
	if (addr >= 1024) || (addr < 0) {
		return panicCode(0x04)
	}
	m.data[addr] = value
	return nil
}

func (m *Memory) LoadUint64(addr int) (uint64,error) {
	if (addr +8 > 1024) || (addr < 0) {
		return 0,panicCode(0x04)
	}
	return binary.BigEndian.Uint64(m.data[addr:addr+8]),nil
}

func (m *Memory) StoreUint64(addr int, value uint64) error {
	if (addr +8 > 1024) || (addr < 0) {
		return panicCode(0x04)
	}
	binary.BigEndian.PutUint64(m.data[addr:addr+8], value)
	return nil
}

func NewMemory() *Memory {
	return &Memory{
	}
}

