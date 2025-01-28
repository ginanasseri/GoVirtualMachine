// Package lexer performs lexical analysis on the user's source code,
// provided as an input string), and returns tokens of the type 'Token'
// as defined in the 'token' package. The tokens are provided to the 
// 'parser' package for further syntax analysis. If an invalid token
// is encountered, the error is propagated to the 'parser' package
// where it is handled. 
package lexer

import (
     "fmt"
     "unicode"
     "strings"
     "strconv"
     "math"
     "gvm/token"
)

type Lexer struct {
    Input string 
    Position int
    CurrentChar byte
}

// Initialize a lexer with an input string. The input string is a 
// line in the user's source code representing on instruction. The
// current character is initialzed as the first character in the 
// input. 
func New(input string) *Lexer {
    var currentChar byte = 0
    if len(input) > 0 {
        currentChar = input[0]
    } 
    return &Lexer{Input: input, Position: 0, CurrentChar: currentChar} 
}

// GetNextChar advance position to the next character in the input for 
// analysis
func (lex *Lexer) GetNextChar() {
    lex.Position += 1
    if lex.Position > len(lex.Input) - 1 {
        lex.CurrentChar = 0
     } else {
        lex.CurrentChar = lex.Input[lex.Position]
    }
    return
}

// IgnoreWhiteSpace advances the current position in the input if a 
// whitespace is encountered 
func (lex *Lexer) IgnoreWhiteSpace() {
    for lex.CurrentChar != 0 && unicode.IsSpace(rune(lex.CurrentChar)) {
        lex.GetNextChar()
    }
    return
}

// Delimiter determines if the character immediately preceeding the 
// current character being processed was a comma (',') or a space (' ').
// This is used to catch the syntax error case where delimiters are 
// missing between valid tokens. (E.g., JUMP2, DRAW$heart). Note that
// it is sufficient to check if a '$' is preceeded by a comma or a space
// as the parser will catch an error if a '$' is preceeded by a comma.
func (lex *Lexer) Delimiter() bool {
    if lex.Position > 0 {
        previousChar := lex.Input[lex.Position - 1]
        return previousChar == ' ' || previousChar == ','
    }
    return false // position = 0
}

// Register checks if the current character in the input is an 'r' or
// 'R', signalling that a register (type REG) token should be processed.
func (lex *Lexer) Register() bool {
    if lex.Position > 0 {
        previousChar := lex.Input[lex.Position - 1]
        return previousChar == 'r' || previousChar == 'R'
    }
    return false
}

// Integer parses a multi-digit integer string and returns it as an int32 value. 
// int32 values are used to be consistent with the ISA of the native OS which 
// the source code was written for as register values in the ISA hold int32 values.
func (lex *Lexer) Integer() (int32, error) {
    if !lex.Register() {
        if !lex.Delimiter() {
            return 0, fmt.Errorf("gvm: syntax error: unexpected INT (missing delimiter)")
        }
    }
    integerString := ""
    for lex.CurrentChar != 0 && unicode.IsDigit(rune(lex.CurrentChar)) {
        integerString += string(lex.CurrentChar)
        lex.GetNextChar()
    }
    integer, err := strconv.ParseInt(integerString, 10, 32) // returns an int64 value 
    if err != nil {
        return 0, fmt.Errorf("gvm: lexer: failed to convert input to integer %w",err)
    }
    // check for integer overflow
    if integer > math.MaxInt32 {
        return 0, fmt.Errorf("gvm: register integer overflow error")
    }
    // check for integer underflow
    if integer < math.MinInt32 {
        return 0, fmt.Errorf("gvm: register integer underflow error")
    }
    integer32 := int32(integer) // convert to int32 
    return integer32, nil    
}

// RegisterIndex attempts to obtain a valid register immediately following an
// 'r' or 'R'. The register must be preceeded by a comma (',') or space (' ')
// and must be less than two digits. 
func (lex *Lexer) RegisterIndex() (int32, error) {

    // if previous char was not a comma or a space, then this r should not be here.
    if !lex.Delimiter() {
        return 0, fmt.Errorf("gvm: unexpected REG (missing delimiter)")
    }
    // otherwise, advance to next character and get the register index.
    lex.GetNextChar()
    if unicode.IsSpace(rune(lex.CurrentChar)) {
        return 0, fmt.Errorf("gvm: missing register index.")
    }
    registerIndex, err := lex.Integer() 
    if err != nil {
        return 0, err
    }
    if registerIndex > 9 {
        return 0, fmt.Errorf("gvm: register indices must be between 0 and 9.")
    }
    return registerIndex, nil
}

// Command builds command string in user input. Commands are uppercase strings 
// with a maximum length of 5.
// If the command string exceeds the maximum length, then an error is returned 
// with an empty string. Otherwise, the command string is returned for further 
// verification. 
func (lex *Lexer) Command() (string, error) {
    var builder strings.Builder
    for lex.CurrentChar != 0 && unicode.IsUpper(rune(lex.CurrentChar)) {
        if builder.Len() > 10 {
            return "", fmt.Errorf("gvm: invalid command: max length reached (10).")
        }
        builder.WriteRune(rune(lex.CurrentChar))
        lex.GetNextChar()
    }
    return builder.String(), nil
}

// Shape is called immediately following a '$' character, which denotes a shape 
// type value. Shapes are lowercase strings with a maximum length of 5. 
// Shape returns an error if the shape string exceeds the maximum length, if 
// the '$' was missing a delimiter, or if no shape was provided. 
// If no error occured, then the shape string is returned for further verification. 
func (lex *Lexer) Shape() (string, error) {
    if !lex.Delimiter() {
        return "", fmt.Errorf("gvm: shape declaration must be preceded by ' '")
    }
    lex.GetNextChar()
    var builder strings.Builder 
    for lex.CurrentChar != 0 && unicode.IsLower(rune(lex.CurrentChar)) {
        if builder.Len() > 6 {
            return "", fmt.Errorf("gvm: invalid shape.")
        }
        builder.WriteRune(rune(lex.CurrentChar))
        lex.GetNextChar()
    }
    if builder.Len() == 0 {
        return "", fmt.Errorf("gvm: missing shape after '$'")
    }
    return builder.String(), nil
}

// GetNextToken is the generating function of the lexer where the input is analyzed
// one character at a time. The input is traversed and tokens are built using the 
// helper functions above. If an error is returned from any helper function, then a
// nil token is returned and the error is propagated to the parser where it is 
// handled. 
//
// REG, INT, and COMMA tokens are immediately returned along with a nil 
// error once successfully created.
//
// COMMAND and SHAPE tokens require further verification. If the command or shape
// is not valid within the source's ISA, then a nil token and error are returned
// to the parser. Otherwise, the COMMAND/SHAPE tokens are returned with a nil 
// error. 
func (lex *Lexer) GetNextToken() (*token.Token, error) {

    for lex.CurrentChar != 0 {
        switch {

        // Whitespace
        case unicode.IsSpace(rune(lex.CurrentChar)): 
            lex.IgnoreWhiteSpace()

        // Register 
        case unicode.ToLower(rune(lex.CurrentChar)) == 'r':  
            regIndex, err := lex.RegisterIndex()
            if err != nil {
               return nil, err
            }
            return token.New(token.REG, regIndex), nil 

        // Integer
        case unicode.IsDigit(rune(lex.CurrentChar)):
            integer, err := lex.Integer()
            if err != nil {
                return nil, err
            }
            return token.New(token.INT, integer), nil

        // Comma 
        case lex.CurrentChar == ',':
            lex.GetNextChar()
            return token.New(token.COMMA,0), nil

        // Shape
        case lex.CurrentChar == '$':
            shape, err := lex.Shape()
            if err != nil {
                return nil, err
            }
            // ensure shape is valid 
            switch shape {
            case "heart":
                return token.New(token.SHAPE,1), nil // 1 denotes heart
            case "bird":
                return token.New(token.SHAPE,2), nil // 2 denotes bird
            default:
                return nil, fmt.Errorf("gvm: invalid shape: '%s'",shape)
            }
                
        // Lowercase letter which is not 'r' - invalid 
        case unicode.IsLower(rune(lex.CurrentChar)):
            return nil, fmt.Errorf("gvm: input is case sensitive: invalid '%c'",rune(lex.CurrentChar))
    
        // Punctiation or symbol which is not ',' - invalid
        case unicode.IsPunct(rune(lex.CurrentChar)) || unicode.IsSymbol(rune(lex.CurrentChar)):
            return nil, fmt.Errorf("gvm: invalid character: %c",rune(lex.CurrentChar))

        // Command 
        case unicode.IsUpper(rune(lex.CurrentChar)):
            command, err := lex.Command()
            if err != nil {
                return nil, err
            }
            // ensure command is valid 
            switch command {
            case "LDI":
                return token.New(token.LDI, 0), nil
            case "STDOUT":
                return token.New(token.STDOUT, 0), nil
            case "PRINTR":
                return token.New(token.PRINTR,0), nil
            case "JUMP":
                return token.New(token.JUMP, 0), nil
            case "ADD":
                return token.New(token.ADD, 0), nil
            case "ADDV":
                return token.New(token.ADDV, 0), nil
            case "DRAW":
                return token.New(token.DRAW, 0), nil
            case "BLINK":
                return token.New(token.BLINK, 0), nil
            default:
                return nil, fmt.Errorf("gvm: undefined: '%s'.",command)
            }
        // something else 
        default:
            return nil, fmt.Errorf("gvm: unrecognized symbol '%c'",rune(lex.CurrentChar))
        }
    }
    // EOF is a dummy-type for final return. 
    return token.New(token.EOF, 0), nil
}
