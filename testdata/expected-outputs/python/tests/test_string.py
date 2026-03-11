"""
String Operations Tests

Tests for string manipulation functions including case conversion,
trimming, splitting, joining, searching, Unicode handling,
and property-based tests.
"""

import pytest
import string
import random
from hypothesis import given, strategies as st, assume
from string_utils import (
    # Basic operations
    reverse,
    length,
    char_at,
    char_code_at,
    
    # Case conversion
    to_upper,
    to_lower,
    to_title_case,
    to_camel_case,
    to_pascal_case,
    to_snake_case,
    to_kebab_case,
    capitalize,
    
    # Trimming
    trim,
    trim_start,
    trim_end,
    
    # Splitting and joining
    split,
    join,
    chars,
    words,
    lines,
    
    # Searching
    contains,
    index_of,
    last_index_of,
    starts_with,
    ends_with,
    count_occurrences,
    
    # Replacement
    replace,
    replace_all,
    remove,
    remove_all,
    
    # Substring extraction
    substring,
    slice,
    
    # Padding
    pad_start,
    pad_end,
    
    # Repetition
    repeat,
    
    # Validation
    is_empty,
    is_blank,
    is_numeric,
    is_alpha,
    is_alphanumeric,
    is_lowercase,
    is_uppercase,
    is_title_case,
    is_palindrome,
    is_anagram,
    
    # Utilities
    truncate,
    ellipsis,
    slugify,
    pluralize,
    singularize,
    camel_case,
    snake_case,
    kebab_case,
    
    # Comparison
    compare,
    compare_ignore_case,
    levenshtein_distance,
    hamming_distance,
    longest_common_substring,
    longest_common_prefix,
    longest_common_suffix,
    
    # Encoding
    to_base64,
    from_base64,
    to_hex,
    from_hex,
    to_binary,
    from_binary,
    rot13,
    caesar_cipher,
    
    # Unicode
    normalize_unicode,
    remove_accents,
    is_unicode,
    count_graphemes,
    reverse_graphemes,
    
    # Formatting
    format_template,
    interpolate,
    mask,
    obfuscate,
    highlight,
    
    # Random generation
    random_string,
    random_alpha,
    random_numeric,
    random_alphanumeric,
    
    # Advanced
    acronym,
    initials,
    word_count,
    sentence_count,
    reading_time,
    complexity_score
)

# ============================================================================
# Basic Operations Tests
# ============================================================================

class TestBasicOperations:
    """Tests for basic string operations."""

    def test_reverse(self):
        """Test string reversal."""
        assert reverse("hello") == "olleh"
        assert reverse("") == ""
        assert reverse("a") == "a"
        assert reverse("racecar") == "racecar"
        assert reverse("hello world") == "dlrow olleh"
        assert reverse("12345") == "54321"
        assert reverse("!@#$%") == "%$#@!"

    def test_reverse_unicode(self):
        """Test reversal with Unicode characters."""
        assert reverse("café") == "éfac"
        assert reverse("résumé") == "émusér"
        assert reverse("Hello 世界") == "界世 olleH"
        assert reverse("🚀✨🌟") == "🌟✨🚀"

    def test_length(self):
        """Test string length."""
        assert length("hello") == 5
        assert length("") == 0
        assert length("hello world") == 11
        assert length("café") == 4
        assert length("Hello 世界") == 8

    def test_char_at(self):
        """Test character at index."""
        assert char_at("hello", 0) == "h"
        assert char_at("hello", 2) == "l"
        assert char_at("hello", 4) == "o"
        assert char_at("hello", -1) == ""
        assert char_at("hello", 5) == ""
        assert char_at("café", 3) == "é"

    def test_char_code_at(self):
        """Test character code at index."""
        assert char_code_at("A", 0) == 65
        assert char_code_at("a", 0) == 97
        assert char_code_at("0", 0) == 48
        assert char_code_at("é", 0) == 233

# ============================================================================
# Case Conversion Tests
# ============================================================================

class TestCaseConversion:
    """Tests for case conversion functions."""

    def test_to_upper(self):
        """Test uppercase conversion."""
        assert to_upper("hello") == "HELLO"
        assert to_upper("Hello World") == "HELLO WORLD"
        assert to_upper("123abc") == "123ABC"
        assert to_upper("") == ""
        assert to_upper("café") == "CAFÉ"
        assert to_upper("straße") == "STRASSE"

    def test_to_lower(self):
        """Test lowercase conversion."""
        assert to_lower("HELLO") == "hello"
        assert to_lower("Hello World") == "hello world"
        assert to_lower("123ABC") == "123abc"
        assert to_lower("") == ""
        assert to_lower("CAFÉ") == "café"

    def test_to_title_case(self):
        """Test title case conversion."""
        assert to_title_case("hello world") == "Hello World"
        assert to_title_case("the quick brown fox") == "The Quick Brown Fox"
        assert to_title_case("hELLO wORLD") == "Hello World"
        assert to_title_case("") == ""

    def test_capitalize(self):
        """Test capitalize first letter."""
        assert capitalize("hello") == "Hello"
        assert capitalize("hello world") == "Hello world"
        assert capitalize("") == ""
        assert capitalize("HELLO") == "Hello"

    def test_to_camel_case(self):
        """Test camelCase conversion."""
        assert to_camel_case("hello world") == "helloWorld"
        assert to_camel_case("hello-world") == "helloWorld"
        assert to_camel_case("hello_world") == "helloWorld"
        assert to_camel_case("HelloWorld") == "helloWorld"
        assert to_camel_case("") == ""

    def test_to_pascal_case(self):
        """Test PascalCase conversion."""
        assert to_pascal_case("hello world") == "HelloWorld"
        assert to_pascal_case("hello-world") == "HelloWorld"
        assert to_pascal_case("hello_world") == "HelloWorld"
        assert to_pascal_case("helloWorld") == "HelloWorld"
        assert to_pascal_case("") == ""

    def test_to_snake_case(self):
        """Test snake_case conversion."""
        assert to_snake_case("helloWorld") == "hello_world"
        assert to_snake_case("HelloWorld") == "hello_world"
        assert to_snake_case("hello-world") == "hello_world"
        assert to_snake_case("hello world") == "hello_world"
        assert to_snake_case("") == ""

    def test_to_kebab_case(self):
        """Test kebab-case conversion."""
        assert to_kebab_case("helloWorld") == "hello-world"
        assert to_kebab_case("HelloWorld") == "hello-world"
        assert to_kebab_case("hello_world") == "hello-world"
        assert to_kebab_case("hello world") == "hello-world"
        assert to_kebab_case("") == ""

# ============================================================================
# Trimming Tests
# ============================================================================

class TestTrimming:
    """Tests for string trimming functions."""

    def test_trim(self):
        """Test trimming whitespace from both ends."""
        assert trim("  hello  ") == "hello"
        assert trim("\t\n hello \n\t") == "hello"
        assert trim("") == ""
        assert trim("   ") == ""
        assert trim("hello") == "hello"
        assert trim("  hello world  ") == "hello world"

    def test_trim_start(self):
        """Test trimming leading whitespace."""
        assert trim_start("  hello") == "hello"
        assert trim_start("\t\nhello") == "hello"
        assert trim_start("  hello  ") == "hello  "
        assert trim_start("hello") == "hello"

    def test_trim_end(self):
        """Test trimming trailing whitespace."""
        assert trim_end("hello  ") == "hello"
        assert trim_end("hello\n\t") == "hello"
        assert trim_end("  hello  ") == "  hello"
        assert trim_end("hello") == "hello"

# ============================================================================
# Splitting and Joining Tests
# ============================================================================

class TestSplittingJoining:
    """Tests for splitting and joining strings."""

    def test_split(self):
        """Test splitting strings."""
        assert split("a,b,c", ",") == ["a", "b", "c"]
        assert split("hello world", " ") == ["hello", "world"]
        assert split("a,,b,c", ",") == ["a", "", "b", "c"]
        assert split(",a,b", ",") == ["", "a", "b"]
        assert split("a,b,", ",") == ["a", "b", ""]
        assert split("hello", ",") == ["hello"]

    def test_split_with_limit(self):
        """Test splitting with limit."""
        assert split("a,b,c,d", ",", 2) == ["a", "b"]
        assert split("a,b,c,d", ",", 0) == []

    def test_split_regex(self):
        """Test splitting with regex pattern."""
        assert split("hello123world", r"\d+") == ["hello", "world"]
        assert split("a,b;c", r"[,;]") == ["a", "b", "c"]

    def test_join(self):
        """Test joining strings."""
        assert join(["a", "b", "c"], ",") == "a,b,c"
        assert join(["hello", "world"], " ") == "hello world"
        assert join([], ",") == ""
        assert join(["hello"], ",") == "hello"
        assert join(["a", "", "c"], ",") == "a,,c"

    def test_chars(self):
        """Test splitting into characters."""
        assert chars("hello") == ["h", "e", "l", "l", "o"]
        assert chars("") == []
        assert chars("café") == ["c", "a", "f", "é"]

    def test_words(self):
        """Test splitting into words."""
        assert words("hello world") == ["hello", "world"]
        assert words("the quick brown fox") == ["the", "quick", "brown", "fox"]
        assert words("hello, world!") == ["hello", "world"]
        assert words("hello   world") == ["hello", "world"]
        assert words("") == []

    def test_lines(self):
        """Test splitting into lines."""
        text = "line1\nline2\nline3"
        assert lines(text) == ["line1", "line2", "line3"]
        assert lines("") == []
        assert lines("single line") == ["single line"]

# ============================================================================
# Searching Tests
# ============================================================================

class TestSearching:
    """Tests for string searching functions."""

    def test_contains(self):
        """Test substring containment."""
        assert contains("hello world", "world") is True
        assert contains("hello world", "xyz") is False
        assert contains("hello", "") is True
        assert contains("Hello World", "hello", case_sensitive=False) is True

    def test_index_of(self):
        """Test finding first occurrence."""
        assert index_of("hello world", "world") == 6
        assert index_of("hello world", "o") == 4
        assert index_of("hello world", "xyz") == -1
        assert index_of("hello", "") == 0
        assert index_of("hello world", "o", 5) == 7

    def test_last_index_of(self):
        """Test finding last occurrence."""
        assert last_index_of("hello world hello", "hello") == 12
        assert last_index_of("hello world", "o") == 7
        assert last_index_of("hello world", "xyz") == -1

    def test_starts_with(self):
        """Test prefix checking."""
        assert starts_with("hello world", "hello") is True
        assert starts_with("hello world", "world") is False
        assert starts_with("hello", "") is True
        assert starts_with("hello world", "world", 6) is True

    def test_ends_with(self):
        """Test suffix checking."""
        assert ends_with("hello world", "world") is True
        assert ends_with("hello world", "hello") is False
        assert ends_with("hello", "") is True
        assert ends_with("hello world", "hello", 5) is True

    def test_count_occurrences(self):
        """Test counting occurrences."""
        assert count_occurrences("hello world hello", "hello") == 2
        assert count_occurrences("hello world", "o") == 2
        assert count_occurrences("aaaaa", "aa") == 2
        assert count_occurrences("hello", "xyz") == 0

# ============================================================================
# Replacement Tests
# ============================================================================

class TestReplacement:
    """Tests for string replacement functions."""

    def test_replace(self):
        """Test replacing first occurrence."""
        assert replace("hello world hello", "hello", "hi") == "hi world hello"
        assert replace("hello world", "xyz", "abc") == "hello world"
        assert replace("hello123world", r"\d+", "!") == "hello!world"

    def test_replace_all(self):
        """Test replacing all occurrences."""
        assert replace_all("hello world hello", "hello", "hi") == "hi world hi"
        assert replace_all("aaa", "a", "b") == "bbb"
        assert replace_all("hello world", "xyz", "abc") == "hello world"

    def test_remove(self):
        """Test removing first occurrence."""
        assert remove("hello world hello", "hello") == " world hello"
        assert remove("aaa", "a") == "aa"

    def test_remove_all(self):
        """Test removing all occurrences."""
        assert remove_all("hello world hello", "hello") == " world "
        assert remove_all("aaa", "a") == ""

# ============================================================================
# Substring Tests
# ============================================================================

class TestSubstring:
    """Tests for substring extraction."""

    def test_substring(self):
        """Test substring extraction."""
        assert substring("hello world", 0, 5) == "hello"
        assert substring("hello world", 6, 11) == "world"
        assert substring("hello world", 6) == "world"
        assert substring("hello world", -5, 5) == "hello"

    def test_slice(self):
        """Test slice extraction."""
        assert slice("hello world", 0, 5) == "hello"
        assert slice("hello world", 6, 11) == "world"
        assert slice("hello world", -5) == "world"
        assert slice("hello world", -5, -1) == "worl"

# ============================================================================
# Padding Tests
# ============================================================================

class TestPadding:
    """Tests for string padding functions."""

    def test_pad_start(self):
        """Test padding at start."""
        assert pad_start("5", 3, "0") == "005"
        assert pad_start("abc", 5, "-") == "--abc"
        assert pad_start("hello", 10) == "     hello"
        assert pad_start("hello", 3) == "hello"

    def test_pad_end(self):
        """Test padding at end."""
        assert pad_end("5", 3, "0") == "500"
        assert pad_end("abc", 5, "-") == "abc--"
        assert pad_end("hello", 10) == "hello     "
        assert pad_end("hello", 3) == "hello"

# ============================================================================
# Repetition Tests
# ============================================================================

class TestRepetition:
    """Tests for string repetition."""

    def test_repeat(self):
        """Test string repetition."""
        assert repeat("ha", 3) == "hahaha"
        assert repeat("abc", 2) == "abcabc"
        assert repeat("hello", 0) == ""
        assert repeat("", 5) == ""
        
        with pytest.raises(ValueError):
            repeat("hello", -1)

# ============================================================================
# Validation Tests
# ============================================================================

class TestValidation:
    """Tests for string validation functions."""

    def test_is_empty(self):
        """Test empty string detection."""
        assert is_empty("") is True
        assert is_empty("   ") is False
        assert is_empty("hello") is False

    def test_is_blank(self):
        """Test blank string detection."""
        assert is_blank("") is True
        assert is_blank("   ") is True
        assert is_blank("\t\n") is True
        assert is_blank("hello") is False

    def test_is_numeric(self):
        """Test numeric string detection."""
        assert is_numeric("123") is True
        assert is_numeric("-123") is True
        assert is_numeric("123.456") is True
        assert is_numeric("abc") is False
        assert is_numeric("123abc") is False

    def test_is_alpha(self):
        """Test alphabetic string detection."""
        assert is_alpha("abc") is True
        assert is_alpha("ABC") is True
        assert is_alpha("abc123") is False
        assert is_alpha("") is False

    def test_is_alphanumeric(self):
        """Test alphanumeric string detection."""
        assert is_alphanumeric("abc123") is True
        assert is_alphanumeric("ABC123") is True
        assert is_alphanumeric("abc-123") is False
        assert is_alphanumeric("") is False

    def test_is_lowercase(self):
        """Test lowercase string detection."""
        assert is_lowercase("hello") is True
        assert is_lowercase("Hello") is False
        assert is_lowercase("123") is True  # Numbers don't affect case

    def test_is_uppercase(self):
        """Test uppercase string detection."""
        assert is_uppercase("HELLO") is True
        assert is_uppercase("Hello") is False
        assert is_uppercase("123") is True

    def test_is_palindrome(self):
        """Test palindrome detection."""
        assert is_palindrome("racecar") is True
        assert is_palindrome("hello") is False
        assert is_palindrome("A man a plan a canal Panama") is True
        assert is_palindrome("") is True

    def test_is_anagram(self):
        """Test anagram detection."""
        assert is_anagram("listen", "silent") is True
        assert is_anagram("hello", "world") is False
        assert is_anagram("debit card", "bad credit") is True
        assert is_anagram("", "") is True

# ============================================================================
# Utility Tests
# ============================================================================

class TestUtilities:
    """Tests for string utility functions."""

    def test_truncate(self):
        """Test string truncation."""
        assert truncate("This is a long string", 10) == "This is..."
        assert truncate("Hello world", 5) == "Hello..."
        assert truncate("Hello", 10) == "Hello"
        assert truncate("", 5) == ""

    def test_ellipsis(self):
        """Test ellipsis formatting."""
        assert ellipsis("This is a long string", 10) == "This is..."
        assert ellipsis("Hello", 10) == "Hello"
        assert ellipsis("This is a long string", 10, "middle") == "This...ring"
        assert ellipsis("This is a long string", 10, "start") == "...g string"

    def test_slugify(self):
        """Test slug generation."""
        assert slugify("Hello World") == "hello-world"
        assert slugify("This is a test") == "this-is-a-test"
        assert slugify("Hello, World!") == "hello-world"
        assert slugify("Café & Restaurant") == "cafe-restaurant"
        assert slugify("") == ""

    def test_pluralize(self):
        """Test pluralization."""
        assert pluralize("cat", 2) == "cats"
        assert pluralize("cat", 1) == "cat"
        assert pluralize("cat", 0) == "cats"
        assert pluralize("child") == "children"
        assert pluralize("person") == "people"
        assert pluralize("mouse") == "mice"

    def test_singularize(self):
        """Test singularization."""
        assert singularize("cats") == "cat"
        assert singularize("children") == "child"
        assert singularize("people") == "person"
        assert singularize("mice") == "mouse"

# ============================================================================
# Comparison Tests
# ============================================================================

class TestComparison:
    """Tests for string comparison functions."""

    def test_compare(self):
        """Test string comparison."""
        assert compare("a", "b") == -1
        assert compare("b", "a") == 1
        assert compare("a", "a") == 0
        assert compare("A", "a") == 1

    def test_compare_ignore_case(self):
        """Test case-insensitive comparison."""
        assert compare_ignore_case("hello", "HELLO") == 0
        assert compare_ignore_case("Hello", "hELLO") == 0
        assert compare_ignore_case("hello", "world") == -1

    def test_levenshtein_distance(self):
        """Test Levenshtein distance."""
        assert levenshtein_distance("kitten", "sitting") == 3
        assert levenshtein_distance("hello", "hello") == 0
        assert levenshtein_distance("", "hello") == 5
        assert levenshtein_distance("hello", "") == 5

    def test_hamming_distance(self):
        """Test Hamming distance."""
        assert hamming_distance("hello", "hallo") == 1
        assert hamming_distance("hello", "hello") == 0
        
        with pytest.raises(ValueError):
            hamming_distance("hello", "world")  # Different lengths

    def test_longest_common_substring(self):
        """Test longest common substring."""
        assert longest_common_substring("abcdef", "zbcdf") == "bcd"
        assert longest_common_substring("hello", "world") == ""
        assert longest_common_substring("", "hello") == ""

    def test_longest_common_prefix(self):
        """Test longest common prefix."""
        assert longest_common_prefix("hello", "helicopter") == "hel"
        assert longest_common_prefix("hello", "world") == ""
        assert longest_common_prefix("", "hello") == ""

    def test_longest_common_suffix(self):
        """Test longest common suffix."""
        assert longest_common_suffix("running", "swimming") == "ing"
        assert longest_common_suffix("hello", "world") == ""
        assert longest_common_suffix("hello", "") == ""

# ============================================================================
# Encoding Tests
# ============================================================================

class TestEncoding:
    """Tests for string encoding functions."""

    def test_base64(self):
        """Test Base64 encoding/decoding."""
        original = "hello world"
        encoded = to_base64(original)
        assert encoded == "aGVsbG8gd29ybGQ="
        assert from_base64(encoded) == original
        
        assert to_base64("") == ""
        assert from_base64("") == ""

    def test_hex(self):
        """Test hex encoding/decoding."""
        original = "hello"
        encoded = to_hex(original)
        assert encoded == "68656c6c6f"
        assert from_hex(encoded) == original
        
        with pytest.raises(ValueError):
            from_hex("xyz")

    def test_binary(self):
        """Test binary encoding/decoding."""
        assert to_binary("A") == "01000001"
        assert from_binary("01000001") == "A"
        assert to_binary("AB") == "0100000101000010"
        assert from_binary("0100000101000010") == "AB"

    def test_rot13(self):
        """Test ROT13 cipher."""
        assert rot13("hello") == "uryyb"
        assert rot13(rot13("hello")) == "hello"
        assert rot13("") == ""

    def test_caesar_cipher(self):
        """Test Caesar cipher."""
        assert caesar_cipher("hello", 3) == "khoor"
        assert caesar_cipher("hello", -3) == "ebiil"
        assert caesar_cipher("Hello, World!", 5) == "Mjqqt, Btwqi!"

# ============================================================================
# Unicode Tests
# ============================================================================

class TestUnicode:
    """Tests for Unicode handling."""

    def test_normalize_unicode(self):
        """Test Unicode normalization."""
        assert normalize_unicode("café") == "café"
        assert normalize_unicode("cafe\u0301") == "café"

    def test_remove_accents(self):
        """Test accent removal."""
        assert remove_accents("café") == "cafe"
        assert remove_accents("résumé") == "resume"
        assert remove_accents("über") == "uber"
        assert remove_accents("Français") == "Francais"

    def test_count_graphemes(self):
        """Test grapheme cluster counting."""
        assert count_graphemes("hello") == 5
        assert count_graphemes("café") == 4
        assert count_graphemes("👨‍👩‍👧‍👦") == 1  # Family emoji
        assert count_graphemes("é") == 1  # e + accent

    def test_reverse_graphemes(self):
        """Test grapheme-aware reversal."""
        assert reverse_graphemes("hello") == "olleh"
        assert reverse_graphemes("café") == "éfac"
        assert reverse_graphemes("👨‍👩‍👧‍👦") == "👨‍👩‍👧‍👦"  # Should stay the same

# ============================================================================
# Formatting Tests
# ============================================================================

class TestFormatting:
    """Tests for string formatting functions."""

    def test_format_template(self):
        """Test template formatting."""
        template = "Hello {{name}}!"
        assert format_template(template, name="World") == "Hello World!"
        assert format_template("{{greet}} {{name}}", greet="Hi", name="John") == "Hi John"

    def test_interpolate(self):
        """Test string interpolation."""
        assert interpolate("Hello {0}!", ["World"]) == "Hello World!"
        assert interpolate("{0} {1}", ["Hello", "World"]) == "Hello World"
        assert interpolate("{0} {0}", ["echo"]) == "echo echo"

    def test_mask(self):
        """Test string masking."""
        assert mask("1234-5678-9012-3456", 4) == "************3456"
        assert mask("password123", 0) == "***********"
        assert mask("secret", 2, "#") == "####et"
        assert mask("1234567890", 4, "*", "start") == "****567890"

    def test_obfuscate(self):
        """Test email/phone obfuscation."""
        assert obfuscate("user@example.com") == "u**r@e*****e.com"
        assert obfuscate("123-456-7890") == "***-***-7890"
        assert obfuscate("a@b.c") == "*@b.c"

    def test_highlight(self):
        """Test highlighting matches."""
        assert highlight("hello world hello", "hello") == "**hello** world **hello**"
        assert highlight("Hello World", "hello", case_sensitive=False) == "**Hello** World"

# ============================================================================
# Random Generation Tests
# ============================================================================

class TestRandomGeneration:
    """Tests for random string generation."""

    def test_random_string(self):
        """Test random string generation."""
        s = random_string(10)
        assert len(s) == 10
        assert isinstance(s, str)

    def test_random_alpha(self):
        """Test random alphabetic string."""
        s = random_alpha(10)
        assert len(s) == 10
        assert all(c.isalpha() for c in s)

    def test_random_numeric(self):
        """Test random numeric string."""
        s = random_numeric(10)
        assert len(s) == 10
        assert all(c.isdigit() for c in s)

    def test_random_alphanumeric(self):
        """Test random alphanumeric string."""
        s = random_alphanumeric(10)
        assert len(s) == 10
        assert all(c.isalnum() for c in s)

# ============================================================================
# Advanced Tests
# ============================================================================

class TestAdvanced:
    """Tests for advanced string functions."""

    def test_acronym(self):
        """Test acronym generation."""
        assert acronym("World Health Organization") == "WHO"
        assert acronym("as soon as possible") == "ASAP"
        assert acronym("") == ""

    def test_initials(self):
        """Test initials extraction."""
        assert initials("John Fitzgerald Kennedy") == "JFK"
        assert initials("john f kennedy") == "JFK"
        assert initials("") == ""

    def test_word_count(self):
        """Test word counting."""
        assert word_count("The quick brown fox") == 4
        assert word_count("") == 0
        assert word_count("Hello, world!") == 2

    def test_sentence_count(self):
        """Test sentence counting."""
        text = "Hello world. This is a test. How are you?"
        assert sentence_count(text) == 3
        assert sentence_count("") == 0

    def test_reading_time(self):
        """Test reading time estimation."""
        text = "word " * 200  # 200 words
        minutes, seconds = reading_time(text)
        assert minutes == 1  # Assuming 200 WPM
        assert 0 <= seconds < 60

    def test_complexity_score(self):
        """Test string complexity scoring."""
        assert complexity_score("hello") < complexity_score("Hello123!")
        assert complexity_score("") == 0

# ============================================================================
# Property-Based Tests
# ============================================================================

class TestProperties:
    """Property-based tests using Hypothesis."""

    @given(st.text())
    def test_reverse_reverse(self, s):
        """Test that reversing twice returns original."""
        assume(len(s) < 1000)  # Avoid extremely long strings
        assert reverse(reverse(s)) == s

    @given(st.text())
    def test_upper_lower_roundtrip(self, s):
        """Test that upper then lower returns original."""
        assert to_lower(to_upper(s)) == to_lower(s)

    @given(st.text())
    def test_trim_doesnt_add(self, s):
        """Test that trim doesn't add characters."""
        trimmed = trim(s)
        assert len(trimmed) <= len(s)

    @given(st.lists(st.text(max_size=5), min_size=1, max_size=5))
    def test_split_join_roundtrip(self, parts):
        """Test that split then join returns original."""
        s = "".join(parts)
        assert join(chars(s), "") == s

    @given(st.text())
    def test_palindrome_property(self, s):
        """Test palindrome property."""
        assert is_palindrome(s + reverse(s)) is True

    @given(st.text(), st.text())
    def test_levenshtein_triangle_inequality(self, a, b):
        """Test Levenshtein distance triangle inequality."""
        assume(len(a) < 50 and len(b) < 50)
        c = a + b
        dist_ab = levenshtein_distance(a, b)
        dist_ac = levenshtein_distance(a, c)
        dist_bc = levenshtein_distance(b, c)
        assert dist_ac <= dist_ab + dist_bc

# ============================================================================
# Edge Cases and Error Handling
# ============================================================================

class TestEdgeCases:
    """Tests for edge cases and error handling."""

    def test_null_input(self):
        """Test handling of None input."""
        with pytest.raises(TypeError):
            reverse(None)

    def test_non_string_input(self):
        """Test handling of non-string input."""
        with pytest.raises(TypeError):
            reverse(123)

    def test_very_long_string(self):
        """Test handling of very long strings."""
        long_string = "a" * 1_000_000
        result = reverse(long_string)
        assert len(result) == 1_000_000
        assert result[0] == "a"  # Should be reversed

    def test_string_with_null_char(self):
        """Test handling of strings with null characters."""
        s = "hello\x00world"
        assert contains(s, "\x00") is True
        assert len(s) == 11

    def test_string_with_control_chars(self):
        """Test handling of control characters."""
        s = "hello\x1b[31mworld\x1b[0m"
        assert len(s) == 22
        assert reverse(s) is not None

# ============================================================================
# Performance Tests
# ============================================================================

class TestPerformance:
    """Performance tests for string operations."""

    @pytest.mark.benchmark
    def test_reverse_performance(self, benchmark):
        """Benchmark string reversal."""
        s = "a" * 10_000
        result = benchmark(reverse, s)
        assert len(result) == 10_000

    @pytest.mark.benchmark
    def test_split_performance(self, benchmark):
        """Benchmark string splitting."""
        s = "a," * 10_000
        result = benchmark(split, s, ",")
        assert len(result) == 10_001

    @pytest.mark.benchmark
    def test_levenshtein_performance(self, benchmark):
        """Benchmark Levenshtein distance."""
        s1 = "a" * 100
        s2 = "b" * 100
        result = benchmark(levenshtein_distance, s1, s2)
        assert result == 100

# ============================================================================
# Fixtures
# ============================================================================

@pytest.fixture
def sample_strings():
    """Provide sample strings for testing."""
    return {
        "empty": "",
        "single": "a",
        "short": "hello",
        "sentence": "The quick brown fox jumps over the lazy dog.",
        "unicode": "café résumé 中文 русский",
        "emoji": "Hello 👋 World 🌍",
        "numbers": "12345",
        "mixed": "abc123!@#",
        "palindrome": "racecar",
        "multiline": "line1\nline2\nline3"
    }

@pytest.fixture
def long_string():
    """Provide a long string for performance testing."""
    return "a" * 100_000

@pytest.fixture
def words_list():
    """Provide a list of words for testing."""
    return ["hello", "world", "python", "testing", "string"]