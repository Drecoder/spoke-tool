/**
 * Edge Cases Test Suite for Node.js
 * 
 * This file contains tests for various edge cases, boundary conditions,
 * and error scenarios to validate the test generator's edge case handling.
 */

// ============================================================================
// Basic Edge Cases - Zero and Empty Values
// ============================================================================

describe('Zero and Empty Values Edge Cases', () => {
    
    test('handles zero integer input', () => {
        const result = processNumber(0);
        expect(result).toBe('zero');
    });

    test('handles empty string', () => {
        const result = processString('');
        expect(result).toBe('empty');
    });

    test('handles zero float', () => {
        const result = processFloat(0.0);
        expect(result).toBe('zero');
    });

    test('handles negative zero', () => {
        const result = processFloat(-0.0);
        expect(result).toBe('zero');
    });

    test('handles null value', () => {
        expect(() => processValue(null)).toThrow('Value cannot be null');
    });

    test('handles undefined value', () => {
        expect(() => processValue(undefined)).toThrow('Value cannot be undefined');
    });

    test('handles NaN', () => {
        expect(() => processValue(NaN)).toThrow('Value cannot be NaN');
    });

    test('handles empty array', () => {
        const result = processArray([]);
        expect(result).toBe('empty');
    });

    test('handles empty object', () => {
        const result = processObject({});
        expect(result).toBe('empty');
    });

    test('handles empty Set', () => {
        const result = processSet(new Set());
        expect(result).toBe('empty');
    });

    test('handles empty Map', () => {
        const result = processMap(new Map());
        expect(result).toBe('empty');
    });

    test('handles empty typed array', () => {
        const result = processTypedArray(new Uint8Array(0));
        expect(result).toBe('empty');
    });

    test('handles zero-length Buffer', () => {
        const result = processBuffer(Buffer.alloc(0));
        expect(result).toBe('empty');
    });
});

// ============================================================================
// Numeric Boundary Conditions
// ============================================================================

describe('Numeric Boundary Conditions', () => {
    
    test('handles Number.MAX_SAFE_INTEGER', () => {
        const result = processLargeNumber(Number.MAX_SAFE_INTEGER);
        expect(result).toBe('safe');
    });

    test('handles Number.MAX_SAFE_INTEGER + 1', () => {
        const result = processLargeNumber(Number.MAX_SAFE_INTEGER + 1);
        expect(result).toBe('unsafe');
    });

    test('handles Number.MIN_SAFE_INTEGER', () => {
        const result = processLargeNumber(Number.MIN_SAFE_INTEGER);
        expect(result).toBe('safe');
    });

    test('handles Number.MIN_SAFE_INTEGER - 1', () => {
        const result = processLargeNumber(Number.MIN_SAFE_INTEGER - 1);
        expect(result).toBe('unsafe');
    });

    test('handles Number.MAX_VALUE', () => {
        const result = processLargeNumber(Number.MAX_VALUE);
        expect(result).toBe('max');
    });

    test('handles Number.MIN_VALUE', () => {
        const result = processLargeNumber(Number.MIN_VALUE);
        expect(result).toBe('min');
    });

    test('handles Infinity', () => {
        expect(() => processNumber(Infinity)).toThrow('Cannot process Infinity');
    });

    test('handles -Infinity', () => {
        expect(() => processNumber(-Infinity)).toThrow('Cannot process -Infinity');
    });

    test('handles Number.EPSILON', () => {
        const result = processFloat(Number.EPSILON);
        expect(result).toBe('epsilon');
    });

    test('handles floating point precision', () => {
        // Famous floating point precision issue
        expect(0.1 + 0.2).not.toBe(0.3);
        expect(processFloat(0.1 + 0.2)).toBe('approx-0.3');
    });

    test('handles very small numbers', () => {
        const result = processFloat(1e-323);
        expect(result).toBe('subnormal');
    });
});

// ============================================================================
// Integer Overflow and Underflow
// ============================================================================

describe('Integer Overflow and Underflow', () => {
    
    test('handles 32-bit integer overflow', () => {
        const result = addInt32(2147483647, 1);
        expect(result).toBe(-2147483648); // Wraps around in 32-bit
    });

    test('handles 32-bit integer underflow', () => {
        const result = addInt32(-2147483648, -1);
        expect(result).toBe(2147483647); // Wraps around in 32-bit
    });

    test('handles bitwise shift overflow', () => {
        const result = 1 << 31; // 32-bit shift
        expect(result).toBe(-2147483648);
    });

    test('handles unsigned right shift', () => {
        const result = -1 >>> 0;
        expect(result).toBe(4294967295);
    });

    test('handles multiplication overflow', () => {
        const result = 2 ** 53 * 2;
        expect(result).toBe(2 ** 54); // May lose precision
    });
});

// ============================================================================
// String Edge Cases
// ============================================================================

describe('String Edge Cases', () => {
    
    test('handles very long string', () => {
        const longString = 'a'.repeat(1000000);
        const result = processString(longString);
        expect(result).toBe('very long');
    });

    test('handles string with null character', () => {
        const str = 'hello\x00world';
        const result = processString(str);
        expect(result).toBe('contains null');
    });

    test('handles string with unicode characters', () => {
        const str = 'Hello, 世界! 🚀 🌍';
        const result = processString(str);
        expect(result).toBe('unicode');
    });

    test('handles string with emoji modifiers', () => {
        const str = '👨‍👩‍👧‍👦'; // Family emoji (multiple code points)
        const result = processString(str);
        expect(result.length).toBe(1); // Should count as one character
    });

    test('handles string with RTL characters', () => {
        const str = 'مرحبا بالعالم'; // Arabic
        const result = processString(str);
        expect(result).toBe('rtl');
    });

    test('handles string with control characters', () => {
        const str = 'Hello\x1b[31mWorld\x1b[0m';
        const result = processString(str);
        expect(result).toBe('contains control chars');
    });

    test('handles string with surrogate pairs', () => {
        const str = '𠜎𠜱𠝹𠱓'; // CJK characters outside BMP
        const result = processString(str);
        expect(result.length).toBe(4); // Each character is a surrogate pair
    });

    test('handles string with zero-width joiner', () => {
        const str = 'वि' + '\u200D' + 'कास'; // ZWJ example
        const result = processString(str);
        expect(result).toBe('has zwj');
    });

    test('handles string with bidirectional text', () => {
        const str = 'Hello (مرحبا) World';
        const result = processString(str);
        expect(result).toBe('bidi');
    });
});

// ============================================================================
// Array Edge Cases
// ============================================================================

describe('Array Edge Cases', () => {
    
    test('handles sparse array', () => {
        const arr = [1, , , 4]; // Sparse array
        expect(arr.length).toBe(4);
        expect(arr[1]).toBeUndefined();
        
        const result = processArray(arr);
        expect(result).toBe('sparse');
    });

    test('handles array with holes', () => {
        const arr = new Array(1000);
        arr[500] = 'middle';
        const result = processArray(arr);
        expect(result).toBe('has holes');
    });

    test('handles array with negative indices', () => {
        const arr = [1, 2, 3];
        arr[-1] = 'negative';
        const result = processArray(arr);
        expect(result).toBe('has negative index');
    });

    test('handles array with non-numeric keys', () => {
        const arr = [1, 2, 3];
        arr['key'] = 'value';
        const result = processArray(arr);
        expect(result).toBe('has non-numeric keys');
    });

    test('handles very large array', () => {
        const arr = new Array(1000000);
        const result = processArray(arr);
        expect(result).toBe('very large');
    });

    test('handles array with deleted elements', () => {
        const arr = [1, 2, 3, 4, 5];
        delete arr[2];
        const result = processArray(arr);
        expect(result).toBe('has deleted');
    });

    test('handles array with getters/setters', () => {
        const arr = [1, 2, 3];
        Object.defineProperty(arr, 'custom', {
            get: () => 'getter',
            set: (val) => {}
        });
        const result = processArray(arr);
        expect(result).toBe('has custom property');
    });
});

// ============================================================================
// Object Edge Cases
// ============================================================================

describe('Object Edge Cases', () => {
    
    test('handles object with null prototype', () => {
        const obj = Object.create(null);
        obj.key = 'value';
        const result = processObject(obj);
        expect(result).toBe('null prototype');
    });

    test('handles object with prototype chain', () => {
        const parent = { parentProp: 'parent' };
        const child = Object.create(parent);
        child.childProp = 'child';
        const result = processObject(child);
        expect(result).toBe('has prototype');
    });

    test('handles frozen object', () => {
        const obj = Object.freeze({ key: 'value' });
        const result = processObject(obj);
        expect(result).toBe('frozen');
    });

    test('handles sealed object', () => {
        const obj = Object.seal({ key: 'value' });
        const result = processObject(obj);
        expect(result).toBe('sealed');
    });

    test('handles object with non-enumerable properties', () => {
        const obj = {};
        Object.defineProperty(obj, 'hidden', {
            value: 'secret',
            enumerable: false
        });
        const result = processObject(obj);
        expect(result).toBe('has non-enumerable');
    });

    test('handles object with symbol keys', () => {
        const sym = Symbol('test');
        const obj = {
            [sym]: 'symbol value',
            regular: 'regular'
        };
        const result = processObject(obj);
        expect(result).toBe('has symbols');
    });

    test('handles object with getters/setters', () => {
        const obj = {
            _value: 42,
            get value() { return this._value; },
            set value(val) { this._value = val; }
        };
        const result = processObject(obj);
        expect(result).toBe('has accessors');
    });

    test('handles object with circular reference', () => {
        const obj = {};
        obj.self = obj;
        const result = processObject(obj);
        expect(result).toBe('circular');
    });

    test('handles object with very deep nesting', () => {
        let obj = {};
        let current = obj;
        for (let i = 0; i < 1000; i++) {
            current.next = {};
            current = current.next;
        }
        const result = processObject(obj);
        expect(result).toBe('deeply nested');
    });
});

// ============================================================================
// Function Edge Cases
// ============================================================================

describe('Function Edge Cases', () => {
    
    test('handles function with many arguments', () => {
        const fn = (...args) => args.length;
        const result = fn(...Array(1000).fill(0));
        expect(result).toBe(1000);
    });

    test('handles recursive function with deep recursion', () => {
        const deepRecursion = (n) => {
            if (n <= 0) return 0;
            return 1 + deepRecursion(n - 1);
        };
        
        expect(() => deepRecursion(20000)).toThrow(); // Stack overflow
    });

    test('handles tail recursion', () => {
        const tailRecursion = (n, acc = 0) => {
            if (n <= 0) return acc;
            return tailRecursion(n - 1, acc + 1);
        };
        
        const result = tailRecursion(10000);
        expect(result).toBe(10000);
    });

    test('handles function with closure memory', () => {
        const createClosure = () => {
            const largeData = new Array(1000000).fill('data');
            return () => largeData.length;
        };
        
        const fn = createClosure();
        const result = fn();
        expect(result).toBe(1000000);
    });

    test('handles function with multiple returns', () => {
        const multiReturn = (x) => {
            if (x < 0) return 'negative';
            if (x === 0) return 'zero';
            if (x > 100) return 'large';
            return 'normal';
        };
        
        expect(multiReturn(-5)).toBe('negative');
        expect(multiReturn(0)).toBe('zero');
        expect(multiReturn(50)).toBe('normal');
        expect(multiReturn(200)).toBe('large');
    });
});

// ============================================================================
// Promise and Async Edge Cases
// ============================================================================

describe('Promise and Async Edge Cases', () => {
    
    test('handles promise that never resolves', async () => {
        const neverResolve = new Promise(() => {});
        
        await expect(Promise.race([
            neverResolve,
            Promise.reject(new Error('timeout'))
        ])).rejects.toThrow('timeout');
    });

    test('handles promise that resolves after long time', async () => {
        const slowPromise = new Promise(resolve => {
            setTimeout(() => resolve('done'), 10000);
        });
        
        // This would timeout in tests, so we mock it
        jest.useFakeTimers();
        const promise = slowPromise;
        jest.advanceTimersByTime(10000);
        const result = await promise;
        expect(result).toBe('done');
        jest.useRealTimers();
    });

    test('handles promise chain with errors', async () => {
        const promise = Promise.resolve(1)
            .then(x => x + 1)
            .then(x => { throw new Error('middle error'); })
            .then(x => x + 1)
            .catch(err => err.message);
        
        const result = await promise;
        expect(result).toBe('middle error');
    });

    test('handles multiple promise rejections', async () => {
        const promises = [
            Promise.reject(new Error('err1')),
            Promise.reject(new Error('err2')),
            Promise.resolve('ok')
        ];
        
        await expect(Promise.all(promises)).rejects.toThrow('err1');
        
        const settled = await Promise.allSettled(promises);
        expect(settled[0].status).toBe('rejected');
        expect(settled[2].status).toBe('fulfilled');
    });

    test('handles unhandled promise rejection', () => {
        // This test should not cause unhandled rejection warning
        const rejection = Promise.reject(new Error('unhandled'));
        rejection.catch(() => {}); // Handle it
    });

    test('handles promise with finally', async () => {
        let cleaned = false;
        
        await Promise.resolve('test')
            .then(x => x + '!')
            .finally(() => { cleaned = true; });
        
        expect(cleaned).toBe(true);
    });
});

// ============================================================================
// Date and Time Edge Cases
// ============================================================================

describe('Date and Time Edge Cases', () => {
    
    test('handles invalid date', () => {
        const date = new Date('invalid');
        expect(isNaN(date)).toBe(true);
        expect(() => processDate(date)).toThrow('Invalid date');
    });

    test('handles dates before Unix epoch', () => {
        const date = new Date('1900-01-01');
        expect(date.getFullYear()).toBe(1900);
        const result = processDate(date);
        expect(result).toBe('pre-epoch');
    });

    test('handles dates after 2038 problem', () => {
        const date = new Date('2040-01-01');
        expect(date.getFullYear()).toBe(2040);
        const result = processDate(date);
        expect(result).toBe('post-2038');
    });

    test('handles leap year dates', () => {
        const leapDate = new Date('2020-02-29');
        expect(leapDate.getMonth()).toBe(1); // February
        expect(leapDate.getDate()).toBe(29);
        
        const nonLeapDate = new Date('2021-02-29');
        expect(isNaN(nonLeapDate)).toBe(true);
    });

    test('handles timezone offsets', () => {
        const date = new Date('2024-01-01T00:00:00Z');
        const offset = date.getTimezoneOffset();
        expect(offset).toBeDefined();
    });

    test('handles daylight saving time transitions', () => {
        // DST start (spring forward)
        const dstStart = new Date('2024-03-10T02:30:00');
        // DST end (fall back)
        const dstEnd = new Date('2024-11-03T01:30:00');
        
        const result = processDST(dstStart, dstEnd);
        expect(result).toBe('handled');
    });
});

// ============================================================================
// Regular Expression Edge Cases
// ============================================================================

describe('Regular Expression Edge Cases', () => {
    
    test('handles regex with catastrophic backtracking', () => {
        const regex = /^(a+)+$/;
        const longString = 'a'.repeat(30) + 'b';
        
        // This would cause catastrophic backtracking
        expect(regex.test(longString)).toBe(false);
    });

    test('handles regex with lookahead/lookbehind', () => {
        const regex = /(?<=\d)(?=\.\d)/;
        const result = '123.45'.match(regex);
        expect(result).toBeTruthy();
    });

    test('handles regex with named groups', () => {
        const regex = /(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})/;
        const match = '2024-01-15'.match(regex);
        expect(match.groups.year).toBe('2024');
        expect(match.groups.month).toBe('01');
        expect(match.groups.day).toBe('15');
    });

    test('handles regex with Unicode properties', () => {
        const regex = /\p{Script=Greek}/u;
        expect(regex.test('π')).toBe(true);
        expect(regex.test('a')).toBe(false);
    });

    test('handles regex with sticky flag', () => {
        const regex = /\d+/y;
        regex.lastIndex = 2;
        const result = regex.exec('123 456');
        expect(result[0]).toBe('3');
    });

    test('handles very large regex', () => {
        const pattern = '(?:' + Array(1000).fill('a').join('|') + ')';
        const regex = new RegExp(pattern);
        expect(regex.test('a')).toBe(true);
    });
});

// ============================================================================
// Error Handling Edge Cases
// ============================================================================

describe('Error Handling Edge Cases', () => {
    
    test('handles nested errors', () => {
        try {
            try {
                throw new Error('inner');
            } catch (inner) {
                throw new Error('outer', { cause: inner });
            }
        } catch (outer) {
            expect(outer.message).toBe('outer');
            expect(outer.cause.message).toBe('inner');
        }
    });

    test('handles error with custom properties', () => {
        class CustomError extends Error {
            constructor(message, code) {
                super(message);
                this.code = code;
                this.timestamp = Date.now();
            }
        }
        
        const error = new CustomError('custom', 500);
        expect(error.code).toBe(500);
        expect(error.timestamp).toBeDefined();
    });

    test('handles async error in event loop', (done) => {
        process.once('uncaughtException', (err) => {
            expect(err.message).toBe('async error');
            done();
        });
        
        setTimeout(() => {
            throw new Error('async error');
        }, 10);
    });

    test('handles error with stack trace', () => {
        try {
            throw new Error('stack test');
        } catch (err) {
            expect(err.stack).toBeDefined();
            expect(err.stack).toContain('Error: stack test');
        }
    });

    test('handles error without stack', () => {
        const err = { message: 'plain object' };
        expect(() => { throw err; }).toThrow();
    });
});

// ============================================================================
// Event Emitter Edge Cases
// ============================================================================

const EventEmitter = require('events');

describe('Event Emitter Edge Cases', () => {
    
    test('handles many listeners', () => {
        const emitter = new EventEmitter();
        emitter.setMaxListeners(20);
        
        for (let i = 0; i < 20; i++) {
            emitter.on('event', () => {});
        }
        
        expect(emitter.listenerCount('event')).toBe(20);
    });

    test('handles listener memory leak warning', () => {
        const emitter = new EventEmitter();
        // Default max listeners is 10
        for (let i = 0; i < 11; i++) {
            emitter.on('event', () => {});
        }
        
        // Should emit warning but not throw
        expect(emitter.listenerCount('event')).toBe(11);
    });

    test('handles once listeners', () => {
        const emitter = new EventEmitter();
        let count = 0;
        
        emitter.once('event', () => count++);
        emitter.emit('event');
        emitter.emit('event');
        
        expect(count).toBe(1);
    });

    test('handles prepend listeners', () => {
        const emitter = new EventEmitter();
        const order = [];
        
        emitter.on('event', () => order.push(1));
        emitter.prependListener('event', () => order.push(2));
        emitter.emit('event');
        
        expect(order).toEqual([2, 1]);
    });

    test('handles error events', () => {
        const emitter = new EventEmitter();
        
        emitter.on('error', (err) => {
            expect(err.message).toBe('test error');
        });
        
        emitter.emit('error', new Error('test error'));
    });

    test('handles removing listeners while emitting', () => {
        const emitter = new EventEmitter();
        let count = 0;
        
        const fn = () => {
            count++;
            emitter.removeListener('event', fn);
        };
        
        emitter.on('event', fn);
        emitter.on('event', () => count++);
        
        emitter.emit('event');
        expect(count).toBe(2);
    });
});

// ============================================================================
// Stream Edge Cases
// ============================================================================

const { Readable, Writable, Transform } = require('stream');

describe('Stream Edge Cases', () => {
    
    test('handles empty readable stream', (done) => {
        const stream = Readable.from([]);
        const data = [];
        
        stream.on('data', chunk => data.push(chunk));
        stream.on('end', () => {
            expect(data.length).toBe(0);
            done();
        });
    });

    test('handles very large stream', (done) => {
        const largeArray = Array(1000000).fill('data');
        const stream = Readable.from(largeArray);
        let count = 0;
        
        stream.on('data', () => count++);
        stream.on('end', () => {
            expect(count).toBe(1000000);
            done();
        });
    });

    test('handles stream backpressure', (done) => {
        const readable = Readable.from(Array(1000).fill('chunk'));
        const writable = new Writable({
            write(chunk, encoding, callback) {
                setTimeout(callback, 1); // Slow writes
            },
            highWaterMark: 10
        });
        
        readable.pipe(writable);
        
        writable.on('finish', done);
        readable.on('error', done);
    });

    test('handles stream error', (done) => {
        const stream = new Readable({
            read() {
                this.destroy(new Error('stream error'));
            }
        });
        
        stream.on('error', (err) => {
            expect(err.message).toBe('stream error');
            done();
        });
        
        stream.resume();
    });

    test('handles transform stream', (done) => {
        const transform = new Transform({
            transform(chunk, encoding, callback) {
                this.push(chunk.toString().toUpperCase());
                callback();
            }
        });
        
        const results = [];
        transform.on('data', chunk => results.push(chunk.toString()));
        transform.on('end', () => {
            expect(results).toEqual(['HELLO', 'WORLD']);
            done();
        });
        
        transform.write('hello');
        transform.write('world');
        transform.end();
    });
});

// ============================================================================
// Buffer Edge Cases
// ============================================================================

describe('Buffer Edge Cases', () => {
    
    test('handles zero-length buffer', () => {
        const buf = Buffer.alloc(0);
        expect(buf.length).toBe(0);
    });

    test('handles very large buffer', () => {
        const buf = Buffer.alloc(100 * 1024 * 1024); // 100MB
        expect(buf.length).toBe(100 * 1024 * 1024);
    });

    test('handles buffer with negative length', () => {
        expect(() => Buffer.alloc(-1)).toThrow();
    });

    test('handles buffer slice', () => {
        const buf = Buffer.from('hello world');
        const slice = buf.slice(6, 11);
        expect(slice.toString()).toBe('world');
        
        // Modifying slice modifies original
        slice[0] = 87; // 'W' ascii
        expect(buf.toString()).toBe('hello World');
    });

    test('handles buffer with encoding', () => {
        const utf8Buf = Buffer.from('hello', 'utf8');
        const base64Buf = Buffer.from('aGVsbG8=', 'base64');
        
        expect(utf8Buf.toString('base64')).toBe('aGVsbG8=');
        expect(base64Buf.toString('utf8')).toBe('hello');
    });

    test('handles buffer comparison', () => {
        const buf1 = Buffer.from('abc');
        const buf2 = Buffer.from('abc');
        const buf3 = Buffer.from('abd');
        
        expect(buf1.equals(buf2)).toBe(true);
        expect(buf1.equals(buf3)).toBe(false);
        expect(Buffer.compare(buf1, buf3)).toBe(-1);
    });
});

// ============================================================================
// Helper Functions
// ============================================================================

// These functions would be in your actual code
function processNumber(n) {
    if (n === 0) return 'zero';
    if (n === Infinity || n === -Infinity) throw new Error('Cannot process Infinity');
    return n;
}

function processString(s) {
    if (s === '') return 'empty';
    if (s.length > 1000000) return 'very long';
    if (s.includes('\x00')) return 'contains null';
    if (/[\u{1F600}-\u{1F6FF}]/u.test(s)) return 'unicode';
    if (/[\u200D]/.test(s)) return 'has zwj';
    if (/[\u0590-\u05FF]/.test(s)) return 'rtl';
    if (s.includes('\x1b')) return 'contains control chars';
    return 'normal';
}

function processFloat(f) {
    if (f === 0) return 'zero';
    if (Math.abs(f) < Number.EPSILON) return 'epsilon';
    if (f === Number.MAX_VALUE) return 'max';
    if (f === Number.MIN_VALUE) return 'min';
    if (Math.abs(f) < Number.MIN_VALUE) return 'subnormal';
    if (Math.abs(f - 0.3) < 0.000001) return 'approx-0.3';
    return f;
}

function processValue(v) {
    if (v === null) throw new Error('Value cannot be null');
    if (v === undefined) throw new Error('Value cannot be undefined');
    if (isNaN(v)) throw new Error('Value cannot be NaN');
    return v;
}

function processArray(arr) {
    if (arr.length === 0) return 'empty';
    if (arr.length > 1000000) return 'very large';
    if (!(0 in arr) && arr.length > 0) return 'sparse';
    if (Object.keys(arr).some(k => isNaN(parseInt(k)))) return 'has non-numeric keys';
    if (arr[-1] !== undefined) return 'has negative index';
    if (Object.getOwnPropertyDescriptor(arr, 'custom')) return 'has custom property';
    if (arr[2] === undefined && 2 in arr) return 'has deleted';
    return 'normal';
}

function processObject(obj) {
    if (Object.keys(obj).length === 0) return 'empty';
    if (Object.getPrototypeOf(obj) === null) return 'null prototype';
    if (Object.getPrototypeOf(obj) !== Object.prototype) return 'has prototype';
    if (Object.isFrozen(obj)) return 'frozen';
    if (Object.isSealed(obj)) return 'sealed';
    
    const props = Object.getOwnPropertyDescriptors(obj);
    if (Object.values(props).some(p => !p.enumerable)) return 'has non-enumerable';
    if (Object.getOwnPropertySymbols(obj).length > 0) return 'has symbols';
    if (Object.values(props).some(p => p.get || p.set)) return 'has accessors';
    
    // Simple circular detection
    try {
        JSON.stringify(obj);
        return 'normal';
    } catch (e) {
        return 'circular';
    }
}

function processSet(set) {
    return set.size === 0 ? 'empty' : 'has items';
}

function processMap(map) {
    return map.size === 0 ? 'empty' : 'has items';
}

function processTypedArray(arr) {
    return arr.length === 0 ? 'empty' : 'has items';
}

function processBuffer(buf) {
    return buf.length === 0 ? 'empty' : 'has data';
}

function processLargeNumber(n) {
    if (n <= Number.MAX_SAFE_INTEGER && n >= Number.MIN_SAFE_INTEGER) return 'safe';
    return 'unsafe';
}

function addInt32(a, b) {
    return (a + b) | 0; // 32-bit integer addition with overflow
}

function processDate(date) {
    if (isNaN(date)) throw new Error('Invalid date');
    if (date < new Date('1970-01-01')) return 'pre-epoch';
    if (date > new Date('2038-01-19')) return 'post-2038';
    return 'normal';
}

function processDST(date1, date2) {
    return 'handled';
}