/**
 * General Error Handling Examples
 * Demonstrates various error handling patterns and best practices
 */

const {
    AppError,
    ValidationError,
    NotFoundError,
    DatabaseError,
    ConfigurationError
} = require('./custom-errors');

// ============================================================================
// Try/Catch Patterns
// ============================================================================

/**
 * Basic try/catch example
 */
function divide(a, b) {
    try {
        if (typeof a !== 'number' || typeof b !== 'number') {
            throw new TypeError('Arguments must be numbers');
        }
        
        if (b === 0) {
            throw new Error('Division by zero');
        }
        
        return a / b;
    } catch (error) {
        if (error instanceof TypeError) {
            return { error: 'Invalid input types' };
        }
        
        if (error.message === 'Division by zero') {
            return Infinity;
        }
        
        throw error; // Re-throw unexpected errors
    }
}

/**
 * Try/catch/finally example
 */
function processFile(filename) {
    let fileHandle = null;
    
    try {
        console.log(`Opening file: ${filename}`);
        fileHandle = { filename, isOpen: true };
        
        if (!filename) {
            throw new Error('Filename is required');
        }
        
        if (filename === 'invalid') {
            throw new Error('Invalid file format');
        }
        
        return `Processed ${filename}`;
    } catch (error) {
        console.error('Error processing file:', error.message);
        throw new AppError(`Failed to process ${filename}`, 500, 'FILE_PROCESS_ERROR');
    } finally {
        if (fileHandle) {
            console.log('Closing file handle');
            fileHandle.isOpen = false;
        }
        console.log('Cleanup complete');
    }
}

/**
 * Multiple catch blocks pattern
 */
function parseUserInput(input) {
    try {
        const parsed = JSON.parse(input);
        
        if (!parsed.name) {
            throw new ValidationError('Missing required field', [
                { field: 'name', message: 'Name is required' }
            ]);
        }
        
        if (parsed.age && (parsed.age < 0 || parsed.age > 150)) {
            throw new ValidationError('Invalid age', [
                { field: 'age', message: 'Age must be between 0 and 150' }
            ]);
        }
        
        return parsed;
    } catch (error) {
        if (error instanceof SyntaxError) {
            return { error: 'Invalid JSON format' };
        }
        
        if (error instanceof ValidationError) {
            return { error: error.message, details: error.errors };
        }
        
        throw error;
    }
}

// ============================================================================
// Error handling in functions
// ============================================================================

/**
 * Guard clause pattern
 */
function getUserName(user) {
    if (!user) {
        throw new NotFoundError('User');
    }
    
    if (!user.name) {
        throw new ValidationError('User missing name', [
            { field: 'name', message: 'Name is required' }
        ]);
    }
    
    return user.name;
}

/**
 * Optional error handling pattern
 */
function findUser(id, options = {}) {
    const { throwIfNotFound = true, defaultValue = null } = options;
    
    // Simulate database lookup
    const users = {
        1: { id: 1, name: 'Alice' },
        2: { id: 2, name: 'Bob' }
    };
    
    const user = users[id];
    
    if (!user && throwIfNotFound) {
        throw new NotFoundError('User', id);
    }
    
    return user || defaultValue;
}

/**
 * Result object pattern (no exceptions)
 */
function divideSafely(a, b) {
    if (typeof a !== 'number' || typeof b !== 'number') {
        return {
            success: false,
            error: 'Invalid input types'
        };
    }
    
    if (b === 0) {
        return {
            success: false,
            error: 'Division by zero'
        };
    }
    
    return {
        success: true,
        result: a / b
    };
}

// ============================================================================
// Error wrapping and chaining
// ============================================================================

/**
 * Wraps lower-level errors with context
 */
async function fetchUserData(userId) {
    try {
        const response = await fetch(`/api/users/${userId}`);
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        return await response.json();
    } catch (error) {
        // Wrap with context
        throw new AppError(
            `Failed to fetch user ${userId}: ${error.message}`,
            500,
            'USER_FETCH_ERROR'
        );
    }
}

/**
 * Error transformation
 */
function handleDatabaseError(error) {
    if (error.code === 'ER_DUP_ENTRY') {
        throw new AppError('Duplicate entry', 409, 'DUPLICATE_ENTRY');
    }
    
    if (error.code === 'ECONNREFUSED') {
        throw new AppError('Database connection refused', 503, 'DB_UNAVAILABLE');
    }
    
    throw new DatabaseError('Unexpected database error', error);
}

// ============================================================================
// Error handling in constructors
// ============================================================================

class UserValidator {
    constructor(userData) {
        if (!userData) {
            throw new ValidationError('User data is required');
        }
        
        this.validateEmail(userData.email);
        this.validateAge(userData.age);
        this.validateName(userData.name);
        
        this.userData = userData;
    }
    
    validateEmail(email) {
        if (!email) {
            throw new ValidationError('Email is required', [
                { field: 'email', message: 'Required' }
            ]);
        }
        
        if (!email.includes('@')) {
            throw new ValidationError('Invalid email format', [
                { field: 'email', message: 'Must contain @' }
            ]);
        }
    }
    
    validateAge(age) {
        if (age !== undefined && (age < 0 || age > 150)) {
            throw new ValidationError('Invalid age', [
                { field: 'age', message: 'Must be between 0 and 150' }
            ]);
        }
    }
    
    validateName(name) {
        if (!name) {
            throw new ValidationError('Name is required', [
                { field: 'name', message: 'Required' }
            ]);
        }
        
        if (name.length < 2) {
            throw new ValidationError('Name too short', [
                { field: 'name', message: 'Must be at least 2 characters' }
            ]);
        }
    }
}

// ============================================================================
// Error handling in async initialization
// ============================================================================

class DatabaseConnection {
    constructor(config) {
        this.config = config;
        this.isConnected = false;
    }
    
    async connect() {
        try {
            await this.establishConnection();
            this.isConnected = true;
            console.log('Database connected');
        } catch (error) {
            throw new DatabaseError('Failed to connect to database', error);
        }
    }
    
    async establishConnection() {
        return new Promise((resolve, reject) => {
            setTimeout(() => {
                if (!this.config.url) {
                    reject(new Error('Missing database URL'));
                    return;
                }
                resolve();
            }, 100);
        });
    }
    
    async query(sql) {
        if (!this.isConnected) {
            throw new AppError('Not connected to database', 500, 'DB_NOT_CONNECTED');
        }
        
        try {
            return await this.executeQuery(sql);
        } catch (error) {
            throw new DatabaseError('Query failed', error);
        }
    }
    
    async executeQuery(sql) {
        // Simulate query execution
        return [{ id: 1, result: 'data' }];
    }
}

// ============================================================================
// Error handling with middleware pattern
// ============================================================================

class ErrorMiddleware {
    constructor() {
        this.handlers = [];
    }
    
    use(handler) {
        this.handlers.push(handler);
    }
    
    async handle(error, context) {
        let currentError = error;
        
        for (const handler of this.handlers) {
            try {
                const result = await handler(currentError, context);
                
                if (result === null) {
                    return null; // Error handled, stop propagation
                }
                
                if (result instanceof Error) {
                    currentError = result; // Transform error
                }
            } catch (handlerError) {
                console.error('Error in middleware:', handlerError);
            }
        }
        
        return currentError; // Unhandled error
    }
}

// Example middleware handlers
const errorLogger = async (error, context) => {
    console.error(`[${new Date().toISOString()}] Error:`, {
        message: error.message,
        code: error.code,
        stack: error.stack,
        context
    });
    return error;
};

const errorTransformer = async (error, context) => {
    if (error.code === 'ER_DUP_ENTRY') {
        return new AppError('Duplicate entry', 409, 'DUPLICATE_ENTRY');
    }
    return error;
};

const notFoundHandler = async (error, context) => {
    if (error instanceof NotFoundError) {
        console.log('Handling not found error, returning null');
        return null; // Handled
    }
    return error;
};

// ============================================================================
// Error aggregation and reporting
// ============================================================================

class ErrorReporter {
    constructor() {
        this.errors = [];
        this.handlers = [];
    }
    
    report(error, context = {}) {
        const errorReport = {
            timestamp: new Date().toISOString(),
            message: error.message,
            name: error.name,
            code: error.code,
            stack: error.stack,
            context
        };
        
        this.errors.push(errorReport);
        this.notifyHandlers(errorReport);
        
        return errorReport;
    }
    
    notifyHandlers(errorReport) {
        this.handlers.forEach(handler => handler(errorReport));
    }
    
    onReport(handler) {
        this.handlers.push(handler);
    }
    
    getReports() {
        return [...this.errors];
    }
    
    clear() {
        this.errors = [];
    }
}

// ============================================================================
// Export everything
// ============================================================================

module.exports = {
    // Basic patterns
    divide,
    processFile,
    parseUserInput,
    
    // Guard patterns
    getUserName,
    findUser,
    divideSafely,
    
    // Error wrapping
    fetchUserData,
    handleDatabaseError,
    
    // Classes with error handling
    UserValidator,
    DatabaseConnection,
    
    // Middleware pattern
    ErrorMiddleware,
    errorLogger,
    errorTransformer,
    notFoundHandler,
    
    // Error reporting
    ErrorReporter
};