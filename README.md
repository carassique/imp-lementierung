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




# Implementation details
## Syntax
* This implementation of IMP enforces strict "stmt; stmt" sequencing rules. such that: `{print 1;}` would be an invalid program, as would `{print 1; print 2;}`. However `{ print 1 }` and `{ print 1; print 2 }` are valid programs, as `;` delimiter is only inserted between two statements in a block.
* Each program must be wrapped by a block `{ }`.
* Empty blocks are not allowed.
* If-then-else expects all blocks to be present. Only `if { }` is not allowed.

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
The parser receives a stack of tokens, which have already been categorized into the following types:
* Integer values (with integer values already parsed and type-cast)
* Boolean values (with boolean values already parsed and type-cast)
* Variable identifiers (with identifier format already checked)
* Terminal symbols (only valid terminal symbols)

The parser then goes through the stack of `n` tokens in `O(n)` time, by popping or peeking each token and deciding on the possible case without backtracking.

Before implementing the recursive descent parser, some aspects of parser grammar had to be normalized. Specifically, binary infix opeartors presented with left-recursion were split into multiple rules that eliminate left-recursion.
For instance, expression part of the grammar was transformed from this:
```
exp ::= 
     | exp "+" exp           -- Addition
     | exp "*" exp           -- Multiplication
     ...
```
into a set of rules which is similar to this:
```
plus ::= mult plusRhs
plusRhs::= + plus |
mult ::= det multRhs
multRhs ::= * mult |
...
```
This normalization behavior was implemented generically at `parseExpressionGeneric()` in `parser.go`, and resulting operator precedence is defined in `parser_globals.go`:
```
	precedence := [...]*InfixOperator{
		&mult,
		&plus,
		&lessThan,
		&equals,
		&and,
		&or,
	}
```
`parseExpressionGeneric()` tries to parse the left-hand-side of the expression with the operator of highest precedence, otherwise falling back to plain value parsing.

The initial rule for parsing programs is implemented in the function `parseProgram()`. Parser is using helper methods from `TokenizerStream` which implement different checks and operations on the token stack, such as `expectTerminal()` (for popping and verifying terminal tokens) or `peekTokenType()` (for checking token types before deciding on applicable rules).

Each parser function which consumes the `TokenizerStream` corresponds to a rule in the normalized grammar. It returns either a respective part of the AST, or an error.

## Variable scoping
This implementation does not use the simplifying assumption of distinct variables, and implements the general specification for standard variable scoping.
Following example sources are therefore valid:

(Example 1 from task definition)
```
{ 
  x := 1; x := x < 2
}
```
...as variable can be redeclared within the same scope. It means that `x` is going to have numeric value `1` until it is redeclared, where it would change it's type for boolean value `true`.

(Example 2 from task definition)
```
x := false;
while x {
  x := 1
};
x := true
```
...as `x` in the outer scope remains of type boolean, and only within the execution block of `while` does it **locally** create variable `x` of type integer, which hides the outer-scoped `x` for the duration of the block (as there is no redeclaration within this block).

### Closures
The above variable scoping rules are upheld by utilizing the generic interface `Closure[T]` as implemented in `closure.go`.   `Closure[T]` has methods for variable assignment, lookup and declaration. These methods take into account a possible tree-like structure that variable environments can take. On each new execution block, both in the typechecker and the evaluator, a new child variable environment is created by calling `closure.makeChild()`. Child closure always retains references to the parent closure, where it redirects variable lookup calls, if variable names have not been redefined locally. This also ensures, that no local variables leak into the outer scope, as they never overwrite the data in the parent closure.

This interface also provides methods for tracking runtime and typechecking errors, by pushing them onto a global error stack, while retaining references to error origins.

## Evaluator
Evaluator employs short-circuit execution within && and || expressions. Evaluator can be allowed to run disregarding type-checking results, by appending the flag `-i` to the end of `impev` command (see [Running](##Running)).

```
./impev "{ print true ||  1 }" -i
Input: { print true ||  1 }
==== Typecheck error:

============== ERROR STACK TYPE-CHECKER ====================
[Or] Expected [TyBool] but received [TyInt]: '1' in '(true || 1)'
[Print] Ill typed parameter for print
============== ERROR STACK END TYPE-CHECKER =================



Interpeted AST: 
{
    print (true || 1)
}

Output:
true
```


# (Initial requirements)
## Type-checker
## Evaluator
(can use variable bound simplification)
## Parser
(can be skipped)