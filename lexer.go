package main

type Token struct {
	Type    string
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNKCIJA"
	LET      = "POSTAVI"
	IF       = "AKO"
	ELSE     = "INACE"
	RETURN   = "VRATI"
	WHILE    = "DOK"
	PRINT    = "ISPISI"
	TRUE     = "TACNO"
	FALSE    = "NETACNO"
)

var keywords = map[string]string{
	"funkcija": FUNCTION,
	"postavi":  LET,
	"ako":      IF,
	"inace":    ELSE,
	"vrati":    RETURN,
	"dok":      WHILE,
	"ispisi":   PRINT,
	"tacno":    TRUE,
	"netacno":  FALSE,
}

func LookupIdent(ident string) string {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
