"""
Simple String Module

Basic string manipulation functions demonstrating Python string operations,
validation, and text processing.
"""

import re
from typing import List


# ============================================================================
# Case Conversion
# ============================================================================

def to_upper(s: str) -> str:
    """
    Convert string to uppercase.
    
    Args:
        s: Input string
    
    Returns:
        Uppercase string
    
    Examples:
        >>> to_upper("hello")
        'HELLO'
    """
    return s.upper()


def to_lower(s: str) -> str:
    """
    Convert string to lowercase.
    
    Args:
        s: Input string
    
    Returns:
        Lowercase string
    """
    return s.lower()


def capitalize(s: str) -> str:
    """
    Capitalize first letter of string.
    
    Args:
        s: Input string
    
    Returns:
        String with first letter capitalized
    
    Examples:
        >>> capitalize("hello")
        'Hello'
    """
    if not s:
        return s
    return s[0].upper() + s[1:].lower()


def title_case(s: str) -> str:
    """
    Convert string to title case (first letter of each word capitalized).
    
    Args:
        s: Input string
    
    Returns:
        Title case string
    """
    return s.title()


def swap_case(s: str) -> str:
    """
    Swap case of each character (upper to lower and vice versa).
    
    Args:
        s: Input string
    
    Returns:
        String with swapped case
    """
    return s.swapcase()


# ============================================================================
# String Reversal
# ============================================================================

def reverse(s: str) -> str:
    """
    Reverse a string.
    
    Args:
        s: Input string
    
    Returns:
        Reversed string
    
    Examples:
        >>> reverse("hello")
        'olleh'
    """
    return s[::-1]


def reverse_words(s: str) -> str:
    """
    Reverse the order of words in a string.
    
    Args:
        s: Input string
    
    Returns:
        String with words reversed
    """
    words = s.split()
    return " ".join(reversed(words))


def reverse_word_order(s: str) -> str:
    """
    Alias for reverse_words.
    """
    return reverse_words(s)


# ============================================================================
# Trimming and Whitespace
# ============================================================================

def trim(s: str) -> str:
    """
    Remove whitespace from both ends of string.
    
    Args:
        s: Input string
    
    Returns:
        Trimmed string
    """
    return s.strip()


def trim_left(s: str) -> str:
    """
    Remove whitespace from left end of string.
    
    Args:
        s: Input string
    
    Returns:
        Left-trimmed string
    """
    return s.lstrip()


def trim_right(s: str) -> str:
    """
    Remove whitespace from right end of string.
    
    Args:
        s: Input string
    
    Returns:
        Right-trimmed string
    """
    return s.rstrip()


def remove_whitespace(s: str) -> str:
    """
    Remove all whitespace from string.
    
    Args:
        s: Input string
    
    Returns:
        String with all whitespace removed
    """
    return "".join(s.split())


def normalize_spaces(s: str) -> str:
    """
    Normalize multiple spaces to single space.
    
    Args:
        s: Input string
    
    Returns:
        String with normalized spaces
    """
    return " ".join(s.split())


# ============================================================================
# Searching and Counting
# ============================================================================

def contains(s: str, substring: str) -> bool:
    """
    Check if string contains substring.
    
    Args:
        s: String to search
        substring: Substring to find
    
    Returns:
        True if substring found, False otherwise
    """
    return substring in s


def count_occurrences(s: str, substring: str) -> int:
    """
    Count occurrences of substring in string.
    
    Args:
        s: String to search
        substring: Substring to count
    
    Returns:
        Number of occurrences
    """
    return s.count(substring)


def count_vowels(s: str) -> int:
    """
    Count vowels in string.
    
    Args:
        s: Input string
    
    Returns:
        Number of vowels (a, e, i, o, u)
    
    Examples:
        >>> count_vowels("hello")
        2
    """
    vowels = "aeiouAEIOU"
    return sum(1 for char in s if char in vowels)


def count_consonants(s: str) -> int:
    """
    Count consonants in string.
    
    Args:
        s: Input string
    
    Returns:
        Number of consonants
    """
    vowels = "aeiouAEIOU"
    return sum(1 for char in s if char.isalpha() and char not in vowels)


def count_words(s: str) -> int:
    """
    Count words in string.
    
    Args:
        s: Input string
    
    Returns:
        Number of words
    """
    return len(s.split())


def count_sentences(s: str) -> int:
    """
    Count sentences in string (based on ., !, ?).
    
    Args:
        s: Input string
    
    Returns:
        Number of sentences
    """
    import re
    return len(re.findall(r'[.!?]+', s))


# ============================================================================
# Validation
# ============================================================================

def is_empty(s: str) -> bool:
    """
    Check if string is empty.
    
    Args:
        s: Input string
    
    Returns:
        True if empty, False otherwise
    """
    return len(s) == 0


def is_blank(s: str) -> bool:
    """
    Check if string is empty or only whitespace.
    
    Args:
        s: Input string
    
    Returns:
        True if blank, False otherwise
    """
    return len(s.strip()) == 0


def is_palindrome(s: str) -> bool:
    """
    Check if string is a palindrome (reads same forwards and backwards).
    
    Args:
        s: Input string
    
    Returns:
        True if palindrome, False otherwise
    
    Examples:
        >>> is_palindrome("racecar")
        True
        >>> is_palindrome("hello")
        False
    """
    # Remove non-alphanumeric and convert to lowercase
    cleaned = "".join(char.lower() for char in s if char.isalnum())
    return cleaned == cleaned[::-1]


def is_alpha(s: str) -> bool:
    """
    Check if string contains only letters.
    
    Args:
        s: Input string
    
    Returns:
        True if alphabetic, False otherwise
    """
    return s.isalpha()


def is_digit(s: str) -> bool:
    """
    Check if string contains only digits.
    
    Args:
        s: Input string
    
    Returns:
        True if digits only, False otherwise
    """
    return s.isdigit()


def is_alnum(s: str) -> bool:
    """
    Check if string contains only letters and numbers.
    
    Args:
        s: Input string
    
    Returns:
        True if alphanumeric, False otherwise
    """
    return s.isalnum()


def is_lower(s: str) -> bool:
    """
    Check if string is all lowercase.
    
    Args:
        s: Input string
    
    Returns:
        True if lowercase, False otherwise
    """
    return s.islower()


def is_upper(s: str) -> bool:
    """
    Check if string is all uppercase.
    
    Args:
        s: Input string
    
    Returns:
        True if uppercase, False otherwise
    """
    return s.isupper()


def is_title(s: str) -> bool:
    """
    Check if string is in title case.
    
    Args:
        s: Input string
    
    Returns:
        True if title case, False otherwise
    """
    return s.istitle()


# ============================================================================
# Substring Operations
# ============================================================================

def left(s: str, n: int) -> str:
    """
    Get first n characters of string.
    
    Args:
        s: Input string
        n: Number of characters
    
    Returns:
        First n characters
    """
    return s[:n]


def right(s: str, n: int) -> str:
    """
    Get last n characters of string.
    
    Args:
        s: Input string
        n: Number of characters
    
    Returns:
        Last n characters
    """
    return s[-n:] if n > 0 else ""


def mid(s: str, start: int, length: int) -> str:
    """
    Get substring starting at position with given length.
    
    Args:
        s: Input string
        start: Starting position (0-indexed)
        length: Length of substring
    
    Returns:
        Substring
    """
    return s[start:start + length]


def before_first(s: str, delimiter: str) -> str:
    """
    Get substring before first occurrence of delimiter.
    
    Args:
        s: Input string
        delimiter: Delimiter string
    
    Returns:
        Substring before first delimiter
    """
    if delimiter not in s:
        return s
    return s.split(delimiter)[0]


def after_first(s: str, delimiter: str) -> str:
    """
    Get substring after first occurrence of delimiter.
    
    Args:
        s: Input string
        delimiter: Delimiter string
    
    Returns:
        Substring after first delimiter
    """
    if delimiter not in s:
        return ""
    return s.split(delimiter, 1)[1]


# ============================================================================
# Replacement and Transformation
# ============================================================================

def replace(s: str, old: str, new: str) -> str:
    """
    Replace all occurrences of old with new.
    
    Args:
        s: Input string
        old: Substring to replace
        new: Replacement string
    
    Returns:
        String with replacements
    """
    return s.replace(old, new)


def remove(s: str, substring: str) -> str:
    """
    Remove all occurrences of substring.
    
    Args:
        s: Input string
        substring: Substring to remove
    
    Returns:
        String with substring removed
    """
    return s.replace(substring, "")


def remove_duplicates(s: str) -> str:
    """
    Remove duplicate characters from string.
    
    Args:
        s: Input string
    
    Returns:
        String with duplicate characters removed
    """
    seen = set()
    result = []
    for char in s:
        if char not in seen:
            seen.add(char)
            result.append(char)
    return "".join(result)


def truncate(s: str, max_length: int, suffix: str = "...") -> str:
    """
    Truncate string to maximum length with suffix.
    
    Args:
        s: Input string
        max_length: Maximum length
        suffix: Suffix to add when truncated
    
    Returns:
        Truncated string
    
    Examples:
        >>> truncate("Hello world", 8)
        'Hello...'
    """
    if len(s) <= max_length:
        return s
    return s[:max_length - len(suffix)] + suffix


def pad_left(s: str, length: int, char: str = " ") -> str:
    """
    Pad string to length by adding characters on left.
    
    Args:
        s: Input string
        length: Desired length
        char: Padding character
    
    Returns:
        Padded string
    """
    if len(s) >= length:
        return s
    return char * (length - len(s)) + s


def pad_right(s: str, length: int, char: str = " ") -> str:
    """
    Pad string to length by adding characters on right.
    
    Args:
        s: Input string
        length: Desired length
        char: Padding character
    
    Returns:
        Padded string
    """
    if len(s) >= length:
        return s
    return s + char * (length - len(s))


def pad_center(s: str, length: int, char: str = " ") -> str:
    """
    Pad string to length by adding characters on both sides.
    
    Args:
        s: Input string
        length: Desired length
        char: Padding character
    
    Returns:
        Padded string
    """
    if len(s) >= length:
        return s
    total_pad = length - len(s)
    left_pad = total_pad // 2
    right_pad = total_pad - left_pad
    return char * left_pad + s + char * right_pad


# ============================================================================
# Splitting and Joining
# ============================================================================

def split(s: str, delimiter: str = " ") -> List[str]:
    """
    Split string by delimiter.
    
    Args:
        s: Input string
        delimiter: Delimiter to split on
    
    Returns:
        List of substrings
    """
    return s.split(delimiter)


def split_lines(s: str) -> List[str]:
    """
    Split string into lines.
    
    Args:
        s: Input string
    
    Returns:
        List of lines
    """
    return s.splitlines()


def join(strings: List[str], delimiter: str = "") -> str:
    """
    Join list of strings with delimiter.
    
    Args:
        strings: List of strings
        delimiter: Delimiter between strings
    
    Returns:
        Joined string
    """
    return delimiter.join(strings)


def join_with_commas(strings: List[str]) -> str:
    """
    Join strings with commas and 'and' for last item.
    
    Args:
        strings: List of strings
    
    Returns:
        Formatted string
    """
    if not strings:
        return ""
    if len(strings) == 1:
        return strings[0]
    if len(strings) == 2:
        return f"{strings[0]} and {strings[1]}"
    
    result = ", ".join(strings[:-1])
    return f"{result}, and {strings[-1]}"


# ============================================================================
# Utility Functions
# ============================================================================

def slugify(s: str) -> str:
    """
    Convert string to URL-friendly slug.
    
    Args:
        s: Input string
    
    Returns:
        URL slug
    
    Examples:
        >>> slugify("Hello World!")
        'hello-world'
    """
    # Convert to lowercase
    s = s.lower()
    # Replace spaces with hyphens
    s = s.replace(" ", "-")
    # Remove non-alphanumeric characters
    s = re.sub(r"[^a-z0-9-]", "", s)
    # Remove multiple hyphens
    s = re.sub(r"-+", "-", s)
    # Remove leading/trailing hyphens
    return s.strip("-")


def initials(s: str) -> str:
    """
    Get initials from name.
    
    Args:
        s: Full name
    
    Returns:
        Initials
    """
    words = s.split()
    return "".join(word[0].upper() for word in words if word)


def acronym(s: str) -> str:
    """
    Create acronym from phrase.
    
    Args:
        s: Input phrase
    
    Returns:
        Acronym
    """
    words = s.split()
    return "".join(word[0].upper() for word in words if word)


def pluralize(word: str, count: int = 2) -> str:
    """
    Simple pluralization of English words.
    
    Args:
        word: Word to pluralize
        count: Number of items
    
    Returns:
        Pluralized word if count != 1
    """
    if count == 1:
        return word
    
    # Simple rules
    if word.endswith(("s", "sh", "ch", "x", "z")):
        return word + "es"
    if word.endswith("y") and not word.endswith(("ay", "ey", "iy", "oy", "uy")):
        return word[:-1] + "ies"
    return word + "s"


def is_email(email: str) -> bool:
    """
    Simple email validation.
    
    Args:
        email: Email address
    
    Returns:
        True if valid email format, False otherwise
    """
    pattern = r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$"
    return bool(re.match(pattern, email))


def extract_emails(text: str) -> List[str]:
    """
    Extract email addresses from text.
    
    Args:
        text: Text to search
    
    Returns:
        List of email addresses found
    """
    pattern = r"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}"
    return re.findall(pattern, text)