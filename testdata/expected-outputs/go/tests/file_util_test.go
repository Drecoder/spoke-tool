package fileutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "existing file",
			path:     testFile,
			expected: true,
		},
		{
			name:     "non-existent file",
			path:     filepath.Join(tmpDir, "nonexistent.txt"),
			expected: false,
		},
		{
			name:     "directory path",
			path:     tmpDir,
			expected: true, // directories also exist
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FileExists(tt.path)
			if result != tt.expected {
				t.Errorf("FileExists(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "directory path",
			path:     tmpDir,
			expected: true,
		},
		{
			name:     "file path",
			path:     testFile,
			expected: false,
		},
		{
			name:     "non-existent path",
			path:     filepath.Join(tmpDir, "nonexistent"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDir(tt.path)
			if result != tt.expected {
				t.Errorf("IsDir(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files with different content
	testFile1 := filepath.Join(tmpDir, "test1.txt")
	content1 := "hello world"
	err := os.WriteFile(testFile1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	testFile2 := filepath.Join(tmpDir, "test2.txt")
	content2 := "line1\nline2\nline3"
	err = os.WriteFile(testFile2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	emptyFile := filepath.Join(tmpDir, "empty.txt")
	err = os.WriteFile(emptyFile, []byte{}, 0644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	tests := []struct {
		name      string
		path      string
		expected  string
		wantError bool
	}{
		{
			name:      "read existing file",
			path:      testFile1,
			expected:  content1,
			wantError: false,
		},
		{
			name:      "read multiline file",
			path:      testFile2,
			expected:  content2,
			wantError: false,
		},
		{
			name:      "read empty file",
			path:      emptyFile,
			expected:  "",
			wantError: false,
		},
		{
			name:      "read non-existent file",
			path:      filepath.Join(tmpDir, "nonexistent.txt"),
			expected:  "",
			wantError: true,
		},
		{
			name:      "read directory",
			path:      tmpDir,
			expected:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := ReadFile(tt.path)

			if tt.wantError {
				if err == nil {
					t.Errorf("ReadFile(%q) expected error, got nil", tt.path)
				}
			} else {
				if err != nil {
					t.Errorf("ReadFile(%q) unexpected error: %v", tt.path, err)
				}
				if content != tt.expected {
					t.Errorf("ReadFile(%q) = %q, want %q", tt.path, content, tt.expected)
				}
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		path      string
		content   string
		setup     func() error
		wantError bool
	}{
		{
			name:      "write new file",
			path:      filepath.Join(tmpDir, "new.txt"),
			content:   "test content",
			wantError: false,
		},
		{
			name:      "write to nested directory",
			path:      filepath.Join(tmpDir, "nested", "deep", "file.txt"),
			content:   "nested file",
			wantError: false,
		},
		{
			name:    "overwrite existing file",
			path:    filepath.Join(tmpDir, "existing.txt"),
			content: "new content",
			setup: func() error {
				return os.WriteFile(filepath.Join(tmpDir, "existing.txt"), []byte("old content"), 0644)
			},
			wantError: false,
		},
		{
			name:      "write empty content",
			path:      filepath.Join(tmpDir, "empty.txt"),
			content:   "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run setup if needed
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Write file
			err := WriteFile(tt.path, tt.content)

			if tt.wantError {
				if err == nil {
					t.Errorf("WriteFile(%q) expected error, got nil", tt.path)
				}
				return
			}

			if err != nil {
				t.Fatalf("WriteFile(%q) unexpected error: %v", tt.path, err)
			}

			// Verify file was written correctly
			content, err := os.ReadFile(tt.path)
			if err != nil {
				t.Fatalf("Failed to read written file: %v", err)
			}

			if string(content) != tt.content {
				t.Errorf("File content = %q, want %q", string(content), tt.content)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	srcFile := filepath.Join(tmpDir, "source.txt")
	srcContent := "source content"
	err := os.WriteFile(srcFile, []byte(srcContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	tests := []struct {
		name      string
		src       string
		dst       string
		setup     func() error
		wantError bool
	}{
		{
			name:      "copy to new file",
			src:       srcFile,
			dst:       filepath.Join(tmpDir, "dest1.txt"),
			wantError: false,
		},
		{
			name:      "copy to nested destination",
			src:       srcFile,
			dst:       filepath.Join(tmpDir, "nested", "dest2.txt"),
			wantError: false,
		},
		{
			name: "copy to existing file (overwrite)",
			src:  srcFile,
			dst:  filepath.Join(tmpDir, "dest3.txt"),
			setup: func() error {
				return os.WriteFile(filepath.Join(tmpDir, "dest3.txt"), []byte("old"), 0644)
			},
			wantError: false,
		},
		{
			name:      "copy non-existent source",
			src:       filepath.Join(tmpDir, "nonexistent.txt"),
			dst:       filepath.Join(tmpDir, "dest4.txt"),
			wantError: true,
		},
		{
			name:      "copy to invalid path",
			src:       srcFile,
			dst:       "/invalid/path/that/cant/be/created/file.txt",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run setup if needed
			if tt.setup != nil {
				if err := tt.setup(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			// Copy file
			err := CopyFile(tt.src, tt.dst)

			if tt.wantError {
				if err == nil {
					t.Errorf("CopyFile(%q, %q) expected error, got nil", tt.src, tt.dst)
				}
				return
			}

			if err != nil {
				t.Fatalf("CopyFile(%q, %q) unexpected error: %v", tt.src, tt.dst, err)
			}

			// Verify destination file exists and has correct content
			if !FileExists(tt.dst) {
				t.Errorf("Destination file %q does not exist after copy", tt.dst)
			}

			dstContent, err := ReadFile(tt.dst)
			if err != nil {
				t.Fatalf("Failed to read destination file: %v", err)
			}

			if dstContent != srcContent {
				t.Errorf("Destination content = %q, want %q", dstContent, srcContent)
			}
		})
	}
}

func TestListFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a mix of files and directories
	files := []string{
		"file1.txt",
		"file2.go",
		"file3_test.go",
		"README.md",
		".hidden.txt",
		"subdir/file4.txt",
		"subdir/file5.go",
		"subdir/nested/file6.txt",
	}

	for _, file := range files {
		fullPath := filepath.Join(tmpDir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		err = os.WriteFile(fullPath, []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	tests := []struct {
		name      string
		dir       string
		pattern   string
		expected  []string
		wantError bool
	}{
		{
			name:    "list all files",
			dir:     tmpDir,
			pattern: "",
			expected: []string{
				".hidden.txt",
				"README.md",
				"file1.txt",
				"file2.go",
				"file3_test.go",
				"subdir/file4.txt",
				"subdir/file5.go",
				"subdir/nested/file6.txt",
			},
			wantError: false,
		},
		{
			name:    "list with .txt pattern",
			dir:     tmpDir,
			pattern: "*.txt",
			expected: []string{
				"file1.txt",
				"subdir/file4.txt",
				"subdir/nested/file6.txt",
			},
			wantError: false,
		},
		{
			name:    "list with .go pattern",
			dir:     tmpDir,
			pattern: "*.go",
			expected: []string{
				"file2.go",
				"subdir/file5.go",
			},
			wantError: false,
		},
		{
			name:    "list with test pattern",
			dir:     tmpDir,
			pattern: "*_test.go",
			expected: []string{
				"file3_test.go",
			},
			wantError: false,
		},
		{
			name:    "list with hidden pattern",
			dir:     tmpDir,
			pattern: ".*",
			expected: []string{
				".hidden.txt",
			},
			wantError: false,
		},
		{
			name:      "list non-existent directory",
			dir:       filepath.Join(tmpDir, "nonexistent"),
			pattern:   "",
			expected:  nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := ListFiles(tt.dir, tt.pattern)

			if tt.wantError {
				if err == nil {
					t.Errorf("ListFiles(%q, %q) expected error, got nil", tt.dir, tt.pattern)
				}
				return
			}

			if err != nil {
				t.Fatalf("ListFiles(%q, %q) unexpected error: %v", tt.dir, tt.pattern, err)
			}

			// Convert full paths to relative for comparison
			var relFiles []string
			for _, f := range files {
				rel, err := filepath.Rel(tmpDir, f)
				if err != nil {
					t.Fatalf("Failed to get relative path: %v", err)
				}
				relFiles = append(relFiles, rel)
			}

			// Compare slices (order doesn't matter)
			if !stringSlicesEqual(relFiles, tt.expected) {
				t.Errorf("ListFiles(%q, %q) = %v, want %v", tt.dir, tt.pattern, relFiles, tt.expected)
			}
		})
	}
}

func TestGetFileInfo(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with known content
	testFile := filepath.Join(tmpDir, "info.txt")
	content := "test content"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a test directory
	testDir := filepath.Join(tmpDir, "testdir")
	err = os.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name      string
		path      string
		check     func(*testing.T, *FileInfo)
		wantError bool
	}{
		{
			name: "get file info",
			path: testFile,
			check: func(t *testing.T, info *FileInfo) {
				if info.Name != "info.txt" {
					t.Errorf("Name = %q, want %q", info.Name, "info.txt")
				}
				if info.Size != int64(len(content)) {
					t.Errorf("Size = %d, want %d", info.Size, len(content))
				}
				if info.IsDir {
					t.Error("IsDir = true, want false")
				}
				if info.Extension != ".txt" {
					t.Errorf("Extension = %q, want %q", info.Extension, ".txt")
				}
			},
			wantError: false,
		},
		{
			name: "get directory info",
			path: testDir,
			check: func(t *testing.T, info *FileInfo) {
				if info.Name != "testdir" {
					t.Errorf("Name = %q, want %q", info.Name, "testdir")
				}
				if !info.IsDir {
					t.Error("IsDir = false, want true")
				}
				if info.Extension != "" {
					t.Errorf("Extension = %q, want empty string", info.Extension)
				}
			},
			wantError: false,
		},
		{
			name:      "get non-existent file info",
			path:      filepath.Join(tmpDir, "nonexistent.txt"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := GetFileInfo(tt.path)

			if tt.wantError {
				if err == nil {
					t.Errorf("GetFileInfo(%q) expected error, got nil", tt.path)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetFileInfo(%q) unexpected error: %v", tt.path, err)
			}
			if info == nil {
				t.Fatal("GetFileInfo returned nil info")
			}

			if tt.check != nil {
				tt.check(t, info)
			}
		})
	}
}

func TestFileInfo_Hash(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two files with same content
	file1 := filepath.Join(tmpDir, "file1.txt")
	content1 := "same content"
	err := os.WriteFile(file1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	file2 := filepath.Join(tmpDir, "file2.txt")
	err = os.WriteFile(file2, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	// Create file with different content
	file3 := filepath.Join(tmpDir, "file3.txt")
	err = os.WriteFile(file3, []byte("different content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file3: %v", err)
	}

	t.Run("same content same hash", func(t *testing.T) {
		info1, err := GetFileInfo(file1)
		if err != nil {
			t.Fatalf("Failed to get file1 info: %v", err)
		}

		info2, err := GetFileInfo(file2)
		if err != nil {
			t.Fatalf("Failed to get file2 info: %v", err)
		}

		if info1.Hash != info2.Hash {
			t.Errorf("Hashes don't match: %q vs %q", info1.Hash, info2.Hash)
		}
	})

	t.Run("different content different hash", func(t *testing.T) {
		info1, err := GetFileInfo(file1)
		if err != nil {
			t.Fatalf("Failed to get file1 info: %v", err)
		}

		info3, err := GetFileInfo(file3)
		if err != nil {
			t.Fatalf("Failed to get file3 info: %v", err)
		}

		if info1.Hash == info3.Hash {
			t.Error("Expected different hashes for different content")
		}
	})
}

// Helper function to compare string slices ignoring order
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Create maps for existence check
	amap := make(map[string]bool)
	for _, s := range a {
		amap[s] = true
	}

	for _, s := range b {
		if !amap[s] {
			return false
		}
	}

	return true
}

func BenchmarkReadFile(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "bench.txt")
	content := "benchmark content\n" + string(make([]byte, 1024))
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadFile(testFile)
	}
}

func BenchmarkWriteFile(b *testing.B) {
	tmpDir := b.TempDir()
	content := "benchmark content\n" + string(make([]byte, 1024))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testFile := filepath.Join(tmpDir, "bench.txt")
		WriteFile(testFile, content)
	}
}

func BenchmarkListFiles(b *testing.B) {
	tmpDir := b.TempDir()

	// Create many files
	for i := 0; i < 100; i++ {
		file := filepath.Join(tmpDir, "file%d.txt")
		err := os.WriteFile(file, []byte("content"), 0644)
		if err != nil {
			b.Fatalf("Failed to create file: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ListFiles(tmpDir, "")
	}
}
