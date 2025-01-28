```text
Susan Sample Programs

susan0: ~ Add Two Numbers ~ [prints 9]
0: LDI r1, 1
1: LDI r2, 8
2: ADD r1,r2
3: STDOUT r1

susan1: ~ Add Three Numbers ~ [prints 2, 4, 6]
0: LDI r1, 2
1: STDOUT r1
2: LDI r8, 2
3: ADD r1,r8
4: STDOUT r1
5: LDI r8, 2
6: ADD r1, r8
7: STDOUT r1

susan2: ~ Jump Past an Add ~ [prints 2, 2, 4]
0: LDI r1, 2
1: STDOUT r1
2: LDI r8, 2
3: JUMP 5
4: ADD r1,r8
5: STDOUT r1
6: LDI r8, 2
7: ADD r1, r8
8: STDOUT r1

susan3: ~ Print Registers ~ [R0 = 6]
0: LDI r1, 3
1: LDI r4, 9
2: LDI r8, 22
3: LDI r3, 19
4: LDI r2, 2
5: PRINTR

susan3: ~ Print Registers ~ [R0 = 6]
0: LDI r1, 3
1: LDI r4, 9
2: LDI r8, 22
3: LDI r3, 19
4: LDI r2, 2
5: PRINTR

susan4: ~ Print Some Stars ~
0: LDI r1, 2
1: LDI r2, 2
2: ADDV r1, r2
3: ADDV r1, r2
4: ADDV r1,r2
5: ADDV r1,r2

susan5: ~ Blink a heart ~
0: BLINK $heart

susan6: ~ Draw a bird ~
0: DRAW $bird
```
