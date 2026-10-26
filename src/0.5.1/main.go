package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tinyvm/vm"
)

var opcodeMap = map[string]byte{
	"PUSH":   vm.PUSH,
	"POP":    vm.POP,
	"ADD":    vm.ADD,
	"SUB":    vm.SUB,
	"MUL":    vm.MUL,
	"DIV":    vm.DIV,
	"LOAD":   vm.LOAD,
	"STORE":  vm.STORE,
	"STORES": vm.STORES,
	"LOOK":   vm.LOOK,
	"LOOC":   vm.LOOC,
	"PUT":    vm.PUT,
	"PUTC":   vm.PUTC,
	"ENDL":   vm.ENDL,
	"LISTEN": vm.LISTEN,
	"JMP":    vm.JMP,
	"JZ":     vm.JZ,
	"JNZ":    vm.JNZ,
	"JEQ":    vm.JEQ,
	"JNEQ":   vm.JNEQ,
	"JGT":    vm.JGT,
	"JLT":    vm.JLT,
	"JGE":    vm.JGE,
	"JLE":    vm.JLE,
	"JPOS":   vm.JPOS,
	"JNEG":   vm.JNEG,
	"CALL":   vm.CALL,
	"RET":    vm.RET,
	"HALT":   vm.HALT,
}

// 每个操作码需要的操作数数量
var operandCount = map[string]int{
	"PUSH":  1, //PUSH <数据>
	"POP":   0, //POP
	"ADD":   0, //ADD
	"SUB":   0, //SUB
	"MUL":   0, //MUL
	"DIV":   0, //DIV
	"LOAD":  1, //LOAD <内存地址>
	"STORE": 2, //STORE <内存地址> <数据>
	"STORES":1, //STORES <内存地址>
	"LOOK":  0, //LOOK
	"LOOC":  0, //LOOC
	"PUT":   0, //PUT
	"PUTC":  0, //PUTC
	"ENDL":  0, //ENDL
	"LISTEN":0, //LISTEN
	"JMP":   1, //JMP <程序计数器指针>
	"JZ":    1, //JZ <程序计数器指针>
	"JNZ":   1, //JNZ <程序计数器指针>
	"JEQ":   1, //JEQ <程序计数器指针>
	"JNEQ":  1, //JNEQ <程序计数器指针>
	"JGT":   1, //JGT <程序计数器指针>
	"JLT":   1, //JLT <程序计数器指针>
	"JGE":   1, //JGE <程序计数器指针>
	"JLE":   1, //JLE <程序计数器指针>
	"JPOS":  1, //JPOS <程序计数器指针> //预留
	"JNEG":  1, //JNEG <程序计数器指针> //预留
	"CALL":  1, //CALL <程序计数器指针>
	"RET":   0, //RET
	"HALT":  0, //HALT
}

var reverseOpcodeMap map[byte]string

func init() {
	reverseOpcodeMap = make(map[byte]string)
	for name, code := range opcodeMap {
		reverseOpcodeMap[code] = name
	}
}

func printUsage() {
	fmt.Println(`TinyVM - 简易栈式虚拟机 v0.5.1
用法:
    -f <文件名>        读入程序文件（必填）
    -r                 执行程序
    -b <文件名>        将程序转为二进制机器码并输出到文件
    -p                 打印程序字节码（反汇编）
    -d                 启用调试模式
`)
}

func isOpcode(b byte) bool {
	_, exists := reverseOpcodeMap[b]
	return exists
}

func disassemble(program []byte) {
	if len(program) == 0 {
		fmt.Println("程序为空")
		return
	}
	fmt.Printf("\n程序字节码 (%d 字节):\n", len(program))
	fmt.Print("  HEX: ")
	for i, b := range program {
		if i > 0 && i%16 == 0 {
			fmt.Println()
			fmt.Print("       ")
		}
		fmt.Printf("%02X ", b)
	}
	fmt.Println()
	fmt.Print("  DEC: ")
	for i, b := range program {
		if i > 0 && i%16 == 0 {
			fmt.Println()
			fmt.Print("       ")
		}
		fmt.Printf("%3d ", b)
	}
	fmt.Println()
	fmt.Println("\n反汇编:")
	i := 0
	for i < len(program) {
		op := program[i]
		i++
		opName, exists := reverseOpcodeMap[op]
		if !exists {
			fmt.Printf("  DATA %d\n", op)
			continue
		}
		count := operandCount[opName]
		var operands []string
		for j := 0; j < count && i < len(program); j++ {
			operands = append(operands, fmt.Sprintf("%d", program[i]))
			i++
		}
		if len(operands) == 0 {
			fmt.Printf("  %s\n", opName)
		} else {
			fmt.Printf("  %s %s\n", opName, strings.Join(operands, " "))
		}
	}
	fmt.Println()
}

func printProgramDebug(program []byte) {
	fmt.Printf("\n程序详细信息:\n")
	fmt.Printf("  总字节数: %d\n", len(program))
	count := 0
	for _, b := range program {
		if isOpcode(b) {
			count++
		}
	}
	fmt.Printf("  操作码数量: %d\n", count)
	fmt.Printf("  数据字节数: %d\n", len(program)-count)
}

// ------------------------------------------------------------
// 标签支持（两遍扫描 + 向前引用 + 大小写不敏感）
// ------------------------------------------------------------

func parseLine(line string) []string {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
		return nil
	}
	if idx := strings.Index(line, ";"); idx >= 0 {
		line = line[:idx]
	}
	if idx := strings.Index(line, "#"); idx >= 0 {
		line = line[:idx]
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	return strings.Fields(line)
}

// 第一遍扫描：先记录所有标签名称，再计算地址并回填
func buildLabelMap(tokensLines [][]string) map[string]int {
	labelMap := make(map[string]int)

	// 阶段1：先记录所有标签名称（占位）
	for _, tokens := range tokensLines {
		if len(tokens) == 0 {
			continue
		}
		first := strings.ToUpper(tokens[0])
		if strings.HasSuffix(first, ":") {
			labelName := strings.TrimSuffix(first, ":")
			labelMap[labelName] = 0 // 暂存占位
		}
	}

	// 阶段2：计算地址，回填标签
	address := 0
	for _, tokens := range tokensLines {
		if len(tokens) == 0 {
			continue
		}
		first := strings.ToUpper(tokens[0])

		if strings.HasSuffix(first, ":") {
			labelName := strings.TrimSuffix(first, ":")
			labelMap[labelName] = address
			continue
		}

		if _, exists := opcodeMap[first]; exists {
			address++ // 操作码占1字节
			count := operandCount[first]
			if count > len(tokens)-1 {
				count = len(tokens) - 1
			}
			address += count
		} else {
			if _, err := strconv.Atoi(first); err == nil {
				address++
			}
		}
	}
	return labelMap
}

func generateBytecode(tokensLines [][]string, labelMap map[string]int) []byte {
	var program []byte

	for _, tokens := range tokensLines {
		if len(tokens) == 0 {
			continue
		}
		first := strings.ToUpper(tokens[0])

		if strings.HasSuffix(first, ":") {
			continue
		}

		if opcode, exists := opcodeMap[first]; exists {
			program = append(program, opcode)
			i := 1
			count := operandCount[first]

			for j := 0; j < count && i < len(tokens); j++ {
				operand := tokens[i]
				i++

				if num, err := strconv.Atoi(operand); err == nil {
					if num < 0 || num > 255 {
						fmt.Printf("警告: 操作数 %d 超出字节范围 (0-255)\n", num)
						continue
					}
					program = append(program, byte(num))
				} else {
					// 标签名大小写不敏感：转大写再查
					if addr, ok := labelMap[strings.ToUpper(operand)]; ok {
						if addr < 0 || addr > 255 {
							fmt.Printf("警告: 标签 %s 地址 %d 超出字节范围 (0-255)\n", operand, addr)
							continue
						}
						program = append(program, byte(addr))
					} else {
						fmt.Printf("警告: 未知标签或无效数字: %s\n", operand)
						program = append(program, 0)
					}
				}
			}
		} else {
			if num, err := strconv.Atoi(first); err == nil {
				if num < 0 || num > 255 {
					fmt.Printf("警告: 数字 %d 超出字节范围 (0-255)\n", num)
				} else {
					program = append(program, byte(num))
				}
			} else {
				fmt.Printf("警告: 未知 token: %s\n", first)
			}
		}
	}
	return program
}

func assembleFile(filename string, debug bool) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	var tokensLines [][]string
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if debug {
			fmt.Printf("[DEBUG] 读取第 %d 行: %s\n", lineNum, line)
		}
		tokens := parseLine(line)
		if tokens != nil {
			tokensLines = append(tokensLines, tokens)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件错误: %v", err)
	}

	labelMap := buildLabelMap(tokensLines)
	if debug {
		fmt.Printf("[DEBUG] 标签映射: %v\n", labelMap)
	}

	program := generateBytecode(tokensLines, labelMap)
	return program, nil
}

// ------------------------------------------------------------
// main
// ------------------------------------------------------------

func main() {
	var filename string
	var run bool
	var binFile string
	var print bool
	var debug bool

	flag.StringVar(&filename, "f", "", "程序文件名")
	flag.BoolVar(&run, "r", false, "执行程序")
	flag.StringVar(&binFile, "b", "", "输出二进制名")
	flag.BoolVar(&print, "p", false, "打印程序字节码")
	flag.BoolVar(&debug, "d", false, "启用调试模式")
	flag.Parse()

	if filename == "" {
		fmt.Println("错误: 请使用 -f 指定程序文件")
		printUsage()
		return
	}

	if debug {
		fmt.Printf("[DEBUG] 打开文件: %s\n", filename)
	}

	program, err := assembleFile(filename, debug)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	if debug {
		printProgramDebug(program)
	}

	if print {
		disassemble(program)
	}

	if binFile != "" {
		if len(program) == 0 {
			fmt.Println("警告: 程序为空，不生成二进制文件")
		} else {
			err = os.WriteFile(binFile, program, 0644)
			if err != nil {
				fmt.Printf("写入二进制文件失败: %v\n", err)
				return
			}
			fmt.Printf("二进制文件已写入: %s (%d 字节)\n", binFile, len(program))
		}
	}

	if run {
		if len(program) == 0 {
			fmt.Println("错误: 程序为空，无法运行")
			return
		}
		fmt.Printf("运行程序 (共 %d 字节)...\n", len(program))
		machine := vm.NewVM()
		if debug {
			fmt.Println("[DEBUG] 开始执行程序")
		}
		err = machine.Run(program)
		if err != nil {
			fmt.Printf("运行错误: %v\n", err)
			return
		}
		fmt.Println("程序执行完成")
	}

	if !run && binFile == "" && !print {
		fmt.Printf("程序已加载: %d 字节\n", len(program))
		fmt.Println("提示: 使用 -r 运行程序，-b 输出二进制文件，-p 打印程序")
	}
}
