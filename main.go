package main

import (
	"bufio"
	"fmt"
	"os"
)

type TokenType string

const (
	// Single-character tokens.
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	COMMA       TokenType = "COMMA"
	DOT         TokenType = "DOT"
	MINUS       TokenType = "MINUS"
	PLUS        TokenType = "PLUS"
	SEMICOLON   TokenType = "SEMICOLON"
	SLASH       TokenType = "SLASH"
	STAR        TokenType = "STAR"

	// One or two character tokens.
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"

	// Literals.
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"

	// Keywords.
	AND    TokenType = "AND"
	CLASS  TokenType = "CLASS"
	ELSE   TokenType = "ELSE"
	FALSE  TokenType = "FALSE"
	FUN    TokenType = "FUN"
	FOR    TokenType = "FOR"
	IF     TokenType = "IF"
	NIL    TokenType = "NIL"
	OR     TokenType = "OR"
	PRINT  TokenType = "PRINT"
	RETURN TokenType = "RETURN"
	SUPER  TokenType = "SUPER"
	THIS   TokenType = "THIS"
	TRUE   TokenType = "TRUE"
	VAR    TokenType = "VAR"
	WHILE  TokenType = "WHILE"

	EOF TokenType = "EOF"
)

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func (t Token) toString() string {
	return fmt.Sprintf("%s %s %s", t.tokenType, t.lexeme, t.literal)
}

type Scanner struct {
	source  []rune
	tokens  []Token
	start   int
	current int
	line    int
}

func makeScanner(source []rune) Scanner {
	return Scanner{
		source: source,
		line:   1,
	}

}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)

}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s *Scanner) advance() rune {
	result := s.source[s.current]
	s.current++

	return result
}

func (s *Scanner) addToken(token TokenType) {
	s.addTokenWithLiteral(token, nil)
}

func (s *Scanner) addTokenWithLiteral(token TokenType, literal any) {
	text := s.source[s.start:s.current]

	s.tokens = append(s.tokens, Token{token, string(text), literal, s.line})
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) matchAssign(expected rune, truthy TokenType, falsy TokenType) TokenType {
	if s.match(expected) {
		return truthy

	} else {
		return falsy
	}
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		// todo: how is this work?
		return '\000'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}

	return s.source[s.current+1]
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	s.addToken(IDENTIFIER)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		// todo: use lox.error
		panic("Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, value)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}
	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	literal := s.source[s.start:s.current]
	// todo: convert to double?
	s.addTokenWithLiteral(NUMBER, literal)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
		break
	case ')':
		s.addToken(RIGHT_PAREN)
		break
	case '{':
		s.addToken(LEFT_BRACE)
		break
	case '}':
		s.addToken(RIGHT_BRACE)
		break
	case ',':
		s.addToken(COMMA)
		break
	case '.':
		s.addToken(DOT)
		break
	case '-':
		s.addToken(MINUS)
		break
	case '+':
		s.addToken(PLUS)
		break
	case ';':
		s.addToken(SEMICOLON)
		break
	case '*':
		s.addToken(STAR)
		break
	case '!':
		s.addToken(s.matchAssign('=', BANG_EQUAL, BANG))
		break
	case '=':
		s.addToken(s.matchAssign('=', EQUAL_EQUAL, EQUAL))
		break
	case '<':
		s.addToken(s.matchAssign('=', LESS_EQUAL, LESS))
		break
	case '>':
		s.addToken(s.matchAssign('=', GREATER_EQUAL, GREATER))
		break
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
		break
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
		break
	case '\n':
		s.line++
		break
	case '"':
		s.string()
		break
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			// todo: call error on Lox
			panic("unexpected token")
		}
	}
}

type Lox struct {
	hadError bool
}

func (l *Lox) runFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Something ain't right", err)
		return
	}
	s := string(file)
	r := []rune(s)
	l.run(r)

	if l.hadError {
		os.Exit(65)
	}
}

func (l *Lox) runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		r := []rune(text)
		l.run(r)
		l.hadError = false
	}
}

func (l *Lox) report(line int, where string, message string) {
	fmt.Printf("[line %s] Error %s : %s\n", line, where, message)
	l.hadError = true
}

func (l *Lox) error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) run(source []rune) {
	scanner := makeScanner(source)
	tokens := scanner.scanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func main() {
	lox := Lox{}

	if len(os.Args) > 2 {
		fmt.Println("Usage: jlox [script]")
	} else if len(os.Args) == 2 {
		lox.runFile(os.Args[1])
	} else {
		lox.runPrompt()
	}
}
