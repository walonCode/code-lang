package parser

import (
	"fmt"
	"strconv"

	"github.com/walonCode/code-lang/ast"
	"github.com/walonCode/code-lang/lexer"
	"github.com/walonCode/code-lang/token"
)

const (
	_ int = iota
	LOWEST 
	EQUALS // ==
	LESSGREATER // > or <
	SUM // +
	PRODUCT // *
	PREFIX // -x or !x
	CALL //myfunction(x)
)

var precendeces = map[token.TokenType]int {
	token.EQ: EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT: LESSGREATER,
	token.GT: LESSGREATER,
	token.PLUS: SUM,
	token.MINUS: SUM,
	token.SLASH: PRODUCT,
	token.ASTERISK: PRODUCT,
}

func(p *Parser)peekPredences()int{
	if p, ok := precendeces[p.peekToken.Type]; ok {
		return p
	}
	
	return LOWEST
}

func(p *Parser)curPrecendence()int{
	if p, ok := precendeces[p.curToken.Type];ok {
		return p
	}
	
	return LOWEST
}

//Pratt Parser
type (
	prefixParseFn func()ast.Expression
	infixParseFn func(ast.Expression)ast.Expression
)


type Parser struct {
	l *lexer.Lexer
	
	errors []string
	
	curToken token.Token
	peekToken token.Token
	
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns map[token.TokenType]infixParseFn
}

func(p *Parser)registerPrefix(tokenType token.TokenType, fn prefixParseFn){
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser)registerInfix(tokenType token.TokenType, fn infixParseFn){
	p.infixParseFns[tokenType] = fn
}

func (p *Parser)nextToken(){
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func(p *Parser)parseStatement()ast.Statement{
	switch p.curToken.Type{
		case token.LET:
			return p.parseLetStatement()
		case token.RETURN:
			return p.parseReturnStatement()
		default:
			return p.parseExpressionStatement()
	}
}

func (p *Parser)parseExpressionStatement()*ast.ExpressionStatement{
	defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	
	if p.peekTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	
	return stmt
}

func (p *Parser)parseIntergerLiteral()ast.Expression {
	defer untrace(trace("parseIntergerLiteral"))
	il := &ast.IntergerLiteral{Token:p.curToken}
	
	value, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	
	il.Value = int64(value)
	
	return il
}

func(p *Parser)noPrefixParseError(t token.TokenType){
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func(p *Parser)parseExpression(predence int)ast.Expression{
	defer untrace(trace("parseExpression"))
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseError(p.curToken.Type)
		return nil
	}
	
	leftExp := prefix()
	
	for !p.peekTokenIs(token.SEMICOLON) && predence < p.peekPredences(){
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		
		p.nextToken()
		
		leftExp = infix(leftExp)
	}
	
	return leftExp
}

func (p *Parser)parseReturnStatement() *ast.ReturnStatement{
	stmt := &ast.ReturnStatement{Token: p.curToken}
	
	p.nextToken()
	for !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	
	return stmt
}

func (p *Parser)parseLetStatement()*ast.LetStatement{
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT){
		return nil
	}
	
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	
	if !p.expectPeek(token.ASSIGN){
		return nil
	}
	
	for !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	
	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser)ParsePrograme()*ast.Program{
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	
	return program
}

func(p *Parser)Errors()[]string{
	return p.errors
}

func(p *Parser)peekError(t token.TokenType){
	msg := fmt.Sprintf("expect next token to be %s, got %s instead", 
	t, p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}

func New(l *lexer.Lexer)*Parser {
	p := &Parser{
		l:l,
		errors: []string{},
	}
	
	p.nextToken()
	p.nextToken()
	
	//prefix
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntergerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	
	//infix
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	return p
}

func(p *Parser)parseInfixExpression(left ast.Expression)ast.Expression{
	defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
		Left: left,
	}
	
	predence := p.curPrecendence()
	p.nextToken()
	expression.Right = p.parseExpression(predence)
	
	return expression
}

func(p *Parser)parsePrefixExpression()ast.Expression{
	defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
	}
	
	p.nextToken()
	
	expression.Right = p.parseExpression(PREFIX)
	
	return expression
}

func (p *Parser)parseIdentifier()ast.Expression{
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}