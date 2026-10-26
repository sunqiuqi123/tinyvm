package vm

import (
	"fmt"
)

func panicCode(code byte) error {
	switch code {
	case 0x00:
		return fmt.Errorf("Panic 0x00:Ok\n") //屁事没有
	case 0x01:
		return fmt.Errorf("Panic 0x01: Stack overflow.\n") //栈上溢
	case 0x02:
		return fmt.Errorf("Panic 0x02: Stack downflow.\n") //栈下溢
	case 0x03:
		return fmt.Errorf("Panic 0x03: Option counter overflow.\n") //读程序的针走丢了
	case 0x04:
		return fmt.Errorf("Panic 0x04: Not a valid memory address.\n") //没有那个内存地址你读滚木呢
	case 0x05:
		return fmt.Errorf("Panic 0x05: Not a valid opcode.\n") //胡言乱语这一块
	case 0x06:
		return fmt.Errorf("Panic 0x06: Divide-by-zero.\n") //除以零了
	default:
		return fmt.Errorf("Unknown panic.") //什么错误代码，闻所未闻
	}
}
