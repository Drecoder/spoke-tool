/**
 * Controllers Module
 * 
 * Handles HTTP request/response logic for the Node.js API.
 * Includes validation, error handling, and response formatting.
 */

const { validationResult } = require('express-validator');
const { User, Product, Order } = require('../models');
const { AppError, catchAsync } = require('../utils/errorHandler');
const { logger } = require('../utils/logger');
const { cache } = require('../utils/cache');
const { queue } = require('../utils/queue');
const { emailService } = require('../services/emailService');
const { stripe } = require('../services/paymentService');
const { s3 } = require('../services/storageService');

// ============================================================================
// Response Formatter
// ============================================================================

/**
 * Format success response
 */
const sendResponse = (res, data, statusCode = 200, message = 'Success') => {
    return res.status(statusCode).json({
        success: true,
        message,
        data,
        timestamp: new Date().toISOString()
    });
};

/**
 * Format paginated response
 */
const sendPaginated = (res, data, page, limit, total, statusCode = 200) => {
    return res.status(statusCode).json({
        success: true,
        data,
        pagination: {
            page: parseInt(page),
            limit: parseInt(limit),
            total,
            pages: Math.ceil(total / limit)
        },
        timestamp: new Date().toISOString()
    });
};

// ============================================================================
// User Controllers
// ============================================================================

/**
 * User Controller
 * Handles all user-related operations
 */
const userController = {
    /**
     * Get all users with pagination and filtering
     */
    getAllUsers: catchAsync(async (req, res) => {
        const { page = 1, limit = 10, role, search, sortBy = 'createdAt', sortOrder = 'desc' } = req.query;
        
        const query = {};
        
        // Apply filters
        if (role) query.role = role;
        if (search) {
            query.$or = [
                { name: { $regex: search, $options: 'i' } },
                { email: { $regex: search, $options: 'i' } }
            ];
        }
        
        // Calculate pagination
        const skip = (parseInt(page) - 1) * parseInt(limit);
        
        // Execute query
        const users = await User.find(query)
            .sort({ [sortBy]: sortOrder === 'desc' ? -1 : 1 })
            .skip(skip)
            .limit(parseInt(limit))
            .select('-password -__v');
        
        const total = await User.countDocuments(query);
        
        // Log the action
        logger.info(`Users fetched by admin: ${req.user.email}`, { count: users.length });
        
        return sendPaginated(res, users, page, limit, total);
    }),

    /**
     * Get user by ID
     */
    getUserById: catchAsync(async (req, res) => {
        const { id } = req.params;
        
        // Check cache first
        const cachedUser = await cache.get(`user:${id}`);
        if (cachedUser) {
            logger.debug(`Cache hit for user: ${id}`);
            return sendResponse(res, JSON.parse(cachedUser));
        }
        
        const user = await User.findById(id).select('-password -__v');
        
        if (!user) {
            throw new AppError('User not found', 404);
        }
        
        // Store in cache (expires in 5 minutes)
        await cache.set(`user:${id}`, JSON.stringify(user), 300);
        
        logger.debug(`User fetched: ${id}`);
        return sendResponse(res, user);
    }),

    /**
     * Create new user
     */
    createUser: catchAsync(async (req, res) => {
        const errors = validationResult(req);
        if (!errors.isEmpty()) {
            return res.status(400).json({ errors: errors.array() });
        }
        
        const { email, password, name, role = 'user' } = req.body;
        
        // Check if user exists
        const existingUser = await User.findOne({ email });
        if (existingUser) {
            throw new AppError('Email already registered', 409);
        }
        
        // Create user
        const user = new User({
            email,
            password,
            name,
            role
        });
        
        await user.save();
        
        // Remove password from response
        user.password = undefined;
        
        // Send welcome email asynchronously
        queue.add('sendWelcomeEmail', { email, name });
        
        logger.info(`New user created: ${email}`);
        
        return sendResponse(res, user, 201, 'User created successfully');
    }),

    /**
     * Update user
     */
    updateUser: catchAsync(async (req, res) => {
        const { id } = req.params;
        const updates = req.body;
        
        // Remove fields that shouldn't be updated directly
        delete updates.password;
        delete updates.role;
        delete updates._id;
        
        const user = await User.findByIdAndUpdate(
            id,
            { ...updates, updatedAt: new Date() },
            { new: true, runValidators: true }
        ).select('-password -__v');
        
        if (!user) {
            throw new AppError('User not found', 404);
        }
        
        // Invalidate cache
        await cache.del(`user:${id}`);
        
        logger.info(`User updated: ${id}`);
        return sendResponse(res, user, 200, 'User updated successfully');
    }),

    /**
     * Delete user
     */
    deleteUser: catchAsync(async (req, res) => {
        const { id } = req.params;
        
        const user = await User.findByIdAndDelete(id);
        
        if (!user) {
            throw new AppError('User not found', 404);
        }
        
        // Clean up user data
        await Promise.all([
            cache.del(`user:${id}`),
            Order.deleteMany({ userId: id }),
            queue.add('cleanupUserData', { userId: id })
        ]);
        
        logger.warn(`User deleted: ${id} by ${req.user.email}`);
        return sendResponse(res, null, 204);
    }),

    /**
     * Update user profile (self)
     */
    updateProfile: catchAsync(async (req, res) => {
        const userId = req.user.id;
        const updates = req.body;
        
        // Remove sensitive fields
        delete updates.role;
        delete updates.email;
        
        const user = await User.findByIdAndUpdate(
            userId,
            { ...updates, updatedAt: new Date() },
            { new: true, runValidators: true }
        ).select('-password -__v');
        
        await cache.del(`user:${userId}`);
        
        logger.info(`Profile updated: ${userId}`);
        return sendResponse(res, user, 200, 'Profile updated successfully');
    }),

    /**
     * Change password
     */
    changePassword: catchAsync(async (req, res) => {
        const { currentPassword, newPassword } = req.body;
        const userId = req.user.id;
        
        const user = await User.findById(userId);
        
        // Verify current password
        const isValid = await user.comparePassword(currentPassword);
        if (!isValid) {
            throw new AppError('Current password is incorrect', 401);
        }
        
        // Update password
        user.password = newPassword;
        await user.save();
        
        // Invalidate all sessions (if using session tokens)
        await cache.delPattern(`session:${userId}:*`);
        
        logger.info(`Password changed: ${userId}`);
        return sendResponse(res, null, 200, 'Password changed successfully');
    })
};

// ============================================================================
// Product Controllers
// ============================================================================

/**
 * Product Controller
 * Handles all product-related operations
 */
const productController = {
    /**
     * Get all products with filtering and pagination
     */
    getAllProducts: catchAsync(async (req, res) => {
        const { 
            page = 1, 
            limit = 20, 
            category, 
            minPrice, 
            maxPrice,
            inStock,
            search,
            sortBy = 'createdAt',
            sortOrder = 'desc'
        } = req.query;
        
        const query = {};
        
        // Build filter query
        if (category) query.category = category;
        if (minPrice || maxPrice) {
            query.price = {};
            if (minPrice) query.price.$gte = parseFloat(minPrice);
            if (maxPrice) query.price.$lte = parseFloat(maxPrice);
        }
        if (inStock === 'true') query.stock = { $gt: 0 };
        if (search) {
            query.$or = [
                { name: { $regex: search, $options: 'i' } },
                { description: { $regex: search, $options: 'i' } }
            ];
        }
        
        const skip = (parseInt(page) - 1) * parseInt(limit);
        
        const products = await Product.find(query)
            .populate('category', 'name slug')
            .sort({ [sortBy]: sortOrder === 'desc' ? -1 : 1 })
            .skip(skip)
            .limit(parseInt(limit));
        
        const total = await Product.countDocuments(query);
        
        return sendPaginated(res, products, page, limit, total);
    }),

    /**
     * Get product by ID
     */
    getProductById: catchAsync(async (req, res) => {
        const { id } = req.params;
        
        // Try cache first
        const cachedProduct = await cache.get(`product:${id}`);
        if (cachedProduct) {
            return sendResponse(res, JSON.parse(cachedProduct));
        }
        
        const product = await Product.findById(id)
            .populate('category', 'name slug')
            .populate('reviews');
        
        if (!product) {
            throw new AppError('Product not found', 404);
        }
        
        // Cache for 10 minutes
        await cache.set(`product:${id}`, JSON.stringify(product), 600);
        
        return sendResponse(res, product);
    }),

    /**
     * Create product (admin only)
     */
    createProduct: catchAsync(async (req, res) => {
        const errors = validationResult(req);
        if (!errors.isEmpty()) {
            return res.status(400).json({ errors: errors.array() });
        }
        
        const productData = req.body;
        
        // Handle image uploads if any
        if (req.files && req.files.length > 0) {
            const uploadPromises = req.files.map(async (file) => {
                const url = await s3.upload(file);
                return url;
            });
            productData.images = await Promise.all(uploadPromises);
        }
        
        const product = new Product(productData);
        await product.save();
        
        logger.info(`Product created: ${product.name} by ${req.user.email}`);
        
        return sendResponse(res, product, 201, 'Product created successfully');
    }),

    /**
     * Update product (admin only)
     */
    updateProduct: catchAsync(async (req, res) => {
        const { id } = req.params;
        const updates = req.body;
        
        const product = await Product.findByIdAndUpdate(
            id,
            { ...updates, updatedAt: new Date() },
            { new: true, runValidators: true }
        );
        
        if (!product) {
            throw new AppError('Product not found', 404);
        }
        
        // Invalidate cache
        await cache.del(`product:${id}`);
        
        logger.info(`Product updated: ${id} by ${req.user.email}`);
        return sendResponse(res, product, 200, 'Product updated successfully');
    }),

    /**
     * Delete product (admin only)
     */
    deleteProduct: catchAsync(async (req, res) => {
        const { id } = req.params;
        
        const product = await Product.findByIdAndDelete(id);
        
        if (!product) {
            throw new AppError('Product not found', 404);
        }
        
        // Clean up related data
        await Promise.all([
            cache.del(`product:${id}`),
            Order.updateMany(
                { 'items.productId': id },
                { $pull: { items: { productId: id } } }
            )
        ]);
        
        logger.warn(`Product deleted: ${id} by ${req.user.email}`);
        return sendResponse(res, null, 204);
    }),

    /**
     * Update product stock
     */
    updateStock: catchAsync(async (req, res) => {
        const { id } = req.params;
        const { quantity, operation = 'set' } = req.body;
        
        const product = await Product.findById(id);
        if (!product) {
            throw new AppError('Product not found', 404);
        }
        
        // Update stock based on operation
        switch (operation) {
            case 'set':
                product.stock = quantity;
                break;
            case 'increment':
                product.stock += quantity;
                break;
            case 'decrement':
                if (product.stock < quantity) {
                    throw new AppError('Insufficient stock', 400);
                }
                product.stock -= quantity;
                break;
            default:
                throw new AppError('Invalid operation', 400);
        }
        
        await product.save();
        await cache.del(`product:${id}`);
        
        logger.info(`Stock updated for product ${id}: ${product.stock}`);
        return sendResponse(res, { stock: product.stock });
    })
};

// ============================================================================
// Order Controllers
// ============================================================================

/**
 * Order Controller
 * Handles all order-related operations
 */
const orderController = {
    /**
     * Get user orders
     */
    getUserOrders: catchAsync(async (req, res) => {
        const userId = req.user.id;
        const { page = 1, limit = 10, status } = req.query;
        
        const query = { userId };
        if (status) query.status = status;
        
        const skip = (parseInt(page) - 1) * parseInt(limit);
        
        const orders = await Order.find(query)
            .populate('items.productId', 'name price')
            .sort({ createdAt: -1 })
            .skip(skip)
            .limit(parseInt(limit));
        
        const total = await Order.countDocuments(query);
        
        return sendPaginated(res, orders, page, limit, total);
    }),

    /**
     * Get order by ID
     */
    getOrderById: catchAsync(async (req, res) => {
        const { id } = req.params;
        const userId = req.user.id;
        const isAdmin = req.user.role === 'admin';
        
        const order = await Order.findById(id)
            .populate('userId', 'name email')
            .populate('items.productId', 'name price images');
        
        if (!order) {
            throw new AppError('Order not found', 404);
        }
        
        // Check authorization (users can only see their own orders)
        if (!isAdmin && order.userId._id.toString() !== userId) {
            throw new AppError('Unauthorized', 403);
        }
        
        return sendResponse(res, order);
    }),

    /**
     * Create order
     */
    createOrder: catchAsync(async (req, res) => {
        const errors = validationResult(req);
        if (!errors.isEmpty()) {
            return res.status(400).json({ errors: errors.array() });
        }
        
        const { items, shippingAddress, paymentMethod } = req.body;
        const userId = req.user.id;
        
        // Validate items and calculate total
        let subtotal = 0;
        const orderItems = [];
        
        for (const item of items) {
            const product = await Product.findById(item.productId);
            if (!product) {
                throw new AppError(`Product ${item.productId} not found`, 404);
            }
            
            if (product.stock < item.quantity) {
                throw new AppError(`Insufficient stock for ${product.name}`, 400);
            }
            
            const itemTotal = product.price * item.quantity;
            subtotal += itemTotal;
            
            orderItems.push({
                productId: product._id,
                name: product.name,
                price: product.price,
                quantity: item.quantity,
                total: itemTotal
            });
            
            // Decrement stock
            product.stock -= item.quantity;
            await product.save();
        }
        
        // Calculate taxes and shipping
        const tax = subtotal * 0.1; // 10% tax
        const shipping = 10; // Flat shipping rate
        const total = subtotal + tax + shipping;
        
        // Create order
        const order = new Order({
            userId,
            items: orderItems,
            subtotal,
            tax,
            shipping,
            total,
            shippingAddress,
            paymentMethod,
            status: 'pending'
        });
        
        await order.save();
        
        // Process payment asynchronously
        queue.add('processPayment', {
            orderId: order._id,
            amount: total,
            paymentMethod,
            userId
        });
        
        logger.info(`Order created: ${order._id} by user ${userId}`);
        
        return sendResponse(res, order, 201, 'Order created successfully');
    }),

    /**
     * Update order status (admin only)
     */
    updateOrderStatus: catchAsync(async (req, res) => {
        const { id } = req.params;
        const { status, trackingNumber } = req.body;
        
        const order = await Order.findById(id);
        if (!order) {
            throw new AppError('Order not found', 404);
        }
        
        order.status = status;
        if (trackingNumber) {
            order.trackingNumber = trackingNumber;
        }
        
        await order.save();
        
        // Send notification to user
        queue.add('sendOrderUpdate', {
            userId: order.userId,
            orderId: order._id,
            status
        });
        
        logger.info(`Order ${id} status updated to ${status}`);
        return sendResponse(res, order, 200, 'Order status updated');
    }),

    /**
     * Cancel order
     */
    cancelOrder: catchAsync(async (req, res) => {
        const { id } = req.params;
        const userId = req.user.id;
        const isAdmin = req.user.role === 'admin';
        
        const order = await Order.findById(id);
        if (!order) {
            throw new AppError('Order not found', 404);
        }
        
        // Check authorization
        if (!isAdmin && order.userId.toString() !== userId) {
            throw new AppError('Unauthorized', 403);
        }
        
        // Check if order can be cancelled
        if (!['pending', 'processing'].includes(order.status)) {
            throw new AppError('Order cannot be cancelled at this stage', 400);
        }
        
        order.status = 'cancelled';
        await order.save();
        
        // Restore stock
        for (const item of order.items) {
            await Product.findByIdAndUpdate(item.productId, {
                $inc: { stock: item.quantity }
            });
        }
        
        // Process refund if payment was made
        if (order.paymentIntentId) {
            queue.add('processRefund', {
                orderId: order._id,
                amount: order.total
            });
        }
        
        logger.info(`Order ${id} cancelled by ${req.user.email}`);
        return sendResponse(res, order, 200, 'Order cancelled successfully');
    })
};

// ============================================================================
// Auth Controllers
// ============================================================================

/**
 * Auth Controller
 * Handles authentication and authorization
 */
const authController = {
    /**
     * Register new user
     */
    register: catchAsync(async (req, res) => {
        const errors = validationResult(req);
        if (!errors.isEmpty()) {
            return res.status(400).json({ errors: errors.array() });
        }
        
        const { email, password, name } = req.body;
        
        // Check if user exists
        const existingUser = await User.findOne({ email });
        if (existingUser) {
            throw new AppError('Email already registered', 409);
        }
        
        // Create user
        const user = new User({
            email,
            password,
            name,
            role: 'user'
        });
        
        await user.save();
        
        // Generate tokens
        const accessToken = user.generateAccessToken();
        const refreshToken = user.generateRefreshToken();
        
        // Store refresh token
        await cache.set(`refresh:${user._id}`, refreshToken, 7 * 24 * 60 * 60); // 7 days
        
        // Send welcome email
        queue.add('sendWelcomeEmail', { email, name });
        
        logger.info(`New user registered: ${email}`);
        
        return sendResponse(res, {
            user: {
                id: user._id,
                email: user.email,
                name: user.name,
                role: user.role
            },
            accessToken,
            refreshToken
        }, 201, 'Registration successful');
    }),

    /**
     * Login user
     */
    login: catchAsync(async (req, res) => {
        const { email, password } = req.body;
        
        // Find user
        const user = await User.findOne({ email }).select('+password');
        if (!user) {
            throw new AppError('Invalid credentials', 401);
        }
        
        // Check password
        const isValid = await user.comparePassword(password);
        if (!isValid) {
            throw new AppError('Invalid credentials', 401);
        }
        
        // Generate tokens
        const accessToken = user.generateAccessToken();
        const refreshToken = user.generateRefreshToken();
        
        // Store refresh token
        await cache.set(`refresh:${user._id}`, refreshToken, 7 * 24 * 60 * 60);
        
        // Update last login
        user.lastLogin = new Date();
        await user.save();
        
        logger.info(`User logged in: ${email}`);
        
        return sendResponse(res, {
            user: {
                id: user._id,
                email: user.email,
                name: user.name,
                role: user.role
            },
            accessToken,
            refreshToken
        });
    }),

    /**
     * Refresh access token
     */
    refreshToken: catchAsync(async (req, res) => {
        const { refreshToken } = req.body;
        
        if (!refreshToken) {
            throw new AppError('Refresh token required', 400);
        }
        
        // Verify refresh token
        const decoded = jwt.verify(refreshToken, process.env.REFRESH_TOKEN_SECRET);
        
        // Check if token exists in cache
        const storedToken = await cache.get(`refresh:${decoded.id}`);
        if (!storedToken || storedToken !== refreshToken) {
            throw new AppError('Invalid refresh token', 401);
        }
        
        // Generate new tokens
        const user = await User.findById(decoded.id);
        const newAccessToken = user.generateAccessToken();
        const newRefreshToken = user.generateRefreshToken();
        
        // Update stored refresh token
        await cache.set(`refresh:${user._id}`, newRefreshToken, 7 * 24 * 60 * 60);
        
        return sendResponse(res, {
            accessToken: newAccessToken,
            refreshToken: newRefreshToken
        });
    }),

    /**
     * Logout user
     */
    logout: catchAsync(async (req, res) => {
        const userId = req.user.id;
        
        // Remove refresh token
        await cache.del(`refresh:${userId}`);
        
        logger.info(`User logged out: ${userId}`);
        return sendResponse(res, null, 200, 'Logout successful');
    }),

    /**
     * Request password reset
     */
    forgotPassword: catchAsync(async (req, res) => {
        const { email } = req.body;
        
        const user = await User.findOne({ email });
        if (!user) {
            // Don't reveal if user exists
            return sendResponse(res, null, 200, 'If the email exists, a reset link has been sent');
        }
        
        // Generate reset token
        const resetToken = crypto.randomBytes(32).toString('hex');
        const resetTokenHash = crypto
            .createHash('sha256')
            .update(resetToken)
            .digest('hex');
        
        // Store reset token (expires in 1 hour)
        await cache.set(`reset:${user._id}`, resetTokenHash, 3600);
        
        // Send reset email
        const resetUrl = `${process.env.FRONTEND_URL}/reset-password?token=${resetToken}&id=${user._id}`;
        queue.add('sendPasswordResetEmail', { email, resetUrl });
        
        logger.info(`Password reset requested for: ${email}`);
        
        return sendResponse(res, null, 200, 'If the email exists, a reset link has been sent');
    }),

    /**
     * Reset password
     */
    resetPassword: catchAsync(async (req, res) => {
        const { userId, token, newPassword } = req.body;
        
        // Get stored token
        const storedTokenHash = await cache.get(`reset:${userId}`);
        if (!storedTokenHash) {
            throw new AppError('Invalid or expired reset token', 400);
        }
        
        // Verify token
        const tokenHash = crypto
            .createHash('sha256')
            .update(token)
            .digest('hex');
        
        if (tokenHash !== storedTokenHash) {
            throw new AppError('Invalid reset token', 400);
        }
        
        // Update password
        const user = await User.findById(userId);
        user.password = newPassword;
        await user.save();
        
        // Remove reset token
        await cache.del(`reset:${userId}`);
        
        // Invalidate all sessions
        await cache.delPattern(`refresh:${userId}:*`);
        
        logger.info(`Password reset completed for user: ${userId}`);
        
        return sendResponse(res, null, 200, 'Password reset successful');
    })
};

// ============================================================================
// Health Controller
// ============================================================================

/**
 * Health Controller
 * Handles health checks and monitoring
 */
const healthController = {
    /**
     * Basic health check
     */
    health: (req, res) => {
        return sendResponse(res, {
            status: 'healthy',
            uptime: process.uptime(),
            timestamp: new Date().toISOString(),
            version: process.env.npm_package_version || '1.0.0'
        });
    },

    /**
     * Detailed health check with component status
     */
    healthDetailed: catchAsync(async (req, res) => {
        const checks = {
            server: 'healthy',
            timestamp: new Date().toISOString()
        };
        
        // Check database
        try {
            await User.findOne();
            checks.database = 'healthy';
        } catch (error) {
            checks.database = 'unhealthy';
            checks.databaseError = error.message;
        }
        
        // Check Redis
        try {
            await cache.ping();
            checks.redis = 'healthy';
        } catch (error) {
            checks.redis = 'unhealthy';
            checks.redisError = error.message;
        }
        
        // Check queue
        try {
            const queueStats = await queue.getStats();
            checks.queue = 'healthy';
            checks.queueStats = queueStats;
        } catch (error) {
            checks.queue = 'unhealthy';
            checks.queueError = error.message;
        }
        
        const isHealthy = checks.database === 'healthy' && 
                         checks.redis === 'healthy' && 
                         checks.queue === 'healthy';
        
        return res.status(isHealthy ? 200 : 503).json({
            success: isHealthy,
            checks
        });
    }),

    /**
     * Readiness probe (for Kubernetes)
     */
    readiness: (req, res) => {
        // Check if service is ready to accept traffic
        return res.status(200).json({
            status: 'ready'
        });
    },

    /**
     * Liveness probe (for Kubernetes)
     */
    liveness: (req, res) => {
        // Check if service is alive
        return res.status(200).json({
            status: 'alive'
        });
    }
};

// ============================================================================
// Export all controllers
// ============================================================================

module.exports = {
    userController,
    productController,
    orderController,
    authController,
    healthController
};