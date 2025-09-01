package ast

import (
	"fmt"
	"strings"
)

// Node represents any node in the AST
type Node interface {
	String() string
	Accept(visitor Visitor) interface{}
}

// Statement represents a statement node
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node
type Expression interface {
	Node
	expressionNode()
}

// Visitor pattern for AST traversal
type Visitor interface {
	VisitProgram(*Program) interface{}
	VisitFunctionDecl(*FunctionDecl) interface{}
	VisitStructDecl(*StructDecl) interface{}
	VisitVarDecl(*VarDecl) interface{}
	VisitIfStmt(*IfStmt) interface{}
	VisitForStmt(*ForStmt) interface{}
	VisitWhileStmt(*WhileStmt) interface{}
	VisitReturnStmt(*ReturnStmt) interface{}
	VisitExpressionStmt(*ExpressionStmt) interface{}
	VisitBlockStmt(*BlockStmt) interface{}
	VisitBinaryExpr(*BinaryExpr) interface{}
	VisitUnaryExpr(*UnaryExpr) interface{}
	VisitCallExpr(*CallExpr) interface{}
	VisitIdentifier(*Identifier) interface{}
	VisitLiteral(*Literal) interface{}
	VisitArrayLiteral(*ArrayLiteral) interface{}
	VisitMapLiteral(*MapLiteral) interface{}
	VisitIndexExpr(*IndexExpr) interface{}
	VisitSelectorExpr(*SelectorExpr) interface{}
}

// Program represents the root of the AST
type Program struct {
	Package    string
	Imports    []*ImportDecl
	Statements []Statement
}

func (p *Program) String() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("package %s\n", p.Package))
	for _, imp := range p.Imports {
		out.WriteString(imp.String() + "\n")
	}
	for _, stmt := range p.Statements {
		out.WriteString(stmt.String() + "\n")
	}
	return out.String()
}

func (p *Program) Accept(visitor Visitor) interface{} {
	return visitor.VisitProgram(p)
}

// ImportDecl represents an import declaration
type ImportDecl struct {
	Path  string
	Alias string
	Items []string // for "from X import Y, Z"
}

func (i *ImportDecl) String() string {
	if len(i.Items) > 0 {
		return fmt.Sprintf("from %s import %s", i.Path, strings.Join(i.Items, ", "))
	}
	if i.Alias != "" {
		return fmt.Sprintf("import %s as %s", i.Path, i.Alias)
	}
	return fmt.Sprintf("import %s", i.Path)
}

// FunctionDecl represents a function declaration
type FunctionDecl struct {
	Name       string
	Parameters []*Parameter
	ReturnType *TypeSpec
	Body       *BlockStmt
	Receiver   *Parameter // for methods
}

func (f *FunctionDecl) String() string {
	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	receiver := ""
	if f.Receiver != nil {
		receiver = fmt.Sprintf("(%s) ", f.Receiver.String())
	}
	returnType := ""
	if f.ReturnType != nil {
		returnType = " " + f.ReturnType.String()
	}
	return fmt.Sprintf("func %s%s(%s)%s:\n%s", receiver, f.Name, strings.Join(params, ", "), returnType, f.Body.String())
}

func (f *FunctionDecl) statementNode() {}
func (f *FunctionDecl) Accept(visitor Visitor) interface{} {
	return visitor.VisitFunctionDecl(f)
}

// Parameter represents a function parameter
type Parameter struct {
	Name string
	Type *TypeSpec
}

func (p *Parameter) String() string {
	if p.Type != nil {
		return fmt.Sprintf("%s %s", p.Name, p.Type.String())
	}
	return p.Name
}

// TypeSpec represents a type specification
type TypeSpec struct {
	Name      string
	IsPointer bool
	IsSlice   bool
	IsArray   bool
	ArraySize int
	KeyType   *TypeSpec // for maps
	ValueType *TypeSpec // for maps, slices, arrays
}

func (t *TypeSpec) String() string {
	result := ""
	if t.IsPointer {
		result += "*"
	}
	if t.IsSlice {
		result += "[]"
	}
	if t.IsArray {
		result += fmt.Sprintf("[%d]", t.ArraySize)
	}
	if t.KeyType != nil && t.ValueType != nil {
		result += fmt.Sprintf("map[%s]%s", t.KeyType.String(), t.ValueType.String())
	} else if t.ValueType != nil {
		result += t.ValueType.String()
	} else {
		result += t.Name
	}
	return result
}

// StructDecl represents a struct declaration
type StructDecl struct {
	Name    string
	Fields  []*Field
	Methods []*FunctionDecl
}

func (s *StructDecl) String() string {
	var fields []string
	for _, f := range s.Fields {
		fields = append(fields, f.String())
	}
	return fmt.Sprintf("struct %s:\n    %s", s.Name, strings.Join(fields, "\n    "))
}

func (s *StructDecl) statementNode() {}
func (s *StructDecl) Accept(visitor Visitor) interface{} {
	return visitor.VisitStructDecl(s)
}

// Field represents a struct field
type Field struct {
	Name string
	Type *TypeSpec
	Tag  string
}

func (f *Field) String() string {
	tag := ""
	if f.Tag != "" {
		tag = fmt.Sprintf(" `%s`", f.Tag)
	}
	return fmt.Sprintf("%s %s%s", f.Name, f.Type.String(), tag)
}

// VarDecl represents a variable declaration
type VarDecl struct {
	Name     string
	Type     *TypeSpec
	Value    Expression
	IsWalrus bool // true for :=, false for =
}

func (v *VarDecl) String() string {
	if v.Type != nil {
		return fmt.Sprintf("var %s %s = %s", v.Name, v.Type.String(), v.Value.String())
	}
	return fmt.Sprintf("%s := %s", v.Name, v.Value.String())
}

func (v *VarDecl) statementNode() {}
func (v *VarDecl) Accept(visitor Visitor) interface{} {
	return visitor.VisitVarDecl(v)
}

// BlockStmt represents a block of statements
type BlockStmt struct {
	Statements []Statement
}

func (b *BlockStmt) String() string {
	var stmts []string
	for _, stmt := range b.Statements {
		stmts = append(stmts, "    "+stmt.String())
	}
	return strings.Join(stmts, "\n")
}

func (b *BlockStmt) statementNode() {}
func (b *BlockStmt) Accept(visitor Visitor) interface{} {
	return visitor.VisitBlockStmt(b)
}

// IfStmt represents an if statement
type IfStmt struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (i *IfStmt) String() string {
	result := fmt.Sprintf("if %s:\n%s", i.Condition.String(), i.ThenBranch.String())
	if i.ElseBranch != nil {
		result += fmt.Sprintf("\nelse:\n%s", i.ElseBranch.String())
	}
	return result
}

func (i *IfStmt) statementNode() {}
func (i *IfStmt) Accept(visitor Visitor) interface{} {
	return visitor.VisitIfStmt(i)
}

// ForStmt represents a for loop
type ForStmt struct {
	Init      Statement
	Condition Expression
	Update    Statement
	Body      *BlockStmt
	IsRange   bool
	RangeVar  string
	RangeExpr Expression
}

func (f *ForStmt) String() string {
	if f.IsRange {
		return fmt.Sprintf("for %s in %s:\n%s", f.RangeVar, f.RangeExpr.String(), f.Body.String())
	}
	return fmt.Sprintf("for %s; %s; %s:\n%s", f.Init.String(), f.Condition.String(), f.Update.String(), f.Body.String())
}

func (f *ForStmt) statementNode() {}
func (f *ForStmt) Accept(visitor Visitor) interface{} {
	return visitor.VisitForStmt(f)
}

// WhileStmt represents a while loop
type WhileStmt struct {
	Condition Expression
	Body      *BlockStmt
}

func (w *WhileStmt) String() string {
	return fmt.Sprintf("while %s:\n%s", w.Condition.String(), w.Body.String())
}

func (w *WhileStmt) statementNode() {}
func (w *WhileStmt) Accept(visitor Visitor) interface{} {
	return visitor.VisitWhileStmt(w)
}

// ReturnStmt represents a return statement
type ReturnStmt struct {
	Value Expression
}

func (r *ReturnStmt) String() string {
	if r.Value != nil {
		return fmt.Sprintf("return %s", r.Value.String())
	}
	return "return"
}

func (r *ReturnStmt) statementNode() {}
func (r *ReturnStmt) Accept(visitor Visitor) interface{} {
	return visitor.VisitReturnStmt(r)
}

// ExpressionStmt represents an expression used as a statement
type ExpressionStmt struct {
	Expression Expression
}

func (e *ExpressionStmt) String() string {
	return e.Expression.String()
}

func (e *ExpressionStmt) statementNode() {}
func (e *ExpressionStmt) Accept(visitor Visitor) interface{} {
	return visitor.VisitExpressionStmt(e)
}

// Expression implementations

// BinaryExpr represents a binary expression
type BinaryExpr struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Operator, b.Right.String())
}

func (b *BinaryExpr) expressionNode() {}
func (b *BinaryExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

// UnaryExpr represents a unary expression
type UnaryExpr struct {
	Operator string
	Operand  Expression
}

func (u *UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", u.Operator, u.Operand.String())
}

func (u *UnaryExpr) expressionNode() {}
func (u *UnaryExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

// CallExpr represents a function call
type CallExpr struct {
	Function  Expression
	Arguments []Expression
}

func (c *CallExpr) String() string {
	var args []string
	for _, arg := range c.Arguments {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", c.Function.String(), strings.Join(args, ", "))
}

func (c *CallExpr) expressionNode() {}
func (c *CallExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitCallExpr(c)
}

// Identifier represents an identifier
type Identifier struct {
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) Accept(visitor Visitor) interface{} {
	return visitor.VisitIdentifier(i)
}

// Literal represents a literal value
type Literal struct {
	Type  string // "int", "float", "string", "bool", "nil"
	Value interface{}
}

func (l *Literal) String() string {
	switch l.Type {
	case "string":
		return fmt.Sprintf(`"%v"`, l.Value)
	default:
		return fmt.Sprintf("%v", l.Value)
	}
}

func (l *Literal) expressionNode() {}
func (l *Literal) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteral(l)
}

// ArrayLiteral represents an array literal
type ArrayLiteral struct {
	Elements []Expression
}

func (a *ArrayLiteral) String() string {
	var elements []string
	for _, elem := range a.Elements {
		elements = append(elements, elem.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (a *ArrayLiteral) expressionNode() {}
func (a *ArrayLiteral) Accept(visitor Visitor) interface{} {
	return visitor.VisitArrayLiteral(a)
}

// MapLiteral represents a map literal
type MapLiteral struct {
	Pairs []MapPair
}

type MapPair struct {
	Key   Expression
	Value Expression
}

func (m *MapLiteral) String() string {
	var pairs []string
	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.String(), pair.Value.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

func (m *MapLiteral) expressionNode() {}
func (m *MapLiteral) Accept(visitor Visitor) interface{} {
	return visitor.VisitMapLiteral(m)
}

// IndexExpr represents an index expression (array[index])
type IndexExpr struct {
	Object Expression
	Index  Expression
}

func (i *IndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", i.Object.String(), i.Index.String())
}

func (i *IndexExpr) expressionNode() {}
func (i *IndexExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitIndexExpr(i)
}

// SelectorExpr represents a selector expression (object.field)
type SelectorExpr struct {
	Object   Expression
	Selector string
}

func (s *SelectorExpr) String() string {
	return fmt.Sprintf("%s.%s", s.Object.String(), s.Selector)
}

func (s *SelectorExpr) expressionNode() {}
func (s *SelectorExpr) Accept(visitor Visitor) interface{} {
	return visitor.VisitSelectorExpr(s)
}
