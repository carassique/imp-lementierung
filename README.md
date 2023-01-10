# imp-lementierung - IMP Language Implementation

## Template and description from
https://sulzmann.github.io/ModelBasedSW/imp.html
## Task definition
https://sulzmann.github.io/ModelBasedSW/notesWiSe22-23.html#(7)

# How to build and run demonstrator utility
## Requirements
* Go 1.19 installation
* If line breaks are used in the  IMP source code, they are expected to be of "newline" type only, not "CRLF". Only tested under Linux.
## Building
Navigate into /src/main and run `go build impev.go` command. Executable will be generated in the same folder.
## Running
Run the executable via command line (e.g. `./impev`). Following options can be used:
* Get help about available flags and parameters: `./impev help`
* Parse, type-check and execute IMP code inline: `./impev "{ print 1234 }"`
    * Notice: this style of execution may conflict with BASH (or other) CLI commands
    * Use flag `-i` to ignore type-checking results and execute the program regardless e.g.:
    `./impev "{ print false && 1 }" -i`
* Parse, type-check and execute IMP code from file: `./impev -f filename`
    * Relative path can be used `./impev -f "src/filename.imp"`
    * Use flag `-i` to ignore type-checking results and execute the program regardless e.g.:
    `./impev -f "src/filename.imp" -i`


## Sample output
### Valid program
```
./impev "{test:=1;print test}"
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


# Implementation details
## Quality assurance
Language implementation makes use of the Go standard testing library https://pkg.go.dev/testing and
includes automated unit and integration tests, which can be run via `go test imp -v` from the root directory.
External test runner can also be used to visualize test completion, such as the Visual Studio Code test runner.
### Test structure
Tests are located under `/src/imp/`
* Evaluator partial integration tests: `evaluator_test.go`
    * Utilizes parser implementation to generate ASTs for testing. Checks for valid output or expected errors.
    * Also runs type checks
* Typechecker unit tests: `typechecker_test.go`
    * Verifies whether (or not) and where any particular AST under test produces errors. Utilizes error stack to verify source of type errors.
* Parser unit tests: `parser_test.go`
    * Assumes generated tokens, does not test tokenization
* Parser integration tests: `parser_src_test.go`
    * Tests parser's integration with lexer, by comparing ASTs generated from source code
* Lexer unit tests: `tokenizer_test.go`
    * Checks generated token lists against expected tokens
* Overall integration tests: `integration_test.go`
    * Performs Tests via IMP source code files from folders `/src/imp/test_source/should_fail/` and `/src/imp/test_source/should_pass`:
        * Tokenizes, parses, type-checks and evaluates each file from the respective folder.
        * If the file resides in the `should_pass` folder, then the test passes only if the file has passed all of the above stages without errors.
        * If the file resides in the `should_fail` folder, then the test passes only if the file causes error on some of the above stages.



## Parser
The parser is a recurisve descent parser. Overall AST generation process comprises of two stages: **tokenization** and **parsing**.
### Tokenization
Tokenization is performed by a tokenizer (lexer), which is implemented in `tokenizer.go`.
Operating principle of the tokenizer is as follows:
1. Input string is accumulated character-by-character into a testable token prefix `tokenCandidate`.
2. After each additional character is appended to the end of the `tokenCandidate`, the `tokenCandidate` is checked by prefix matching functions of form `(prefix string) -> bool`. If any such function returns true, it indicates that the token prefix can be a part of a valid token, so accumulation continues with the next step.
3. Once `tokenCandidate` is no longer matched by any prefix matcher, then it indicates, that the possible token must be at least one character shorter.
4. `tokenCandidate` is shortened by one character, and token creation is attempted, following the order of precedence defined in `allMatchers()`.
    * Order of precedence is important to not interpret e.g. `true` as a variable with identifier `true`, but rather as a boolean literal.
5. Since the only ambiguous case involves `=` and `==` tokens, one-character-backtrack is sufficient to distinguish all types of tokens.
### Parsing

## Type-checker
## Evaluator

## Type-checker
## Evaluator
(can use variable bound simplification)
## Parser
(can be skipped)