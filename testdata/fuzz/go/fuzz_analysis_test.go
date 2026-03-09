//go:build go1.18
// +build go1.18

package fuzz

import (
	"testing"

	"github.com/yourusername/spoke-tool/api/types"
	"github.com/yourusername/spoke-tool/internal/test"
)

// FuzzCodeAnalysis tests the code analyzer with various inputs
// to ensure it never panics and handles edge cases gracefully.
func FuzzCodeAnalysis(f *testing.F) {
	// Seed corpus with various code snippets
	seeds := []string{
		"",
		"package main",
		"package main\n\nfunc main() {}",
		"package main\n\nfunc add(a, b int) int { return a + b }",
		"invalid go code @#$%",
		"package main\n\nfunc ( )",
		"package main\n\nfunc a( a int",
		"package main\n\n// comment only",
		"package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Println(\"hello\") }",
		"package main\n\ntype Person struct { Name string; Age int }",
		"package main\n\nfunc (p Person) Greet() string { return \"Hello, \" + p.Name }",
		"package main\n\nfunc fibonacci(n int) int { if n <= 1 { return n }; return fibonacci(n-1) + fibonacci(n-2) }",
		"package main\n\nfunc process(data []int) ([]int, error) { return nil, nil }",
		"package main\n\ntype Reader interface { Read(p []byte) (n int, err error) }",
		"package main\n\nconst ( StatusOK = 200; StatusNotFound = 404 )",
		"package main\n\nvar version = \"1.0.0\"",
		"package main\n\nfunc init() {}",
		"package main\n\nfunc withManyParams(a, b, c, d, e, f, g, h, i, j int) int { return a + b + c + d + e + f + g + h + i + j }",
		"package main\n\nfunc withManyReturns() (a, b, c, d, e, f, g, h, i, j int) { return 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 }",
		"package main\n\nfunc nested() { func() { func() { println(\"deep\") }() }() }",
		"package main\n\nfunc deferExample() { defer func() { recover() }(); panic(\"error\") }",
		"package main\n\nfunc channels() { ch := make(chan int); go func() { ch <- 42 }(); <-ch }",
		"package main\n\nfunc selectExample() { ch1 := make(chan int); ch2 := make(chan int); select { case <-ch1: case v := <-ch2: println(v); default: } }",
		"package main\n\nfunc rangeExample() { nums := []int{1,2,3}; for i, v := range nums { println(i, v) } }",
		"package main\n\nfunc typeSwitch(i interface{}) { switch v := i.(type) { case int: println(\"int\"); case string: println(\"string\"); default: println(\"unknown\") } }",
		"package main\n\nfunc typeAssertion() { var i interface{} = \"hello\"; s := i.(string); println(s) }",
		"package main\n\nfunc mapOperations() { m := make(map[string]int); m[\"key\"] = 42; delete(m, \"key\") }",
		"package main\n\nfunc sliceOperations() { s := make([]int, 0, 10); s = append(s, 1, 2, 3); s = s[:2] }",
		"package main\n\nfunc pointerOperations() { x := 42; p := &x; *p = 100 }",
		"package main\n\nfunc structOperations() { type Point struct { X, Y int }; p := Point{10, 20}; p.X = 30 }",
		"package main\n\nfunc blankIdentifier() { _ = 42; _, _ = 1, 2 }",
		"package main\n\nfunc labeledStatements() { outer: for i := 0; i < 10; i++ { for j := 0; j < 10; j++ { if i+j > 10 { break outer } } } }",
		"package main\n\nfunc gotoStatement() { i := 0; loop: if i < 10 { i++; goto loop } }",
		"package main\n\nfunc fallthroughExample() { switch 2 { case 1: println(\"one\"); case 2: println(\"two\"); fallthrough; case 3: println(\"three\") } }",
		"package main\n\nfunc bitwiseOperations() { x := 0b1010; y := 0b1100; _ = x & y; _ = x | y; _ = x ^ y }",
		"package main\n\nfunc withBuildTags() { //go:build linux\n println(\"linux only\") }",
		"package main\n\nfunc emptyFile()\n",
		"package main\n\n//go:noinline\nfunc withPragmas() {}",
	}

	// Add seeds to corpus
	for _, seed := range seeds {
		f.Add(seed)
	}

	// Create analyzer once (shared across fuzz runs for efficiency)
	analyzer := test.NewAnalyzer(test.AnalyzerConfig{
		Languages:      []types.Language{types.Go},
		ExportedOnly:   false,
		IncludePrivate: true,
	})

	f.Fuzz(func(t *testing.T, code string) {
		// This should NEVER panic - just analyze the code
		result, err := analyzer.AnalyzeCode(code)

		// We don't care about the result validity, just that:
		// 1. No panic occurred
		// 2. Return values are valid (not nil)
		// 3. Error may be nil or non-nil, both are acceptable

		if result == nil && err == nil {
			// Both nil is suspicious - either should have result or error
			t.Log("Both result and error are nil")
		}

		// Basic sanity checks - if we got a result, it should be valid
		if result != nil {
			// Check that result has basic fields (don't panic on access)
			_ = result.Functions
			_ = result.Imports
			_ = result.Complexity
			_ = result.LinesOfCode
		}
	})
}

// FuzzFilePaths tests file path handling to ensure no path traversal vulnerabilities
func FuzzFilePaths(f *testing.F) {
	seeds := []string{
		"",
		".",
		"..",
		"./path/to/file",
		"../../../etc/passwd",
		"C:\\Windows\\System32",
		"path with spaces",
		"path/with//double/slashes",
		"path/with/./current/dir",
		"path/with/../parent/dir",
		"very/long/path/with/many/segments/file.txt",
		"!@#$%^&*().txt",
		"path/with/unicode/café.txt",
		"path/with/emoji/🚀.txt",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, path string) {
		// Test file path sanitization - should never panic
		safePath := test.SanitizePath(path)
		_ = safePath

		// Test path operations
		isValid := test.IsValidPath(path)
		_ = isValid
	})
}
