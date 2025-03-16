package parser

import (
	"ferment/ast"
	"ferment/lexer"
	"fmt"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  lexer.Token
	peekToken lexer.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	pos, tok, _ := p.l.Lex()

	fmt.Printf("nextToken - pos: %+v, token: %v (%s)\n", pos, tok, tok.String())

	p.curToken = p.peekToken
	p.peekToken = tok
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken {
	case lexer.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	p.nextToken()

	if p.curToken != lexer.IDENT {
		p.errors = append(p.errors, "expected indetifier after POSTAVI")
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.String(),
	}

	p.nextToken()

	if p.curToken != lexer.ASSIGN {
		p.errors = append(p.errors, "expected '=' after identifier")
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression()

	if p.peekToken == lexer.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	switch p.curToken {
	case lexer.IDENT:
		return &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.String(),
		}
	case lexer.INT:
		return &ast.IntegerLiteral{
			Token: p.curToken,
			Value: p.curToken.String(),
		}
	}

	return nil
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken != lexer.EOF {
		fmt.Printf("ParseProgram loop - curToken: %v (%s)\n", p.curToken, p.curToken.String())

		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		} else {
			fmt.Println("nil statement returned")
		}

		p.nextToken()
	}

	return program
}
