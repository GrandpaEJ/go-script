package lexer

// No additional imports needed for basic lexer functionality

// Lexer represents the lexical analyzer
type Lexer struct {
	input        string
	position     int   // current position in input (points to current char)
	readPosition int   // current reading position in input (after current char)
	ch           byte  // current char under examination
	line         int   // current line number
	column       int   // current column number
	indentStack  []int // stack to track indentation levels
}

// New creates a new lexer instance
func New(input string) *Lexer {
	l := &Lexer{
		input:       input,
		line:        1,
		column:      0,
		indentStack: []int{0}, // start with 0 indentation
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL character represents EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// peekCharAt returns the character at offset positions ahead
func (l *Lexer) peekCharAt(offset int) byte {
	pos := l.readPosition + offset - 1
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

// skipWhitespace skips whitespace characters except newlines
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (integer or float)
func (l *Lexer) readNumber() (string, TokenType) {
	position := l.position
	tokenType := INT

	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		tokenType = FLOAT
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Check for scientific notation
	if l.ch == 'e' || l.ch == 'E' {
		tokenType = FLOAT
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position], tokenType
}

// readString reads a string literal
func (l *Lexer) readString(delimiter byte) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == delimiter || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar() // skip escaped character
		}
	}
	return l.input[position:l.position]
}

// readComment reads a comment
func (l *Lexer) readComment() string {
	position := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

// handleIndentation processes indentation at the beginning of a line
func (l *Lexer) handleIndentation() []Token {
	var tokens []Token
	indentLevel := 0

	// Count spaces/tabs for indentation
	for l.ch == ' ' || l.ch == '\t' {
		if l.ch == '\t' {
			indentLevel += 4 // treat tab as 4 spaces
		} else {
			indentLevel++
		}
		l.readChar()
	}

	// Skip empty lines and comments
	if l.ch == '\n' || l.ch == '#' {
		return tokens
	}

	currentIndent := l.indentStack[len(l.indentStack)-1]

	if indentLevel > currentIndent {
		// Increased indentation - INDENT token
		l.indentStack = append(l.indentStack, indentLevel)
		tokens = append(tokens, Token{
			Type:     INDENT,
			Literal:  "",
			Line:     l.line,
			Column:   l.column - indentLevel,
			Position: l.position - indentLevel,
		})
	} else if indentLevel < currentIndent {
		// Decreased indentation - DEDENT tokens
		for len(l.indentStack) > 1 && l.indentStack[len(l.indentStack)-1] > indentLevel {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			tokens = append(tokens, Token{
				Type:     DEDENT,
				Literal:  "",
				Line:     l.line,
				Column:   l.column - indentLevel,
				Position: l.position - indentLevel,
			})
		}
	}

	return tokens
}

// NextToken returns the next token
func (l *Lexer) NextToken() Token {
	var tok Token

	// Skip indentation handling for now - simplified approach
	// TODO: Implement proper indentation handling with token queue

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else {
			tok = newToken(ASSIGN, l.ch, l.line, l.column, l.position)
		}
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: PLUS_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: INCREMENT, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else {
			tok = newToken(PLUS, l.ch, l.line, l.column, l.position)
		}
	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: MINUS_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: DECREMENT, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: ARROW, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else {
			tok = newToken(MINUS, l.ch, l.line, l.column, l.position)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: MULT_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: POWER, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else {
			tok = newToken(MULTIPLY, l.ch, l.line, l.column, l.position)
		}
	case '/':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: DIV_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else {
			tok = newToken(DIVIDE, l.ch, l.line, l.column, l.position)
		}
	case '%':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: MOD_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else {
			tok = newToken(MODULO, l.ch, l.line, l.column, l.position)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column, l.position)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LT_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else if l.peekChar() == '<' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LEFT_SHIFT, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: CHANNEL, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else {
			tok = newToken(LT, l.ch, l.line, l.column, l.position)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GT_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: RIGHT_SHIFT, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else {
			tok = newToken(GT, l.ch, l.line, l.column, l.position)
		}
	case '&':
		if l.peekChar() == '^' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: BIT_CLEAR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else {
			tok = newToken(BITWISE_AND, l.ch, l.line, l.column, l.position)
		}
	case '|':
		tok = newToken(BITWISE_OR, l.ch, l.line, l.column, l.position)
	case '^':
		tok = newToken(BITWISE_XOR, l.ch, l.line, l.column, l.position)
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: WALRUS, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1, Position: l.position - 1}
		} else {
			tok = newToken(COLON, l.ch, l.line, l.column, l.position)
		}
	case ';':
		tok = newToken(SEMICOLON, l.ch, l.line, l.column, l.position)
	case ',':
		tok = newToken(COMMA, l.ch, l.line, l.column, l.position)
	case '.':
		tok = newToken(DOT, l.ch, l.line, l.column, l.position)
	case '(':
		tok = newToken(LPAREN, l.ch, l.line, l.column, l.position)
	case ')':
		tok = newToken(RPAREN, l.ch, l.line, l.column, l.position)
	case '[':
		tok = newToken(LBRACKET, l.ch, l.line, l.column, l.position)
	case ']':
		tok = newToken(RBRACKET, l.ch, l.line, l.column, l.position)
	case '{':
		tok = newToken(LBRACE, l.ch, l.line, l.column, l.position)
	case '}':
		tok = newToken(RBRACE, l.ch, l.line, l.column, l.position)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString('"')
		tok.Line = l.line
		tok.Column = l.column
		tok.Position = l.position
	case '\'':
		tok.Type = CHAR
		tok.Literal = l.readString('\'')
		tok.Line = l.line
		tok.Column = l.column
		tok.Position = l.position
	case '#':
		tok.Type = COMMENT
		tok.Literal = l.readComment()
		tok.Line = l.line
		tok.Column = l.column
		tok.Position = l.position
		return tok // Don't advance past comment
	case '\n':
		tok = newToken(NEWLINE, l.ch, l.line, l.column, l.position)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = l.line
		tok.Column = l.column
		tok.Position = l.position
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Column = l.column - len(tok.Literal)
			tok.Position = l.position - len(tok.Literal)
			return tok // Don't advance past identifier
		} else if isDigit(l.ch) {
			tok.Literal, tok.Type = l.readNumber()
			tok.Line = l.line
			tok.Column = l.column - len(tok.Literal)
			tok.Position = l.position - len(tok.Literal)
			return tok // Don't advance past number
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column, l.position)
		}
	}

	l.readChar()
	return tok
}

// newToken creates a new token
func newToken(tokenType TokenType, ch byte, line, column, position int) Token {
	return Token{
		Type:     tokenType,
		Literal:  string(ch),
		Line:     line,
		Column:   column,
		Position: position,
	}
}

// isLetter checks if a character is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch > 127
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isAlphaNumeric checks if a character is alphanumeric
func isAlphaNumeric(ch byte) bool {
	return isLetter(ch) || isDigit(ch)
}
