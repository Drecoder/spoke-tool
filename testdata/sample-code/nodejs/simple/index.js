/**
 * Node.js Simple Examples - Main Entry Point
 * 
 * This file exports all the simple examples for easy importing.
 * It serves as the main entry point for the package.
 */

// Import all modules
const math = require('./math');
const strings = require('./strings');
const objects = require('./objects');
const async = require('./async');

// Re-export everything
module.exports = {
    // Math functions
    ...math,
    
    // String functions
    ...strings,
    
    // Object/Class examples
    ...objects,
    
    // Async functions
    ...async,
    
    // Convenience groupings
    math,
    strings,
    objects,
    async
};

// Export a version constant
module.exports.version = '1.0.0';

// Export a simple greeting function
module.exports.greet = (name) => {
    return `Hello, ${name || 'World'}!`;
};

// Export a utility to run all examples (for testing)
module.exports.runExamples = () => {
    console.log('=== Math Examples ===');
    console.log('math.add(5, 3):', math.add(5, 3));
    console.log('math.subtract(10, 4):', math.subtract(10, 4));
    console.log('math.multiply(6, 7):', math.multiply(6, 7));
    console.log('math.divide(10, 2):', math.divide(10, 2));
    console.log('math.factorial(5):', math.factorial(5));
    console.log('math.fibonacci(10):', math.fibonacci(10));
    
    console.log('\n=== String Examples ===');
    console.log('strings.capitalize("hello"):', strings.capitalize('hello'));
    console.log('strings.toUpperCase("world"):', strings.toUpperCase('world'));
    console.log('strings.reverse("hello"):', strings.reverse('hello'));
    console.log('strings.countOccurrences("hello world hello", "hello"):', 
        strings.countOccurrences('hello world hello', 'hello'));
    
    console.log('\n=== Object Examples ===');
    const rect = new objects.Rectangle(5, 3);
    console.log('Rectangle(5,3) area:', rect.area());
    console.log('Rectangle(5,3) perimeter:', rect.perimeter());
    
    const dog = new objects.Dog('Rex', 'German Shepherd');
    console.log('Dog speak:', dog.speak());
    
    console.log('\n=== Async Examples ===');
    console.log('See async.test.js for async examples');
};