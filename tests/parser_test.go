package tests

import (
	"testing"

	"github.com/GrandpaEJ/go-script/pkg/ast"
	"github.com/GrandpaEJ/go-script/pkg/lexer"
	"github.com/GrandpaEJ/go-script/pkg/parser"
)

func TestFunctionDeclaration(t *testing.T) {
	input := `func add(x int, y int) int:
    return x + y`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionDecl)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.FunctionDecl. got=%T",
			program.Statements[0])
	}

	if stmt.Name != "add" {
		t.Fatalf("function name wrong. expected='add', got='%s'", stmt.Name)
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("function parameters wrong. expected=2, got=%d", len(stmt.Parameters))
	}

	if stmt.Parameters[0].Name != "x" || stmt.Parameters[0].Type.Name != "int" {
		t.Fatalf("parameter 0 wrong. expected='x int', got='%s %s'",
			stmt.Parameters[0].Name, stmt.Parameters[0].Type.Name)
	}

	if stmt.Parameters[1].Name != "y" || stmt.Parameters[1].Type.Name != "int" {
		t.Fatalf("parameter 1 wrong. expected='y int', got='%s %s'",
			stmt.Parameters[1].Name, stmt.Parameters[1].Type.Name)
	}

	if stmt.ReturnType.Name != "int" {
		t.Fatalf("return type wrong. expected='int', got='%s'", stmt.ReturnType.Name)
	}
}

func TestVariableDeclaration(t *testing.T) {
	tests := []struct {
		input         string
		expectedName  string
		expectedValue interface{}
	}{
		{"x := 5", "x", 5},
		{"y := true", "y", true},
		{"name := \"hello\"", "name", "hello"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.VarDecl)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.VarDecl. got=%T",
				program.Statements[0])
		}

		if stmt.Name != tt.expectedName {
			t.Fatalf("variable name wrong. expected='%s', got='%s'",
				tt.expectedName, stmt.Name)
		}

		if !testLiteralExpression(t, stmt.Value, tt.expectedValue) {
			return
		}
	}
}

func TestIfStatement(t *testing.T) {
	input := `if x > 5:
    return true
else:
    return false`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStmt. got=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", ">", 5) {
		return
	}

	thenStmt, ok := stmt.ThenBranch.(*ast.BlockStmt)
	if !ok {
		t.Fatalf("stmt.ThenBranch is not *ast.BlockStmt. got=%T", stmt.ThenBranch)
	}

	if len(thenStmt.Statements) != 1 {
		t.Fatalf("thenStmt.Statements does not contain 1 statement. got=%d",
			len(thenStmt.Statements))
	}

	elseStmt, ok := stmt.ElseBranch.(*ast.BlockStmt)
	if !ok {
		t.Fatalf("stmt.ElseBranch is not *ast.BlockStmt. got=%T", stmt.ElseBranch)
	}

	if len(elseStmt.Statements) != 1 {
		t.Fatalf("elseStmt.Statements does not contain 1 statement. got=%d",
			len(elseStmt.Statements))
	}
}

func TestForStatement(t *testing.T) {
	input := `for i in range(10):
    print(i)`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForStmt)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ForStmt. got=%T",
			program.Statements[0])
	}

	if !stmt.IsRange {
		t.Fatalf("stmt.IsRange is not true. got=%t", stmt.IsRange)
	}

	if stmt.RangeVar != "i" {
		t.Fatalf("stmt.RangeVar wrong. expected='i', got='%s'", stmt.RangeVar)
	}

	callExpr, ok := stmt.RangeExpr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("stmt.RangeExpr is not *ast.CallExpr. got=%T", stmt.RangeExpr)
	}

	if !testIdentifier(t, callExpr.Function, "range") {
		return
	}

	if len(callExpr.Arguments) != 1 {
		t.Fatalf("wrong number of arguments. expected=1, got=%d", len(callExpr.Arguments))
	}

	if !testLiteralExpression(t, callExpr.Arguments[0], 10) {
		return
	}
}

func TestStructDeclaration(t *testing.T) {
	input := `struct Person:
    name string
    age int
    
    func greet(self) string:
        return "Hello"`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.StructDecl)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.StructDecl. got=%T",
			program.Statements[0])
	}

	if stmt.Name != "Person" {
		t.Fatalf("struct name wrong. expected='Person', got='%s'", stmt.Name)
	}

	if len(stmt.Fields) != 2 {
		t.Fatalf("struct fields wrong. expected=2, got=%d", len(stmt.Fields))
	}

	if stmt.Fields[0].Name != "name" || stmt.Fields[0].Type.Name != "string" {
		t.Fatalf("field 0 wrong. expected='name string', got='%s %s'",
			stmt.Fields[0].Name, stmt.Fields[0].Type.Name)
	}

	if stmt.Fields[1].Name != "age" || stmt.Fields[1].Type.Name != "int" {
		t.Fatalf("field 1 wrong. expected='age int', got='%s %s'",
			stmt.Fields[1].Name, stmt.Fields[1].Type.Name)
	}

	if len(stmt.Methods) != 1 {
		t.Fatalf("struct methods wrong. expected=1, got=%d", len(stmt.Methods))
	}

	if stmt.Methods[0].Name != "greet" {
		t.Fatalf("method name wrong. expected='greet', got='%s'", stmt.Methods[0].Name)
	}
}

func TestExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"5", 5},
		{"10", 10},
		{"true", true},
		{"false", false},
		{"\"hello world\"", "hello world"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStmt. got=%T",
				program.Statements[0])
		}

		if !testLiteralExpression(t, stmt.Expression, tt.expected) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testStringLiteral(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.Literal)
	if !ok {
		t.Errorf("il not *ast.Literal. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	str, ok := exp.(*ast.Literal)
	if !ok {
		t.Errorf("exp not *ast.Literal. got=%T", exp)
		return false
	}

	if str.Value != value {
		t.Errorf("str.Value not %q. got=%q", value, str.Value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Literal)
	if !ok {
		t.Errorf("exp not *ast.Literal. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("exp is not *ast.BinaryExpr. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}
