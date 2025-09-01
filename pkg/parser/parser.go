package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GrandpaEJ/go-script/pkg/ast"
	"github.com/GrandpaEJ/go-script/pkg/lexer"
	"github.com/GrandpaEJ/go-script/pkg/stdlib"
)

// Parser represents the parser
type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token

	errors []string

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Precedence levels
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[lexer.TokenType]int{
	lexer.EQ:       EQUALS,
	lexer.NOT_EQ:   EQUALS,
	lexer.LT:       LESSGREATER,
	lexer.GT:       LESSGREATER,
	lexer.LT_EQ:    LESSGREATER,
	lexer.GT_EQ:    LESSGREATER,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.DIVIDE:   PRODUCT,
	lexer.MULTIPLY: PRODUCT,
	lexer.MODULO:   PRODUCT,
	lexer.POWER:    PRODUCT,
	lexer.LPAREN:   CALL,
	lexer.LBRACKET: INDEX,
	lexer.DOT:      INDEX,
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.NIL, p.parseNilLiteral)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.NOT, p.parsePrefixExpression)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(lexer.LBRACE, p.parseMapLiteral)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.DIVIDE, p.parseInfixExpression)
	p.registerInfix(lexer.MULTIPLY, p.parseInfixExpression)
	p.registerInfix(lexer.MODULO, p.parseInfixExpression)
	p.registerInfix(lexer.POWER, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.GT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.AND, p.parseInfixExpression)
	p.registerInfix(lexer.OR, p.parseInfixExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.LBRACKET, p.parseIndexExpression)
	p.registerInfix(lexer.DOT, p.parseSelectorExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	// Skip comments
	for p.peekToken.Type == lexer.COMMENT {
		p.peekToken = p.l.NextToken()
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		lexer.TokenTypeString(t), lexer.TokenTypeString(p.peekToken.Type))
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", lexer.TokenTypeString(t))
	p.errors = append(p.errors, msg)
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// Parse package declaration
	if p.curTokenIs(lexer.PACKAGE) {
		p.nextToken()
		if p.curTokenIs(lexer.IDENT) {
			program.Package = p.curToken.Literal
			p.nextToken()
		}
	} else {
		program.Package = "main" // default package
	}

	// Skip newlines after package declaration
	for p.curTokenIs(lexer.NEWLINE) {
		p.nextToken()
	}

	// Parse imports
	for p.curTokenIs(lexer.IMPORT) || p.curTokenIs(lexer.FROM) {
		importDecl := p.parseImportDeclaration()
		if importDecl != nil {
			program.Imports = append(program.Imports, importDecl)
		}
		// Skip newlines after imports
		for p.curTokenIs(lexer.NEWLINE) {
			p.nextToken()
		}
	}

	// Parse statements
	for !p.curTokenIs(lexer.EOF) {
		// Skip newlines and indentation tokens at the top level
		if p.curTokenIs(lexer.NEWLINE) || p.curTokenIs(lexer.INDENT) || p.curTokenIs(lexer.DEDENT) {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		// Always advance to the next token unless we're at EOF
		if !p.curTokenIs(lexer.EOF) {
			p.nextToken()
		}
	}

	return program
}

func (p *Parser) parseImportDeclaration() *ast.ImportDecl {
	importDecl := &ast.ImportDecl{}

	if p.curTokenIs(lexer.FROM) {
		// from "path" import item1, item2
		p.nextToken()
		if !p.curTokenIs(lexer.STRING) {
			return nil
		}
		importDecl.Path = p.curToken.Literal
		p.nextToken()

		if !p.expectPeek(lexer.IMPORT) {
			return nil
		}
		p.nextToken()

		// Parse import items
		for {
			if p.curTokenIs(lexer.IDENT) {
				importDecl.Items = append(importDecl.Items, p.curToken.Literal)
			}
			if !p.peekTokenIs(lexer.COMMA) {
				break
			}
			p.nextToken() // consume comma
			p.nextToken() // move to next item
		}
	} else if p.curTokenIs(lexer.IMPORT) {
		p.nextToken()

		// Handle different import formats
		if p.curTokenIs(lexer.LPAREN) {
			// import ("os", "fmt", "time")
			p.nextToken()
			for !p.curTokenIs(lexer.RPAREN) && !p.curTokenIs(lexer.EOF) {
				if p.curTokenIs(lexer.STRING) {
					// Remove quotes and resolve alias
					rawPath := strings.Trim(p.curToken.Literal, `"`)
					resolvedPath := `"` + stdlib.GetRealPackagePath(rawPath) + `"`
					importDecl.Items = append(importDecl.Items, resolvedPath)
				}
				p.nextToken()
				if p.curTokenIs(lexer.COMMA) {
					p.nextToken()
				}
			}
		} else if p.curTokenIs(lexer.STRING) {
			// import "path" [as alias]
			// Remove quotes from the path
			rawPath := strings.Trim(p.curToken.Literal, `"`)
			// Resolve alias to actual Go package path
			importDecl.Path = `"` + stdlib.GetRealPackagePath(rawPath) + `"`
			p.nextToken()

			// Check for alias
			if p.curTokenIs(lexer.IDENT) && p.curToken.Literal == "as" {
				p.nextToken()
				if p.curTokenIs(lexer.IDENT) {
					importDecl.Alias = p.curToken.Literal
				}
			}
		} else if p.curTokenIs(lexer.IDENT) {
			// import os (without quotes for standard library)
			importDecl.Path = p.curToken.Literal
			// Convert to quoted format for standard library
			if isStandardLibrary(p.curToken.Literal) {
				importDecl.Path = p.curToken.Literal
			}
		}
	}

	return importDecl
}

// Helper function to check if a package is in Go's standard library
func isStandardLibrary(pkg string) bool {
	standardLibs := map[string]bool{
		"os": true, "fmt": true, "time": true, "json": true, "math": true,
		"strings": true, "strconv": true, "io": true, "bufio": true,
		"net": true, "http": true, "url": true, "path": true, "filepath": true,
		"sort": true, "sync": true, "context": true, "errors": true,
		"log": true, "regexp": true, "crypto": true, "encoding": true,
		"database": true, "html": true, "image": true, "mime": true,
		"reflect": true, "runtime": true, "testing": true, "unsafe": true,
	}
	return standardLibs[pkg]
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.FUNC:
		return p.parseFunctionDeclaration()
	case lexer.STRUCT:
		return p.parseStructDeclaration()
	case lexer.VAR:
		return p.parseVarDeclaration()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.IDENT:
		// Check if this is a variable assignment (identifier := value or identifier = value)
		if p.peekTokenIs(lexer.WALRUS) || p.peekTokenIs(lexer.ASSIGN) {
			return p.parseVarDeclaration()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseFunctionDeclaration() *ast.FunctionDecl {
	stmt := &ast.FunctionDecl{}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	// Optional return type
	if p.peekTokenIs(lexer.IDENT) || p.peekTokenIs(lexer.LPAREN) {
		p.nextToken()
		stmt.ReturnType = p.parseTypeSpec()
	}

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	// Skip newlines
	for p.peekTokenIs(lexer.NEWLINE) {
		p.nextToken()
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	params := []*ast.Parameter{}

	if p.peekTokenIs(lexer.RPAREN) {
		return params
	}

	p.nextToken()

	param := &ast.Parameter{Name: p.curToken.Literal}
	if p.peekTokenIs(lexer.IDENT) {
		p.nextToken()
		param.Type = p.parseTypeSpec()
	}
	params = append(params, param)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		param := &ast.Parameter{Name: p.curToken.Literal}
		if p.peekTokenIs(lexer.IDENT) {
			p.nextToken()
			param.Type = p.parseTypeSpec()
		}
		params = append(params, param)
	}

	return params
}

func (p *Parser) parseTypeSpec() *ast.TypeSpec {
	typeSpec := &ast.TypeSpec{}

	// Handle pointer types
	if p.curTokenIs(lexer.MULTIPLY) {
		typeSpec.IsPointer = true
		p.nextToken()
	}

	// Handle slice/array types
	if p.curTokenIs(lexer.LBRACKET) {
		p.nextToken()
		if p.curTokenIs(lexer.RBRACKET) {
			// Slice type []T
			typeSpec.IsSlice = true
			p.nextToken()
			typeSpec.ValueType = p.parseTypeSpec()
		} else {
			// Array type [N]T
			typeSpec.IsArray = true
			if p.curTokenIs(lexer.INT) {
				size, _ := strconv.Atoi(p.curToken.Literal)
				typeSpec.ArraySize = size
			}
			if !p.expectPeek(lexer.RBRACKET) {
				return nil
			}
			p.nextToken()
			typeSpec.ValueType = p.parseTypeSpec()
		}
		return typeSpec
	}

	// Handle map types
	if p.curTokenIs(lexer.IDENT) && p.curToken.Literal == "map" {
		if !p.expectPeek(lexer.LBRACKET) {
			return nil
		}
		p.nextToken()
		typeSpec.KeyType = p.parseTypeSpec()
		if !p.expectPeek(lexer.RBRACKET) {
			return nil
		}
		p.nextToken()
		typeSpec.ValueType = p.parseTypeSpec()
		return typeSpec
	}

	// Basic type
	typeSpec.Name = p.curToken.Literal
	return typeSpec
}

func (p *Parser) parseStructDeclaration() *ast.StructDecl {
	stmt := &ast.StructDecl{}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	// Skip newlines
	for p.peekTokenIs(lexer.NEWLINE) {
		p.nextToken()
	}

	// Parse fields
	for !p.curTokenIs(lexer.EOF) && p.curTokenIs(lexer.IDENT) || p.curTokenIs(lexer.FUNC) {
		if p.curTokenIs(lexer.NEWLINE) {
			p.nextToken()
			continue
		}

		if p.curTokenIs(lexer.FUNC) {
			// Method declaration
			method := p.parseFunctionDeclaration()
			if method != nil {
				// Add receiver
				method.Receiver = &ast.Parameter{
					Name: "self",
					Type: &ast.TypeSpec{Name: stmt.Name},
				}
				stmt.Methods = append(stmt.Methods, method)
			}
		} else if p.curTokenIs(lexer.IDENT) {
			// Field declaration
			field := &ast.Field{Name: p.curToken.Literal}
			p.nextToken()
			if p.curTokenIs(lexer.IDENT) {
				field.Type = p.parseTypeSpec()
			}
			stmt.Fields = append(stmt.Fields, field)
		}
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseVarDeclaration() *ast.VarDecl {
	stmt := &ast.VarDecl{}

	if p.curTokenIs(lexer.VAR) {
		// var name type = value OR var name = value
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		stmt.Name = p.curToken.Literal

		if p.peekTokenIs(lexer.IDENT) {
			// var name type = value
			p.nextToken()
			stmt.Type = p.parseTypeSpec()
			if !p.expectPeek(lexer.ASSIGN) {
				return nil
			}
			p.nextToken()
			stmt.Value = p.parseExpression(LOWEST)
		} else if p.peekTokenIs(lexer.ASSIGN) {
			// var name = value
			p.nextToken()
			p.nextToken()
			stmt.Value = p.parseExpression(LOWEST)
		}
	} else if p.curTokenIs(lexer.IDENT) {
		// name := value (walrus operator) or name = value (assignment)
		stmt.Name = p.curToken.Literal
		if p.peekTokenIs(lexer.WALRUS) {
			stmt.IsWalrus = true
			p.nextToken() // consume :=
			p.nextToken() // move to value
			stmt.Value = p.parseExpression(LOWEST)
		} else if p.peekTokenIs(lexer.ASSIGN) {
			stmt.IsWalrus = false
			p.nextToken() // consume =
			p.nextToken() // move to value
			stmt.Value = p.parseExpression(LOWEST)
		}
	}

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStmt {
	stmt := &ast.IfStmt{}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	// Skip newlines
	for p.peekTokenIs(lexer.NEWLINE) {
		p.nextToken()
	}

	stmt.ThenBranch = p.parseBlockStatement()

	// Handle elif/else
	if p.peekTokenIs(lexer.ELIF) || p.peekTokenIs(lexer.ELSE) {
		p.nextToken()
		if p.curTokenIs(lexer.ELIF) {
			stmt.ElseBranch = p.parseIfStatement()
		} else {
			// else
			if !p.expectPeek(lexer.COLON) {
				return nil
			}
			// Skip newlines
			for p.peekTokenIs(lexer.NEWLINE) {
				p.nextToken()
			}
			stmt.ElseBranch = p.parseBlockStatement()
		}
	}

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStmt {
	stmt := &ast.ForStmt{}

	p.nextToken()

	// Check for range-based for loop
	if p.peekTokenIs(lexer.IN) {
		stmt.IsRange = true
		stmt.RangeVar = p.curToken.Literal
		p.nextToken() // consume 'in'
		p.nextToken()
		stmt.RangeExpr = p.parseExpression(LOWEST)
	} else {
		// Traditional for loop
		stmt.Init = p.parseStatement()
		if !p.expectPeek(lexer.SEMICOLON) {
			return nil
		}
		p.nextToken()
		stmt.Condition = p.parseExpression(LOWEST)
		if !p.expectPeek(lexer.SEMICOLON) {
			return nil
		}
		p.nextToken()
		stmt.Update = p.parseStatement()
	}

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	// Skip newlines
	for p.peekTokenIs(lexer.NEWLINE) {
		p.nextToken()
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStmt {
	stmt := &ast.WhileStmt{}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	// Skip newlines
	for p.peekTokenIs(lexer.NEWLINE) {
		p.nextToken()
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{}

	if !p.peekTokenIs(lexer.NEWLINE) && !p.peekTokenIs(lexer.EOF) {
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStmt {
	stmt := &ast.ExpressionStmt{}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStmt {
	block := &ast.BlockStmt{}
	block.Statements = []ast.Statement{}

	// Parse all statements until we hit EOF or a function declaration
	for !p.curTokenIs(lexer.EOF) {
		// Skip newlines
		if p.curTokenIs(lexer.NEWLINE) {
			p.nextToken()
			continue
		}

		// Stop if we encounter another function (top-level)
		if p.curTokenIs(lexer.FUNC) {
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		if !p.curTokenIs(lexer.EOF) {
			p.nextToken()
		}
	}

	return block
}

// Expression parsing methods

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(lexer.NEWLINE) && !p.peekTokenIs(lexer.EOF) && !p.peekTokenIs(lexer.INDENT) && !p.peekTokenIs(lexer.DEDENT) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.Literal{Type: "int"}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.Literal{Type: "float"}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.Literal{Type: "string", Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.Literal{Type: "bool", Value: p.curTokenIs(lexer.TRUE)}
}

func (p *Parser) parseNilLiteral() ast.Expression {
	return &ast.Literal{Type: "nil", Value: nil}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.UnaryExpr{
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Operand = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.BinaryExpr{
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{}
	array.Elements = p.parseExpressionList(lexer.RBRACKET)
	return array
}

func (p *Parser) parseMapLiteral() ast.Expression {
	mapLit := &ast.MapLiteral{}
	mapLit.Pairs = []ast.MapPair{}

	if p.peekTokenIs(lexer.RBRACE) {
		p.nextToken()
		return mapLit
	}

	p.nextToken()

	for {
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(lexer.COLON) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(LOWEST)

		mapLit.Pairs = append(mapLit.Pairs, ast.MapPair{Key: key, Value: value})

		if !p.peekTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken()
		p.nextToken()
	}

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	return mapLit
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpr{Function: fn}
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)
	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpr{Object: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseSelectorExpression(left ast.Expression) ast.Expression {
	exp := &ast.SelectorExpr{Object: left}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	exp.Selector = p.curToken.Literal
	return exp
}

func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return args
}
