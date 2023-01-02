### Implemented
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
