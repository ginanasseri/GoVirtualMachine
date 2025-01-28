// Package token provides data struture used to represent tokens
// in lexical analysis performed by the 'lexer' package. These
// tokens serve as the building blocks for parsing source code.
package token

import (
    "fmt"
)

// Defining token types as constants 
const (
    INT     = "INT"
    REG     = "REG"
    COMMA   = "COMMA"
    LDI     = "LDI"
    STDOUT  = "STDOUT"
    JUMP    = "JUMP"
    ADD     = "ADD"
    ADDV    = "ADDV"
    DRAW    = "DRAW"
    BLINK   = "BLINK"
    SHAPE   = "SHAPE"
    PRINTR  = "PRINTR"
    EOF     = "EOF"
)

// Type int32 is used to be consistent with the source's ISA
type Token struct {
    TokenType string
    Value     int32
}

func New(tokenType string, value int32) *Token {
    return &Token{TokenType: tokenType, Value: value}
}

func (t Token) String() string {
    return fmt.Sprintf("{Type: %s | Value: %d}", t.TokenType, t.Value)
}
