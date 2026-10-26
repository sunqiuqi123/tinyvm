package vm

type VM struct {
	stack *Stack
	memory *Memory
	pc int
	running bool
}

const (
	//0x01-0x09: 栈操作
	PUSH = 0x01
	POP = 0x02
	ADD = 0x03
	SUB = 0x04
	MUL = 0x05
	DIV = 0x06
	LOAD = 0x07
	STORE = 0x08
	STORES = 0x09
	//0x0A-0x0F: IO操作
	LOOK = 0x0A
	LOOC = 0x0B
	PUT = 0x0C
	PUTC = 0x0D
	ENDL = 0x0E
	LISTEN = 0x0F
	//0x10-0x1F: 跳转操作
	JMP = 0x10
	JZ = 0x11
	JNZ = 0x12
	JEQ = 0x13
	JNEQ = 0x14
	JGT = 0x15
	JLT = 0x16
	JGE = 0x17
	JLE = 0x18
	JPOS = 0x19
	JNEG = 0x1A
	CALL = 0x1E
	RET = 0x1F
	//0x20-0x2F: 内存操作
	//停机
	HALT = 0xFF
)

func NewVM() *VM {
	return &VM{
		stack: NewStack(),
		memory: NewMemory(),
		pc: 0,
		running: false, //这里我改成了默认关机，和halt配合
	}
}

func (vm *VM) Run(program []byte) error {
	vm.pc = 0; vm.running = true //开机
	for vm.running && vm.pc < len(program) { //这里pc是从0开始，len是从1开始，所以是<而不是<=
		op := program[vm.pc]; vm.pc++ //和栈差不多的指针逻辑
		switch op { //开始读入程序
		case PUSH:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value := program[vm.pc]; vm.pc++ //往下走一步找数据
			vm.stack.Push(value)
		case POP:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			vm.stack.Pop()
		case ADD:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			vm.stack.Add()
		case SUB:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			vm.stack.Sub()
		case MUL:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			vm.stack.Mul()
		case DIV:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			vm.stack.Div()
		case LOAD:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			addr := int(program[vm.pc]); vm.pc++
			value, _ := vm.memory.LoadByte(addr)
			vm.stack.Push(value)
		case STORE:
			if vm.pc + 1  > len(program) {
				return panicCode(0x03)
			}
			addr := int(program[vm.pc]); vm.pc++
			value := byte(program[vm.pc]); vm.pc++
			vm.memory.StoreByte(addr,value)
		case STORES:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			addr := int(program[vm.pc]); vm.pc++
			value, _ := vm.stack.Pop()
			vm.memory.StoreByte(addr,value)
		case LOOK:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			vm.stack.Look()
		case LOOC:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			vm.stack.Looc()
		case PUT:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			vm.stack.Put()
		case PUTC:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			vm.stack.Putc()
		case ENDL:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			vm.stack.Endl()
		case LISTEN:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value, _ := vm.stack.Listen()
			vm.stack.Push(value)
		case JMP:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			addr := int(program[vm.pc])
			vm.pc = addr
		case JZ:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value, _ := vm.stack.Pop()
			if value == 0 {
				addr := int(program[vm.pc])
				vm.pc = addr
			}else{
				vm.pc++
			}
		case JNZ:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value, _ := vm.stack.Pop()
			if value != 0 {
				addr := int(program[vm.pc])
				vm.pc = addr
			}else{
				vm.pc++
			}
		case JEQ:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value1, _ := vm.stack.Pop()
			value2, _ := vm.stack.Pop()
			if value2 == value1 {
				addr := int(program[vm.pc])
				vm.pc = addr
			}else{
				vm.pc++
			}
		case JNEQ:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value1, _ := vm.stack.Pop()
			value2, _ := vm.stack.Pop()
			if value2 != value1 {
				addr := int(program[vm.pc])
				vm.pc = addr
			}else{
				vm.pc++
			}
		case JGT:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value1, _ := vm.stack.Pop()
			value2, _ := vm.stack.Pop()
			if value2 > value1 {
				addr := int(program[vm.pc])
				vm.pc = addr
			}else{
				vm.pc++
			}
		case JLT:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value1, _ := vm.stack.Pop()
			value2, _ := vm.stack.Pop()
			if value2 < value1 {
				addr := int(program[vm.pc])
				vm.pc = addr
			}else{
				vm.pc++
			}
		case JGE:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value1, _ := vm.stack.Pop()
			value2, _ := vm.stack.Pop()
			if value2 >= value1 {
				addr := int(program[vm.pc])
				vm.pc = addr
			}else{
				vm.pc++
			}
		case JLE:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value1, _ := vm.stack.Pop()
			value2, _ := vm.stack.Pop()
			if value2 <= value1 {
				addr := int(program[vm.pc])
				vm.pc = addr
			}else{
				vm.pc++
			}
		case JPOS:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value, _ := vm.stack.Pop()
			if value > 0 {
				addr := int(program[vm.pc])
				vm.pc = addr
			}else{
				vm.pc++
			}
		case JNEG:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			value, _ := vm.stack.Pop()
			if value < 0 {
				addr := int(program[vm.pc])
				vm.pc = addr
			}else{
				vm.pc++
			}
		case CALL:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			addr := int(program[vm.pc]); vm.pc++
			vm.stack.Push(byte(vm.pc))
			vm.pc = int(addr)
		case RET:
			if vm.pc > len(program) {
				return panicCode(0x03)
			}
			addr, err := vm.stack.Pop()
			if err != nil {
				return err
			}
			vm.pc = int(addr)
		case HALT:
			vm.running = false; //关机
		default:
			return panicCode(0x05)
		}
	}
	return nil
}
