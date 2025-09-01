package codegen

import (
	"fmt"
	"strings"

	"github.com/GrandpaEJ/go-script/pkg/ast"
)

// Generator represents the code generator
type Generator struct {
	output      strings.Builder
	indentLevel int
}

// New creates a new code generator
func New() *Generator {
	return &Generator{}
}

// Generate generates Go code from the AST
func (g *Generator) Generate(program *ast.Program) string {
	g.output.Reset()
	g.indentLevel = 0

	// Generate package declaration
	g.writeLine(fmt.Sprintf("package %s", program.Package))
	g.writeLine("")

	// Generate imports
	if len(program.Imports) > 0 {
		g.writeLine("import (")
		g.indentLevel++
		for _, imp := range program.Imports {
			g.generateImport(imp)
		}
		g.indentLevel--
		g.writeLine(")")
		g.writeLine("")
	}

	// Generate statements
	for i, stmt := range program.Statements {
		g.generateStatement(stmt)
		// Add blank line between top-level statements, but not after the last one
		if i < len(program.Statements)-1 {
			g.writeLine("")
		}
	}

	return g.output.String()
}

func (g *Generator) generateImport(imp *ast.ImportDecl) {
	if len(imp.Items) > 0 {
		if imp.Path != "" {
			// Handle "from X import Y, Z" style imports
			for _, item := range imp.Items {
				g.writeLine(fmt.Sprintf(`%s "%s"`, item, imp.Path))
			}
		} else {
			// Handle import ("os", "fmt", "time") style
			for _, item := range imp.Items {
				// Remove quotes if already present
				cleanItem := strings.Trim(item, `"`)
				g.writeLine(fmt.Sprintf(`"%s"`, cleanItem))
			}
		}
	} else if imp.Alias != "" {
		g.writeLine(fmt.Sprintf(`%s "%s"`, imp.Alias, imp.Path))
	} else {
		// Handle both quoted and unquoted imports
		path := imp.Path
		if !strings.HasPrefix(path, `"`) && !strings.HasSuffix(path, `"`) {
			// Add quotes if not present
			path = fmt.Sprintf(`"%s"`, path)
		}
		g.writeLine(path)
	}
}

func (g *Generator) generateStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.FunctionDecl:
		g.generateFunctionDecl(s)
	case *ast.StructDecl:
		g.generateStructDecl(s)
	case *ast.VarDecl:
		g.generateVarDecl(s)
	case *ast.IfStmt:
		g.generateIfStmt(s)
	case *ast.ForStmt:
		g.generateForStmt(s)
	case *ast.WhileStmt:
		g.generateWhileStmt(s)
	case *ast.ReturnStmt:
		g.generateReturnStmt(s)
	case *ast.ExpressionStmt:
		g.generateExpressionStmt(s)
	case *ast.BlockStmt:
		g.generateBlockStmt(s)
	}
}

func (g *Generator) generateFunctionDecl(fn *ast.FunctionDecl) {
	// Generate function signature
	signature := "func "

	// Add receiver if it's a method
	if fn.Receiver != nil {
		signature += fmt.Sprintf("(%s) ", g.generateParameter(fn.Receiver))
	}

	signature += fn.Name + "("

	// Add parameters
	var params []string
	for _, param := range fn.Parameters {
		params = append(params, g.generateParameter(param))
	}
	signature += strings.Join(params, ", ")
	signature += ")"

	// Add return type
	if fn.ReturnType != nil {
		signature += " " + g.generateTypeSpec(fn.ReturnType)
	}

	g.writeLine(signature + " {")
	g.indentLevel++
	g.generateBlockStmt(fn.Body)
	g.indentLevel--
	g.writeLine("}")
}

func (g *Generator) generateStructDecl(s *ast.StructDecl) {
	g.writeLine(fmt.Sprintf("type %s struct {", s.Name))
	g.indentLevel++

	for _, field := range s.Fields {
		g.generateField(field)
	}

	g.indentLevel--
	g.writeLine("}")

	// Generate methods separately
	for _, method := range s.Methods {
		g.writeLine("")
		g.generateFunctionDecl(method)
	}
}

func (g *Generator) generateField(field *ast.Field) {
	line := field.Name
	if field.Type != nil {
		line += " " + g.generateTypeSpec(field.Type)
	}
	if field.Tag != "" {
		line += " `" + field.Tag + "`"
	}
	g.writeLine(line)
}

func (g *Generator) generateVarDecl(v *ast.VarDecl) {
	if v.Type != nil {
		// var name type = value
		line := fmt.Sprintf("var %s %s", v.Name, g.generateTypeSpec(v.Type))
		if v.Value != nil {
			line += " = " + g.generateExpression(v.Value)
		}
		g.writeLine(line)
	} else if v.IsWalrus {
		// name := value (new variable)
		g.writeLine(fmt.Sprintf("%s := %s", v.Name, g.generateExpression(v.Value)))
	} else {
		// name = value (assignment to existing variable)
		g.writeLine(fmt.Sprintf("%s = %s", v.Name, g.generateExpression(v.Value)))
	}
}

func (g *Generator) generateIfStmt(i *ast.IfStmt) {
	g.writeLine(fmt.Sprintf("if %s {", g.generateExpression(i.Condition)))
	g.indentLevel++
	g.generateStatement(i.ThenBranch)
	g.indentLevel--

	if i.ElseBranch != nil {
		g.writeLine("} else {")
		g.indentLevel++
		g.generateStatement(i.ElseBranch)
		g.indentLevel--
	}

	g.writeLine("}")
}

func (g *Generator) generateForStmt(f *ast.ForStmt) {
	if f.IsRange {
		// Convert "for x in range(n)" to "for x := 0; x < n; x++"
		if callExpr, ok := f.RangeExpr.(*ast.CallExpr); ok {
			if ident, ok := callExpr.Function.(*ast.Identifier); ok && ident.Value == "range" {
				if len(callExpr.Arguments) > 0 {
					g.writeLine(fmt.Sprintf("for %s := 0; %s < %s; %s++ {",
						f.RangeVar, f.RangeVar, g.generateExpression(callExpr.Arguments[0]), f.RangeVar))
				}
			}
		} else {
			// Regular range over slice/map
			g.writeLine(fmt.Sprintf("for %s := range %s {", f.RangeVar, g.generateExpression(f.RangeExpr)))
		}
	} else {
		// Traditional for loop
		init := ""
		if f.Init != nil {
			init = g.generateStatementInline(f.Init)
		}
		condition := ""
		if f.Condition != nil {
			condition = g.generateExpression(f.Condition)
		}
		update := ""
		if f.Update != nil {
			update = g.generateStatementInline(f.Update)
		}
		g.writeLine(fmt.Sprintf("for %s; %s; %s {", init, condition, update))
	}

	g.indentLevel++
	g.generateBlockStmt(f.Body)
	g.indentLevel--
	g.writeLine("}")
}

func (g *Generator) generateWhileStmt(w *ast.WhileStmt) {
	g.writeLine(fmt.Sprintf("for %s {", g.generateExpression(w.Condition)))
	g.indentLevel++
	g.generateBlockStmt(w.Body)
	g.indentLevel--
	g.writeLine("}")
}

func (g *Generator) generateReturnStmt(r *ast.ReturnStmt) {
	if r.Value != nil {
		g.writeLine(fmt.Sprintf("return %s", g.generateExpression(r.Value)))
	} else {
		g.writeLine("return")
	}
}

func (g *Generator) generateExpressionStmt(e *ast.ExpressionStmt) {
	g.writeLine(g.generateExpression(e.Expression))
}

func (g *Generator) generateBlockStmt(b *ast.BlockStmt) {
	for _, stmt := range b.Statements {
		g.generateStatement(stmt)
	}
}

func (g *Generator) generateExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.Identifier:
		return e.Value
	case *ast.Literal:
		return g.generateLiteral(e)
	case *ast.BinaryExpr:
		return g.generateBinaryExpr(e)
	case *ast.UnaryExpr:
		return g.generateUnaryExpr(e)
	case *ast.CallExpr:
		return g.generateCallExpr(e)
	case *ast.ArrayLiteral:
		return g.generateArrayLiteral(e)
	case *ast.MapLiteral:
		return g.generateMapLiteral(e)
	case *ast.IndexExpr:
		return g.generateIndexExpr(e)
	case *ast.SelectorExpr:
		return g.generateSelectorExpr(e)
	default:
		return ""
	}
}

func (g *Generator) generateLiteral(l *ast.Literal) string {
	switch l.Type {
	case "string":
		return fmt.Sprintf(`"%s"`, l.Value)
	case "int", "float":
		return fmt.Sprintf("%v", l.Value)
	case "bool":
		if l.Value.(bool) {
			return "true"
		}
		return "false"
	case "nil":
		return "nil"
	default:
		return fmt.Sprintf("%v", l.Value)
	}
}

func (g *Generator) generateBinaryExpr(b *ast.BinaryExpr) string {
	// Convert Go-script operators to Go operators
	operator := b.Operator
	switch operator {
	case "and":
		operator = "&&"
	case "or":
		operator = "||"
	case "**":
		// Power operator - need to use math.Pow
		return fmt.Sprintf("math.Pow(%s, %s)", g.generateExpression(b.Left), g.generateExpression(b.Right))
	}

	return fmt.Sprintf("(%s %s %s)", g.generateExpression(b.Left), operator, g.generateExpression(b.Right))
}

func (g *Generator) generateUnaryExpr(u *ast.UnaryExpr) string {
	operator := u.Operator
	if operator == "not" {
		operator = "!"
	}
	return fmt.Sprintf("%s%s", operator, g.generateExpression(u.Operand))
}

func (g *Generator) generateCallExpr(c *ast.CallExpr) string {
	var args []string
	for _, arg := range c.Arguments {
		args = append(args, g.generateExpression(arg))
	}

	// Handle special functions
	if ident, ok := c.Function.(*ast.Identifier); ok {
		switch ident.Value {
		case "print":
			return fmt.Sprintf("fmt.Println(%s)", strings.Join(args, ", "))
		case "len":
			return fmt.Sprintf("len(%s)", strings.Join(args, ", "))
		case "range":
			// This should be handled in for loop context
			if len(args) > 0 {
				return args[0]
			}
		}
	}

	return fmt.Sprintf("%s(%s)", g.generateExpression(c.Function), strings.Join(args, ", "))
}

func (g *Generator) generateArrayLiteral(a *ast.ArrayLiteral) string {
	var elements []string
	for _, elem := range a.Elements {
		elements = append(elements, g.generateExpression(elem))
	}
	return fmt.Sprintf("[]interface{}{%s}", strings.Join(elements, ", "))
}

func (g *Generator) generateMapLiteral(m *ast.MapLiteral) string {
	var pairs []string
	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", g.generateExpression(pair.Key), g.generateExpression(pair.Value)))
	}
	return fmt.Sprintf("map[interface{}]interface{}{%s}", strings.Join(pairs, ", "))
}

func (g *Generator) generateIndexExpr(i *ast.IndexExpr) string {
	return fmt.Sprintf("%s[%s]", g.generateExpression(i.Object), g.generateExpression(i.Index))
}

func (g *Generator) generateSelectorExpr(s *ast.SelectorExpr) string {
	return fmt.Sprintf("%s.%s", g.generateExpression(s.Object), s.Selector)
}

func (g *Generator) generateParameter(p *ast.Parameter) string {
	if p.Type != nil {
		return fmt.Sprintf("%s %s", p.Name, g.generateTypeSpec(p.Type))
	}
	return p.Name
}

func (g *Generator) generateTypeSpec(t *ast.TypeSpec) string {
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
		result += fmt.Sprintf("map[%s]%s", g.generateTypeSpec(t.KeyType), g.generateTypeSpec(t.ValueType))
	} else if t.ValueType != nil {
		result += g.generateTypeSpec(t.ValueType)
	} else {
		result += t.Name
	}
	return result
}

func (g *Generator) generateStatementInline(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case *ast.VarDecl:
		if s.Type != nil {
			return fmt.Sprintf("var %s %s = %s", s.Name, g.generateTypeSpec(s.Type), g.generateExpression(s.Value))
		}
		return fmt.Sprintf("%s := %s", s.Name, g.generateExpression(s.Value))
	case *ast.ExpressionStmt:
		return g.generateExpression(s.Expression)
	default:
		return ""
	}
}

func (g *Generator) writeLine(line string) {
	if line == "" {
		g.output.WriteString("\n")
		return
	}

	// Add indentation
	for i := 0; i < g.indentLevel; i++ {
		g.output.WriteString("\t")
	}
	g.output.WriteString(line)
	g.output.WriteString("\n")
}
