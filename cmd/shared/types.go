package shared

import (
	"fmt"
	"time"
)

// ExitCode represents program exit codes
type ExitCode int

const (
	ExitCodeSuccess ExitCode = iota
	ExitCodeConfigError
	ExitCodeModelError
	ExitCodeAnalysisError
	ExitCodeGenerationError
	ExitCodeWriteError
	ExitCodeWatchError
)

// String returns the string representation of the exit code
func (e ExitCode) String() string {
	switch e {
	case ExitCodeSuccess:
		return "success"
	case ExitCodeConfigError:
		return "configuration error"
	case ExitCodeModelError:
		return "model error"
	case ExitCodeAnalysisError:
		return "analysis error"
	case ExitCodeGenerationError:
		return "generation error"
	case ExitCodeWriteError:
		return "write error"
	case ExitCodeWatchError:
		return "watch error"
	default:
		return "unknown error"
	}
}

// LogLevel represents logging verbosity
type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
	LogLevelTrace
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelError:
		return "ERROR"
	case LogLevelWarn:
		return "WARN"
	case LogLevelInfo:
		return "INFO"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	default:
		return "UNKNOWN"
	}
}

// CommandFlags represents common CLI flags
type CommandFlags struct {
	ConfigPath  string
	ProjectPath string
	Watch       bool
	Force       bool
	Verbose     bool
	Version     bool
	LogLevel    LogLevel
	Timeout     time.Duration
}

// VersionInfo represents build version information
type VersionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// String returns formatted version string
func (v VersionInfo) String() string {
	return fmt.Sprintf("%s (commit: %s, built: %s, %s/%s, go: %s)",
		v.Version, v.Commit, v.BuildDate, v.OS, v.Arch, v.GoVersion)
}

// Progress represents progress information for long operations
type Progress struct {
	Current int    `json:"current"`
	Total   int    `json:"total"`
	Message string `json:"message"`
	Done    bool   `json:"done"`
}

// WatchEvent represents a file system watch event
type WatchEvent struct {
	Path      string    `json:"path"`
	Operation string    `json:"operation"` // "create", "write", "remove", "rename"
	Timestamp time.Time `json:"timestamp"`
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	Success   bool          `json:"success"`
	Message   string        `json:"message"`
	Duration  time.Duration `json:"duration_ms"`
	ExitCode  ExitCode      `json:"exit_code"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// ModelStatus represents the status of a model
type ModelStatus struct {
	Name      string    `json:"name"`
	Available bool      `json:"available"`
	Size      string    `json:"size,omitempty"`
	Modified  time.Time `json:"modified,omitempty"`
}

// Stats represents command statistics
type Stats struct {
	FilesAnalyzed    int           `json:"files_analyzed"`
	FunctionsFound   int           `json:"functions_found"`
	TestsGenerated   int           `json:"tests_generated"`
	DocsGenerated    int           `json:"docs_generated"`
	ModelsQueried    int           `json:"models_queried"`
	TotalDuration    time.Duration `json:"total_duration_ms"`
	ModelDuration    time.Duration `json:"model_duration_ms"`
	AnalysisDuration time.Duration `json:"analysis_duration_ms"`
	WriteDuration    time.Duration `json:"write_duration_ms"`
}