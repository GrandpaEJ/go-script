package tests

import (
	"testing"

	"github.com/GrandpaEJ/go-script/pkg/lexer"
)

func TestLexerBasicTokens(t *testing.T) {
	input := `func main():
    x := 42
    print("Hello")
    return x + 1`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.FUNC, "func"},
		{lexer.IDENT, "main"},
		{lexer.LPAREN, "("},
		{lexer.RPAREN, ")"},
		{lexer.COLON, ":"},
		{lexer.NEWLINE, "\n"},
		{lexer.IDENT, "x"},
		{lexer.WALRUS, ":="},
		{lexer.INT, "42"},
		{lexer.NEWLINE, "\n"},
		{lexer.IDENT, "print"},
		{lexer.LPAREN, "("},
		{lexer.STRING, "Hello"},
		{lexer.RPAREN, ")"},
		{lexer.NEWLINE, "\n"},
		{lexer.RETURN, "return"},
		{lexer.IDENT, "x"},
		{lexer.PLUS, "+"},
		{lexer.INT, "1"},
		{lexer.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, lexer.TokenTypeString(tt.expectedType), lexer.TokenTypeString(tok.Type))
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexerOperators(t *testing.T) {
	input := `+ - * / % ** == != < <= > >= and or not := = += -= *= /= %=`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.PLUS, "+"},
		{lexer.MINUS, "-"},
		{lexer.MULTIPLY, "*"},
		{lexer.DIVIDE, "/"},
		{lexer.MODULO, "%"},
		{lexer.POWER, "**"},
		{lexer.EQ, "=="},
		{lexer.NOT_EQ, "!="},
		{lexer.LT, "<"},
		{lexer.LT_EQ, "<="},
		{lexer.GT, ">"},
		{lexer.GT_EQ, ">="},
		{lexer.AND, "and"},
		{lexer.OR, "or"},
		{lexer.NOT, "not"},
		{lexer.WALRUS, ":="},
		{lexer.ASSIGN, "="},
		{lexer.PLUS_EQ, "+="},
		{lexer.MINUS_EQ, "-="},
		{lexer.MULT_EQ, "*="},
		{lexer.DIV_EQ, "/="},
		{lexer.MOD_EQ, "%="},
		{lexer.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, lexer.TokenTypeString(tt.expectedType), lexer.TokenTypeString(tok.Type))
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexerKeywords(t *testing.T) {
	input := `func if elif else for while return import from struct var const true false nil`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.FUNC, "func"},
		{lexer.IF, "if"},
		{lexer.ELIF, "elif"},
		{lexer.ELSE, "else"},
		{lexer.FOR, "for"},
		{lexer.WHILE, "while"},
		{lexer.RETURN, "return"},
		{lexer.IMPORT, "import"},
		{lexer.FROM, "from"},
		{lexer.STRUCT, "struct"},
		{lexer.VAR, "var"},
		{lexer.CONST, "const"},
		{lexer.TRUE, "true"},
		{lexer.FALSE, "false"},
		{lexer.NIL, "nil"},
		{lexer.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, lexer.TokenTypeString(tt.expectedType), lexer.TokenTypeString(tok.Type))
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexerNumbers(t *testing.T) {
	input := `42 3.14 1.5e10 2.5E-3`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.INT, "42"},
		{lexer.FLOAT, "3.14"},
		{lexer.FLOAT, "1.5e10"},
		{lexer.FLOAT, "2.5E-3"},
		{lexer.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, lexer.TokenTypeString(tt.expectedType), lexer.TokenTypeString(tok.Type))
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexerStrings(t *testing.T) {
	input := `"hello world" "escaped \"quote\"" 'single quote'`

	tests := []struct {
		expectedType    lexer.TokenType
		expectedLiteral string
	}{
		{lexer.STRING, "hello world"},
		{lexer.STRING, "escaped \\\"quote\\\""},
		{lexer.CHAR, "single quote"},
		{lexer.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, lexer.TokenTypeString(tt.expectedType), lexer.TokenTypeString(tok.Type))
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexerComments(t *testing.T) {
	input := `# This is a comment
func main():
    # Another comment
    print("hello")`

	l := lexer.New(input)

	// Should skip comments and get the actual tokens
	expectedTokens := []lexer.TokenType{
		lexer.NEWLINE,
		lexer.FUNC,
		lexer.IDENT,
		lexer.LPAREN,
		lexer.RPAREN,
		lexer.COLON,
		lexer.NEWLINE,
		lexer.IDENT,
		lexer.LPAREN,
		lexer.STRING,
		lexer.RPAREN,
		lexer.EOF,
	}

	for i, expectedType := range expectedTokens {
		tok := l.NextToken()
		if tok.Type != expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, lexer.TokenTypeString(expectedType), lexer.TokenTypeString(tok.Type))
		}
	}
}
