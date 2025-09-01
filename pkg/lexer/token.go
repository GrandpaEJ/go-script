package lexer

import "fmt"

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	// Literals
	IDENT  // identifiers
	INT    // integers
	FLOAT  // floating point numbers
	STRING // string literals
	CHAR   // character literals

	// Keywords
	AND
	OR
	NOT
	IF
	ELIF
	ELSE
	FOR
	WHILE
	FUNC
	RETURN
	IMPORT
	FROM
	STRUCT
	INTERFACE
	VAR
	CONST
	TRUE
	FALSE
	NIL
	IN
	RANGE
	BREAK
	CONTINUE
	DEFER
	GO
	CHAN
	SELECT
	CASE
	DEFAULT
	SWITCH
	TYPE
	PACKAGE

	// Operators
	ASSIGN    // =
	WALRUS    // :=
	PLUS      // +
	MINUS     // -
	MULTIPLY  // *
	DIVIDE    // /
	MODULO    // %
	POWER     // **
	EQ        // ==
	NOT_EQ    // !=
	LT        // <
	LT_EQ     // <=
	GT        // >
	GT_EQ     // >=
	PLUS_EQ   // +=
	MINUS_EQ  // -=
	MULT_EQ   // *=
	DIV_EQ    // /=
	MOD_EQ    // %=
	BITWISE_AND // &
	BITWISE_OR  // |
	BITWISE_XOR // ^
	LEFT_SHIFT  // <<
	RIGHT_SHIFT // >>
	BIT_CLEAR   // &^
	INCREMENT   // ++
	DECREMENT   // --

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	DOT       // .
	ARROW     // ->
	CHANNEL   // <-

	// Brackets
	LPAREN   // (
	RPAREN   // )
	LBRACKET // [
	RBRACKET // ]
	LBRACE   // {
	RBRACE   // }

	// Special
	NEWLINE
	INDENT
	DEDENT
)

// Token represents a single token
type Token struct {
	Type     TokenType
	Literal  string
	Line     int
	Column   int
	Position int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Literal: %q, Line: %d, Column: %d}",
		TokenTypeString(t.Type), t.Literal, t.Line, t.Column)
}

// TokenTypeString returns the string representation of a token type
func TokenTypeString(tokenType TokenType) string {
	switch tokenType {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case COMMENT:
		return "COMMENT"
	case IDENT:
		return "IDENT"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case CHAR:
		return "CHAR"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case IF:
		return "IF"
	case ELIF:
		return "ELIF"
	case ELSE:
		return "ELSE"
	case FOR:
		return "FOR"
	case WHILE:
		return "WHILE"
	case FUNC:
		return "FUNC"
	case RETURN:
		return "RETURN"
	case IMPORT:
		return "IMPORT"
	case FROM:
		return "FROM"
	case STRUCT:
		return "STRUCT"
	case INTERFACE:
		return "INTERFACE"
	case VAR:
		return "VAR"
	case CONST:
		return "CONST"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NIL:
		return "NIL"
	case IN:
		return "IN"
	case RANGE:
		return "RANGE"
	case BREAK:
		return "BREAK"
	case CONTINUE:
		return "CONTINUE"
	case DEFER:
		return "DEFER"
	case GO:
		return "GO"
	case CHAN:
		return "CHAN"
	case SELECT:
		return "SELECT"
	case CASE:
		return "CASE"
	case DEFAULT:
		return "DEFAULT"
	case SWITCH:
		return "SWITCH"
	case TYPE:
		return "TYPE"
	case PACKAGE:
		return "PACKAGE"
	case ASSIGN:
		return "ASSIGN"
	case WALRUS:
		return "WALRUS"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case MODULO:
		return "MODULO"
	case POWER:
		return "POWER"
	case EQ:
		return "EQ"
	case NOT_EQ:
		return "NOT_EQ"
	case LT:
		return "LT"
	case LT_EQ:
		return "LT_EQ"
	case GT:
		return "GT"
	case GT_EQ:
		return "GT_EQ"
	case PLUS_EQ:
		return "PLUS_EQ"
	case MINUS_EQ:
		return "MINUS_EQ"
	case MULT_EQ:
		return "MULT_EQ"
	case DIV_EQ:
		return "DIV_EQ"
	case MOD_EQ:
		return "MOD_EQ"
	case BITWISE_AND:
		return "BITWISE_AND"
	case BITWISE_OR:
		return "BITWISE_OR"
	case BITWISE_XOR:
		return "BITWISE_XOR"
	case LEFT_SHIFT:
		return "LEFT_SHIFT"
	case RIGHT_SHIFT:
		return "RIGHT_SHIFT"
	case BIT_CLEAR:
		return "BIT_CLEAR"
	case INCREMENT:
		return "INCREMENT"
	case DECREMENT:
		return "DECREMENT"
	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case COLON:
		return "COLON"
	case DOT:
		return "DOT"
	case ARROW:
		return "ARROW"
	case CHANNEL:
		return "CHANNEL"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case NEWLINE:
		return "NEWLINE"
	case INDENT:
		return "INDENT"
	case DEDENT:
		return "DEDENT"
	default:
		return "UNKNOWN"
	}
}

// Keywords maps keyword strings to their token types
var Keywords = map[string]TokenType{
	"and":       AND,
	"or":        OR,
	"not":       NOT,
	"if":        IF,
	"elif":      ELIF,
	"else":      ELSE,
	"for":       FOR,
	"while":     WHILE,
	"func":      FUNC,
	"return":    RETURN,
	"import":    IMPORT,
	"from":      FROM,
	"struct":    STRUCT,
	"interface": INTERFACE,
	"var":       VAR,
	"const":     CONST,
	"true":      TRUE,
	"false":     FALSE,
	"nil":       NIL,
	"in":        IN,
	"range":     RANGE,
	"break":     BREAK,
	"continue":  CONTINUE,
	"defer":     DEFER,
	"go":        GO,
	"chan":      CHAN,
	"select":    SELECT,
	"case":      CASE,
	"default":   DEFAULT,
	"switch":    SWITCH,
	"type":      TYPE,
	"package":   PACKAGE,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}
