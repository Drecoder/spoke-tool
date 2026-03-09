/**
 * Math Operations Tests
 * 
 * Tests for mathematical functions including basic arithmetic,
 * advanced math, statistics, and edge cases.
 */

const {
    // Basic arithmetic
    add,
    subtract,
    multiply,
    divide,
    
    // Advanced math
    power,
    sqrt,
    cbrt,
    abs,
    round,
    floor,
    ceil,
    trunc,
    
    // Trigonometry
    sin,
    cos,
    tan,
    asin,
    acos,
    atan,
    
    // Logarithms
    log,
    log10,
    log2,
    exp,
    
    // Statistics
    sum,
    average,
    median,
    mode,
    variance,
    standardDeviation,
    
    // Combinatorics
    factorial,
    fibonacci,
    gcd,
    lcm,
    isPrime,
    
    // Utilities
    clamp,
    lerp,
    remap,
    random,
    randomInt,
    roundTo
} = require('../src/math');

// ============================================================================
// Basic Arithmetic Tests
// ============================================================================

describe('Basic Arithmetic', () => {
    
    describe('add', () => {
        test('should add positive numbers correctly', () => {
            expect(add(2, 3)).toBe(5);
            expect(add(10, 20)).toBe(30);
            expect(add(1.5, 2.5)).toBe(4.0);
        });

        test('should add negative numbers correctly', () => {
            expect(add(-2, -3)).toBe(-5);
            expect(add(-10, 5)).toBe(-5);
            expect(add(10, -5)).toBe(5);
        });

        test('should add zero correctly', () => {
            expect(add(0, 5)).toBe(5);
            expect(add(5, 0)).toBe(5);
            expect(add(0, 0)).toBe(0);
        });

        test('should handle decimal precision', () => {
            expect(add(0.1, 0.2)).toBeCloseTo(0.3, 12);
            expect(add(1.234, 2.345)).toBeCloseTo(3.579, 12);
        });

        test('should handle large numbers', () => {
            expect(add(1e10, 2e10)).toBe(3e10);
            expect(add(Number.MAX_SAFE_INTEGER, 1)).toBe(Number.MAX_SAFE_INTEGER + 1);
        });

        test('should handle multiple arguments', () => {
            expect(add(1, 2, 3, 4, 5)).toBe(15);
            expect(add(10)).toBe(10);
            expect(add()).toBe(0);
        });
    });

    describe('subtract', () => {
        test('should subtract positive numbers correctly', () => {
            expect(subtract(10, 4)).toBe(6);
            expect(subtract(4, 10)).toBe(-6);
            expect(subtract(1.5, 0.5)).toBe(1.0);
        });

        test('should subtract negative numbers correctly', () => {
            expect(subtract(-10, -4)).toBe(-6);
            expect(subtract(-10, 4)).toBe(-14);
            expect(subtract(10, -4)).toBe(14);
        });

        test('should subtract zero correctly', () => {
            expect(subtract(5, 0)).toBe(5);
            expect(subtract(0, 5)).toBe(-5);
            expect(subtract(0, 0)).toBe(0);
        });

        test('should handle decimal precision', () => {
            expect(subtract(5.5, 2.2)).toBeCloseTo(3.3, 12);
            expect(subtract(0.3, 0.1)).toBeCloseTo(0.2, 12);
        });

        test('should handle multiple arguments', () => {
            expect(subtract(10, 2, 3)).toBe(5);
            expect(subtract(100, 10, 20, 30)).toBe(40);
        });
    });

    describe('multiply', () => {
        test('should multiply positive numbers correctly', () => {
            expect(multiply(2, 3)).toBe(6);
            expect(multiply(10, 0.5)).toBe(5);
            expect(multiply(1.5, 2)).toBe(3);
        });

        test('should multiply negative numbers correctly', () => {
            expect(multiply(-2, 3)).toBe(-6);
            expect(multiply(2, -3)).toBe(-6);
            expect(multiply(-2, -3)).toBe(6);
        });

        test('should multiply by zero correctly', () => {
            expect(multiply(5, 0)).toBe(0);
            expect(multiply(0, 5)).toBe(0);
            expect(multiply(0, 0)).toBe(0);
        });

        test('should multiply by one correctly', () => {
            expect(multiply(5, 1)).toBe(5);
            expect(multiply(1, 5)).toBe(5);
        });

        test('should handle decimal precision', () => {
            expect(multiply(0.1, 0.1)).toBeCloseTo(0.01, 12);
            expect(multiply(2.5, 1.5)).toBeCloseTo(3.75, 12);
        });

        test('should handle large numbers', () => {
            expect(multiply(1e6, 1e6)).toBe(1e12);
        });

        test('should handle multiple arguments', () => {
            expect(multiply(2, 3, 4)).toBe(24);
            expect(multiply(2, 3, 4, 5)).toBe(120);
        });
    });

    describe('divide', () => {
        test('should divide positive numbers correctly', () => {
            expect(divide(10, 2)).toBe(5);
            expect(divide(9, 3)).toBe(3);
            expect(divide(7, 2)).toBe(3.5);
        });

        test('should divide negative numbers correctly', () => {
            expect(divide(-10, 2)).toBe(-5);
            expect(divide(10, -2)).toBe(-5);
            expect(divide(-10, -2)).toBe(5);
        });

        test('should divide zero correctly', () => {
            expect(divide(0, 5)).toBe(0);
        });

        test('should throw error when dividing by zero', () => {
            expect(() => divide(10, 0)).toThrow('Division by zero');
            expect(() => divide(0, 0)).toThrow('Division by zero');
        });

        test('should handle decimal precision', () => {
            expect(divide(1, 3)).toBeCloseTo(0.3333333333, 10);
            expect(divide(0.1, 0.2)).toBeCloseTo(0.5, 12);
        });

        test('should handle multiple arguments', () => {
            expect(divide(100, 2, 5)).toBe(10);
            expect(divide(100, 2, 2, 5)).toBe(5);
        });
    });
});

// ============================================================================
// Advanced Math Tests
// ============================================================================

describe('Advanced Math', () => {
    
    describe('power', () => {
        test('should calculate positive exponents', () => {
            expect(power(2, 3)).toBe(8);
            expect(power(5, 2)).toBe(25);
            expect(power(10, 0)).toBe(1);
            expect(power(2, 1)).toBe(2);
        });

        test('should calculate negative exponents', () => {
            expect(power(2, -1)).toBeCloseTo(0.5, 12);
            expect(power(4, -2)).toBeCloseTo(0.0625, 12);
        });

        test('should calculate fractional exponents', () => {
            expect(power(9, 0.5)).toBeCloseTo(3, 12);
            expect(power(27, 1/3)).toBeCloseTo(3, 12);
        });

        test('should handle zero base', () => {
            expect(power(0, 5)).toBe(0);
            expect(power(0, 0)).toBe(1); // Mathematical convention
        });

        test('should handle negative base', () => {
            expect(power(-2, 2)).toBe(4);
            expect(power(-2, 3)).toBe(-8);
        });

        test('should handle large exponents', () => {
            expect(power(2, 10)).toBe(1024);
            expect(power(2, 20)).toBe(1048576);
        });
    });

    describe('sqrt', () => {
        test('should calculate square root of perfect squares', () => {
            expect(sqrt(4)).toBe(2);
            expect(sqrt(9)).toBe(3);
            expect(sqrt(16)).toBe(4);
            expect(sqrt(25)).toBe(5);
            expect(sqrt(100)).toBe(10);
        });

        test('should calculate square root of non-perfect squares', () => {
            expect(sqrt(2)).toBeCloseTo(1.41421356, 8);
            expect(sqrt(3)).toBeCloseTo(1.73205081, 8);
            expect(sqrt(5)).toBeCloseTo(2.23606798, 8);
        });

        test('should handle zero', () => {
            expect(sqrt(0)).toBe(0);
        });

        test('should handle negative numbers', () => {
            expect(() => sqrt(-1)).toThrow('Cannot calculate square root of negative number');
        });

        test('should handle large numbers', () => {
            expect(sqrt(1e10)).toBeCloseTo(100000, 0);
        });
    });

    describe('cbrt', () => {
        test('should calculate cube root of perfect cubes', () => {
            expect(cbrt(8)).toBe(2);
            expect(cbrt(27)).toBe(3);
            expect(cbrt(64)).toBe(4);
            expect(cbrt(125)).toBe(5);
        });

        test('should calculate cube root of non-perfect cubes', () => {
            expect(cbrt(2)).toBeCloseTo(1.25992105, 8);
            expect(cbrt(10)).toBeCloseTo(2.15443469, 8);
        });

        test('should handle negative numbers', () => {
            expect(cbrt(-8)).toBe(-2);
            expect(cbrt(-27)).toBe(-3);
        });

        test('should handle zero', () => {
            expect(cbrt(0)).toBe(0);
        });
    });

    describe('abs', () => {
        test('should return absolute value of positive numbers', () => {
            expect(abs(5)).toBe(5);
            expect(abs(0)).toBe(0);
            expect(abs(3.14)).toBe(3.14);
        });

        test('should return absolute value of negative numbers', () => {
            expect(abs(-5)).toBe(5);
            expect(abs(-3.14)).toBe(3.14);
            expect(abs(-0)).toBe(0);
        });

        test('should handle edge cases', () => {
            expect(abs(Number.MAX_SAFE_INTEGER)).toBe(Number.MAX_SAFE_INTEGER);
            expect(abs(Number.MIN_SAFE_INTEGER)).toBe(Number.MAX_SAFE_INTEGER);
        });
    });

    describe('rounding functions', () => {
        test('round should round to nearest integer', () => {
            expect(round(3.2)).toBe(3);
            expect(round(3.5)).toBe(4);
            expect(round(3.8)).toBe(4);
            expect(round(-3.2)).toBe(-3);
            expect(round(-3.5)).toBe(-3); // Note: JavaScript rounds to nearest, ties away from zero
            expect(round(-3.8)).toBe(-4);
        });

        test('floor should round down', () => {
            expect(floor(3.2)).toBe(3);
            expect(floor(3.8)).toBe(3);
            expect(floor(-3.2)).toBe(-4);
            expect(floor(-3.8)).toBe(-4);
        });

        test('ceil should round up', () => {
            expect(ceil(3.2)).toBe(4);
            expect(ceil(3.8)).toBe(4);
            expect(ceil(-3.2)).toBe(-3);
            expect(ceil(-3.8)).toBe(-3);
        });

        test('trunc should remove decimal part', () => {
            expect(trunc(3.2)).toBe(3);
            expect(trunc(3.8)).toBe(3);
            expect(trunc(-3.2)).toBe(-3);
            expect(trunc(-3.8)).toBe(-3);
        });

        test('roundTo should round to specified precision', () => {
            expect(roundTo(3.14159, 2)).toBeCloseTo(3.14, 12);
            expect(roundTo(3.14159, 4)).toBeCloseTo(3.1416, 12);
            expect(roundTo(123.456, -1)).toBeCloseTo(120, 12);
        });
    });
});

// ============================================================================
// Trigonometry Tests
// ============================================================================

describe('Trigonometry', () => {
    
    describe('sin', () => {
        test('should calculate sine of common angles', () => {
            expect(sin(0)).toBeCloseTo(0, 12);
            expect(sin(Math.PI / 6)).toBeCloseTo(0.5, 12);
            expect(sin(Math.PI / 2)).toBeCloseTo(1, 12);
            expect(sin(Math.PI)).toBeCloseTo(0, 12);
            expect(sin(3 * Math.PI / 2)).toBeCloseTo(-1, 12);
        });

        test('should handle negative angles', () => {
            expect(sin(-Math.PI / 2)).toBeCloseTo(-1, 12);
        });

        test('should handle large angles', () => {
            expect(sin(100 * Math.PI)).toBeCloseTo(0, 10);
        });
    });

    describe('cos', () => {
        test('should calculate cosine of common angles', () => {
            expect(cos(0)).toBeCloseTo(1, 12);
            expect(cos(Math.PI / 3)).toBeCloseTo(0.5, 12);
            expect(cos(Math.PI / 2)).toBeCloseTo(0, 12);
            expect(cos(Math.PI)).toBeCloseTo(-1, 12);
        });

        test('should handle negative angles', () => {
            expect(cos(-Math.PI / 2)).toBeCloseTo(0, 12);
        });
    });

    describe('tan', () => {
        test('should calculate tangent of common angles', () => {
            expect(tan(0)).toBeCloseTo(0, 12);
            expect(tan(Math.PI / 4)).toBeCloseTo(1, 12);
            expect(tan(Math.PI / 3)).toBeCloseTo(1.73205081, 8);
        });

        test('should handle asymptotes', () => {
            expect(() => tan(Math.PI / 2)).toThrow(); // Should be undefined/infinite
        });
    });

    describe('inverse trig functions', () => {
        test('asin should calculate inverse sine', () => {
            expect(asin(0)).toBeCloseTo(0, 12);
            expect(asin(1)).toBeCloseTo(Math.PI / 2, 12);
            expect(asin(0.5)).toBeCloseTo(Math.PI / 6, 12);
        });

        test('asin should handle invalid inputs', () => {
            expect(() => asin(2)).toThrow();
            expect(() => asin(-2)).toThrow();
        });

        test('acos should calculate inverse cosine', () => {
            expect(acos(1)).toBeCloseTo(0, 12);
            expect(acos(0)).toBeCloseTo(Math.PI / 2, 12);
            expect(acos(0.5)).toBeCloseTo(Math.PI / 3, 12);
        });

        test('atan should calculate inverse tangent', () => {
            expect(atan(0)).toBeCloseTo(0, 12);
            expect(atan(1)).toBeCloseTo(Math.PI / 4, 12);
            expect(atan(Infinity)).toBeCloseTo(Math.PI / 2, 12);
        });
    });
});

// ============================================================================
// Logarithm Tests
// ============================================================================

describe('Logarithms', () => {
    
    describe('log', () => {
        test('should calculate natural logarithm', () => {
            expect(log(1)).toBeCloseTo(0, 12);
            expect(log(Math.E)).toBeCloseTo(1, 12);
            expect(log(10)).toBeCloseTo(2.30258509, 8);
        });

        test('should handle invalid inputs', () => {
            expect(() => log(0)).toThrow();
            expect(() => log(-1)).toThrow();
        });
    });

    describe('log10', () => {
        test('should calculate base-10 logarithm', () => {
            expect(log10(1)).toBeCloseTo(0, 12);
            expect(log10(10)).toBeCloseTo(1, 12);
            expect(log10(100)).toBeCloseTo(2, 12);
            expect(log10(1000)).toBeCloseTo(3, 12);
        });

        test('should handle invalid inputs', () => {
            expect(() => log10(0)).toThrow();
            expect(() => log10(-1)).toThrow();
        });
    });

    describe('log2', () => {
        test('should calculate base-2 logarithm', () => {
            expect(log2(1)).toBeCloseTo(0, 12);
            expect(log2(2)).toBeCloseTo(1, 12);
            expect(log2(4)).toBeCloseTo(2, 12);
            expect(log2(8)).toBeCloseTo(3, 12);
            expect(log2(16)).toBeCloseTo(4, 12);
        });

        test('should handle invalid inputs', () => {
            expect(() => log2(0)).toThrow();
            expect(() => log2(-1)).toThrow();
        });
    });

    describe('exp', () => {
        test('should calculate exponential function', () => {
            expect(exp(0)).toBeCloseTo(1, 12);
            expect(exp(1)).toBeCloseTo(Math.E, 12);
            expect(exp(2)).toBeCloseTo(7.3890561, 8);
        });

        test('should handle negative arguments', () => {
            expect(exp(-1)).toBeCloseTo(1 / Math.E, 12);
        });
    });
});

// ============================================================================
// Statistics Tests
// ============================================================================

describe('Statistics', () => {
    
    describe('sum', () => {
        test('should calculate sum of numbers', () => {
            expect(sum([1, 2, 3, 4, 5])).toBe(15);
            expect(sum([-1, -2, -3])).toBe(-6);
            expect(sum([1.5, 2.5, 3.5])).toBe(7.5);
        });

        test('should handle empty array', () => {
            expect(sum([])).toBe(0);
        });

        test('should handle single element', () => {
            expect(sum([42])).toBe(42);
        });
    });

    describe('average', () => {
        test('should calculate average of numbers', () => {
            expect(average([1, 2, 3, 4, 5])).toBe(3);
            expect(average([10, 20, 30])).toBe(20);
            expect(average([1.5, 2.5, 3.5])).toBe(2.5);
        });

        test('should handle empty array', () => {
            expect(() => average([])).toThrow();
        });

        test('should handle single element', () => {
            expect(average([42])).toBe(42);
        });
    });

    describe('median', () => {
        test('should calculate median of odd-length array', () => {
            expect(median([1, 3, 5])).toBe(3);
            expect(median([10, 20, 30, 40, 50])).toBe(30);
        });

        test('should calculate median of even-length array', () => {
            expect(median([1, 2, 3, 4])).toBe(2.5);
            expect(median([10, 20, 30, 40])).toBe(25);
        });

        test('should handle unsorted arrays', () => {
            expect(median([5, 1, 4, 2, 3])).toBe(3);
            expect(median([5, 1, 4, 2])).toBe(3);
        });

        test('should handle empty array', () => {
            expect(() => median([])).toThrow();
        });
    });

    describe('mode', () => {
        test('should find most frequent value', () => {
            expect(mode([1, 2, 2, 3, 4])).toBe(2);
            expect(mode([1, 1, 2, 2, 3])).toEqual([1, 2]); // Multiple modes
        });

        test('should handle all unique values', () => {
            expect(mode([1, 2, 3, 4])).toBeNull(); // No mode
        });

        test('should handle single element', () => {
            expect(mode([42])).toBe(42);
        });
    });

    describe('variance', () => {
        test('should calculate population variance', () => {
            expect(variance([1, 2, 3, 4, 5])).toBe(2);
            expect(variance([10, 20, 30, 40])).toBe(125);
        });

        test('should handle sample variance', () => {
            expect(variance([1, 2, 3, 4, 5], true)).toBe(2.5);
        });

        test('should handle empty array', () => {
            expect(() => variance([])).toThrow();
        });

        test('should handle single element', () => {
            expect(variance([42])).toBe(0);
        });
    });

    describe('standardDeviation', () => {
        test('should calculate population standard deviation', () => {
            expect(standardDeviation([1, 2, 3, 4, 5])).toBeCloseTo(1.41421356, 8);
            expect(standardDeviation([10, 20, 30, 40])).toBeCloseTo(11.1803399, 8);
        });

        test('should handle sample standard deviation', () => {
            expect(standardDeviation([1, 2, 3, 4, 5], true)).toBeCloseTo(1.58113883, 8);
        });
    });
});

// ============================================================================
// Combinatorics Tests
// ============================================================================

describe('Combinatorics', () => {
    
    describe('factorial', () => {
        test('should calculate factorial of small numbers', () => {
            expect(factorial(0)).toBe(1);
            expect(factorial(1)).toBe(1);
            expect(factorial(2)).toBe(2);
            expect(factorial(3)).toBe(6);
            expect(factorial(4)).toBe(24);
            expect(factorial(5)).toBe(120);
            expect(factorial(6)).toBe(720);
            expect(factorial(7)).toBe(5040);
            expect(factorial(8)).toBe(40320);
            expect(factorial(9)).toBe(362880);
            expect(factorial(10)).toBe(3628800);
        });

        test('should throw error for negative numbers', () => {
            expect(() => factorial(-1)).toThrow();
        });

        test('should handle large numbers', () => {
            expect(factorial(20)).toBe(2432902008176640000);
        });
    });

    describe('fibonacci', () => {
        test('should calculate Fibonacci sequence', () => {
            expect(fibonacci(0)).toBe(0);
            expect(fibonacci(1)).toBe(1);
            expect(fibonacci(2)).toBe(1);
            expect(fibonacci(3)).toBe(2);
            expect(fibonacci(4)).toBe(3);
            expect(fibonacci(5)).toBe(5);
            expect(fibonacci(6)).toBe(8);
            expect(fibonacci(7)).toBe(13);
            expect(fibonacci(8)).toBe(21);
            expect(fibonacci(9)).toBe(34);
            expect(fibonacci(10)).toBe(55);
        });

        test('should throw error for negative numbers', () => {
            expect(() => fibonacci(-1)).toThrow();
        });

        test('should handle large n', () => {
            expect(fibonacci(20)).toBe(6765);
            expect(fibonacci(30)).toBe(832040);
        });
    });

    describe('gcd', () => {
        test('should calculate greatest common divisor', () => {
            expect(gcd(48, 18)).toBe(6);
            expect(gcd(12, 8)).toBe(4);
            expect(gcd(17, 19)).toBe(1);
            expect(gcd(0, 5)).toBe(5);
            expect(gcd(5, 0)).toBe(5);
        });

        test('should handle negative numbers', () => {
            expect(gcd(-48, 18)).toBe(6);
            expect(gcd(48, -18)).toBe(6);
            expect(gcd(-48, -18)).toBe(6);
        });
    });

    describe('lcm', () => {
        test('should calculate least common multiple', () => {
            expect(lcm(12, 18)).toBe(36);
            expect(lcm(8, 12)).toBe(24);
            expect(lcm(17, 19)).toBe(323);
            expect(lcm(0, 5)).toBe(0);
        });

        test('should handle negative numbers', () => {
            expect(lcm(-12, 18)).toBe(36);
        });
    });

    describe('isPrime', () => {
        test('should identify prime numbers', () => {
            expect(isPrime(2)).toBe(true);
            expect(isPrime(3)).toBe(true);
            expect(isPrime(5)).toBe(true);
            expect(isPrime(7)).toBe(true);
            expect(isPrime(11)).toBe(true);
            expect(isPrime(13)).toBe(true);
            expect(isPrime(17)).toBe(true);
            expect(isPrime(19)).toBe(true);
            expect(isPrime(23)).toBe(true);
            expect(isPrime(29)).toBe(true);
            expect(isPrime(31)).toBe(true);
            expect(isPrime(37)).toBe(true);
            expect(isPrime(41)).toBe(true);
            expect(isPrime(43)).toBe(true);
            expect(isPrime(47)).toBe(true);
        });

        test('should identify composite numbers', () => {
            expect(isPrime(1)).toBe(false);
            expect(isPrime(4)).toBe(false);
            expect(isPrime(6)).toBe(false);
            expect(isPrime(8)).toBe(false);
            expect(isPrime(9)).toBe(false);
            expect(isPrime(10)).toBe(false);
            expect(isPrime(15)).toBe(false);
            expect(isPrime(21)).toBe(false);
            expect(isPrime(25)).toBe(false);
            expect(isPrime(27)).toBe(false);
            expect(isPrime(33)).toBe(false);
            expect(isPrime(35)).toBe(false);
            expect(isPrime(39)).toBe(false);
            expect(isPrime(45)).toBe(false);
            expect(isPrime(49)).toBe(false);
        });

        test('should handle negative numbers', () => {
            expect(isPrime(-2)).toBe(false);
            expect(isPrime(-3)).toBe(false);
        });

        test('should handle zero', () => {
            expect(isPrime(0)).toBe(false);
        });

        test('should handle large primes', () => {
            expect(isPrime(997)).toBe(true);
            expect(isPrime(99991)).toBe(true);
            expect(isPrime(100003)).toBe(true);
        });

        test('should handle large composites', () => {
            expect(isPrime(99999)).toBe(false);
            expect(isPrime(100000)).toBe(false);
        });
    });
});

// ============================================================================
// Utility Function Tests
// ============================================================================

describe('Math Utilities', () => {
    
    describe('clamp', () => {
        test('should clamp values within range', () => {
            expect(clamp(5, 1, 10)).toBe(5);
            expect(clamp(0, 1, 10)).toBe(1);
            expect(clamp(15, 1, 10)).toBe(10);
        });

        test('should handle negative ranges', () => {
            expect(clamp(-5, -10, -1)).toBe(-5);
            expect(clamp(-15, -10, -1)).toBe(-10);
            expect(clamp(0, -10, -1)).toBe(-1);
        });

        test('should handle equal min/max', () => {
            expect(clamp(5, 5, 5)).toBe(5);
            expect(clamp(10, 5, 5)).toBe(5);
        });
    });

    describe('lerp', () => {
        test('should linearly interpolate', () => {
            expect(lerp(0, 10, 0)).toBe(0);
            expect(lerp(0, 10, 0.5)).toBe(5);
            expect(lerp(0, 10, 1)).toBe(10);
        });

        test('should handle negative interpolation', () => {
            expect(lerp(0, 10, -1)).toBe(-10);
        });

        test('should handle extrapolation', () => {
            expect(lerp(0, 10, 2)).toBe(20);
        });

        test('should handle decimal values', () => {
            expect(lerp(1.5, 3.5, 0.5)).toBe(2.5);
        });
    });

    describe('remap', () => {
        test('should remap values from one range to another', () => {
            expect(remap(5, 0, 10, 0, 100)).toBe(50);
            expect(remap(0, -10, 10, 0, 100)).toBe(50);
            expect(remap(10, 0, 10, 0, 100)).toBe(100);
        });

        test('should handle clamped remapping', () => {
            expect(remap(15, 0, 10, 0, 100, true)).toBe(100);
            expect(remap(-5, 0, 10, 0, 100, true)).toBe(0);
        });

        test('should handle inverse ranges', () => {
            expect(remap(5, 0, 10, 100, 0)).toBe(50);
        });
    });

    describe('random', () => {
        test('should generate random number between 0 and 1', () => {
            for (let i = 0; i < 100; i++) {
                const r = random();
                expect(r).toBeGreaterThanOrEqual(0);
                expect(r).toBeLessThan(1);
            }
        });

        test('should generate random number in range', () => {
            for (let i = 0; i < 100; i++) {
                const r = random(5, 10);
                expect(r).toBeGreaterThanOrEqual(5);
                expect(r).toBeLessThan(10);
            }
        });

        test('should handle reversed range', () => {
            for (let i = 0; i < 100; i++) {
                const r = random(10, 5);
                expect(r).toBeGreaterThanOrEqual(5);
                expect(r).toBeLessThan(10);
            }
        });
    });

    describe('randomInt', () => {
        test('should generate random integer in range', () => {
            for (let i = 0; i < 100; i++) {
                const r = randomInt(1, 10);
                expect(Number.isInteger(r)).toBe(true);
                expect(r).toBeGreaterThanOrEqual(1);
                expect(r).toBeLessThanOrEqual(10);
            }
        });

        test('should handle inclusive min/max', () => {
            const values = new Set();
            for (let i = 0; i < 100; i++) {
                values.add(randomInt(1, 2));
            }
            expect(values.has(1)).toBe(true);
            expect(values.has(2)).toBe(true);
        });

        test('should handle single value range', () => {
            expect(randomInt(5, 5)).toBe(5);
        });
    });
});

// ============================================================================
// Edge Cases and Error Handling
// ============================================================================

describe('Edge Cases and Error Handling', () => {
    
    test('should handle NaN inputs', () => {
        expect(() => add(NaN, 5)).toThrow();
        expect(() => sqrt(NaN)).toThrow();
        expect(() => log(NaN)).toThrow();
    });

    test('should handle infinite inputs', () => {
        expect(() => add(Infinity, 5)).toThrow();
        expect(() => subtract(Infinity, 5)).toThrow();
        expect(() => multiply(Infinity, 5)).toThrow();
        expect(() => divide(Infinity, 5)).toThrow();
    });

    test('should handle very large numbers', () => {
        expect(add(1e308, 1e308)).toBe(Infinity); // Overflow
        expect(multiply(1e200, 1e200)).toBe(1e400);
    });

    test('should handle very small numbers', () => {
        expect(add(1e-308, 1e-308)).toBeCloseTo(2e-308, 308);
        expect(multiply(1e-200, 1e-200)).toBe(1e-400);
    });

    test('should handle precision edge cases', () => {
        expect(add(0.1, 0.2)).not.toBe(0.3); // Floating point precision
        expect(add(0.1, 0.2)).toBeCloseTo(0.3, 15);
    });

    test('should handle integer overflow', () => {
        const maxInt = Number.MAX_SAFE_INTEGER;
        expect(add(maxInt, 1)).toBe(maxInt + 1); // Becomes unsafe integer
        expect(Number.isSafeInteger(maxInt + 1)).toBe(false);
    });

    test('should handle division by very small numbers', () => {
        expect(divide(1, 1e-308)).toBe(1e308);
        expect(divide(1, 0)).toThrow();
    });
});

// ============================================================================
// Performance Tests
// ============================================================================

describe('Math Performance', () => {
    
    test('should calculate factorial efficiently', () => {
        const start = performance.now();
        factorial(20);
        const duration = performance.now() - start;
        expect(duration).toBeLessThan(10); // Should be fast
    });

    test('should check primality efficiently', () => {
        const start = performance.now();
        isPrime(99991);
        const duration = performance.now() - start;
        expect(duration).toBeLessThan(10);
    });

    test('should calculate Fibonacci efficiently', () => {
        const start = performance.now();
        fibonacci(40);
        const duration = performance.now() - start;
        expect(duration).toBeLessThan(1000); // Might be slow with recursion
    });
});

// ============================================================================
// Benchmark Tests
// ============================================================================

describe('Math Benchmarks', () => {
    
    test('add benchmark', () => {
        const start = performance.now();
        for (let i = 0; i < 1000000; i++) {
            add(i, i);
        }
        const duration = performance.now() - start;
        console.log(`1M additions: ${duration}ms`);
    });

    test('sqrt benchmark', () => {
        const start = performance.now();
        for (let i = 0; i < 100000; i++) {
            sqrt(i);
        }
        const duration = performance.now() - start;
        console.log(`100K square roots: ${duration}ms`);
    });

    test('prime check benchmark', () => {
        const start = performance.now();
        for (let i = 0; i < 10000; i++) {
            isPrime(i);
        }
        const duration = performance.now() - start;
        console.log(`10K prime checks: ${duration}ms`);
    });
});