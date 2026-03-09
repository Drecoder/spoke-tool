package analyzers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"

	"github.com/yourusername/spoke-tool/api/types"
)

// DepsAnalyzer analyzes dependencies between code components
// This helps understand test impact and relationships
type DepsAnalyzer struct {
	// No state - pure functions
}

// DependencyGraph represents dependencies between code elements
type DependencyGraph struct {
	// Nodes in the graph (functions, files, packages)
	Nodes []*DepNode `json:"nodes"`

	// Edges between nodes
	Edges []*DepEdge `json:"edges"`

	// Language of the code
	Language types.Language `json:"language"`
}

// DepNode represents a node in the dependency graph
type DepNode struct {
	ID        string `json:"id"`
	Type      string `json:"type"` // "file", "function", "package", "class"
	Name      string `json:"name"`
	FilePath  string `json:"file_path,omitempty"`
	Package   string `json:"package,omitempty"`
	LineStart int    `json:"line_start"`
	LineEnd   int    `json:"line_end"`
}

// DepEdge represents an edge in the dependency graph
type DepEdge struct {
	From string `json:"from"` // Source node ID
	To   string `json:"to"`   // Target node ID
	Type string `json:"type"` // "imports", "calls", "uses", "implements"
	Line int    `json:"line,omitempty"`
}

// ImpactAnalysis represents the impact of a change
type ImpactAnalysis struct {
	// Nodes directly affected
	DirectlyAffected []string `json:"directly_affected"`

	// Nodes indirectly affected (transitive dependencies)
	IndirectlyAffected []string `json:"indirectly_affected"`

	// Test files that need to be run
	AffectedTests []string `json:"affected_tests"`

	// Impact score (0-100)
	ImpactScore int `json:"impact_score"`
}

// NewDepsAnalyzer creates a new dependency analyzer
func NewDepsAnalyzer() *DepsAnalyzer {
	return &DepsAnalyzer{}
}

// AnalyzeGoDeps analyzes dependencies in Go code
func (d *DepsAnalyzer) AnalyzeGoDeps(content string, filePath string) (*DependencyGraph, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, content, 0)
	if err != nil {
		return nil, err
	}

	graph := &DependencyGraph{
		Language: types.Go,
		Nodes:    []*DepNode{},
		Edges:    []*DepEdge{},
	}

	// Add file node
	fileNode := &DepNode{
		ID:       "file:" + filePath,
		Type:     "file",
		Name:     filePath,
		FilePath: filePath,
		Package:  node.Name.Name,
	}
	graph.Nodes = append(graph.Nodes, fileNode)

	// Add package node
	pkgNode := &DepNode{
		ID:      "package:" + node.Name.Name,
		Type:    "package",
		Name:    node.Name.Name,
		Package: node.Name.Name,
	}
	graph.Nodes = append(graph.Nodes, pkgNode)

	// Add import edges
	for _, imp := range node.Imports {
		if imp.Path != nil {
			importPath := strings.Trim(imp.Path.Value, "\"")
			impNode := &DepNode{
				ID:   "import:" + importPath,
				Type: "import",
				Name: importPath,
			}
			graph.Nodes = append(graph.Nodes, impNode)

			// File imports this package
			graph.Edges = append(graph.Edges, &DepEdge{
				From: fileNode.ID,
				To:   impNode.ID,
				Type: "imports",
				Line: fset.Position(imp.Pos()).Line,
			})
		}
	}

	// Analyze functions and their calls
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// Add function node
			funcNode := &DepNode{
				ID:        "func:" + filePath + ":" + x.Name.Name,
				Type:      "function",
				Name:      x.Name.Name,
				FilePath:  filePath,
				Package:   node.Name.Name,
				LineStart: fset.Position(x.Pos()).Line,
				LineEnd:   fset.Position(x.End()).Line,
			}
			graph.Nodes = append(graph.Nodes, funcNode)

			// Function belongs to file
			graph.Edges = append(graph.Edges, &DepEdge{
				From: funcNode.ID,
				To:   fileNode.ID,
				Type: "defined_in",
			})

			// Function belongs to package
			graph.Edges = append(graph.Edges, &DepEdge{
				From: funcNode.ID,
				To:   pkgNode.ID,
				Type: "part_of",
			})

			// Analyze function body for calls
			if x.Body != nil {
				ast.Inspect(x.Body, func(y ast.Node) bool {
					if call, ok := y.(*ast.CallExpr); ok {
						d.analyzeGoCall(call, funcNode.ID, filePath, fset, graph)
					}
					return true
				})
			}

		case *ast.CallExpr:
			// Already handled in FuncDecl
		}
		return true
	})

	return graph, nil
}

// AnalyzeNodeJSDeps analyzes dependencies in Node.js code
func (d *DepsAnalyzer) AnalyzeNodeJSDeps(content string, filePath string) (*DependencyGraph, error) {
	graph := &DependencyGraph{
		Language: types.NodeJS,
		Nodes:    []*DepNode{},
		Edges:    []*DepEdge{},
	}

	// Add file node
	fileNode := &DepNode{
		ID:       "file:" + filePath,
		Type:     "file",
		Name:     filePath,
		FilePath: filePath,
	}
	graph.Nodes = append(graph.Nodes, fileNode)

	// Find imports/requires
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lineNum := i + 1

		// Match require statements
		requireRegex := regexp.MustCompile(`require\s*\(\s*['"]([^'"]+)['"]\s*\)`)
		if matches := requireRegex.FindStringSubmatch(line); len(matches) > 1 {
			importPath := matches[1]

			impNode := &DepNode{
				ID:   "import:" + importPath,
				Type: "import",
				Name: importPath,
			}
			graph.Nodes = append(graph.Nodes, impNode)

			graph.Edges = append(graph.Edges, &DepEdge{
				From: fileNode.ID,
				To:   impNode.ID,
				Type: "imports",
				Line: lineNum,
			})
		}

		// Match import statements
		importRegex := regexp.MustCompile(`import\s+(?:.*\s+from\s+)?['"]([^'"]+)['"]`)
		if matches := importRegex.FindStringSubmatch(line); len(matches) > 1 {
			importPath := matches[1]

			impNode := &DepNode{
				ID:   "import:" + importPath,
				Type: "import",
				Name: importPath,
			}
			graph.Nodes = append(graph.Nodes, impNode)

			graph.Edges = append(graph.Edges, &DepEdge{
				From: fileNode.ID,
				To:   impNode.ID,
				Type: "imports",
				Line: lineNum,
			})
		}

		// Find function calls (simplified)
		callRegex := regexp.MustCompile(`(\w+)\s*\(`)
		if matches := callRegex.FindAllStringSubmatch(line, -1); len(matches) > 0 {
			for _, match := range matches {
				if len(match) > 1 {
					funcName := match[1]
					// Skip common keywords
					if d.isJSKeyword(funcName) {
						continue
					}

					// Add call edge
					graph.Edges = append(graph.Edges, &DepEdge{
						From: fileNode.ID,
						To:   "func:" + funcName,
						Type: "calls",
						Line: lineNum,
					})
				}
			}
		}
	}

	return graph, nil
}

// AnalyzePythonDeps analyzes dependencies in Python code
func (d *DepsAnalyzer) AnalyzePythonDeps(content string, filePath string) (*DependencyGraph, error) {
	graph := &DependencyGraph{
		Language: types.Python,
		Nodes:    []*DepNode{},
		Edges:    []*DepEdge{},
	}

	// Add file node
	fileNode := &DepNode{
		ID:       "file:" + filePath,
		Type:     "file",
		Name:     filePath,
		FilePath: filePath,
	}
	graph.Nodes = append(graph.Nodes, fileNode)

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Match import statements
		if strings.HasPrefix(trimmed, "import ") {
			parts := strings.Fields(trimmed)
			if len(parts) > 1 {
				importPath := parts[1]
				impNode := &DepNode{
					ID:   "import:" + importPath,
					Type: "import",
					Name: importPath,
				}
				graph.Nodes = append(graph.Nodes, impNode)

				graph.Edges = append(graph.Edges, &DepEdge{
					From: fileNode.ID,
					To:   impNode.ID,
					Type: "imports",
					Line: lineNum,
				})
			}
		}

		// Match from ... import statements
		if strings.HasPrefix(trimmed, "from ") {
			parts := strings.Fields(trimmed)
			if len(parts) > 1 {
				importPath := parts[1]
				impNode := &DepNode{
					ID:   "import:" + importPath,
					Type: "import",
					Name: importPath,
				}
				graph.Nodes = append(graph.Nodes, impNode)

				graph.Edges = append(graph.Edges, &DepEdge{
					From: fileNode.ID,
					To:   impNode.ID,
					Type: "imports",
					Line: lineNum,
				})
			}
		}

		// Find function calls (simplified)
		callRegex := regexp.MustCompile(`(\w+)\s*\(`)
		if matches := callRegex.FindAllStringSubmatch(line, -1); len(matches) > 0 {
			for _, match := range matches {
				if len(match) > 1 {
					funcName := match[1]
					// Skip common keywords
					if d.isPythonKeyword(funcName) {
						continue
					}

					// Add call edge
					graph.Edges = append(graph.Edges, &DepEdge{
						From: fileNode.ID,
						To:   "func:" + funcName,
						Type: "calls",
						Line: lineNum,
					})
				}
			}
		}
	}

	return graph, nil
}

// Helper methods for Go analysis

func (d *DepsAnalyzer) analyzeGoCall(call *ast.CallExpr, callerID string, filePath string, fset *token.FileSet, graph *DependencyGraph) {
	var funcName string

	switch fun := call.Fun.(type) {
	case *ast.Ident:
		funcName = fun.Name
	case *ast.SelectorExpr:
		// Method call or package function
		if ident, ok := fun.X.(*ast.Ident); ok {
			funcName = ident.Name + "." + fun.Sel.Name
		} else {
			funcName = fun.Sel.Name
		}
	default:
		return
	}

	calleeID := "func:" + funcName

	// Add edge
	graph.Edges = append(graph.Edges, &DepEdge{
		From: callerID,
		To:   calleeID,
		Type: "calls",
		Line: fset.Position(call.Pos()).Line,
	})
}

// FindAffectedFunctions finds functions that would be affected by changes
func (d *DepsAnalyzer) FindAffectedFunctions(graph *DependencyGraph, changedNodes []string) *ImpactAnalysis {
	impact := &ImpactAnalysis{
		DirectlyAffected:   []string{},
		IndirectlyAffected: []string{},
		AffectedTests:      []string{},
	}

	// Build adjacency maps
	incoming := make(map[string][]string)
	outgoing := make(map[string][]string)

	for _, edge := range graph.Edges {
		incoming[edge.To] = append(incoming[edge.To], edge.From)
		outgoing[edge.From] = append(outgoing[edge.From], edge.To)
	}

	// Find directly affected (reverse dependencies)
	for _, node := range changedNodes {
		if deps, ok := incoming[node]; ok {
			impact.DirectlyAffected = append(impact.DirectlyAffected, deps...)
		}
	}

	// Find indirectly affected (transitive)
	visited := make(map[string]bool)
	for _, node := range impact.DirectlyAffected {
		d.collectTransitive(node, incoming, visited)
	}

	for node := range visited {
		if !d.contains(impact.DirectlyAffected, node) {
			impact.IndirectlyAffected = append(impact.IndirectlyAffected, node)
		}
	}

	// Find affected test files
	for _, node := range impact.DirectlyAffected {
		if strings.Contains(node, "_test.") || strings.Contains(node, "test_") {
			impact.AffectedTests = append(impact.AffectedTests, node)
		}
	}

	// Calculate impact score
	impact.ImpactScore = d.calculateImpactScore(impact)

	return impact
}

// FindTestImpact finds which tests need to run based on changes
func (d *DepsAnalyzer) FindTestImpact(graph *DependencyGraph, changedFiles []string) []string {
	affectedTests := make(map[string]bool)

	// Build reverse dependency map
	reverseDeps := make(map[string][]string)
	for _, edge := range graph.Edges {
		if edge.Type == "calls" || edge.Type == "imports" {
			reverseDeps[edge.To] = append(reverseDeps[edge.To], edge.From)
		}
	}

	// For each changed file, find dependent tests
	for _, file := range changedFiles {
		fileID := "file:" + file
		d.collectTests(fileID, reverseDeps, affectedTests)
	}

	// Convert to slice
	result := make([]string, 0, len(affectedTests))
	for test := range affectedTests {
		result = append(result, test)
	}

	return result
}

// Helper methods for impact analysis

func (d *DepsAnalyzer) collectTransitive(node string, deps map[string][]string, visited map[string]bool) {
	if visited[node] {
		return
	}
	visited[node] = true

	if dependents, ok := deps[node]; ok {
		for _, dep := range dependents {
			d.collectTransitive(dep, deps, visited)
		}
	}
}

func (d *DepsAnalyzer) collectTests(node string, deps map[string][]string, tests map[string]bool) {
	if strings.Contains(node, "test") || strings.Contains(node, "Test") {
		tests[node] = true
	}

	if dependents, ok := deps[node]; ok {
		for _, dep := range dependents {
			d.collectTests(dep, deps, tests)
		}
	}
}

func (d *DepsAnalyzer) calculateImpactScore(impact *ImpactAnalysis) int {
	score := 0

	// Each directly affected function adds 5 points
	score += len(impact.DirectlyAffected) * 5

	// Each indirectly affected adds 2 points
	score += len(impact.IndirectlyAffected) * 2

	// Each affected test adds 3 points
	score += len(impact.AffectedTests) * 3

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

func (d *DepsAnalyzer) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// FormatImpact formats impact analysis for display
func (d *DepsAnalyzer) FormatImpact(impact *ImpactAnalysis) string {
	var sb strings.Builder

	sb.WriteString("📊 Change Impact Analysis\n")
	sb.WriteString(strings.Repeat("=", 50))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("Impact Score: %d/100\n", impact.ImpactScore))

	if impact.ImpactScore >= 70 {
		sb.WriteString("⚠️  High impact - careful review needed\n")
	} else if impact.ImpactScore >= 40 {
		sb.WriteString("📝 Medium impact - moderate testing needed\n")
	} else {
		sb.WriteString("✅ Low impact - isolated change\n")
	}

	sb.WriteString("\n")

	if len(impact.DirectlyAffected) > 0 {
		sb.WriteString("Directly Affected:\n")
		sb.WriteString("------------------\n")
		for _, node := range impact.DirectlyAffected {
			sb.WriteString(fmt.Sprintf("  • %s\n", d.formatNodeID(node)))
		}
		sb.WriteString("\n")
	}

	if len(impact.IndirectlyAffected) > 0 {
		sb.WriteString("Indirectly Affected:\n")
		sb.WriteString("--------------------\n")
		// Show first 10
		for i, node := range impact.IndirectlyAffected {
			if i >= 10 {
				sb.WriteString(fmt.Sprintf("  • ... and %d more\n", len(impact.IndirectlyAffected)-10))
				break
			}
			sb.WriteString(fmt.Sprintf("  • %s\n", d.formatNodeID(node)))
		}
		sb.WriteString("\n")
	}

	if len(impact.AffectedTests) > 0 {
		sb.WriteString("Tests to Run:\n")
		sb.WriteString("-------------\n")
		for _, test := range impact.AffectedTests {
			sb.WriteString(fmt.Sprintf("  • %s\n", d.formatNodeID(test)))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (d *DepsAnalyzer) formatNodeID(nodeID string) string {
	// Remove prefix for cleaner display
	nodeID = strings.TrimPrefix(nodeID, "func:")
	nodeID = strings.TrimPrefix(nodeID, "file:")
	nodeID = strings.TrimPrefix(nodeID, "import:")

	// Just show filename if it's a path
	if strings.Contains(nodeID, "/") {
		parts := strings.Split(nodeID, "/")
		nodeID = parts[len(parts)-1]
	}

	return nodeID
}

// Keyword lists

func (d *DepsAnalyzer) isJSKeyword(word string) bool {
	keywords := map[string]bool{
		"if": true, "else": true, "for": true, "while": true,
		"switch": true, "case": true, "break": true, "continue": true,
		"return": true, "try": true, "catch": true, "finally": true,
		"throw": true, "new": true, "delete": true, "typeof": true,
		"instanceof": true, "void": true, "this": true, "super": true,
		"class": true, "extends": true, "import": true, "export": true,
		"console": true, "require": true, "module": true, "exports": true,
	}
	return keywords[word]
}

func (d *DepsAnalyzer) isPythonKeyword(word string) bool {
	keywords := map[string]bool{
		"if": true, "else": true, "elif": true, "for": true,
		"while": true, "break": true, "continue": true, "return": true,
		"try": true, "except": true, "finally": true, "raise": true,
		"import": true, "from": true, "as": true, "def": true,
		"class": true, "with": true, "lambda": true, "yield": true,
		"assert": true, "pass": true, "del": true, "global": true,
		"nonlocal": true, "True": true, "False": true, "None": true,
		"and": true, "or": true, "not": true, "in": true, "is": true,
		"print": true, "len": true, "range": true, "open": true,
	}
	return keywords[word]
}

// ExportGraph exports the dependency graph in various formats
func (d *DepsAnalyzer) ExportGraph(graph *DependencyGraph, format string) (string, error) {
	switch format {
	case "dot":
		return d.exportToDOT(graph), nil
	case "json":
		// Would return JSON
		return "", fmt.Errorf("JSON export not implemented")
	case "mermaid":
		return d.exportToMermaid(graph), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func (d *DepsAnalyzer) exportToDOT(graph *DependencyGraph) string {
	var sb strings.Builder

	sb.WriteString("digraph Dependencies {\n")
	sb.WriteString("  rankdir=LR;\n")
	sb.WriteString("  node [shape=box];\n\n")

	// Add nodes
	for _, node := range graph.Nodes {
		color := d.getNodeColor(node.Type)
		sb.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\", color=%s];\n",
			node.ID, node.Name, color))
	}

	sb.WriteString("\n")

	// Add edges
	for _, edge := range graph.Edges {
		style := d.getEdgeStyle(edge.Type)
		sb.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\", %s];\n",
			edge.From, edge.To, edge.Type, style))
	}

	sb.WriteString("}\n")

	return sb.String()
}

func (d *DepsAnalyzer) exportToMermaid(graph *DependencyGraph) string {
	var sb strings.Builder

	sb.WriteString("graph TD\n")

	// Add nodes
	nodeIDs := make(map[string]string)
	for i, node := range graph.Nodes {
		nodeID := fmt.Sprintf("N%d", i)
		nodeIDs[node.ID] = nodeID

		shape := d.getMermaidShape(node.Type)
		sb.WriteString(fmt.Sprintf("    %s%s[\"%s\"]\n", nodeID, shape, node.Name))
	}

	sb.WriteString("\n")

	// Add edges
	for _, edge := range graph.Edges {
		fromID := nodeIDs[edge.From]
		toID := nodeIDs[edge.To]

		if fromID != "" && toID != "" {
			sb.WriteString(fmt.Sprintf("    %s -->|%s| %s\n", fromID, edge.Type, toID))
		}
	}

	return sb.String()
}

func (d *DepsAnalyzer) getNodeColor(nodeType string) string {
	switch nodeType {
	case "file":
		return "lightblue"
	case "function":
		return "lightgreen"
	case "import":
		return "lightyellow"
	case "package":
		return "lightpink"
	default:
		return "lightgray"
	}
}

func (d *DepsAnalyzer) getEdgeStyle(edgeType string) string {
	switch edgeType {
	case "imports":
		return "style=dashed"
	case "calls":
		return "style=solid"
	default:
		return ""
	}
}

func (d *DepsAnalyzer) getMermaidShape(nodeType string) string {
	switch nodeType {
	case "file":
		return "(["
	case "function":
		return "(("
	case "import":
		return "[["
	default:
		return "["
	}
}

func (d *DepsAnalyzer) getMermaidStyle(edgeType string) string {
	switch edgeType {
	case "imports":
		return "-.->"
	default:
		return "-->"
	}
}
