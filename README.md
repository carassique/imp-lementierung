# imp-lementierung
IMP Language Implementation

# How to build and run demonstrator utility
## Requirements
* Go 1.19 installation
* If line breaks are used in the  IMP source code, they are expected to be of "newline" type only, not "CRLF". Only tested under Linux.
## Building
Navigate into /src/main and run `go build impev.go` command. Executable will be generated in the same folder.
## Running
Run the executable via command line (e.g. `./impev`). Following options can be used:
* Parse, type-check and execute IMP code inline: `./impev "{ print 1234 }"`
    * Notice: this style of execution may conflict with BASH (or other) CLI commands
* Parse, type-check and execute IMP code from file: `./impev -f filename`
    * Relative path can be used `./impev -f "src/filename.imp"`

## Sample output
### Valid program
```
src/main$ ./impev "{test:=1;print test}"
Input: {test:=1;print test}


Interpeted AST: 
{
    test := 1;
    print test
}

Output:
1
```
### Invalid program
```
./impev "{test=1; print test}"
Input: {test=1; print test}
==== Typecheck error:

============== ERROR STACK TYPE-CHECKER ====================
[Assign] Variable "test" does not exist in this scope
[Seq] First statement of the sequence did not pass type checking
============== ERROR STACK END TYPE-CHECKER =================
```

**This implementation of IMP enforces strict "stmt;stmt" sequencing rules. such that: `{print 1;}` would be an invalid program, as would `{print 1; print 2;}`.**


# Template and description from
https://sulzmann.github.io/ModelBasedSW/imp.html

# Task definition
https://sulzmann.github.io/ModelBasedSW/notesWiSe22-23.html#(7)

## Type-checker
## Evaluator
(can use variable bound simplification)
## Parser
(can be skipped)