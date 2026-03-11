/**
 * Simple Math Operations
 * Basic arithmetic functions for testing
 */

// Basic arithmetic
function add(a, b) {
    return a + b;
}

function subtract(a, b) {
    return a - b;
}

function multiply(a, b) {
    return a * b;
}

function divide(a, b) {
    if (b === 0) {
        throw new Error('Division by zero');
    }
    return a / b;
}

// Power and roots
function square(x) {
    return x * x;
}

function cube(x) {
    return x * x * x;
}

function power(base, exponent) {
    return Math.pow(base, exponent);
}

function sqrt(x) {
    if (x < 0) {
        throw new Error('Cannot calculate square root of negative number');
    }
    return Math.sqrt(x);
}

// Number theory
function isEven(n) {
    return n % 2 === 0;
}

function isOdd(n) {
    return n % 2 !== 0;
}

function factorial(n) {
    if (n < 0) {
        throw new Error('Factorial not defined for negative numbers');
    }
    if (n <= 1) return 1;
    return n * factorial(n - 1);
}

function fibonacci(n) {
    if (n < 0) {
        throw new Error('Fibonacci not defined for negative numbers');
    }
    if (n <= 1) return n;
    return fibonacci(n - 1) + fibonacci(n - 2);
}

// Min/max
function min(a, b) {
    return a < b ? a : b;
}

function max(a, b) {
    return a > b ? a : b;
}

function minOfArray(arr) {
    if (!Array.isArray(arr) || arr.length === 0) {
        throw new Error('Invalid array');
    }
    return Math.min(...arr);
}

function maxOfArray(arr) {
    if (!Array.isArray(arr) || arr.length === 0) {
        throw new Error('Invalid array');
    }
    return Math.max(...arr);
}

// Sum and average
function sum(arr) {
    if (!Array.isArray(arr)) {
        throw new Error('Invalid array');
    }
    return arr.reduce((acc, val) => acc + val, 0);
}

function average(arr) {
    if (!Array.isArray(arr) || arr.length === 0) {
        throw new Error('Invalid array');
    }
    return sum(arr) / arr.length;
}

// Geometry
function areaOfCircle(radius) {
    if (radius < 0) {
        throw new Error('Radius cannot be negative');
    }
    return Math.PI * radius * radius;
}

function circumference(radius) {
    if (radius < 0) {
        throw new Error('Radius cannot be negative');
    }
    return 2 * Math.PI * radius;
}

// Export all functions
module.exports = {
    add,
    subtract,
    multiply,
    divide,
    square,
    cube,
    power,
    sqrt,
    isEven,
    isOdd,
    factorial,
    fibonacci,
    min,
    max,
    minOfArray,
    maxOfArray,
    sum,
    average,
    areaOfCircle,
    circumference
};