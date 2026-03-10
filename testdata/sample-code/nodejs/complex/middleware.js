/**
 * Middleware Module
 * 
 * Custom middleware for Express.js application including:
 * - Authentication & Authorization
 * - Request Logging
 * - Rate Limiting
 * - Error Handling
 * - Request Validation
 * - Security Headers
 * - Compression
 * - CORS
 */

const jwt = require('jsonwebtoken');
const rateLimit = require('express-rate-limit');
const compression = require('compression');
const helmet = require('helmet');
const cors = require('cors');
const mongoSanitize = require('express-mongo-sanitize');
const xss = require('xss-clean');
const hpp = require('hpp');
const useragent = require('express-useragent');
const requestIp = require('request-ip');
const { v4: uuidv4 } = require('uuid');
const { logger } = require('../utils/logger');
const { cache } = require('../utils/cache');
const { metrics } = require('../utils/metrics');
const { AppError } = require('../utils/errors');
const User = require('../models/User');

// ============================================================================
// Request ID Middleware
// ============================================================================

/**
 * Generate and attach unique request ID
 */
const requestId = (req, res, next) => {
    req.id = uuidv4();
    res.setHeader('X-Request-ID', req.id);
    next();
};

// ============================================================================
// Request Logging Middleware
// ============================================================================

/**
 * Detailed request logging with performance metrics
 */
const requestLogger = (req, res, next) => {
    const start = Date.now();
    
    // Log on response finish
    res.on('finish', () => {
        const duration = Date.now() - start;
        const logData = {
            requestId: req.id,
            method: req.method,
            url: req.originalUrl || req.url,
            status: res.statusCode,
            duration: `${duration}ms`,
            ip: req.clientIp,
            userAgent: req.headers['user-agent'],
            userId: req.user?.id || 'anonymous'
        };
        
        // Track metrics
        metrics.httpRequestDuration
            .labels(req.method, req.route?.path || req.path, res.statusCode)
            .observe(duration / 1000);
        
        metrics.httpRequestsTotal
            .labels(req.method, req.route?.path || req.path, res.statusCode)
            .inc();
        
        if (res.statusCode >= 500) {
            logger.error('Request failed', logData);
        } else if (res.statusCode >= 400) {
            logger.warn('Request warning', logData);
        } else {
            logger.info('Request completed', logData);
        }
    });
    
    next();
};

// ============================================================================
// Authentication Middleware
// ============================================================================

/**
 * Verify JWT token and attach user to request
 */
const authenticate = async (req, res, next) => {
    try {
        const authHeader = req.headers.authorization;
        
        if (!authHeader || !authHeader.startsWith('Bearer ')) {
            throw new AppError('No token provided', 401);
        }
        
        const token = authHeader.split(' ')[1];
        
        // Check if token is blacklisted
        const isBlacklisted = await cache.get(`blacklist:${token}`);
        if (isBlacklisted) {
            throw new AppError('Token has been revoked', 401);
        }
        
        // Verify token
        const decoded = jwt.verify(token, process.env.JWT_SECRET);
        
        // Check if user still exists
        const user = await User.findById(decoded.id).select('-password');
        if (!user) {
            throw new AppError('User no longer exists', 401);
        }
        
        // Check if user is active
        if (!user.isActive) {
            throw new AppError('Account is deactivated', 401);
        }
        
        // Check if password was changed after token was issued
        if (user.passwordChangedAt && decoded.iat < user.passwordChangedAt.getTime() / 1000) {
            throw new AppError('Password was changed, please login again', 401);
        }
        
        req.user = user;
        next();
    } catch (error) {
        if (error.name === 'JsonWebTokenError') {
            return next(new AppError('Invalid token', 401));
        }
        if (error.name === 'TokenExpiredError') {
            return next(new AppError('Token expired', 401));
        }
        next(error);
    }
};

/**
 * Optional authentication - doesn't error if no token
 */
const optionalAuth = async (req, res, next) => {
    try {
        const authHeader = req.headers.authorization;
        
        if (authHeader && authHeader.startsWith('Bearer ')) {
            const token = authHeader.split(' ')[1];
            const decoded = jwt.verify(token, process.env.JWT_SECRET);
            
            const user = await User.findById(decoded.id).select('-password');
            if (user && user.isActive) {
                req.user = user;
            }
        }
        next();
    } catch (error) {
        // Silently fail auth for optional routes
        next();
    }
};

// ============================================================================
// Authorization Middleware
// ============================================================================

/**
 * Restrict access to specific roles
 */
const authorize = (...roles) => {
    return (req, res, next) => {
        if (!req.user) {
            return next(new AppError('Authentication required', 401));
        }
        
        if (!roles.includes(req.user.role)) {
            return next(new AppError('Insufficient permissions', 403));
        }
        
        next();
    };
};

/**
 * Check if user owns the resource
 */
const checkOwnership = (resourceField = 'userId') => {
    return async (req, res, next) => {
        try {
            if (!req.user) {
                return next(new AppError('Authentication required', 401));
            }
            
            // Admin can access all
            if (req.user.role === 'admin') {
                return next();
            }
            
            const resourceId = req.params.id;
            const Model = req.model;
            
            if (!Model) {
                return next(new AppError('Model not defined for ownership check', 500));
            }
            
            const resource = await Model.findById(resourceId);
            if (!resource) {
                return next(new AppError('Resource not found', 404));
            }
            
            // Check if user owns the resource
            const ownerId = resource[resourceField]?.toString();
            if (ownerId !== req.user.id) {
                return next(new AppError('You do not own this resource', 403));
            }
            
            req.resource = resource;
            next();
        } catch (error) {
            next(error);
        }
    };
};

// ============================================================================
// Rate Limiting Middleware
// ============================================================================

/**
 * General rate limiter
 */
const limiter = rateLimit({
    windowMs: 15 * 60 * 1000, // 15 minutes
    max: 100, // Limit each IP to 100 requests per windowMs
    message: {
        success: false,
        message: 'Too many requests, please try again later.'
    },
    standardHeaders: true,
    legacyHeaders: false,
    keyGenerator: (req) => {
        return req.user?.id || req.clientIp;
    },
    skip: (req) => {
        // Skip rate limiting for health checks
        return req.path === '/health' || req.path === '/metrics';
    },
    handler: (req, res) => {
        logger.warn(`Rate limit exceeded for IP: ${req.clientIp}`);
        res.status(429).json({
            success: false,
            message: 'Too many requests, please try again later.'
        });
    }
});

/**
 * Strict rate limiter for auth endpoints
 */
const authLimiter = rateLimit({
    windowMs: 15 * 60 * 1000, // 15 minutes
    max: 5, // 5 attempts per 15 minutes
    skipSuccessfulRequests: true, // Don't count successful logins
    message: {
        success: false,
        message: 'Too many authentication attempts, please try again later.'
    }
});

/**
 * API key rate limiter
 */
const apiKeyLimiter = rateLimit({
    windowMs: 60 * 60 * 1000, // 1 hour
    max: 1000,
    keyGenerator: (req) => req.apiKey,
    message: {
        success: false,
        message: 'API key rate limit exceeded'
    }
});

// ============================================================================
// API Key Authentication
// ============================================================================

/**
 * Validate API key
 */
const validateApiKey = async (req, res, next) => {
    try {
        const apiKey = req.headers['x-api-key'];
        
        if (!apiKey) {
            return next(new AppError('API key required', 401));
        }
        
        // Check cache first
        let apiKeyData = await cache.get(`apikey:${apiKey}`);
        
        if (!apiKeyData) {
            // Verify from database (simplified)
            const validKeys = {
                'test-key-123': { client: 'Test Client', tier: 'basic' },
                'premium-key-456': { client: 'Premium Client', tier: 'premium' }
            };
            
            apiKeyData = validKeys[apiKey];
            
            if (!apiKeyData) {
                return next(new AppError('Invalid API key', 401));
            }
            
            // Cache for 1 hour
            await cache.set(`apikey:${apiKey}`, JSON.stringify(apiKeyData), 3600);
        } else {
            apiKeyData = JSON.parse(apiKeyData);
        }
        
        req.apiKey = apiKey;
        req.client = apiKeyData.client;
        req.tier = apiKeyData.tier;
        
        next();
    } catch (error) {
        next(error);
    }
};

// ============================================================================
// Security Middleware
// ============================================================================

/**
 * Security headers with Helmet
 */
const securityHeaders = helmet({
    contentSecurityPolicy: {
        directives: {
            defaultSrc: ["'self'"],
            styleSrc: ["'self'", "'unsafe-inline'"],
            scriptSrc: ["'self'"],
            imgSrc: ["'self'", "data:", "https:"],
            connectSrc: ["'self'", "https://api.example.com"]
        }
    },
    hsts: {
        maxAge: 31536000,
        includeSubDomains: true,
        preload: true
    }
});

/**
 * CORS configuration
 */
const corsOptions = {
    origin: (origin, callback) => {
        const allowedOrigins = [
            'http://localhost:3000',
            'http://localhost:3001',
            'https://app.example.com',
            'https://admin.example.com'
        ];
        
        if (!origin || allowedOrigins.includes(origin) || process.env.NODE_ENV === 'development') {
            callback(null, true);
        } else {
            callback(new Error('Not allowed by CORS'));
        }
    },
    credentials: true,
    optionsSuccessStatus: 200,
    methods: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS'],
    allowedHeaders: ['Content-Type', 'Authorization', 'X-Request-ID', 'X-API-Key'],
    exposedHeaders: ['X-Request-ID', 'X-RateLimit-Limit', 'X-RateLimit-Remaining']
};

// ============================================================================
// Request Validation Middleware
// ============================================================================

/**
 * Validate request body against schema
 */
const validateRequest = (schema, property = 'body') => {
    return (req, res, next) => {
        const { error, value } = schema.validate(req[property], {
            abortEarly: false,
            stripUnknown: true
        });
        
        if (error) {
            const errors = error.details.map(detail => ({
                field: detail.path.join('.'),
                message: detail.message
            }));
            
            logger.debug('Validation error', { errors, path: req.path });
            
            return res.status(400).json({
                success: false,
                message: 'Validation failed',
                errors
            });
        }
        
        // Replace with validated value
        req[property] = value;
        next();
    };
};

/**
 * Validate query parameters
 */
const validateQuery = (schema) => validateRequest(schema, 'query');

/**
 * Validate URL parameters
 */
const validateParams = (schema) => validateRequest(schema, 'params');

// ============================================================================
// Request Parsing Middleware
// ============================================================================

/**
 * Parse user agent
 */
const parseUserAgent = useragent.express();

/**
 * Get client IP
 */
const getClientIp = requestIp.mw();

/**
 * Add request timestamp
 */
const requestTime = (req, res, next) => {
    req.requestTime = new Date().toISOString();
    next();
};

// ============================================================================
// Response Compression
// ============================================================================

/**
 * Compress responses
 */
const compress = compression({
    level: 6, // Compression level
    threshold: 1024, // Only compress responses > 1kb
    filter: (req, res) => {
        // Don't compress if client doesn't accept gzip
        if (!req.headers['accept-encoding']?.includes('gzip')) {
            return false;
        }
        // Use compression filter from compression module
        return compression.filter(req, res);
    }
});

// ============================================================================
// Data Sanitization
// ============================================================================

/**
 * Sanitize against NoSQL injection
 */
const sanitizeData = mongoSanitize({
    replaceWith: '_',
    onSanitize: ({ req, key }) => {
        logger.warn(`NoSQL injection attempt blocked on ${key}`, {
            ip: req.clientIp,
            path: req.path
        });
    }
});

/**
 * Sanitize against XSS
 */
const preventXss = xss();

/**
 * Prevent parameter pollution
 */
const preventHpp = hpp({
    whitelist: [
        'sort',
        'fields',
        'page',
        'limit',
        'search',
        'category',
        'minPrice',
        'maxPrice'
    ]
});

// ============================================================================
// Cache Middleware
// ============================================================================

/**
 * Cache GET responses
 */
const cacheResponse = (duration = 300) => {
    return async (req, res, next) => {
        // Only cache GET requests
        if (req.method !== 'GET') {
            return next();
        }
        
        // Skip if user is authenticated (don't cache personalized data)
        if (req.user) {
            return next();
        }
        
        const key = `cache:${req.originalUrl || req.url}`;
        
        try {
            const cachedResponse = await cache.get(key);
            
            if (cachedResponse) {
                const parsed = JSON.parse(cachedResponse);
                return res.status(200).json(parsed);
            }
            
            // Store original send function
            const originalSend = res.json;
            
            // Override json method to cache response
            res.json = function(data) {
                // Cache the response
                cache.set(key, JSON.stringify(data), duration);
                
                // Call original send
                originalSend.call(this, data);
            };
            
            next();
        } catch (error) {
            logger.error('Cache middleware error', { error: error.message });
            next();
        }
    };
};

/**
 * Clear cache for patterns
 */
const clearCache = (patterns) => {
    return async (req, res, next) => {
        // Store original send function
        const originalSend = res.json;
        
        res.json = async function(data) {
            // Clear cache patterns after successful write operations
            if (res.statusCode >= 200 && res.statusCode < 300) {
                for (const pattern of patterns) {
                    await cache.delPattern(pattern);
                }
            }
            
            // Call original send
            originalSend.call(this, data);
        };
        
        next();
    };
};

// ============================================================================
// Metrics Middleware
// ============================================================================

/**
 * Track active requests
 */
const trackActiveRequests = (req, res, next) => {
    metrics.activeRequests.inc();
    
    res.on('finish', () => {
        metrics.activeRequests.dec();
    });
    
    next();
};

/**
 * Track response size
 */
const trackResponseSize = (req, res, next) => {
    const originalWrite = res.write;
    const originalEnd = res.end;
    const chunks = [];
    
    res.write = function(chunk) {
        chunks.push(chunk);
        originalWrite.apply(res, arguments);
    };
    
    res.end = function(chunk) {
        if (chunk) {
            chunks.push(chunk);
        }
        
        const responseSize = Buffer.concat(chunks).length;
        metrics.responseSize
            .labels(req.method, req.route?.path || req.path)
            .observe(responseSize);
        
        originalEnd.apply(res, arguments);
    };
    
    next();
};

// ============================================================================
// Error Handling Middleware
// ============================================================================

/**
 * 404 handler
 */
const notFound = (req, res, next) => {
    const error = new AppError(`Route ${req.originalUrl} not found`, 404);
    next(error);
};

/**
 * Global error handler
 */
const errorHandler = (err, req, res, next) => {
    const statusCode = err.statusCode || 500;
    const isOperational = err.isOperational || false;
    
    // Log error
    const logData = {
        requestId: req.id,
        method: req.method,
        url: req.originalUrl,
        statusCode,
        message: err.message,
        stack: err.stack,
        userId: req.user?.id,
        ip: req.clientIp
    };
    
    if (statusCode >= 500) {
        logger.error('Server error', logData);
        
        // Send alert for critical errors
        if (process.env.NODE_ENV === 'production') {
            // Send to error tracking service
            // Sentry.captureException(err, { extra: logData });
        }
    } else {
        logger.warn('Client error', logData);
    }
    
    // Track error metrics
    metrics.errorsTotal
        .labels(req.method, req.route?.path || req.path, statusCode)
        .inc();
    
    // Send response
    res.status(statusCode).json({
        success: false,
        message: isOperational ? err.message : 'Internal server error',
        error: process.env.NODE_ENV === 'development' ? {
            stack: err.stack,
            details: err.details
        } : undefined,
        requestId: req.id,
        timestamp: new Date().toISOString()
    });
};

// ============================================================================
// Maintenance Mode
// ============================================================================

/**
 * Check if application is in maintenance mode
 */
const maintenanceMode = (req, res, next) => {
    if (process.env.MAINTENANCE_MODE === 'true') {
        // Allow health checks during maintenance
        if (req.path === '/health' || req.path === '/ready' || req.path === '/live') {
            return next();
        }
        
        return res.status(503).json({
            success: false,
            message: 'Service temporarily unavailable due to maintenance',
            estimatedDowntime: process.env.MAINTENANCE_ETA || 'unknown'
        });
    }
    next();
};

// ============================================================================
// Request Throttling
// ============================================================================

/**
 * Throttle requests based on user tier
 */
const throttleByTier = (req, res, next) => {
    const limits = {
        free: { windowMs: 60 * 1000, max: 10 },
        basic: { windowMs: 60 * 1000, max: 60 },
        premium: { windowMs: 60 * 1000, max: 300 },
        enterprise: { windowMs: 60 * 1000, max: 1000 }
    };
    
    const tier = req.tier || 'free';
    const limit = limits[tier];
    
    if (!limit) {
        return next();
    }
    
    // Implement sliding window throttling
    const key = `throttle:${tier}:${req.user?.id || req.clientIp}`;
    
    cache.get(key, (err, count) => {
        if (err) return next();
        
        const currentCount = parseInt(count) || 0;
        
        if (currentCount >= limit.max) {
            return res.status(429).json({
                success: false,
                message: `Rate limit exceeded for ${tier} tier`
            });
        }
        
        cache.incr(key);
        cache.expire(key, limit.windowMs / 1000);
        
        next();
    });
};

// ============================================================================
// Request Timeout
// ============================================================================

/**
 * Set timeout for requests
 */
const timeout = (duration = 30000) => {
    return (req, res, next) => {
        // Set timeout
        req.setTimeout(duration, () => {
            const err = new AppError('Request timeout', 408);
            next(err);
        });
        
        next();
    };
};

// ============================================================================
// Export all middleware
// ============================================================================

module.exports = {
    // Core middleware
    requestId,
    requestLogger,
    requestTime,
    getClientIp,
    parseUserAgent,
    compress,
    
    // Security
    securityHeaders,
    cors: () => cors(corsOptions),
    sanitizeData,
    preventXss,
    preventHpp,
    
    // Auth
    authenticate,
    optionalAuth,
    authorize,
    checkOwnership,
    validateApiKey,
    
    // Rate limiting
    limiter,
    authLimiter,
    apiKeyLimiter,
    throttleByTier,
    
    // Validation
    validateRequest,
    validateQuery,
    validateParams,
    
    // Caching
    cacheResponse,
    clearCache,
    
    // Metrics
    trackActiveRequests,
    trackResponseSize,
    
    // Error handling
    notFound,
    errorHandler,
    
    // Utility
    maintenanceMode,
    timeout
};