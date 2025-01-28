// Package instruction provides data structure used to represent bytecode instructions
// generated by the 'parser' package and executed by the 'interpreter' package. The 
// interface provides helper functions used by the interpreter to obtain the 
// commands and argument values represented by the bytecode instructions. 
// executed by the 'interpreter' package. 
//
// Instruction types: 
//  - BinaryInstructions contain 2 parameters  - e.g. ADD r1,r2
//  - UnaryInstructions contain 1 parameter    - e.g. JUMP 3
//  - NullaryInstructions contain 0 parameters - e.g. PRINTR 
//  - ErrorInstructions provide a consistent return type in the case of errors
//    and preserve the error type which occured during a failed instruction 
//    execution 
//
// Type int32 values are used to be consistent with the registers of the ISA
// being emulatd. 
package instructions

import (
    "fmt"
)

type Instruction interface {
    GetOpCode() int32
    GetArg1()   int32
    GetArg2()   int32
    String()    string
}

// Binary instructions: instructions with 2 parameters 
type BinaryInstruction struct {
    OpCode, Arg1, Arg2 int32
}

func NewBinaryInstruction(opcode, arg1, arg2 int32) Instruction {
    binaryInstruction := &BinaryInstruction{
        OpCode: opcode,
        Arg1: arg1,
        Arg2: arg2,
    }
    return binaryInstruction
}

func (bi *BinaryInstruction) GetOpCode() int32 {
    return int32(bi.OpCode)
}

func (bi *BinaryInstruction) GetArg1() int32 {
    return int32(bi.Arg1)
}

func (bi *BinaryInstruction) GetArg2() int32 {
    return int32(bi.Arg2)
}

func (bi *BinaryInstruction) String() string {
    return fmt.Sprintf("/%#02x,%d,%d}", bi.OpCode, bi.Arg1, bi.Arg2)
}

// Unary Instructions: instructions with a single argument 
type UnaryInstruction struct {
    OpCode int32
    Arg1   int32
}

func NewUnaryInstruction(opcode,arg1 int32) Instruction {
    return &UnaryInstruction{OpCode: opcode, Arg1: arg1}
}

func (ui *UnaryInstruction) String() string {
    return fmt.Sprintf("/%#02x,%d}", ui.OpCode, ui.Arg1)
}

func (ui *UnaryInstruction) GetOpCode() int32 {
    return ui.OpCode
}

func (ui * UnaryInstruction) GetArg1() int32 {
    return ui.Arg1
}

func (ui * UnaryInstruction) GetArg2() int32 {
    return 0
}

// Nullary Instructions: instructions with no arguments  
type NullaryInstruction struct {
    OpCode int32
}

func NewNullaryInstruction(opcode int32) Instruction {
    return &NullaryInstruction{OpCode: opcode}
}

func (ni *NullaryInstruction) String() string {
    return fmt.Sprintf("/%02x}", ni.OpCode)
}

func (ni *NullaryInstruction) GetOpCode() int32 {
    return ni.OpCode
}

func (ni * NullaryInstruction) GetArg1() int32 {
    return 0 
}

func (ni * NullaryInstruction) GetArg2() int32 {
    return 0
}

// Error 
type Error struct {
    ErrorType error
}

func NewError(err error) Instruction {
    return &Error{ErrorType: err}
}

func (e *Error) String() string {
    return fmt.Sprintf("E: %v", e.ErrorType)
}

func (e *Error) GetOpCode() int32 {
    return 0
}

func (e *Error) GetArg1() int32 {
    return 0 
}

func (e *Error) GetArg2() int32 {
    return 0
}
