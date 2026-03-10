/**
 * Server Module
 * 
 * Main entry point for the Node.js API server.
 * Configures Express, middleware, database connections,
 * error handling, and server lifecycle management.
 */

// ============================================================================
// Core Dependencies
// ============================================================================

const express = require('express');
const http = require('http');
const https = require('https');
const fs = require('fs');
const path = require('path');
const cluster = require('cluster');
const os = require('os');

// ============================================================================
// Environment Configuration
// ============================================================================

require('dotenv').config({ path: path.join(__dirname, '../../.env') });
require('dotenv-expand').expand(process.env);

const env = process.env.NODE_ENV || 'development';
const isDev = env === 'development';
const isProd = env === 'production';
const isTest = env === 'test';
const isStaging = env === 'staging';

// Load environment-specific config
const config = require('../config')[env];

// ============================================================================
// Application Metrics & Monitoring
// ============================================================================

// APM (Application Performance Monitoring)
if (isProd && process.env.ELASTIC_APM_ENABLED === 'true') {
    require('elastic-apm-node').start({
        serviceName: process.env.APM_SERVICE_NAME || 'nodejs-api',
        serverUrl: process.env.APM_SERVER_URL,
        environment: env
    });
}

// New Relic (optional)
if (isProd && process.env.NEW_RELIC_ENABLED === 'true') {
    require('newrelic');
}

// ============================================================================
// OpenTelemetry Setup
// ============================================================================

if (process.env.OTEL_ENABLED === 'true') {
    const { NodeSDK } = require('@opentelemetry/sdk-node');
    const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');
    const { JaegerExporter } = require('@opentelemetry/exporter-jaeger');
    
    const sdk = new NodeSDK({
        traceExporter: new JaegerExporter({
            endpoint: process.env.JAEGER_ENDPOINT || 'http://localhost:14268/api/traces',
        }),
        instrumentations: [getNodeAutoInstrumentations()],
        serviceName: process.env.OTEL_SERVICE_NAME || 'nodejs-api',
    });
    
    sdk.start();
    
    process.on('SIGTERM', () => {
        sdk.shutdown()
            .then(() => console.log('Tracing terminated'))
            .catch((error) => console.log('Error terminating tracing', error))
            .finally(() => process.exit(0));
    });
}

// ============================================================================
// Import Dependencies
// ============================================================================

// Third-party middleware
const compression = require('compression');
const helmet = require('helmet');
const cors = require('cors');
const morgan = require('morgan');
const cookieParser = require('cookie-parser');
const session = require('express-session');
const RedisStore = require('connect-redis').default;
const methodOverride = require('method-override');
const useragent = require('express-useragent');
const requestIp = require('request-ip');
const responseTime = require('response-time');
const slowDown = require('express-slow-down');

// Database
const mongoose = require('mongoose');
const redis = require('redis');
const { PrismaClient } = require('@prisma/client');

// Queue
const Queue = require('bull');
const { createBullBoard } = require('@bull-board/api');
const { BullAdapter } = require('@bull-board/api/bullAdapter');
const { ExpressAdapter } = require('@bull-board/express');

// Security
const rateLimit = require('express-rate-limit');
const mongoSanitize = require('express-mongo-sanitize');
const xss = require('xss-clean');
const hpp = require('hpp');
const { expressCspHeader, INLINE, NONE, SELF } = require('express-csp-header');

// Utilities
const { v4: uuidv4 } = require('uuid');
const moment = require('moment');
const chalk = require('chalk');

// ============================================================================
// Local Imports
// ============================================================================

const { router, socketRoutes } = require('./routes');
const { errorHandler, notFound } = require('./middleware/errorHandler');
const { requestLogger, requestId, requestTime } = require('./middleware');
const { logger, stream } = require('./utils/logger');
const { metrics, metricsMiddleware } = require('./utils/metrics');
const { cache } = require('./utils/cache');
const { queue } = require('./utils/queue');
const { emailService } = require('./services/emailService');
const { stripe } = require('./services/paymentService');
const { s3 } = require('./services/storageService');

// ============================================================================
// Database Connections
// ============================================================================

/**
 * Connect to MongoDB
 */
const connectMongoDB = async () => {
    try {
        const conn = await mongoose.connect(process.env.MONGODB_URI || config.mongodb.uri, {
            useNewUrlParser: true,
            useUnifiedTopology: true,
            maxPoolSize: config.mongodb.maxPoolSize || 10,
            minPoolSize: config.mongodb.minPoolSize || 2,
            socketTimeoutMS: config.mongodb.socketTimeoutMS || 45000,
            connectTimeoutMS: config.mongodb.connectTimeoutMS || 10000,
            serverSelectionTimeoutMS: config.mongodb.serverSelectionTimeoutMS || 5000,
            heartbeatFrequencyMS: config.mongodb.heartbeatFrequencyMS || 10000,
            retryWrites: true,
            retryReads: true,
        });
        
        logger.info(`MongoDB connected: ${conn.connection.host}`);
        
        mongoose.connection.on('error', (err) => {
            logger.error('MongoDB connection error:', err);
        });
        
        mongoose.connection.on('disconnected', () => {
            logger.warn('MongoDB disconnected');
        });
        
        mongoose.connection.on('reconnected', () => {
            logger.info('MongoDB reconnected');
        });
        
        return conn;
    } catch (error) {
        logger.error('MongoDB connection failed:', error);
        if (isProd) {
            // In production, retry connection
            setTimeout(connectMongoDB, 5000);
        } else {
            process.exit(1);
        }
    }
};

/**
 * Connect to Redis
 */
const connectRedis = async () => {
    try {
        const redisClient = redis.createClient({
            url: process.env.REDIS_URL || config.redis.url,
            socket: {
                reconnectStrategy: (retries) => {
                    if (retries > 10) {
                        logger.error('Redis max retries reached');
                        return new Error('Redis max retries reached');
                    }
                    return Math.min(retries * 100, 3000);
                },
                connectTimeout: config.redis.connectTimeout || 10000,
                keepAlive: config.redis.keepAlive || 30000,
            },
            password: process.env.REDIS_PASSWORD,
            database: config.redis.db || 0,
        });
        
        redisClient.on('error', (err) => {
            logger.error('Redis error:', err);
        });
        
        redisClient.on('connect', () => {
            logger.info('Redis connected');
        });
        
        redisClient.on('ready', () => {
            logger.info('Redis ready');
        });
        
        await redisClient.connect();
        
        // Initialize cache with Redis client
        cache.init(redisClient);
        
        return redisClient;
    } catch (error) {
        logger.error('Redis connection failed:', error);
        if (isProd) {
            setTimeout(connectRedis, 5000);
        }
        return null;
    }
};

/**
 * Connect to Prisma (for PostgreSQL)
 */
const initPrisma = () => {
    const prisma = new PrismaClient({
        log: isDev ? ['query', 'info', 'warn', 'error'] : ['error'],
        errorFormat: 'pretty',
    });
    
    prisma.$on('beforeExit', async () => {
        logger.info('Prisma before exit');
    });
    
    return prisma;
};

// ============================================================================
// Express App Initialization
// ============================================================================

const app = express();

// ============================================================================
// Security Middleware
// ============================================================================

// Helmet for security headers
app.use(helmet({
    contentSecurityPolicy: {
        directives: {
            defaultSrc: ["'self'"],
            styleSrc: ["'self'", "'unsafe-inline'"],
            scriptSrc: ["'self'"],
            imgSrc: ["'self'", "data:", "https:"],
            connectSrc: ["'self'", "https://api.example.com"],
            fontSrc: ["'self'"],
            objectSrc: ["'none'"],
            mediaSrc: ["'self'"],
            frameSrc: ["'none'"],
        },
    },
    crossOriginEmbedderPolicy: false,
    crossOriginResourcePolicy: { policy: "cross-origin" },
}));

// CSP headers (alternative to helmet's CSP)
app.use(expressCspHeader({
    directives: {
        'default-src': [SELF],
        'style-src': [SELF, INLINE],
        'script-src': [SELF],
        'img-src': [SELF, 'data:', 'https://*.example.com'],
        'font-src': [SELF],
        'connect-src': [SELF],
        'frame-src': [NONE],
    },
}));

// CORS configuration
const corsOptions = {
    origin: (origin, callback) => {
        const allowedOrigins = config.cors.origins || [
            'http://localhost:3000',
            'http://localhost:3001',
            'https://app.example.com',
        ];
        
        if (!origin || allowedOrigins.includes(origin) || isDev) {
            callback(null, true);
        } else {
            callback(new Error('Not allowed by CORS'));
        }
    },
    credentials: true,
    optionsSuccessStatus: 200,
    methods: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS'],
    allowedHeaders: [
        'Content-Type',
        'Authorization',
        'X-Request-ID',
        'X-API-Key',
        'X-CSRF-Token',
        'Accept',
        'Origin',
        'User-Agent',
    ],
    exposedHeaders: [
        'X-Request-ID',
        'X-RateLimit-Limit',
        'X-RateLimit-Remaining',
        'X-RateLimit-Reset',
    ],
    maxAge: 86400, // 24 hours
};

app.use(cors(corsOptions));

// Rate limiting
const limiter = rateLimit({
    windowMs: config.rateLimit.windowMs || 15 * 60 * 1000, // 15 minutes
    max: config.rateLimit.max || 100,
    message: {
        success: false,
        message: 'Too many requests, please try again later.',
    },
    standardHeaders: true,
    legacyHeaders: false,
    keyGenerator: (req) => {
        return req.user?.id || req.clientIp;
    },
    skip: (req) => {
        return req.path === '/health' || req.path === '/metrics';
    },
    handler: (req, res) => {
        logger.warn(`Rate limit exceeded for IP: ${req.clientIp}`);
        res.status(429).json({
            success: false,
            message: 'Too many requests, please try again later.',
        });
    },
});

// Apply rate limiting to all routes
app.use('/api', limiter);

// Slower down (gradual rate limiting)
const speedLimiter = slowDown({
    windowMs: 15 * 60 * 1000,
    delayAfter: 50,
    delayMs: 500,
});

app.use('/api', speedLimiter);

// Data sanitization against NoSQL injection
app.use(mongoSanitize({
    replaceWith: '_',
    onSanitize: ({ req, key }) => {
        logger.warn(`NoSQL injection attempt blocked on ${key}`, {
            ip: req.clientIp,
            path: req.path,
        });
    },
}));

// Data sanitization against XSS
app.use(xss());

// Prevent parameter pollution
app.use(hpp({
    whitelist: [
        'sort',
        'fields',
        'page',
        'limit',
        'search',
        'category',
        'minPrice',
        'maxPrice',
        'status',
        'role',
    ],
}));

// ============================================================================
// Standard Middleware
// ============================================================================

// Compression
app.use(compression({
    level: 6,
    threshold: 1024,
    filter: (req, res) => {
        if (req.headers['x-no-compression']) {
            return false;
        }
        return compression.filter(req, res);
    },
}));

// Body parsing
app.use(express.json({ limit: config.bodyLimit || '10mb' }));
app.use(express.urlencoded({ extended: true, limit: config.bodyLimit || '10mb' }));
app.use(express.raw({ type: 'application/octet-stream', limit: '50mb' }));
app.use(express.text({ type: 'text/plain' }));

// Cookie parsing
app.use(cookieParser(process.env.COOKIE_SECRET));

// Session management
if (config.session.enabled) {
    const sessionStore = new RedisStore({
        client: redisClient,
        prefix: 'sess:',
        ttl: config.session.ttl || 86400,
    });
    
    app.use(session({
        store: sessionStore,
        secret: process.env.SESSION_SECRET || config.session.secret,
        resave: false,
        saveUninitialized: false,
        name: 'sessionId',
        cookie: {
            secure: isProd,
            httpOnly: true,
            maxAge: config.session.ttl || 86400 * 1000,
            sameSite: 'lax',
            domain: config.session.domain || undefined,
        },
        rolling: true,
    }));
}

// Method override (for PUT/DELETE from forms)
app.use(methodOverride('_method'));

// User agent parsing
app.use(useragent.express());

// Request IP
app.use(requestIp.mw());

// Request ID
app.use(requestId);

// Request timestamp
app.use(requestTime);

// Response time tracking
app.use(responseTime((req, res, time) => {
    metrics.httpRequestDuration
        .labels(req.method, req.route?.path || req.path, res.statusCode)
        .observe(time / 1000);
}));

// ============================================================================
// Logging Middleware
// ============================================================================

// Morgan for HTTP request logging
const morganFormat = isDev ? 'dev' : 'combined';
app.use(morgan(morganFormat, { stream }));

// Custom request logger
app.use(requestLogger);

// ============================================================================
// Metrics Middleware
// ============================================================================

app.use(metricsMiddleware);

// Expose metrics endpoint
app.get('/metrics', async (req, res) => {
    res.set('Content-Type', metrics.register.contentType);
    res.end(await metrics.register.metrics());
});

// ============================================================================
// Health Check Endpoints
// ============================================================================

app.get('/health', (req, res) => {
    res.status(200).json({
        status: 'healthy',
        timestamp: new Date().toISOString(),
        uptime: process.uptime(),
        memory: process.memoryUsage(),
        cpu: process.cpuUsage(),
        version: process.env.npm_package_version || '1.0.0',
    });
});

app.get('/health/readiness', (req, res) => {
    // Check database connections
    const mongoState = mongoose.connection.readyState;
    const redisState = redisClient?.isReady ? 1 : 0;
    
    const isReady = mongoState === 1 && redisState === 1;
    
    res.status(isReady ? 200 : 503).json({
        status: isReady ? 'ready' : 'not ready',
        checks: {
            mongodb: mongoState === 1 ? 'connected' : 'disconnected',
            redis: redisState === 1 ? 'connected' : 'disconnected',
        },
        timestamp: new Date().toISOString(),
    });
});

app.get('/health/liveness', (req, res) => {
    res.status(200).json({
        status: 'alive',
        timestamp: new Date().toISOString(),
    });
});

app.get('/health/detailed', async (req, res) => {
    const checks = {
        server: 'healthy',
        timestamp: new Date().toISOString(),
    };
    
    // Check MongoDB
    try {
        await mongoose.connection.db.admin().ping();
        checks.mongodb = 'healthy';
    } catch (error) {
        checks.mongodb = 'unhealthy';
        checks.mongodbError = error.message;
    }
    
    // Check Redis
    try {
        await redisClient.ping();
        checks.redis = 'healthy';
    } catch (error) {
        checks.redis = 'unhealthy';
        checks.redisError = error.message;
    }
    
    // Check queue
    try {
        const queueStats = await queue.getJobCounts();
        checks.queue = 'healthy';
        checks.queueStats = queueStats;
    } catch (error) {
        checks.queue = 'unhealthy';
        checks.queueError = error.message;
    }
    
    const isHealthy = checks.mongodb === 'healthy' && 
                     checks.redis === 'healthy' && 
                     checks.queue === 'healthy';
    
    res.status(isHealthy ? 200 : 503).json({
        status: isHealthy ? 'healthy' : 'degraded',
        checks,
    });
});

// ============================================================================
// Bull Board for Queue Monitoring
// ============================================================================

const serverAdapter = new ExpressAdapter();
serverAdapter.setBasePath('/admin/queues');

createBullBoard({
    queues: [
        new BullAdapter(queue),
        // Add more queues as needed
    ],
    serverAdapter,
});

app.use('/admin/queues', authenticate, authorize('admin'), serverAdapter.getRouter());

// ============================================================================
// API Routes
// ============================================================================

// Mount API routes
app.use('/api/v1', router);

// ============================================================================
// Static Files
// ============================================================================

app.use('/static', express.static(path.join(__dirname, '../../public'), {
    maxAge: isProd ? '30d' : 0,
    etag: true,
    lastModified: true,
    setHeaders: (res, path) => {
        if (path.endsWith('.html')) {
            res.setHeader('Cache-Control', 'no-cache');
        }
    },
}));

app.use('/uploads', express.static(path.join(__dirname, '../../uploads'), {
    maxAge: isProd ? '7d' : 0,
}));

// ============================================================================
// 404 Handler
// ============================================================================

app.use(notFound);

// ============================================================================
// Global Error Handler
// ============================================================================

app.use(errorHandler);

// ============================================================================
// Server Creation
// ============================================================================

const createServer = () => {
    let server;
    
    if (config.https.enabled && isProd) {
        // HTTPS server
        const httpsOptions = {
            key: fs.readFileSync(config.https.keyPath),
            cert: fs.readFileSync(config.https.certPath),
            ca: config.https.caPath ? fs.readFileSync(config.https.caPath) : undefined,
            requestCert: config.https.requestCert || false,
            rejectUnauthorized: config.https.rejectUnauthorized || false,
        };
        
        server = https.createServer(httpsOptions, app);
    } else {
        // HTTP server
        server = http.createServer(app);
    }
    
    return server;
};

// ============================================================================
// Socket.IO Setup
// ============================================================================

const setupSocketIO = (server) => {
    const io = require('socket.io')(server, {
        cors: {
            origin: config.cors.origins,
            credentials: true,
        },
        path: '/socket.io',
        serveClient: false,
        pingInterval: 10000,
        pingTimeout: 5000,
        cookie: false,
        transports: ['websocket', 'polling'],
    });
    
    // Authentication middleware for socket.io
    io.use(async (socket, next) => {
        try {
            const token = socket.handshake.auth.token;
            // Verify token logic
            // const user = await verifyToken(token);
            // socket.user = user;
            next();
        } catch (error) {
            next(new Error('Authentication error'));
        }
    });
    
    // Initialize socket routes
    socketRoutes(io);
    
    return io;
};

// ============================================================================
// Graceful Shutdown
// ============================================================================

const gracefulShutdown = (server, io) => {
    logger.info('Received shutdown signal');
    
    // Stop accepting new connections
    server.close(() => {
        logger.info('HTTP server closed');
    });
    
    // Close Socket.IO
    if (io) {
        io.close(() => {
            logger.info('Socket.IO closed');
        });
    }
    
    // Close database connections
    mongoose.connection.close(false, () => {
        logger.info('MongoDB connection closed');
    });
    
    if (redisClient) {
        redisClient.quit().then(() => {
            logger.info('Redis connection closed');
        });
    }
    
    // Close queue connections
    queue.close().then(() => {
        logger.info('Queue connections closed');
    });
    
    // Force exit after timeout
    setTimeout(() => {
        logger.error('Could not close connections in time, forcefully shutting down');
        process.exit(1);
    }, config.shutdownTimeout || 30000);
};

// ============================================================================
// Cluster Mode
// ============================================================================

const startCluster = () => {
    const numCPUs = os.cpus().length;
    const workers = Math.min(numCPUs, config.maxWorkers || numCPUs);
    
    if (cluster.isMaster) {
        logger.info(`Master ${process.pid} is running`);
        logger.info(`Starting ${workers} workers`);
        
        // Fork workers
        for (let i = 0; i < workers; i++) {
            cluster.fork();
        }
        
        cluster.on('exit', (worker, code, signal) => {
            logger.warn(`Worker ${worker.process.pid} died. Restarting...`);
            cluster.fork();
        });
        
        cluster.on('online', (worker) => {
            logger.info(`Worker ${worker.process.pid} is online`);
        });
        
    } else {
        startServer();
    }
};

// ============================================================================
// Server Startup
// ============================================================================

const startServer = async () => {
    try {
        // Connect to databases
        await connectMongoDB();
        const redisClient = await connectRedis();
        const prisma = initPrisma();
        
        // Create server
        const server = createServer();
        const io = setupSocketIO(server);
        
        // Store instances in app locals
        app.locals.db = mongoose.connection;
        app.locals.redis = redisClient;
        app.locals.prisma = prisma;
        app.locals.io = io;
        app.locals.queue = queue;
        
        // Start server
        const port = process.env.PORT || config.port || 3000;
        const host = process.env.HOST || config.host || '0.0.0.0';
        
        server.listen(port, host, () => {
            logger.info(`Server running on ${host}:${port} in ${env} mode`);
            logger.info(`Worker ${process.pid} started`);
            
            // Send ready signal to PM2
            if (typeof process.send === 'function') {
                process.send('ready');
            }
        });
        
        // Handle graceful shutdown
        process.on('SIGTERM', () => gracefulShutdown(server, io));
        process.on('SIGINT', () => gracefulShutdown(server, io));
        process.on('SIGQUIT', () => gracefulShutdown(server, io));
        
        // Handle uncaught exceptions
        process.on('uncaughtException', (error) => {
            logger.error('Uncaught Exception:', error);
            // In production, you might want to restart the process
            if (isProd) {
                gracefulShutdown(server, io);
            }
        });
        
        process.on('unhandledRejection', (reason, promise) => {
            logger.error('Unhandled Rejection at:', promise, 'reason:', reason);
        });
        
    } catch (error) {
        logger.error('Failed to start server:', error);
        process.exit(1);
    }
};

// ============================================================================
// Start Application
// ============================================================================

if (config.clusterMode && !isTest) {
    startCluster();
} else {
    startServer();
}

// ============================================================================
// Export for testing
// ============================================================================

module.exports = { app, startServer };