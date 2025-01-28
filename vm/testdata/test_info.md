```text
test0: Expected output: 9 
    0: LDI r1, 1
    1: LDI r2, 8
    2: ADD r1,r2
    3: STDOUT r1

test1: Expected output: syntax error propagated from parser 
    0: LDI r1, 1
    1: LDI r2,,8
    2: ADD r1,r2
    3: STDOUT r1

test2: Expected output: syntax error propagated from lexer
    0: LDIr1, 1
    1: LDI r2,8
    2: ADD r1,r2
    3: STDOUT r1

test3: Expected output: 4 (N instructions)
    0: LDI r1, 1
    1: LDI r2,8
    2: ADD r1,r2
    3: STDOUT r0

test4: Expected output: syntax error propagated from parser  
    0: LDI r7, 99
    1: LDI r6, 
    2: ADD r7,r6
    3: STDOUT r7

test5: Expected output: 112
    0: LDI r7, 99
    1: LDI r6, 13 
    2: ADD r7,r6
    3: STDOUT r7

test6: Expected output: 13
    0: LDI r7, 99
    1: LDI r6, 13
    2: ADD r7,r6
    3: STDOUT r6

test7: Expected output: invalid register (out of range)
    0: LDI r10, 1
    1: LDI r6, 2
    2: ADD r7,r6
    3: STDOUT r10

test8: Expected output: invalid register write permission denied
    0: LDI r6, 5
    1: LDI r0, 2
    2: ADD r6,r0
    3: STDOUT r6

test9: Expected output: 0
    0: LDI r1, 12
    1: LDI r3, 13
    2: ADD r7,r6
    3: STDOUT r6

test10: Expected output: 1
    0: LDI r1, 1
    1: LDI r2, 8
    2: JUMP 4
    3: ADD r1,r2
    4: STDOUT r1

test11: Expected output: infinite loop warning 
    0: LDI r1, 1
    1: LDI r2, 8
    2: JUMP 1
    3: ADD r1,r2
    4: STDOUT r1

test12: Expected output: segmentation violation 
    0: LDI r1, 1
    1: LDI r2, 8
    2: JUMP 5
    3: ADD r1,r2
    4: STDOUT r1
```
