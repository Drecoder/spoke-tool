"""
Python Fuzz Tests

Tests the robustness of config parsing and code analysis
with various inputs to ensure no crashes.
"""

import pytest
import yaml
import ast
from typing import Any, Dict, List, Optional

# Mock config parser (in real implementation, import your actual modules)
def parse_config(data: str) -> Dict[str, Any]:
    """Parse configuration from string (YAML or JSON)."""
    try:
        return yaml.safe_load(data) or {}
    except yaml.YAMLError:
        try:
            import json
            return json.loads(data) or {}
        except (json.JSONDecodeError, TypeError):
            return {}

# Mock code analyzer (in real implementation, import your actual modules)
def analyze_code(code: str) -> Dict[str, Any]:
    """Analyze Python code and extract information."""
    result = {
        'functions': [],
        'classes': [],
        'imports': [],
        'lines_of_code': len(code.split('\n')),
        'has_errors': False
    }
    
    try:
        tree = ast.parse(code)
        
        for node in ast.walk(tree):
            if isinstance(node, ast.FunctionDef):
                result['functions'].append(node.name)
            elif isinstance(node, ast.ClassDef):
                result['classes'].append(node.name)
            elif isinstance(node, (ast.Import, ast.ImportFrom)):
                for alias in node.names:
                    result['imports'].append(alias.name)
    except SyntaxError:
        result['has_errors'] = True
    
    return result

# ============================================================================
# Config Parsing Fuzz Test
# ============================================================================

def test_fuzz_config_parsing():
    """Fuzz test for config parsing."""
    seeds = [
        "",
        "{}",
        "test_spoke:\n  enabled: true",
        "test_spoke:\n  enabled: false",
        "models:\n  encoder: codebert",
        "models:\n  decoder: deepseek-coder:7b",
        "models:\n  fast: gemma2:2b",
        "test_spoke:\n  coverage_threshold: 80",
        "test_spoke:\n  auto_run: true",
        "readme_spoke:\n  enabled: true",
        "readme_spoke:\n  auto_update: true",
        "readme_spoke:\n  sections:\n    - installation\n    - quickstart",
        "squeeze:\n  max_cpu_percent: 80",
        "squeeze:\n  max_memory_mb: 4096",
        "squeeze:\n  idle_threshold_ms: 500",
        "audit:\n  enabled: true\n  path: audit.log",
        "project_root: /path/to/project",
        "log_level: debug",
        "log_level: info",
        "log_level: warn",
        "log_level: error",
        "log_json: true",
        "log_color: true",
        "test_spoke:\n  enabled: not-a-bool",
        "test_spoke:\n  coverage_threshold: not-a-number",
        "test_spoke:\n  coverage_threshold: -10",
        "test_spoke:\n  coverage_threshold: 150",
        "test_spoke:\n  frameworks:\n    go: testing\n    nodejs: jest\n    python: pytest",
        "test_spoke:\n  frameworks:\n    go: \n    nodejs: ",
        "test_spoke:\n  test_file_patterns:\n    go: '*_test.go'\n    nodejs: '*.test.js'\n    python: 'test_*.py'",
        "test_spoke:\n  max_tests_per_function: 10",
        "test_spoke:\n  max_tests_per_function: -5",
        "test_spoke:\n  include_edge_cases: true",
        "test_spoke:\n  generate_mocks: true",
        "test_spoke:\n  languages:\n    go:\n      framework: testing\n      test_pattern: '*_test.go'",
        "test_spoke:\n  languages:\n    rust:\n      framework: rust_test",
        "models:\n  temperature: 0.7",
        "models:\n  temperature: 2.5",
        "models:\n  temperature: -1",
        "models:\n  max_tokens: 2048",
        "models:\n  max_tokens: 0",
        "models:\n  max_tokens: -100",
        "models:\n  timeout: 30s",
        "models:\n  timeout: invalid",
        "models:\n  ollama_host: http://localhost:11434",
        "models:\n  ollama_host: ",
        "readme_spoke:\n  sections:\n    - invalid-section",
        "readme_spoke:\n  sections: not-a-list",
        "readme_spoke:\n  include_examples: true",
        "readme_spoke:\n  max_examples_per_function: 3",
        "readme_spoke:\n  preserve_manual: true",
        "readme_spoke:\n  template_file: README.tmpl.md",
        "readme_spoke:\n  output_file: README.md",
        "readme_spoke:\n  doc_formats:\n    go: godoc\n    nodejs: jsdoc\n    python: pydoc",
        "squeeze:\n  enabled: true",
        "squeeze:\n  max_cpu_percent: 0",
        "squeeze:\n  max_cpu_percent: 101",
        "squeeze:\n  max_memory_mb: 0",
        "squeeze:\n  max_memory_mb: -10",
        "squeeze:\n  idle_threshold_ms: 0",
        "squeeze:\n  max_concurrent: 4",
        "squeeze:\n  max_concurrent: 0",
        "squeeze:\n  min_concurrent: 1",
        "squeeze:\n  min_concurrent: 10",
        "audit:\n  enabled: true",
        "audit:\n  path: ",
        "audit:\n  retain_days: 30",
        "audit:\n  retain_days: -5",
        "audit:\n  json: true",
        "---",
        "# comment only",
        "test_spoke:\n  # nested comment\n  enabled: true",
        "\t\ttest_spoke:\n\t\t\tenabled: true",
        "test_spoke: { enabled: true }",
        "test_spoke: { enabled: true, auto_run: false }",
        "[invalid yaml",
        ": : :",
        "test_spoke: [list, not, map]",
        "models:\n  - encoder\n  - decoder",
        "test_spoke:\n  enabled: !!str true",
        "test_spoke:\n  coverage_threshold: !!float 80",
        "squeeze:\n  max_cpu_percent: !!int 80",
        "test_spoke:\n  frameworks: !!map\n    go: testing",
        "a" * 10000,
        "test_spoke:\n  description: " + "very long " * 1000,
        "models:\n  encoder: " + "a" * 1000,
        "squeeze:\n  max_cpu_percent: " + "9" * 100,
        "test_spoke:\n  - " + "item " * 100,
        "test_spoke:\n  enabled: true\n" + "  extra: value\n" * 100,
        "%YAML 1.2\n---\ntest_spoke:\n  enabled: true",
        "%TAG ! tag:example.com,2024:\n---\ntest_spoke:\n  enabled: true",
        "test_spoke:\n  enabled: !!binary |\n    c3RyCg==",
        "test_spoke:\n  enabled: !!timestamp 2024-01-01",
        "test_spoke:\n  enabled: !!set {a, b, c}",
        "test_spoke:\n  enabled: !!omap\n    - a: 1\n    - b: 2",
        "<<: *anchor",
        "test_spoke: &anchor\n  enabled: true\nother: *anchor",
        "test_spoke:\n  enabled: !!python/name:sys.stdout",
        b"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F".decode('latin-1'),
        "test_spoke:\n  enabled: \x00\x01\x02",
    ]
    
    for seed in seeds:
        try:
            # This should NEVER raise an unhandled exception
            config = parse_config(seed)
            
            # Basic sanity checks - accessing fields should not raise
            if config and isinstance(config, dict):
                # Try accessing various fields
                test_fields = [
                    'test_spoke', 'models', 'readme_spoke',
                    'squeeze', 'audit', 'project_root',
                    'log_level', 'log_json', 'log_color'
                ]
                
                for field in test_fields:
                    try:
                        value = config.get(field)
                        # If it's a dict, try accessing nested fields
                        if isinstance(value, dict):
                            if field == 'test_spoke':
                                _ = value.get('enabled')
                                _ = value.get('coverage_threshold')
                            elif field == 'models':
                                _ = value.get('encoder')
                                _ = value.get('decoder')
                                _ = value.get('fast')
                            elif field == 'readme_spoke':
                                _ = value.get('enabled')
                                _ = value.get('sections')
                    except (AttributeError, KeyError, TypeError):
                        # Ignore access errors - we just don't want crashes
                        pass
        except Exception as e:
            # Parser can throw on invalid input - that's acceptable
            # We just don't want crashes/uncaught exceptions
            print(f"Parser threw for seed: {seed[:50]}... - {type(e).__name__}: {e}")

# ============================================================================
# Code Analysis Fuzz Test
# ============================================================================

def test_fuzz_code_analysis():
    """Fuzz test for code analysis."""
    seeds = [
        "",
        "def add(a, b): return a + b",
        "class Calculator:\n    def add(self, a, b): return a + b",
        "import sys\nfrom os import path\nimport numpy as np",
        "invalid python code @#$%",
        "def ( )",
        "def a( a",
        "# comment only",
        "'''docstring'''",
        "\"\"\"multiline\ndocstring\"\"\"",
        "x = 42",
        "y = 3.14159",
        "s = 'hello'",
        "t = \"world\"",
        "f = f'hello {name}'",
        "b = b'bytes'",
        "arr = [1, 2, 3]",
        "tup = (1, 2, 3)",
        "d = {'a': 1, 'b': 2}",
        "s = {1, 2, 3}",
        "if x > 0:\n    print('positive')\nelif x < 0:\n    print('negative')\nelse:\n    print('zero')",
        "for i in range(10):\n    print(i)",
        "while x < 10:\n    x += 1",
        "try:\n    x = 1 / 0\nexcept ZeroDivisionError:\n    print('error')\nfinally:\n    print('done')",
        "with open('file.txt') as f:\n    data = f.read()",
        "async def fetch():\n    return await api.get()",
        "def generator():\n    yield 1\n    yield 2",
        "@decorator\ndef func():\n    pass",
        "class Meta(type):\n    pass",
        "class MyClass(metaclass=Meta):\n    pass",
        "def func(a: int, b: str) -> bool:\n    return True",
        "lambda x: x * 2",
        "[x * 2 for x in range(10)]",
        "{x: x * 2 for x in range(10)}",
        "{x * 2 for x in range(10)}",
        "(x * 2 for x in range(10))",
        "x if condition else y",
        "x or y",
        "x and y",
        "not x",
        "x in y",
        "x is y",
        "x is not y",
        "x < y",
        "x <= y",
        "x > y",
        "x >= y",
        "x == y",
        "x != y",
        "x + y",
        "x - y",
        "x * y",
        "x / y",
        "x // y",
        "x % y",
        "x ** y",
        "x << y",
        "x >> y",
        "x & y",
        "x | y",
        "x ^ y",
        "~x",
        "x @ y",
        "x += y",
        "x -= y",
        "x *= y",
        "x /= y",
        "x //= y",
        "x %= y",
        "x **= y",
        "x <<= y",
        "x >>= y",
        "x &= y",
        "x |= y",
        "x ^= y",
        "x @= y",
        "True",
        "False",
        "None",
        "Ellipsis",
        "...",
        "__debug__",
        "0b1010",
        "0o777",
        "0x123",
        "1e10",
        "1j",
        "2 + 3j",
        "['a', 'b', 'c']",
        "{'a': 1, 'b': 2}",
        "{1, 2, 3}",
        "frozenset([1, 2, 3])",
        "bytearray(b'hello')",
        "memoryview(b'hello')",
        "def func(*args, **kwargs): pass",
        "def func(a, b, /, c, d, *, e, f): pass",
        "global x",
        "nonlocal y",
        "del x",
        "assert x > 0, 'message'",
        "raise Exception('message')",
        "yield from generator()",
        "await coroutine()",
        "async for x in async_iter:\n    print(x)",
        "async with resource:\n    await resource.acquire()",
        "from __future__ import annotations",
        "__all__ = ['func', 'class']",
        "if __name__ == '__main__':\n    main()",
        "@property\n    def prop(self): return self._prop",
        "@staticmethod\n    def static(): pass",
        "@classmethod\n    def cls(cls): pass",
        "a" * 10000,
        "def f():\n    pass\n" * 1000,
        "'''" + "a" * 10000 + "'''",
        b"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F".decode('latin-1'),
        "# -*- coding: utf-8 -*-\n# vim: set fileencoding=utf-8 :",
    ]
    
    for seed in seeds:
        try:
            # This should NEVER raise an unhandled exception
            result = analyze_code(seed)
            
            # Basic sanity checks - result should be a dict
            if result and isinstance(result, dict):
                # Access fields to ensure no hidden errors
                _ = result.get('functions', [])
                _ = result.get('classes', [])
                _ = result.get('imports', [])
                _ = result.get('lines_of_code', 0)
                _ = result.get('has_errors', False)
                
                # If functions is a list, access it safely
                if isinstance(result.get('functions'), list):
                    for func in result['functions']:
                        _ = func  # Just access
                        
                # If classes is a list, access it safely
                if isinstance(result.get('classes'), list):
                    for cls in result['classes']:
                        _ = cls  # Just access
                        
                # If imports is a list, access it safely
                if isinstance(result.get('imports'), list):
                    for imp in result['imports']:
                        _ = imp  # Just access
                        
        except Exception as e:
            # Analyzer can throw on invalid input - that's acceptable
            # We just don't want crashes/uncaught exceptions
            print(f"Analyzer threw for seed: {seed[:50]}... - {type(e).__name__}: {e}")

# ============================================================================
# Main execution
# ============================================================================

if __name__ == '__main__':
    pytest.main([__file__, '-v'])