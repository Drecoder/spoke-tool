package extractors

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"

	"example.com/spoke-tool/api/types"
	apitypes "example.com/spoke-tool/api/types"
)

// SignatureExtractor extracts function and type signatures from code
// This is PURELY EXTRACTIVE - no modifications
type SignatureExtractor struct {
	// No state - pure functions
}

// FunctionSignature represents a complete function signature
type FunctionSignature struct {
	// Function name
	Name string `json:"name"`

	// Language of the function
	Language types.Language `json:"language"`

	// Full signature string
	Signature string `json:"signature"`

	// Parameters
	Parameters []*Parameter `json:"parameters"`

	// Return types
	Returns []*ReturnType `json:"returns"`

	// Receiver (for methods)
	Receiver *Receiver `json:"receiver,omitempty"`

	// Type parameters (generics)
	TypeParams []*TypeParameter `json:"type_params,omitempty"`

	// Whether the function is exported/public
	IsExported bool `json:"is_exported"`

	// File where function is defined
	FilePath string `json:"file_path"`

	// Line number
	Line int `json:"line"`

	// Documentation comment
	Doc string `json:"doc,omitempty"`
}

// Parameter represents a function parameter
type Parameter struct {
	// Parameter name
	Name string `json:"name"`

	// Parameter type
	Type string `json:"type"`

	// Whether the parameter is variadic
	IsVariadic bool `json:"is_variadic,omitempty"`

	// Default value (if any)
	DefaultValue string `json:"default_value,omitempty"`

	// Description (from comments)
	Description string `json:"description,omitempty"`
}

// ReturnType represents a return type
type ReturnType struct {
	// Type name
	Type string `json:"type"`

	// Name (if named return)
	Name string `json:"name,omitempty"`

	// Description (from comments)
	Description string `json:"description,omitempty"`
}

// Receiver represents a method receiver
type Receiver struct {
	// Receiver name
	Name string `json:"name"`

	// Receiver type
	Type string `json:"type"`

	// Whether it's a pointer receiver
	IsPointer bool `json:"is_pointer"`
}

// TypeParameter represents a generic type parameter
type TypeParameter struct {
	// Parameter name
	Name string `json:"name"`

	// Constraint
	Constraint string `json:"constraint,omitempty"`
}

// TypeSignature represents a type/struct/interface signature
type TypeSignature struct {
	// Type name
	Name string `json:"name"`

	// Kind (struct, interface, class, etc.)
	Kind string `json:"kind"`

	// Language
	Language apitypes.Language `json:"language"`

	// Fields (for structs/classes)
	Fields []*Field `json:"fields,omitempty"`

	// Methods
	Methods []*FunctionSignature `json:"methods,omitempty"`

	// Implemented interfaces
	Implements []string `json:"implements,omitempty"`

	// Whether the type is exported
	IsExported bool `json:"is_exported"`

	// File path
	FilePath string `json:"file_path"`

	// Line number
	Line int `json:"line"`

	// Documentation
	Doc string `json:"doc,omitempty"`
}

// Field represents a struct/class field
type Field struct {
	// Field name
	Name string `json:"name"`

	// Field type
	Type string `json:"type"`

	// Tags (like Go struct tags)
	Tags map[string]string `json:"tags,omitempty"`

	// Whether the field is exported
	IsExported bool `json:"is_exported"`

	// Default value (if any)
	DefaultValue string `json:"default_value,omitempty"`

	// Description
	Description string `json:"description,omitempty"`
}

// NewSignatureExtractor creates a new signature extractor
func NewSignatureExtractor() *SignatureExtractor {
	return &SignatureExtractor{}
}

// ExtractFromGo extracts signatures from Go code
func (s *SignatureExtractor) ExtractFromGo(content string, filePath string) ([]*FunctionSignature, []*TypeSignature, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	var functions []*FunctionSignature
	var types []*TypeSignature

	// Process declarations
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			// Function declaration
			fn := s.extractGoFunction(d, fset, filePath)
			if fn != nil {
				functions = append(functions, fn)
			}

		case *ast.GenDecl:
			// Type, var, const declarations
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						ts := s.extractGoType(typeSpec, fset, filePath, node)
						if ts != nil {
							types = append(types, ts)
						}
					}
				}
			}
		}
	}

	return functions, types, nil
}

// ExtractFromNodeJS extracts signatures from Node.js/JavaScript code
func (s *SignatureExtractor) ExtractFromNodeJS(content string, filePath string) ([]*FunctionSignature, []*TypeSignature, error) {
	fmt.Printf("Types package: %T\n", apitypes.NodeJS)
	var functions []*FunctionSignature
	var types []*TypeSignature

	lines := strings.Split(content, "\n")

	// Extract functions
	funcRegex := regexp.MustCompile(`(?:function\s+)?(\w+)\s*\(([^)]*)\)\s*(?::\s*(\w+))?\s*{`)
	classRegex := regexp.MustCompile(`class\s+(\w+)(?:\s+extends\s+(\w+))?\s*{`)
	methodRegex := regexp.MustCompile(`(\w+)\s*\(([^)]*)\)\s*(?::\s*(\w+))?\s*{`)

	for i, line := range lines {
		lineNum := i + 1

		// Extract functions
		if matches := funcRegex.FindStringSubmatch(line); len(matches) > 1 {
			fn := &FunctionSignature{
				Name:       matches[1],
				Language:   apitypes.NodeJS,
				Parameters: []*Parameter{},
				Returns:    []*ReturnType{},
				FilePath:   filePath,
				Line:       lineNum,
				IsExported: !strings.HasPrefix(matches[1], "_"),
			}

			// Parse parameters
			if len(matches) > 2 {
				fn.Parameters = s.parseNodeJSParams(matches[2])
			}

			// Parse return type
			if len(matches) > 3 && matches[3] != "" {
				fn.Returns = append(fn.Returns, &ReturnType{Type: matches[3]})
			}

			// Build full signature
			fn.Signature = s.buildNodeJSSignature(fn)

			functions = append(functions, fn)
		}

		// Extract classes
		if matches := classRegex.FindStringSubmatch(line); len(matches) > 1 {
			ts := &TypeSignature{
				Name:       matches[1],
				Kind:       "class",
				Language:   apitypes.NodeJS,
				Fields:     []*Field{},
				Methods:    []*FunctionSignature{},
				FilePath:   filePath,
				Line:       lineNum,
				IsExported: true,
			}

			// Add inheritance
			if len(matches) > 2 && matches[2] != "" {
				ts.Implements = append(ts.Implements, matches[2])
			}

			types = append(types, ts)
		}

		// Extract methods (inside classes)
		if strings.Contains(line, "class") {
			continue // Skip class lines
		}

		if matches := methodRegex.FindStringSubmatch(line); len(matches) > 1 {
			// This could be a method - would need context to know if it's inside a class
			// For now, just note it as a potential method
		}
	}

	return functions, types, nil
}

// ExtractFromPython extracts signatures from Python code
func (s *SignatureExtractor) ExtractFromPython(content string, filePath string) ([]*FunctionSignature, []*TypeSignature, error) {
	var functions []*FunctionSignature
	var types []*TypeSignature

	lines := strings.Split(content, "\n")

	// Extract functions
	funcRegex := regexp.MustCompile(`def\s+(\w+)\s*\(([^)]*)\)(?:\s*->\s*([^:]+))?\s*:`)
	classRegex := regexp.MustCompile(`class\s+(\w+)(?:\s*\(([^)]+)\))?\s*:`)

	for i, line := range lines {
		lineNum := i + 1

		// Extract functions
		if matches := funcRegex.FindStringSubmatch(line); len(matches) > 1 {
			fn := &FunctionSignature{
				Name:       matches[1],
				Language:   apitypes.Python,
				Parameters: []*Parameter{},
				Returns:    []*ReturnType{},
				FilePath:   filePath,
				Line:       lineNum,
				IsExported: !strings.HasPrefix(matches[1], "_"),
			}

			// Parse parameters
			if len(matches) > 2 {
				fn.Parameters = s.parsePythonParams(matches[2])
			}

			// Parse return type
			if len(matches) > 3 && matches[3] != "" {
				fn.Returns = append(fn.Returns, &ReturnType{Type: strings.TrimSpace(matches[3])})
			}

			// Build full signature
			fn.Signature = s.buildPythonSignature(fn)

			functions = append(functions, fn)
		}

		// Extract classes
		if matches := classRegex.FindStringSubmatch(line); len(matches) > 1 {
			ts := &TypeSignature{
				Name:       matches[1],
				Kind:       "class",
				Language:   apitypes.Python,
				Fields:     []*Field{},
				Methods:    []*FunctionSignature{},
				FilePath:   filePath,
				Line:       lineNum,
				IsExported: !strings.HasPrefix(matches[1], "_"),
			}

			// Add inheritance
			if len(matches) > 2 && matches[2] != "" {
				parents := strings.Split(matches[2], ",")
				for _, p := range parents {
					ts.Implements = append(ts.Implements, strings.TrimSpace(p))
				}
			}

			types = append(types, ts)
		}
	}

	return functions, types, nil
}

// ExtractSignatures automatically detects language and extracts signatures
func (s *SignatureExtractor) ExtractSignatures(content string, filePath string, lang apitypes.Language) ([]*FunctionSignature, []*TypeSignature, error) {
	switch lang {
	case apitypes.Go:
		return s.ExtractFromGo(content, filePath)
	case apitypes.NodeJS:
		return s.ExtractFromNodeJS(content, filePath)
	case apitypes.Python:
		return s.ExtractFromPython(content, filePath)
	default:
		return nil, nil, fmt.Errorf("unsupported language: %s", lang)
	}
}

// Go-specific extraction methods

func (s *SignatureExtractor) extractGoFunction(fn *ast.FuncDecl, fset *token.FileSet, filePath string) *FunctionSignature {
	sig := &FunctionSignature{
		Name:       fn.Name.Name,
		Language:   apitypes.Go,
		Parameters: []*Parameter{},
		Returns:    []*ReturnType{},
		FilePath:   filePath,
		Line:       fset.Position(fn.Pos()).Line,
		IsExported: fn.Name.IsExported(),
	}

	// Extract receiver (for methods)
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		recv := fn.Recv.List[0]
		sig.Receiver = s.extractGoReceiver(recv)
	}

	// Extract parameters
	if fn.Type.Params != nil {
		for _, param := range fn.Type.Params.List {
			params := s.extractGoParameters(param)
			sig.Parameters = append(sig.Parameters, params...)
		}
	}

	// Extract return types
	if fn.Type.Results != nil {
		for _, result := range fn.Type.Results.List {
			returns := s.extractGoReturns(result)
			sig.Returns = append(sig.Returns, returns...)
		}
	}

	// Extract type parameters (generics)
	if fn.Type.TypeParams != nil {
		for _, param := range fn.Type.TypeParams.List {
			typeParams := s.extractGoTypeParams(param)
			sig.TypeParams = append(sig.TypeParams, typeParams...)
		}
	}

	// Build signature string
	sig.Signature = s.buildGoSignature(sig, fn)

	// Extract doc comment
	if fn.Doc != nil {
		sig.Doc = fn.Doc.Text()
	}

	return sig
}

func (s *SignatureExtractor) extractGoType(typeSpec *ast.TypeSpec, fset *token.FileSet, filePath string, node *ast.File) *TypeSignature {
	ts := &TypeSignature{
		Name:       typeSpec.Name.Name,
		Language:   apitypes.Go,
		Fields:     []*Field{},
		Methods:    []*FunctionSignature{},
		FilePath:   filePath,
		Line:       fset.Position(typeSpec.Pos()).Line,
		IsExported: typeSpec.Name.IsExported(),
	}

	// Determine kind and extract fields
	switch t := typeSpec.Type.(type) {
	case *ast.StructType:
		ts.Kind = "struct"
		if t.Fields != nil {
			for _, field := range t.Fields.List {
				fields := s.extractGoFields(field)
				ts.Fields = append(ts.Fields, fields...)
			}
		}

	case *ast.InterfaceType:
		ts.Kind = "interface"
		// Would extract methods here

	case *ast.Ident:
		ts.Kind = "alias"
		ts.Implements = append(ts.Implements, t.Name)

	case *ast.ArrayType:
		ts.Kind = "array"
		ts.Implements = append(ts.Implements, s.exprToString(t.Elt))

	case *ast.MapType:
		ts.Kind = "map"
	}

	// Find methods for this type
	if ts.Kind == "struct" || ts.Kind == "interface" {
		// Methods would be found by scanning all functions
		// This would need to be done separately
	}

	// Extract doc comment
	if typeSpec.Doc != nil {
		ts.Doc = typeSpec.Doc.Text()
	} else if node != nil && node.Doc != nil {
		// Could be attached to the type declaration
	}

	return ts
}

func (s *SignatureExtractor) extractGoReceiver(field *ast.Field) *Receiver {
	recv := &Receiver{}

	// Get receiver name
	if len(field.Names) > 0 {
		recv.Name = field.Names[0].Name
	}

	// Get receiver type
	switch t := field.Type.(type) {
	case *ast.StarExpr:
		recv.IsPointer = true
		if ident, ok := t.X.(*ast.Ident); ok {
			recv.Type = ident.Name
		}
	case *ast.Ident:
		recv.IsPointer = false
		recv.Type = t.Name
	}

	return recv
}

func (s *SignatureExtractor) extractGoParameters(field *ast.Field) []*Parameter {
	var params []*Parameter

	typeStr := s.exprToString(field.Type)
	isVariadic := strings.HasPrefix(typeStr, "...")

	if len(field.Names) > 0 {
		for _, name := range field.Names {
			params = append(params, &Parameter{
				Name:       name.Name,
				Type:       typeStr,
				IsVariadic: isVariadic,
			})
		}
	} else {
		// Unnamed parameter
		params = append(params, &Parameter{
			Name:       "",
			Type:       typeStr,
			IsVariadic: isVariadic,
		})
	}

	return params
}

func (s *SignatureExtractor) extractGoReturns(field *ast.Field) []*ReturnType {
	var returns []*ReturnType

	typeStr := s.exprToString(field.Type)

	if len(field.Names) > 0 {
		for _, name := range field.Names {
			returns = append(returns, &ReturnType{
				Name: name.Name,
				Type: typeStr,
			})
		}
	} else {
		returns = append(returns, &ReturnType{
			Type: typeStr,
		})
	}

	return returns
}

func (s *SignatureExtractor) extractGoTypeParams(field *ast.Field) []*TypeParameter {
	var params []*TypeParameter

	// Simplified - would need to parse constraints
	if len(field.Names) > 0 {
		for _, name := range field.Names {
			params = append(params, &TypeParameter{
				Name: name.Name,
			})
		}
	}

	return params
}

func (s *SignatureExtractor) extractGoFields(field *ast.Field) []*Field {
	var fields []*Field

	typeStr := s.exprToString(field.Type)

	// Extract tags
	tags := make(map[string]string)
	if field.Tag != nil {
		tagStr := strings.Trim(field.Tag.Value, "`")
		pairs := strings.Split(tagStr, " ")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, ":", 2)
			if len(kv) == 2 {
				tags[kv[0]] = strings.Trim(kv[1], "\"")
			}
		}
	}

	if len(field.Names) > 0 {
		for _, name := range field.Names {
			fields = append(fields, &Field{
				Name:       name.Name,
				Type:       typeStr,
				Tags:       tags,
				IsExported: name.IsExported(),
			})
		}
	} else {
		// Embedded field
		fields = append(fields, &Field{
			Name:       typeStr,
			Type:       typeStr,
			Tags:       tags,
			IsExported: true,
		})
	}

	return fields
}

func (s *SignatureExtractor) exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + s.exprToString(e.X)
	case *ast.ArrayType:
		return "[]" + s.exprToString(e.Elt)
	case *ast.MapType:
		return "map[" + s.exprToString(e.Key) + "]" + s.exprToString(e.Value)
	case *ast.SelectorExpr:
		return s.exprToString(e.X) + "." + e.Sel.Name
	case *ast.FuncType:
		return "func"
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.Ellipsis:
		return "..." + s.exprToString(e.Elt)
	default:
		return "unknown"
	}
}

func (s *SignatureExtractor) buildGoSignature(sig *FunctionSignature, fn *ast.FuncDecl) string {
	var sb strings.Builder

	sb.WriteString("func ")

	// Receiver
	if sig.Receiver != nil {
		sb.WriteString("(")
		if sig.Receiver.Name != "" {
			sb.WriteString(sig.Receiver.Name)
			sb.WriteString(" ")
		}
		if sig.Receiver.IsPointer {
			sb.WriteString("*")
		}
		sb.WriteString(sig.Receiver.Type)
		sb.WriteString(") ")
	}

	// Function name
	sb.WriteString(sig.Name)

	// Type parameters (generics)
	if len(sig.TypeParams) > 0 {
		sb.WriteString("[")
		for i, tp := range sig.TypeParams {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(tp.Name)
		}
		sb.WriteString("]")
	}

	// Parameters
	sb.WriteString("(")
	for i, p := range sig.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		if p.Name != "" {
			sb.WriteString(p.Name)
			sb.WriteString(" ")
		}
		sb.WriteString(p.Type)
	}
	sb.WriteString(")")

	// Returns
	if len(sig.Returns) > 0 {
		if len(sig.Returns) == 1 && sig.Returns[0].Name == "" {
			sb.WriteString(" ")
			sb.WriteString(sig.Returns[0].Type)
		} else {
			sb.WriteString(" (")
			for i, r := range sig.Returns {
				if i > 0 {
					sb.WriteString(", ")
				}
				if r.Name != "" {
					sb.WriteString(r.Name)
					sb.WriteString(" ")
				}
				sb.WriteString(r.Type)
			}
			sb.WriteString(")")
		}
	}

	return sb.String()
}

// Node.js-specific methods

func (s *SignatureExtractor) parseNodeJSParams(paramStr string) []*Parameter {
	if paramStr == "" {
		return nil
	}

	var params []*Parameter
	parts := strings.Split(paramStr, ",")

	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}

		param := &Parameter{
			Type: "any", // JavaScript is dynamic
		}

		// Handle default values
		if strings.Contains(p, "=") {
			eq := strings.SplitN(p, "=", 2)
			param.Name = strings.TrimSpace(eq[0])
			param.DefaultValue = strings.TrimSpace(eq[1])
		} else {
			param.Name = p
		}

		// Handle destructuring (simplified)
		if strings.Contains(param.Name, "{") {
			param.Name = "options"
		}

		params = append(params, param)
	}

	return params
}

func (s *SignatureExtractor) buildNodeJSSignature(fn *FunctionSignature) string {
	var sb strings.Builder

	sb.WriteString("function ")
	sb.WriteString(fn.Name)
	sb.WriteString("(")

	for i, p := range fn.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(p.Name)
		if p.DefaultValue != "" {
			sb.WriteString(" = ")
			sb.WriteString(p.DefaultValue)
		}
	}

	sb.WriteString(")")

	if len(fn.Returns) > 0 {
		sb.WriteString(": ")
		sb.WriteString(fn.Returns[0].Type)
	}

	return sb.String()
}

// Python-specific methods

func (s *SignatureExtractor) parsePythonParams(paramStr string) []*Parameter {
	if paramStr == "" || paramStr == "self" {
		return nil
	}

	var params []*Parameter
	parts := strings.Split(paramStr, ",")

	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" || p == "self" {
			continue
		}

		param := &Parameter{
			Type: "Any", // Python is dynamic
		}

		// Handle type hints
		if strings.Contains(p, ":") {
			typeParts := strings.SplitN(p, ":", 2)
			param.Name = strings.TrimSpace(typeParts[0])
			if len(typeParts) > 1 {
				// Extract type, ignoring default values for now
				typeStr := strings.TrimSpace(typeParts[1])
				if strings.Contains(typeStr, "=") {
					typeStr = strings.SplitN(typeStr, "=", 2)[0]
				}
				param.Type = strings.TrimSpace(typeStr)
			}
		} else {
			param.Name = p
		}

		// Handle default values
		if strings.Contains(p, "=") {
			eq := strings.SplitN(p, "=", 2)
			param.Name = strings.TrimSpace(eq[0])
			param.DefaultValue = strings.TrimSpace(eq[1])
		}

		params = append(params, param)
	}

	return params
}

func (s *SignatureExtractor) buildPythonSignature(fn *FunctionSignature) string {
	var sb strings.Builder

	sb.WriteString("def ")
	sb.WriteString(fn.Name)
	sb.WriteString("(")

	hasSelf := false
	for i, p := range fn.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(p.Name)

		// Add type hint if available
		if p.Type != "Any" && p.Type != "" {
			sb.WriteString(": ")
			sb.WriteString(p.Type)
		}

		if p.Name == "self" {
			hasSelf = true
		}
	}

	// Add self if it's a method and not present
	if !hasSelf && fn.Receiver != nil {
		if len(fn.Parameters) > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("self")
	}

	sb.WriteString(")")

	// Add return type hint
	if len(fn.Returns) > 0 {
		sb.WriteString(" -> ")
		sb.WriteString(fn.Returns[0].Type)
	}

	return sb.String()
}

// Utility methods

// FormatFunctionSignature formats a function signature for documentation
func (s *SignatureExtractor) FormatFunctionSignature(sig *FunctionSignature) string {
	var sb strings.Builder

	// Function name with link
	sb.WriteString(fmt.Sprintf("### `%s`\n\n", sig.Name))

	// Signature in code block
	sb.WriteString("```" + string(sig.Language) + "\n")
	sb.WriteString(sig.Signature)
	sb.WriteString("\n```\n\n")

	// Parameters
	if len(sig.Parameters) > 0 {
		sb.WriteString("**Parameters:**\n\n")
		sb.WriteString("| Name | Type | Description |\n")
		sb.WriteString("|------|------|-------------|\n")

		for _, p := range sig.Parameters {
			desc := p.Description
			if desc == "" {
				desc = "-"
			}
			sb.WriteString(fmt.Sprintf("| `%s` | `%s` | %s |\n",
				p.Name, p.Type, desc))
		}
		sb.WriteString("\n")
	}

	// Returns
	if len(sig.Returns) > 0 {
		sb.WriteString("**Returns:**\n\n")
		for _, r := range sig.Returns {
			sb.WriteString(fmt.Sprintf("- `%s`", r.Type))
			if r.Description != "" {
				sb.WriteString(fmt.Sprintf(": %s", r.Description))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Location
	sb.WriteString(fmt.Sprintf("*Defined in %s:%d*\n",
		filepath.Base(sig.FilePath), sig.Line))

	return sb.String()
}

// FormatTypeSignature formats a type signature for documentation
func (s *SignatureExtractor) FormatTypeSignature(ts *TypeSignature) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("### %s `%s`\n\n", ts.Kind, ts.Name))

	if len(ts.Fields) > 0 {
		sb.WriteString("**Fields:**\n\n")
		sb.WriteString("| Name | Type | Description |\n")
		sb.WriteString("|------|------|-------------|\n")

		for _, f := range ts.Fields {
			desc := f.Description
			if desc == "" {
				desc = "-"
			}
			sb.WriteString(fmt.Sprintf("| `%s` | `%s` | %s |\n",
				f.Name, f.Type, desc))
		}
		sb.WriteString("\n")
	}

	if len(ts.Methods) > 0 {
		sb.WriteString("**Methods:**\n\n")
		for _, m := range ts.Methods {
			sb.WriteString(fmt.Sprintf("- `%s`\n", m.Signature))
		}
		sb.WriteString("\n")
	}

	if len(ts.Implements) > 0 {
		sb.WriteString(fmt.Sprintf("*Implements: %s*\n",
			strings.Join(ts.Implements, ", ")))
	}

	return sb.String()
}
