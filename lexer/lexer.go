package lexer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

type Token int

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func (p Position) Line() int {
	return p.line
}

func (p Position) Column() int {
	return p.column
}

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT
	STRING

	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH

	EQ
	NOT_EQ

	LT
	GT

	COMMA
	SEMICOLON

	LPAREN
	RPAREN
	LBRACE
	RBRACE

	FUNCTION
	LET
	IF
	ELSE
	RETURN
	WHILE
	PRINT
	TRUE
	FALSE
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	INT:     "INT",

	ASSIGN:   "=",
	PLUS:     "+",
	MINUS:    "-",
	BANG:     "!",
	ASTERISK: "*",
	SLASH:    "/",

	EQ:     "==",
	NOT_EQ: "!=",

	LT: "<",
	GT: ">",

	COMMA:     ",",
	SEMICOLON: ";",

	LPAREN: "(",
	RPAREN: ")",
	LBRACE: "{",
	RBRACE: "}",

	FUNCTION: "FUNKCIJA",
	LET:      "POSTAVI",
	IF:       "AKO",
	ELSE:     "INACE",
	RETURN:   "VRATI",
	WHILE:    "DOK",
	PRINT:    "ISPISI",
	TRUE:     "TACNO",
	FALSE:    "NETACNO",
}

var keywords = map[string]Token{
	"POSTAVI":  LET,
	"VRATI":    RETURN,
	"FUNKCIJA": FUNCTION,
	"AKO":      IF,
	"INACE":    ELSE,
	"DOK":      WHILE,
	"ISPISI":   PRINT,
	"TACNO":    TRUE,
	"NETACNO":  FALSE,
}

func LookupIdent(ident string) Token {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (t Token) String() string {
	return tokens[t]
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}
			panic(err)
		}
		l.pos.column++

		switch r {
		case ';':
			return l.pos, SEMICOLON, ";"
		case '+':
			return l.pos, PLUS, "+"
		case '-':
			return l.pos, MINUS, "-"
		case '*':
			return l.pos, ASTERISK, "*"
		case '/':
			return l.pos, SLASH, "/"
		case '=':
			nr, err := l.peek()
			if err == nil && nr == '=' {
				l.reader.ReadRune()
				l.pos.column++
				return l.pos, EQ, "=="
			}
			return l.pos, ASSIGN, "="
		case '(':
			return l.pos, LPAREN, "("
		case ')':
			return l.pos, RPAREN, ")"
		case '{':
			return l.pos, LBRACE, "{"
		case '}':
			return l.pos, RBRACE, "}"
		case '<':
			return l.pos, LT, "<"
		case '>':
			return l.pos, GT, ">"
		case '!':
			nr, err := l.peek()
			if err == nil && nr == '=' {
				l.reader.ReadRune()
				l.pos.column++
				return l.pos, NOT_EQ, "!="
			}
			return l.pos, BANG, "!"
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return startPos, INT, lit
			} else if unicode.IsLetter(r) {
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()

				tok := LookupIdent(lit)

				return startPos, tok, lit
			} else {
				return l.pos, ILLEGAL, string(r)
			}
		}
	}
}

func (l *Lexer) peek() (rune, error) {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		return 0, nil
	}

	if err := l.reader.UnreadRune(); err != nil {
		return 0, nil
	}

	return r, nil
}

func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}

		l.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			l.backup()
			return lit
		}
	}
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

func (l *Lexer) resetPostion() {
	l.pos.line++
	l.pos.column = 0
}

func TestLexer() {
	file, err := os.Open("input.test")
	if err != nil {
		panic(err)
	}

	lexer := NewLexer(file)
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}

		fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.column, tok, lit)
	}
}
