# Open questions

## Syntax
- [x] No { } parenthesis programs allowed or not?
        - implemented to spec, -> only allow {} wrapped programs
- [x] Line breaks and extra symbols? Unicode?
- [x] Tokenizer non-whitespace delimited token recognition buggy (see simple_print_loop.imp)
- [x] while whileSomething tokenizer
- [x] precedence of <  
        0 < iterator + 1 < iterator + 5 
        results in
        (0 < ( iterator + 1 < ( iterator + 5 ))
        should be caught by typechecker, as
        int < bool)
        - precedence rules were adjusted to match TypeScript

## Semantics
- [ ] check assignment before declaration
        - should invarably cause error
- [x] == should work for any (equal) type
        - is already the case, unless undefined is checked
- [x] short cirquiting OR and AND seems to be required
### Check
- [ ] Verify language specification
- [x] (-1) wrapping minus mandatory or not, validate implementation
        - it is not, and it works without wrapping
        - notice: prettifier wraps every expression with () to visualize AST
- [x] TODO: check if default parsed order of operations matches requirements, see TestExpressionParser
        - no precedence requirement was given, using the standard precedence as defined for TypeScript
- [x] Tokenizer "12,4" integer test case -> should return error
        - see TestNumbers
- [ ] Check isSkippable for CRLF/LF
- [x] if {} else {} else {} should not cause panic
        - was only causing panic due to test setup running the prettify()

- [x] Implement "not" evaluator
- [x] Stop execution on error

## Architecture
- [x] Separate ExecutionContext from parser
- [ ] Remove unused code

## General-TODOs
- [ ] Check build system, allow running as CLI, compile
- [ ] Add rigid test cases for evaluator and typechecker
- [ ] Request review

## Deadlines
- [x] Fetch deadlines
    W15, 17.1-23.1 Projektabgabe und Präsentation

