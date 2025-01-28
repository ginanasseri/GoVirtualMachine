package lexer

import (
    "testing"
)

type TestCase struct {
    input string
    shouldPass bool
}

func TestLexer(t *testing.T) {
    testCases := []TestCase{

        // register test
        {"DRAW",true},
        {" $heart",true},
        {" r1",true},
        {" R1", true},
        {" r99", false},
        {" r12",false},
        {",r1",true},
        {" r999",false},  // out of bounds register
        {"dr1", false},  // no space before register
        {"Dr22", false}, 
        {" R", false},    // missing register index
        {" r", false},
        {" r 2", false},
    

        // integers 
        {" 0", true},
        {" 1243", true},
        {" 10983",true},
        {"a123", false},
        {"123", false},

        // symbols 
        {",", true},
        {".", false},
        {"&", false},

        // commands
        {"ADD", true},
        {"LDI", true},
        {"STDOUT", true},
        {"    ADD  ",true},
        {"AD", false},  // unknown command
        {"add", false},  // lowercase
        {"AdD", false}, 
        {"aDD", false},
        {"PEW", false},
        {"A1DD", false},
        {"ADD1", true}, // these are true because INT and REG have their own checks.
        {"ADDr", true},
        {" r", false},
        {"r",false},
        {"1",false},
    }

    for _, testCase := range testCases {
        lex := New(testCase.input)
        _, err := lex.GetNextToken()

        if err == nil && !testCase.shouldPass {
            t.Errorf("FAIL: no error returned from invalid input: %s", testCase.input)
        }
        if err != nil && testCase.shouldPass {
            t.Errorf("FAIL: error returned from valid input: %s: error message: %v", testCase.input, err)
        }
    }
}
