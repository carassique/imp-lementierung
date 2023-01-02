## Readiness status by rule
### Implemented
- [X] vars       Variable names, start with lower-case letter
* Check format

- [x] prog      ::= block
- [x] block     ::= "{" statement "}"
- [x] statement ::=  statement ";" statement           -- Command sequence
- [x]            |  vars ":=" exp                     -- Variable declaration
- [x]            |  vars "=" exp                      -- Variable assignment
- [x]            |  "while" exp block                 -- While
- [x]            |  "if" exp block "else" block       -- If-then-else
- [x]            |  "print" exp                       -- Print

- [x] exp ::= 0 | 1 | -1 | ...     -- Integers
- [x]     | "true" | "false"      -- Booleans
- [x]     | exp "+" exp           -- Addition
- [x]     | exp "*" exp           -- Multiplication
- [x]     | exp "||" exp          -- Disjunction
- [x]     | exp "&&" exp          -- Conjunction
- [x]     | "!" exp               -- Negation
- [x]     | exp "==" exp          -- Equality test
- [x]     | exp "<" exp           -- Lesser test
- [x]     | "(" exp ")"           -- Grouping of expressions
- [x]     | vars                  -- Variables

### Tested
- [ ] vars       Variable names, start with lower-case letter

- [ ] prog      ::= block
- [ ] block     ::= "{" statement "}"
- [x] statement ::=  statement ";" statement           -- Command sequence
- [x]            |  vars ":=" exp                     -- Variable declaration
- [x]            |  vars "=" exp                      -- Variable assignment
- [x]            |  "while" exp block                 -- While
- [x]            |  "if" exp block "else" block       -- If-then-else
- [x]            |  "print" exp                       -- Print

- [x] exp ::= 0 | 1 | -1 | ...     -- Integers
- [x]     | "true" | "false"      -- Booleans
- [x]     | exp "+" exp           -- Addition
- [x]     | exp "*" exp           -- Multiplication
- [x]     | exp "||" exp          -- Disjunction
- [x]     | exp "&&" exp          -- Conjunction
- [x]     | "!" exp               -- Negation
- [x]     | exp "==" exp          -- Equality test
- [x]     | exp "<" exp           -- Lesser test
- [x]     | "(" exp ")"           -- Grouping of expressions
- [x]     | vars                  -- Variables
