/**
 * Node.js Benchmark Tests
 * 
 * This file contains benchmarks for various operations to measure performance
 * and validate the test generator's benchmark support.
 */

// ============================================================================
// Setup and Utilities
// ============================================================================

const { performance, PerformanceObserver } = require('perf_hooks');
const { promisify } = require('util');
const fs = require('fs');
const path = require('path');

// Simple benchmarking utility
class Benchmark {
    constructor(name) {
        this.name = name;
        this.results = [];
    }

    async run(iterations, fn, ...args) {
        const times = [];
        const start = performance.now();
        
        for (let i = 0; i < iterations; i++) {
            const iterationStart = performance.now();
            await fn(...args);
            const iterationEnd = performance.now();
            times.push(iterationEnd - iterationStart);
        }
        
        const end = performance.now();
        const total = end - start;
        
        const result = {
            name: this.name,
            iterations,
            total,
            average: total / iterations,
            min: Math.min(...times),
            max: Math.max(...times),
            opsPerSecond: (iterations / total) * 1000
        };
        
        this.results.push(result);
        return result;
    }

    report() {
        console.log(`\n📊 Benchmark Results: ${this.name}`);
        console.log('='.repeat(60));
        this.results.forEach(r => {
            console.log(`  Iterations: ${r.iterations}`);
            console.log(`  Total time: ${r.total.toFixed(2)}ms`);
            console.log(`  Average: ${r.average.toFixed(3)}ms/op`);
            console.log(`  Min/Max: ${r.min.toFixed(3)}ms / ${r.max.toFixed(3)}ms`);
            console.log(`  Ops/sec: ${r.opsPerSecond.toFixed(0)}`);
            console.log('-'.repeat(40));
        });
    }
}

// ============================================================================
// String Operation Benchmarks
// ============================================================================

describe('String Operations Benchmarks', () => {
    const iterations = 10000;
    const testString = 'The quick brown fox jumps over the lazy dog';
    const parts = ['hello', 'world', 'this', 'is', 'a', 'test', 'string'];
    
    test('String concatenation - plus operator', () => {
        const benchmark = new Benchmark('String Concat (+)');
        
        for (let i = 0; i < iterations; i++) {
            let result = '';
            for (const part of parts) {
                result += part;
            }
            expect(result).toBeDefined();
        }
        
        benchmark.report();
    });
    
    test('String concatenation - join method', () => {
        const benchmark = new Benchmark('String Join');
        
        for (let i = 0; i < iterations; i++) {
            const result = parts.join('');
            expect(result).toBeDefined();
        }
        
        benchmark.report();
    });
    
    test('String concatenation - template literals', () => {
        const benchmark = new Benchmark('Template Literals');
        
        for (let i = 0; i < iterations; i++) {
            const result = `${parts[0]}${parts[1]}${parts[2]}${parts[3]}`;
            expect(result).toBeDefined();
        }
        
        benchmark.report();
    });
    
    test('String concatenation - concat method', () => {
        const benchmark = new Benchmark('String.concat()');
        
        for (let i = 0; i < iterations; i++) {
            let result = '';
            for (const part of parts) {
                result = result.concat(part);
            }
            expect(result).toBeDefined();
        }
        
        benchmark.report();
    });
    
    test('String manipulation - substring', () => {
        const benchmark = new Benchmark('String.substring()');
        
        for (let i = 0; i < iterations; i++) {
            const result = testString.substring(10, 20);
            expect(result).toBe('brown fox');
        }
        
        benchmark.report();
    });
    
    test('String manipulation - slice', () => {
        const benchmark = new Benchmark('String.slice()');
        
        for (let i = 0; i < iterations; i++) {
            const result = testString.slice(10, 20);
            expect(result).toBe('brown fox');
        }
        
        benchmark.report();
    });
    
    test('String manipulation - split/join', () => {
        const benchmark = new Benchmark('String.split().join()');
        
        for (let i = 0; i < iterations; i++) {
            const result = testString.split(' ').join('-');
            expect(result).toBe('The-quick-brown-fox-jumps-over-the-lazy-dog');
        }
        
        benchmark.report();
    });
    
    test('String manipulation - replace', () => {
        const benchmark = new Benchmark('String.replace()');
        
        for (let i = 0; i < iterations; i++) {
            const result = testString.replace(/o/g, '0');
            expect(result).toContain('0');
        }
        
        benchmark.report();
    });
    
    test('String manipulation - replaceAll', () => {
        const benchmark = new Benchmark('String.replaceAll()');
        
        for (let i = 0; i < iterations; i++) {
            const result = testString.replaceAll('o', '0');
            expect(result).toContain('0');
        }
        
        benchmark.report();
    });
});

// ============================================================================
// Array Operation Benchmarks
// ============================================================================

describe('Array Operations Benchmarks', () => {
    const iterations = 10000;
    const smallArray = Array.from({ length: 100 }, (_, i) => i);
    const mediumArray = Array.from({ length: 1000 }, (_, i) => i);
    const largeArray = Array.from({ length: 10000 }, (_, i) => i);
    
    test('Array creation - literal', () => {
        const benchmark = new Benchmark('Array Literal');
        
        for (let i = 0; i < iterations; i++) {
            const arr = [1, 2, 3, 4, 5];
            expect(arr.length).toBe(5);
        }
        
        benchmark.report();
    });
    
    test('Array creation - constructor', () => {
        const benchmark = new Benchmark('Array Constructor');
        
        for (let i = 0; i < iterations; i++) {
            const arr = new Array(5);
            expect(arr.length).toBe(5);
        }
        
        benchmark.report();
    });
    
    test('Array creation - from', () => {
        const benchmark = new Benchmark('Array.from()');
        
        for (let i = 0; i < iterations; i++) {
            const arr = Array.from({ length: 5 }, (_, i) => i);
            expect(arr.length).toBe(5);
        }
        
        benchmark.report();
    });
    
    test('Array iteration - for loop', () => {
        const benchmark = new Benchmark('For Loop');
        
        for (let i = 0; i < iterations; i++) {
            let sum = 0;
            for (let j = 0; j < smallArray.length; j++) {
                sum += smallArray[j];
            }
            expect(sum).toBe(4950);
        }
        
        benchmark.report();
    });
    
    test('Array iteration - forEach', () => {
        const benchmark = new Benchmark('Array.forEach()');
        
        for (let i = 0; i < iterations; i++) {
            let sum = 0;
            smallArray.forEach(val => { sum += val; });
            expect(sum).toBe(4950);
        }
        
        benchmark.report();
    });
    
    test('Array iteration - for...of', () => {
        const benchmark = new Benchmark('For...of');
        
        for (let i = 0; i < iterations; i++) {
            let sum = 0;
            for (const val of smallArray) {
                sum += val;
            }
            expect(sum).toBe(4950);
        }
        
        benchmark.report();
    });
    
    test('Array iteration - map', () => {
        const benchmark = new Benchmark('Array.map()');
        
        for (let i = 0; i < iterations; i++) {
            const doubled = smallArray.map(x => x * 2);
            expect(doubled[0]).toBe(0);
        }
        
        benchmark.report();
    });
    
    test('Array iteration - filter', () => {
        const benchmark = new Benchmark('Array.filter()');
        
        for (let i = 0; i < iterations; i++) {
            const evens = smallArray.filter(x => x % 2 === 0);
            expect(evens.length).toBe(50);
        }
        
        benchmark.report();
    });
    
    test('Array iteration - reduce', () => {
        const benchmark = new Benchmark('Array.reduce()');
        
        for (let i = 0; i < iterations; i++) {
            const sum = smallArray.reduce((acc, val) => acc + val, 0);
            expect(sum).toBe(4950);
        }
        
        benchmark.report();
    });
    
    test('Array search - indexOf', () => {
        const benchmark = new Benchmark('Array.indexOf()');
        
        for (let i = 0; i < iterations; i++) {
            const index = smallArray.indexOf(50);
            expect(index).toBe(50);
        }
        
        benchmark.report();
    });
    
    test('Array search - includes', () => {
        const benchmark = new Benchmark('Array.includes()');
        
        for (let i = 0; i < iterations; i++) {
            const found = smallArray.includes(50);
            expect(found).toBe(true);
        }
        
        benchmark.report();
    });
    
    test('Array search - find', () => {
        const benchmark = new Benchmark('Array.find()');
        
        for (let i = 0; i < iterations; i++) {
            const found = smallArray.find(x => x === 50);
            expect(found).toBe(50);
        }
        
        benchmark.report();
    });
    
    test('Array push vs unshift', () => {
        test('push - end', () => {
            const benchmark = new Benchmark('Array.push()');
            const arr = [];
            
            for (let i = 0; i < iterations; i++) {
                arr.push(i);
            }
            
            expect(arr.length).toBe(iterations);
            benchmark.report();
        });
        
        test('unshift - beginning', () => {
            const benchmark = new Benchmark('Array.unshift()');
            const arr = [];
            
            for (let i = 0; i < iterations; i++) {
                arr.unshift(i);
            }
            
            expect(arr.length).toBe(iterations);
            benchmark.report();
        });
    });
});

// ============================================================================
// Object Operation Benchmarks
// ============================================================================

describe('Object Operations Benchmarks', () => {
    const iterations = 10000;
    const testObj = { a: 1, b: 2, c: 3, d: 4, e: 5 };
    const largeObj = Object.fromEntries(
        Array.from({ length: 1000 }, (_, i) => [`key${i}`, i])
    );
    
    test('Object creation - literal', () => {
        const benchmark = new Benchmark('Object Literal');
        
        for (let i = 0; i < iterations; i++) {
            const obj = { x: 1, y: 2, z: 3 };
            expect(obj.x).toBe(1);
        }
        
        benchmark.report();
    });
    
    test('Object creation - constructor', () => {
        const benchmark = new Benchmark('Object Constructor');
        
        for (let i = 0; i < iterations; i++) {
            const obj = new Object();
            obj.x = 1;
            obj.y = 2;
            obj.z = 3;
            expect(obj.x).toBe(1);
        }
        
        benchmark.report();
    });
    
    test('Object creation - Object.create()', () => {
        const benchmark = new Benchmark('Object.create()');
        
        for (let i = 0; i < iterations; i++) {
            const obj = Object.create(null);
            obj.x = 1;
            obj.y = 2;
            obj.z = 3;
            expect(obj.x).toBe(1);
        }
        
        benchmark.report();
    });
    
    test('Object property access - dot notation', () => {
        const benchmark = new Benchmark('Dot Notation');
        
        for (let i = 0; i < iterations; i++) {
            const val = testObj.a;
            expect(val).toBe(1);
        }
        
        benchmark.report();
    });
    
    test('Object property access - bracket notation', () => {
        const benchmark = new Benchmark('Bracket Notation');
        
        for (let i = 0; i < iterations; i++) {
            const val = testObj['a'];
            expect(val).toBe(1);
        }
        
        benchmark.report();
    });
    
    test('Object iteration - for...in', () => {
        const benchmark = new Benchmark('For...in');
        
        for (let i = 0; i < iterations; i++) {
            let sum = 0;
            for (const key in testObj) {
                sum += testObj[key];
            }
            expect(sum).toBe(15);
        }
        
        benchmark.report();
    });
    
    test('Object iteration - Object.keys()', () => {
        const benchmark = new Benchmark('Object.keys()');
        
        for (let i = 0; i < iterations; i++) {
            const keys = Object.keys(testObj);
            expect(keys.length).toBe(5);
        }
        
        benchmark.report();
    });
    
    test('Object iteration - Object.values()', () => {
        const benchmark = new Benchmark('Object.values()');
        
        for (let i = 0; i < iterations; i++) {
            const values = Object.values(testObj);
            expect(values.length).toBe(5);
        }
        
        benchmark.report();
    });
    
    test('Object iteration - Object.entries()', () => {
        const benchmark = new Benchmark('Object.entries()');
        
        for (let i = 0; i < iterations; i++) {
            const entries = Object.entries(testObj);
            expect(entries.length).toBe(5);
        }
        
        benchmark.report();
    });
});

// ============================================================================
// Function Call Benchmarks
// ============================================================================

describe('Function Call Benchmarks', () => {
    const iterations = 100000;
    
    function regularFunction(x) {
        return x * 2;
    }
    
    const arrowFunction = (x) => x * 2;
    
    class TestClass {
        method(x) {
            return x * 2;
        }
        
        static staticMethod(x) {
            return x * 2;
        }
    }
    
    test('Regular function call', () => {
        const benchmark = new Benchmark('Regular Function');
        
        for (let i = 0; i < iterations; i++) {
            const result = regularFunction(i);
            expect(result).toBe(i * 2);
        }
        
        benchmark.report();
    });
    
    test('Arrow function call', () => {
        const benchmark = new Benchmark('Arrow Function');
        
        for (let i = 0; i < iterations; i++) {
            const result = arrowFunction(i);
            expect(result).toBe(i * 2);
        }
        
        benchmark.report();
    });
    
    test('Method call', () => {
        const benchmark = new Benchmark('Method Call');
        const instance = new TestClass();
        
        for (let i = 0; i < iterations; i++) {
            const result = instance.method(i);
            expect(result).toBe(i * 2);
        }
        
        benchmark.report();
    });
    
    test('Static method call', () => {
        const benchmark = new Benchmark('Static Method');
        
        for (let i = 0; i < iterations; i++) {
            const result = TestClass.staticMethod(i);
            expect(result).toBe(i * 2);
        }
        
        benchmark.report();
    });
    
    test('Function call with apply', () => {
        const benchmark = new Benchmark('Function.apply()');
        
        for (let i = 0; i < iterations; i++) {
            const result = regularFunction.apply(null, [i]);
            expect(result).toBe(i * 2);
        }
        
        benchmark.report();
    });
    
    test('Function call with call', () => {
        const benchmark = new Benchmark('Function.call()');
        
        for (let i = 0; i < iterations; i++) {
            const result = regularFunction.call(null, i);
            expect(result).toBe(i * 2);
        }
        
        benchmark.report();
    });
    
    test('Function call with bind', () => {
        const benchmark = new Benchmark('Function.bind()');
        
        for (let i = 0; i < iterations; i++) {
            const bound = regularFunction.bind(null, i);
            const result = bound();
            expect(result).toBe(i * 2);
        }
        
        benchmark.report();
    });
});

// ============================================================================
// Async Operation Benchmarks
// ============================================================================

describe('Async Operations Benchmarks', () => {
    const iterations = 1000;
    
    function delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
    
    test('Promise creation', async () => {
        const benchmark = new Benchmark('Promise Creation');
        
        for (let i = 0; i < iterations; i++) {
            const promise = Promise.resolve(i);
            const result = await promise;
            expect(result).toBe(i);
        }
        
        benchmark.report();
    });
    
    test('Promise chain', async () => {
        const benchmark = new Benchmark('Promise Chain');
        
        for (let i = 0; i < iterations; i++) {
            const result = await Promise.resolve(i)
                .then(x => x * 2)
                .then(x => x + 1)
                .then(x => x.toString());
            expect(result).toBe(String(i * 2 + 1));
        }
        
        benchmark.report();
    });
    
    test('Async/await', async () => {
        const benchmark = new Benchmark('Async/Await');
        
        const asyncFunction = async (x) => {
            const doubled = x * 2;
            const added = doubled + 1;
            return added.toString();
        };
        
        for (let i = 0; i < iterations; i++) {
            const result = await asyncFunction(i);
            expect(result).toBe(String(i * 2 + 1));
        }
        
        benchmark.report();
    });
    
    test('Promise.all', async () => {
        const benchmark = new Benchmark('Promise.all');
        
        for (let i = 0; i < iterations; i++) {
            const promises = [Promise.resolve(i), Promise.resolve(i * 2)];
            const [a, b] = await Promise.all(promises);
            expect(a).toBe(i);
            expect(b).toBe(i * 2);
        }
        
        benchmark.report();
    });
    
    test('Promise.race', async () => {
        const benchmark = new Benchmark('Promise.race');
        
        for (let i = 0; i < iterations; i++) {
            const result = await Promise.race([
                Promise.resolve(i),
                delay(10).then(() => i * 2)
            ]);
            expect(result).toBe(i);
        }
        
        benchmark.report();
    });
});

// ============================================================================
// JSON Operation Benchmarks
// ============================================================================

describe('JSON Operations Benchmarks', () => {
    const iterations = 10000;
    const testData = {
        id: 42,
        name: 'Test Object',
        tags: ['node', 'javascript', 'benchmark'],
        active: true,
        score: 3.14159,
        metadata: {
            created: '2024-01-01',
            version: 1,
            source: 'benchmark'
        },
        nested: {
            array: [1, 2, 3, 4, 5],
            object: { a: 1, b: 2 }
        }
    };
    
    test('JSON.stringify', () => {
        const benchmark = new Benchmark('JSON.stringify()');
        
        for (let i = 0; i < iterations; i++) {
            const json = JSON.stringify(testData);
            expect(json).toContain('Test Object');
        }
        
        benchmark.report();
    });
    
    test('JSON.parse', () => {
        const benchmark = new Benchmark('JSON.parse()');
        const json = JSON.stringify(testData);
        
        for (let i = 0; i < iterations; i++) {
            const obj = JSON.parse(json);
            expect(obj.name).toBe('Test Object');
        }
        
        benchmark.report();
    });
    
    test('JSON round trip', () => {
        const benchmark = new Benchmark('JSON Round Trip');
        
        for (let i = 0; i < iterations; i++) {
            const json = JSON.stringify(testData);
            const obj = JSON.parse(json);
            expect(obj.name).toBe('Test Object');
        }
        
        benchmark.report();
    });
});

// ============================================================================
// File System Benchmarks
// ============================================================================

describe('File System Benchmarks', () => {
    const iterations = 100;
    const testFile = path.join(__dirname, 'test-file.txt');
    const testContent = 'Hello, World!\n'.repeat(1000);
    
    beforeAll(() => {
        fs.writeFileSync(testFile, testContent);
    });
    
    afterAll(() => {
        fs.unlinkSync(testFile);
    });
    
    test('fs.readFileSync', () => {
        const benchmark = new Benchmark('fs.readFileSync');
        
        for (let i = 0; i < iterations; i++) {
            const content = fs.readFileSync(testFile, 'utf8');
            expect(content.length).toBe(testContent.length);
        }
        
        benchmark.report();
    });
    
    test('fs.readFile (callback)', (done) => {
        const benchmark = new Benchmark('fs.readFile callback');
        let completed = 0;
        
        for (let i = 0; i < iterations; i++) {
            fs.readFile(testFile, 'utf8', (err, content) => {
                expect(err).toBeNull();
                expect(content.length).toBe(testContent.length);
                completed++;
                if (completed === iterations) {
                    benchmark.report();
                    done();
                }
            });
        }
    });
    
    test('fs.readFile (promise)', async () => {
        const benchmark = new Benchmark('fs.readFile promise');
        const readFile = promisify(fs.readFile);
        
        for (let i = 0; i < iterations; i++) {
            const content = await readFile(testFile, 'utf8');
            expect(content.length).toBe(testContent.length);
        }
        
        benchmark.report();
    });
});

// ============================================================================
// Data Structure Benchmarks
// ============================================================================

describe('Data Structure Benchmarks', () => {
    const iterations = 10000;
    const data = Array.from({ length: 1000 }, (_, i) => i);
    
    test('Set operations', () => {
        const benchmark = new Benchmark('Set Operations');
        
        for (let i = 0; i < iterations; i++) {
            const set = new Set(data);
            expect(set.has(500)).toBe(true);
            expect(set.size).toBe(1000);
        }
        
        benchmark.report();
    });
    
    test('Map operations', () => {
        const benchmark = new Benchmark('Map Operations');
        
        for (let i = 0; i < iterations; i++) {
            const map = new Map(data.map(x => [x, x * 2]));
            expect(map.get(500)).toBe(1000);
            expect(map.size).toBe(1000);
        }
        
        benchmark.report();
    });
    
    test('WeakSet operations', () => {
        const benchmark = new Benchmark('WeakSet Operations');
        const objects = data.map(x => ({ id: x }));
        
        for (let i = 0; i < iterations; i++) {
            const weakSet = new WeakSet(objects);
            expect(weakSet.has(objects[500])).toBe(true);
        }
        
        benchmark.report();
    });
    
    test('WeakMap operations', () => {
        const benchmark = new Benchmark('WeakMap Operations');
        const objects = data.map(x => ({ id: x }));
        const values = data.map(x => x * 2);
        
        for (let i = 0; i < iterations; i++) {
            const weakMap = new WeakMap();
            objects.forEach((obj, idx) => weakMap.set(obj, values[idx]));
            expect(weakMap.get(objects[500])).toBe(1000);
        }
        
        benchmark.report();
    });
});

// ============================================================================
// Regular Expression Benchmarks
// ============================================================================

describe('Regular Expression Benchmarks', () => {
    const iterations = 10000;
    const text = 'The quick brown fox jumps over the lazy dog. '.repeat(100);
    
    test('Simple regex match', () => {
        const benchmark = new Benchmark('Simple Regex Match');
        
        for (let i = 0; i < iterations; i++) {
            const matches = text.match(/fox/g);
            expect(matches.length).toBe(100);
        }
        
        benchmark.report();
    });
    
    test('Complex regex match', () => {
        const benchmark = new Benchmark('Complex Regex Match');
        const regex = /\b\w{5}\b/g; // 5-letter words
        
        for (let i = 0; i < iterations; i++) {
            const matches = text.match(regex);
            expect(matches.length).toBeGreaterThan(0);
        }
        
        benchmark.report();
    });
    
    test('Regex replace', () => {
        const benchmark = new Benchmark('Regex Replace');
        
        for (let i = 0; i < iterations; i++) {
            const result = text.replace(/fox/g, 'cat');
            expect(result).toContain('cat');
        }
        
        benchmark.report();
    });
    
    test('Regex test', () => {
        const benchmark = new Benchmark('Regex.test()');
        const regex = /fox/;
        
        for (let i = 0; i < iterations; i++) {
            const found = regex.test(text);
            expect(found).toBe(true);
        }
        
        benchmark.report();
    });
});

// ============================================================================
// Math Operation Benchmarks
// ============================================================================

describe('Math Operations Benchmarks', () => {
    const iterations = 100000;
    
    test('Math.abs', () => {
        const benchmark = new Benchmark('Math.abs()');
        
        for (let i = 0; i < iterations; i++) {
            const result = Math.abs(-i);
            expect(result).toBe(i);
        }
        
        benchmark.report();
    });
    
    test('Math.floor vs Math.trunc', () => {
        test('Math.floor', () => {
            const benchmark = new Benchmark('Math.floor()');
            
            for (let i = 0; i < iterations; i++) {
                const result = Math.floor(i / 3.7);
                expect(result).toBeDefined();
            }
            
            benchmark.report();
        });
        
        test('Math.trunc', () => {
            const benchmark = new Benchmark('Math.trunc()');
            
            for (let i = 0; i < iterations; i++) {
                const result = Math.trunc(i / 3.7);
                expect(result).toBeDefined();
            }
            
            benchmark.report();
        });
    });
    
    test('Math.random', () => {
        const benchmark = new Benchmark('Math.random()');
        
        for (let i = 0; i < iterations; i++) {
            const rand = Math.random();
            expect(rand).toBeGreaterThanOrEqual(0);
            expect(rand).toBeLessThan(1);
        }
        
        benchmark.report();
    });
    
    test('Math.pow vs exponentiation', () => {
        test('Math.pow', () => {
            const benchmark = new Benchmark('Math.pow()');
            
            for (let i = 0; i < iterations; i++) {
                const result = Math.pow(i, 2);
                expect(result).toBe(i * i);
            }
            
            benchmark.report();
        });
        
        test('** operator', () => {
            const benchmark = new Benchmark('** Operator');
            
            for (let i = 0; i < iterations; i++) {
                const result = i ** 2;
                expect(result).toBe(i * i);
            }
            
            benchmark.report();
        });
    });
});