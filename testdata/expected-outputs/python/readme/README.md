Perfect! Let's create a comprehensive README for a Python project that demonstrates proper documentation format and serves as a reference for the readme generator.

## 📁 **For `testdata/expected-outputs/python/readme/README.md`:**

```markdown
# Python Utility Library

A comprehensive, high-performance utility library for Python applications with minimal dependencies.

[![Python Version](https://img.shields.io/badge/python-3.11%2B-blue)](https://python.org)
[![PyPI version](https://img.shields.io/pypi/v/python-utils.svg)](https://pypi.org/project/python-utils/)
[![License](https://img.shields.io/pypi/l/python-utils.svg)](https://github.com/username/python-utils/blob/main/LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/username/python-utils/ci.yml)](https://github.com/username/python-utils/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/username/python-utils)](https://codecov.io/gh/username/python-utils)
[![Downloads](https://img.shields.io/pypi/dm/python-utils.svg)](https://pypi.org/project/python-utils/)

## ✨ Features

- ✅ **Math utilities** - Advanced mathematical operations with precision control
- ✅ **String manipulation** - Comprehensive string processing functions
- ✅ **File system helpers** - Safe and efficient file operations
- ✅ **Data validation** - Robust input validation utilities
- ✅ **Type conversion** - Reliable data type converters
- ✅ **Async utilities** - Coroutine helpers and control flow
- ✅ **Caching mechanisms** - In-memory and Redis cache implementations
- ✅ **DateTime utilities** - Date and time manipulation
- ✅ **Encryption helpers** - Hashing and encryption utilities
- ✅ **100% test coverage** - Fully tested with pytest
- ✅ **Type hints** - Complete type annotations
- ✅ **Zero external dependencies** - Lightweight and secure

## 📦 Installation

### Using pip

```bash
# Install from PyPI
pip install python-utils

# Install with specific extras
pip install python-utils[redis]  # Redis cache support
pip install python-utils[crypto] # Encryption support
pip install python-utils[all]    # All extras

# Install from source
git clone https://github.com/username/python-utils.git
cd python-utils
pip install -e .
```

### Using Poetry

```bash
poetry add python-utils
```

### Using pipenv

```bash
pipenv install python-utils
```

## 🚀 Quick Start

```python
from python_utils import (
    math_utils, string_utils, file_utils,
    validator, converter, async_utils
)

# Math operations
result = math_utils.add(5.2, 3.1)
print(f'5.2 + 3.1 = {result}')  # 8.3

result = math_utils.divide(10, 2)
print(f'10 / 2 = {result}')  # 5.0

# String manipulation
text = string_utils.reverse('hello')
print(f'hello reversed = {text}')  # 'olleh'

text = string_utils.to_camel_case('user_name')
print(f'user_name in camelCase = {text}')  # 'userName'

# File operations
file_utils.write_file('output.txt', 'Hello World')
content = file_utils.read_file('output.txt')
print(f'File content: {content}')

# Validation
is_valid = validator.is_email('user@example.com')
print(f'Is valid email? {is_valid}')  # True

# Async operations
import asyncio

async def main():
    result = await async_utils.delay(1.5)
    print('Waited 1.5 seconds')

asyncio.run(main())
```

## 📚 Modules

The library is organized into the following modules:

| Module | Description | Documentation |
|--------|-------------|---------------|
| **math_utils** | Mathematical operations | [API Docs](./API.md#module-math_utils) |
| **string_utils** | String manipulation | [API Docs](./API.md#module-string_utils) |
| **file_utils** | File system operations | [API Docs](./API.md#module-file_utils) |
| **validator** | Input validation | [API Docs](./API.md#module-validator) |
| **converter** | Data conversion | [API Docs](./API.md#module-converter) |
| **async_utils** | Asynchronous utilities | [API Docs](./API.md#module-async_utils) |
| **cache** | Caching mechanisms | [API Docs](./API.md#module-cache) |
| **datetime_utils** | Date and time utilities | [API Docs](./API.md#module-datetime_utils) |
| **encryption** | Encryption and hashing | [API Docs](./API.md#module-encryption) |

## 📖 Detailed Usage

### Math Utilities

```python
from python_utils import math_utils

# Basic arithmetic
print(math_utils.add(5, 3))          # 8
print(math_utils.subtract(10, 4))    # 6
print(math_utils.multiply(6, 7))     # 42
print(math_utils.divide(10, 2))      # 5.0

# Advanced math
print(math_utils.power(2, 8))        # 256
print(math_utils.sqrt(16))            # 4.0
print(math_utils.factorial(5))        # 120
print(math_utils.fibonacci(10))       # 55

# Number theory
print(math_utils.gcd(48, 18))         # 6
print(math_utils.lcm(12, 18))         # 36
print(math_utils.is_prime(17))        # True

# Statistics
numbers = [1, 2, 3, 4, 5]
print(math_utils.sum(numbers))        # 15
print(math_utils.average(numbers))    # 3.0
print(math_utils.median(numbers))     # 3.0
print(math_utils.mode([1, 2, 2, 3]))  # 2
print(math_utils.variance(numbers))   # 2.0

# Utilities
print(math_utils.clamp(15, 1, 10))    # 10
print(math_utils.lerp(0, 10, 0.5))    # 5.0
print(math_utils.remap(5, 0, 10, 0, 100))  # 50.0
```

### String Utilities

```python
from python_utils import string_utils

# Basic operations
print(string_utils.reverse('hello'))           # 'olleh'
print(string_utils.to_upper('hello'))          # 'HELLO'
print(string_utils.to_lower('HELLO'))          # 'hello'
print(string_utils.capitalize('hello'))        # 'Hello'
print(string_utils.title_case('hello world'))  # 'Hello World'

# Trimming
print(string_utils.trim('  hello  '))          # 'hello'
print(string_utils.trim_start('  hello'))      # 'hello'
print(string_utils.trim_end('hello  '))        # 'hello'

# Splitting and joining
print(string_utils.split('a,b,c', ','))        # ['a', 'b', 'c']
print(string_utils.join(['a', 'b', 'c'], '-')) # 'a-b-c'
print(string_utils.chars('hello'))             # ['h', 'e', 'l', 'l', 'o']
print(string_utils.words('hello world'))       # ['hello', 'world']

# Searching
print(string_utils.contains('hello world', 'world'))   # True
print(string_utils.index_of('hello world', 'world'))   # 6
print(string_utils.starts_with('hello', 'he'))         # True
print(string_utils.ends_with('world', 'ld'))           # True
print(string_utils.count('hello hello', 'hello'))      # 2

# Replacement
print(string_utils.replace('hello world', 'world', 'there'))  # 'hello there'
print(string_utils.replace_all('aaa', 'a', 'b'))              # 'bbb'

# Case conversion
print(string_utils.to_camel_case('user_name'))       # 'userName'
print(string_utils.to_pascal_case('user_name'))      # 'UserName'
print(string_utils.to_snake_case('userName'))        # 'user_name'
print(string_utils.to_kebab_case('userName'))        # 'user-name'

# Utilities
print(string_utils.truncate('This is a long string', 10))  # 'This is...'
print(string_utils.slugify('Hello World!'))                 # 'hello-world'
print(string_utils.pluralize('cat', 3))                     # 'cats'
```

### File Utilities

```python
from python_utils import file_utils

# Check existence
print(file_utils.exists('config.json'))        # True/False
print(file_utils.is_file('config.json'))       # True
print(file_utils.is_dir('docs'))               # True

# Read/write
content = file_utils.read_file('data.txt')
file_utils.write_file('output.txt', 'Hello World')
file_utils.append_file('log.txt', 'New entry\n')

# Copy/move/delete
file_utils.copy_file('source.txt', 'backup.txt')
file_utils.move_file('old.txt', 'new.txt')
file_utils.delete_file('temp.txt')

# Directory operations
files = file_utils.list_files('.', '*.py')
file_utils.ensure_directory('data/subdir')

# File info
size = file_utils.get_size('data.txt')
mtime = file_utils.get_modified_time('config.json')
print(f'Size: {size} bytes')
print(f'Modified: {mtime}')
```

### Validation

```python
from python_utils import validator

# Email validation
print(validator.is_email('user@example.com'))      # True
print(validator.is_email('invalid'))               # False

# URL validation
print(validator.is_url('https://example.com'))     # True

# Phone numbers
print(validator.is_phone('+1-555-123-4567'))       # True

# IP addresses
print(validator.is_ipv4('192.168.1.1'))            # True
print(validator.is_ipv6('2001:db8::1'))            # True

# UUID
uuid = '123e4567-e89b-12d3-a456-426614174000'
print(validator.is_uuid(uuid))                     # True

# Credit cards
print(validator.is_credit_card('4111111111111111'))  # True

# String validation
print(validator.is_alpha('abc'))                    # True
print(validator.is_numeric('123'))                  # True
print(validator.is_alphanumeric('abc123'))          # True

# Password strength
print(validator.is_strong_password('Passw0rd!'))    # True
```

### Type Conversion

```python
from python_utils import converter

# Number conversion
print(converter.to_int('42'))          # 42
print(converter.to_float('3.14'))      # 3.14

# Boolean conversion
print(converter.to_bool('true'))       # True
print(converter.to_bool('false'))      # False

# String conversion
print(converter.to_str(42))            # '42'
print(converter.to_str([1, 2, 3]))     # '[1, 2, 3]'

# JSON
data = {'name': 'John', 'age': 30}
json_str = converter.to_json(data)
print(json_str)                         # '{"name": "John", "age": 30}'
parsed = converter.from_json(json_str)
print(parsed['name'])                   # 'John'

# Encoding
print(converter.to_base64('hello'))     # 'aGVsbG8='
print(converter.from_base64('aGVsbG8=').decode())  # 'hello'

print(converter.to_hex('hello'))        # '68656c6c6f'
print(converter.from_hex('68656c6c6f').decode())  # 'hello'
```

### Async Utilities

```python
import asyncio
from python_utils import async_utils

async def main():
    # Delay
    print('Start')
    await async_utils.delay(1.5)
    print('After 1.5 seconds')
    
    # Timeout
    try:
        result = await async_utils.timeout(
            asyncio.sleep(5),
            2.0
        )
    except TimeoutError:
        print('Operation timed out')
    
    # Retry
    async def unstable():
        import random
        if random.random() < 0.7:
            raise ValueError('Random failure')
        return 'success'
    
    result = await async_utils.retry(
        unstable,
        attempts=5,
        delay=0.5,
        backoff=2.0
    )
    print(f'Result after retries: {result}')
    
    # Parallel execution
    async def task(i):
        await asyncio.sleep(0.1)
        return i
    
    results = await async_utils.parallel(
        [task(i) for i in range(10)],
        max_concurrency=3
    )
    print(f'Parallel results: {results}')

asyncio.run(main())
```

### Caching

```python
from python_utils.cache import Cache

# In-memory cache
cache = Cache(ttl=300, max_size=1000)

# Set values
cache.set('user:1', {'name': 'John', 'age': 30})
cache.set('user:2', {'name': 'Jane'}, ttl=60)

# Get values
user = cache.get('user:1')
print(user)  # {'name': 'John', 'age': 30}

# Check existence
if cache.has('user:1'):
    print('User exists in cache')

# Cache statistics
print(f'Size: {cache.size()}')
print(f'Keys: {cache.keys()}')
print(f'Hit rate: {cache.hit_rate():.2f}')

# Delete
cache.delete('user:2')
cache.clear()

# Cache decorator
from python_utils.cache import cached

@cached(ttl=60)
def expensive_function(x, y):
    # Simulate expensive computation
    import time
    time.sleep(2)
    return x * y

# First call - slow
result = expensive_function(5, 10)

# Second call - fast (cached)
result = expensive_function(5, 10)

# Redis cache (requires redis-py)
from python_utils.cache import RedisCache

redis_cache = RedisCache(
    host='localhost',
    port=6379,
    db=0,
    ttl=300
)
```

### DateTime Utilities

```python
from python_utils import datetime_utils
from datetime import datetime, timedelta

# Current time
now = datetime_utils.now()
utc_now = datetime_utils.utc_now()
print(f'Local: {now}, UTC: {utc_now}')

# Formatting
dt = datetime(2024, 1, 15, 14, 30, 0)
print(datetime_utils.format_date(dt))           # '2024-01-15'
print(datetime_utils.format_datetime(dt))       # '2024-01-15 14:30:00'
print(datetime_utils.format_iso(dt))            # '2024-01-15T14:30:00'

# Parsing
dt = datetime_utils.parse_date('2024-01-15')
dt = datetime_utils.parse_datetime('2024-01-15 14:30:00')
dt = datetime_utils.parse_iso('2024-01-15T14:30:00Z')

# Timezone conversion
dt_ny = datetime_utils.to_timezone(now, 'America/New_York')
dt_london = datetime_utils.to_timezone(now, 'Europe/London')

# Date arithmetic
tomorrow = datetime_utils.add_days(now, 1)
next_week = datetime_utils.add_weeks(now, 1)
next_month = datetime_utils.add_months(now, 1)
next_year = datetime_utils.add_years(now, 1)

# Difference
diff = datetime_utils.days_between(now, tomorrow)          # 1
diff = datetime_utils.hours_between(now, next_week)        # 168
diff = datetime_utils.seconds_between(now, tomorrow)       # 86400

# Human readable
print(datetime_utils.humanize(now - timedelta(seconds=30)))  # '30 seconds ago'
print(datetime_utils.humanize(now + timedelta(days=1)))      # 'in 1 day'

# Age calculation
birth_date = datetime(1990, 5, 15)
age = datetime_utils.age(birth_date)
print(f'Age: {age} years')

# Range generation
dates = datetime_utils.date_range(
    datetime(2024, 1, 1),
    datetime(2024, 1, 7),
    'days'
)

for dt in dates:
    print(dt.date())
```

### Encryption Utilities

```python
from python_utils import encryption

# Hashing
password = "my_secret_password"
hashed = encryption.hash_password(password)
print(f'Hashed: {hashed}')

# Verify password
is_valid = encryption.verify_password(password, hashed)
print(f'Valid: {is_valid}')

# MD5
md5_hash = encryption.md5('hello world')
print(f'MD5: {md5_hash}')

# SHA256
sha256_hash = encryption.sha256('hello world')
print(f'SHA256: {sha256_hash}')

# Base64 encoding
encoded = encryption.base64_encode('hello')
decoded = encryption.base64_decode(encoded)
print(f'Base64: {encoded}')
print(f'Decoded: {decoded}')

# AES encryption (requires cryptography)
key = encryption.generate_key()
ciphertext = encryption.encrypt_aes('secret message', key)
plaintext = encryption.decrypt_aes(ciphertext, key)
print(f'Decrypted: {plaintext}')

# JWT tokens
token = encryption.create_jwt({'user_id': 123}, 'secret')
payload = encryption.verify_jwt(token, 'secret')
print(f'Token payload: {payload}')
```

## 📊 Performance

### Benchmarks

```bash
# Run benchmarks
pytest tests/benchmarks/ -v

# Results (ran on Python 3.11, Intel i7)
```

| Operation | Ops/sec | Time (μs) |
|-----------|---------|-----------|
| math_utils.add | 10,000,000 | 0.1 |
| math_utils.divide | 8,000,000 | 0.125 |
| string_utils.reverse (short) | 5,000,000 | 0.2 |
| string_utils.reverse (long) | 100,000 | 10 |
| file_utils.read_file (1KB) | 10,000 | 100 |
| validator.is_email | 1,000,000 | 1 |
| cache.get | 5,000,000 | 0.2 |
| encryption.md5 | 500,000 | 2 |

## 🧪 Testing

```bash
# Install test dependencies
pip install pytest pytest-cov pytest-asyncio pytest-benchmark

# Run all tests
pytest

# Run with coverage
pytest --cov=src/ --cov-report=html

# Run specific test file
pytest tests/test_math_utils.py

# Run with verbose output
pytest -v

# Run benchmarks
pytest tests/benchmarks/ --benchmark-only

# Run with parallel execution
pytest -n auto
```

## 📖 API Reference

For detailed API documentation, see the [API Reference](./API.md).

## 🛠️ Development

### Setup

```bash
# Clone repository
git clone https://github.com/username/python-utils.git
cd python-utils

# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install development dependencies
pip install -e .[dev]
pip install -r requirements-dev.txt

# Run tests
pytest

# Run linters
flake8 src/ tests/
mypy src/
black --check src/ tests/
isort --check src/ tests/
```

### Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

### Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/username/python-utils/tags).

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Jane Smith** - *Initial work* - [@janesmith](https://github.com/janesmith)
- **John Doe** - *Documentation* - [@johndoe](https://github.com/johndoe)

See also the list of [contributors](https://github.com/username/python-utils/contributors) who participated in this project.

## 🙏 Acknowledgments

- The Python community for excellent tooling
- pytest developers for testing framework
- All contributors who have helped shape this project

## 📚 Additional Resources

- [API Documentation](./API.md)
- [Installation Guide](./INSTALL.md)
- [Contributing Guide](./CONTRIBUTING.md)
- [Changelog](./CHANGELOG.md)
- [Examples](./examples)

## 📞 Support

- 📧 Email: support@example.com
- 💬 Discord: [Join our Discord](https://discord.gg/example)
- 🐛 GitHub Issues: [Create an issue](https://github.com/username/python-utils/issues)
- 📖 Stack Overflow: Tag questions with `python-utils`

## 📊 Stats

- **PyPI Downloads**: 50,000+/month
- **GitHub Stars**: 1,200
- **Contributors**: 42
- **Test Coverage**: 100%

---

**Built with ❤️ for the Python community**

[![Star on GitHub](https://img.shields.io/github/stars/username/python-utils?style=social)](https://github.com/username/python-utils/stargazers)
[![Follow on Twitter](https://img.shields.io/twitter/follow/username?style=social)](https://twitter.com/username)
```

## ✅ **What this README provides:**

| Section | Purpose |
|---------|---------|
| **Badges** | PyPI version, Python versions, license, build status, coverage, downloads |
| **Features** | Quick overview of what the library offers |
| **Installation** | pip, poetry, pipenv install commands |
| **Quick Start** | Complete working example to get started |
| **Modules** | Table of contents linking to detailed API docs |
| **Detailed Usage** | Comprehensive examples for each module |
| **Performance** | Benchmarks and metrics |
| **Testing** | How to run tests and benchmarks |
| **API Reference** | Link to detailed API documentation |
| **Development** | Setup instructions and available scripts |
| **Contributing** | How to contribute |
| **Support** | How to get help |

## 🎯 **Purpose as Test Data**

This file serves as an **expected output** for validating that your readme generator produces:
- ✅ Proper Python package README format
- ✅ Comprehensive documentation structure
- ✅ Clear code examples
- ✅ Badges and metadata
- ✅ Module documentation
- ✅ Performance metrics
- ✅ Contributing guidelines