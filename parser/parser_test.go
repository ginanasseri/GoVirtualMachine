/*
Commands:

STDOUT: COMMAND
LDI
JUMP
ADD
ADDV
DRAW
BLINK

*/
package parser

import (
    "testing"
)

type TestCase struct {
    input string
    shouldPass bool
}

func TestParser(t *testing.T) {
    testCases := []TestCase{

        // STDOUT REG
        {"STDOUT r1", true},
        {"STDOUTr1", false},
        {"STDOUT3", false},
        {"STDOUT 3", false},
        {" r1 STDOUT", false},
        {" ,STDOUT", false},
        {"STDOUT,", false},


        // LDI REG COMMA INT
        {"LDI r1,3", true},
        {"LDI r1 , 3 ", true},
        {"LDI , 3", false},
        {"LDIr1, 3", false},
        {"LDI r1 3", false},
        {"LDI r1, ", false},
        {"LDI 3, r1", false},
        {" 3LDI r1,3", false},
       // {" $LDI r1,3", false},
        {"LDI r8, 89", true},


        // JUMP INT
        {"JUMP 4", true},
        {"JUMP r4", false},
        {"JUMP4", false},
        {"JUMP , 4", false},
        {"JUMP ", false},
        {" 4 JUMP", false},
        
        // ADD REG COMMA REG 
        {"ADD r1,r2",true},
        {"  ADD R1,   R2 ", true},
        {"ADDr1, r2", false},
        {"ADD r1 r2", false},
        {"ADD rr2, r1", false},
        {"ADD , r2", false},
        {"ADD r1 r2", false},
        {"ADD r1 ,", false},
        {" r1 ADD , r2", false},
        {"ADD 1,r2", false},
        {"ADD r1, 2", false},
        {"ADD 1, 22", false},
        {" $ADD r1,r2", false},

        // ADDV REG COMMA REG
        {"ADDV r1,r2",true},
        {"  ADDV R1,   R2 ", true},
        {"ADDV , r2", false},
        {"ADDV r1 r2", false},
        {"ADDV r1 ,", false},
        {" r1 ADDV , r2", false},
        {"ADDV 1,r2", false},
        {"ADDV r1, 2", false},
        {"ADDV 1, 22", false},
        {" $ADDV r1,r2", false},

        // DRAW SHAPE 
        {"DRAW $heart",true},
        {"DRAW $bird", true},
        {"BLINK $heart", true},
        {"BLINK $bird", true},
        {"DRAW $hi", false},
        {"DRAW", false},
        {"BLINK", false},
        {"DRAW$heart", false},
        {"BLINK$heart", false},
        {"DRAW $HEART", false},
    }

    for _, testCase := range testCases {

        parse, err1 := New(testCase.input)
        
        // ignore returned error from first token as this was already tested in the lexer_test
        if err1 != nil {
            continue 
        }

        _, err2 := parse.Instruction()

        if err2 == nil && !testCase.shouldPass {
            t.Errorf("FAIL: no error returned from invalid input: %s", testCase.input)
        }
        if err2 != nil && testCase.shouldPass {
            t.Errorf("FAIL: error returned from valid input: %s: error message: %v", testCase.input, err2)
        }
    }
}
