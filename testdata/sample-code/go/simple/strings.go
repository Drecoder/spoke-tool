package main

import (
	"fmt"
	"strings"
	"unicode"
)

// ============================================================================
// Basic String Functions
// ============================================================================

// Length returns the length of a string
func Length(s string) int {
	return len(s)
}

// CharAt returns the character at index i
func CharAt(s string, i int) byte {
	if i < 0 || i >= len(s) {
		return 0
	}
	return s[i]
}

// Concat concatenates two strings
func Concat(a, b string) string {
	return a + b
}

// Repeat returns a string repeated n times
func Repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// ============================================================================
// Case Conversion
// ============================================================================

// ToUpper converts a string to uppercase
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToLower converts a string to lowercase
func ToLower(s string) string {
	return strings.ToLower(s)
}

// Capitalize capitalizes the first letter of a string
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// TitleCase converts a string to title case
func TitleCase(s string) string {
	return strings.Title(s)
}

// SwapCase swaps the case of each letter
func SwapCase(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if unicode.IsUpper(r) {
			result[i] = unicode.ToLower(r)
		} else if unicode.IsLower(r) {
			result[i] = unicode.ToUpper(r)
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// ============================================================================
// Trimming
// ============================================================================

// Trim removes whitespace from both ends
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// TrimLeft removes whitespace from the left
func TrimLeft(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

// TrimRight removes whitespace from the right
func TrimRight(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// TrimPrefix removes a prefix if present
func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// TrimSuffix removes a suffix if present
func TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

// ============================================================================
// Searching
// ============================================================================

// Contains checks if a string contains a substring
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// ContainsAny checks if a string contains any of the characters
func ContainsAny(s, chars string) bool {
	return strings.ContainsAny(s, chars)
}

// Index returns the first index of a substring
func Index(s, substr string) int {
	return strings.Index(s, substr)
}

// LastIndex returns the last index of a substring
func LastIndex(s, substr string) int {
	return strings.LastIndex(s, substr)
}

// StartsWith checks if a string starts with a prefix
func StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// EndsWith checks if a string ends with a suffix
func EndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// Count counts occurrences of a substring
func Count(s, substr string) int {
	return strings.Count(s, substr)
}

// ============================================================================
// Splitting and Joining
// ============================================================================

// Split splits a string by a separator
func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

// SplitN splits a string into at most n parts
func SplitN(s, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
}

// SplitAfter splits after each separator
func SplitAfter(s, sep string) []string {
	return strings.SplitAfter(s, sep)
}

// Fields splits by whitespace
func Fields(s string) []string {
	return strings.Fields(s)
}

// Join concatenates strings with a separator
func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

// ============================================================================
// Replacement
// ============================================================================

// Replace replaces occurrences of old with new
func Replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

// ReplaceAll replaces all occurrences
func ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// Map applies a function to each character
func Map(s string, mapping func(rune) rune) string {
	return strings.Map(mapping, s)
}

// ============================================================================
// Comparison
// ============================================================================

// Compare compares two strings lexicographically
func Compare(a, b string) int {
	return strings.Compare(a, b)
}

// EqualFold checks equality ignoring case
func EqualFold(a, b string) bool {
	return strings.EqualFold(a, b)
}

// ============================================================================
// Padding
// ============================================================================

// PadLeft pads a string on the left
func PadLeft(s string, length int, pad byte) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(pad), length-len(s))
	return padding + s
}

// PadRight pads a string on the right
func PadRight(s string, length int, pad byte) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(pad), length-len(s))
	return s + padding
}

// PadBoth pads a string on both sides
func PadBoth(s string, length int, pad byte) string {
	if len(s) >= length {
		return s
	}
	totalPad := length - len(s)
	leftPad := totalPad / 2
	rightPad := totalPad - leftPad
	return strings.Repeat(string(pad), leftPad) + s + strings.Repeat(string(pad), rightPad)
}

// ============================================================================
// Substring Extraction
// ============================================================================

// Substring returns a substring from start to end
func Substring(s string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	if start >= end {
		return ""
	}
	return s[start:end]
}

// Left returns the first n characters
func Left(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if n >= len(s) {
		return s
	}
	return s[:n]
}

// Right returns the last n characters
func Right(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if n >= len(s) {
		return s
	}
	return s[len(s)-n:]
}

// ============================================================================
// Validation
// ============================================================================

// IsEmpty checks if a string is empty
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsBlank checks if a string is empty or whitespace
func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNumeric checks if a string contains only digits
func IsNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// IsAlpha checks if a string contains only letters
func IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsAlphanumeric checks if a string contains only letters and digits
func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsLowerCase checks if all letters are lowercase
func IsLowerCase(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsLower(r) {
			return false
		}
	}
	return true
}

// IsUpperCase checks if all letters are uppercase
func IsUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// ============================================================================
// Transformations
// ============================================================================

// Reverse reverses a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if i == 0 {
			words[i] = strings.ToLower(word)
		} else {
			words[i] = strings.Title(word)
		}
	}
	return strings.Join(words, "")
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' || r == '-' {
			return '_'
		}
		return r
	}, s)
}

// ToKebabCase converts a string to kebab-case
func ToKebabCase(s string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' || r == '_' {
			return '-'
		}
		return r
	}, s)
}

// ============================================================================
// Utility Functions
// ============================================================================

// Truncate truncates a string to a maximum length
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// Center centers a string in a field of given width
func Center(s string, width int, pad byte) string {
	if len(s) >= width {
		return s
	}
	totalPad := width - len(s)
	leftPad := totalPad / 2
	rightPad := totalPad - leftPad
	return strings.Repeat(string(pad), leftPad) + s + strings.Repeat(string(pad), rightPad)
}

// WordCount returns the number of words in a string
func WordCount(s string) int {
	return len(strings.Fields(s))
}

// Unique returns a string with duplicate characters removed
func Unique(s string) string {
	seen := make(map[rune]bool)
	result := make([]rune, 0)
	for _, r := range s {
		if !seen[r] {
			seen[r] = true
			result = append(result, r)
		}
	}
	return string(result)
}

// ============================================================================
// Main Function
// ============================================================================

func main() {
	// Basic operations
	fmt.Println("=== Basic Operations ===")
	str := "Hello, World!"
	fmt.Printf("Original: %q\n", str)
	fmt.Printf("Length: %d\n", Length(str))
	fmt.Printf("CharAt(7): %c\n", CharAt(str, 7))
	fmt.Printf("Concat: %q\n", Concat("Hello", " World"))
	fmt.Printf("Repeat: %q\n", Repeat("Ha", 3))

	// Case conversion
	fmt.Println("\n=== Case Conversion ===")
	str2 := "Hello, World!"
	fmt.Printf("ToUpper: %q\n", ToUpper(str2))
	fmt.Printf("ToLower: %q\n", ToLower(str2))
	fmt.Printf("Capitalize: %q\n", Capitalize(str2))
	fmt.Printf("TitleCase: %q\n", TitleCase(str2))
	fmt.Printf("SwapCase: %q\n", SwapCase(str2))

	// Trimming
	fmt.Println("\n=== Trimming ===")
	str3 := "  \tHello, World!  \n"
	fmt.Printf("Original: %q\n", str3)
	fmt.Printf("Trim: %q\n", Trim(str3))
	fmt.Printf("TrimLeft: %q\n", TrimLeft(str3))
	fmt.Printf("TrimRight: %q\n", TrimRight(str3))

	// Searching
	fmt.Println("\n=== Searching ===")
	str4 := "Hello, World! Hello, Go!"
	fmt.Printf("Contains 'World': %v\n", Contains(str4, "World"))
	fmt.Printf("ContainsAny 'xyz': %v\n", ContainsAny(str4, "xyz"))
	fmt.Printf("Index of 'Hello': %d\n", Index(str4, "Hello"))
	fmt.Printf("LastIndex of 'Hello': %d\n", LastIndex(str4, "Hello"))
	fmt.Printf("StartsWith 'Hello': %v\n", StartsWith(str4, "Hello"))
	fmt.Printf("EndsWith 'Go!': %v\n", EndsWith(str4, "Go!"))
	fmt.Printf("Count of 'Hello': %d\n", Count(str4, "Hello"))

	// Splitting and joining
	fmt.Println("\n=== Splitting and Joining ===")
	str5 := "apple,banana,orange,grape"
	parts := Split(str5, ",")
	fmt.Printf("Split: %v\n", parts)
	fmt.Printf("Join: %q\n", Join(parts, " | "))
	fmt.Printf("Fields: %v\n", Fields("  one   two   three  "))

	// Replacement
	fmt.Println("\n=== Replacement ===")
	str6 := "Hello, Hello, Hello!"
	fmt.Printf("Replace (2 times): %q\n", Replace(str6, "Hello", "Hi", 2))
	fmt.Printf("ReplaceAll: %q\n", ReplaceAll(str6, "Hello", "Hi"))

	// Comparison
	fmt.Println("\n=== Comparison ===")
	fmt.Printf("Compare 'a' vs 'b': %d\n", Compare("a", "b"))
	fmt.Printf("EqualFold 'Hello' vs 'HELLO': %v\n", EqualFold("Hello", "HELLO"))

	// Padding
	fmt.Println("\n=== Padding ===")
	str7 := "Go"
	fmt.Printf("PadLeft: %q\n", PadLeft(str7, 5, '*'))
	fmt.Printf("PadRight: %q\n", PadRight(str7, 5, '*'))
	fmt.Printf("PadBoth: %q\n", PadBoth(str7, 5, '*'))

	// Substring
	fmt.Println("\n=== Substring ===")
	str8 := "Hello, World!"
	fmt.Printf("Substring(0,5): %q\n", Substring(str8, 0, 5))
	fmt.Printf("Left(5): %q\n", Left(str8, 5))
	fmt.Printf("Right(6): %q\n", Right(str8, 6))

	// Validation
	fmt.Println("\n=== Validation ===")
	fmt.Printf("IsEmpty '': %v\n", IsEmpty(""))
	fmt.Printf("IsBlank '   ': %v\n", IsBlank("   "))
	fmt.Printf("IsNumeric '123': %v\n", IsNumeric("123"))
	fmt.Printf("IsAlpha 'abc': %v\n", IsAlpha("abc"))
	fmt.Printf("IsAlphanumeric 'abc123': %v\n", IsAlphanumeric("abc123"))
	fmt.Printf("IsLowerCase 'hello': %v\n", IsLowerCase("hello"))
	fmt.Printf("IsUpperCase 'HELLO': %v\n", IsUpperCase("HELLO"))

	// Transformations
	fmt.Println("\n=== Transformations ===")
	str9 := "Hello, World!"
	fmt.Printf("Reverse: %q\n", Reverse(str9))
	fmt.Printf("ToCamelCase: %q\n", ToCamelCase("hello world from go"))
	fmt.Printf("ToSnakeCase: %q\n", ToSnakeCase("hello-world-from-go"))
	fmt.Printf("ToKebabCase: %q\n", ToKebabCase("hello_world_from_go"))

	// Utility functions
	fmt.Println("\n=== Utility ===")
	str10 := "This is a very long string that needs truncation"
	fmt.Printf("Truncate: %q\n", Truncate(str10, 20))
	fmt.Printf("Center: %q\n", Center("Go", 11, '*'))
	fmt.Printf("WordCount: %d\n", WordCount("The quick brown fox jumps over the lazy dog"))
	fmt.Printf("Unique: %q\n", Unique("hello world"))
}
