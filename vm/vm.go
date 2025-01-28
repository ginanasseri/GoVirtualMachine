// Package vm implements a Go Virtual Machine which provies an isolated 
// execution environment for executing programs written for a Susan; a
// hypothetical ISA with 10 32-bit CPU registers which is capable of
// executing basic instructions including read, write, add, jump, and 
// printing operations. 
//
// The Virtual Machine acts as the control transfer point between 
// the host OS and Susan processes, complete with its own Virtual 
// Memory data structures. 
package vm

import (
    "fmt"
    "os"
    "bufio"
    "gvm/parser"
    "gvm/instructions"
    "gvm/interpreter"
)

const ( 
    NUM_REGISTERS = 10 // constant. Does not change
    INIT  = 10 // initial codeblock size
)

// VirtualMemory defines the memory architecture of the virtual machine. It
// contains the registers, the executatble code block, and the size of the
// code block so that the last address in the code block is immediately
// accessable
type VirtualMemory struct {
    // Registers is a slice of int32 values which serve as a one-to-one mapping
    // between the virtual memory and Susan's 32-bit CPU registers. Registers
    // are used by the interpreter when executing instructions. Registers 
    // 1:9 are read and write registers. Register 0 is a special purpose 
    // register which holds the last address within the code block. This 
    // register is read-only when control is transferred to the interpreter. 
    Registers []int32

    // Code is a slice of Instruction structs which represent the executable
    // bytecode instructions of the program loaded into the virtual machine.
    Code []instructions.Instruction

    // CodeSize in an integer value which holds the number of executable 
    // instructions within the Code slice. It serves as an immediate 
    // access point to determine the last address within the executable 
    // code block to prevent writing to or branching to restricted 
    // or invalid memory addresses. 
    CodeSize int
}

// NewVirtualMemory initializes a new VirtualMemory instance to be used to 
// represent Susan's memory image within the virtual machine. It initializes
// the registers, the code block Code with an initial size of INIT_CB_SIZE,
// and initializes CodeSize as 0 indicating that no instructions have been
// written yet. The virtual memory provides an isolated environment for 
// the virtual machine to load and execute programs with.

func NewVirtualMemory() *VirtualMemory {
    return &VirtualMemory{
        Registers: make([]int32, NUM_REGISTERS),
        Code: make([]instructions.Instruction, INIT), 
        CodeSize: 0,
    }
}



// VirtualMachine represents the virtual machine main architecture. 
// It is the control point between each stage of the simulation;
// loading the program code, writing the executable to the virtual
// memory code section, passing control to and from the interpreter,
// and handling any errors returned at any point during the 
// execution process. 
type VirtualMachine struct {
    // VMem points to a VirtualMemory instance as described above. 
    // During execution, once the program code is loaded, as each
    // instruction is parsed into bytecode, it is written into the 
    // Code block and the CodeSize is updated. After the parsing 
    // is completed, the final CodeSize is written into Register[0].
    VMem *VirtualMemory

    // The interpreter is a pointer to an Interpreter instance as
    // defined in the 'interpreter' package. Upon initialization,
    // the interpreter gets a reference to the virtual memory. 
    // Once parsing is completed, the interpreter decodes each
    // instruction, extracting the opcode, and then executes 
    // the instructions by dispatching to the routine indicated by 
    // each instructions opcode. 
    Interpreter *interpreter.Interpreter
}

// NewVirtualMachine initializes a new VirtualMachine instance. It 
// sets up the guest memory image for the source ISA and code block
// for the source executable. Once the VirtualMemory is allocated,
// an interpreter instance is initialized with a reference to the 
// virtual machine's registers and code block. This allows the
// interpreter to execute instructions in an isolated environment
// with control is returned back to the VM 
// 
// Initialize a virtual machine with pre-allocated virtual memory and 
// an interpreter with a reference to the virtual memory. 
func NewVirtualMachine() *VirtualMachine {
    vMem := NewVirtualMemory()
    return &VirtualMachine{
        VMem: vMem,
        Interpreter: interpreter.New(vMem.Registers, vMem.Code),
    }
}

// ParseInstructions reads each line in the source code file, 
// where each line is an instruction within Susan's instruction 
// set. Each instruction is tokenized, parsed into represent-
// ative bytecode instructions, and written into the virtual
// memory code block section. 


func (vm *VirtualMachine) ParseInstructions(sourceCode *os.File) error {
    scanner := bufio.NewScanner(sourceCode)
    for scanner.Scan() {
        sourceInstruction := scanner.Text()
        parser, err := parser.New(sourceInstruction)
        if err != nil {
            return err
        }
        // get bytecode instruction from source instruction
        byteCodeInstr, err := parser.Instruction() 
        if err != nil {
            return err
        } 
        // write bytecode instruction to virtual memory 
        vm.VMem.Code[vm.VMem.CodeSize] = byteCodeInstr
        vm.VMem.CodeSize += 1
    }
    return nil
}

// Execute is the control transfer center of the virtual 
// machine. Within execute, the source program is loaded 
// via a call to the host OS. Once loaded, the program 
// code is parsed into bytecode instructions which are 
// written into the virtual memory code block. The file 
// is then closed and control is transferred to the
// interpreter which executes each instruction using a 
// basic decode and dispatch routine to execute each
// Susan instruction by performing register write and 
// read operations, internal state management via its 
// program counter, and system calls to the host OS 
// for printing to the screen. 
//
// Any errors which occur are propagated from the source 
// and returned and handled here.


func (vm *VirtualMachine) Execute(file string) error {    

    // Load program code
    sourceCode, err := os.Open(file)
    if err != nil {
        return fmt.Errorf("gvm: vm.Execute: failed to open file: '%s'",file)
    }
    defer sourceCode.Close()

    // Parse source instructions as bytecode into the virtual
    // memory code block 
    if err := vm.ParseInstructions(sourceCode); err != nil {
        return err
    }
    // Write last address of code block to register 0
    vm.VMem.Registers[0] = int32(vm.VMem.CodeSize)

    // Invoke interpreter to execute program
    if err := vm.Interpreter.Interpret(); err != nil {
        return err
    }
    return nil
}
