# Go Virtual Machine 

The Go Virtual Machine provides an isolated execution environment for executing programs written for Susan; a hypothetical ISA with 10 32-bit CPU registers which is capable of executing basic instructions including read, write, add, jump, and printing operations. 

**Usage:** 
- To run directly, use: `go run main.go`
- To compile and create an executable use: `go build -o gvm` and `./gvm` to run program. 

Once loaded, enter 'run [file]' to execute instructions, where file is a program written in the Susan instruction set. Some samples programs are included within the `sun` directory. There are also several sample files within the `vm/testdata/` directory which demonstrate error handling, including branching error handling. The user can continue running programs or use `exit` to exit.

E.g., run sun/susan5 

### Instruction Set Summary

|Mnemonic|Operands|Description | Operation  |
|:--------|:--------|:-------------|:------------|
| STDOUT |  Rd    |Print register value |   |
|  LDI   |  Rd,K |Load Immediate   | Rd ← K  |
|  JUMP |  K | Jump  | PC ← K  |
|  ADD | Rd,Rr  |Add   | Rd ← Rd + Rk  |
| PRINTR | | Print all registers


### Registers 
Susan has 10 32-bit registers for read and write operations
- Register 0 is a special purpose register which stores the address of the last instruction in the program code. This is used to check if JUMP instructions are valid. This Register is read-only when in execution mode.
- Registers 1:9 are general purpose read-and-write registers. 


### Additional Features: Visual Mode 
|Mnemonic|Operands|Description | Operation  |
|:--------|:--------|:-------------|:------------|
|  ADDV | Rd,Rr  |Visual mode add   | Rd ← Rd + Rk |
| DRAW  | \$s  | Draw shape  |   |
| BLINK  | \$s  | Blink shape  |   |

- ADDV: Rd and Rr values must be ≤ 10. ADDV is an educational feature to visualize how two numbers are added together
- DRAW: Susan can print two shapes: a bird and a heart. Use \\$heart to print a heart and \\$bird to print a bird
- BLINK: functions the same as DRAW, but the image blinks on the screen 

---

# Package Contents and Control Flow

## Package Contents 
- `vm`: The vm package implements the Go Virtual Machine. It contains the virtual machine architecture including the virtual memory structures and the interpreter. It is the point of control transfer between the host OS and the Susan process. When a New Virtual Machine instance is initialized, memory is allocated in a Virtual Memory data structure to hold the executable code and Susan registers. The Virtual Machine is initialized with a pointer to the virtual memory, and an interpreter which is passed a reference to the virtual memory. 
    - **Also includes** `vm_test.go` and the directory `testdata` that contains test cases of valid programs and cases of programs with errors. (tests errors raised by the lexer or parser are correctly propagated to `main` and exception handling is behaving as expected. Fails if any error in a program is not detected.)
- `interpreter`: the interpreter package executes the bytecode instructions contained in the virtual memory executable code block section using the decode and dispatch method. 
- `parser`: the parser package implements the lexer to obtain a token stream from a line of input, where each line is a Susan instruction, and creates a bytecode representation of each instruction. 
    - **Also includes** `parser_test.go` (tests that the parser correctly accepted token streams with valid syntax (e.g., ADD r1, r2) and rejecting token streams with invalid syntax). 
- `lexer`: the lexer breaks down the Susan source code file into tokens. 
    - **Also includes** `lexer_test.go` (tests that the lexer is correctly accepting all valid token types (e.g., INT, REG, etc.) and rejecting any token not defined in the language.)
- `instructions`: instructions defines the data type Instruction which represent the bytecode instructions created by the parser
- `token`: token defines the data type Token which represent the input tokens created by the lexer. 


## Control Flow 
After the user enters 'run file' where 'file' is a valid Susan program:
1. `main`: The file is verified and, once verified, a new Virtual Machine instance is initialized and control is transferred to the `vm` package by a call to `vm.Execute()` 
2. `vm.Execute()`: GVM loads the source program via a call to the host OS and, if the file is successfully loaded, then control is transferred to `vm.ParseInstructions()` to get the source program's executable. 
3. `vm.ParseInstructions()`: Still working through host system calls, the source program is scanned line by line and, for each line, the `parser` is invoked to tokenize the input while simultaneously checking the syntax of each instruction. After an instruction is sucessfully parsed, a representative bytecode for the instruction is created which contains an encoding of the instruction's opcode and location of its operands, if applicable. If any syntax errors are detected, the error is propagated back to `main` and printed. The user can then terminate the machine or fix their mistake and run the program/run a different program. If no syntax errors are detected, then the bytecode instructions are written into the VirtualMemory executable code block and control is transferred back to `vm.Execute()`
4. `vm.Execute()`: now working entirely within the virtual machine environment, the last address of the executable codeblock is written into register 0 and the interpreter is invoked. 
5. `vm.interpreter.Interpret()`: The interpreter executes each instruction using the decode and dispatch method and working, operating solely on virtual memory structures. If any error occurs, it is propagated back to main similarly as described in step 3. If a JUMP instruction is included, then not all instructions need be executed. The interpreter manages its own Program Counter to traverse the code block and reads register 0 to know when it has reached the last instruction. Once the program counter value is equal to the value at register 0, control is transferred back to `vm.Execute()`
6. `vm.Execute()`: returns control back to main
7. `main`: user can exit or run another program. The Virtual Memory is cleared and reset for each program run by initializing a new VirtualMachine instance for each program executed.
