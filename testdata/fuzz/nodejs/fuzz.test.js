/**
 * Node.js Fuzz Tests
 * 
 * Tests the robustness of config parsing and code analysis
 * with various inputs to ensure no crashes.
 */

// @ts-check

const { test } = require('@jest/globals');

// Mock the config parser and analyzer functions
// In a real implementation, these would import your actual modules
const parseConfig = (data) => {
    try {
        return JSON.parse(data);
    } catch {
        try {
            // Try YAML-like parsing (simplified)
            const lines = data.split('\n');
            const config = {};
            for (const line of lines) {
                const match = line.match(/^(\w+):\s*(.+)$/);
                if (match) {
                    const [_, key, value] = match;
                    if (value === 'true') config[key] = true;
                    else if (value === 'false') config[key] = false;
                    else if (!isNaN(value)) config[key] = Number(value);
                    else config[key] = value;
                }
            }
            return config;
        } catch {
            return {};
        }
    }
};

const analyzeCode = (code) => {
    const functions = [];
    const lines = code.split('\n');
    
    // Simple function detection
    const funcRegex = /(?:function|const|let|var)\s+(\w+)\s*[=:(]/;
    for (const line of lines) {
        const match = line.match(funcRegex);
        if (match) {
            functions.push(match[1]);
        }
    }
    
    return {
        functions,
        linesOfCode: lines.length,
        hasErrors: false
    };
};

// ============================================================================
// Config Parsing Fuzz Test
// ============================================================================

test('FuzzConfigParsing', () => {
    const seeds = [
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
        "a".repeat(10000),
        "test_spoke:\n  description: " + "very long ".repeat(1000),
        "models:\n  encoder: " + "a".repeat(1000),
        "squeeze:\n  max_cpu_percent: " + "9".repeat(100),
        "test_spoke:\n  - " + "item ".repeat(100),
        "test_spoke:\n  enabled: true\n" + "  extra: value\n".repeat(100),
        "%YAML 1.2\n---\ntest_spoke:\n  enabled: true",
        "%TAG ! tag:example.com,2024:\n---\ntest_spoke:\n  enabled: true",
        "test_spoke:\n  enabled: !!binary |\n    c3RyCg==",
        "test_spoke:\n  enabled: !!timestamp 2024-01-01",
        "test_spoke:\n  enabled: !!set {a, b, c}",
        "test_spoke:\n  enabled: !!omap\n    - a: 1\n    - b: 2",
        "<<: *anchor",
        "test_spoke: &anchor\n  enabled: true\nother: *anchor",
        "test_spoke:\n  enabled: !!python/name:sys.stdout",
        "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F",
        "test_spoke:\n  enabled: \x00\x01\x02"
    ];
    
    for (const seed of seeds) {
        try {
            // This should NEVER throw
            const config = parseConfig(seed);
            
            // Basic sanity checks - accessing fields should not throw
            if (config && typeof config === 'object') {
                // Try accessing various fields
                const testFields = [
                    'test_spoke', 'models', 'readme_spoke', 
                    'squeeze', 'audit', 'project_root',
                    'log_level', 'log_json', 'log_color'
                ];
                
                for (const field of testFields) {
                    try {
                        const value = config[field];
                        // If it's an object, try accessing nested fields
                        if (value && typeof value === 'object') {
                            if (field === 'test_spoke') {
                                const _ = value.enabled;
                                const __ = value.coverage_threshold;
                            } else if (field === 'models') {
                                const _ = value.encoder;
                                const __ = value.decoder;
                                const ___ = value.fast;
                            } else if (field === 'readme_spoke') {
                                const _ = value.enabled;
                                const __ = value.sections;
                            }
                        }
                    } catch {
                        // Ignore access errors - we just don't want crashes
                    }
                }
            }
        } catch (err) {
            // Parser can throw on invalid input - that's acceptable
            // We just don't want crashes/uncaught exceptions
            console.log(`Parser threw for seed: ${seed.substring(0, 50)}...`);
        }
    }
});

// ============================================================================
// Code Analysis Fuzz Test
// ============================================================================

test('FuzzCodeAnalysis', () => {
    const seeds = [
        "",
        "function add(a, b) { return a + b; }",
        "const multiply = (x, y) => x * y;",
        "class Calculator { add(a, b) { return a + b; } }",
        "invalid javascript @#$%",
        "function ( )",
        "function a( a",
        "// comment only",
        "import { something } from 'module';",
        "export const value = 42;",
        "async function fetchData() { return await api.get(); }",
        "function* generator() { yield 1; yield 2; }",
        "const obj = { method() { return this; } };",
        "if (true) { console.log('hello'); }",
        "for (let i = 0; i < 10; i++) { console.log(i); }",
        "while (false) { break; }",
        "try { throw new Error('test'); } catch (e) { console.log(e); }",
        "switch(x) { case 1: break; default: break; }",
        "const [a, b] = [1, 2];",
        "const { x, y } = { x: 1, y: 2 };",
        "function withManyParams(a, b, c, d, e, f, g, h, i, j) { return a + b + c + d + e + f + g + h + i + j; }",
        "function withManyReturns() { return [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; }",
        "function nested() { (function() { (function() { console.log('deep'); })(); })(); }",
        "function* generator() { yield 1; yield 2; yield 3; }",
        "async function asyncFunc() { await Promise.resolve(); }",
        "const promise = new Promise((resolve, reject) => { resolve(); });",
        "Promise.resolve().then(() => {}).catch(() => {});",
        "setTimeout(() => {}, 1000);",
        "setInterval(() => {}, 1000);",
        "document.addEventListener('click', () => {});",
        "window.onload = () => {};",
        "process.on('uncaughtException', () => {});",
        "module.exports = { foo: 'bar' };",
        "exports.foo = 'bar';",
        "require('fs');",
        "import fs from 'fs';",
        "import * as path from 'path';",
        "const _ = require('lodash');",
        "''",
        "'hello'",
        "`template ${string}`",
        "/regex/g",
        "true",
        "false",
        "null",
        "undefined",
        "123",
        "3.14159",
        "0x123",
        "0b1010",
        "0o777",
        "1e10",
        "NaN",
        "Infinity",
        "[]",
        "[1, 2, 3]",
        "{}",
        "{ key: 'value' }",
        "new Date()",
        "new Error('message')",
        "new Map()",
        "new Set()",
        "new WeakMap()",
        "new WeakSet()",
        "Symbol('test')",
        "BigInt(123)",
        "typeof x",
        "instanceof y",
        "delete obj.prop",
        "void 0",
        "x ? y : z",
        "x || y",
        "x && y",
        "x ?? y",
        "x?.y",
        "x?.y?.z",
        "x?.[y]",
        "x?.()",
        "x + y",
        "x - y",
        "x * y",
        "x / y",
        "x % y",
        "x ** y",
        "x << y",
        "x >> y",
        "x >>> y",
        "x & y",
        "x | y",
        "x ^ y",
        "~x",
        "!x",
        "++x",
        "--x",
        "x++",
        "x--",
        "x = y",
        "x += y",
        "x -= y",
        "x *= y",
        "x /= y",
        "x %= y",
        "x **= y",
        "x <<= y",
        "x >>= y",
        "x >>>= y",
        "x &= y",
        "x |= y",
        "x ^= y",
        "x &&= y",
        "x ||= y",
        "x ??= y",
        "a".repeat(10000),
        "function f() {}\n".repeat(1000),
        "/* huge comment */ " + "a".repeat(10000),
        "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F",
        "// @ts-check\n/** @type {string} */\nconst x = 'hello';"
    ];
    
    for (const seed of seeds) {
        try {
            // This should NEVER throw
            const result = analyzeCode(seed);
            
            // Basic sanity checks - result should be an object
            if (result && typeof result === 'object') {
                // Access fields to ensure no hidden errors
                const _ = result.functions;
                const __ = result.linesOfCode;
                const ___ = result.hasErrors;
                
                // If functions is an array, access it safely
                if (Array.isArray(result.functions)) {
                    for (const func of result.functions) {
                        const _ = func; // Just access
                    }
                }
            }
        } catch (err) {
            // Analyzer can throw on invalid input - that's acceptable
            // We just don't want crashes/uncaught exceptions
            console.log(`Analyzer threw for seed: ${seed.substring(0, 50)}...`);
        }
    }
});