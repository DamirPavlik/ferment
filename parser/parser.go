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
	_, tok, _ := p.l.Lex()

	p.curToken = p.peekToken
	p.peekToken = tok
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.PRINT:
		return p.parsePrintStatement()
	default:
		p.nextToken()
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.Value = p.parseExpression()

	if p.peekToken == lexer.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}

	if p.peekToken != lexer.LPAREN {
		p.errors = append(p.errors, "expected '(' after AKO")
		return nil
	}

	p.nextToken()
	p.nextToken()

	stmt.Condition = p.parseExpression()

	if p.peekToken != lexer.RPAREN {
		p.errors = append(p.errors, "expected ')' after condition")
		return nil
	}

	p.nextToken()
	p.nextToken()

	if p.curToken != lexer.LBRACE {
		p.errors = append(p.errors, "expected '{' after condition")
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekToken == lexer.ELSE {
		p.nextToken()
		p.nextToken()
		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	stmt := &ast.PrintStatement{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression()

	if p.peekToken == lexer.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for p.curToken != lexer.RBRACE && p.curToken != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression()

	if p.curToken != lexer.LBRACE {
		p.errors = append(p.errors, "expected '{' after condition")
		return nil
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	left := p.parsePrimary()

	for p.peekToken == lexer.PLUS || p.peekToken == lexer.MINUS || p.peekToken == lexer.ASTERISK || p.peekToken == lexer.SLASH || p.peekToken == lexer.GT || p.peekToken == lexer.LT {
		p.nextToken()
		operator := p.curToken

		p.nextToken()
		right := p.parsePrimary()

		left = &ast.InfixExpression{
			Token:    operator,
			Left:     left,
			Operator: operator.String(),
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parsePrimary() ast.Expression {
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
	case lexer.STRING:
		return &ast.StringLiteral{
			Token: p.curToken,
			Value: p.curToken.String(),
		}
	case lexer.LPAREN:
		p.nextToken()
		exp := p.parseExpression()

		if p.curToken != lexer.RPAREN {
			p.errors = append(p.errors, "expected closing )")
		}

		return exp
	}

	return nil
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken != lexer.EOF {
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
