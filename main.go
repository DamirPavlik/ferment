package main

import (
	"ferment/ast"
	"ferment/lexer"
	"ferment/parser"
	"fmt"
	"strings"
)

func main() {
	reader := strings.NewReader(`
POSTAVI x = 5;
VRATI x + y;
AKO (x > y) { ISPISI "x je veće"; } INACE { ISPISI "y je veće"; }
`)

	lex := lexer.NewLexer(reader)
	p := parser.NewParser(lex)

	program := p.ParseProgram()

	fmt.Println("Parsed AST:")
	for _, stmt := range program.Statements {
		switch stmt := stmt.(type) {
		case *ast.LetStatement:
			fmt.Printf("LetStatement: Name=%s, Value=%v\n", stmt.Name.Value, stmt.Value)
		case *ast.ReturnStatement:
			fmt.Printf("ReturnStatement: Value=%v\n", stmt.Value)
		case *ast.IfStatement:
			fmt.Printf("IfStatement: Condition=%v, Consequence=%v, Alternative=%v\n", stmt.Condition, stmt.Consequence, stmt.Alternative)
		}
	}

	if len(p.Errors()) > 0 {
		fmt.Println("\nParser Errors:")
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
	}
}
