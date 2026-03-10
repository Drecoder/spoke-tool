/**
 * Async Error Handling Examples
 * Demonstrates error handling in promises, async/await, and callbacks
 */

const { 
    AppError, 
    NotFoundError, 
    ValidationError,
    DatabaseError 
} = require('./custom-errors');

// ============================================================================
// Promise-based error handling
// ============================================================================

/**
 * Simulates fetching a user from a database
 */
function fetchUser(id) {
    return new Promise((resolve, reject) => {
        setTimeout(() => {
            if (!id || id <= 0) {
                reject(new ValidationError('Invalid user ID', [
                    { field: 'id', message: 'ID must be positive' }
                ]));
                return;
            }
            
            if (id === 404) {
                reject(new NotFoundError('User', id));
                return;
            }
            
            if (id === 500) {
                reject(new DatabaseError('Database connection failed'));
                return;
            }
            
            resolve({
                id,
                name: `User ${id}`,
                email: `user${id}@example.com`
            });
        }, 100);
    });
}

/**
 * Fetches user and handles errors with .catch()
 */
function getUserWithCatch(id) {
    return fetchUser(id)
        .then(user => {
            console.log('User found:', user);
            return user;
        })
        .catch(error => {
            console.error('Error fetching user:', error.message);
            
            if (error instanceof NotFoundError) {
                return { id, name: 'Default User', email: 'default@example.com' };
            }
            
            throw error;
        });
}

/**
 * Fetches user and handles specific error types
 */
function getUserWithSpecificHandling(id) {
    return fetchUser(id)
        .then(user => {
            return { success: true, data: user };
        })
        .catch(error => {
            if (error instanceof ValidationError) {
                return { success: false, errors: error.errors };
            }
            
            if (error instanceof NotFoundError) {
                return { success: false, message: 'User not found, creating default' };
            }
            
            return { success: false, message: 'Unexpected error', error: error.message };
        });
}

// ============================================================================
// Async/await error handling
// ============================================================================

/**
 * Fetches user using async/await with try/catch
 */
async function getUserAsync(id) {
    try {
        const user = await fetchUser(id);
        return { success: true, data: user };
    } catch (error) {
        if (error instanceof AppError) {
            return {
                success: false,
                code: error.code,
                message: error.message,
                statusCode: error.statusCode
            };
        }
        
        return {
            success: false,
            message: 'Unexpected error',
            error: error.message
        };
    }
}

/**
 * Fetches user and posts with error handling
 */
async function getUserWithPosts(userId) {
    try {
        const user = await fetchUser(userId);
        
        try {
            const posts = await fetchPosts(user.id);
            return { ...user, posts };
        } catch (postError) {
            console.warn('Failed to fetch posts:', postError.message);
            return { ...user, posts: [] };
        }
    } catch (error) {
        console.error('Failed to fetch user:', error.message);
        throw new AppError('Failed to load user data', 500, 'USER_DATA_ERROR');
    }
}

/**
 * Simulates fetching posts for a user
 */
function fetchPosts(userId) {
    return new Promise((resolve, reject) => {
        setTimeout(() => {
            if (userId === 999) {
                reject(new Error('Posts service unavailable'));
                return;
            }
            
            resolve([
                { id: 1, title: 'Post 1' },
                { id: 2, title: 'Post 2' }
            ]);
        }, 100);
    });
}

// ============================================================================
// Multiple async operations with error handling
// ============================================================================

/**
 * Fetches multiple users with Promise.all
 */
async function fetchMultipleUsers(userIds) {
    try {
        const promises = userIds.map(id => fetchUser(id));
        const users = await Promise.all(promises);
        return { success: true, data: users };
    } catch (error) {
        return {
            success: false,
            message: 'Failed to fetch some users',
            error: error.message
        };
    }
}

/**
 * Fetches multiple users with Promise.allSettled
 */
async function fetchMultipleUsersSettled(userIds) {
    const promises = userIds.map(id => fetchUser(id));
    const results = await Promise.allSettled(promises);
    
    return {
        fulfilled: results
            .filter(r => r.status === 'fulfilled')
            .map(r => r.value),
        rejected: results
            .filter(r => r.status === 'rejected')
            .map(r => ({
                reason: r.reason.message,
                code: r.reason.code
            }))
    };
}

/**
 * Fetches with timeout
 */
async function fetchWithTimeout(id, timeoutMs = 1000) {
    const timeoutPromise = new Promise((_, reject) => {
        setTimeout(() => reject(new Error('Request timed out')), timeoutMs);
    });
    
    try {
        const user = await Promise.race([fetchUser(id), timeoutPromise]);
        return { success: true, data: user };
    } catch (error) {
        if (error.message === 'Request timed out') {
            return {
                success: false,
                message: 'Request timed out',
                timeout: timeoutMs
            };
        }
        throw error;
    }
}

// ============================================================================
// Retry logic with error handling
// ============================================================================

/**
 * Retries an async operation with exponential backoff
 */
async function retry(operation, maxRetries = 3, delay = 100) {
    let lastError;
    
    for (let attempt = 1; attempt <= maxRetries; attempt++) {
        try {
            return await operation();
        } catch (error) {
            lastError = error;
            
            // Don't retry validation errors
            if (error instanceof ValidationError) {
                throw error;
            }
            
            // Don't retry not found errors
            if (error instanceof NotFoundError) {
                throw error;
            }
            
            if (attempt < maxRetries) {
                const waitTime = delay * Math.pow(2, attempt - 1);
                console.log(`Attempt ${attempt} failed, retrying in ${waitTime}ms...`);
                await new Promise(resolve => setTimeout(resolve, waitTime));
            }
        }
    }
    
    throw new Error(`Operation failed after ${maxRetries} attempts: ${lastError.message}`);
}

// ============================================================================
// Callback-style error handling
// ============================================================================

/**
 * Callback-style function with error-first pattern
 */
function readConfig(path, callback) {
    setTimeout(() => {
        if (!path) {
            callback(new Error('Path is required'), null);
            return;
        }
        
        if (path === 'invalid') {
            callback(new Error('Invalid config path'), null);
            return;
        }
        
        callback(null, { port: 3000, host: 'localhost' });
    }, 100);
}

/**
 * Promisified version of readConfig
 */
function readConfigPromise(path) {
    return new Promise((resolve, reject) => {
        readConfig(path, (err, config) => {
            if (err) reject(err);
            else resolve(config);
        });
    });
}

// ============================================================================
// Error handling in streams/events
// ============================================================================

class EventEmitter {
    constructor() {
        this.events = {};
        this.errorHandlers = [];
    }
    
    on(event, handler) {
        if (!this.events[event]) {
            this.events[event] = [];
        }
        this.events[event].push(handler);
    }
    
    onError(handler) {
        this.errorHandlers.push(handler);
    }
    
    emit(event, data) {
        if (!this.events[event]) return;
        
        try {
            this.events[event].forEach(handler => {
                try {
                    handler(data);
                } catch (handlerError) {
                    this.handleError(new Error(`Handler error: ${handlerError.message}`));
                }
            });
        } catch (error) {
            this.handleError(error);
        }
    }
    
    handleError(error) {
        this.errorHandlers.forEach(handler => handler(error));
        
        if (this.errorHandlers.length === 0) {
            console.error('Unhandled error in EventEmitter:', error);
        }
    }
}

// ============================================================================
// Error aggregation
// ============================================================================

class ErrorAggregator {
    constructor() {
        this.errors = [];
        this.warnings = [];
    }
    
    addError(error) {
        this.errors.push({
            error: error.message,
            code: error.code,
            timestamp: new Date().toISOString()
        });
    }
    
    addWarning(message) {
        this.warnings.push({
            message,
            timestamp: new Date().toISOString()
        });
    }
    
    hasErrors() {
        return this.errors.length > 0;
    }
    
    hasWarnings() {
        return this.warnings.length > 0;
    }
    
    clear() {
        this.errors = [];
        this.warnings = [];
    }
    
    getReport() {
        return {
            errorCount: this.errors.length,
            warningCount: this.warnings.length,
            errors: this.errors,
            warnings: this.warnings
        };
    }
}

// ============================================================================
// Export everything
// ============================================================================

module.exports = {
    // Core functions
    fetchUser,
    fetchPosts,
    
    // Promise-based handling
    getUserWithCatch,
    getUserWithSpecificHandling,
    
    // Async/await handling
    getUserAsync,
    getUserWithPosts,
    
    // Multiple operations
    fetchMultipleUsers,
    fetchMultipleUsersSettled,
    fetchWithTimeout,
    
    // Retry logic
    retry,
    
    // Callback patterns
    readConfig,
    readConfigPromise,
    
    // Event handling
    EventEmitter,
    
    // Error aggregation
    ErrorAggregator
};