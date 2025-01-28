// Package interpreter executes the bytecode instructions contained 
// in the executable code block in the VirtualMachine's Virtual
// Memory as defined in the 'vm' package. The interpreter executes
// instructions one at a time using the decode and dispatch method.
package interpreter

import (
    "fmt"
    "os"
    "time"
    "strings"
    "gvm/instructions"
    "github.com/fatih/color"
)

// Defining OPCODES as constants  
const (
    OPCODE_STDOUT = 0x00
    OPCODE_LDI    = 0x01
    OPCODE_JUMP   = 0x02
    OPCODE_ADD    = 0x17
    OPCODE_ADDV   = 0x18
    OPCODE_DRAW   = 0x19
    OPCODE_BLINK  = 0x20
    OPCODE_PRINTR = 0x21
)

type Interpreter struct {
    PC int32
    Registers []int32
    Code []instructions.Instruction
}

// Initialze an interpreter with pre-allocated virtual memory provided by 
// the 'vm': this includes 10 registers and a codeblock containing the
// bytecode representation of the source program's asm instructions. 
func New(vregisters []int32, code []instructions.Instruction) *Interpreter {
    return &Interpreter{
        PC: 0,
        Registers: vregisters, 
        Code: code,
    }
}

// Interpret reads each bytecode instruction in the virtual memory codeblock 
// and calls the decode and dispatch routine for each instruction. After 
// the source code is parsed as bytecode into the virtual memory codeblock, 
// the VM writes the size of the codeblock into register 0 so the interpreter
// knows the start and end addresses of the address space where the bytecode
// is stored, i.e., where the interpreter has permission to access. 
func (interp *Interpreter) Interpret() error {
    lastAddr,_ := interp.ReadFrom(0)
    for interp.PC < lastAddr {
        if err := interp.DecodeAndDispatch(interp.Code[interp.PC]); err != nil {
            return err
        }
        interp.PC++
    }
    return nil 
}


// DecodeAndDispatch reads a bytecode instruction to obtain the OpCode for the 
// current instruction and, if a valid opcode is obtained, the interpreter 
// gets any additional information needed from the bytecode, depending on the 
// type of instruction, and then dispatches to the appropriate routine to 
// execute the instruction. 
func (interp *Interpreter) DecodeAndDispatch(instr instructions.Instruction) error {
    switch instr.GetOpCode() {
    
     // LDI
    case OPCODE_LDI:
        if err := interp.LoadImmediate(instr.GetArg1(),instr.GetArg2()); err != nil {
            return err 
        }
        return nil

    // JUMP
    case OPCODE_JUMP:
        if err := interp.JumpTo(instr.GetArg1()); err != nil {
            return err
        }
        return nil
    
    // ADD
    case OPCODE_ADD:
        if err := interp.Add(instr.GetArg1(),instr.GetArg2()); err != nil {
            return err
        }
        return nil

    // ADDV 
    case OPCODE_ADDV:
        if err := interp.AddV(instr.GetArg1(),instr.GetArg2()); err != nil {
            return err
        }
        return nil

    // DRAW
    case OPCODE_DRAW:
        if err := interp.Draw(instr.GetArg1()); err != nil {
            return err
        }
        return nil

    // BLINK
    case OPCODE_BLINK:
        if err := interp.Blink(instr.GetArg1()); err != nil {
            return err
        }
        return nil
        
    // STDOUT
    case OPCODE_STDOUT:
        if err := interp.PrintToStdOut(instr.GetArg1()); err != nil {
            return err
        }
        return nil 

    // PRINTR
    case OPCODE_PRINTR:
        if err := interp.PrintRegisters(); err != nil {
            return err
        }
        return nil
     
    // invalid opcode     
    default:
        return fmt.Errorf("interpreter/interpreter.go: invalid command %v",instr)
    }
}

// WriteTo writes an int32 value to a regiseter. Where register holds the index of the 
// register being written to. If the register is r0, then a permission denied error 
// is returned as r0 is read only. If the register index is > 9, then an 
// invalid register index error is returned as the ISA being emulated contains 
// registers R0:R9. 
func (interp *Interpreter) WriteTo(register, value int32) error {
    if register == 0 {
        return fmt.Errorf("gvm: write to R0: permission denied [R0 is read-only]")
    }
    if register > 9 {
        return fmt.Errorf("gvm: invalid register: R%d [use registers R0:R9]",register)
    }
    interp.Registers[register] = value
        return nil
}

// ReadFrom returns the value from the indicated register. If the register index is out 
// of range then an error is returned. 
func (interp *Interpreter) ReadFrom(register int32) (int32, error) {
    if register > 9 {
        return 0, fmt.Errorf("gvm: invalid register: R%d [use registers R0:R9]",register)
    }
    return interp.Registers[register], nil
}

// CheckJump validates the requested jump address provided as an argument to a 
// JUMP instruction. 
//
// If the address provided is less than the current address, then an infinite loop
// warning error is returned. 
//
// If the address provided is outside of the virtual memory codeblock, where the last 
// address of the codeblock is provided in register 0, then a segmentation violation
// error is returned as the JUMP address is not a valid memory address. 
func (interp *Interpreter) CheckJump(addr int32) error {
    lastAddr,_ := interp.ReadFrom(0)
    if addr <= interp.PC {
        return fmt.Errorf("gvm: JUMP at addr %d to %d: infinite loop warning.",interp.PC,addr)
    }
    if addr > lastAddr - 1 {
        return fmt.Errorf("gvm: JUMP addr invalid: segmentation violation.")
    }
    return nil
}


// ****** Instruction function list starts here ***** 

// LDI routine: LoadImmediate
// LoadImmediate writes an int32 literal value into the register at index 
// regiseter. 
func (interp *Interpreter) LoadImmediate(register, value int32) error {
    if err := interp.WriteTo(register,value); err != nil {
        return err
    }
    return nil
}

// JUMP routine: JumpTo
// JumpTo validates the jump address with the CheckJump function and,
// if no error is returned, then the PC is updated to the jump to 
// address - 1 (as the PC is incremented when returned). The address
// is converted to an integer 
func (interp *Interpreter) JumpTo(address int32) error {
    if err := interp.CheckJump(address); err != nil {
        return err
    }
    interp.PC = address-1
    return nil
}

// ADD routine: Add
// Add reads the values in registers ri and rj, adds them, and
// writes the result to register ri. 
func (interp *Interpreter) Add(ri, rj int32) error {
    value1, err := interp.ReadFrom(ri)
    if err != nil {
        return err
    }
    value2, err := interp.ReadFrom(rj)
    if err != nil {
        return err
    }
    result := value1 + value2
    interp.WriteTo(ri,result)
    return nil
}

// STDOUT routine: PrintToStdOut
// Prints the value at register to terminal 
func (interp *Interpreter) PrintToStdOut(register int32) error {
    value, err := interp.ReadFrom(register)
    if err != nil {
        return err 
    }
    fmt.Printf("%d\n",value)
    return nil 
}

// PRINTR routine: PrintRegisters()
// Prints all registers and their corresponding values
func (interp *Interpreter) PrintRegisters() error {
    for i := 0; i < 10; i++ {
        i32 := int32(i)
        value, err := interp.ReadFrom(i32) 
        if err != nil {
            return err
        }
        fmt.Printf("R%d: %d\n",i,value)
    }
    return nil
}

// ADDV routine: AddV
// AddV is a "visual add" function which reads the values in ri and rj and 
// prints the values out in *'s one at a time to provide a visual learning
// tool for learning addition. Result must be values must be less than 20
// limit the amount of starts printed to the screen. Stars are printed
// in different colours. 
func (interp *Interpreter) AddV(ri, rj int32) error {
    value1, err := interp.ReadFrom(ri)
    if err != nil {
        return err
    }
    value2, err := interp.ReadFrom(rj)
    if err != nil {
        return err
    }
    if value1 > 10 || value2 > 10 {
        return fmt.Errorf("gvm: ADDV instruction: please use values less than 10 for visual add mode")
    }
    i := int(value1)
    j := int(value2)
    k := int32(i + j)
    interp.WriteTo(ri, k)

    fmt.Printf("%d + %d ",i,j)

    // colour string functions for the first i stars
    starColor1 := color.New(color.FgRed).SprintFunc()
    message1 := strings.Repeat("* ",i)
    for _, char := range message1 {
        fmt.Print(starColor1(string(char))) // applies function to string
        time.Sleep(100 * time.Millisecond)
        os.Stdout.Sync() 
    }

    time.Sleep(100 * time.Millisecond)
    fmt.Printf("+ ")
    time.Sleep(100 * time.Millisecond)

    // the j stars 
    message2 := strings.Repeat("* ",j)
    starColor2 := color.New(color.FgBlue).SprintFunc()
    for _, char := range message2 {
        fmt.Print(starColor2(string(char)))
        time.Sleep(100 * time.Millisecond)
        os.Stdout.Sync()
    }

    time.Sleep(100 * time.Millisecond)
    fmt.Printf("= ")

    // the i + j stars 
    message3 := strings.Repeat("* ",i+j)
    starColor3 := color.New(color.FgGreen).SprintFunc()
    for _, char := range message3 {
        fmt.Print(starColor3(string(char)))
        time.Sleep(100 * time.Millisecond)
        os.Stdout.Sync()
    }
    fmt.Println("")
    return nil
}


// DRAW routine: Draw
// If shape = 1 then a heart is drawn. If shape = 2
// then a bird is drawn. The argument of passed into 
// DrawHeart and DrawBird sets "Blink" to false. 
func (interp *Interpreter) Draw(shape int32) error {
    switch shape {
    case 1:
        interp.DrawHeart(0)
        return nil
    case 2:
        interp.DrawBird(0)
        return nil
    default:
        return fmt.Errorf("interpreter/interpreter.go: invalid shape id: %d",shape)
    }
}

// BLINK routine: Draw
// If shape = 1 then a heart is blinked. If shape = 2
// then a bird is blinked. The argument of passed into 
// DrawHeart and DrawBird sets "Blink" to true. 
func (interp *Interpreter) Blink(shape int32) error {
    switch shape {
    case 1:
        interp.DrawHeart(1)
        return nil
    case 2:
        interp.DrawBird(1)
        return nil
    default:
        return fmt.Errorf("interpreter/interpreter.go: invalid shape id: %d",shape)
    }
}

// DrawHeart prints a heart shape to the screen. If blink is
// set to 1 then the heart will blink on the screen. 
func (interp *Interpreter) DrawHeart(blink int) {
    heart := `
    ."". ."".
    |   '   |
     \     /
      '. .'
        '
`
    if blink == 0 {
        color.Red(heart)
    } else {
        blinkHeart := color.New(color.FgRed, color.BlinkSlow).SprintFunc()
        print(blinkHeart(heart))
    }
    return
}

// DrawBird prints a bird shape to the screen. If blink is set to 
// 1 then the bird image will blink on the screen. 
func (interp *Interpreter) DrawBird(blink int) {
    bird := `

       \\
       (o>
    \\_//)
     \_/_)
      _|_
      `
    if blink == 0 {
        color.Blue(bird)
    } else {
        blinkBird := color.New(color.FgBlue, color.BlinkSlow).SprintFunc()
        print(blinkBird(bird))
    }
    return
}
