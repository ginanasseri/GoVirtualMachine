// Package main initializes the Virtual Machine. To execute a program, 
// use 'run file' where 'file' is the name of your program. Enter 
// exit to exit. 
package main

import (
    "fmt"
    "strings"
    "time"
    "os"
    "bufio"
    "gvm/vm"
)

func hello() {
//    fmt.Printf("Starting GVM: \n")
    for i := 0; i <= 100; i += 10 {
        fmt.Printf("\r[%-10s] %d%% Complete", strings.Repeat("#", i/10), i)
        time.Sleep(150 * time.Millisecond)
    }
    fmt.Printf("\nWelcome! Use 'run [filename]' to execute a Susan program or EXIT to exit.\n")
    return
}

func main() {

    hello()
    scanner := bufio.NewScanner(os.Stdin)

    // this code section prompts for filename of program user wishes to run 
    for {
        fmt.Print(">> ")
        if !scanner.Scan() {
            fmt.Println("gvm: error reading from STDIN channel: exiting program.\n")
            break // EOF or error
        }
        input := scanner.Text()
        // user hits enter
        if input == "" {
            continue 
        }
        // splitting input 
        parts := strings.Fields(input)
        // user exits (case insensitive)
        if strings.EqualFold(parts[0], "EXIT") {
            break
        }
        if parts[0] != "run" {
            fmt.Printf("gvm: invalid input: use 'run [filename]' to execute program or EXIT to exit.\n")
            continue
        }
        // if we're here, then parts[0] is run
        if len(parts) < 2 {
            fmt.Printf("gvm: missing filename\n")
            continue 
        }
        if len(parts) > 2 {
            fmt.Printf("gvm: too many arguments\n")
            continue
        }
        // if we're here then we have exactly 2 arguments and can verify file 
        filename := parts[1]
        if _, err := os.Stat(filename); err != nil {
            if os.IsNotExist(err) {
                fmt.Printf("gvm: file not found in directory: %s\n",filename)
            } else {
                fmt.Printf("gvm: file error: %v\n",err)
            }
            continue 
        }
        // if we're here, we have a valid file and can initialize the VM and 
        // execute the source program 
        vm := vm.NewVirtualMachine()
        if err := vm.Execute(filename); err != nil {
            fmt.Printf("%v\n",err)
        }
     }
}
