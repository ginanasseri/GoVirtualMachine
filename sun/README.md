# Example Programs

```text
susan0: Load 1 into R1 and 8 into R8, add value of R8 to R1, and display R1 value. [Output: 9]
0: LDI r1, 1
1: LDI r2, 8
2: ADD r1,r2
3: STDOUT r1

susan1: Load 3 values to R8 and add to R1 [Output: 2, 6, 14]
0: LDI r1, 2
1: STDOUT r1
2: LDI r8, 4
3: ADD r1,r8
4: STDOUT r1
5: LDI r8, 8
6: ADD r1, r8
7: STDOUT r1

susan2: JUMP example: jump at line 3 to line 5, skips 2nd ADD instruction [output: 2, 2, 10]
0: LDI r1, 2
1: STDOUT r1
2: LDI r8, 4
3: JUMP 5
4: ADD r1,r8
5: STDOUT r1
6: LDI r8, 8
7: ADD r1, r8
8: STDOUT r1

susan3: Print all register values with PRINTR. R0 stores the PC which starts at 1 and points to next instruction [R0 = 7]. 
        Unused register values are 0. 
0: LDI r1, 3
1: LDI r4, 9
2: LDI r8, 22
3: LDI r3, 19
4: LDI r2, 2
5: PRINTR

susan4: Print some stars
0: LDI r1, 2
1: LDI r2, 2
2: ADDV r1, r2
3: ADDV r1, r2
4: ADDV r1,r2
5: ADDV r1,r2

susan5: Display a heart with blinking on
0: BLINK $heart

susan6: Display a bird with blinking off
0: DRAW $bird
```

## Error Handling Examples:
```text
write_permission: Write attempt on read-only register.
0: LDI r5, 12
1: STDOUT r5
2: LDI r0, 1
3: ADD r5, r1
4: STDOUT r5

    Output:
            12
            gvm: write to R0: permission denied [R0 is read-only]

infinite: Invalid JUMP exception handling: INFINITE LOOP WARNING
0: LDI r1, 2
1: STDOUT r1
2: LDI r8, 4
3: JUMP 3
4: ADD r1,r8
5: STDOUT r1
6: LDI r8, 8
7: ADD r1, r8
8: STDOUT r1

    Output: 
            2
            gvm: JUMP at addr 3 to 3: infinite loop warning.


offpage: Invalid JUMP segmentation violation: Jump to non-existant address
0: LDI r1, 2
1: STDOUT r1
2: LDI r8, 4
3: JUMP 10
4: ADD r1,r8
5: STDOUT r1
6: LDI r8, 8
7: ADD r1, r8
8: STDOUT r1

    Output:
            2
            gvm: JUMP addr invalid: segmentation violation.


case0: Invalid token error
0: LID r1, 2
1: STDOUT r1

    Output: 
            gvm: undefined: 'LID'.

case1: Missing space between tokens (Note that whitespace is ignored. E.g., 'LDI    r1,   9' is accepted)
0: LDI r1, 1
1: LDI r2, 8
2: ADDr1,r2
3: STDOUT r1

    Output:
            gvm: unexpected REG (missing delimiter)    

Note that separate error messages are displayed for any token out of order. E.g.:
- LDI r1,,8 will output: unexpected ','
- LDI 8 will output:     expected REG

Any symbols in program not defined in the language, e.g., '.' will output: gvm: invalid character .
``` 
