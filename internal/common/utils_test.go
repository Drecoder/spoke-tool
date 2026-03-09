package common

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileUtils_FileExists(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if !Files.FileExists(tmpfile.Name()) {
		t.Error("FileExists should return true for existing file")
	}

	if Files.FileExists("nonexistent") {
		t.Error("FileExists should return false for nonexistent file")
	}
}

func TestFileUtils_ReadWriteFile(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	testFile := filepath.Join(tmpdir, "test.txt")
	testContent := "hello world"

	err = Files.WriteFile(testFile, testContent)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	content, err := Files.ReadFile(testFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if content != testContent {
		t.Errorf("expected %q, got %q", testContent, content)
	}
}

func TestHashUtils_MD5(t *testing.T) {
	hash := Hashes.MD5("hello")
	expected := "5d41402abc4b2a76b9719d911017c592"
	if hash != expected {
		t.Errorf("MD5: expected %s, got %s", expected, hash)
	}
}

func TestHashUtils_SHA256(t *testing.T) {
	hash := Hashes.SHA256("hello")
	expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if hash != expected {
		t.Errorf("SHA256: expected %s, got %s", expected, hash)
	}
}

func TestStringUtils_SnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"helloWorld", "hello_world"},
		{"HelloWorld", "hello_world"},
		{"hello_world", "hello_world"},
		{"hello-world", "hello-world"}, // keeps hyphens
	}

	for _, tt := range tests {
		result := Strings.SnakeCase(tt.input)
		if result != tt.expected {
			t.Errorf("SnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestStringUtils_CamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "helloWorld"},
		{"hello-world", "helloWorld"},
		{"hello world", "helloWorld"},
		{"HelloWorld", "helloWorld"},
	}

	for _, tt := range tests {
		result := Strings.CamelCase(tt.input)
		if result != tt.expected {
			t.Errorf("CamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestStringUtils_PascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"hello-world", "HelloWorld"},
		{"hello world", "HelloWorld"},
		{"helloWorld", "HelloWorld"},
	}

	for _, tt := range tests {
		result := Strings.PascalCase(tt.input)
		if result != tt.expected {
			t.Errorf("PascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestStringUtils_Truncate(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 8, "hello..."},
		{"hello", 3, "hello"}, // doesn't truncate if max < len
	}

	for _, tt := range tests {
		result := Strings.Truncate(tt.input, tt.max)
		if result != tt.expected {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.max, result, tt.expected)
		}
	}
}

func TestSliceUtils_Contains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if !Slices.Contains(slice, "b") {
		t.Error("Contains should return true for existing item")
	}

	if Slices.Contains(slice, "d") {
		t.Error("Contains should return false for non-existing item")
	}
}

func TestSliceUtils_Unique(t *testing.T) {
	slice := []string{"a", "b", "a", "c", "b"}
	result := Slices.Unique(slice)

	expected := []string{"a", "b", "c"}
	if len(result) != len(expected) {
		t.Fatalf("Unique: expected %v, got %v", expected, result)
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Unique: expected %v, got %v", expected, result)
		}
	}
}

func TestSliceUtils_Intersection(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"b", "c", "d"}

	result := Slices.Intersection(a, b)
	expected := []string{"b", "c"}

	if len(result) != len(expected) {
		t.Fatalf("Intersection: expected %v, got %v", expected, result)
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Intersection: expected %v, got %v", expected, result)
		}
	}
}

func TestSliceUtils_Difference(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"b", "c", "d"}

	result := Slices.Difference(a, b)
	expected := []string{"a"}

	if len(result) != len(expected) {
		t.Fatalf("Difference: expected %v, got %v", expected, result)
	}

	if result[0] != expected[0] {
		t.Errorf("Difference: expected %v, got %v", expected, result)
	}
}

func TestTimeUtils_FormatDuration(t *testing.T) {
	tests := []struct {
		d        time.Duration
		expected string
	}{
		{100 * time.Nanosecond, "100 ns"},
		{1500 * time.Nanosecond, "1.50 µs"},
		{2 * time.Millisecond, "2.00 ms"},
		{1500 * time.Millisecond, "1.50 s"},
		{90 * time.Second, "1m30s"},
	}

	for _, tt := range tests {
		result := Times.FormatDuration(tt.d)
		if result != tt.expected {
			t.Errorf("FormatDuration(%v) = %q, want %q", tt.d, result, tt.expected)
		}
	}
}

func TestMathUtils_MinMax(t *testing.T) {
	if Math.Min(5, 10) != 5 {
		t.Error("Min(5,10) should be 5")
	}
	if Math.Max(5, 10) != 10 {
		t.Error("Max(5,10) should be 10")
	}
}

func TestMathUtils_Clamp(t *testing.T) {
	tests := []struct {
		value    int
		min      int
		max      int
		expected int
	}{
		{5, 0, 10, 5},
		{-5, 0, 10, 0},
		{15, 0, 10, 10},
	}

	for _, tt := range tests {
		result := Math.Clamp(tt.value, tt.min, tt.max)
		if result != tt.expected {
			t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.value, tt.min, tt.max, result, tt.expected)
		}
	}
}

func TestVersionUtils_CompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.0.1", -1},
		{"1.0.1", "1.0.0", 1},
		{"v1.0.0", "1.0.0", 0},
		{"1.0", "1.0.0", 0},
		{"2.0.0", "1.9.9", 1},
	}

	for _, tt := range tests {
		result := Versions.CompareVersions(tt.v1, tt.v2)
		if result != tt.expected {
			t.Errorf("CompareVersions(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
		}
	}
}

func TestSafeFileName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello.txt", "hello.txt"},
		{"hello/world:test", "hello_world_test"},
		{"  test  ", "test"},
		{".test.", "test"},
	}

	for _, tt := range tests {
		result := SafeFileName(tt.input)
		if result != tt.expected {
			t.Errorf("SafeFileName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
