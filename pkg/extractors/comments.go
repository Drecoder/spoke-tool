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
)

// CommentExtractor extracts comments and documentation from code
// This is PURELY EXTRACTIVE - no modifications
type CommentExtractor struct {
	// No state - pure functions
}

// Comment represents a extracted comment with context
type Comment struct {
	// The comment text
	Text string `json:"text"`

	// Type of comment (line, block, doc)
	Type string `json:"type"` // "line", "block", "doc"

	// Language of the code
	Language types.Language `json:"language"`

	// File where comment was found
	FilePath string `json:"file_path"`

	// Line number
	Line int `json:"line"`

	// The code element this comment is associated with
	AssociatedWith string `json:"associated_with,omitempty"` // function name, struct name, etc.

	// Whether this is a TODO/FIXME/NOTE
	IsNote bool `json:"is_note"`

	// Note type if IsNote is true
	NoteType string `json:"note_type,omitempty"` // "TODO", "FIXME", "NOTE", "XXX"
}

// DocComment represents a structured documentation comment
type DocComment struct {
	// The full comment
	Raw string `json:"raw"`

	// Summary (first line/sentence)
	Summary string `json:"summary"`

	// Description (rest of comment)
	Description string `json:"description"`

	// Parameters documented
	Parameters []DocParam `json:"parameters,omitempty"`

	// Return value documented
	Returns string `json:"returns,omitempty"`

	// Examples found
	Examples []string `json:"examples,omitempty"`

	// Deprecated notice
	Deprecated string `json:"deprecated,omitempty"`

	// See also references
	SeeAlso []string `json:"see_also,omitempty"`
}

// DocParam represents a documented parameter
type DocParam struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
	Desc string `json:"desc"`
}

// NewCommentExtractor creates a new comment extractor
func NewCommentExtractor() *CommentExtractor {
	return &CommentExtractor{}
}

// ExtractFromGo extracts comments from Go code
func (c *CommentExtractor) ExtractFromGo(content string, filePath string) ([]*Comment, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var comments []*Comment

	// Extract file-level comments
	if node.Doc != nil {
		for _, comment := range node.Doc.List {
			comments = append(comments, &Comment{
				Text:           comment.Text,
				Type:           c.getCommentType(comment.Text),
				Language:       types.Go,
				FilePath:       filePath,
				Line:           fset.Position(comment.Pos()).Line,
				AssociatedWith: "file:" + filepath.Base(filePath),
				IsNote:         c.isNote(comment.Text),
				NoteType:       c.getNoteType(comment.Text),
			})
		}
	}

	// Extract comments from AST
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// Function comments
			if x.Doc != nil {
				for _, comment := range x.Doc.List {
					comments = append(comments, &Comment{
						Text:           comment.Text,
						Type:           c.getCommentType(comment.Text),
						Language:       types.Go,
						FilePath:       filePath,
						Line:           fset.Position(comment.Pos()).Line,
						AssociatedWith: "func:" + x.Name.Name,
						IsNote:         c.isNote(comment.Text),
						NoteType:       c.getNoteType(comment.Text),
					})
				}
			}

		case *ast.GenDecl:
			// Type, var, const comments
			if x.Doc != nil {
				for _, comment := range x.Doc.List {
					// Find the name of what's being declared
					name := ""
					if len(x.Specs) > 0 {
						if spec, ok := x.Specs[0].(*ast.TypeSpec); ok {
							name = "type:" + spec.Name.Name
						} else if spec, ok := x.Specs[0].(*ast.ValueSpec); ok && len(spec.Names) > 0 {
							name = "var:" + spec.Names[0].Name
						}
					}

					comments = append(comments, &Comment{
						Text:           comment.Text,
						Type:           c.getCommentType(comment.Text),
						Language:       types.Go,
						FilePath:       filePath,
						Line:           fset.Position(comment.Pos()).Line,
						AssociatedWith: name,
						IsNote:         c.isNote(comment.Text),
						NoteType:       c.getNoteType(comment.Text),
					})
				}
			}

			// Line comments inside declarations
			for _, spec := range x.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok && ts.Comment != nil {
					for _, comment := range ts.Comment.List {
						comments = append(comments, &Comment{
							Text:           comment.Text,
							Type:           c.getCommentType(comment.Text),
							Language:       types.Go,
							FilePath:       filePath,
							Line:           fset.Position(comment.Pos()).Line,
							AssociatedWith: "type:" + ts.Name.Name,
							IsNote:         c.isNote(comment.Text),
							NoteType:       c.getNoteType(comment.Text),
						})
					}
				}
			}
		}

		return true
	})

	return comments, nil
}

// ExtractFromNodeJS extracts comments from Node.js/JavaScript code
func (c *CommentExtractor) ExtractFromNodeJS(content string, filePath string) ([]*Comment, error) {
	var comments []*Comment

	lines := strings.Split(content, "\n")

	// Track multiline comments
	inMultilineComment := false
	var multilineComment strings.Builder
	multilineStartLine := 0

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Handle multiline comments
		if strings.Contains(trimmed, "/*") && !inMultilineComment {
			inMultilineComment = true
			multilineStartLine = lineNum
			multilineComment.Reset()

			// Check if comment ends on same line
			if strings.Contains(trimmed, "*/") {
				// Extract comment content
				startIdx := strings.Index(trimmed, "/*")
				endIdx := strings.Index(trimmed, "*/")
				if endIdx > startIdx {
					commentText := trimmed[startIdx : endIdx+2]
					comments = append(comments, &Comment{
						Text:     commentText,
						Type:     "block",
						Language: types.NodeJS,
						FilePath: filePath,
						Line:     lineNum,
						IsNote:   c.isNote(commentText),
						NoteType: c.getNoteType(commentText),
					})
				}
				inMultilineComment = false
			} else {
				// Start collecting multiline
				multilineComment.WriteString(trimmed + "\n")
			}
			continue
		}

		if inMultilineComment {
			multilineComment.WriteString(trimmed + "\n")
			if strings.Contains(trimmed, "*/") {
				// End of multiline comment
				comments = append(comments, &Comment{
					Text:     multilineComment.String(),
					Type:     "block",
					Language: types.NodeJS,
					FilePath: filePath,
					Line:     multilineStartLine,
					IsNote:   c.isNote(multilineComment.String()),
					NoteType: c.getNoteType(multilineComment.String()),
				})
				inMultilineComment = false
			}
			continue
		}

		// Handle line comments
		if strings.HasPrefix(trimmed, "//") {
			comments = append(comments, &Comment{
				Text:     trimmed,
				Type:     "line",
				Language: types.NodeJS,
				FilePath: filePath,
				Line:     lineNum,
				IsNote:   c.isNote(trimmed),
				NoteType: c.getNoteType(trimmed),
			})
		}

		// Handle JSDoc comments
		if strings.HasPrefix(trimmed, "/**") {
			// Find end of JSDoc
			jsdocContent := trimmed + "\n"
			j := i + 1
			for j < len(lines) {
				nextLine := strings.TrimSpace(lines[j])
				jsdocContent += nextLine + "\n"
				if strings.Contains(nextLine, "*/") {
					break
				}
				j++
			}

			comments = append(comments, &Comment{
				Text:     jsdocContent,
				Type:     "doc",
				Language: types.NodeJS,
				FilePath: filePath,
				Line:     lineNum,
				IsNote:   c.isNote(jsdocContent),
				NoteType: c.getNoteType(jsdocContent),
			})
		}
	}

	// Try to associate comments with functions (simplified)
	c.associateNodeJSComments(comments, content)

	return comments, nil
}

// ExtractFromPython extracts comments from Python code
func (c *CommentExtractor) ExtractFromPython(content string, filePath string) ([]*Comment, error) {
	var comments []*Comment

	lines := strings.Split(content, "\n")

	// Track multiline strings/docstrings
	inMultilineString := false
	var multilineString strings.Builder
	multilineStartLine := 0
	quoteChar := ""

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Handle docstrings and multiline strings
		if !inMultilineString && (strings.HasPrefix(trimmed, "\"\"\"") || strings.HasPrefix(trimmed, "'''")) {
			inMultilineString = true
			multilineStartLine = lineNum
			multilineString.Reset()
			quoteChar = trimmed[:3]

			// Check if docstring ends on same line
			if strings.Contains(trimmed[3:], quoteChar) {
				comments = append(comments, &Comment{
					Text:     trimmed,
					Type:     "doc",
					Language: types.Python,
					FilePath: filePath,
					Line:     lineNum,
					IsNote:   c.isNote(trimmed),
					NoteType: c.getNoteType(trimmed),
				})
				inMultilineString = false
			} else {
				multilineString.WriteString(trimmed + "\n")
			}
			continue
		}

		if inMultilineString {
			multilineString.WriteString(trimmed + "\n")
			if strings.Contains(trimmed, quoteChar) {
				// End of docstring
				comments = append(comments, &Comment{
					Text:     multilineString.String(),
					Type:     "doc",
					Language: types.Python,
					FilePath: filePath,
					Line:     multilineStartLine,
					IsNote:   c.isNote(multilineString.String()),
					NoteType: c.getNoteType(multilineString.String()),
				})
				inMultilineString = false
			}
			continue
		}

		// Handle line comments
		if strings.Contains(line, "#") {
			// Find the comment part
			parts := strings.SplitN(line, "#", 2)
			if len(parts) > 1 {
				commentText := "#" + parts[1]
				comments = append(comments, &Comment{
					Text:     commentText,
					Type:     "line",
					Language: types.Python,
					FilePath: filePath,
					Line:     lineNum,
					IsNote:   c.isNote(commentText),
					NoteType: c.getNoteType(commentText),
				})
			}
		}
	}

	// Try to associate comments with functions/classes
	c.associatePythonComments(comments, content)

	return comments, nil
}

// ExtractNotes extracts only TODO/FIXME/NOTE comments
func (c *CommentExtractor) ExtractNotes(comments []*Comment) []*Comment {
	var notes []*Comment

	for _, comment := range comments {
		if comment.IsNote {
			notes = append(notes, comment)
		}
	}

	return notes
}

// ParseDocComment parses a documentation comment into structured format
func (c *CommentExtractor) ParseDocComment(comment *Comment) *DocComment {
	if comment.Type != "doc" {
		return nil
	}

	doc := &DocComment{
		Raw:        comment.Text,
		Parameters: []DocParam{},
		SeeAlso:    []string{},
	}

	// Clean the comment
	text := c.cleanComment(comment.Text)
	lines := strings.Split(text, "\n")

	// Extract summary (first line or first sentence)
	if len(lines) > 0 {
		firstLine := lines[0]
		if strings.Contains(firstLine, ".") {
			parts := strings.SplitN(firstLine, ".", 2)
			doc.Summary = strings.TrimSpace(parts[0]) + "."
			if len(parts) > 1 {
				doc.Description = strings.TrimSpace(parts[1])
			}
		} else {
			doc.Summary = firstLine
			if len(lines) > 1 {
				doc.Description = strings.Join(lines[1:], "\n")
			}
		}
	}

	// Extract language-specific elements
	switch comment.Language {
	case types.Go:
		c.parseGoDoc(doc, text)
	case types.NodeJS:
		c.parseJSDoc(doc, text)
	case types.Python:
		c.parsePyDoc(doc, text)
	}

	return doc
}

// Helper methods for Go docs

func (c *CommentExtractor) parseGoDoc(doc *DocComment, text string) {
	// Look for deprecation
	if strings.Contains(strings.ToLower(text), "deprecated") {
		doc.Deprecated = "yes"
	}

	// Look for examples (simplified)
	if strings.Contains(text, "Example:") {
		parts := strings.Split(text, "Example:")
		if len(parts) > 1 {
			doc.Examples = append(doc.Examples, strings.TrimSpace(parts[1]))
		}
	}
}

// Helper methods for JSDoc

func (c *CommentExtractor) parseJSDoc(doc *DocComment, text string) {
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// @param {type} name - description
		if strings.HasPrefix(line, "@param") {
			param := DocParam{}

			// Extract type
			typeRegex := regexp.MustCompile(`{([^}]+)}`)
			if matches := typeRegex.FindStringSubmatch(line); len(matches) > 1 {
				param.Type = matches[1]
			}

			// Extract name
			nameRegex := regexp.MustCompile(`}\s*(\w+)`)
			if matches := nameRegex.FindStringSubmatch(line); len(matches) > 1 {
				param.Name = matches[1]
			}

			// Extract description
			descRegex := regexp.MustCompile(`-\s*(.+)`)
			if matches := descRegex.FindStringSubmatch(line); len(matches) > 1 {
				param.Desc = matches[1]
			}

			doc.Parameters = append(doc.Parameters, param)
		}

		// @returns {type} description
		if strings.HasPrefix(line, "@returns") || strings.HasPrefix(line, "@return") {
			typeRegex := regexp.MustCompile(`{([^}]+)}`)
			if matches := typeRegex.FindStringSubmatch(line); len(matches) > 1 {
				doc.Returns = matches[1]
			}
		}

		// @example
		if strings.HasPrefix(line, "@example") {
			example := strings.TrimPrefix(line, "@example")
			doc.Examples = append(doc.Examples, strings.TrimSpace(example))
		}

		// @deprecated
		if strings.HasPrefix(line, "@deprecated") {
			doc.Deprecated = strings.TrimPrefix(line, "@deprecated")
		}

		// @see
		if strings.HasPrefix(line, "@see") {
			see := strings.TrimPrefix(line, "@see")
			doc.SeeAlso = append(doc.SeeAlso, strings.TrimSpace(see))
		}
	}
}

// Helper methods for Python docstrings

func (c *CommentExtractor) parsePyDoc(doc *DocComment, text string) {
	lines := strings.Split(text, "\n")
	inArgs := false
	inReturns := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Google-style Args:
		if strings.Contains(line, "Args:") {
			inArgs = true
			inReturns = false
			continue
		}

		// Google-style Returns:
		if strings.Contains(line, "Returns:") {
			inArgs = false
			inReturns = true
			continue
		}

		if inArgs && line != "" {
			// Parse "name (type): description"
			paramRegex := regexp.MustCompile(`(\w+)\s+\(([^)]+)\):\s*(.+)`)
			if matches := paramRegex.FindStringSubmatch(line); len(matches) > 3 {
				doc.Parameters = append(doc.Parameters, DocParam{
					Name: matches[1],
					Type: matches[2],
					Desc: matches[3],
				})
			}
		}

		if inReturns && line != "" {
			doc.Returns = line
			inReturns = false
		}

		// Look for Examples:
		if strings.Contains(line, "Example:") || strings.Contains(line, "Examples:") {
			// Collect following indented lines
		}
	}
}

// Helper methods for comment association

func (c *CommentExtractor) associateNodeJSComments(comments []*Comment, content string) {
	lines := strings.Split(content, "\n")

	for _, comment := range comments {
		// Look for function definition after comment
		if comment.Line < len(lines) {
			for i := comment.Line; i < len(lines) && i < comment.Line+5; i++ {
				line := lines[i-1] // Lines are 1-indexed

				// Check for function declaration
				if strings.Contains(line, "function ") {
					funcRegex := regexp.MustCompile(`function\s+(\w+)`)
					if matches := funcRegex.FindStringSubmatch(line); len(matches) > 1 {
						comment.AssociatedWith = "func:" + matches[1]
						break
					}
				}

				// Check for arrow function assignment
				if strings.Contains(line, "=") && strings.Contains(line, "=>") {
					nameRegex := regexp.MustCompile(`(\w+)\s*=`)
					if matches := nameRegex.FindStringSubmatch(line); len(matches) > 1 {
						comment.AssociatedWith = "func:" + matches[1]
						break
					}
				}

				// Check for class
				if strings.Contains(line, "class ") {
					classRegex := regexp.MustCompile(`class\s+(\w+)`)
					if matches := classRegex.FindStringSubmatch(line); len(matches) > 1 {
						comment.AssociatedWith = "class:" + matches[1]
						break
					}
				}
			}
		}
	}
}

func (c *CommentExtractor) associatePythonComments(comments []*Comment, content string) {
	lines := strings.Split(content, "\n")

	for _, comment := range comments {
		// Look for function/class definition after comment
		if comment.Line < len(lines) {
			for i := comment.Line; i < len(lines) && i < comment.Line+5; i++ {
				line := lines[i-1]
				trimmed := strings.TrimSpace(line)

				// Check for function definition
				if strings.HasPrefix(trimmed, "def ") {
					funcRegex := regexp.MustCompile(`def\s+(\w+)`)
					if matches := funcRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
						comment.AssociatedWith = "func:" + matches[1]
						break
					}
				}

				// Check for class definition
				if strings.HasPrefix(trimmed, "class ") {
					classRegex := regexp.MustCompile(`class\s+(\w+)`)
					if matches := classRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
						comment.AssociatedWith = "class:" + matches[1]
						break
					}
				}
			}
		}
	}
}

// Utility methods

func (c *CommentExtractor) getCommentType(comment string) string {
	if strings.HasPrefix(comment, "//") {
		return "line"
	}
	if strings.HasPrefix(comment, "/*") {
		if strings.HasPrefix(comment, "/**") {
			return "doc"
		}
		return "block"
	}
	if strings.HasPrefix(comment, "\"\"\"") || strings.HasPrefix(comment, "'''") {
		return "doc"
	}
	return "unknown"
}

func (c *CommentExtractor) isNote(comment string) bool {
	notes := []string{"TODO", "FIXME", "NOTE", "XXX", "HACK", "BUG", "OPTIMIZE"}
	upper := strings.ToUpper(comment)

	for _, note := range notes {
		if strings.Contains(upper, note) {
			return true
		}
	}
	return false
}

func (c *CommentExtractor) getNoteType(comment string) string {
	notes := []string{"TODO", "FIXME", "NOTE", "XXX", "HACK", "BUG", "OPTIMIZE"}
	upper := strings.ToUpper(comment)

	for _, note := range notes {
		if strings.Contains(upper, note) {
			return note
		}
	}
	return ""
}

func (c *CommentExtractor) cleanComment(comment string) string {
	// Remove comment markers
	comment = strings.TrimPrefix(comment, "//")
	comment = strings.TrimPrefix(comment, "/*")
	comment = strings.TrimSuffix(comment, "*/")
	comment = strings.TrimPrefix(comment, "/**")
	comment = strings.TrimPrefix(comment, "\"\"\"")
	comment = strings.TrimSuffix(comment, "\"\"\"")
	comment = strings.TrimPrefix(comment, "'''")
	comment = strings.TrimSuffix(comment, "'''")

	// Trim spaces
	comment = strings.TrimSpace(comment)

	// Remove leading * from each line in block comments
	lines := strings.Split(comment, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimPrefix(strings.TrimSpace(line), "*")
	}

	return strings.Join(lines, "\n")
}

// GetSummary returns a human-readable summary of comments
func (c *CommentExtractor) GetSummary(comments []*Comment) string {
	if len(comments) == 0 {
		return "No comments found"
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("📝 Found %d comments\n", len(comments)))
	sb.WriteString(strings.Repeat("=", 50))
	sb.WriteString("\n\n")

	// Group by file
	byFile := make(map[string][]*Comment)
	for _, comment := range comments {
		byFile[comment.FilePath] = append(byFile[comment.FilePath], comment)
	}

	for file, fileComments := range byFile {
		sb.WriteString(fmt.Sprintf("📁 %s:\n", filepath.Base(file)))

		// Count by type
		lineCount, blockCount, docCount := 0, 0, 0
		noteCount := 0

		for _, c := range fileComments {
			switch c.Type {
			case "line":
				lineCount++
			case "block":
				blockCount++
			case "doc":
				docCount++
			}
			if c.IsNote {
				noteCount++
			}
		}

		sb.WriteString(fmt.Sprintf("   Line comments: %d\n", lineCount))
		sb.WriteString(fmt.Sprintf("   Block comments: %d\n", blockCount))
		sb.WriteString(fmt.Sprintf("   Doc comments: %d\n", docCount))
		sb.WriteString(fmt.Sprintf("   Notes (TODO/FIXME): %d\n", noteCount))

		// Show notes
		if noteCount > 0 {
			sb.WriteString("   Notes:\n")
			for _, comment := range fileComments {
				if comment.IsNote {
					sb.WriteString(fmt.Sprintf("     • [%s] Line %d: %s\n",
						comment.NoteType, comment.Line, c.cleanComment(comment.Text)))
				}
			}
		}

		sb.WriteString("\n")
	}

	return sb.String()
}
