### Implemented
- [x] vars       Variable names, start with lower-case letter

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
- [ ] statement ::=  statement ";" statement           -- Command sequence
- [ ]            |  vars ":=" exp                     -- Variable declaration
- [ ]            |  vars "=" exp                      -- Variable assignment
- [ ]            |  "while" exp block                 -- While
- [ ]            |  "if" exp block "else" block       -- If-then-else
- [ ]            |  "print" exp                       -- Print

- [ ] exp ::= 0 | 1 | -1 | ...     -- Integers
- [ ]     | "true" | "false"      -- Booleans
- [ ]     | exp "+" exp           -- Addition
- [ ]     | exp "*" exp           -- Multiplication
- [ ]     | exp "||" exp          -- Disjunction
- [ ]     | exp "&&" exp          -- Conjunction
- [ ]     | "!" exp               -- Negation
- [ ]     | exp "==" exp          -- Equality test
- [ ]     | exp "<" exp           -- Lesser test
- [ ]     | "(" exp ")"           -- Grouping of expressions
- [ ]     | vars                  -- Variables
