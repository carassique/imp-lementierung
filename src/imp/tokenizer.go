package imp

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var terminalTokensPriority = []string{
	OPEN_BLOCK_GROUPING, CLOSE_BLOCK_GROUPING,
	PRINT, WHILE, IF, ELSE, STATEMENT_DELIMITER,
	DECLARATION, ASSIGNMENT, ADD, MULTIPLY,
	OR, AND, NOT, EQUALS, LESS_THAN, OPEN_EXPRESSION_GROUPING,
	CLOSE_EXPRESSION_GROUPING,
}

var terminalTokens = toSet(terminalTokensPriority)

type TokenCandidatePrefixMatcher func(tokenCandidate string) bool

type TokenCandidatePrefixMatchers []TokenCandidatePrefixMatcher

type TokenCandidateMatcher func(tokenCandidate string) (bool, Token)

type TokenCandidateMatchers []TokenCandidateMatcher

type TokenizerResultData []Token

type StringSet map[string]struct{}

func toSet(tokens []string) StringSet {
	tokenSet := make(map[string]struct{})
	for _, token := range tokens {
		tokenSet[token] = struct{}{}
	}
	return tokenSet
}

func isInteger(tokenCandidate string) (bool, Token) {
	// TODO: consider returning parsed value as a tuple or struct
	if value, err := strconv.Atoi(tokenCandidate); err == nil {
		return true, Token{token: tokenCandidate, tokenType: IntegerValue, integerValue: value}
	}
	return false, noMatch(tokenCandidate)
}

func isBoolean(tokenCandidate string) (bool, Token) {
	if tokenCandidate == BOOLEAN_TRUE {
		return true, Token{token: tokenCandidate, tokenType: BooleanValue, booleanValue: true}
	}
	if tokenCandidate == BOOLEAN_FALSE {
		return true, Token{token: tokenCandidate, tokenType: BooleanValue, booleanValue: false}
	}
	return false, noMatch(tokenCandidate)
}

func isVariableName(tokenCandidate string) (bool, Token) {
	// TODO: implement variable format
	match, _ := regexp.MatchString("^[a-z]([A-Za-z]|[0-9])*$", tokenCandidate)
	if match {
		return true, Token{token: tokenCandidate, tokenType: VariableName}
	}
	return false, noMatch(tokenCandidate)
}

func isValue(tokenCandidate string) (bool, Token) {
	isInteger, integerToken := isInteger(tokenCandidate)
	if isInteger {
		return true, integerToken
	}
	isBoolean, booleanToken := isBoolean(tokenCandidate)
	if isBoolean {
		return true, booleanToken
	}
	// Order matters
	isVariableName, variableNameToken := isVariableName(tokenCandidate)
	if isVariableName {
		return true, variableNameToken
	}
	return false, Token{}
}

func isAmbiguous(tokenCandidate string) bool {
	return tokenCandidate == ASSIGNMENT
}

func isTerminal(tokenCandidate string, terminalTokens StringSet) bool {
	if _, ok := terminalTokens[tokenCandidate]; ok {
		return ok
	}
	return false
}

func matchCandidateToken(tokenCandidate string, matchers TokenCandidateMatchers) (bool, TokenCandidateMatchers) {
	candidateMatchers := make([]TokenCandidateMatcher, 0)
	if len(tokenCandidate) == 0 {
		return false, candidateMatchers
	}
	hasMatches := false
	for _, matcher := range matchers {
		matched, _ := matcher(tokenCandidate)
		if matched {
			hasMatches = true
			candidateMatchers = append(candidateMatchers, matcher)
		}
	}
	return hasMatches, candidateMatchers
}

func makeToken(tokenCandidate string, matchers TokenCandidateMatchers) (bool, Token) {
	if len(tokenCandidate) == 0 || len(matchers) == 0 {
		return false, noMatch("")
	}

	for _, matcher := range matchers {
		ok, token := matcher(tokenCandidate) //TODO: consider match preference
		if ok && token.token == tokenCandidate {
			return ok, token
		}
	}
	return false, noMatch(tokenCandidate)
}

func terminalPrefixMatcher(tokenCandidate string) bool {
	for _, token := range terminalTokensPriority {
		return strings.HasPrefix(token, tokenCandidate)
	}
	return false
}

func terminalPrefixMatcher2(tokenCandidate string) (bool, Token) {
	for _, token := range terminalTokensPriority {
		if strings.HasPrefix(token, tokenCandidate) {
			return true, terminal(token)
		}
	}
	return false, noMatch(tokenCandidate)
}

func integerPrefixMatcher(tokenCandidate string) (bool, Token) {
	if tokenCandidate == "-" {
		// Account for negative number prefix
		return true, integer(-0)
	}
	return isInteger(tokenCandidate)
}

func variablePrefixMatcher(tokenCandidate string) (bool, Token) {
	return isVariableName(tokenCandidate)
}

func booleanPrefixMatcher(tokenCandidate string) (bool, Token) {
	ok, token := isBoolean(tokenCandidate)
	if !ok {
		prefixTrue := strings.HasPrefix(BOOLEAN_TRUE, tokenCandidate)
		prefixFalse := strings.HasPrefix(BOOLEAN_FALSE, tokenCandidate)

		if prefixTrue || prefixFalse {
			return true, booleanToken(prefixTrue)
		}
	}
	return ok, token
}

func terminalMatcher(tokenCandidate string) (bool, Token) {
	// What is actually required is a prefix (possibility) matcher, not identity matcher
	if isTerminal(tokenCandidate, terminalTokens) {
		return true, terminal(tokenCandidate)
	}
	return false, noMatch(tokenCandidate)
}

func allPrefixMatchers() TokenCandidatePrefixMatchers {
	return TokenCandidatePrefixMatchers{
		terminalPrefixMatcher,
	}
}

func allMatchers() TokenCandidateMatchers {
	matchers := TokenCandidateMatchers{
		terminalPrefixMatcher2,
		integerPrefixMatcher,
		booleanPrefixMatcher, //Booleans should match before variable names
		variablePrefixMatcher,
	}
	return matchers
}

func anyPrefixMatches(word string, matchers ...TokenCandidatePrefixMatcher) bool {
	for _, matcher := range matchers {
		if matcher(word) {
			return true
		}
	}
	return false
}

func anyFullMatches(word string, matchers ...TokenCandidateMatcher) (bool, Token) {
	for _, matcher := range matchers {
		ok, token := matcher(word)
		if ok {
			return ok, token
		}
	}
	return false, noMatch(word)
}

func tokenize(sourcecode string) TokenizerResultData {
	tokenList := make([]Token, 0)
	candidateMatchers := allMatchers()
	//fullMatchers := allMatchers()
	//prefixMatchers := allPrefixMatchers()
	paddedSourceCode := sourcecode + " " // To match last character easier
	tokenCandidate := ""
	currentToken := ""
	for _, character := range paddedSourceCode {
		currentToken += (string)(character)
		hasMatch, newCandidateMatchers := matchCandidateToken(currentToken, candidateMatchers)

		// if anyPrefixMatches(currentToken, prefixMatchers...) {
		// 	// consider token still possible

		// } else {
		// 	// evaluate results

		// 	// has full matches? -> jump back to the longest
		// 	// has no full matches -> skip one character, try again
		// }

		// if ok, _ := anyFullMatches(currentToken, fullMatchers...); ok {

		// }
		if !hasMatch {
			// TODO: make full match
			success, nextToken := makeToken(tokenCandidate, candidateMatchers)
			candidateMatchers = allMatchers()
			if success {
				tokenList = append(tokenList, nextToken)
			} else {
				// TODO error

			}
			// Rematch with this character as first
			tokenCandidate = ""
			currentToken = (string)(character)
			hasMatch, newCandidateMatchers = matchCandidateToken(currentToken, candidateMatchers)
			if !hasMatch {
				// Single-character word does not have matching tokens -> skip
				currentToken = ""
				tokenCandidate = ""
				candidateMatchers = allMatchers()
			} else {
				candidateMatchers = newCandidateMatchers
				tokenCandidate = currentToken
			}
		} else {
			tokenCandidate = currentToken
			candidateMatchers = newCandidateMatchers
		}
	}
	return tokenList
}

func tokenize2(sourceCode string, terminalTokens StringSet) TokenizerResultData {
	// TODO: simplify code, use scanner or generic lexer
	// TODO: write tests and fix issues
	tokenList := make([]Token, 0)

	currentToken := ""
	tokenCandidate := ""
	for _, character := range sourceCode {
		if unicode.IsSpace(character) {
			if len(currentToken) > 0 {
				if len(tokenCandidate) > 0 {
					tokenList = append(tokenList, Token{
						tokenType: Terminal,
						token:     currentToken,
					})
				} else {
					//TODO: is non-terminal?
					if constitutesValue, value := isValue(currentToken); constitutesValue {
						tokenList = append(tokenList, value)
					}

					//TODO: error - no token recognized!
				}
				tokenCandidate = ""
				currentToken = ""
			}
		} else {
			// Ignore spaces between tokens
			currentToken += (string)(character)
			if isTerminal(currentToken, terminalTokens) {
				tokenCandidate = currentToken
				if !isAmbiguous(currentToken) {
					tokenList = append(tokenList, Token{
						tokenType: Terminal,
						token:     currentToken,
					})
					currentToken = ""
					tokenCandidate = ""
				}
			} else {
				if len(tokenCandidate) > 0 {
					tokenList = append(tokenList, Token{
						tokenType: Terminal,
						token:     tokenCandidate,
					})
					currentToken = (string)(character)
				}
				if isTerminal(currentToken, terminalTokens) {
					tokenCandidate = currentToken
					if !isAmbiguous(currentToken) {
						tokenList = append(tokenList, Token{
							tokenType: Terminal,
							token:     currentToken,
						})
						currentToken = ""
						tokenCandidate = ""
					}
				}
			}
		}

	}
	// Anything remains after the last character, it should be matched
	if len(tokenCandidate) > 0 {
		tokenList = append(tokenList, Token{
			tokenType: Terminal,
			token:     tokenCandidate,
		})
	}

	return tokenList
}
