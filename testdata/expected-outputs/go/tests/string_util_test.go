package stringutil

import (
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal string",
			input:    "hello",
			expected: "olleh",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "a",
		},
		{
			name:     "palindrome",
			input:    "racecar",
			expected: "racecar",
		},
		{
			name:     "string with spaces",
			input:    "hello world",
			expected: "dlrow olleh",
		},
		{
			name:     "string with punctuation",
			input:    "hello!",
			expected: "!olleh",
		},
		{
			name:     "unicode string",
			input:    "Hello, 世界",
			expected: "界世 ,olleH",
		},
		{
			name:     "string with numbers",
			input:    "abc123",
			expected: "321cba",
		},
		{
			name:     "string with mixed case",
			input:    "HelloWorld",
			expected: "dlroWolleH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reverse(tt.input)
			if result != tt.expected {
				t.Errorf("Reverse(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase string",
			input:    "hello",
			expected: "HELLO",
		},
		{
			name:     "uppercase string",
			input:    "HELLO",
			expected: "HELLO",
		},
		{
			name:     "mixed case",
			input:    "HeLLo WoRLd",
			expected: "HELLO WORLD",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "string with numbers",
			input:    "abc123",
			expected: "ABC123",
		},
		{
			name:     "string with punctuation",
			input:    "hello!@#",
			expected: "HELLO!@#",
		},
		{
			name:     "unicode string",
			input:    "hello 世界",
			expected: "HELLO 世界", // Unicode characters remain unchanged
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToUpper(tt.input)
			if result != tt.expected {
				t.Errorf("ToUpper(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "uppercase string",
			input:    "HELLO",
			expected: "hello",
		},
		{
			name:     "lowercase string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "mixed case",
			input:    "HeLLo WoRLd",
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "string with numbers",
			input:    "ABC123",
			expected: "abc123",
		},
		{
			name:     "string with punctuation",
			input:    "HELLO!@#",
			expected: "hello!@#",
		},
		{
			name:     "unicode string",
			input:    "HELLO 世界",
			expected: "hello 世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToLower(tt.input)
			if result != tt.expected {
				t.Errorf("ToLower(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTrimSpace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no spaces",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "leading spaces",
			input:    "  hello",
			expected: "hello",
		},
		{
			name:     "trailing spaces",
			input:    "hello  ",
			expected: "hello",
		},
		{
			name:     "both sides",
			input:    "  hello  ",
			expected: "hello",
		},
		{
			name:     "internal spaces",
			input:    "hello world",
			expected: "hello world", // Internal spaces preserved
		},
		{
			name:     "tabs and newlines",
			input:    "\t\n hello \t\n",
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "     ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TrimSpace(tt.input)
			if result != tt.expected {
				t.Errorf("TrimSpace(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		sep      string
		expected []string
	}{
		{
			name:     "split by space",
			input:    "hello world",
			sep:      " ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "split by comma",
			input:    "a,b,c",
			sep:      ",",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "multiple separators",
			input:    "a,,b,c",
			sep:      ",",
			expected: []string{"a", "", "b", "c"},
		},
		{
			name:     "no separator",
			input:    "hello",
			sep:      ",",
			expected: []string{"hello"},
		},
		{
			name:     "empty string",
			input:    "",
			sep:      ",",
			expected: []string{""},
		},
		{
			name:     "separator at start",
			input:    ",a,b",
			sep:      ",",
			expected: []string{"", "a", "b"},
		},
		{
			name:     "separator at end",
			input:    "a,b,",
			sep:      ",",
			expected: []string{"a", "b", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Split(tt.input, tt.sep)
			if !stringSlicesEqual(result, tt.expected) {
				t.Errorf("Split(%q, %q) = %v, want %v", tt.input, tt.sep, result, tt.expected)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		sep      string
		expected string
	}{
		{
			name:     "join with space",
			input:    []string{"hello", "world"},
			sep:      " ",
			expected: "hello world",
		},
		{
			name:     "join with comma",
			input:    []string{"a", "b", "c"},
			sep:      ",",
			expected: "a,b,c",
		},
		{
			name:     "join with empty separator",
			input:    []string{"a", "b", "c"},
			sep:      "",
			expected: "abc",
		},
		{
			name:     "single element",
			input:    []string{"hello"},
			sep:      ",",
			expected: "hello",
		},
		{
			name:     "empty slice",
			input:    []string{},
			sep:      ",",
			expected: "",
		},
		{
			name:     "slice with empty strings",
			input:    []string{"a", "", "c"},
			sep:      ",",
			expected: "a,,c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Join(tt.input, tt.sep)
			if result != tt.expected {
				t.Errorf("Join(%v, %q) = %q, want %q", tt.input, tt.sep, result, tt.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "contains substring",
			s:        "hello world",
			substr:   "world",
			expected: true,
		},
		{
			name:     "does not contain",
			s:        "hello world",
			substr:   "xyz",
			expected: false,
		},
		{
			name:     "empty substring",
			s:        "hello",
			substr:   "",
			expected: true, // Empty string is contained in any string
		},
		{
			name:     "empty string",
			s:        "",
			substr:   "hello",
			expected: false,
		},
		{
			name:     "both empty",
			s:        "",
			substr:   "",
			expected: true,
		},
		{
			name:     "case sensitive",
			s:        "Hello World",
			substr:   "hello",
			expected: false,
		},
		{
			name:     "at start",
			s:        "hello world",
			substr:   "hello",
			expected: true,
		},
		{
			name:     "at end",
			s:        "hello world",
			substr:   "world",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("Contains(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		old      string
		new      string
		n        int
		expected string
	}{
		{
			name:     "replace all",
			s:        "hello world hello",
			old:      "hello",
			new:      "hi",
			n:        -1,
			expected: "hi world hi",
		},
		{
			name:     "replace once",
			s:        "hello world hello",
			old:      "hello",
			new:      "hi",
			n:        1,
			expected: "hi world hello",
		},
		{
			name:     "replace twice",
			s:        "hello world hello",
			old:      "hello",
			new:      "hi",
			n:        2,
			expected: "hi world hi",
		},
		{
			name:     "no matches",
			s:        "hello world",
			old:      "xyz",
			new:      "abc",
			n:        -1,
			expected: "hello world",
		},
		{
			name:     "empty old string",
			s:        "hello",
			old:      "",
			new:      "x",
			n:        -1,
			expected: "xhellox", // Empty matches between each character
		},
		{
			name:     "replace with empty",
			s:        "hello world",
			old:      "o",
			new:      "",
			n:        -1,
			expected: "hell wrld",
		},
		{
			name:     "empty string",
			s:        "",
			old:      "a",
			new:      "b",
			n:        -1,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Replace(tt.s, tt.old, tt.new, tt.n)
			if result != tt.expected {
				t.Errorf("Replace(%q, %q, %q, %d) = %q, want %q",
					tt.s, tt.old, tt.new, tt.n, result, tt.expected)
			}
		})
	}
}

func TestReplaceAll(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		old      string
		new      string
		expected string
	}{
		{
			name:     "replace all occurrences",
			s:        "hello world hello",
			old:      "hello",
			new:      "hi",
			expected: "hi world hi",
		},
		{
			name:     "no matches",
			s:        "hello world",
			old:      "xyz",
			new:      "abc",
			expected: "hello world",
		},
		{
			name:     "multiple matches",
			s:        "aaaaa",
			old:      "aa",
			new:      "b",
			expected: "bba", // Overlapping not replaced
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceAll(tt.s, tt.old, tt.new)
			if result != tt.expected {
				t.Errorf("ReplaceAll(%q, %q, %q) = %q, want %q",
					tt.s, tt.old, tt.new, result, tt.expected)
			}
		})
	}
}

func TestIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected int
	}{
		{
			name:     "substring present",
			s:        "hello world",
			substr:   "world",
			expected: 6,
		},
		{
			name:     "substring not present",
			s:        "hello world",
			substr:   "xyz",
			expected: -1,
		},
		{
			name:     "at beginning",
			s:        "hello world",
			substr:   "hello",
			expected: 0,
		},
		{
			name:     "empty substring",
			s:        "hello",
			substr:   "",
			expected: 0,
		},
		{
			name:     "empty string",
			s:        "",
			substr:   "hello",
			expected: -1,
		},
		{
			name:     "both empty",
			s:        "",
			substr:   "",
			expected: 0,
		},
		{
			name:     "case sensitive",
			s:        "Hello World",
			substr:   "hello",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IndexOf(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("IndexOf(%q, %q) = %d, want %d", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestLastIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected int
	}{
		{
			name:     "single occurrence",
			s:        "hello world",
			substr:   "world",
			expected: 6,
		},
		{
			name:     "multiple occurrences",
			s:        "hello world hello",
			substr:   "hello",
			expected: 12,
		},
		{
			name:     "not present",
			s:        "hello world",
			substr:   "xyz",
			expected: -1,
		},
		{
			name:     "at beginning",
			s:        "hello world",
			substr:   "hello",
			expected: 0,
		},
		{
			name:     "empty substring",
			s:        "hello",
			substr:   "",
			expected: 5, // Length of string
		},
		{
			name:     "empty string",
			s:        "",
			substr:   "hello",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LastIndexOf(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("LastIndexOf(%q, %q) = %d, want %d", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		prefix   string
		expected bool
	}{
		{
			name:     "has prefix",
			s:        "hello world",
			prefix:   "hello",
			expected: true,
		},
		{
			name:     "does not have prefix",
			s:        "hello world",
			prefix:   "world",
			expected: false,
		},
		{
			name:     "empty prefix",
			s:        "hello",
			prefix:   "",
			expected: true,
		},
		{
			name:     "empty string",
			s:        "",
			prefix:   "hello",
			expected: false,
		},
		{
			name:     "both empty",
			s:        "",
			prefix:   "",
			expected: true,
		},
		{
			name:     "case sensitive",
			s:        "Hello World",
			prefix:   "hello",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasPrefix(tt.s, tt.prefix)
			if result != tt.expected {
				t.Errorf("HasPrefix(%q, %q) = %v, want %v", tt.s, tt.prefix, result, tt.expected)
			}
		})
	}
}

func TestHasSuffix(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		suffix   string
		expected bool
	}{
		{
			name:     "has suffix",
			s:        "hello world",
			suffix:   "world",
			expected: true,
		},
		{
			name:     "does not have suffix",
			s:        "hello world",
			suffix:   "hello",
			expected: false,
		},
		{
			name:     "empty suffix",
			s:        "hello",
			suffix:   "",
			expected: true,
		},
		{
			name:     "empty string",
			s:        "",
			suffix:   "hello",
			expected: false,
		},
		{
			name:     "both empty",
			s:        "",
			suffix:   "",
			expected: true,
		},
		{
			name:     "case sensitive",
			s:        "Hello World",
			suffix:   "world",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasSuffix(tt.s, tt.suffix)
			if result != tt.expected {
				t.Errorf("HasSuffix(%q, %q) = %v, want %v", tt.s, tt.suffix, result, tt.expected)
			}
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected int
	}{
		{
			name:     "single occurrence",
			s:        "hello world",
			substr:   "world",
			expected: 1,
		},
		{
			name:     "multiple occurrences",
			s:        "hello hello hello",
			substr:   "hello",
			expected: 3,
		},
		{
			name:     "overlapping",
			s:        "aaaaa",
			substr:   "aa",
			expected: 2, // Non-overlapping count
		},
		{
			name:     "no occurrences",
			s:        "hello world",
			substr:   "xyz",
			expected: 0,
		},
		{
			name:     "empty substring",
			s:        "hello",
			substr:   "",
			expected: 6, // Length + 1
		},
		{
			name:     "empty string",
			s:        "",
			substr:   "hello",
			expected: 0,
		},
		{
			name:     "both empty",
			s:        "",
			substr:   "",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Count(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("Count(%q, %q) = %d, want %d", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		count     int
		expected  string
		wantError bool
	}{
		{
			name:      "repeat multiple times",
			s:         "hello",
			count:     3,
			expected:  "hellohellohello",
			wantError: false,
		},
		{
			name:      "repeat once",
			s:         "hello",
			count:     1,
			expected:  "hello",
			wantError: false,
		},
		{
			name:      "repeat zero",
			s:         "hello",
			count:     0,
			expected:  "",
			wantError: false,
		},
		{
			name:      "negative count",
			s:         "hello",
			count:     -1,
			expected:  "",
			wantError: true,
		},
		{
			name:      "empty string",
			s:         "",
			count:     5,
			expected:  "",
			wantError: false,
		},
		{
			name:      "large count",
			s:         "a",
			count:     1000,
			expected:  string(make([]byte, 1000)),
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Repeat(tt.s, tt.count)

			if tt.wantError {
				if err == nil {
					t.Errorf("Repeat(%q, %d) expected error, got nil", tt.s, tt.count)
				}
			} else {
				if err != nil {
					t.Errorf("Repeat(%q, %d) unexpected error: %v", tt.s, tt.count, err)
				}
				if result != tt.expected {
					t.Errorf("Repeat(%q, %d) = %q, want %q", tt.s, tt.count, result, tt.expected)
				}
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "camelCase",
			input:    "camelCase",
			expected: "camel_case",
		},
		{
			name:     "PascalCase",
			input:    "PascalCase",
			expected: "pascal_case",
		},
		{
			name:     "already snake_case",
			input:    "snake_case",
			expected: "snake_case",
		},
		{
			name:     "with numbers",
			input:    "userID123",
			expected: "user_id123",
		},
		{
			name:     "multiple words",
			input:    "thisIsALongVariableName",
			expected: "this_is_a_long_variable_name",
		},
		{
			name:     "acronyms",
			input:    "parseJSON",
			expected: "parse_json",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "hello",
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "snake_case",
			input:    "snake_case",
			expected: "snakeCase",
		},
		{
			name:     "multiple words",
			input:    "this_is_a_long_variable_name",
			expected: "thisIsALongVariableName",
		},
		{
			name:     "already camelCase",
			input:    "camelCase",
			expected: "camelCase",
		},
		{
			name:     "with numbers",
			input:    "user_id_123",
			expected: "userId123",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "leading underscore",
			input:    "_private_var",
			expected: "PrivateVar", // or "_privateVar"? Implementation dependent
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "snake_case",
			input:    "snake_case",
			expected: "SnakeCase",
		},
		{
			name:     "camelCase",
			input:    "camelCase",
			expected: "CamelCase",
		},
		{
			name:     "multiple words",
			input:    "this_is_a_long_variable_name",
			expected: "ThisIsALongVariableName",
		},
		{
			name:     "with numbers",
			input:    "user_id_123",
			expected: "UserId123",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "hello",
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		maxLen   int
		expected string
	}{
		{
			name:     "shorter than max",
			s:        "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "equal to max",
			s:        "hello",
			maxLen:   5,
			expected: "hello",
		},
		{
			name:     "longer than max",
			s:        "hello world",
			maxLen:   8,
			expected: "hello...",
		},
		{
			name:     "maxLen too small",
			s:        "hello",
			maxLen:   2,
			expected: "he", // or "..."? Implementation dependent
		},
		{
			name:     "empty string",
			s:        "",
			maxLen:   5,
			expected: "",
		},
		{
			name:     "maxLen zero",
			s:        "hello",
			maxLen:   0,
			expected: "",
		},
		{
			name:     "maxLen negative",
			s:        "hello",
			maxLen:   -1,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Truncate(tt.s, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.s, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no duplicates",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "consecutive duplicates",
			input:    "hello  world",
			expected: "hello world", // Only one space
		},
		{
			name:     "word duplicates",
			input:    "hello hello world",
			expected: "hello world",
		},
		{
			name:     "line duplicates",
			input:    "line1\nline2\nline1\nline3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single line",
			input:    "hello",
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveDuplicates(tt.input)
			if result != tt.expected {
				t.Errorf("RemoveDuplicates(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIndent(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		indent   string
		expected string
	}{
		{
			name:     "single line",
			s:        "hello",
			indent:   "  ",
			expected: "  hello",
		},
		{
			name:     "multiple lines",
			s:        "line1\nline2\nline3",
			indent:   "  ",
			expected: "  line1\n  line2\n  line3",
		},
		{
			name:     "empty string",
			s:        "",
			indent:   "  ",
			expected: "",
		},
		{
			name:     "empty indent",
			s:        "hello",
			indent:   "",
			expected: "hello",
		},
		{
			name:     "with trailing newline",
			s:        "hello\n",
			indent:   "  ",
			expected: "  hello\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Indent(tt.s, tt.indent)
			if result != tt.expected {
				t.Errorf("Indent(%q, %q) = %q, want %q", tt.s, tt.indent, result, tt.expected)
			}
		})
	}
}

func TestContainsAny(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		substrings []string
		expected   bool
	}{
		{
			name:       "contains one",
			s:          "hello world",
			substrings: []string{"foo", "world", "bar"},
			expected:   true,
		},
		{
			name:       "contains multiple",
			s:          "hello world",
			substrings: []string{"hello", "world"},
			expected:   true,
		},
		{
			name:       "contains none",
			s:          "hello world",
			substrings: []string{"foo", "bar", "baz"},
			expected:   false,
		},
		{
			name:       "empty substrings",
			s:          "hello",
			substrings: []string{},
			expected:   false,
		},
		{
			name:       "empty string",
			s:          "",
			substrings: []string{"hello"},
			expected:   false,
		},
		{
			name:       "contains empty string",
			s:          "hello",
			substrings: []string{""},
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsAny(tt.s, tt.substrings)
			if result != tt.expected {
				t.Errorf("ContainsAny(%q, %v) = %v, want %v", tt.s, tt.substrings, result, tt.expected)
			}
		})
	}
}

// Helper function for slice comparison
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Benchmarks
func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Reverse("hello world this is a test string")
	}
}

func BenchmarkToUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToUpper("hello world")
	}
}

func BenchmarkSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Split("a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p", ",")
	}
}

func BenchmarkJoin(b *testing.B) {
	parts := []string{"a", "b", "c", "d", "e", "f", "g"}
	for i := 0; i < b.N; i++ {
		Join(parts, ",")
	}
}

func BenchmarkContains(b *testing.B) {
	s := "the quick brown fox jumps over the lazy dog"
	for i := 0; i < b.N; i++ {
		Contains(s, "fox")
	}
}

func BenchmarkReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReplaceAll("the quick brown fox jumps over the lazy dog", "the", "a")
	}
}

func BenchmarkToSnakeCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToSnakeCase("thisIsAVeryLongVariableNameWithManyWords")
	}
}

func BenchmarkToCamelCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToCamelCase("this_is_a_very_long_variable_name_with_many_words")
	}
}
