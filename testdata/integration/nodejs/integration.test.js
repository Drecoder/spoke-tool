/**
 * Node.js Integration Tests
 * 
 * Tests the interaction between multiple components of the system
 * in a real Node.js environment.
 */

const fs = require('fs').promises;
const path = require('path');
const os = require('os');
const { exec } = require('child_process');
const util = require('util');
const execPromise = util.promisify(exec);

// Mock implementations of your components
// In a real implementation, these would import your actual modules
const { analyzeProject } = require('../src/analyzer');
const { generateTests } = require('../src/generator');
const { extractDocs } = require('../src/extractor');
const { formatReadme } = require('../src/formatter');
const { updateReadme } = require('../src/updater');

// ============================================================================
// Test Environment Setup
// ============================================================================

let testDir;
let tempDir;

beforeAll(async () => {
    // Create temporary test directory
    tempDir = await fs.mkdtemp(path.join(os.tmpdir(), 'spoke-tool-node-'));
    testDir = path.join(tempDir, 'test-project');
    await fs.mkdir(testDir, { recursive: true });
});

afterAll(async () => {
    // Clean up temporary directory
    await fs.rm(tempDir, { recursive: true, force: true });
});

beforeEach(async () => {
    // Clear test directory before each test
    const files = await fs.readdir(testDir);
    for (const file of files) {
        await fs.rm(path.join(testDir, file), { recursive: true, force: true });
    }
});

// ============================================================================
// Helper Functions
// ============================================================================

async function createTestFile(filename, content) {
    const filePath = path.join(testDir, filename);
    await fs.writeFile(filePath, content);
    return filePath;
}

async function fileExists(filename) {
    try {
        await fs.access(path.join(testDir, filename));
        return true;
    } catch {
        return false;
    }
}

async function readTestFile(filename) {
    return await fs.readFile(path.join(testDir, filename), 'utf8');
}

// ============================================================================
// Test Spoke Integration Tests
// ============================================================================

describe('Test Spoke Integration', () => {
    
    test('should analyze JavaScript code and find functions', async () => {
        // Create test JavaScript file
        const jsCode = `
// math.js
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

const processData = (data) => {
    return data.map(x => x * 2);
};

class Calculator {
    add(a, b) {
        return a + b;
    }
    
    subtract(a, b) {
        return a - b;
    }
}

module.exports = {
    add,
    subtract,
    multiply,
    divide,
    processData,
    Calculator
};
`;
        await createTestFile('math.js', jsCode);

        // Analyze project
        const analysis = await analyzeProject(testDir, { language: 'nodejs' });

        // Verify analysis results
        expect(analysis.functions).toBeDefined();
        expect(analysis.functions.length).toBeGreaterThan(0);
        
        // Check for specific functions
        const functionNames = analysis.functions.map(f => f.name);
        expect(functionNames).toContain('add');
        expect(functionNames).toContain('subtract');
        expect(functionNames).toContain('multiply');
        expect(functionNames).toContain('divide');
        expect(functionNames).toContain('processData');
        
        // Check for class methods
        expect(functionNames).toContain('add'); // Calculator.add
        expect(functionNames).toContain('subtract'); // Calculator.subtract
    });

    test('should analyze TypeScript code with type annotations', async () => {
        // Create test TypeScript file
        const tsCode = `
// user.ts
interface User {
    id: number;
    name: string;
    email?: string;
    readonly createdAt: Date;
}

type UserID = string | number;

class UserService {
    private users: Map<UserID, User> = new Map();
    
    createUser(id: UserID, name: string, email?: string): User {
        const user: User = {
            id: typeof id === 'string' ? parseInt(id) : id,
            name,
            email,
            createdAt: new Date()
        };
        this.users.set(id, user);
        return user;
    }
    
    getUser(id: UserID): User | undefined {
        return this.users.get(id);
    }
    
    async fetchUser(id: number): Promise<User> {
        const response = await fetch(\`/api/users/\${id}\`);
        return await response.json();
    }
}

export { User, UserID, UserService };
`;
        await createTestFile('user.ts', tsCode);

        // Analyze project
        const analysis = await analyzeProject(testDir, { language: 'typescript' });

        // Verify analysis found TypeScript constructs
        expect(analysis.functions).toBeDefined();
        expect(analysis.classes).toBeDefined();
        expect(analysis.interfaces).toBeDefined();
        
        // Check for UserService class
        const userService = analysis.classes.find(c => c.name === 'UserService');
        expect(userService).toBeDefined();
        expect(userService.methods).toContain('createUser');
        expect(userService.methods).toContain('getUser');
        expect(userService.methods).toContain('fetchUser');
        
        // Check for User interface
        const userInterface = analysis.interfaces.find(i => i.name === 'User');
        expect(userInterface).toBeDefined();
        expect(userInterface.properties).toContain('id');
        expect(userInterface.properties).toContain('name');
    });

    test('should find untested functions', async () => {
        // Create source file with functions
        const sourceCode = `
// calculator.js
export function add(a, b) {
    return a + b;
}

export function subtract(a, b) {
    return a - b;
}

export function multiply(a, b) {
    return a * b;
}
`;
        await createTestFile('calculator.js', sourceCode);

        // Create test file (only tests add and subtract)
        const testCode = `
// calculator.test.js
import { add, subtract } from './calculator';

describe('Calculator', () => {
    test('add', () => {
        expect(add(2, 3)).toBe(5);
    });
    
    test('subtract', () => {
        expect(subtract(10, 4)).toBe(6);
    });
});
`;
        await createTestFile('calculator.test.js', testCode);

        // Analyze project
        const analysis = await analyzeProject(testDir, { language: 'nodejs' });

        // Find untested functions
        const untested = analysis.functions.filter(f => !f.hasTest);
        
        // multiply should be untested
        expect(untested.length).toBe(1);
        expect(untested[0].name).toBe('multiply');
    });

    test('should generate test file for untested functions', async () => {
        // Create source file
        const sourceCode = `
// math.js
export function add(a, b) {
    return a + b;
}

export function subtract(a, b) {
    return a - b;
}

export function multiply(a, b) {
    return a * b;
}

export function divide(a, b) {
    if (b === 0) throw new Error('Division by zero');
    return a / b;
}
`;
        await createTestFile('math.js', sourceCode);

        // Analyze project
        const analysis = await analyzeProject(testDir, { language: 'nodejs' });

        // Find untested functions
        const untested = analysis.functions.filter(f => !f.hasTest);

        // Generate tests
        const generatedTests = await generateTests(analysis, untested);

        // Verify test generation
        expect(generatedTests.length).toBeGreaterThan(0);
        
        // Write test file
        const testFilePath = path.join(testDir, 'math.test.js');
        await fs.writeFile(testFilePath, generatedTests[0].code);

        // Verify file was created
        const exists = await fileExists('math.test.js');
        expect(exists).toBe(true);
    });
});

// ============================================================================
// Readme Spoke Integration Tests
// ============================================================================

describe('Readme Spoke Integration', () => {
    
    test('should extract documentation from code and tests', async () => {
        // Create source file with JSDoc comments
        const sourceCode = `
/**
 * Calculator class for performing arithmetic operations
 * @class
 */
class Calculator {
    /**
     * Adds two numbers together
     * @param {number} a - First number
     * @param {number} b - Second number
     * @returns {number} The sum of a and b
     */
    add(a, b) {
        return a + b;
    }
    
    /**
     * Subtracts two numbers
     * @param {number} a - First number
     * @param {number} b - Second number
     * @returns {number} The difference a - b
     */
    subtract(a, b) {
        return a - b;
    }
    
    /**
     * Multiplies two numbers
     * @param {number} a - First number
     * @param {number} b - Second number
     * @returns {number} The product a * b
     */
    multiply(a, b) {
        return a * b;
    }
}

module.exports = Calculator;
`;
        await createTestFile('calculator.js', sourceCode);

        // Create test file with examples
        const testCode = `
const Calculator = require('./calculator');

describe('Calculator', () => {
    test('add should return sum', () => {
        const calc = new Calculator();
        expect(calc.add(2, 3)).toBe(5);
    });
    
    test('subtract should return difference', () => {
        const calc = new Calculator();
        expect(calc.subtract(10, 4)).toBe(6);
    });
});
`;
        await createTestFile('calculator.test.js', testCode);

        // Extract documentation
        const docs = await extractDocs(testDir);

        // Verify extraction
        expect(docs.functions).toBeDefined();
        expect(docs.functions.length).toBe(3);
        
        // Check JSDoc extraction
        const addFunc = docs.functions.find(f => f.name === 'add');
        expect(addFunc.description).toContain('Adds two numbers');
        expect(addFunc.params).toContainEqual(expect.objectContaining({ name: 'a', type: 'number' }));
        expect(addFunc.params).toContainEqual(expect.objectContaining({ name: 'b', type: 'number' }));
        expect(addFunc.returns.type).toBe('number');
        
        // Check example extraction from tests
        expect(addFunc.examples).toBeDefined();
        expect(addFunc.examples.length).toBeGreaterThan(0);
    });

    test('should generate README from extracted content', async () => {
        // Create source files
        await createTestFile('index.js', `
/**
 * Main entry point
 * @module mylib
 */

/**
 * Returns the current version
 * @returns {string} Version string
 */
function version() {
    return '1.0.0';
}

module.exports = { version };
`);

        // Extract documentation
        const docs = await extractDocs(testDir);

        // Generate README
        const readme = await formatReadme({
            title: 'My Library',
            description: 'A sample library for testing',
            docs: docs,
            sections: ['installation', 'usage', 'api', 'examples']
        });

        // Verify README content
        expect(readme).toContain('# My Library');
        expect(readme).toContain('A sample library for testing');
        expect(readme).toContain('## Installation');
        expect(readme).toContain('## Usage');
        expect(readme).toContain('## API Reference');
        
        // Write README
        await fs.writeFile(path.join(testDir, 'README.md'), readme);

        // Verify file was created
        const exists = await fileExists('README.md');
        expect(exists).toBe(true);
    });

    test('should update existing README while preserving manual content', async () => {
        // Create existing README with manual content
        const existingReadme = `# My Project

This is a manually written description.

## Installation

\`\`\`bash
npm install my-project
\`\`\`

<!-- GENERATED SECTION - DO NOT EDIT MANUALLY -->
## API Reference

This section will be auto-generated.

<!-- END GENERATED SECTION -->

## Manual Section

This section was written by hand and should be preserved.
`;
        await createTestFile('README.md', existingReadme);

        // Generate new API content
        const newApiContent = `## API Reference

### \`add(a, b)\`

Adds two numbers together.

**Parameters:**
- \`a\` (number): First number
- \`b\` (number): Second number

**Returns:** number - The sum
`;

        // Update README
        const updated = await updateReadme(testDir, {
            generatedContent: newApiContent,
            marker: '<!-- GENERATED SECTION -->'
        });

        // Verify manual content preserved
        expect(updated).toContain('This is a manually written description');
        expect(updated).toContain('## Manual Section');
        expect(updated).toContain('This section was written by hand');
        
        // Verify generated content updated
        expect(updated).toContain('### `add(a, b)`');
        expect(updated).toContain('Adds two numbers together');
        
        // Verify installation section preserved
        expect(updated).toContain('npm install my-project');
    });
});

// ============================================================================
// End-to-End Integration Tests
// ============================================================================

describe('End-to-End Integration', () => {
    
    test('complete workflow: analyze → generate tests → extract → generate readme', async () => {
        // Create a complete project
        const projectCode = `
// user.js
/**
 * Represents a user in the system
 */
class User {
    /**
     * Create a new user
     * @param {number} id - User ID
     * @param {string} name - User name
     * @param {number} age - User age
     */
    constructor(id, name, age) {
        this.id = id;
        this.name = name;
        this.age = age;
    }
    
    /**
     * Check if user is an adult
     * @returns {boolean} True if age >= 18
     */
    isAdult() {
        return this.age >= 18;
    }
    
    /**
     * Get user greeting
     * @returns {string} Greeting message
     */
    greet() {
        return \`Hello, my name is \${this.name}\`;
    }
}

module.exports = User;
`;
        await createTestFile('user.js', projectCode);

        // Step 1: Analyze project
        const analysis = await analyzeProject(testDir, { language: 'nodejs' });
        
        expect(analysis.functions).toBeDefined();
        expect(analysis.classes).toBeDefined();
        
        const userClass = analysis.classes.find(c => c.name === 'User');
        expect(userClass).toBeDefined();
        expect(userClass.methods).toContain('isAdult');
        expect(userClass.methods).toContain('greet');

        // Step 2: Find untested functions and generate tests
        const untested = analysis.functions.filter(f => !f.hasTest);
        
        if (untested.length > 0) {
            const generatedTests = await generateTests(analysis, untested);
            
            if (generatedTests.length > 0) {
                await fs.writeFile(
                    path.join(testDir, 'user.test.js'),
                    generatedTests[0].code
                );
            }
        }

        // Step 3: Extract documentation
        const docs = await extractDocs(testDir);
        
        // Step 4: Generate README
        const readme = await formatReadme({
            title: 'User Library',
            description: 'A simple user management library',
            docs: docs,
            sections: ['installation', 'usage', 'api', 'examples']
        });

        await fs.writeFile(path.join(testDir, 'README.md'), readme);

        // Step 5: Verify all files exist
        const files = ['user.js', 'README.md'];
        for (const file of files) {
            const exists = await fileExists(file);
            expect(exists).toBe(true);
        }
        
        // Optionally check for test file if it was created
        const testExists = await fileExists('user.test.js');
        if (testExists) {
            const testContent = await readTestFile('user.test.js');
            expect(testContent).toBeDefined();
        }
    }, 30000); // 30 second timeout
});

// ============================================================================
// Performance Tests
// ============================================================================

describe('Performance Tests', () => {
    
    test('should analyze large project within time limit', async () => {
        // Generate many files
        for (let i = 0; i < 50; i++) {
            const content = `
// file${i}.js
export function func${i}() {
    return ${i};
}

export class Class${i} {
    method() {
        return ${i};
    }
}
`;
            await createTestFile(`file${i}.js`, content);
        }

        // Measure analysis time
        const start = Date.now();
        const analysis = await analyzeProject(testDir, { language: 'nodejs' });
        const duration = Date.now() - start;

        // Verify analysis completed
        expect(analysis.files).toBeDefined();
        expect(analysis.files.length).toBe(50);
        
        // Performance assertion (adjust threshold as needed)
        console.log(`Analyzed 50 files in ${duration}ms`);
        expect(duration).toBeLessThan(5000); // Should complete within 5 seconds
    });

    test('should handle concurrent analysis requests', async () => {
        // Create multiple projects
        const projects = [];
        for (let i = 0; i < 5; i++) {
            const projDir = path.join(testDir, `project-${i}`);
            await fs.mkdir(projDir, { recursive: true });
            
            const content = `
// index.js
export function func${i}() {
    return ${i};
}
`;
            await fs.writeFile(path.join(projDir, 'index.js'), content);
            projects.push(projDir);
        }

        // Analyze concurrently
        const start = Date.now();
        const results = await Promise.all(
            projects.map(dir => analyzeProject(dir, { language: 'nodejs' }))
        );
        const duration = Date.now() - start;

        // Verify all analyses completed
        expect(results.length).toBe(5);
        results.forEach(r => {
            expect(r.functions).toBeDefined();
        });

        console.log(`Concurrent analysis of 5 projects took ${duration}ms`);
    });
});

// ============================================================================
// Error Handling Tests
// ============================================================================

describe('Error Handling', () => {
    
    test('should handle invalid JavaScript syntax gracefully', async () => {
        const invalidCode = `
function broken( {
    return } @#$%
`;
        await createTestFile('broken.js', invalidCode);

        // Should not throw, but return analysis with error info
        const analysis = await analyzeProject(testDir, { language: 'nodejs' });
        
        expect(analysis.errors).toBeDefined();
        expect(analysis.errors.length).toBeGreaterThan(0);
    });

    test('should handle missing files gracefully', async () => {
        // Analyze empty directory
        const analysis = await analyzeProject(testDir, { language: 'nodejs' });
        
        expect(analysis.files).toEqual([]);
        expect(analysis.functions).toEqual([]);
        expect(analysis.warnings).toBeDefined();
    });

    test('should handle very large files without crashing', async () => {
        // Create a very large file
        const largeContent = '// Large file\n' + 'x = 42;\n'.repeat(100000);
        await createTestFile('large.js', largeContent);

        // Should not crash
        const analysis = await analyzeProject(testDir, { language: 'nodejs' });
        expect(analysis).toBeDefined();
    });
});