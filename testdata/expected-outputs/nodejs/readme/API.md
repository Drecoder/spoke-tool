```markdown
# Node.js API Reference

This document provides a comprehensive reference for the Node.js modules and their APIs.

## 📦 Modules

- [math](#module-math) - Mathematical operations
- [string](#module-string) - String manipulation utilities
- [file](#module-file) - File system operations
- [validator](#module-validator) - Input validation functions
- [converter](#module-converter) - Data conversion utilities
- [async](#module-async) - Asynchronous utilities
- [cache](#module-cache) - Caching mechanisms

---

## Module `math`

Provides mathematical operations with error handling and precision control.

### Installation

```javascript
const math = require('./math');
// or
import * as math from './math.js';
```

### Functions

#### `add(a, b)`

Adds two numbers together.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `a` | number | First number |
| `b` | number | Second number |

**Returns:** `number` - The sum of a and b

**Example:**
```javascript
const result = math.add(5.2, 3.1);
console.log(result); // 8.3
```

#### `subtract(a, b)`

Subtracts the second number from the first.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `a` | number | First number |
| `b` | number | Second number |

**Returns:** `number` - The difference a - b

**Example:**
```javascript
const result = math.subtract(10.5, 4.2);
console.log(result); // 6.3
```

#### `multiply(a, b)`

Multiplies two numbers together.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `a` | number | First number |
| `b` | number | Second number |

**Returns:** `number` - The product a * b

**Example:**
```javascript
const result = math.multiply(3.0, 4.5);
console.log(result); // 13.5
```

#### `divide(a, b)`

Divides the first number by the second.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `a` | number | Dividend |
| `b` | number | Divisor |

**Returns:** `number` - The quotient a / b

**Throws:** `Error` - If b is zero

**Example:**
```javascript
try {
    const result = math.divide(10.0, 2.0);
    console.log(result); // 5.0
} catch (err) {
    console.error(err.message);
}
```

#### `power(base, exponent)`

Raises a base to the power of an exponent.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `base` | number | The base number |
| `exponent` | number | The exponent |

**Returns:** `number` - base raised to exponent

**Example:**
```javascript
const result = math.power(2.0, 3.0);
console.log(result); // 8.0
```

#### `sqrt(x)`

Calculates the square root of a number.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `x` | number | The number (must be non-negative) |

**Returns:** `number` - The square root of x

**Throws:** `Error` - If x is negative

**Example:**
```javascript
const result = math.sqrt(16.0);
console.log(result); // 4.0
```

#### `abs(x)`

Returns the absolute value of a number.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `x` | number | The input number |

**Returns:** `number` - The absolute value of x

**Example:**
```javascript
console.log(math.abs(-42)); // 42
console.log(math.abs(42));  // 42
```

#### `max(...numbers)`

Returns the largest of the given numbers.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `...numbers` | number[] | Variable number of arguments |

**Returns:** `number` - The maximum value

**Example:**
```javascript
const result = math.max(5, 10, 3, 8, 2);
console.log(result); // 10
```

#### `min(...numbers)`

Returns the smallest of the given numbers.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `...numbers` | number[] | Variable number of arguments |

**Returns:** `number` - The minimum value

**Example:**
```javascript
const result = math.min(5, 10, 3, 8, 2);
console.log(result); // 2
```

#### `round(value, precision)`

Rounds a number to the specified precision.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `value` | number | The number to round |
| `precision` | number | Number of decimal places (default: 0) |

**Returns:** `number` - The rounded number

**Example:**
```javascript
console.log(math.round(3.14159, 2)); // 3.14
console.log(math.round(3.14159, 0)); // 3
```

---

## Module `string`

Provides string manipulation utilities.

### Functions

#### `reverse(str)`

Reverses a string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to reverse |

**Returns:** `string` - The reversed string

**Example:**
```javascript
const result = string.reverse('hello');
console.log(result); // 'olleh'
```

#### `toUpper(str)`

Converts a string to uppercase.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to convert |

**Returns:** `string` - Uppercase string

**Example:**
```javascript
const result = string.toUpper('hello');
console.log(result); // 'HELLO'
```

#### `toLower(str)`

Converts a string to lowercase.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to convert |

**Returns:** `string` - Lowercase string

**Example:**
```javascript
const result = string.toLower('HELLO');
console.log(result); // 'hello'
```

#### `trim(str)`

Removes whitespace from both ends of a string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to trim |

**Returns:** `string` - Trimmed string

**Example:**
```javascript
const result = string.trim('  hello world  ');
console.log(result); // 'hello world'
```

#### `split(str, separator)`

Splits a string into an array of substrings.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to split |
| `separator` | string | The separator to use |

**Returns:** `string[]` - Array of substrings

**Example:**
```javascript
const result = string.split('a,b,c', ',');
console.log(result); // ['a', 'b', 'c']
```

#### `join(arr, separator)`

Joins an array of strings into a single string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `arr` | string[] | Array of strings to join |
| `separator` | string | The separator to use (default: '') |

**Returns:** `string` - Joined string

**Example:**
```javascript
const result = string.join(['a', 'b', 'c'], ',');
console.log(result); // 'a,b,c'
```

#### `contains(str, substring)`

Checks if a string contains a substring.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to search |
| `substring` | string | The substring to find |

**Returns:** `boolean` - True if substring is found

**Example:**
```javascript
const result = string.contains('hello world', 'world');
console.log(result); // true
```

#### `replace(str, search, replacement)`

Replaces occurrences of a substring.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The original string |
| `search` | string\|RegExp | The pattern to search for |
| `replacement` | string | The replacement string |

**Returns:** `string` - String with replacements

**Example:**
```javascript
const result = string.replace('hello world', 'world', 'there');
console.log(result); // 'hello there'
```

#### `replaceAll(str, search, replacement)`

Replaces all occurrences of a substring.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The original string |
| `search` | string | The substring to search for |
| `replacement` | string | The replacement string |

**Returns:** `string` - String with all replacements

**Example:**
```javascript
const result = string.replaceAll('hello hello hello', 'hello', 'hi');
console.log(result); // 'hi hi hi'
```

#### `indexOf(str, substring)`

Returns the index of the first occurrence of a substring.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to search |
| `substring` | string | The substring to find |

**Returns:** `number` - The index, or -1 if not found

**Example:**
```javascript
const result = string.indexOf('hello world', 'world');
console.log(result); // 6
```

#### `lastIndexOf(str, substring)`

Returns the index of the last occurrence of a substring.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to search |
| `substring` | string | The substring to find |

**Returns:** `number` - The index, or -1 if not found

**Example:**
```javascript
const result = string.lastIndexOf('hello world hello', 'hello');
console.log(result); // 12
```

#### `startsWith(str, prefix)`

Checks if a string starts with a given prefix.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to check |
| `prefix` | string | The prefix to look for |

**Returns:** `boolean` - True if string starts with prefix

**Example:**
```javascript
const result = string.startsWith('hello world', 'hello');
console.log(result); // true
```

#### `endsWith(str, suffix)`

Checks if a string ends with a given suffix.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to check |
| `suffix` | string | The suffix to look for |

**Returns:** `boolean` - True if string ends with suffix

**Example:**
```javascript
const result = string.endsWith('hello world', 'world');
console.log(result); // true
```

#### `count(str, substring)`

Counts occurrences of a substring.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to search |
| `substring` | string | The substring to count |

**Returns:** `number` - Number of occurrences

**Example:**
```javascript
const result = string.count('hello hello hello', 'hello');
console.log(result); // 3
```

#### `repeat(str, count)`

Repeats a string a specified number of times.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to repeat |
| `count` | number | Number of repetitions |

**Returns:** `string` - Repeated string

**Throws:** `Error` - If count is negative

**Example:**
```javascript
const result = string.repeat('ha', 3);
console.log(result); // 'hahaha'
```

#### `toCamelCase(str)`

Converts a string to camelCase.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The input string (snake_case or kebab-case) |

**Returns:** `string` - camelCase version

**Example:**
```javascript
console.log(string.toCamelCase('user_name'));     // 'userName'
console.log(string.toCamelCase('first-name'));    // 'firstName'
console.log(string.toCamelCase('hello world'));   // 'helloWorld'
```

#### `toPascalCase(str)`

Converts a string to PascalCase.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The input string |

**Returns:** `string` - PascalCase version

**Example:**
```javascript
console.log(string.toPascalCase('user_name'));    // 'UserName'
console.log(string.toPascalCase('first-name'));   // 'FirstName'
console.log(string.toPascalCase('hello world'));  // 'HelloWorld'
```

#### `toSnakeCase(str)`

Converts a string to snake_case.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The input string (camelCase or PascalCase) |

**Returns:** `string` - snake_case version

**Example:**
```javascript
console.log(string.toSnakeCase('userName'));      // 'user_name'
console.log(string.toSnakeCase('FirstName'));     // 'first_name'
console.log(string.toSnakeCase('helloWorld'));    // 'hello_world'
```

#### `truncate(str, maxLength)`

Truncates a string to the specified length.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to truncate |
| `maxLength` | number | Maximum length |

**Returns:** `string` - Truncated string with '...' if needed

**Example:**
```javascript
const result = string.truncate('This is a long string', 10);
console.log(result); // 'This is...'
```

#### `padStart(str, targetLength, padString)`

Pads the beginning of a string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to pad |
| `targetLength` | number | Desired length |
| `padString` | string | String to pad with (default: ' ') |

**Returns:** `string` - Padded string

**Example:**
```javascript
const result = string.padStart('5', 3, '0');
console.log(result); // '005'
```

#### `padEnd(str, targetLength, padString)`

Pads the end of a string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | The string to pad |
| `targetLength` | number | Desired length |
| `padString` | string | String to pad with (default: ' ') |

**Returns:** `string` - Padded string

**Example:**
```javascript
const result = string.padEnd('5', 3, '0');
console.log(result); // '500'
```

---

## Module `file`

Provides file system operations.

### Functions

#### `exists(path)`

Checks if a file or directory exists.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | string | Path to check |

**Returns:** `Promise<boolean>` - True if path exists

**Example:**
```javascript
const exists = await file.exists('./config.json');
console.log(exists); // true or false
```

#### `readFile(path, encoding)`

Reads a file asynchronously.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | string | Path to the file |
| `encoding` | string | File encoding (default: 'utf8') |

**Returns:** `Promise<string>` - File contents

**Throws:** `Error` - If file doesn't exist or can't be read

**Example:**
```javascript
try {
    const content = await file.readFile('./data.txt');
    console.log(content);
} catch (err) {
    console.error(err.message);
}
```

#### `readFileSync(path, encoding)`

Reads a file synchronously.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | string | Path to the file |
| `encoding` | string | File encoding (default: 'utf8') |

**Returns:** `string` - File contents

**Throws:** `Error` - If file doesn't exist or can't be read

**Example:**
```javascript
try {
    const content = file.readFileSync('./data.txt');
    console.log(content);
} catch (err) {
    console.error(err.message);
}
```

#### `writeFile(path, data)`

Writes data to a file asynchronously.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | string | Path to the file |
| `data` | string\|Buffer | Data to write |

**Returns:** `Promise<void>`

**Throws:** `Error` - If file can't be written

**Example:**
```javascript
try {
    await file.writeFile('./output.txt', 'Hello World');
    console.log('File written successfully');
} catch (err) {
    console.error(err.message);
}
```

#### `writeFileSync(path, data)`

Writes data to a file synchronously.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | string | Path to the file |
| `data` | string\|Buffer | Data to write |

**Throws:** `Error` - If file can't be written

**Example:**
```javascript
try {
    file.writeFileSync('./output.txt', 'Hello World');
    console.log('File written successfully');
} catch (err) {
    console.error(err.message);
}
```

#### `copyFile(src, dest)`

Copies a file asynchronously.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `src` | string | Source path |
| `dest` | string | Destination path |

**Returns:** `Promise<void>`

**Throws:** `Error` - If source doesn't exist or copy fails

**Example:**
```javascript
try {
    await file.copyFile('./source.txt', './backup/source.txt');
    console.log('File copied successfully');
} catch (err) {
    console.error(err.message);
}
```

#### `mkdir(path, options)`

Creates a directory asynchronously.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | string | Directory path |
| `options` | Object | Options (recursive, mode) |

**Returns:** `Promise<void>`

**Example:**
```javascript
await file.mkdir('./data/nested', { recursive: true });
```

#### `readdir(path)`

Reads directory contents asynchronously.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | string | Directory path |

**Returns:** `Promise<string[]>` - Array of filenames

**Example:**
```javascript
const files = await file.readdir('./data');
console.log(files);
```

#### `stat(path)`

Gets file statistics asynchronously.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | string | Path to file or directory |

**Returns:** `Promise<Stats>` - File statistics

**Example:**
```javascript
const stats = await file.stat('./data.txt');
console.log(stats.size);
console.log(stats.isFile());
```

#### `unlink(path)`

Deletes a file asynchronously.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | string | Path to file |

**Returns:** `Promise<void>`

**Example:**
```javascript
await file.unlink('./temp.txt');
```

#### `rmdir(path, options)`

Deletes a directory asynchronously.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `path` | string | Directory path |
| `options` | Object | Options (recursive) |

**Returns:** `Promise<void>`

**Example:**
```javascript
await file.rmdir('./temp', { recursive: true });
```

---

## Module `validator`

Provides input validation functions.

### Functions

#### `isEmail(str)`

Validates an email address.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | Email address to validate |

**Returns:** `boolean` - True if valid email

**Example:**
```javascript
console.log(validator.isEmail('user@example.com')); // true
console.log(validator.isEmail('invalid'));          // false
```

#### `isURL(str)`

Validates a URL.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | URL to validate |

**Returns:** `boolean` - True if valid URL

**Example:**
```javascript
console.log(validator.isURL('https://example.com')); // true
console.log(validator.isURL('not-a-url'));           // false
```

#### `isPhoneNumber(str)`

Validates a phone number (US format).

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | Phone number to validate |

**Returns:** `boolean` - True if valid phone number

**Example:**
```javascript
console.log(validator.isPhoneNumber('555-123-4567')); // true
console.log(validator.isPhoneNumber('123'));          // false
```

#### `isZipCode(str)`

Validates a US ZIP code.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | ZIP code to validate |

**Returns:** `boolean` - True if valid ZIP code

**Example:**
```javascript
console.log(validator.isZipCode('12345'));      // true
console.log(validator.isZipCode('12345-6789')); // true
console.log(validator.isZipCode('123'));        // false
```

#### `isCreditCard(str)`

Validates a credit card number (Luhn algorithm).

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | Credit card number |

**Returns:** `boolean` - True if valid credit card

**Example:**
```javascript
console.log(validator.isCreditCard('4111111111111111')); // true
console.log(validator.isCreditCard('1234567890123456')); // false
```

#### `isIP(str)`

Validates an IP address (v4 or v6).

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | IP address to validate |

**Returns:** `boolean` - True if valid IP address

**Example:**
```javascript
console.log(validator.isIP('192.168.1.1'));           // true
console.log(validator.isIP('2001:db8::1'));           // true
console.log(validator.isIP('not-an-ip'));             // false
```

#### `isUUID(str)`

Validates a UUID.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | UUID to validate |

**Returns:** `boolean` - True if valid UUID

**Example:**
```javascript
const uuid = '123e4567-e89b-12d3-a456-426614174000';
console.log(validator.isUUID(uuid)); // true
```

#### `isAlphanumeric(str)`

Checks if string contains only letters and numbers.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | String to check |

**Returns:** `boolean` - True if alphanumeric

**Example:**
```javascript
console.log(validator.isAlphanumeric('abc123')); // true
console.log(validator.isAlphanumeric('abc-123')); // false
```

#### `isNumeric(str)`

Checks if string contains only numbers.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | String to check |

**Returns:** `boolean` - True if numeric

**Example:**
```javascript
console.log(validator.isNumeric('12345')); // true
console.log(validator.isNumeric('123a'));  // false
```

#### `isAlpha(str)`

Checks if string contains only letters.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | String to check |

**Returns:** `boolean` - True if alphabetic

**Example:**
```javascript
console.log(validator.isAlpha('abc'));   // true
console.log(validator.isAlpha('abc123')); // false
```

#### `isLength(str, min, max)`

Checks if string length is within range.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | String to check |
| `min` | number | Minimum length |
| `max` | number | Maximum length |

**Returns:** `boolean` - True if length is valid

**Example:**
```javascript
console.log(validator.isLength('hello', 3, 10)); // true
console.log(validator.isLength('hi', 3, 10));    // false
```

#### `isInRange(value, min, max)`

Checks if a number is within range.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `value` | number | Number to check |
| `min` | number | Minimum value |
| `max` | number | Maximum value |

**Returns:** `boolean` - True if value is in range

**Example:**
```javascript
console.log(validator.isInRange(5, 1, 10));  // true
console.log(validator.isInRange(15, 1, 10)); // false
```

#### `matches(str, pattern)`

Checks if string matches a regex pattern.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | String to check |
| `pattern` | RegExp\|string | Regex pattern |

**Returns:** `boolean` - True if string matches

**Example:**
```javascript
console.log(validator.matches('abc123', /^[a-z0-9]+$/)); // true
```

#### `isRequired(value)`

Checks if a value is not empty.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `value` | any | Value to check |

**Returns:** `boolean` - True if value is not empty

**Example:**
```javascript
console.log(validator.isRequired('hello')); // true
console.log(validator.isRequired(''));       // false
console.log(validator.isRequired(null));     // false
```

---

## Module `converter`

Provides data conversion utilities.

### Functions

#### `toNumber(str)`

Converts a string to a number.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | String to convert |

**Returns:** `number` - Converted number, or NaN if invalid

**Example:**
```javascript
console.log(converter.toNumber('42'));     // 42
console.log(converter.toNumber('3.14'));   // 3.14
console.log(converter.toNumber('abc'));    // NaN
```

#### `toInt(str)`

Converts a string to an integer.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | String to convert |

**Returns:** `number` - Integer value, or NaN if invalid

**Example:**
```javascript
console.log(converter.toInt('42'));     // 42
console.log(converter.toInt('3.14'));   // 3
console.log(converter.toInt('abc'));    // NaN
```

#### `toFloat(str)`

Converts a string to a float.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | String to convert |

**Returns:** `number` - Float value, or NaN if invalid

**Example:**
```javascript
console.log(converter.toFloat('3.14'));  // 3.14
console.log(converter.toFloat('42'));    // 42.0
```

#### `toBoolean(str)`

Converts a string to a boolean.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | String to convert |

**Returns:** `boolean` - Boolean value

**Example:**
```javascript
console.log(converter.toBoolean('true'));   // true
console.log(converter.toBoolean('false'));  // false
console.log(converter.toBoolean('1'));      // true
console.log(converter.toBoolean('0'));      // false
console.log(converter.toBoolean('yes'));    // true
console.log(converter.toBoolean('no'));     // false
```

#### `toString(value)`

Converts any value to a string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `value` | any | Value to convert |

**Returns:** `string` - String representation

**Example:**
```javascript
console.log(converter.toString(42));      // '42'
console.log(converter.toString(true));    // 'true'
console.log(converter.toString({a:1}));   // '{"a":1}'
```

#### `toJSON(value)`

Converts a value to a JSON string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `value` | any | Value to convert |

**Returns:** `string` - JSON string

**Throws:** `Error` - If value cannot be stringified

**Example:**
```javascript
const obj = { name: 'John', age: 30 };
console.log(converter.toJSON(obj)); // '{"name":"John","age":30}'
```

#### `fromJSON(json)`

Parses a JSON string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `json` | string | JSON string to parse |

**Returns:** `any` - Parsed value

**Throws:** `Error` - If JSON is invalid

**Example:**
```javascript
const obj = converter.fromJSON('{"name":"John","age":30}');
console.log(obj.name); // 'John'
```

#### `toBase64(str)`

Converts a string to Base64.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | String to encode |

**Returns:** `string` - Base64 encoded string

**Example:**
```javascript
console.log(converter.toBase64('hello')); // 'aGVsbG8='
```

#### `fromBase64(str)`

Decodes a Base64 string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | Base64 string to decode |

**Returns:** `string` - Decoded string

**Example:**
```javascript
console.log(converter.fromBase64('aGVsbG8=')); // 'hello'
```

#### `toHex(str)`

Converts a string to hexadecimal.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `str` | string | String to encode |

**Returns:** `string` - Hexadecimal representation

**Example:**
```javascript
console.log(converter.toHex('hello')); // '68656c6c6f'
```

#### `fromHex(hex)`

Decodes a hexadecimal string.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `hex` | string | Hexadecimal string to decode |

**Returns:** `string` - Decoded string

**Example:**
```javascript
console.log(converter.fromHex('68656c6c6f')); // 'hello'
```

---

## Module `async`

Provides asynchronous utilities.

### Functions

#### `delay(ms)`

Creates a promise that resolves after a delay.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `ms` | number | Delay in milliseconds |

**Returns:** `Promise<void>`

**Example:**
```javascript
await async.delay(1000);
console.log('1 second later');
```

#### `timeout(promise, ms)`

Adds a timeout to a promise.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `promise` | Promise | The promise to add timeout to |
| `ms` | number | Timeout in milliseconds |

**Returns:** `Promise` - Original promise or timeout error

**Throws:** `Error` - If promise times out

**Example:**
```javascript
try {
    const result = await async.timeout(fetch('/api/data'), 5000);
    console.log(result);
} catch (err) {
    console.error('Request timed out');
}
```

#### `retry(fn, options)`

Retries a function on failure.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `fn` | Function | Async function to retry |
| `options` | Object | Retry options (attempts, delay, backoff) |

**Returns:** `Promise` - Result of the function

**Example:**
```javascript
const result = await async.retry(
    () => fetch('/api/unstable'),
    { attempts: 3, delay: 1000, backoff: 2 }
);
```

#### `parallel(tasks, concurrency)`

Runs async tasks in parallel with concurrency limit.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `tasks` | Function[] | Array of async functions |
| `concurrency` | number | Max concurrent tasks |

**Returns:** `Promise<Array>` - Results of all tasks

**Example:**
```javascript
const results = await async.parallel([
    () => fetch('/api/user/1'),
    () => fetch('/api/user/2'),
    () => fetch('/api/user/3')
], 2);
```

#### `series(tasks)`

Runs async tasks in series.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `tasks` | Function[] | Array of async functions |

**Returns:** `Promise<Array>` - Results of all tasks

**Example:**
```javascript
const results = await async.series([
    () => db.insert(user1),
    () => db.insert(user2),
    () => db.insert(user3)
]);
```

#### `waterfall(tasks)`

Runs async tasks in a waterfall (each passes result to next).

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `tasks` | Function[] | Array of async functions |

**Returns:** `Promise<any>` - Final result

**Example:**
```javascript
const result = await async.waterfall([
    () => getUser(1),
    (user) => getOrders(user.id),
    (orders) => calculateTotal(orders)
]);
```

---

## Module `cache`

Provides caching utilities.

### Classes

#### `class Cache`

In-memory cache with TTL support.

**Constructor:**
```javascript
const cache = new Cache(options);
```

**Options:**
| Name | Type | Description |
|------|------|-------------|
| `ttl` | number | Default TTL in milliseconds |
| `maxSize` | number | Maximum number of items |
| `checkPeriod` | number | Cleanup interval |

**Methods:**

##### `set(key, value, ttl)`

Sets a value in the cache.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `key` | string | Cache key |
| `value` | any | Value to store |
| `ttl` | number | Optional TTL (overrides default) |

**Returns:** `boolean` - True if successful

**Example:**
```javascript
cache.set('user:1', { name: 'John' }, 60000);
```

##### `get(key)`

Gets a value from the cache.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `key` | string | Cache key |

**Returns:** `any` - Stored value or undefined

**Example:**
```javascript
const user = cache.get('user:1');
```

##### `has(key)`

Checks if a key exists in the cache.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `key` | string | Cache key |

**Returns:** `boolean` - True if key exists

**Example:**
```javascript
if (cache.has('user:1')) {
    // Use cached value
}
```

##### `delete(key)`

Deletes a key from the cache.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| `key` | string | Cache key |

**Returns:** `boolean` - True if deleted

**Example:**
```javascript
cache.delete('user:1');
```

##### `clear()`

Clears all items from the cache.

**Example:**
```javascript
cache.clear();
```

##### `size()`

Returns the number of items in the cache.

**Returns:** `number` - Cache size

**Example:**
```javascript
console.log(cache.size()); // 42
```

##### `keys()`

Returns all keys in the cache.

**Returns:** `string[]` - Array of keys

**Example:**
```javascript
const keys = cache.keys();
```

##### `values()`

Returns all values in the cache.

**Returns:** `any[]` - Array of values

**Example:**
```javascript
const values = cache.values();
```

##### `entries()`

Returns all key-value pairs in the cache.

**Returns:** `Array<[string, any]>` - Array of entries

**Example:**
```javascript
for (const [key, value] of cache.entries()) {
    console.log(key, value);
}
```

#### `class RedisCache`

Redis-backed cache implementation.

**Constructor:**
```javascript
const cache = new RedisCache(options);
```

**Options:**
| Name | Type | Description |
|------|------|-------------|
| `host` | string | Redis host |
| `port` | number | Redis port |
| `password` | string | Redis password |
| `ttl` | number | Default TTL |

**Methods:** Same as `Cache` class but async.

**Example:**
```javascript
const cache = new RedisCache({ host: 'localhost', port: 6379 });
await cache.set('user:1', { name: 'John' });
const user = await cache.get('user:1');
```

---

## 📊 Type Index

| Class/Type | Module | Description |
|------------|--------|-------------|
| `Cache` | cache | In-memory cache |
| `RedisCache` | cache | Redis cache client |
| `Stats` | file | File statistics |
| `ValidationError` | validator | Validation error |

## 🔧 Error Types

| Error | Module | Description |
|-------|--------|-------------|
| `FileNotFoundError` | file | File not found |
| `PermissionError` | file | Permission denied |
| `ValidationError` | validator | Validation failed |
| `TimeoutError` | async | Operation timed out |
| `CacheError` | cache | Cache operation failed |

## 📈 Performance Characteristics

| Function | Time Complexity | Space Complexity |
|----------|----------------|------------------|
| `math.add` | O(1) | O(1) |
| `math.divide` | O(1) | O(1) |
| `string.reverse` | O(n) | O(n) |
| `file.readFile` | O(n) | O(n) |
| `validator.isEmail` | O(n) | O(1) |
| `cache.get` | O(1) | O(1) |
| `async.parallel` | O(n) | O(n) |

## 🧪 Examples

See the [examples](../examples/) directory for complete runnable examples.

---

*Last Updated: 2024*
```

## ✅ **What this API documentation provides:**

| Section | Description |
|---------|-------------|
| **Module Overview** | Description of each module's purpose |
| **Installation** | Import statements for each module |
| **Functions** | Complete function signatures with parameters and return values |
| **Classes** | Class definitions with constructor options and methods |
| **Parameters** | Detailed parameter tables with types and descriptions |
| **Return Values** | Clear return type documentation |
| **Throws** | Error conditions documented |
| **Examples** | Code examples for each function |
| **Type Index** | Quick reference of types by module |
| **Error Types** | Error classes and their meanings |
| **Performance** | Time and space complexity notes |

## 🎯 **Purpose as Test Data**

This file serves as an **expected output** for validating that your readme generator produces:
- ✅ Proper JSDoc-style documentation
- ✅ Comprehensive module documentation
- ✅ Clear parameter and return type documentation
- ✅ Error condition documentation
- ✅ Code examples
- ✅ Consistent formatting across modules
- ✅ Performance characteristics