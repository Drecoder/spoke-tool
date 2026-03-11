package analyzers

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"

	"example.com/spoke-tool/api/types"
)

// CodeAnalyzer provides language-specific code analysis
// This is a reusable package for analyzing code structure
type CodeAnalyzer struct {
	// No state - pure functions
}

// NewCodeAnalyzer creates a new code analyzer
func NewCodeAnalyzer() *CodeAnalyzer {
	return &CodeAnalyzer{}
}

// AnalysisResult represents the result of code analysis
type AnalysisResult struct {
	Language    types.Language  `json:"language"`
	Functions   []*FunctionInfo `json:"functions"`
	Types       []*TypeInfo     `json:"types"`
	Imports     []string        `json:"imports"`
	Complexity  int             `json:"complexity"`
	LinesOfCode int             `json:"lines_of_code"`
}

// FunctionInfo represents information about a function
type FunctionInfo struct {
	Name       string       `json:"name"`
	Signature  string       `json:"signature"`
	Parameters []*ParamInfo `json:"parameters"`
	Returns    []string     `json:"returns"`
	LineStart  int          `json:"line_start"`
	LineEnd    int          `json:"line_end"`
	IsExported bool         `json:"is_exported"`
	IsMethod   bool         `json:"is_method"`
	Receiver   string       `json:"receiver,omitempty"`
	Complexity int          `json:"complexity"`
	Doc        string       `json:"doc,omitempty"`
}

// ParamInfo represents function parameter information
type ParamInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// TypeInfo represents type/struct/class information
type TypeInfo struct {
	Name       string          `json:"name"`
	Kind       string          `json:"kind"` // struct, interface, class
	Fields     []*FieldInfo    `json:"fields"`
	Methods    []*FunctionInfo `json:"methods"`
	LineStart  int             `json:"line_start"`
	LineEnd    int             `json:"line_end"`
	IsExported bool            `json:"is_exported"`
	Doc        string          `json:"doc,omitempty"`
}

// FieldInfo represents struct/class field information
type FieldInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Tags string `json:"tags,omitempty"`
}

// AnalyzeGo analyzes Go code and returns structured information
func (a *CodeAnalyzer) AnalyzeGo(content string) (*AnalysisResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	result := &AnalysisResult{
		Language:  types.Go,
		Functions: []*FunctionInfo{},
		Types:     []*TypeInfo{},
		Imports:   []string{},
	}

	// Collect imports
	for _, imp := range node.Imports {
		if imp.Path != nil {
			result.Imports = append(result.Imports, strings.Trim(imp.Path.Value, "\""))
		}
	}

	// Analyze declarations
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			// Function declaration
			fn := a.parseGoFunction(d, fset)
			result.Functions = append(result.Functions, fn)
			result.Complexity += fn.Complexity

		case *ast.GenDecl:
			// Type, var, const declarations
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						typeInfo := a.parseGoType(typeSpec, fset)
						if typeInfo != nil {
							result.Types = append(result.Types, typeInfo)
						}
					}
				}
			}
		}
	}

	// Count lines of code
	result.LinesOfCode = strings.Count(content, "\n") + 1

	return result, nil
}

// AnalyzeNodeJS analyzes Node.js/JavaScript code (simplified regex-based)
// For production, consider using @babel/parser or similar
func (a *CodeAnalyzer) AnalyzeNodeJS(content string) (*AnalysisResult, error) {
	result := &AnalysisResult{
		Language:  types.NodeJS,
		Functions: []*FunctionInfo{},
		Types:     []*TypeInfo{},
		Imports:   []string{},
	}

	// Extract imports (simplified)
	importRegex := regexp.MustCompile(`(?:import|require)\s*\(?['"]([^'"]+)['"]`)
	importMatches := importRegex.FindAllStringSubmatch(content, -1)
	for _, match := range importMatches {
		if len(match) > 1 {
			result.Imports = append(result.Imports, match[1])
		}
	}

	// Extract functions (simplified)
	funcRegex := regexp.MustCompile(`(?:function\s+)?(\w+)\s*\(([^)]*)\)\s*{`)
	funcMatches := funcRegex.FindAllStringSubmatchIndex(content, -1)

	for _, match := range funcMatches {
		if len(match) < 4 {
			continue
		}

		funcName := content[match[2]:match[3]]

		// Skip if it looks like a test
		if strings.HasPrefix(funcName, "test") || strings.HasPrefix(funcName, "it") {
			continue
		}

		// Extract parameters
		paramStr := ""
		if len(match) > 5 {
			paramStr = content[match[4]:match[5]]
		}
		params := a.parseNodeJSParams(paramStr)

		// Find end of function (simplified)
		start := match[0]
		end := a.findBlockEnd(content, start)

		fn := &FunctionInfo{
			Name:       funcName,
			Signature:  content[start:match[1]],
			Parameters: params,
			LineStart:  strings.Count(content[:start], "\n") + 1,
			LineEnd:    strings.Count(content[:end], "\n") + 1,
			IsExported: !strings.HasPrefix(funcName, "_"),
			Complexity: a.estimateNodeJSComplexity(content[start:end]),
		}

		result.Functions = append(result.Functions, fn)
		result.Complexity += fn.Complexity
	}

	// Extract classes (simplified)
	classRegex := regexp.MustCompile(`class\s+(\w+)`)
	classMatches := classRegex.FindAllStringSubmatchIndex(content, -1)

	for _, match := range classMatches {
		if len(match) < 4 {
			continue
		}

		className := content[match[2]:match[3]]

		typeInfo := &TypeInfo{
			Name:       className,
			Kind:       "class",
			Fields:     []*FieldInfo{},
			Methods:    []*FunctionInfo{},
			IsExported: true,
		}

		result.Types = append(result.Types, typeInfo)
	}

	result.LinesOfCode = strings.Count(content, "\n") + 1
	return result, nil
}

// AnalyzePython analyzes Python code (simplified regex-based)
// For production, consider using the ast module via python AST parser
func (a *CodeAnalyzer) AnalyzePython(content string) (*AnalysisResult, error) {
	result := &AnalysisResult{
		Language:  types.Python,
		Functions: []*FunctionInfo{},
		Types:     []*TypeInfo{},
		Imports:   []string{},
	}

	// Extract imports
	importRegex := regexp.MustCompile(`(?:import|from)\s+(\w+)`)
	importMatches := importRegex.FindAllStringSubmatch(content, -1)
	for _, match := range importMatches {
		if len(match) > 1 {
			result.Imports = append(result.Imports, match[1])
		}
	}

	// Extract functions
	funcRegex := regexp.MustCompile(`def\s+(\w+)\s*\(([^)]*)\):`)
	funcMatches := funcRegex.FindAllStringSubmatchIndex(content, -1)

	for _, match := range funcMatches {
		if len(match) < 4 {
			continue
		}

		funcName := content[match[2]:match[3]]

		// Skip test functions
		if strings.HasPrefix(funcName, "test_") {
			continue
		}

		// Extract parameters
		paramStr := ""
		if len(match) > 5 {
			paramStr = content[match[4]:match[5]]
		}
		params := a.parsePythonParams(paramStr)

		// Find end of function (by indentation)
		start := match[0]
		end := a.findPythonFunctionEnd(content, start)

		fn := &FunctionInfo{
			Name:       funcName,
			Signature:  content[start:match[1]],
			Parameters: params,
			LineStart:  strings.Count(content[:start], "\n") + 1,
			LineEnd:    strings.Count(content[:end], "\n") + 1,
			IsExported: !strings.HasPrefix(funcName, "_"),
			Complexity: a.estimatePythonComplexity(content[start:end]),
		}

		result.Functions = append(result.Functions, fn)
		result.Complexity += fn.Complexity
	}

	// Extract classes
	classRegex := regexp.MustCompile(`class\s+(\w+)`)
	classMatches := classRegex.FindAllStringSubmatchIndex(content, -1)

	for _, match := range classMatches {
		if len(match) < 4 {
			continue
		}

		className := content[match[2]:match[3]]

		typeInfo := &TypeInfo{
			Name:       className,
			Kind:       "class",
			Fields:     []*FieldInfo{},
			Methods:    []*FunctionInfo{},
			IsExported: !strings.HasPrefix(className, "_"),
		}

		result.Types = append(result.Types, typeInfo)
	}

	result.LinesOfCode = strings.Count(content, "\n") + 1
	return result, nil
}

// Helper methods for Go analysis

func (a *CodeAnalyzer) parseGoFunction(fn *ast.FuncDecl, fset *token.FileSet) *FunctionInfo {
	info := &FunctionInfo{
		Name:       fn.Name.Name,
		IsExported: fn.Name.IsExported(),
		IsMethod:   fn.Recv != nil,
		Parameters: []*ParamInfo{},
		Returns:    []string{},
	}

	// Get position
	info.LineStart = fset.Position(fn.Pos()).Line
	info.LineEnd = fset.Position(fn.End()).Line

	// Get receiver if method
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		recv := fn.Recv.List[0]
		// Extract receiver type
		if star, ok := recv.Type.(*ast.StarExpr); ok {
			if ident, ok := star.X.(*ast.Ident); ok {
				info.Receiver = "*" + ident.Name
			}
		} else if ident, ok := recv.Type.(*ast.Ident); ok {
			info.Receiver = ident.Name
		}
	}

	// Parse parameters
	if fn.Type.Params != nil {
		for _, param := range fn.Type.Params.List {
			typeStr := a.exprToString(param.Type)
			for _, name := range param.Names {
				info.Parameters = append(info.Parameters, &ParamInfo{
					Name: name.Name,
					Type: typeStr,
				})
			}
			// Handle unnamed parameters
			if len(param.Names) == 0 {
				info.Parameters = append(info.Parameters, &ParamInfo{
					Name: "",
					Type: typeStr,
				})
			}
		}
	}

	// Parse returns
	if fn.Type.Results != nil {
		for _, result := range fn.Type.Results.List {
			typeStr := a.exprToString(result.Type)
			if len(result.Names) > 0 {
				for range result.Names {
					info.Returns = append(info.Returns, typeStr)
				}
			} else {
				info.Returns = append(info.Returns, typeStr)
			}
		}
	}

	// Estimate complexity
	if fn.Body != nil {
		info.Complexity = a.estimateGoComplexity(fn.Body)
	}

	// Get doc comment
	if fn.Doc != nil {
		info.Doc = fn.Doc.Text()
	}

	// Build signature
	info.Signature = a.buildGoSignature(fn)

	return info
}

func (a *CodeAnalyzer) parseGoType(typeSpec *ast.TypeSpec, fset *token.FileSet) *TypeInfo {
	typeInfo := &TypeInfo{
		Name:       typeSpec.Name.Name,
		IsExported: typeSpec.Name.IsExported(),
		Fields:     []*FieldInfo{},
		Methods:    []*FunctionInfo{},
	}

	// Get position
	typeInfo.LineStart = fset.Position(typeSpec.Pos()).Line
	typeInfo.LineEnd = fset.Position(typeSpec.End()).Line

	// Determine kind and parse fields
	switch t := typeSpec.Type.(type) {
	case *ast.StructType:
		typeInfo.Kind = "struct"
		if t.Fields != nil {
			for _, field := range t.Fields.List {
				fieldInfo := &FieldInfo{
					Type: a.exprToString(field.Type),
				}
				// Get field names
				if len(field.Names) > 0 {
					for _, name := range field.Names {
						fieldInfo.Name = name.Name
						typeInfo.Fields = append(typeInfo.Fields, fieldInfo)
					}
				} else {
					// Embedded field
					fieldInfo.Name = a.exprToString(field.Type)
					typeInfo.Fields = append(typeInfo.Fields, fieldInfo)
				}

				// Get struct tags
				if field.Tag != nil {
					fieldInfo.Tags = field.Tag.Value
				}
			}
		}

	case *ast.InterfaceType:
		typeInfo.Kind = "interface"
		// Could parse methods here

	case *ast.Ident:
		typeInfo.Kind = "alias"
	}

	// Get doc comment
	if typeSpec.Doc != nil {
		typeInfo.Doc = typeSpec.Doc.Text()
	}

	return typeInfo
}

func (a *CodeAnalyzer) exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + a.exprToString(e.X)
	case *ast.ArrayType:
		return "[]" + a.exprToString(e.Elt)
	case *ast.MapType:
		return "map[" + a.exprToString(e.Key) + "]" + a.exprToString(e.Value)
	case *ast.SelectorExpr:
		return a.exprToString(e.X) + "." + e.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func"
	default:
		return "unknown"
	}
}

func (a *CodeAnalyzer) estimateGoComplexity(body *ast.BlockStmt) int {
	complexity := 1

	ast.Inspect(body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt,
			*ast.CaseClause, *ast.CommClause,
			*ast.BinaryExpr: // For logical operators (&&, ||)
			complexity++
		}
		return true
	})

	return complexity
}

func (a *CodeAnalyzer) buildGoSignature(fn *ast.FuncDecl) string {
	var sb strings.Builder

	sb.WriteString("func ")
	if fn.Recv != nil {
		sb.WriteString("(")
		// Add receiver (simplified)
		if len(fn.Recv.List) > 0 {
			recv := fn.Recv.List[0]
			if len(recv.Names) > 0 {
				sb.WriteString(recv.Names[0].Name)
			}
			sb.WriteString(" ")
			sb.WriteString(a.exprToString(recv.Type))
		}
		sb.WriteString(") ")
	}
	sb.WriteString(fn.Name.Name)
	sb.WriteString("(")

	// Add parameters
	if fn.Type.Params != nil {
		for i, param := range fn.Type.Params.List {
			if i > 0 {
				sb.WriteString(", ")
			}
			typeStr := a.exprToString(param.Type)
			if len(param.Names) > 0 {
				for j, name := range param.Names {
					if j > 0 {
						sb.WriteString(", ")
					}
					sb.WriteString(name.Name)
				}
				sb.WriteString(" ")
				sb.WriteString(typeStr)
			} else {
				sb.WriteString(typeStr)
			}
		}
	}
	sb.WriteString(")")

	// Add returns
	if fn.Type.Results != nil {
		sb.WriteString(" ")
		if len(fn.Type.Results.List) > 1 {
			sb.WriteString("(")
		}
		for i, result := range fn.Type.Results.List {
			if i > 0 {
				sb.WriteString(", ")
			}
			typeStr := a.exprToString(result.Type)
			sb.WriteString(typeStr)
		}
		if len(fn.Type.Results.List) > 1 {
			sb.WriteString(")")
		}
	}

	return sb.String()
}

// Helper methods for Node.js analysis

func (a *CodeAnalyzer) parseNodeJSParams(paramStr string) []*ParamInfo {
	if paramStr == "" {
		return nil
	}

	var params []*ParamInfo
	parts := strings.Split(paramStr, ",")

	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p != "" {
			// Handle default values
			if strings.Contains(p, "=") {
				p = strings.TrimSpace(strings.Split(p, "=")[0])
			}
			params = append(params, &ParamInfo{
				Name: p,
				Type: "any", // JavaScript is dynamic
			})
		}
	}

	return params
}

func (a *CodeAnalyzer) estimateNodeJSComplexity(code string) int {
	complexity := 1

	// Count control structures
	controlStructures := []string{
		"if", "else if", "for", "while", "switch",
		"&&", "||", "case", "catch",
	}

	for _, cs := range controlStructures {
		count := strings.Count(code, cs)
		complexity += count
	}

	return complexity
}

// Helper methods for Python analysis

func (a *CodeAnalyzer) parsePythonParams(paramStr string) []*ParamInfo {
	if paramStr == "" || paramStr == "self" {
		return nil
	}

	var params []*ParamInfo
	parts := strings.Split(paramStr, ",")

	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p != "" && p != "self" {
			// Handle default values
			if strings.Contains(p, "=") {
				p = strings.TrimSpace(strings.Split(p, "=")[0])
			}
			// Handle type hints (simplified)
			if strings.Contains(p, ":") {
				p = strings.TrimSpace(strings.Split(p, ":")[0])
			}
			params = append(params, &ParamInfo{
				Name: p,
				Type: "Any", // Python is dynamic
			})
		}
	}

	return params
}

func (a *CodeAnalyzer) estimatePythonComplexity(code string) int {
	complexity := 1

	// Count control structures
	controlStructures := []string{
		"if ", "elif ", "else:", "for ", "while ",
		"and ", "or ", "except:", "with ",
	}

	for _, cs := range controlStructures {
		count := strings.Count(code, cs)
		complexity += count
	}

	return complexity
}

// Common helper methods

func (a *CodeAnalyzer) findBlockEnd(content string, start int) int {
	braceCount := 0
	inString := false
	escape := false

	for i := start; i < len(content); i++ {
		c := content[i]

		if escape {
			escape = false
			continue
		}

		if c == '\\' && inString {
			escape = true
			continue
		}

		if c == '"' || c == '\'' || c == '`' {
			inString = !inString
			continue
		}

		if !inString {
			if c == '{' {
				braceCount++
			} else if c == '}' {
				braceCount--
				if braceCount == 0 {
					return i + 1
				}
			}
		}
	}

	return len(content)
}

func (a *CodeAnalyzer) findPythonFunctionEnd(content string, start int) int {
	lines := strings.Split(content[start:], "\n")
	if len(lines) == 0 {
		return len(content)
	}

	// Get indentation of first line
	firstLine := lines[0]
	indent := len(firstLine) - len(strings.TrimLeft(firstLine, " "))

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		currentIndent := len(line) - len(strings.TrimLeft(line, " "))
		if currentIndent <= indent && trimmed != "" {
			// Found end of function
			return start + len(strings.Join(lines[:i], "\n"))
		}
	}

	return len(content)
}

// GetUntestedFunctions identifies functions that don't have corresponding tests
func (a *CodeAnalyzer) GetUntestedFunctions(functions []*FunctionInfo, testFunctions []string) []*FunctionInfo {
	var untested []*FunctionInfo

	for _, fn := range functions {
		// Skip if it's a test helper or internal
		if !fn.IsExported {
			continue
		}

		// Check if function has a test
		hasTest := false
		for _, tf := range testFunctions {
			if tf == "Test"+fn.Name || tf == "test_"+fn.Name || tf == fn.Name {
				hasTest = true
				break
			}
		}

		if !hasTest {
			untested = append(untested, fn)
		}
	}

	return untested
}
