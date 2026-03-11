```markdown
# Node.js Utility Library

A comprehensive, high-performance utility library for Node.js applications with zero dependencies.

[![npm version](https://img.shields.io/npm/v/node-utils.svg)](https://www.npmjs.com/package/node-utils)
[![License](https://img.shields.io/npm/l/node-utils.svg)](https://github.com/username/node-utils/blob/main/LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/username/node-utils/ci.yml)](https://github.com/username/node-utils/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/username/node-utils)](https://codecov.io/gh/username/node-utils)
[![npm downloads](https://img.shields.io/npm/dm/node-utils.svg)](https://www.npmjs.com/package/node-utils)

## Features

- ✅ **Math utilities** - Advanced mathematical operations with precision control
- ✅ **String manipulation** - Comprehensive string processing functions
- ✅ **File system helpers** - Safe and efficient file operations
- ✅ **Data validation** - Robust input validation utilities
- ✅ **Type conversion** - Reliable data type converters
- ✅ **Async utilities** - Promise helpers and control flow
- ✅ **Caching mechanisms** - In-memory and Redis cache implementations
- ✅ **Zero dependencies** - Lightweight and secure
- ✅ **100% test coverage** - Fully tested with Jest
- ✅ **TypeScript support** - Includes type definitions

## Installation

```bash
# Using npm
npm install node-utils

# Using yarn
yarn add node-utils

# Using pnpm
pnpm add node-utils
```

## Quick Start

```javascript
// CommonJS
const { math, string, file, validator } = require('node-utils');

// ES Modules
import { math, string, file, validator } from 'node-utils';

// Math operations
const sum = math.add(5.2, 3.1);
console.log('5.2 + 3.1 =', sum); // 8.3

const quotient = math.divide(10, 2);
console.log('10 / 2 =', quotient); // 5

// String manipulation
const reversed = string.reverse('hello');
console.log('hello reversed =', reversed); // 'olleh'

const camelCase = string.toCamelCase('user_name');
console.log('user_name in camelCase =', camelCase); // 'userName'

// File operations
await file.writeFile('./output.txt', 'Hello World');
const content = await file.readFile('./output.txt');
console.log('File content:', content);

// Validation
const isValidEmail = validator.isEmail('user@example.com');
console.log('Is valid email?', isValidEmail); // true
```

## 📦 Modules

The library is organized into the following modules:

| Module | Description | Documentation |
|--------|-------------|---------------|
| **math** | Mathematical operations | [API Docs](./API.md#module-math) |
| **string** | String manipulation | [API Docs](./API.md#module-string) |
| **file** | File system operations | [API Docs](./API.md#module-file) |
| **validator** | Input validation | [API Docs](./API.md#module-validator) |
| **converter** | Data conversion | [API Docs](./API.md#module-converter) |
| **async** | Asynchronous utilities | [API Docs](./API.md#module-async) |
| **cache** | Caching mechanisms | [API Docs](./API.md#module-cache) |

## 📚 Detailed Usage

### Math Module

```javascript
const { math } = require('node-utils');

// Basic arithmetic
console.log(math.add(5, 3));        // 8
console.log(math.subtract(10, 4));  // 6
console.log(math.multiply(6, 7));   // 42
console.log(math.divide(10, 2));    // 5

// Advanced math
console.log(math.power(2, 8));      // 256
console.log(math.sqrt(16));         // 4
console.log(math.abs(-42));         // 42
console.log(math.round(3.14159, 2)); // 3.14

// Statistics
console.log(math.max(1, 5, 3, 9, 2));    // 9
console.log(math.min(1, 5, 3, 9, 2));    // 1
console.log(math.sum([1, 2, 3, 4, 5]));  // 15
console.log(math.average([1, 2, 3, 4])); // 2.5

// Error handling
try {
    math.divide(10, 0);
} catch (err) {
    console.error(err.message); // 'Division by zero'
}
```

### String Module

```javascript
const { string } = require('node-utils');

// Basic string operations
console.log(string.reverse('hello'));           // 'olleh'
console.log(string.toUpper('hello'));           // 'HELLO'
console.log(string.toLower('HELLO'));           // 'hello'
console.log(string.trim('  hello  '));          // 'hello'

// Searching
console.log(string.contains('hello world', 'world')); // true
console.log(string.indexOf('hello world', 'world'));  // 6
console.log(string.count('hello hello', 'hello'));    // 2

// Case conversion
console.log(string.toCamelCase('user_name'));     // 'userName'
console.log(string.toPascalCase('user_name'));    // 'UserName'
console.log(string.toSnakeCase('userName'));      // 'user_name'

// Manipulation
console.log(string.repeat('ha', 3));              // 'hahaha'
console.log(string.truncate('Long string', 8));   // 'Long...'
console.log(string.padStart('5', 3, '0'));        // '005'
console.log(string.padEnd('5', 3, '0'));          // '500'

// Splitting and joining
const parts = string.split('a,b,c', ',');
console.log(parts);                               // ['a', 'b', 'c']
console.log(string.join(['a', 'b', 'c'], '-'));   // 'a-b-c'
```

### File Module

```javascript
const { file } = require('node-utils');

// Check if file exists
const exists = await file.exists('./config.json');
console.log('File exists:', exists);

// Read file
try {
    const content = await file.readFile('./data.txt', 'utf8');
    console.log('File content:', content);
} catch (err) {
    console.error('Error reading file:', err.message);
}

// Write file
await file.writeFile('./output.txt', 'Hello World');

// Copy file
await file.copyFile('./source.txt', './backup/source.txt');

// Directory operations
await file.mkdir('./data/nested', { recursive: true });
const files = await file.readdir('./data');
console.log('Directory contents:', files);

// File information
const stats = await file.stat('./data.txt');
console.log('File size:', stats.size);
console.log('Is file:', stats.isFile());
console.log('Is directory:', stats.isDirectory());

// Delete files
await file.unlink('./temp.txt');
await file.rmdir('./temp', { recursive: true });

// Synchronous versions
const syncContent = file.readFileSync('./data.txt');
file.writeFileSync('./output.txt', 'Hello World');
```

### Validator Module

```javascript
const { validator } = require('node-utils');

// Email validation
console.log(validator.isEmail('user@example.com'));     // true
console.log(validator.isEmail('invalid'));              // false

// URL validation
console.log(validator.isURL('https://example.com'));    // true
console.log(validator.isURL('not-a-url'));              // false

// Phone number
console.log(validator.isPhoneNumber('555-123-4567'));   // true

// ZIP code
console.log(validator.isZipCode('12345'));              // true
console.log(validator.isZipCode('12345-6789'));         // true

// Credit card
console.log(validator.isCreditCard('4111111111111111')); // true

// IP address
console.log(validator.isIP('192.168.1.1'));              // true
console.log(validator.isIP('2001:db8::1'));              // true

// UUID
const uuid = '123e4567-e89b-12d3-a456-426614174000';
console.log(validator.isUUID(uuid));                     // true

// String validation
console.log(validator.isAlphanumeric('abc123'));         // true
console.log(validator.isNumeric('12345'));               // true
console.log(validator.isAlpha('abc'));                   // true

// Length validation
console.log(validator.isLength('hello', 3, 10));         // true

// Range validation
console.log(validator.isInRange(5, 1, 10));              // true

// Required fields
console.log(validator.isRequired('hello'));              // true
console.log(validator.isRequired(''));                    // false
```

### Converter Module

```javascript
const { converter } = require('node-utils');

// Number conversion
console.log(converter.toNumber('42'));           // 42
console.log(converter.toInt('3.14'));            // 3
console.log(converter.toFloat('3.14'));          // 3.14

// Boolean conversion
console.log(converter.toBoolean('true'));        // true
console.log(converter.toBoolean('false'));       // false
console.log(converter.toBoolean('1'));           // true
console.log(converter.toBoolean('0'));           // false

// String conversion
console.log(converter.toString(42));              // '42'
console.log(converter.toString(true));            // 'true'

// JSON conversion
const obj = { name: 'John', age: 30 };
const json = converter.toJSON(obj);
console.log(json);                                // '{"name":"John","age":30}'
console.log(converter.fromJSON(json));            // { name: 'John', age: 30 }

// Encoding
console.log(converter.toBase64('hello'));         // 'aGVsbG8='
console.log(converter.fromBase64('aGVsbG8='));    // 'hello'

console.log(converter.toHex('hello'));            // '68656c6c6f'
console.log(converter.fromHex('68656c6c6f'));     // 'hello'
```

### Async Module

```javascript
const { async } = require('node-utils');

// Delay
console.log('Start');
await async.delay(1000);
console.log('1 second later');

// Timeout
try {
    const result = await async.timeout(
        fetch('/api/data'),
        5000
    );
    console.log(result);
} catch (err) {
    console.error('Request timed out');
}

// Retry
const result = await async.retry(
    () => fetch('/api/unstable'),
    { attempts: 3, delay: 1000, backoff: 2 }
);

// Parallel execution
const results = await async.parallel([
    () => fetch('/api/user/1'),
    () => fetch('/api/user/2'),
    () => fetch('/api/user/3')
], 2);

// Series execution
const seriesResults = await async.series([
    () => db.insert(user1),
    () => db.insert(user2),
    () => db.insert(user3)
]);

// Waterfall
const finalResult = await async.waterfall([
    () => getUser(1),
    (user) => getOrders(user.id),
    (orders) => calculateTotal(orders)
]);
```

### Cache Module

```javascript
const { cache } = require('node-utils');

// In-memory cache
const memCache = new cache.Cache({ ttl: 60000, maxSize: 100 });

memCache.set('user:1', { name: 'John', age: 30 });
const user = memCache.get('user:1');
console.log(user); // { name: 'John', age: 30 }

console.log(memCache.has('user:1')); // true
console.log(memCache.size());        // 1

memCache.delete('user:1');
memCache.clear();

// Redis cache (if Redis is available)
const redisCache = new cache.RedisCache({
    host: 'localhost',
    port: 6379,
    ttl: 60000
});

await redisCache.set('user:1', { name: 'John' });
const redisUser = await redisCache.get('user:1');
console.log(redisUser); // { name: 'John' }
```

## 📊 Performance

### Benchmarks

```bash
# Run benchmarks
npm run benchmark

# Results (ran on Node.js 20, Intel i7)
```

| Operation | Ops/sec | Latency (p99) |
|-----------|---------|---------------|
| math.add | 50,000,000 | 0.02ms |
| math.divide | 40,000,000 | 0.025ms |
| string.reverse (short) | 10,000,000 | 0.1ms |
| string.reverse (long) | 1,000,000 | 1ms |
| file.readFile (1KB) | 10,000 | 0.1ms |
| validator.isEmail | 5,000,000 | 0.2ms |
| cache.get | 20,000,000 | 0.05ms |

## 🧪 Testing

```bash
# Run all tests
npm test

# Run with coverage
npm run test:coverage

# Run in watch mode
npm run test:watch

# Run specific test file
npm test -- math.test.js

# Run benchmarks
npm run benchmark
```

## 📖 API Reference

For detailed API documentation, see the [API Reference](./API.md).

## 🚀 Performance Optimization

### Tips for Best Performance

1. **Use specific imports** to reduce bundle size:
```javascript
// Good - only import what you need
const { add, divide } = require('node-utils/math');

// Avoid - imports everything
const utils = require('node-utils');
```

2. **Cache results** for expensive operations:
```javascript
const cache = new cache.Cache({ ttl: 60000 });
const result = cache.get('key') ?? expensiveOperation();
```

3. **Use async versions** for I/O operations:
```javascript
// Good - non-blocking
await file.readFile('./data.txt');

// Avoid - blocking
file.readFileSync('./data.txt');
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NODE_UTILS_CACHE_TTL` | Default cache TTL (ms) | 60000 |
| `NODE_UTILS_CACHE_MAX_SIZE` | Maximum cache size | 1000 |
| `NODE_UTILS_LOG_LEVEL` | Log level (error/warn/info/debug) | info |

### Programmatic Configuration

```javascript
const { config } = require('node-utils');

config.set({
    cache: {
        ttl: 30000,
        maxSize: 500
    },
    logging: {
        level: 'debug',
        format: 'json'
    }
});
```

## 📦 Bundle Size

| Module | Size (minified) | Size (gzipped) |
|--------|----------------|----------------|
| math | 5KB | 1.5KB |
| string | 8KB | 2.5KB |
| file | 6KB | 2KB |
| validator | 10KB | 3KB |
| converter | 4KB | 1.2KB |
| async | 7KB | 2.2KB |
| cache | 5KB | 1.8KB |
| **Full** | 45KB | 14KB |

## 🛠️ Development

### Setup

```bash
# Clone repository
git clone https://github.com/username/node-utils.git
cd node-utils

# Install dependencies
npm install

# Build
npm run build

# Run tests
npm test
```

### Scripts

| Script | Description |
|--------|-------------|
| `npm run build` | Build the library |
| `npm test` | Run tests |
| `npm run test:coverage` | Run tests with coverage |
| `npm run lint` | Lint code |
| `npm run format` | Format code with Prettier |
| `npm run benchmark` | Run benchmarks |
| `npm run docs` | Generate documentation |

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📈 Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/username/node-utils/tags).

## 🐛 Known Issues

- File operations on Windows paths with spaces need escaping
- Cache keys are limited to strings
- Redis cache requires Redis server v5+

## 🗺️ Roadmap

### Version 2.0 (Planned)
- [ ] Full TypeScript rewrite
- [ ] Stream utilities
- [ ] WebAssembly modules for performance
- [ ] Plugin system
- [ ] CLI tool

### Version 1.x
- [x] Core utilities
- [x] File operations
- [x] Validation
- [x] Caching
- [x] Async helpers

## 📚 Additional Resources

- [API Documentation](./API.md)
- [Contributing Guide](./CONTRIBUTING.md)
- [Changelog](./CHANGELOG.md)
- [Examples](./examples)
- [Benchmarks](./benchmarks)

## 📞 Support

- 📧 Email: support@example.com
- 💬 Discord: [Join our Discord](https://discord.gg/example)
- 🐛 GitHub Issues: [Create an issue](https://github.com/username/node-utils/issues)
- 📖 Stack Overflow: Tag questions with `node-utils`

## 🙏 Acknowledgments

- The Node.js community for inspiration
- All our contributors
- Open source projects that made this possible

## 📊 Stats

- **Weekly Downloads**: 50,000+
- **GitHub Stars**: 1,200
- **Contributors**: 42
- **Test Coverage**: 100%

---

**Built with ❤️ for the Node.js community**

[![Star on GitHub](https://img.shields.io/github/stars/username/node-utils?style=social)](https://github.com/username/node-utils/stargazers)
[![Follow on Twitter](https://img.shields.io/twitter/follow/username?style=social)](https://twitter.com/username)
```

## ✅ **What this README provides:**

| Section | Purpose |
|---------|---------|
| **Badges** | npm version, license, build status, coverage, downloads |
| **Features** | Quick overview of what the library offers |
| **Installation** | npm/yarn/pnpm install commands |
| **Quick Start** | Complete working example to get started |
| **Modules** | Table of contents linking to detailed API docs |
| **Detailed Usage** | Comprehensive examples for each module |
| **Performance** | Benchmarks and optimization tips |
| **Testing** | How to run tests and benchmarks |
| **Bundle Size** | Size breakdown by module |
| **Development** | Setup instructions and available scripts |
| **Contributing** | How to contribute |
| **Roadmap** | Future plans |
| **Support** | How to get help |

## 🎯 **Purpose as Test Data**

This file serves as an **expected output** for validating that your readme generator produces:
- ✅ Proper npm package README format
- ✅ Comprehensive documentation structure
- ✅ Clear code examples
- ✅ Badges and metadata
- ✅ Module documentation
- ✅ Performance metrics
- ✅ Contributing guidelines