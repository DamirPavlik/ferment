package lexer

import (
	"bufio"
	"io"
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
	STRING:  "STRING",

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
		case '\n':
			l.resetPostion()
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
			return l.pos, EQ, "="
		default:
			if unicode.IsSpace(r) {
				continue
			}
		}
	}
}

func (l *Lexer) resetPostion() {
	l.pos.line++
	l.pos.column = 0
}
