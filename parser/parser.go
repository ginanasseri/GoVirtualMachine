// Package parser implements the 'lexer' package to obtain a token stream 
// for a line of input, where each line of input is a Susan instruction.
//
// The parser verifies the syntax of the current instruction and, if 
// no syntax errors are detected, creates a bytecode representation of the 
// instruction using the data structure Instruction defined in the 
// 'instruction' package along with the opcodes listed below. 
//
// If a syntax error is detected by the parser or propagated from a return
// from the 'lexer' package, then the error is propagated to the 'vm'
// package where it is handled. 
package parser

import (
    "fmt"
    "gvm/token"
    "gvm/lexer"
    "gvm/instructions"
)

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

type Parser struct {
    Lex *lexer.Lexer
    CurrentToken *token.Token
}

// Initialize a parser with a lexer, giving it as an argument the current 
// instruction to be parsed in the user's source code, and set the 
// current token as the first token in the input returned from the lexer. 
func New(input string) (*Parser, error) {
    lex := lexer.New(input)
    currentToken, err := lex.GetNextToken()
    return &Parser{Lex: lex, CurrentToken: currentToken}, err
}

// Consume ensures that the token type of the current token being parsed
// is consistent with the expected type according to the syntax of the 
// language. If the types match, then Consume attempts to get the next
// token from the lexer. If the types do not match or the lexer returns
// an error, then a type mismatch or lexer error is returned.
// Otherwise, the parser's current token advances to the next token in 
// the input to be parsed. 
func (p *Parser) Consume(expectedType string) error {
    if p.CurrentToken.TokenType != expectedType {
        return fmt.Errorf("gvm: syntax error: unexpected %s", p.CurrentToken.TokenType)
    }
    var err error
    p.CurrentToken, err = p.Lex.GetNextToken()
    if err != nil {
        return err
    }
    return nil
}

// Instruction creates bytecode instructions from a stream of tokens
// as they are parsed. If the current token type is a specific 
// instruction, then Instruction checks that the instruction syntax
// is correct and all required parameters have been provided. If the
// instruction syntax is correct, then a bytecode instruction is 
// returned to the interpreter
// reads the current type of instruction being parsed and 
// attempts to create a bytecode instruction
//
// The expected syntax for each instruction INSTR is provided 
// as a comment following the 'case INSTR:' line. 
//
// Instruction types:
//  - BinaryInstruction types require two arguments
//  - UnaryInstruction types require one argument
//  - NullaryInstruction have zero 

func (p *Parser) Instruction() (instructions.Instruction, error) {
    currentToken := p.CurrentToken
    switch currentToken.TokenType {

    case token.JUMP: 
        // JUMP INT (address)
        if err := p.Consume(token.JUMP); err != nil {
            return instructions.NewError(err), err
        }
        currentToken = p.CurrentToken
        if err := p.Consume(token.INT); err != nil {
            return instructions.NewError(err), err
        }
        jumpTo := currentToken.Value
        return instructions.NewUnaryInstruction(int32(OPCODE_JUMP), jumpTo), nil

    case token.LDI:
        // LDI REG COMMA INT
        if err := p.Consume(token.LDI); err != nil {
            return instructions.NewError(err), err
        }
        // REG (load to)
        currentToken = p.CurrentToken
        if err := p.Consume(token.REG); err != nil {
            return instructions.NewError(err),err
        }
        regIndex := currentToken.Value
        
        // COMMA 
        if err := p.Consume(token.COMMA); err != nil {
            return instructions.NewError(err), err
        }
        // INT (load value)
        currentToken = p.CurrentToken
        if err := p.Consume(token.INT); err != nil {
            return instructions.NewError(err), err
        }
        loadValue := currentToken.Value

        // return LDI instruction bytecode 
        return instructions.NewBinaryInstruction(int32(OPCODE_LDI),regIndex,loadValue), nil
    
    case token.ADD: 
        // ADD REG COMMA REG 
        if err := p.Consume(token.ADD); err != nil {
            return instructions.NewError(err), err
        }
        // REG - first add register 
        currentToken = p.CurrentToken
        if err := p.Consume(token.REG); err != nil {
            return instructions.NewError(err), err
        }
        regIndex1 := currentToken.Value

        // Check for comma 
        if err := p.Consume(token.COMMA); err != nil {
            return instructions.NewError(err), err
        }
        // REG - second add register 
        currentToken = p.CurrentToken
        if err := p.Consume(token.REG); err != nil {
            return instructions.NewError(err), err
        }
        regIndex2 := currentToken.Value
        
        // Return ADD instruction bytecode 
        add_instr := instructions.NewBinaryInstruction(int32(OPCODE_ADD), regIndex1, regIndex2)
        return add_instr, nil

    case token.ADDV:
        // ADD REG COMMA REG
        if err := p.Consume(token.ADDV); err != nil {
            return instructions.NewError(err), err
        }
        // REG
        currentToken = p.CurrentToken
        if err := p.Consume(token.REG); err != nil {
            return instructions.NewError(err), err
        }
        regIndex1 := currentToken.Value

        // COMMA
        if err := p.Consume(token.COMMA); err != nil {
            return instructions.NewError(err), err
        }
        // REG
        currentToken = p.CurrentToken
        if err := p.Consume(token.REG); err != nil {
            return instructions.NewError(err), err
        }
        regIndex2 := currentToken.Value

        addv_instr := instructions.NewBinaryInstruction(int32(OPCODE_ADDV),regIndex1,regIndex2)
        return addv_instr, nil

    case token.DRAW:
        // DRAW SHAPE
        if err := p.Consume(token.DRAW); err != nil {
            return instructions.NewError(err), err
        }
        currentToken = p.CurrentToken
        if err := p.Consume(token.SHAPE); err != nil {
            return instructions.NewError(err), err
        }
        shape := currentToken.Value
        return instructions.NewUnaryInstruction(int32(OPCODE_DRAW),shape), nil

    case token.BLINK:
        // BLINK SHAPE
        if err := p.Consume(token.BLINK); err != nil {
            return instructions.NewError(err), err
        }
        // SHAPE
        currentToken = p.CurrentToken
        if err := p.Consume(token.SHAPE); err != nil {
            return instructions.NewError(err), err
        }
        shape := currentToken.Value
        
        // Return BLINK bytecode instruction
        return instructions.NewUnaryInstruction(int32(OPCODE_BLINK),shape), nil
    
    case token.STDOUT:
        // STDOUT REG
        if err := p.Consume(token.STDOUT); err != nil {
            return instructions.NewError(err), err
        }
        // REG - print from register 
        currentToken = p.CurrentToken
        if err := p.Consume(token.REG); err != nil {
            return instructions.NewError(err),err
        }
        regIndex := currentToken.Value

        // Return STDOUT bytecode instruction
        return instructions.NewUnaryInstruction(int32(OPCODE_STDOUT), regIndex), nil

    case token.PRINTR:
        if err := p.Consume(token.PRINTR); err != nil {
            return instructions.NewError(err), err
        }
        return instructions.NewNullaryInstruction(int32(OPCODE_PRINTR)),nil

    default:
        err := fmt.Errorf("gvm: default case: invalid '%v'",currentToken)
        return instructions.NewError(err),err
    }
}
