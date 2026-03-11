```markdown
# Python API Reference

This document provides a comprehensive reference for the Python modules and their APIs.

## 📦 Modules

- [math_utils](#module-math_utils) - Mathematical operations
- [string_utils](#module-string_utils) - String manipulation utilities
- [file_utils](#module-file_utils) - File system operations
- [validator](#module-validator) - Input validation functions
- [converter](#module-converter) - Data conversion utilities
- [async_utils](#module-async_utils) - Asynchronous utilities
- [cache](#module-cache) - Caching mechanisms
- [datetime_utils](#module-datetime_utils) - Date and time utilities
- [encryption](#module-encryption) - Encryption and hashing

---

## Module `math_utils`

Provides mathematical operations with error handling and precision control.

### Installation

```python
from math_utils import (
    add, subtract, multiply, divide,
    power, sqrt, factorial, fibonacci,
    gcd, lcm, is_prime, clamp
)
```

### Functions

#### `add(a: float, b: float, *args: float) -> float`

Adds two or more numbers together.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `a` | float | First number |
| `b` | float | Second number |
| `*args` | float | Additional numbers (optional) |

**Returns:** `float` - The sum of all numbers

**Raises:** `TypeError` - If any argument is not a number

**Example:**
```python
>>> add(5.2, 3.1)
8.3
>>> add(1, 2, 3, 4, 5)
15
```

#### `subtract(a: float, b: float, *args: float) -> float`

Subtracts numbers sequentially.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `a` | float | First number |
| `b` | float | Second number |
| `*args` | float | Additional numbers to subtract |

**Returns:** `float` - The result of a - b - args...

**Example:**
```python
>>> subtract(10.5, 4.2)
6.3
>>> subtract(100, 20, 30, 10)
40
```

#### `multiply(a: float, b: float, *args: float) -> float`

Multiplies two or more numbers together.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `a` | float | First number |
| `b` | float | Second number |
| `*args` | float | Additional numbers |

**Returns:** `float` - The product of all numbers

**Example:**
```python
>>> multiply(3.0, 4.5)
13.5
>>> multiply(2, 3, 4)
24
```

#### `divide(a: float, b: float, *args: float) -> float`

Divides numbers sequentially.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `a` | float | Dividend |
| `b` | float | Divisor |
| `*args` | float | Additional divisors |

**Returns:** `float` - The result of a / b / args...

**Raises:** `ZeroDivisionError` - If any divisor is zero

**Example:**
```python
>>> divide(10.0, 2.0)
5.0
>>> divide(100, 2, 5)
10.0
```

#### `power(base: float, exponent: float) -> float`

Raises a base to the power of an exponent.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `base` | float | The base number |
| `exponent` | float | The exponent |

**Returns:** `float` - base raised to exponent

**Example:**
```python
>>> power(2.0, 3.0)
8.0
>>> power(4.0, 0.5)
2.0
```

#### `sqrt(x: float) -> float`

Calculates the square root of a number.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `x` | float | The number (must be non-negative) |

**Returns:** `float` - The square root of x

**Raises:** `ValueError` - If x is negative

**Example:**
```python
>>> sqrt(16.0)
4.0
>>> sqrt(2.0)
1.4142135623730951
```

#### `factorial(n: int) -> int`

Calculates the factorial of a non-negative integer.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `n` | int | Non-negative integer |

**Returns:** `int` - n! (n factorial)

**Raises:** `ValueError` - If n is negative

**Example:**
```python
>>> factorial(5)
120
>>> factorial(0)
1
```

#### `fibonacci(n: int) -> int`

Calculates the nth Fibonacci number.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `n` | int | Position in Fibonacci sequence (non-negative) |

**Returns:** `int` - The nth Fibonacci number

**Raises:** `ValueError` - If n is negative

**Example:**
```python
>>> fibonacci(0)
0
>>> fibonacci(1)
1
>>> fibonacci(10)
55
```

#### `gcd(a: int, b: int) -> int`

Calculates the greatest common divisor of two integers.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `a` | int | First integer |
| `b` | int | Second integer |

**Returns:** `int` - GCD of a and b

**Example:**
```python
>>> gcd(48, 18)
6
>>> gcd(17, 19)
1
```

#### `lcm(a: int, b: int) -> int`

Calculates the least common multiple of two integers.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `a` | int | First integer |
| `b` | int | Second integer |

**Returns:** `int` - LCM of a and b

**Example:**
```python
>>> lcm(12, 18)
36
>>> lcm(17, 19)
323
```

#### `is_prime(n: int) -> bool`

Determines if a number is prime.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `n` | int | Number to check |

**Returns:** `bool` - True if n is prime, False otherwise

**Example:**
```python
>>> is_prime(17)
True
>>> is_prime(21)
False
```

#### `clamp(value: float, min_val: float, max_val: float) -> float`

Clamps a value between a minimum and maximum range.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `value` | float | Value to clamp |
| `min_val` | float | Minimum allowed value |
| `max_val` | float | Maximum allowed value |

**Returns:** `float` - Clamped value

**Example:**
```python
>>> clamp(5, 1, 10)
5
>>> clamp(0, 1, 10)
1
>>> clamp(15, 1, 10)
10
```

#### `lerp(a: float, b: float, t: float) -> float`

Linearly interpolates between two values.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `a` | float | Start value |
| `b` | float | End value |
| `t` | float | Interpolation factor (0-1) |

**Returns:** `float` - Interpolated value

**Example:**
```python
>>> lerp(0, 10, 0.5)
5.0
>>> lerp(0, 10, 0.75)
7.5
```

---

## Module `string_utils`

Provides string manipulation utilities.

### Functions

#### `reverse(s: str) -> str`

Reverses a string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to reverse |

**Returns:** `str` - The reversed string

**Example:**
```python
>>> reverse('hello')
'olleh'
>>> reverse('café')
'éfac'
```

#### `to_upper(s: str) -> str`

Converts a string to uppercase.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to convert |

**Returns:** `str` - Uppercase string

**Example:**
```python
>>> to_upper('hello')
'HELLO'
```

#### `to_lower(s: str) -> str`

Converts a string to lowercase.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to convert |

**Returns:** `str` - Lowercase string

**Example:**
```python
>>> to_lower('HELLO')
'hello'
```

#### `capitalize(s: str) -> str`

Capitalizes the first letter of a string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to capitalize |

**Returns:** `str` - Capitalized string

**Example:**
```python
>>> capitalize('hello')
'Hello'
>>> capitalize('hello world')
'Hello world'
```

#### `title_case(s: str) -> str`

Converts a string to title case (each word capitalized).

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to convert |

**Returns:** `str` - Title case string

**Example:**
```python
>>> title_case('hello world')
'Hello World'
>>> title_case('the quick brown fox')
'The Quick Brown Fox'
```

#### `trim(s: str) -> str`

Removes whitespace from both ends of a string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to trim |

**Returns:** `str` - Trimmed string

**Example:**
```python
>>> trim('  hello world  ')
'hello world'
```

#### `split(s: str, delimiter: str = ' ', maxsplit: int = -1) -> List[str]`

Splits a string into a list of substrings.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to split |
| `delimiter` | str | The delimiter to use (default: space) |
| `maxsplit` | int | Maximum number of splits (default: -1 = unlimited) |

**Returns:** `List[str]` - List of substrings

**Example:**
```python
>>> split('a,b,c', ',')
['a', 'b', 'c']
>>> split('a,b,c,d', ',', 2)
['a', 'b', 'c,d']
```

#### `join(items: List[str], delimiter: str = '') -> str`

Joins a list of strings into a single string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `items` | List[str] | List of strings to join |
| `delimiter` | str | The delimiter to use (default: empty) |

**Returns:** `str` - Joined string

**Example:**
```python
>>> join(['a', 'b', 'c'], ',')
'a,b,c'
>>> join(['hello', 'world'], ' ')
'hello world'
```

#### `contains(s: str, substring: str, case_sensitive: bool = True) -> bool`

Checks if a string contains a substring.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to search |
| `substring` | str | The substring to find |
| `case_sensitive` | bool | Whether to respect case (default: True) |

**Returns:** `bool` - True if substring is found

**Example:**
```python
>>> contains('hello world', 'world')
True
>>> contains('Hello World', 'hello', False)
True
```

#### `starts_with(s: str, prefix: str) -> bool`

Checks if a string starts with a given prefix.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to check |
| `prefix` | str | The prefix to look for |

**Returns:** `bool` - True if string starts with prefix

**Example:**
```python
>>> starts_with('hello world', 'hello')
True
```

#### `ends_with(s: str, suffix: str) -> bool`

Checks if a string ends with a given suffix.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to check |
| `suffix` | str | The suffix to look for |

**Returns:** `bool` - True if string ends with suffix

**Example:**
```python
>>> ends_with('hello world', 'world')
True
```

#### `count(s: str, substring: str) -> int`

Counts occurrences of a substring.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to search |
| `substring` | str | The substring to count |

**Returns:** `int` - Number of occurrences

**Example:**
```python
>>> count('hello hello hello', 'hello')
3
```

#### `replace(s: str, old: str, new: str, count: int = -1) -> str`

Replaces occurrences of a substring.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The original string |
| `old` | str | The substring to replace |
| `new` | str | The replacement string |
| `count` | int | Maximum replacements (default: -1 = all) |

**Returns:** `str` - String with replacements

**Example:**
```python
>>> replace('hello world hello', 'hello', 'hi')
'hi world hi'
>>> replace('aaa', 'a', 'b', 2)
'bba'
```

#### `to_camel_case(s: str) -> str`

Converts a string to camelCase.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The input string (snake_case or kebab-case) |

**Returns:** `str` - camelCase version

**Example:**
```python
>>> to_camel_case('user_name')
'userName'
>>> to_camel_case('first-name')
'firstName'
```

#### `to_snake_case(s: str) -> str`

Converts a string to snake_case.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The input string (camelCase or PascalCase) |

**Returns:** `str` - snake_case version

**Example:**
```python
>>> to_snake_case('userName')
'user_name'
>>> to_snake_case('FirstName')
'first_name'
```

#### `to_kebab_case(s: str) -> str`

Converts a string to kebab-case.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The input string |

**Returns:** `str` - kebab-case version

**Example:**
```python
>>> to_kebab_case('userName')
'user-name'
```

#### `truncate(s: str, max_length: int, ellipsis: str = '...') -> str`

Truncates a string to the specified length.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to truncate |
| `max_length` | int | Maximum length |
| `ellipsis` | str | Ellipsis to append (default: '...') |

**Returns:** `str` - Truncated string

**Example:**
```python
>>> truncate('This is a long string', 10)
'This is...'
```

#### `slugify(s: str) -> str`

Creates a URL-friendly slug from a string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | The string to slugify |

**Returns:** `str` - URL-friendly slug

**Example:**
```python
>>> slugify('Hello World!')
'hello-world'
>>> slugify('Café & Restaurant')
'cafe-restaurant'
```

---

## Module `file_utils`

Provides file system operations.

### Functions

#### `exists(path: str) -> bool`

Checks if a file or directory exists.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | str | Path to check |

**Returns:** `bool` - True if path exists

**Example:**
```python
>>> exists('config.json')
True
```

#### `is_file(path: str) -> bool`

Checks if a path is a file.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | str | Path to check |

**Returns:** `bool` - True if path is a file

**Example:**
```python
>>> is_file('config.json')
True
>>> is_file('docs')
False
```

#### `is_dir(path: str) -> bool`

Checks if a path is a directory.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | str | Path to check |

**Returns:** `bool` - True if path is a directory

**Example:**
```python
>>> is_dir('docs')
True
```

#### `read_file(path: str, encoding: str = 'utf-8') -> str`

Reads a file and returns its contents.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | str | Path to the file |
| `encoding` | str | File encoding (default: 'utf-8') |

**Returns:** `str` - File contents

**Raises:** `FileNotFoundError` - If file doesn't exist

**Example:**
```python
>>> content = read_file('data.txt')
>>> print(content)
Hello World
```

#### `write_file(path: str, content: str, encoding: str = 'utf-8') -> None`

Writes content to a file.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | str | Path to the file |
| `content` | str | Content to write |
| `encoding` | str | File encoding (default: 'utf-8') |

**Raises:** `PermissionError` - If write permission is denied

**Example:**
```python
>>> write_file('output.txt', 'Hello World')
```

#### `append_file(path: str, content: str, encoding: str = 'utf-8') -> None`

Appends content to a file.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | str | Path to the file |
| `content` | str | Content to append |
| `encoding` | str | File encoding (default: 'utf-8') |

**Example:**
```python
>>> append_file('log.txt', 'New log entry\n')
```

#### `copy_file(src: str, dst: str) -> None`

Copies a file from source to destination.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `src` | str | Source file path |
| `dst` | str | Destination file path |

**Raises:** `FileNotFoundError` - If source doesn't exist

**Example:**
```python
>>> copy_file('source.txt', 'backup/source.txt')
```

#### `move_file(src: str, dst: str) -> None`

Moves or renames a file.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `src` | str | Source file path |
| `dst` | str | Destination file path |

**Example:**
```python
>>> move_file('old.txt', 'new.txt')
```

#### `delete_file(path: str) -> None`

Deletes a file.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | str | Path to the file |

**Example:**
```python
>>> delete_file('temp.txt')
```

#### `list_files(directory: str, pattern: str = '*') -> List[str]`

Lists files in a directory matching a pattern.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `directory` | str | Directory to search |
| `pattern` | str | File pattern (e.g., '*.txt') |

**Returns:** `List[str]` - List of file paths

**Example:**
```python
>>> list_files('.', '*.py')
['main.py', 'utils.py', 'test.py']
```

#### `get_size(path: str) -> int`

Gets the size of a file in bytes.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | str | Path to the file |

**Returns:** `int` - File size in bytes

**Example:**
```python
>>> size = get_size('data.txt')
>>> print(f'{size} bytes')
1024 bytes
```

#### `get_modified_time(path: str) -> datetime`

Gets the last modified time of a file.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | str | Path to the file |

**Returns:** `datetime` - Last modified time

**Example:**
```python
>>> mtime = get_modified_time('config.json')
>>> print(mtime)
2024-01-15 14:30:00
```

#### `ensure_directory(path: str) -> None`

Ensures a directory exists, creating it if necessary.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | str | Directory path |

**Example:**
```python
>>> ensure_directory('data/subdir/nested')
```

---

## Module `validator`

Provides input validation functions.

### Functions

#### `is_email(email: str) -> bool`

Validates an email address.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `email` | str | Email address to validate |

**Returns:** `bool` - True if valid email

**Example:**
```python
>>> is_email('user@example.com')
True
>>> is_email('invalid')
False
```

#### `is_url(url: str) -> bool`

Validates a URL.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `url` | str | URL to validate |

**Returns:** `bool` - True if valid URL

**Example:**
```python
>>> is_url('https://example.com')
True
>>> is_url('not-a-url')
False
```

#### `is_phone(phone: str) -> bool`

Validates a phone number (international format).

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `phone` | str | Phone number to validate |

**Returns:** `bool` - True if valid phone number

**Example:**
```python
>>> is_phone('+1-555-123-4567')
True
>>> is_phone('123-456-7890')
True
```

#### `is_zip_code(zip_code: str) -> bool`

Validates a US ZIP code.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `zip_code` | str | ZIP code to validate |

**Returns:** `bool` - True if valid ZIP code

**Example:**
```python
>>> is_zip_code('12345')
True
>>> is_zip_code('12345-6789')
True
```

#### `is_ipv4(ip: str) -> bool`

Validates an IPv4 address.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `ip` | str | IPv4 address to validate |

**Returns:** `bool` - True if valid IPv4 address

**Example:**
```python
>>> is_ipv4('192.168.1.1')
True
>>> is_ipv4('256.0.0.1')
False
```

#### `is_ipv6(ip: str) -> bool`

Validates an IPv6 address.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `ip` | str | IPv6 address to validate |

**Returns:** `bool` - True if valid IPv6 address

**Example:**
```python
>>> is_ipv6('2001:db8::1')
True
```

#### `is_uuid(uuid_str: str) -> bool`

Validates a UUID.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `uuid_str` | str | UUID to validate |

**Returns:** `bool` - True if valid UUID

**Example:**
```python
>>> is_uuid('123e4567-e89b-12d3-a456-426614174000')
True
```

#### `is_credit_card(card_number: str) -> bool`

Validates a credit card number using Luhn algorithm.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `card_number` | str | Credit card number |

**Returns:** `bool` - True if valid credit card

**Example:**
```python
>>> is_credit_card('4111111111111111')
True
```

#### `is_alpha(s: str) -> bool`

Checks if string contains only letters.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | String to check |

**Returns:** `bool` - True if alphabetic

**Example:**
```python
>>> is_alpha('abc')
True
>>> is_alpha('abc123')
False
```

#### `is_numeric(s: str) -> bool`

Checks if string contains only numbers.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | String to check |

**Returns:** `bool` - True if numeric

**Example:**
```python
>>> is_numeric('123')
True
>>> is_numeric('12.3')
False
```

#### `is_alphanumeric(s: str) -> bool`

Checks if string contains only letters and numbers.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `s` | str | String to check |

**Returns:** `bool` - True if alphanumeric

**Example:**
```python
>>> is_alphanumeric('abc123')
True
>>> is_alphanumeric('abc-123')
False
```

#### `is_strong_password(password: str, min_length: int = 8) -> bool`

Checks if a password meets strength requirements.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `password` | str | Password to check |
| `min_length` | int | Minimum length (default: 8) |

**Returns:** `bool` - True if password is strong

**Requirements:**
- At least `min_length` characters
- At least one uppercase letter
- At least one lowercase letter
- At least one digit
- At least one special character

**Example:**
```python
>>> is_strong_password('Passw0rd!')
True
>>> is_strong_password('password')
False
```

---

## Module `converter`

Provides data conversion utilities.

### Functions

#### `to_int(value: Any, default: int = 0) -> int`

Converts a value to an integer.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `value` | Any | Value to convert |
| `default` | int | Default if conversion fails |

**Returns:** `int` - Converted integer

**Example:**
```python
>>> to_int('42')
42
>>> to_int('3.14')
3
>>> to_int('abc', -1)
-1
```

#### `to_float(value: Any, default: float = 0.0) -> float`

Converts a value to a float.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `value` | Any | Value to convert |
| `default` | float | Default if conversion fails |

**Returns:** `float` - Converted float

**Example:**
```python
>>> to_float('3.14')
3.14
>>> to_float('42')
42.0
```

#### `to_bool(value: Any, default: bool = False) -> bool`

Converts a value to a boolean.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `value` | Any | Value to convert |
| `default` | bool | Default if conversion fails |

**Returns:** `bool` - Converted boolean

**Truthy values:** 'true', 'yes', '1', 'on', True, 1
**Falsy values:** 'false', 'no', '0', 'off', False, 0, None

**Example:**
```python
>>> to_bool('true')
True
>>> to_bool('false')
False
>>> to_bool('1')
True
```

#### `to_str(value: Any) -> str`

Converts any value to a string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `value` | Any | Value to convert |

**Returns:** `str` - String representation

**Example:**
```python
>>> to_str(42)
'42'
>>> to_str(3.14)
'3.14'
>>> to_str([1, 2, 3])
'[1, 2, 3]'
```

#### `to_json(data: Any, indent: int = None) -> str`

Converts data to a JSON string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `data` | Any | Data to serialize |
| `indent` | int | Pretty print indentation |

**Returns:** `str` - JSON string

**Raises:** `TypeError` - If data is not JSON serializable

**Example:**
```python
>>> to_json({'name': 'John', 'age': 30})
'{"name": "John", "age": 30}'
```

#### `from_json(json_str: str) -> Any`

Parses a JSON string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `json_str` | str | JSON string to parse |

**Returns:** `Any` - Parsed data

**Raises:** `JSONDecodeError` - If JSON is invalid

**Example:**
```python
>>> data = from_json('{"name": "John", "age": 30}')
>>> print(data['name'])
John
```

#### `to_base64(data: Union[str, bytes]) -> str`

Encodes data to Base64.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `data` | Union[str, bytes] | Data to encode |

**Returns:** `str` - Base64 encoded string

**Example:**
```python
>>> to_base64('hello')
'aGVsbG8='
```

#### `from_base64(encoded: str) -> bytes`

Decodes a Base64 string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `encoded` | str | Base64 string to decode |

**Returns:** `bytes` - Decoded bytes

**Example:**
```python
>>> from_base64('aGVsbG8=').decode()
'hello'
```

#### `to_hex(data: Union[str, bytes]) -> str`

Encodes data to hexadecimal.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `data` | Union[str, bytes] | Data to encode |

**Returns:** `str` - Hexadecimal string

**Example:**
```python
>>> to_hex('hello')
'68656c6c6f'
```

#### `from_hex(hex_str: str) -> bytes`

Decodes a hexadecimal string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `hex_str` | str | Hexadecimal string to decode |

**Returns:** `bytes` - Decoded bytes

**Example:**
```python
>>> from_hex('68656c6c6f').decode()
'hello'
```

---

## Module `async_utils`

Provides asynchronous utilities.

### Functions

#### `async def delay(seconds: float) -> None`

Delays execution for a specified number of seconds.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `seconds` | float | Delay in seconds |

**Example:**
```python
await delay(1.5)  # Wait 1.5 seconds
```

#### `async def timeout(coro, seconds: float) -> Any`

Adds a timeout to an awaitable.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `coro` | Awaitable | Coroutine or awaitable |
| `seconds` | float | Timeout in seconds |

**Returns:** `Any` - Result of the awaitable

**Raises:** `TimeoutError` - If operation times out

**Example:**
```python
try:
    result = await timeout(fetch_data(), 5.0)
except TimeoutError:
    print('Operation timed out')
```

#### `async def retry(coro, attempts: int = 3, delay: float = 1.0, backoff: float = 2.0) -> Any`

Retries a coroutine on failure.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `coro` | Callable | Coroutine function to retry |
| `attempts` | int | Maximum number of attempts |
| `delay` | float | Initial delay between attempts |
| `backoff` | float | Multiplier for exponential backoff |

**Returns:** `Any` - Result of the coroutine

**Example:**
```python
result = await retry(
    lambda: fetch_unstable_api(),
    attempts=5,
    delay=1.0,
    backoff=2.0
)
```

#### `async def parallel(tasks: List[Callable], max_concurrency: int = None) -> List[Any]`

Runs async tasks in parallel with concurrency limit.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `tasks` | List[Callable] | List of async functions |
| `max_concurrency` | int | Maximum concurrent tasks |

**Returns:** `List[Any]` - Results in original order

**Example:**
```python
results = await parallel([
    fetch_user(1),
    fetch_user(2),
    fetch_user(3)
], max_concurrency=2)
```

#### `async def gather_with_limit(tasks: List[Callable], limit: int) -> List[Any]`

Gathers tasks with a concurrency limit.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `tasks` | List[Callable] | List of async functions |
| `limit` | int | Concurrency limit |

**Returns:** `List[Any]` - Results in original order

**Example:**
```python
results = await gather_with_limit(
    [process_item(i) for i in items],
    limit=10
)
```

---

## Module `cache`

Provides caching utilities.

### Classes

#### `class Cache`

In-memory cache with TTL support.

**Constructor:**
```python
cache = Cache(ttl: int = 300, max_size: int = 1000)
```

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `ttl` | int | Time-to-live in seconds (default: 300) |
| `max_size` | int | Maximum number of items (default: 1000) |

**Methods:**

##### `set(key: str, value: Any, ttl: int = None) -> None`

Sets a value in the cache.

**Parameters:**
| Name | Type | Description |
|------|------|-------------