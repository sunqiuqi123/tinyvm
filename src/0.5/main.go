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

// 操作码映射
var opcodeMap = map[string]byte{
	"PUSH":  vm.PUSH,
	"POP":   vm.POP,
	"ADD":   vm.ADD,
	"SUB":   vm.SUB,
	"MUL":   vm.MUL,
	"DIV":   vm.DIV,
	"LOAD":  vm.LOAD,
	"STORE": vm.STORE,
	"STORES":vm.STORES,
	"LOOK":  vm.LOOK,
	"LOOC":  vm.LOOC,
	"PUT":   vm.PUT,
	"PUTC":  vm.PUTC,
	"ENDL":  vm.ENDL,
	"LISTEN":vm.LISTEN,
	"JMP":   vm.JMP,
	"JZ":    vm.JZ,
	"JNZ":   vm.JNZ,
	"JEQ":   vm.JEQ,
	"JNEQ":  vm.JNEQ,
	"JGT":   vm.JGT,
	"JLT":   vm.JLT,
	"JGE":   vm.JGE,
	"JLE":   vm.JLE,
	"JPOS":  vm.JPOS,
	"JNEG":  vm.JNEG,
	"RET":   vm.RET,
	"HALT":  vm.HALT,
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
	"RET":   0, //RET
	"HALT":  0, //HALT
}

// 反向映射：字节码 -> 操作码名称
var reverseOpcodeMap map[byte]string

func init() {
	reverseOpcodeMap = make(map[byte]string)
	for name, code := range opcodeMap {
		reverseOpcodeMap[code] = name
	}
}

func printUsage() {
	fmt.Println(`TinyVM - 简易栈式虚拟机 v0.5
用法:
    -f <文件名>        读入程序文件（必填）
    -r                 执行程序
    -b <文件名>        将程序转为二进制机器码并输出到文件
    -p                 打印程序字节码（反汇编）
    -d                 启用调试模式

示例:
    tinyvm -f program.tvm -p          # 加载并打印程序
    tinyvm -f program.tvm -r          # 加载并运行程序
    tinyvm -f program.tvm -b out.bin  # 生成二进制文件
    tinyvm -f program.tvm -r -p -d    # 调试模式运行并打印
`)
}

// 判断一个字节是否是操作码
func isOpcode(b byte) bool {
	_, exists := reverseOpcodeMap[b]
	return exists
}

// 处理一行代码（同步版本）
func processLine(line string, program *[]byte) {
	// 去除首尾空格
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// 处理注释（支持 ; 和 #）
	if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
		return
	}

	// 移除行内注释
	if idx := strings.Index(line, ";"); idx >= 0 {
		line = line[:idx]
	}
	if idx := strings.Index(line, "#"); idx >= 0 {
		line = line[:idx]
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// 按空格分割 token
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return
	}

	i := 0
	for i < len(tokens) {
		token := strings.ToUpper(tokens[i])

		// 检查是否是操作码
		opcode, exists := opcodeMap[token]
		if exists {
			// 添加操作码
			*program = append(*program, opcode)
			i++

			// 获取需要的操作数数量
			count := operandCount[token]

			// 处理操作数
			for j := 0; j < count && i < len(tokens); j++ {
				operand := tokens[i]

				// 尝试解析为数字
				if num, err := strconv.Atoi(operand); err == nil {
					if num < 0 || num > 255 {
						fmt.Printf("警告: 操作数 %d 超出字节范围 (0-255)\n", num)
						i++
						continue
					}
					*program = append(*program, byte(num))
				} else {
					fmt.Printf("警告: 无效的操作数 %s (需要数字)\n", operand)
					i++
				}
				i++
			}
		} else {
			// 不是操作码，尝试作为数字处理
			if num, err := strconv.Atoi(tokens[i]); err == nil {
				if num < 0 || num > 255 {
					fmt.Printf("警告: 数字 %d 超出字节范围 (0-255)\n", num)
				} else {
					*program = append(*program, byte(num))
				}
			} else {
				fmt.Printf("警告: 未知 token: %s\n", tokens[i])
			}
			i++
		}
	}
}

// 反汇编程序
func disassemble(program []byte) {
	if len(program) == 0 {
		fmt.Println("程序为空")
		return
	}

	fmt.Printf("\n程序字节码 (%d 字节):\n", len(program))

	// 十六进制显示
	fmt.Print("  HEX: ")
	for i, b := range program {
		if i > 0 && i%16 == 0 {
			fmt.Println()
			fmt.Print("       ")
		}
		fmt.Printf("%02X ", b)
	}
	fmt.Println()

	// 十进制显示
	fmt.Print("  DEC: ")
	for i, b := range program {
		if i > 0 && i%16 == 0 {
			fmt.Println()
			fmt.Print("       ")
		}
		fmt.Printf("%3d ", b)
	}
	fmt.Println()

	// 反汇编显示
	fmt.Println("\n反汇编:")
	i := 0
	for i < len(program) {
		op := program[i]
		i++

		// 查找操作码名称
		opName, exists := reverseOpcodeMap[op]
		if !exists {
			fmt.Printf("  DATA %d\n", op)
			continue
		}

		// 获取操作数数量
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

// 打印程序的详细信息（调试用）
func printProgramDebug(program []byte) {
	fmt.Printf("\n程序详细信息:\n")
	fmt.Printf("  总字节数: %d\n", len(program))
	fmt.Printf("  操作码数量: ")
	count := 0
	for _, b := range program {
		if isOpcode(b) {
			count++
		}
	}
	fmt.Printf("%d\n", count)
	fmt.Printf("  数据字节数: %d\n", len(program)-count)
}

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

	// 检查是否提供了文件名
	if filename == "" {
		fmt.Println("错误: 请使用 -f 指定程序文件")
		printUsage()
		return
	}

	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	if debug {
		fmt.Printf("[DEBUG] 打开文件: %s\n", filename)
	}

	var program []byte

	scanner := bufio.NewScanner(file)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if debug {
			fmt.Printf("[DEBUG] 读取第 %d 行: %s\n", lineNum, line)
		}

		// ✅ 顺序处理，保证顺序
		processLine(line, &program)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件错误: %v\n", err)
		return
	}

	if debug {
		fmt.Printf("[DEBUG] 程序加载完成，共 %d 字节\n", len(program))
		printProgramDebug(program)
	}

	// 打印程序
	if print {
		disassemble(program)
	}

	// 输出二进制文件
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

	// 运行程序
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

	// 如果没有任何操作，显示提示
	if !run && binFile == "" && !print {
		fmt.Printf("程序已加载: %d 字节\n", len(program))
		fmt.Println("提示: 使用 -r 运行程序，-b 输出二进制文件，-p 打印程序")
	}
}
