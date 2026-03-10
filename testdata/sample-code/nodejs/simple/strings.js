/**
 * Simple String Operations
 * Basic string manipulation functions for testing
 */

// Case conversion
function toUpperCase(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.toUpperCase();
}

function toLowerCase(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.toLowerCase();
}

function capitalize(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    if (str.length === 0) return str;
    return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
}

function capitalizeWords(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.split(' ').map(word => capitalize(word)).join(' ');
}

// Trimming
function trim(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.trim();
}

function trimStart(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.trimStart();
}

function trimEnd(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.trimEnd();
}

// Searching
function contains(str, substring) {
    if (typeof str !== 'string' || typeof substring !== 'string') {
        throw new Error('Inputs must be strings');
    }
    return str.includes(substring);
}

function startsWith(str, prefix) {
    if (typeof str !== 'string' || typeof prefix !== 'string') {
        throw new Error('Inputs must be strings');
    }
    return str.startsWith(prefix);
}

function endsWith(str, suffix) {
    if (typeof str !== 'string' || typeof suffix !== 'string') {
        throw new Error('Inputs must be strings');
    }
    return str.endsWith(suffix);
}

function indexOf(str, substring) {
    if (typeof str !== 'string' || typeof substring !== 'string') {
        throw new Error('Inputs must be strings');
    }
    return str.indexOf(substring);
}

function lastIndexOf(str, substring) {
    if (typeof str !== 'string' || typeof substring !== 'string') {
        throw new Error('Inputs must be strings');
    }
    return str.lastIndexOf(substring);
}

function countOccurrences(str, substring) {
    if (typeof str !== 'string' || typeof substring !== 'string') {
        throw new Error('Inputs must be strings');
    }
    if (substring === '') return str.length + 1;
    return str.split(substring).length - 1;
}

// Substring extraction
function substring(str, start, end) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.substring(start, end);
}

function slice(str, start, end) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.slice(start, end);
}

function left(str, n) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.slice(0, n);
}

function right(str, n) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.slice(-n);
}

// Replacement
function replace(str, search, replacement) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.replace(search, replacement);
}

function replaceAll(str, search, replacement) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.replaceAll(search, replacement);
}

// Splitting and joining
function split(str, separator) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.split(separator);
}

function join(arr, separator) {
    if (!Array.isArray(arr)) {
        throw new Error('Input must be an array');
    }
    return arr.join(separator);
}

// Validation
function isEmpty(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.length === 0;
}

function isBlank(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return str.trim().length === 0;
}

function isNumeric(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return /^\d+$/.test(str);
}

function isAlpha(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return /^[a-zA-Z]+$/.test(str);
}

function isAlphanumeric(str) {
    if (typeof str !== 'string') {
        throw new Error('Input must be a string');
    }
    return /^[a-zA-Z0-9]+$/.test(str);
}

// Export all functions
module.exports = {
    toUpperCase,
    toLowerCase,
    capitalize,
    capitalizeWords,
    trim,
    trimStart,
    trimEnd,
    contains,
    startsWith,
    endsWith,
    indexOf,
    lastIndexOf,
    countOccurrences,
    substring,
    slice,
    left,
    right,
    replace,
    replaceAll,
    split,
    join,
    isEmpty,
    isBlank,
    isNumeric,
    isAlpha,
    isAlphanumeric
};