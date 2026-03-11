package common

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// FileUtils provides file system utilities
type FileUtils struct{}

// FileInfo represents file information
type FileInfo struct {
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"mod_time"`
	IsDir     bool      `json:"is_dir"`
	Extension string    `json:"extension"`
	Hash      string    `json:"hash,omitempty"`
}

// FileExists checks if a file exists
func (fu *FileUtils) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDir checks if a path is a directory
func (fu *FileUtils) IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ReadFile reads a file and returns its contents
func (fu *FileUtils) ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return string(data), nil
}

// WriteFile writes data to a file
func (fu *FileUtils) WriteFile(path string, data string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// GetFileInfo returns information about a file
func (fu *FileUtils) GetFileInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Path:      path,
		Name:      info.Name(),
		Size:      info.Size(),
		ModTime:   info.ModTime(),
		IsDir:     info.IsDir(),
		Extension: filepath.Ext(path),
	}, nil
}

// ListFiles returns all files in a directory matching the pattern
func (fu *FileUtils) ListFiles(dir string, pattern string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if pattern == "" || matchesPattern(path, pattern) {
				files = append(files, path)
			}
		}
		return nil
	})

	return files, err
}

// CopyFile copies a file from src to dst
func (fu *FileUtils) CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// HashUtils provides hashing utilities
type HashUtils struct{}

// MD5 calculates MD5 hash of a string
func (hu *HashUtils) MD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// MD5File calculates MD5 hash of a file
func (hu *HashUtils) MD5File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// SHA256 calculates SHA256 hash of a string
func (hu *HashUtils) SHA256(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// SHA256File calculates SHA256 hash of a file
func (hu *HashUtils) SHA256File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// StringUtils provides string manipulation utilities
type StringUtils struct{}

// SnakeCase converts a string to snake_case
func (su *StringUtils) SnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// CamelCase converts a string to camelCase
func (su *StringUtils) CamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// Replace separators with spaces
	re := regexp.MustCompile(`[-_\s]+`)
	s = re.ReplaceAllString(s, " ")

	// Split and capitalize
	parts := strings.Split(s, " ")
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}

	// Join and lower first
	result := strings.Join(parts, "")
	if len(result) > 0 {
		result = strings.ToLower(result[:1]) + result[1:]
	}

	return result
}

// PascalCase converts a string to PascalCase
func (su *StringUtils) PascalCase(s string) string {
	s = su.CamelCase(s)
	if len(s) > 0 {
		s = strings.ToUpper(s[:1]) + s[1:]
	}
	return s
}

// Truncate truncates a string to max length
func (su *StringUtils) Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// Indent adds indentation to each line
func (su *StringUtils) Indent(s string, indent string) string {
	if s == "" {
		return s
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}

// ContainsAny checks if string contains any of the substrings
func (su *StringUtils) ContainsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicate lines from a string
func (su *StringUtils) RemoveDuplicates(s string) string {
	lines := strings.Split(s, "\n")
	seen := make(map[string]bool)
	var result []string

	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// SliceUtils provides slice manipulation utilities
type SliceUtils struct{}

// Contains checks if a slice contains a string
func (su *SliceUtils) Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsInt checks if a slice contains an int
func (su *SliceUtils) ContainsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Unique returns a slice with duplicate strings removed
func (su *SliceUtils) Unique(slice []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// Intersection returns strings common to both slices
func (su *SliceUtils) Intersection(a, b []string) []string {
	seen := make(map[string]bool)
	for _, item := range a {
		seen[item] = true
	}

	var result []string
	for _, item := range b {
		if seen[item] {
			result = append(result, item)
		}
	}

	return result
}

// Difference returns strings in a but not in b
func (su *SliceUtils) Difference(a, b []string) []string {
	seen := make(map[string]bool)
	for _, item := range b {
		seen[item] = true
	}

	var result []string
	for _, item := range a {
		if !seen[item] {
			result = append(result, item)
		}
	}

	return result
}

// Chunk splits a slice into chunks of size n
func (su *SliceUtils) Chunk(slice []string, n int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += n {
		end := i + n
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// TimeUtils provides time utilities
type TimeUtils struct{}

// FormatDuration formats a duration in a human-readable way
func (tu *TimeUtils) FormatDuration(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return fmt.Sprintf("%d ns", d.Nanoseconds())
	case d < time.Millisecond:
		return fmt.Sprintf("%.2f µs", float64(d.Nanoseconds())/1000)
	case d < time.Second:
		return fmt.Sprintf("%.2f ms", float64(d.Nanoseconds())/1_000_000)
	case d < time.Minute:
		return fmt.Sprintf("%.2f s", d.Seconds())
	default:
		return d.Round(time.Second).String()
	}
}

// ParseTime safely parses a time string
func (tu *TimeUtils) ParseTime(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", s)
}

// JSONUtils provides JSON utilities
type JSONUtils struct{}

// PrettyJSON returns pretty-printed JSON
func (ju *JSONUtils) PrettyJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// CompactJSON returns compact JSON
func (ju *JSONUtils) CompactJSON(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MergeJSON merges multiple JSON objects
func (ju *JSONUtils) MergeJSON(objs ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, obj := range objs {
		for k, v := range obj {
			result[k] = v
		}
	}
	return result
}

// OSUtils provides operating system utilities
type OSUtils struct{}

// GetEnv gets an environment variable with a default
func (ou *OSUtils) GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// HomeDir returns the user's home directory
func (ou *OSUtils) HomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return home, nil
}

// TempDir creates a temporary directory
func (ou *OSUtils) TempDir(prefix string) (string, error) {
	return os.MkdirTemp("", prefix)
}

// IsWindows returns true if running on Windows
func (ou *OSUtils) IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux returns true if running on Linux
func (ou *OSUtils) IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsMacOS returns true if running on macOS
func (ou *OSUtils) IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// SystemInfo returns information about the system
func (ou *OSUtils) SystemInfo() map[string]string {
	return map[string]string{
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"cpus":       fmt.Sprintf("%d", runtime.NumCPU()),
		"go_version": runtime.Version(),
		"hostname":   getHostname(),
	}
}

// MathUtils provides mathematical utilities
type MathUtils struct{}

// Min returns the minimum of two ints
func (mu *MathUtils) Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two ints
func (mu *MathUtils) Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Clamp clamps a value between min and max
func (mu *MathUtils) Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Percentage calculates a percentage
func (mu *MathUtils) Percentage(value, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (value / total) * 100
}

// VersionUtils provides version comparison utilities
type VersionUtils struct{}

// CompareVersions compares two semantic versions
// Returns -1 if v1 < v2, 0 if equal, 1 if v1 > v2
func (vu *VersionUtils) CompareVersions(v1, v2 string) int {
	v1Parts := strings.Split(strings.TrimPrefix(v1, "v"), ".")
	v2Parts := strings.Split(strings.TrimPrefix(v2, "v"), ".")

	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 string
		if i < len(v1Parts) {
			p1 = v1Parts[i]
		}
		if i < len(v2Parts) {
			p2 = v2Parts[i]
		}

		// Compare numeric parts
		n1, e1 := parseInt(p1)
		n2, e2 := parseInt(p2)

		if e1 == nil && e2 == nil {
			if n1 < n2 {
				return -1
			}
			if n1 > n2 {
				return 1
			}
		} else {
			// Compare as strings
			if p1 < p2 {
				return -1
			}
			if p1 > p2 {
				return 1
			}
		}
	}

	return 0
}

// IsVersionGreater checks if v1 > v2
func (vu *VersionUtils) IsVersionGreater(v1, v2 string) bool {
	return vu.CompareVersions(v1, v2) > 0
}

// IsVersionLess checks if v1 < v2
func (vu *VersionUtils) IsVersionLess(v1, v2 string) bool {
	return vu.CompareVersions(v1, v2) < 0
}

// Global instances for easy access
var (
	Files    = &FileUtils{}
	Hashes   = &HashUtils{}
	Strings  = &StringUtils{}
	Slices   = &SliceUtils{}
	Times    = &TimeUtils{}
	JSON     = &JSONUtils{}
	OS       = &OSUtils{}
	Math     = &MathUtils{}
	Versions = &VersionUtils{}
)

// Helper functions

func matchesPattern(path, pattern string) bool {
	matched, err := filepath.Match(pattern, filepath.Base(path))
	if err != nil {
		return false
	}
	return matched
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}

// SafeFileName converts a string to a safe filename
func SafeFileName(s string) string {
	// Replace invalid characters with underscore
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	s = re.ReplaceAllString(s, "_")

	// Trim spaces and dots
	s = strings.TrimSpace(s)
	s = strings.Trim(s, ".")

	// Limit length
	if len(s) > 255 {
		s = s[:255]
	}

	return s
}

// Retry retries a function with backoff
func Retry(attempts int, sleep time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			time.Sleep(sleep * time.Duration(i))
		}
		err = fn()
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

// Parallel processes items in parallel with a limit
func Parallel(items []interface{}, fn func(interface{}) error, concurrency int) []error {
	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}

	sem := make(chan struct{}, concurrency)
	errs := make([]error, len(items))

	for i, item := range items {
		sem <- struct{}{}
		go func(i int, item interface{}) {
			defer func() { <-sem }()
			errs[i] = fn(item)
		}(i, item)
	}

	// Wait for all to complete
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}

	return errs
}
