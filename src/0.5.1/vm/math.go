package vm

func ByteCombine(a byte, b byte) int { //预留
	return int(a)<<8 | int(b)
}

func PowTwo(exp int) int { //预留
	result := 1; base := 2
    for exp > 0 {
        if exp&1 == 1 {
            result *= base
        }
        base *= base
        exp >>= 1
    }
    return result
}

