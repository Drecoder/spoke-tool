/**
 * String Operations Tests
 * 
 * Tests for string manipulation functions including case conversion,
 * trimming, splitting, joining, searching, and Unicode handling.
 */

const {
    // Basic operations
    reverse,
    length,
    charAt,
    charCodeAt,
    
    // Case conversion
    toUpper,
    toLower,
    toTitleCase,
    toCamelCase,
    toPascalCase,
    toSnakeCase,
    toKebabCase,
    
    // Trimming
    trim,
    trimStart,
    trimEnd,
    
    // Splitting and joining
    split,
    join,
    chars,
    words,
    
    // Searching
    contains,
    indexOf,
    lastIndexOf,
    startsWith,
    endsWith,
    count,
    
    // Replacement
    replace,
    replaceAll,
    remove,
    
    // Substring extraction
    substring,
    slice,
    substr,
    
    // Padding
    padStart,
    padEnd,
    
    // Repetition
    repeat,
    
    // Utilities
    truncate,
    capitalize,
    camelCase,
    snakeCase,
    kebabCase,
    pluralize,
    singularize,
    slugify,
    ellipsis,
    
    // Validation
    isEmpty,
    isBlank,
    isNumeric,
    isAlpha,
    isAlphanumeric,
    isEmail,
    isURL,
    isPhoneNumber,
    isPostalCode,
    
    // Unicode
    normalize,
    toAscii,
    countGraphemes,
    
    // Encoding
    toBase64,
    fromBase64,
    toHex,
    fromHex,
    toBinary,
    fromBinary,
    
    // Comparison
    compare,
    equalsIgnoreCase,
    levenshteinDistance,
    
    // Advanced
    template,
    interpolate,
    mask,
    obfuscate,
    highlight
} = require('../src/string');

// ============================================================================
// Basic Operations Tests
// ============================================================================

describe('Basic Operations', () => {
    
    describe('reverse', () => {
        test('should reverse normal strings', () => {
            expect(reverse('hello')).toBe('olleh');
            expect(reverse('world')).toBe('dlrow');
            expect(reverse('javascript')).toBe('tpircsavaj');
        });

        test('should handle empty string', () => {
            expect(reverse('')).toBe('');
        });

        test('should handle single character', () => {
            expect(reverse('a')).toBe('a');
        });

        test('should handle palindromes', () => {
            expect(reverse('racecar')).toBe('racecar');
            expect(reverse('madam')).toBe('madam');
        });

        test('should handle strings with spaces', () => {
            expect(reverse('hello world')).toBe('dlrow olleh');
        });

        test('should handle strings with punctuation', () => {
            expect(reverse('hello!')).toBe('!olleh');
            expect(reverse('a,b,c')).toBe('c,b,a');
        });

        test('should handle Unicode characters', () => {
            expect(reverse('café')).toBe('éfac');
            expect(reverse('résumé')).toBe('émusér');
            expect(reverse('Hello 世界')).toBe('界世 olleH');
        });

        test('should handle emoji', () => {
            expect(reverse('hello 👋')).toBe('👋 olleh');
            expect(reverse('🚀✨🌟')).toBe('🌟✨🚀');
        });

        test('should handle surrogate pairs', () => {
            expect(reverse('𠜎𠜱𠝹𠱓')).toBe('𠱓𠝹𠜱𠜎');
        });
    });

    describe('length', () => {
        test('should return correct length', () => {
            expect(length('hello')).toBe(5);
            expect(length('')).toBe(0);
            expect(length('hello world')).toBe(11);
        });

        test('should handle Unicode characters', () => {
            expect(length('café')).toBe(4);
            expect(length('Hello 世界')).toBe(8);
            expect(length('👨‍👩‍👧‍👦')).toBe(1); // Single grapheme
        });
    });

    describe('charAt', () => {
        test('should return character at index', () => {
            expect(charAt('hello', 0)).toBe('h');
            expect(charAt('hello', 2)).toBe('l');
            expect(charAt('hello', 4)).toBe('o');
        });

        test('should handle out of bounds', () => {
            expect(charAt('hello', -1)).toBe('');
            expect(charAt('hello', 5)).toBe('');
            expect(charAt('hello', 10)).toBe('');
        });

        test('should handle Unicode', () => {
            expect(charAt('café', 3)).toBe('é');
        });
    });

    describe('charCodeAt', () => {
        test('should return character code', () => {
            expect(charCodeAt('A', 0)).toBe(65);
            expect(charCodeAt('a', 0)).toBe(97);
            expect(charCodeAt('0', 0)).toBe(48);
        });

        test('should handle out of bounds', () => {
            expect(charCodeAt('hello', -1)).toBe(NaN);
            expect(charCodeAt('hello', 5)).toBe(NaN);
        });

        test('should handle Unicode', () => {
            expect(charCodeAt('é', 0)).toBe(233);
        });
    });
});

// ============================================================================
// Case Conversion Tests
// ============================================================================

describe('Case Conversion', () => {
    
    describe('toUpper', () => {
        test('should convert to uppercase', () => {
            expect(toUpper('hello')).toBe('HELLO');
            expect(toUpper('Hello World')).toBe('HELLO WORLD');
            expect(toUpper('123abc')).toBe('123ABC');
        });

        test('should handle empty string', () => {
            expect(toUpper('')).toBe('');
        });

        test('should handle Unicode', () => {
            expect(toUpper('café')).toBe('CAFÉ');
            expect(toUpper('straße')).toBe('STRASSE'); // German sharp s
        });
    });

    describe('toLower', () => {
        test('should convert to lowercase', () => {
            expect(toLower('HELLO')).toBe('hello');
            expect(toLower('Hello World')).toBe('hello world');
            expect(toLower('123ABC')).toBe('123abc');
        });

        test('should handle empty string', () => {
            expect(toLower('')).toBe('');
        });

        test('should handle Unicode', () => {
            expect(toLower('CAFÉ')).toBe('café');
            expect(toLower('İstanbul')).toBe('i̇stanbul'); // Turkish dotless i
        });
    });

    describe('toTitleCase', () => {
        test('should convert to title case', () => {
            expect(toTitleCase('hello world')).toBe('Hello World');
            expect(toTitleCase('the quick brown fox')).toBe('The Quick Brown Fox');
        });

        test('should handle mixed case', () => {
            expect(toTitleCase('hELLo wORLD')).toBe('Hello World');
        });

        test('should handle articles correctly', () => {
            expect(toTitleCase('the lord of the rings')).toBe('The Lord of the Rings');
            expect(toTitleCase('war and peace')).toBe('War and Peace');
        });

        test('should handle single word', () => {
            expect(toTitleCase('hello')).toBe('Hello');
        });
    });

    describe('toCamelCase', () => {
        test('should convert to camelCase', () => {
            expect(toCamelCase('hello world')).toBe('helloWorld');
            expect(toCamelCase('hello-world')).toBe('helloWorld');
            expect(toCamelCase('hello_world')).toBe('helloWorld');
        });

        test('should handle single word', () => {
            expect(toCamelCase('hello')).toBe('hello');
        });

        test('should handle empty string', () => {
            expect(toCamelCase('')).toBe('');
        });
    });

    describe('toPascalCase', () => {
        test('should convert to PascalCase', () => {
            expect(toPascalCase('hello world')).toBe('HelloWorld');
            expect(toPascalCase('hello-world')).toBe('HelloWorld');
            expect(toPascalCase('hello_world')).toBe('HelloWorld');
        });

        test('should handle single word', () => {
            expect(toPascalCase('hello')).toBe('Hello');
        });
    });

    describe('toSnakeCase', () => {
        test('should convert to snake_case', () => {
            expect(toSnakeCase('helloWorld')).toBe('hello_world');
            expect(toSnakeCase('HelloWorld')).toBe('hello_world');
            expect(toSnakeCase('hello-world')).toBe('hello_world');
        });

        test('should handle single word', () => {
            expect(toSnakeCase('hello')).toBe('hello');
        });

        test('should handle numbers', () => {
            expect(toSnakeCase('hello123World')).toBe('hello123_world');
        });
    });

    describe('toKebabCase', () => {
        test('should convert to kebab-case', () => {
            expect(toKebabCase('helloWorld')).toBe('hello-world');
            expect(toKebabCase('HelloWorld')).toBe('hello-world');
            expect(toKebabCase('hello_world')).toBe('hello-world');
        });

        test('should handle single word', () => {
            expect(toKebabCase('hello')).toBe('hello');
        });
    });
});

// ============================================================================
// Trimming Tests
// ============================================================================

describe('Trimming', () => {
    
    describe('trim', () => {
        test('should trim whitespace from both ends', () => {
            expect(trim('  hello  ')).toBe('hello');
            expect(trim('\t\n hello \n\t')).toBe('hello');
        });

        test('should handle empty string', () => {
            expect(trim('')).toBe('');
        });

        test('should handle string with only whitespace', () => {
            expect(trim('   ')).toBe('');
            expect(trim('\t\n')).toBe('');
        });

        test('should preserve internal spaces', () => {
            expect(trim('  hello world  ')).toBe('hello world');
        });
    });

    describe('trimStart', () => {
        test('should trim leading whitespace', () => {
            expect(trimStart('  hello')).toBe('hello');
            expect(trimStart('\t\nhello')).toBe('hello');
        });

        test('should preserve trailing spaces', () => {
            expect(trimStart('  hello  ')).toBe('hello  ');
        });
    });

    describe('trimEnd', () => {
        test('should trim trailing whitespace', () => {
            expect(trimEnd('hello  ')).toBe('hello');
            expect(trimEnd('hello\n\t')).toBe('hello');
        });

        test('should preserve leading spaces', () => {
            expect(trimEnd('  hello  ')).toBe('  hello');
        });
    });
});

// ============================================================================
// Splitting and Joining Tests
// ============================================================================

describe('Splitting and Joining', () => {
    
    describe('split', () => {
        test('should split by delimiter', () => {
            expect(split('a,b,c', ',')).toEqual(['a', 'b', 'c']);
            expect(split('hello world', ' ')).toEqual(['hello', 'world']);
        });

        test('should handle multiple delimiters', () => {
            expect(split('a,,b,c', ',')).toEqual(['a', '', 'b', 'c']);
        });

        test('should handle delimiter at start', () => {
            expect(split(',a,b', ',')).toEqual(['', 'a', 'b']);
        });

        test('should handle delimiter at end', () => {
            expect(split('a,b,', ',')).toEqual(['a', 'b', '']);
        });

        test('should handle empty string', () => {
            expect(split('', ',')).toEqual(['']);
        });

        test('should handle no delimiter match', () => {
            expect(split('hello', ',')).toEqual(['hello']);
        });

        test('should handle regex delimiter', () => {
            expect(split('hello123world', /\d+/)).toEqual(['hello', 'world']);
            expect(split('a,b;c', /[,;]/)).toEqual(['a', 'b', 'c']);
        });

        test('should handle limit parameter', () => {
            expect(split('a,b,c,d', ',', 2)).toEqual(['a', 'b']);
            expect(split('a,b,c,d', ',', 0)).toEqual([]);
        });
    });

    describe('join', () => {
        test('should join array with delimiter', () => {
            expect(join(['a', 'b', 'c'], ',')).toBe('a,b,c');
            expect(join(['hello', 'world'], ' ')).toBe('hello world');
        });

        test('should handle empty array', () => {
            expect(join([], ',')).toBe('');
        });

        test('should handle single element', () => {
            expect(join(['hello'], ',')).toBe('hello');
        });

        test('should handle empty strings', () => {
            expect(join(['a', '', 'c'], ',')).toBe('a,,c');
        });

        test('should handle different separators', () => {
            expect(join(['a', 'b', 'c'], '-')).toBe('a-b-c');
            expect(join(['a', 'b', 'c'], '')).toBe('abc');
        });
    });

    describe('chars', () => {
        test('should split into characters', () => {
            expect(chars('hello')).toEqual(['h', 'e', 'l', 'l', 'o']);
        });

        test('should handle empty string', () => {
            expect(chars('')).toEqual([]);
        });

        test('should handle Unicode', () => {
            expect(chars('café')).toEqual(['c', 'a', 'f', 'é']);
        });

        test('should handle emoji', () => {
            expect(chars('🚀✨')).toEqual(['🚀', '✨']);
        });
    });

    describe('words', () => {
        test('should split into words', () => {
            expect(words('hello world')).toEqual(['hello', 'world']);
            expect(words('the quick brown fox')).toEqual(['the', 'quick', 'brown', 'fox']);
        });

        test('should handle punctuation', () => {
            expect(words('hello, world!')).toEqual(['hello', 'world']);
            expect(words('one.two.three')).toEqual(['one', 'two', 'three']);
        });

        test('should handle multiple spaces', () => {
            expect(words('hello   world')).toEqual(['hello', 'world']);
        });

        test('should handle empty string', () => {
            expect(words('')).toEqual([]);
        });
    });
});

// ============================================================================
// Searching Tests
// ============================================================================

describe('Searching', () => {
    
    describe('contains', () => {
        test('should find substring', () => {
            expect(contains('hello world', 'world')).toBe(true);
            expect(contains('hello world', 'xyz')).toBe(false);
        });

        test('should handle empty substring', () => {
            expect(contains('hello', '')).toBe(true);
        });

        test('should handle case sensitivity', () => {
            expect(contains('Hello World', 'hello')).toBe(false);
        });

        test('should support case-insensitive mode', () => {
            expect(contains('Hello World', 'hello', true)).toBe(true);
        });

        test('should handle start position', () => {
            expect(contains('hello world', 'world', 6)).toBe(true);
            expect(contains('hello world', 'world', 7)).toBe(false);
        });
    });

    describe('indexOf', () => {
        test('should find first occurrence', () => {
            expect(indexOf('hello world', 'world')).toBe(6);
            expect(indexOf('hello world', 'o')).toBe(4);
        });

        test('should return -1 if not found', () => {
            expect(indexOf('hello world', 'xyz')).toBe(-1);
        });

        test('should handle empty substring', () => {
            expect(indexOf('hello', '')).toBe(0);
        });

        test('should handle start position', () => {
            expect(indexOf('hello world', 'o', 5)).toBe(7);
        });
    });

    describe('lastIndexOf', () => {
        test('should find last occurrence', () => {
            expect(lastIndexOf('hello world hello', 'hello')).toBe(12);
            expect(lastIndexOf('hello world', 'o')).toBe(7);
        });

        test('should return -1 if not found', () => {
            expect(lastIndexOf('hello world', 'xyz')).toBe(-1);
        });

        test('should handle empty substring', () => {
            expect(lastIndexOf('hello', '')).toBe(5);
        });
    });

    describe('startsWith', () => {
        test('should check prefix', () => {
            expect(startsWith('hello world', 'hello')).toBe(true);
            expect(startsWith('hello world', 'world')).toBe(false);
        });

        test('should handle empty prefix', () => {
            expect(startsWith('hello', '')).toBe(true);
        });

        test('should handle position', () => {
            expect(startsWith('hello world', 'world', 6)).toBe(true);
        });
    });

    describe('endsWith', () => {
        test('should check suffix', () => {
            expect(endsWith('hello world', 'world')).toBe(true);
            expect(endsWith('hello world', 'hello')).toBe(false);
        });

        test('should handle empty suffix', () => {
            expect(endsWith('hello', '')).toBe(true);
        });

        test('should handle position', () => {
            expect(endsWith('hello world', 'hello', 5)).toBe(true);
        });
    });

    describe('count', () => {
        test('should count occurrences', () => {
            expect(count('hello world hello', 'hello')).toBe(2);
            expect(count('hello world', 'o')).toBe(2);
        });

        test('should handle overlapping', () => {
            expect(count('aaaaa', 'aa')).toBe(2); // Non-overlapping
        });

        test('should handle empty substring', () => {
            expect(count('hello', '')).toBe(6); // Length + 1
        });

        test('should handle not found', () => {
            expect(count('hello', 'xyz')).toBe(0);
        });
    });
});

// ============================================================================
// Replacement Tests
// ============================================================================

describe('Replacement', () => {
    
    describe('replace', () => {
        test('should replace first occurrence', () => {
            expect(replace('hello world hello', 'hello', 'hi')).toBe('hi world hello');
        });

        test('should handle no matches', () => {
            expect(replace('hello world', 'xyz', 'abc')).toBe('hello world');
        });

        test('should handle regex pattern', () => {
            expect(replace('hello123world', /\d+/, '!')).toBe('hello!world');
        });

        test('should handle function replacement', () => {
            expect(replace('hello123', /\d+/, (match) => parseInt(match) * 2)).toBe('hello246');
        });
    });

    describe('replaceAll', () => {
        test('should replace all occurrences', () => {
            expect(replaceAll('hello world hello', 'hello', 'hi')).toBe('hi world hi');
            expect(replaceAll('aaa', 'a', 'b')).toBe('bbb');
        });

        test('should handle no matches', () => {
            expect(replaceAll('hello world', 'xyz', 'abc')).toBe('hello world');
        });

        test('should handle regex with global flag', () => {
            expect(replaceAll('hello123world456', /\d+/g, '!')).toBe('hello!world!');
        });

        test('should handle special characters', () => {
            expect(replaceAll('a.b.c.d', '.', '-')).toBe('a-b-c-d');
        });
    });

    describe('remove', () => {
        test('should remove substring', () => {
            expect(remove('hello world hello', 'hello')).toBe(' world ');
            expect(remove('aaa', 'a')).toBe('');
        });

        test('should handle multiple removals', () => {
            expect(remove('hello world', 'o')).toBe('hell wrld');
        });

        test('should handle no matches', () => {
            expect(remove('hello world', 'xyz')).toBe('hello world');
        });
    });
});

// ============================================================================
// Substring Extraction Tests
// ============================================================================

describe('Substring Extraction', () => {
    
    describe('substring', () => {
        test('should extract substring', () => {
            expect(substring('hello world', 0, 5)).toBe('hello');
            expect(substring('hello world', 6, 11)).toBe('world');
        });

        test('should handle negative indices', () => {
            expect(substring('hello world', -5, 5)).toBe('hello');
        });

        test('should handle swapped indices', () => {
            expect(substring('hello world', 5, 0)).toBe('hello');
        });

        test('should handle single index', () => {
            expect(substring('hello world', 6)).toBe('world');
        });
    });

    describe('slice', () => {
        test('should extract slice', () => {
            expect(slice('hello world', 0, 5)).toBe('hello');
            expect(slice('hello world', 6, 11)).toBe('world');
        });

        test('should handle negative indices', () => {
            expect(slice('hello world', -5)).toBe('world');
            expect(slice('hello world', -5, -1)).toBe('worl');
        });

        test('should handle out of bounds', () => {
            expect(slice('hello', 10)).toBe('');
        });
    });

    describe('substr', () => {
        test('should extract substring by length', () => {
            expect(substr('hello world', 0, 5)).toBe('hello');
            expect(substr('hello world', 6, 5)).toBe('world');
        });

        test('should handle negative start', () => {
            expect(substr('hello world', -5)).toBe('world');
        });

        test('should handle length exceeding string', () => {
            expect(substr('hello', 2, 10)).toBe('llo');
        });
    });
});

// ============================================================================
// Padding Tests
// ============================================================================

describe('Padding', () => {
    
    describe('padStart', () => {
        test('should pad start', () => {
            expect(padStart('5', 3, '0')).toBe('005');
            expect(padStart('abc', 5, '-')).toBe('--abc');
        });

        test('should handle string longer than target', () => {
            expect(padStart('hello', 3, '-')).toBe('hello');
        });

        test('should handle empty pad string', () => {
            expect(padStart('hello', 10, '')).toBe('hello');
        });

        test('should handle default pad string', () => {
            expect(padStart('hello', 10)).toBe('     hello');
        });
    });

    describe('padEnd', () => {
        test('should pad end', () => {
            expect(padEnd('5', 3, '0')).toBe('500');
            expect(padEnd('abc', 5, '-')).toBe('abc--');
        });

        test('should handle string longer than target', () => {
            expect(padEnd('hello', 3, '-')).toBe('hello');
        });

        test('should handle default pad string', () => {
            expect(padEnd('hello', 10)).toBe('hello     ');
        });
    });
});

// ============================================================================
// Repetition Tests
// ============================================================================

describe('Repetition', () => {
    
    describe('repeat', () => {
        test('should repeat string', () => {
            expect(repeat('ha', 3)).toBe('hahaha');
            expect(repeat('abc', 2)).toBe('abcabc');
        });

        test('should handle zero count', () => {
            expect(repeat('hello', 0)).toBe('');
        });

        test('should handle negative count', () => {
            expect(() => repeat('hello', -1)).toThrow();
        });

        test('should handle empty string', () => {
            expect(repeat('', 5)).toBe('');
        });

        test('should handle count 1', () => {
            expect(repeat('hello', 1)).toBe('hello');
        });
    });
});

// ============================================================================
// Utility Functions Tests
// ============================================================================

describe('Utilities', () => {
    
    describe('truncate', () => {
        test('should truncate long strings', () => {
            expect(truncate('This is a long string', 10)).toBe('This is...');
            expect(truncate('Hello world', 5)).toBe('Hello...');
        });

        test('should not truncate short strings', () => {
            expect(truncate('Hello', 10)).toBe('Hello');
        });

        test('should handle custom ellipsis', () => {
            expect(truncate('Hello world', 5, '***')).toBe('Hello***');
        });

        test('should handle empty string', () => {
            expect(truncate('', 5)).toBe('');
        });
    });

    describe('capitalize', () => {
        test('should capitalize first letter', () => {
            expect(capitalize('hello')).toBe('Hello');
            expect(capitalize('hello world')).toBe('Hello world');
        });

        test('should handle empty string', () => {
            expect(capitalize('')).toBe('');
        });

        test('should handle already capitalized', () => {
            expect(capitalize('Hello')).toBe('Hello');
        });
    });

    describe('pluralize', () => {
        test('should pluralize regular words', () => {
            expect(pluralize('cat')).toBe('cats');
            expect(pluralize('dog')).toBe('dogs');
        });

        test('should handle irregular words', () => {
            expect(pluralize('child')).toBe('children');
            expect(pluralize('person')).toBe('people');
            expect(pluralize('mouse')).toBe('mice');
            expect(pluralize('ox')).toBe('oxen');
        });

        test('should handle words ending in y', () => {
            expect(pluralize('city')).toBe('cities');
            expect(pluralize('key')).toBe('keys'); // vowel + y
        });

        test('should handle words ending in s, sh, ch', () => {
            expect(pluralize('bus')).toBe('buses');
            expect(pluralize('brush')).toBe('brushes');
            expect(pluralize('match')).toBe('matches');
        });

        test('should handle count', () => {
            expect(pluralize('cat', 1)).toBe('cat');
            expect(pluralize('cat', 2)).toBe('cats');
            expect(pluralize('cat', 0)).toBe('cats');
        });
    });

    describe('singularize', () => {
        test('should singularize regular words', () => {
            expect(singularize('cats')).toBe('cat');
            expect(singularize('dogs')).toBe('dog');
        });

        test('should handle irregular words', () => {
            expect(singularize('children')).toBe('child');
            expect(singularize('people')).toBe('person');
            expect(singularize('mice')).toBe('mouse');
        });
    });

    describe('slugify', () => {
        test('should create URL-friendly slugs', () => {
            expect(slugify('Hello World')).toBe('hello-world');
            expect(slugify('This is a test')).toBe('this-is-a-test');
        });

        test('should handle special characters', () => {
            expect(slugify('Hello, World!')).toBe('hello-world');
            expect(slugify('Café & Restaurant')).toBe('cafe-restaurant');
        });

        test('should handle multiple spaces', () => {
            expect(slugify('hello   world')).toBe('hello-world');
        });

        test('should handle accents', () => {
            expect(slugify('Café Français')).toBe('cafe-francais');
            expect(slugify('über cool')).toBe('uber-cool');
        });
    });

    describe('ellipsis', () => {
        test('should add ellipsis to long strings', () => {
            expect(ellipsis('This is a long string', 10)).toBe('This is...');
        });

        test('should handle custom position', () => {
            expect(ellipsis('This is a long string', 10, 'middle')).toBe('This...ing');
            expect(ellipsis('This is a long string', 10, 'start')).toBe('...g string');
        });
    });
});

// ============================================================================
// Validation Tests
// ============================================================================

describe('Validation', () => {
    
    describe('isEmpty', () => {
        test('should detect empty string', () => {
            expect(isEmpty('')).toBe(true);
            expect(isEmpty('hello')).toBe(false);
            expect(isEmpty('   ')).toBe(false);
        });
    });

    describe('isBlank', () => {
        test('should detect blank string', () => {
            expect(isBlank('')).toBe(true);
            expect(isBlank('   ')).toBe(true);
            expect(isBlank('\t\n')).toBe(true);
            expect(isBlank('hello')).toBe(false);
        });
    });

    describe('isNumeric', () => {
        test('should detect numeric strings', () => {
            expect(isNumeric('123')).toBe(true);
            expect(isNumeric('-123')).toBe(true);
            expect(isNumeric('123.456')).toBe(true);
            expect(isNumeric('abc')).toBe(false);
            expect(isNumeric('123abc')).toBe(false);
        });

        test('should handle empty string', () => {
            expect(isNumeric('')).toBe(false);
        });
    });

    describe('isAlpha', () => {
        test('should detect alphabetic strings', () => {
            expect(isAlpha('abc')).toBe(true);
            expect(isAlpha('ABC')).toBe(true);
            expect(isAlpha('abc123')).toBe(false);
            expect(isAlpha('')).toBe(false);
        });
    });

    describe('isAlphanumeric', () => {
        test('should detect alphanumeric strings', () => {
            expect(isAlphanumeric('abc123')).toBe(true);
            expect(isAlphanumeric('ABC123')).toBe(true);
            expect(isAlphanumeric('abc-123')).toBe(false);
            expect(isAlphanumeric('')).toBe(false);
        });
    });

    describe('isEmail', () => {
        test('should validate email addresses', () => {
            expect(isEmail('user@example.com')).toBe(true);
            expect(isEmail('user.name@example.co.uk')).toBe(true);
            expect(isEmail('user+tag@example.com')).toBe(true);
            expect(isEmail('user@example')).toBe(false);
            expect(isEmail('@example.com')).toBe(false);
            expect(isEmail('user@.com')).toBe(false);
            expect(isEmail('')).toBe(false);
        });
    });

    describe('isURL', () => {
        test('should validate URLs', () => {
            expect(isURL('https://example.com')).toBe(true);
            expect(isURL('http://example.com/path')).toBe(true);
            expect(isURL('https://sub.example.com')).toBe(true);
            expect(isURL('example.com')).toBe(false);
            expect(isURL('https://')).toBe(false);
            expect(isURL('')).toBe(false);
        });
    });

    describe('isPhoneNumber', () => {
        test('should validate phone numbers', () => {
            expect(isPhoneNumber('123-456-7890')).toBe(true);
            expect(isPhoneNumber('(123) 456-7890')).toBe(true);
            expect(isPhoneNumber('123.456.7890')).toBe(true);
            expect(isPhoneNumber('+1 123-456-7890')).toBe(true);
            expect(isPhoneNumber('1234567890')).toBe(true);
            expect(isPhoneNumber('abc')).toBe(false);
        });
    });

    describe('isPostalCode', () => {
        test('should validate postal codes', () => {
            expect(isPostalCode('12345')).toBe(true);
            expect(isPostalCode('12345-6789')).toBe(true);
            expect(isPostalCode('A1B 2C3')).toBe(true); // Canadian
            expect(isPostalCode('1234')).toBe(false);
            expect(isPostalCode('')).toBe(false);
        });
    });
});

// ============================================================================
// Unicode Tests
// ============================================================================

describe('Unicode Handling', () => {
    
    describe('normalize', () => {
        test('should normalize Unicode strings', () => {
            expect(normalize('café')).toBe('café');
            expect(normalize('cafe\u0301')).toBe('café'); // Combined form
        });

        test('should handle different forms', () => {
            expect(normalize('\u1E9B\u0323', 'NFC')).toBe('ẛ̣');
            expect(normalize('\u1E9B\u0323', 'NFD')).toBe('ẛ̣');
        });
    });

    describe('toAscii', () => {
        test('should convert to ASCII', () => {
            expect(toAscii('café')).toBe('cafe');
            expect(toAscii('über')).toBe('uber');
            expect(toAscii('Français')).toBe('Francais');
            expect(toAscii('你好')).toBe(''); // Non-ASCII removed
        });
    });

    describe('countGraphemes', () => {
        test('should count grapheme clusters', () => {
            expect(countGraphemes('hello')).toBe(5);
            expect(countGraphemes('café')).toBe(4);
            expect(countGraphemes('👨‍👩‍👧‍👦')).toBe(1); // Family emoji
            expect(countGraphemes('é')).toBe(1); // e + accent
        });
    });
});

// ============================================================================
// Encoding Tests
// ============================================================================

describe('Encoding', () => {
    
    describe('toBase64 / fromBase64', () => {
        test('should encode/decode Base64', () => {
            const encoded = toBase64('hello world');
            expect(encoded).toBe('aGVsbG8gd29ybGQ=');
            expect(fromBase64(encoded)).toBe('hello world');
        });

        test('should handle empty string', () => {
            expect(toBase64('')).toBe('');
            expect(fromBase64('')).toBe('');
        });

        test('should handle Unicode', () => {
            const encoded = toBase64('café');
            expect(fromBase64(encoded)).toBe('café');
        });
    });

    describe('toHex / fromHex', () => {
        test('should encode/decode hex', () => {
            expect(toHex('hello')).toBe('68656c6c6f');
            expect(fromHex('68656c6c6f')).toBe('hello');
        });

        test('should handle empty string', () => {
            expect(toHex('')).toBe('');
            expect(fromHex('')).toBe('');
        });

        test('should handle invalid hex', () => {
            expect(() => fromHex('xyz')).toThrow();
        });
    });

    describe('toBinary / fromBinary', () => {
        test('should encode/decode binary', () => {
            expect(toBinary('A')).toBe('01000001');
            expect(fromBinary('01000001')).toBe('A');
        });

        test('should handle empty string', () => {
            expect(toBinary('')).toBe('');
            expect(fromBinary('')).toBe('');
        });

        test('should handle multiple characters', () => {
            expect(toBinary('AB')).toBe('0100000101000010');
            expect(fromBinary('0100000101000010')).toBe('AB');
        });
    });
});

// ============================================================================
// Comparison Tests
// ============================================================================

describe('Comparison', () => {
    
    describe('compare', () => {
        test('should compare strings', () => {
            expect(compare('a', 'b')).toBe(-1);
            expect(compare('b', 'a')).toBe(1);
            expect(compare('a', 'a')).toBe(0);
        });

        test('should handle case sensitivity', () => {
            expect(compare('A', 'a')).toBe(1); // 'A' < 'a' in ASCII
        });

        test('should handle different lengths', () => {
            expect(compare('aa', 'a')).toBe(1);
        });
    });

    describe('equalsIgnoreCase', () => {
        test('should compare ignoring case', () => {
            expect(equalsIgnoreCase('hello', 'HELLO')).toBe(true);
            expect(equalsIgnoreCase('Hello', 'hELLO')).toBe(true);
            expect(equalsIgnoreCase('hello', 'world')).toBe(false);
        });

        test('should handle Unicode', () => {
            expect(equalsIgnoreCase('café', 'CAFÉ')).toBe(true);
        });
    });

    describe('levenshteinDistance', () => {
        test('should calculate edit distance', () => {
            expect(levenshteinDistance('kitten', 'sitting')).toBe(3);
            expect(levenshteinDistance('hello', 'hello')).toBe(0);
            expect(levenshteinDistance('', 'hello')).toBe(5);
            expect(levenshteinDistance('hello', '')).toBe(5);
        });

        test('should handle different lengths', () => {
            expect(levenshteinDistance('cat', 'cats')).toBe(1);
        });

        test('should handle case sensitivity', () => {
            expect(levenshteinDistance('Hello', 'hello')).toBe(1);
        });
    });
});

// ============================================================================
// Advanced Tests
// ============================================================================

describe('Advanced Functions', () => {
    
    describe('template', () => {
        test('should replace template variables', () => {
            expect(template('Hello {{name}}!', { name: 'World' })).toBe('Hello World!');
            expect(template('{{greet}} {{name}}', { greet: 'Hi', name: 'John' })).toBe('Hi John');
        });

        test('should handle missing variables', () => {
            expect(template('Hello {{name}}!', {})).toBe('Hello !');
        });

        test('should handle custom delimiters', () => {
            expect(template('Hello ${name}!', { name: 'World' }, '${', '}')).toBe('Hello World!');
        });

        test('should handle nested templates', () => {
            expect(template('{{greet}} {{user.name}}', { 
                greet: 'Hello', 
                user: { name: 'John' } 
            })).toBe('Hello John');
        });
    });

    describe('interpolate', () => {
        test('should interpolate values', () => {
            expect(interpolate('Hello {0}!', ['World'])).toBe('Hello World!');
            expect(interpolate('{0} {1}', ['Hello', 'World'])).toBe('Hello World');
        });

        test('should handle multiple occurrences', () => {
            expect(interpolate('{0} {0} {0}', ['echo'])).toBe('echo echo echo');
        });

        test('should handle missing values', () => {
            expect(interpolate('Hello {0} {1}!', ['World'])).toBe('Hello World !');
        });
    });

    describe('mask', () => {
        test('should mask parts of string', () => {
            expect(mask('1234-5678-9012-3456', 4)).toBe('************3456');
            expect(mask('password123', 0)).toBe('***********');
        });

        test('should handle custom mask character', () => {
            expect(mask('secret', 2, '#')).toBe('####et');
        });

        test('should handle different masking positions', () => {
            expect(mask('1234567890', 4, '*', 'start')).toBe('****567890');
            expect(mask('1234567890', 4, '*', 'end')).toBe('123456****');
        });
    });

    describe('obfuscate', () => {
        test('should obfuscate email', () => {
            expect(obfuscate('user@example.com')).toBe('u**r@e*****e.com');
        });

        test('should obfuscate phone number', () => {
            expect(obfuscate('123-456-7890')).toBe('***-***-7890');
        });

        test('should handle short strings', () => {
            expect(obfuscate('a@b.c')).toBe('*@b.c');
        });
    });

    describe('highlight', () => {
        test('should highlight matches', () => {
            expect(highlight('hello world hello', 'hello')).toBe('**hello** world **hello**');
        });

        test('should handle custom markers', () => {
            expect(highlight('hello world', 'hello', '<mark>', '</mark>')).toBe('<mark>hello</mark> world');
        });

        test('should handle case-insensitive', () => {
            expect(highlight('Hello World', 'hello', '**', '**', true)).toBe('**Hello** World');
        });
    });
});

// ============================================================================
// Edge Cases and Error Handling
// ============================================================================

describe('Edge Cases', () => {
    
    test('should handle null input', () => {
        expect(() => reverse(null)).toThrow();
        expect(() => toUpper(null)).toThrow();
        expect(() => split(null, ',')).toThrow();
    });

    test('should handle undefined input', () => {
        expect(() => reverse(undefined)).toThrow();
        expect(() => toUpper(undefined)).toThrow();
    });

    test('should handle non-string input', () => {
        expect(() => reverse(123)).toThrow();
        expect(() => toUpper({})).toThrow();
        expect(() => split([], ',')).toThrow();
    });

    test('should handle very long strings', () => {
        const longString = 'a'.repeat(1000000);
        expect(reverse(longString)).toBeDefined();
        expect(toUpper(longString)).toBeDefined();
    });

    test('should handle strings with control characters', () => {
        const str = 'hello\x00world';
        expect(reverse(str)).toBe('dlrow\x00olleh');
        expect(contains(str, '\x00')).toBe(true);
    });

    test('should handle strings with emoji', () => {
        const str = 'Hello 👋 World 🌍';
        expect(reverse(str)).toBe('🌍 dlroW 👋 olleH');
        expect(length(str)).toBe(13); // Spaces count
    });
});

// ============================================================================
// Performance Tests
// ============================================================================

describe('Performance', () => {
    
    test('reverse performance', () => {
        const str = 'a'.repeat(10000);
        const start = performance.now();
        reverse(str);
        const duration = performance.now() - start;
        expect(duration).toBeLessThan(50);
    });

    test('toUpper performance', () => {
        const str = 'a'.repeat(10000);
        const start = performance.now();
        toUpper(str);
        const duration = performance.now() - start;
        expect(duration).toBeLessThan(20);
    });

    test('split/join performance', () => {
        const str = 'a,b,'.repeat(1000);
        const start = performance.now();
        const arr = split(str, ',');
        join(arr, ',');
        const duration = performance.now() - start;
        expect(duration).toBeLessThan(100);
    });

    test('regex operations performance', () => {
        const str = 'a'.repeat(10000);
        const start = performance.now();
        replaceAll(str, 'a', 'b');
        const duration = performance.now() - start;
        expect(duration).toBeLessThan(100);
    });
});