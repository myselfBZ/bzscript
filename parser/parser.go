package parser

import (
	"fmt"
	"strconv"

	"github.com/myselfBZ/bzscript/ast"
	"github.com/myselfBZ/bzscript/lexer"
	"github.com/myselfBZ/bzscript/token"
)

type Precedence int

const (
	LOWEST      Precedence = iota
	COMPARISION            // <,>,<=,>=,!=,==
	ADD_SUB
	MULTI_DIV
	PREFIX
	CALL
)

var precedences = map[token.TokenType]Precedence{
	token.LPAREN: 		  CALL,
	token.PLUS:           ADD_SUB,
	token.MINUS:          ADD_SUB,
	token.MULTIPLICATION: MULTI_DIV,
	token.DIVISION:       MULTI_DIV,
	token.GT:             COMPARISION,
	token.LT:             COMPARISION,
	token.GTOREQ:         COMPARISION,
	token.LTOREQ:         COMPARISION,
	token.EQ:             COMPARISION,
	token.NOT_EQ:         COMPARISION,
}

type PrefixParseFunc func() ast.Expression
type InfixParseFunc func(left ast.Expression) ast.Expression

type Parser struct {
	lexer    *lexer.Lexer
	curToken *token.Token
	peekTok  *token.Token
	errs     []error

	prefixFunc map[token.TokenType]PrefixParseFunc
	infixFunc  map[token.TokenType]InfixParseFunc
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:      l,
		errs:       make([]error, 0),
		prefixFunc: make(map[token.TokenType]PrefixParseFunc),
		infixFunc:  make(map[token.TokenType]InfixParseFunc),
	}
	p.curToken = p.lexer.NextToken()
	p.peekTok = p.lexer.NextToken()
	p.registerPrefixFunc(token.INT, p.parseIntiger)
	p.registerPrefixFunc(token.FUNCTION, p.parseAnonymousFunc)
	p.registerPrefixFunc(token.LPAREN, p.parseGroupedExpressions)
	p.registerPrefixFunc(token.FLOAT, p.parseFloat)
	p.registerPrefixFunc(token.STRING, p.parseString)
	p.registerPrefixFunc(token.TRUE, p.parseBool)
	p.registerPrefixFunc(token.FALSE, p.parseBool)
	p.registerPrefixFunc(token.IDENT, p.parseIdent)
	p.registerPrefixFunc(token.BANG, p.parsePrefix)
	p.registerPrefixFunc(token.MINUS, p.parsePrefix)

	p.registerInfixFunc(token.PLUS, p.parseInfixExpr)
	p.registerInfixFunc(token.MINUS, p.parseInfixExpr)
	p.registerInfixFunc(token.MULTIPLICATION, p.parseInfixExpr)
	p.registerInfixFunc(token.LPAREN, p.parseFunctionCall)
	p.registerInfixFunc(token.DIVISION, p.parseInfixExpr)
	p.registerInfixFunc(token.EQ, p.parseInfixExpr)
	p.registerInfixFunc(token.NOT_EQ, p.parseInfixExpr)
	p.registerInfixFunc(token.GTOREQ, p.parseInfixExpr)
	p.registerInfixFunc(token.LTOREQ, p.parseInfixExpr)
	p.registerInfixFunc(token.LT, p.parseInfixExpr)
	p.registerInfixFunc(token.GT, p.parseInfixExpr)
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for p.curToken.Type != token.EOF {
		stmnt := p.parse()
		if stmnt != nil {
			program.AddStatement(stmnt)
		}
		p.next()
	}
	return program
}

func (p *Parser) parse() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.FUNCTION:
		return p.parseFunctionLiteral()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	node := &ast.ExpressionStatement{Token: p.curToken}
	node.Expression = p.parseExpression(LOWEST)
	return node
}

func (p *Parser) parseExpression(prec Precedence) ast.Expression {
	fn, ok := p.prefixFunc[p.curToken.Type]
	if !ok {
		p.onError(fmt.Errorf("expected expression, got: %s", p.curToken.Literal))
		return nil
	}
	left := fn()
	for prec < p.peekTokenPrecedence() {
		infixFn, ok := p.infixFunc[p.peekTok.Type]
		if !ok {
			break
		}
		p.next()
		left = infixFn(left)
	}
	return left
}

func (p *Parser) parseString() ast.Expression {
	if p.curToken.Type != token.STRING {
		p.onError(fmt.Errorf("expected string got '%s'", p.curToken.Type))
		return nil
	}

	return &ast.String{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseFloat() ast.Expression {
	node := &ast.Float{Token: p.curToken}
	validFloat, err := strconv.ParseFloat(p.curToken.Literal, 64)

	if err != nil {
		p.onError(fmt.Errorf("expected a valid float, got %s", p.curToken.Literal))
		return nil
	}

	node.Value = validFloat
	return node
}

func (p *Parser) parseBool() ast.Expression {
	node := &ast.Bool{Token: p.curToken}
	val, err := strconv.ParseBool(p.curToken.Literal)
	if err != nil {
		p.onError(fmt.Errorf("expected a valid boolean value, got '%s'", p.curToken.Literal))
		return nil
	}
	node.Value = val
	return node
}

func (p *Parser) parseIntiger() ast.Expression {
	node := &ast.Intiger{Token: p.curToken}
	validInt, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.onError(fmt.Errorf("expected a valid intiger, got %s", p.curToken.Literal))
		return nil
	}
	node.Value = validInt
	return node
}

func (p *Parser) parseIdent() ast.Expression {
	return &ast.Ident{Value: p.curToken.Literal, Token: p.curToken}
}

func (p *Parser) parseAnonymousFunc() ast.Expression {
	function := &ast.AnonymousFuncLiteral{Token: p.curToken}
	if !p.peekTokenIs(token.LPAREN) {
		p.onError(p.expectedError(token.LPAREN, p.peekTok.Type))
		return nil
	}
	p.next()
	function.Params = p.parseParams()
	if !p.peekTokenIs(token.LBRACE) {
		p.onError(p.expectedError(token.LBRACE, p.peekTok.Type))
		return nil
	}
	p.next()
	function.Body = p.parseBlock()
	return function
}

func (p *Parser) parseFunctionCall(function ast.Expression) ast.Expression {
	node := &ast.FunctionCall{Token: p.curToken, Function: function}
	node.Args = p.parseCallArguements()
	return node
}

func (p *Parser) parseCallArguements() []ast.Expression {
	var arguments []ast.Expression
	if p.peekTokenIs(token.RPAREN) {
		p.next()
		return arguments
	}
	p.next()
	arguments = append(arguments, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.next()
		if p.peekTokenIs(token.RPAREN) {
			break
		}
		p.next()
		node := p.parseExpression(LOWEST)
		arguments = append(arguments, node)
	}
	if !p.peekTokenIs(token.RPAREN) {
		p.onError(p.expectedError(token.RPAREN, p.peekTok.Type))
		return nil
	}
	p.next()
	return arguments
}

func (p *Parser) parseFunctionLiteral() ast.Statement {
	node := &ast.ExpressionStatement{Token: p.curToken}
	function := &ast.FunctionLiteral{Token: p.curToken}
	if !p.peekTokenIs(token.IDENT) {
		p.onError(p.expectedError(token.IDENT, p.peekTok.Type))
		return nil
	}
	p.next()
	function.Ident = p.parseIdent().(*ast.Ident)
	if !p.peekTokenIs(token.LPAREN) {
		p.onError(p.expectedError(token.LPAREN, p.peekTok.Type))
		return nil
	}
	p.next()
	function.Params = p.parseParams()
	if !p.peekTokenIs(token.LBRACE) {
		p.onError(p.expectedError(token.LBRACE, p.peekTok.Type))
		return nil
	}
	p.next()
	function.Body = p.parseBlock()
	node.Expression = function
	return node
}

func (p *Parser) parseParams() []*ast.Ident {
	idents := []*ast.Ident{}
	if p.peekTokenIs(token.RPAREN) {
		p.next()
		return idents
	}
	p.next()
	ident1 := &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
	idents = append(idents, ident1)
	for p.peekTokenIs(token.COMMA) {
		p.next()
		if !p.peekTokenIs(token.IDENT) {
			break
		}
		p.next()
		ident := &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
		idents = append(idents, ident)
	}
	if !p.peekTokenIs(token.RPAREN) {
		p.onError(p.expectedError(token.RPAREN, p.peekTok.Type))
		return nil
	}
	p.next()
	return idents
}

func (p *Parser) parseIfStatement() ast.Statement {
	node := &ast.IfStatement{Token: p.curToken}
	p.next()
	condition := p.parseExpression(LOWEST)
	if condition == nil {
		return nil
	}
	node.Condition = condition
	if !p.peekTokenIs(token.LBRACE) {
		p.onError(fmt.Errorf("missing { on if statement block, got: %s", p.peekTok.Literal))
		return nil
	}
	p.next()
	consequence := p.parseBlock()
	if consequence == nil {
		return nil
	}
	node.Consequence = consequence
	if p.peekTokenIs(token.ELSE) {
		p.next()
		if !p.peekTokenIs(token.LBRACE) {
			p.onError(fmt.Errorf("missing { on else statement block"))
			return nil
		}
		p.next()
		alternative := p.parseBlock()
		if alternative == nil {
			return nil
		}
		node.Alternative = alternative
		return node
	}
	return node
}

func (p *Parser) parseBlock() *ast.Block {
	node := &ast.Block{Token: p.curToken}
	p.next()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmnt := p.parse()
		node.Statements = append(node.Statements, stmnt)
		p.next()
	}
	if !p.curTokenIs(token.RBRACE) {
		p.onError(fmt.Errorf("missing } at the end of the block"))
		return nil
	}
	return node
}

func (p *Parser) parseVarStatement() ast.Statement {
	node := &ast.VarStatement{Token: p.curToken}
	if !p.peekTokenIs(token.IDENT) {
		p.onError(p.expectedError(token.IDENT, p.peekTok.Type))
		p.next()
		return nil
	}
	p.next()
	ident := &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
	node.Ident = ident
	p.next()
	if !p.curTokenIs(token.ASSIGN) {
		p.onError(p.expectedError(token.ASSIGN, p.peekTok.Type))
		return nil
	}
	p.next()
	node.Value = p.parseExpression(LOWEST)
	return node
}

func (p *Parser) parseGroupedExpressions() ast.Expression {
	p.next()
	node := p.parseExpression(LOWEST)
	if !p.peekTokenIs(token.RPAREN) {
		p.onError(fmt.Errorf("missing ')' at the end of the expression"))
		return nil
	}
	p.next()
	return node
}

func (p *Parser) parsePrefix() ast.Expression {
	node := &ast.PrefixExpression{Operator: p.curToken.Literal}
	p.next()
	node.Expression = p.parseExpression(PREFIX)
	return node
}

func (p *Parser) parseInfixExpr(left ast.Expression) ast.Expression {
	node := &ast.InfixExpression{Operator: p.curToken.Literal, Left: left}
	prec, ok := precedences[p.curToken.Type]
	if !ok {
		prec = LOWEST
	}
	p.next()
	right := p.parseExpression(prec)
	node.Right = right
	return node
}

func (p *Parser) Errors() []error {
	return p.errs
}

func (p *Parser) onError(err error) {
	p.errs = append(p.errs, err)
}

func (p *Parser) next() {
	p.curToken = p.peekTok
	p.peekTok = p.lexer.NextToken()
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekTok.Type == t
}
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) expectedError(exp token.TokenType, got token.TokenType) error {
	return fmt.Errorf("expected token type %s, got %s", exp, got)
}
func (p *Parser) registerPrefixFunc(t token.TokenType, fn PrefixParseFunc) {
	p.prefixFunc[t] = fn
}
func (p *Parser) registerInfixFunc(t token.TokenType, fn InfixParseFunc) {
	p.infixFunc[t] = fn
}
func (p *Parser) peekTokenPrecedence() Precedence {
	if prec, ok := precedences[p.peekTok.Type]; ok {
		return prec
	}
	return LOWEST
}
