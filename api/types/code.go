cat > api/types/code.go << 'EOF'
package types

// Language represents the programming language
type Language string

const (
	Go     Language = "go"
	NodeJS Language = "nodejs"
	Python Language = "python"
)

// CodeFile represents a source code file
type CodeFile struct {
	Path     string   `json:"path"`
	Language Language `json:"language"`
	Content  string   `json:"content"`
	Hash     string   `json:"hash,omitempty"` // For change detection
}

// Function represents a function/method in the code
type Function struct {
	Name       string   `json:"name"`
	Language   Language `json:"language"`
	FilePath   string   `json:"file_path"`
	Signature  string   `json:"signature"`
	Content    string   `json:"content"`
	LineStart  int      `json:"line_start"`
	LineEnd    int      `json:"line_end"`
	Complexity int      `json:"complexity"` // 1-10
	HasTest    bool     `json:"has_test"`
	TestName   string   `json:"test_name,omitempty"`
	TestFile   string   `json:"test_file,omitempty"`
	IsExported bool     `json:"is_exported"` // Public API
}

// Class represents a class in the code (for OOP languages)
type Class struct {
	Name       string     `json:"name"`
	Language   Language   `json:"language"`
	FilePath   string     `json:"file_path"`
	Methods    []Function `json:"methods"`
	LineStart  int        `json:"line_start"`
	LineEnd    int        `json:"line_end"`
	HasTest    bool       `json:"has_test"`
	IsExported bool       `json:"is_exported"`
}

// Import represents an imported dependency
type Import struct {
	Path    string `json:"path"`
	Alias   string `json:"alias,omitempty"`
	IsLocal bool   `json:"is_local"` // Internal vs external
}

// CodeAnalysis represents the complete analysis of code
type CodeAnalysis struct {
	Language    Language   `json:"language"`
	Files       []CodeFile `json:"files"`
	Functions   []Function `json:"functions"`
	Classes     []Class    `json:"classes,omitempty"`
	Imports     []Import   `json:"imports"`
	TestFiles   []string   `json:"test_files"`
	Coverage    float64    `json:"coverage,omitempty"`
	Summary     string     `json:"summary"`
	Timestamp   string     `json:"timestamp"`
}

// ChangeType represents what changed in the code
type ChangeType string

const (
	Added    ChangeType = "added"
	Modified ChangeType = "modified"
	Deleted  ChangeType = "deleted"
	Renamed  ChangeType = "renamed"
)

// CodeChange represents a change to the codebase
type CodeChange struct {
	File        string     `json:"file"`
	Language    Language   `json:"language"`
	Type        ChangeType `json:"type"`
	Functions   []string   `json:"functions_affected"`
	Classes     []string   `json:"classes_affected,omitempty"`
	Content     string     `json:"content,omitempty"`
	PreviousHash string    `json:"previous_hash,omitempty"`
	NewHash     string     `json:"new_hash,omitempty"`
	Timestamp   string     `json:"timestamp"`
}
EOF