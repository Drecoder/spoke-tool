package types

// DocFormat represents the documentation format
type DocFormat string

const (
	DocFormatGodoc    DocFormat = "godoc"
	DocFormatJSDoc    DocFormat = "jsdoc"
	DocFormatPyDoc    DocFormat = "pydoc"
	DocFormatMarkdown DocFormat = "markdown"
	DocFormatRST      DocFormat = "rst" // reStructuredText for Python
)

// DocSection represents a section of documentation
type DocSection string

const (
	DocSectionTitle        DocSection = "title"
	DocSectionInstallation DocSection = "installation"
	DocSectionQuickStart   DocSection = "quickstart"
	DocSectionAPI          DocSection = "api"
	DocSectionExamples     DocSection = "examples"
	DocSectionContributing DocSection = "contributing"
	DocSectionLicense      DocSection = "license"
)

// DocContent represents generated documentation content
type DocContent struct {
	Language    Language      `json:"language"`
	Function    string        `json:"function,omitempty"`
	Class       string        `json:"class,omitempty"`
	Format      DocFormat     `json:"format"`
	Content     string        `json:"content"`
	Examples    []DocExample  `json:"examples,omitempty"`
	Parameters  []DocParam    `json:"parameters,omitempty"`
	Returns     string        `json:"returns,omitempty"`
	Exceptions  []string      `json:"exceptions,omitempty"`
}

// DocParam represents a function parameter in documentation
type DocParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Optional    bool   `json:"optional"`
	Default     string `json:"default,omitempty"`
}

// DocExample represents a code example in documentation
type DocExample struct {
	Language Language `json:"language"`
	Code     string   `json:"code"`
	Description string `json:"description,omitempty"`
	IsTest    bool    `json:"is_test"` // Whether this comes from a test
}

// ReadmeSection represents a section in the README
type ReadmeSection struct {
	Type    DocSection `json:"type"`
	Title   string     `json:"title"`
	Content string     `json:"content"`
	Order   int        `json:"order"`
}

// Readme represents the complete README file
type Readme struct {
	Path        string          `json:"path"`
	Language    Language        `json:"language"`
	ProjectName string          `json:"project_name"`
	Description string          `json:"description"`
	Sections    []ReadmeSection `json:"sections"`
	Content     string          `json:"content"` // Full markdown
	LastUpdated string          `json:"last_updated"`
}

// DocSuggestion represents a documentation suggestion
type DocSuggestion struct {
	Language     Language   `json:"language"`
	FunctionName string     `json:"function_name,omitempty"`
	ClassName    string     `json:"class_name,omitempty"`
	Content      DocContent `json:"content"`
	Section      DocSection `json:"section,omitempty"`
	Confidence   float64    `json:"confidence"`
}