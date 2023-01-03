# Open questions

## Syntax
- [ ] No { } parenthesis programs allowed or not?
- [ ] Line breaks and extra symbols? Unicode?
- [ ] Tokenizer non-whitespace delimited token recognition buggy (see simple_print_loop.imp)
- [ ] while whileSomething tokenizer
- [ ] precedence of <  
        0 < iterator + 1 < iterator + 5 
        results in
        (0 < ( iterator + 1 < ( iterator + 5 ))
        should be caught by typechecker, as
        int < bool


)

## Semantics
- [ ] check assignment before declaration

### Check
- [ ] Verify language specification
- [ ] (-1) wrapping minus mandatory or not, validate implementation
- [ ] TODO: check if default parsed order of operations matches requirements, see TestExpressionParser
- [ ] Tokenizer "12,4" integer test case -> should return error
- [ ] Check isSkippable for CRLF/LF
- [ ] if {} else {} else {} should not cause panic

- [ ] Implement "not" evaluator

## Architecture
- [ ] Separate ExecutionContext from parser
- [ ] Remove unused code

## Deadlines
- [x] Fetch deadlines
    W15, 17.1-23.1 Projektabgabe und Präsentation

