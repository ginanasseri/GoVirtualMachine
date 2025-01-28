// Tests instruction functioning and error handling
// Note that Lexer and Parser tests already 100% passed,
// so testing for syntax errors is not necessary. 
package vm

import (
    "testing"
)

type TestCase struct {
    input string
    shouldPass bool
}

func TestParser(t *testing.T) {
    testCases := []TestCase{
    {"testdata/test0", true},
    {"testdata/test1", false},
    {"testdata/test2", false},
    {"testdata/test3", true},
    {"testdata/test4", false},
    {"testdata/test5", true},
    {"testdata/test6", true},
    {"testdata/test7", false},
    {"testdata/test8", false},
    {"testdata/test9", true},
    {"testdata/test10", true},
    {"testdata/test11", false},
    {"testdata/test12", false},
    }

    for _, testCase := range testCases {

        vm := NewVirtualMachine()
        err := vm.Execute(testCase.input)

        if err == nil && !testCase.shouldPass {
            t.Errorf("No error returned from invalid input file: %s", testCase.input)
        }
        if err != nil && testCase.shouldPass {
            t.Errorf("Error returned from valid input file: %s: error message: %v", testCase.input, err)
        }

        if err != nil && !testCase.shouldPass {
            t.Logf("%s got it: %v", testCase.input, err)
        }
    }
}
