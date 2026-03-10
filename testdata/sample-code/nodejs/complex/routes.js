/**
 * Routes Module
 * 
 * Express route definitions for the Node.js API including:
 * - RESTful endpoints
 * - Route grouping and versioning
 * - Middleware integration
 * - Input validation
 * - Documentation
 */

const express = require('express');
const router = express.Router();

// ============================================================================
// Controller Imports
// ============================================================================

const {
    userController,
    productController,
    orderController,
    authController,
    healthController
} = require('../controllers');

const {
    webhookController
} = require('../controllers/webhook.controller');

const {
    analyticsController
} = require('../controllers/analytics.controller');

const {
    adminController
} = require('../controllers/admin.controller');

// ============================================================================
// Middleware Imports
// ============================================================================

const {
    authenticate,
    authorize,
    optionalAuth,
    checkOwnership,
    validateApiKey,
    limiter,
    authLimiter,
    apiKeyLimiter,
    cacheResponse,
    clearCache,
    validateRequest,
    validateQuery,
    validateParams,
    maintenanceMode,
    timeout
} = require('../middleware');

// ============================================================================
// Validation Schemas
// ============================================================================

const {
    userValidation,
    productValidation,
    orderValidation,
    authValidation,
    commonValidation
} = require('../validations');

// ============================================================================
// API Documentation (Swagger)
// ============================================================================

/**
 * @swagger
 * components:
 *   securitySchemes:
 *     bearerAuth:
 *       type: http
 *       scheme: bearer
 *       bearerFormat: JWT
 *     apiKeyAuth:
 *       type: apiKey
 *       in: header
 *       name: X-API-Key
 * 
 *   responses:
 *     UnauthorizedError:
 *       description: Access token is missing or invalid
 *     ForbiddenError:
 *       description: Insufficient permissions
 *     NotFoundError:
 *       description: Resource not found
 *     ValidationError:
 *       description: Request validation failed
 */

// ============================================================================
// Health Check Routes
// ============================================================================

/**
 * @swagger
 * /health:
 *   get:
 *     summary: Basic health check
 *     tags: [Health]
 *     responses:
 *       200:
 *         description: Service is healthy
 */
router.get('/health', healthController.health);

/**
 * @swagger
 * /health/detailed:
 *   get:
 *     summary: Detailed health check with component status
 *     tags: [Health]
 *     responses:
 *       200:
 *         description: Health check passed
 *       503:
 *         description: Service degraded
 */
router.get('/health/detailed', 
    authenticate, 
    authorize('admin'), 
    healthController.healthDetailed
);

/**
 * @swagger
 * /ready:
 *   get:
 *     summary: Readiness probe for Kubernetes
 *     tags: [Health]
 *     responses:
 *       200:
 *         description: Service is ready
 */
router.get('/ready', healthController.readiness);

/**
 * @swagger
 * /live:
 *   get:
 *     summary: Liveness probe for Kubernetes
 *     tags: [Health]
 *     responses:
 *       200:
 *         description: Service is alive
 */
router.get('/live', healthController.liveness);

// ============================================================================
// Metrics Route
// ============================================================================

/**
 * @swagger
 * /metrics:
 *   get:
 *     summary: Prometheus metrics endpoint
 *     tags: [Monitoring]
 *     security:
 *       - apiKeyAuth: []
 *     responses:
 *       200:
 *         description: Metrics in Prometheus format
 */
router.get('/metrics', 
    validateApiKey,
    (req, res) => {
        res.set('Content-Type', 'text/plain');
        res.send(require('../utils/metrics').getMetrics());
    }
);

// ============================================================================
// Auth Routes
// ============================================================================

const authRouter = express.Router();

/**
 * @swagger
 * /auth/register:
 *   post:
 *     summary: Register a new user
 *     tags: [Authentication]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/RegisterRequest'
 *     responses:
 *       201:
 *         description: User registered successfully
 *       400:
 *         description: Validation error
 *       409:
 *         description: Email already exists
 */
authRouter.post('/register',
    authLimiter,
    validateRequest(authValidation.register),
    authController.register
);

/**
 * @swagger
 * /auth/login:
 *   post:
 *     summary: Login user
 *     tags: [Authentication]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/LoginRequest'
 *     responses:
 *       200:
 *         description: Login successful
 *       401:
 *         description: Invalid credentials
 */
authRouter.post('/login',
    authLimiter,
    validateRequest(authValidation.login),
    authController.login
);

/**
 * @swagger
 * /auth/refresh:
 *   post:
 *     summary: Refresh access token
 *     tags: [Authentication]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               refreshToken:
 *                 type: string
 *     responses:
 *       200:
 *         description: New tokens generated
 *       401:
 *         description: Invalid refresh token
 */
authRouter.post('/refresh',
    validateRequest(authValidation.refreshToken),
    authController.refreshToken
);

/**
 * @swagger
 * /auth/logout:
 *   post:
 *     summary: Logout user
 *     tags: [Authentication]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: Logout successful
 */
authRouter.post('/logout',
    authenticate,
    authController.logout
);

/**
 * @swagger
 * /auth/forgot-password:
 *   post:
 *     summary: Request password reset
 *     tags: [Authentication]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               email:
 *                 type: string
 *                 format: email
 *     responses:
 *       200:
 *         description: Reset email sent if email exists
 */
authRouter.post('/forgot-password',
    authLimiter,
    validateRequest(authValidation.forgotPassword),
    authController.forgotPassword
);

/**
 * @swagger
 * /auth/reset-password:
 *   post:
 *     summary: Reset password with token
 *     tags: [Authentication]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               userId:
 *                 type: string
 *               token:
 *                 type: string
 *               newPassword:
 *                 type: string
 *     responses:
 *       200:
 *         description: Password reset successful
 *       400:
 *         description: Invalid or expired token
 */
authRouter.post('/reset-password',
    authLimiter,
    validateRequest(authValidation.resetPassword),
    authController.resetPassword
);

/**
 * @swagger
 * /auth/verify-email/{token}:
 *   get:
 *     summary: Verify email address
 *     tags: [Authentication]
 *     parameters:
 *       - in: path
 *         name: token
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Email verified successfully
 *       400:
 *         description: Invalid token
 */
authRouter.get('/verify-email/:token',
    validateParams(commonValidation.token),
    authController.verifyEmail
);

/**
 * @swagger
 * /auth/oauth/{provider}:
 *   get:
 *     summary: OAuth authentication redirect
 *     tags: [Authentication]
 *     parameters:
 *       - in: path
 *         name: provider
 *         required: true
 *         schema:
 *           type: string
 *           enum: [google, github, facebook]
 *     responses:
 *       302:
 *         description: Redirect to OAuth provider
 */
authRouter.get('/oauth/:provider',
    validateParams(commonValidation.provider),
    authController.oauthRedirect
);

/**
 * @swagger
 * /auth/oauth/{provider}/callback:
 *   get:
 *     summary: OAuth callback
 *     tags: [Authentication]
 *     parameters:
 *       - in: path
 *         name: provider
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       302:
 *         description: Redirect to frontend with tokens
 */
authRouter.get('/oauth/:provider/callback',
    validateParams(commonValidation.provider),
    authController.oauthCallback
);

// ============================================================================
// User Routes
// ============================================================================

const userRouter = express.Router();

// Apply authentication to all user routes
userRouter.use(authenticate);

/**
 * @swagger
 * /users:
 *   get:
 *     summary: Get all users (paginated)
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           default: 1
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *           default: 10
 *       - in: query
 *         name: role
 *         schema:
 *           type: string
 *           enum: [user, admin, manager]
 *       - in: query
 *         name: search
 *         schema:
 *           type: string
 *       - in: query
 *         name: sortBy
 *         schema:
 *           type: string
 *           default: createdAt
 *       - in: query
 *         name: sortOrder
 *         schema:
 *           type: string
 *           enum: [asc, desc]
 *           default: desc
 *     responses:
 *       200:
 *         description: List of users
 */
userRouter.get('/',
    authorize('admin'),
    validateQuery(userValidation.list),
    cacheResponse(300), // Cache for 5 minutes
    userController.getAllUsers
);

/**
 * @swagger
 * /users/{id}:
 *   get:
 *     summary: Get user by ID
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: User details
 *       404:
 *         description: User not found
 */
userRouter.get('/:id',
    validateParams(commonValidation.id),
    userController.getUserById
);

/**
 * @swagger
 * /users/{id}:
 *   put:
 *     summary: Update user
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/UpdateUserRequest'
 *     responses:
 *       200:
 *         description: User updated successfully
 *       403:
 *         description: Not authorized to update this user
 */
userRouter.put('/:id',
    validateParams(commonValidation.id),
    validateRequest(userValidation.update),
    checkOwnership('userId', 'User'),
    userController.updateUser
);

/**
 * @swagger
 * /users/{id}:
 *   delete:
 *     summary: Delete user
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       204:
 *         description: User deleted successfully
 *       403:
 *         description: Not authorized to delete this user
 */
userRouter.delete('/:id',
    validateParams(commonValidation.id),
    authorize('admin'),
    clearCache(['users:*']),
    userController.deleteUser
);

/**
 * @swagger
 * /users/profile:
 *   get:
 *     summary: Get current user profile
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: Current user profile
 */
userRouter.get('/profile/me',
    userController.getUserById
);

/**
 * @swagger
 * /users/profile:
 *   put:
 *     summary: Update current user profile
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/UpdateProfileRequest'
 *     responses:
 *       200:
 *         description: Profile updated successfully
 */
userRouter.put('/profile/me',
    validateRequest(userValidation.updateProfile),
    userController.updateProfile
);

/**
 * @swagger
 * /users/change-password:
 *   post:
 *     summary: Change user password
 *     tags: [Users]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               currentPassword:
 *                 type: string
 *               newPassword:
 *                 type: string
 *     responses:
 *       200:
 *         description: Password changed successfully
 */
userRouter.post('/change-password',
    validateRequest(userValidation.changePassword),
    userController.changePassword
);

// ============================================================================
// Product Routes
// ============================================================================

const productRouter = express.Router();

/**
 * @swagger
 * /products:
 *   get:
 *     summary: Get all products
 *     tags: [Products]
 *     parameters:
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           default: 1
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *           default: 20
 *       - in: query
 *         name: category
 *         schema:
 *           type: string
 *       - in: query
 *         name: minPrice
 *         schema:
 *           type: number
 *       - in: query
 *         name: maxPrice
 *         schema:
 *           type: number
 *       - in: query
 *         name: inStock
 *         schema:
 *           type: boolean
 *       - in: query
 *         name: search
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: List of products
 */
productRouter.get('/',
    optionalAuth,
    validateQuery(productValidation.list),
    cacheResponse(60), // Cache for 1 minute
    productController.getAllProducts
);

/**
 * @swagger
 * /products/{id}:
 *   get:
 *     summary: Get product by ID
 *     tags: [Products]
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Product details
 *       404:
 *         description: Product not found
 */
productRouter.get('/:id',
    optionalAuth,
    validateParams(commonValidation.id),
    cacheResponse(60),
    productController.getProductById
);

/**
 * @swagger
 * /products/{slug}:
 *   get:
 *     summary: Get product by slug
 *     tags: [Products]
 *     parameters:
 *       - in: path
 *         name: slug
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Product details
 */
productRouter.get('/slug/:slug',
    optionalAuth,
    validateParams(commonValidation.slug),
    productController.getProductBySlug
);

/**
 * @swagger
 * /products:
 *   post:
 *     summary: Create new product
 *     tags: [Products]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/CreateProductRequest'
 *     responses:
 *       201:
 *         description: Product created successfully
 */
productRouter.post('/',
    authenticate,
    authorize('admin', 'manager'),
    validateRequest(productValidation.create),
    clearCache(['products:*']),
    productController.createProduct
);

/**
 * @swagger
 * /products/{id}:
 *   put:
 *     summary: Update product
 *     tags: [Products]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/UpdateProductRequest'
 *     responses:
 *       200:
 *         description: Product updated successfully
 */
productRouter.put('/:id',
    authenticate,
    authorize('admin', 'manager'),
    validateParams(commonValidation.id),
    validateRequest(productValidation.update),
    clearCache(['products:*', `product:${req.params.id}`]),
    productController.updateProduct
);

/**
 * @swagger
 * /products/{id}:
 *   delete:
 *     summary: Delete product
 *     tags: [Products]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       204:
 *         description: Product deleted successfully
 */
productRouter.delete('/:id',
    authenticate,
    authorize('admin'),
    validateParams(commonValidation.id),
    clearCache(['products:*', `product:${req.params.id}`]),
    productController.deleteProduct
);

/**
 * @swagger
 * /products/{id}/stock:
 *   patch:
 *     summary: Update product stock
 *     tags: [Products]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               quantity:
 *                 type: integer
 *               operation:
 *                 type: string
 *                 enum: [set, increment, decrement]
 *     responses:
 *       200:
 *         description: Stock updated
 */
productRouter.patch('/:id/stock',
    authenticate,
    authorize('admin', 'manager'),
    validateParams(commonValidation.id),
    validateRequest(productValidation.updateStock),
    clearCache([`product:${req.params.id}`]),
    productController.updateStock
);

// ============================================================================
// Order Routes
// ============================================================================

const orderRouter = express.Router();

// Apply authentication to all order routes
orderRouter.use(authenticate);

/**
 * @swagger
 * /orders:
 *   get:
 *     summary: Get user orders
 *     tags: [Orders]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           default: 1
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *           default: 10
 *       - in: query
 *         name: status
 *         schema:
 *           type: string
 *           enum: [pending, processing, shipped, delivered, cancelled]
 *     responses:
 *       200:
 *         description: List of user orders
 */
orderRouter.get('/',
    validateQuery(orderValidation.list),
    orderController.getUserOrders
);

/**
 * @swagger
 * /orders/all:
 *   get:
 *     summary: Get all orders (admin only)
 *     tags: [Orders]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           default: 1
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *           default: 20
 *     responses:
 *       200:
 *         description: List of all orders
 */
orderRouter.get('/all',
    authorize('admin'),
    orderController.getAllOrders
);

/**
 * @swagger
 * /orders/{id}:
 *   get:
 *     summary: Get order by ID
 *     tags: [Orders]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Order details
 *       404:
 *         description: Order not found
 */
orderRouter.get('/:id',
    validateParams(commonValidation.id),
    orderController.getOrderById
);

/**
 * @swagger
 * /orders:
 *   post:
 *     summary: Create new order
 *     tags: [Orders]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/CreateOrderRequest'
 *     responses:
 *       201:
 *         description: Order created successfully
 */
orderRouter.post('/',
    validateRequest(orderValidation.create),
    clearCache(['orders:*']),
    orderController.createOrder
);

/**
 * @swagger
 * /orders/{id}/status:
 *   patch:
 *     summary: Update order status
 *     tags: [Orders]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               status:
 *                 type: string
 *                 enum: [processing, shipped, delivered, cancelled]
 *               trackingNumber:
 *                 type: string
 *     responses:
 *       200:
 *         description: Order status updated
 */
orderRouter.patch('/:id/status',
    authenticate,
    authorize('admin'),
    validateParams(commonValidation.id),
    validateRequest(orderValidation.updateStatus),
    clearCache([`order:${req.params.id}`, 'orders:*']),
    orderController.updateOrderStatus
);

/**
 * @swagger
 * /orders/{id}/cancel:
 *   post:
 *     summary: Cancel order
 *     tags: [Orders]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Order cancelled successfully
 */
orderRouter.post('/:id/cancel',
    validateParams(commonValidation.id),
    clearCache([`order:${req.params.id}`, 'orders:*']),
    orderController.cancelOrder
);

/**
 * @swagger
 * /orders/{id}/track:
 *   get:
 *     summary: Track order shipment
 *     tags: [Orders]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Tracking information
 */
orderRouter.get('/:id/track',
    validateParams(commonValidation.id),
    orderController.trackOrder
);

// ============================================================================
// Admin Routes
// ============================================================================

const adminRouter = express.Router();

adminRouter.use(authenticate);
adminRouter.use(authorize('admin'));

/**
 * @swagger
 * /admin/stats:
 *   get:
 *     summary: Get system statistics
 *     tags: [Admin]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: System statistics
 */
adminRouter.get('/stats',
    adminController.getStats
);

/**
 * @swagger
 * /admin/users:
 *   get:
 *     summary: Get all users with details
 *     tags: [Admin]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: List of all users
 */
adminRouter.get('/users',
    adminController.getAllUsers
);

/**
 * @swagger
 * /admin/users/{id}/suspend:
 *   post:
 *     summary: Suspend user account
 *     tags: [Admin]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: User suspended
 */
adminRouter.post('/users/:id/suspend',
    validateParams(commonValidation.id),
    adminController.suspendUser
);

/**
 * @swagger
 * /admin/reports/sales:
 *   get:
 *     summary: Get sales report
 *     tags: [Admin]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: startDate
 *         required: true
 *         schema:
 *           type: string
 *           format: date
 *       - in: query
 *         name: endDate
 *         required: true
 *         schema:
 *           type: string
 *           format: date
 *     responses:
 *       200:
 *         description: Sales report data
 */
adminRouter.get('/reports/sales',
    validateQuery(adminValidation.salesReport),
    adminController.getSalesReport
);

// ============================================================================
// Webhook Routes
// ============================================================================

const webhookRouter = express.Router();

// No authentication for webhooks (they use signature verification)

/**
 * @swagger
 * /webhooks/stripe:
 *   post:
 *     summary: Stripe webhook handler
 *     tags: [Webhooks]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *     responses:
 *       200:
 *         description: Webhook received
 */
webhookRouter.post('/stripe',
    express.raw({ type: 'application/json' }),
    timeout(10000),
    webhookController.stripeWebhook
);

/**
 * @swagger
 * /webhooks/paypal:
 *   post:
 *     summary: PayPal webhook handler
 *     tags: [Webhooks]
 *     responses:
 *       200:
 *         description: Webhook received
 */
webhookRouter.post('/paypal',
    webhookController.paypalWebhook
);

/**
 * @swagger
 * /webhooks/sendgrid:
 *   post:
 *     summary: SendGrid webhook for email events
 *     tags: [Webhooks]
 *     responses:
 *       200:
 *         description: Webhook received
 */
webhookRouter.post('/sendgrid',
    webhookController.sendgridWebhook
);

// ============================================================================
// Analytics Routes
// ============================================================================

const analyticsRouter = express.Router();

analyticsRouter.use(authenticate);

/**
 * @swagger
 * /analytics/dashboard:
 *   get:
 *     summary: Get analytics dashboard data
 *     tags: [Analytics]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: Dashboard analytics
 */
analyticsRouter.get('/dashboard',
    authorize('admin', 'manager'),
    analyticsController.getDashboardData
);

/**
 * @swagger
 * /analytics/events:
 *   post:
 *     summary: Track user event
 *     tags: [Analytics]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               event:
 *                 type: string
 *               data:
 *                 type: object
 *     responses:
 *       202:
 *         description: Event accepted
 */
analyticsRouter.post('/events',
    analyticsController.trackEvent
);

// ============================================================================
// Search Routes
// ============================================================================

const searchRouter = express.Router();

/**
 * @swagger
 * /search:
 *   get:
 *     summary: Global search across resources
 *     tags: [Search]
 *     parameters:
 *       - in: query
 *         name: q
 *         required: true
 *         schema:
 *           type: string
 *       - in: query
 *         name: type
 *         schema:
 *           type: string
 *           enum: [products, users, orders]
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           default: 1
 *     responses:
 *       200:
 *         description: Search results
 */
searchRouter.get('/',
    optionalAuth,
    validateQuery(searchValidation.global),
    searchController.globalSearch
);

// ============================================================================
// Static Routes
// ============================================================================

const staticRouter = express.Router();

/**
 * Serve static files
 */
staticRouter.use('/uploads', express.static('uploads'));
staticRouter.use('/images', express.static('public/images'));
staticRouter.use('/docs', express.static('docs'));

// ============================================================================
// Mount all routers
// ============================================================================

// API version 1
const v1Router = express.Router();

v1Router.use('/auth', authRouter);
v1Router.use('/users', userRouter);
v1Router.use('/products', productRouter);
v1Router.use('/orders', orderRouter);
v1Router.use('/admin', adminRouter);
v1Router.use('/webhooks', webhookRouter);
v1Router.use('/analytics', analyticsRouter);
v1Router.use('/search', searchRouter);

// API version 2 (future)
const v2Router = express.Router();
// v2Router.use('/users', userRouterV2);

// Mount API versions
router.use('/api/v1', v1Router);
router.use('/api/v2', v2Router);

// Mount static routes
router.use('/', staticRouter);

// ============================================================================
// GraphQL endpoint (if using GraphQL)
// ============================================================================

/**
 * @swagger
 * /graphql:
 *   post:
 *     summary: GraphQL endpoint
 *     tags: [GraphQL]
 *     responses:
 *       200:
 *         description: GraphQL response
 */
if (process.env.ENABLE_GRAPHQL === 'true') {
    const { graphqlHTTP } = require('express-graphql');
    const { schema } = require('../graphql/schema');
    
    router.use('/graphql',
        authenticate,
        graphqlHTTP((req) => ({
            schema,
            graphiql: process.env.NODE_ENV === 'development',
            context: { user: req.user }
        }))
    );
}

// ============================================================================
// WebSocket routes (Socket.io)
// ============================================================================

/**
 * Socket.io namespace for real-time updates
 */
const socketRoutes = (io) => {
    const notifications = io.of('/notifications');
    
    notifications.use((socket, next) => {
        const token = socket.handshake.auth.token;
        // Verify token middleware
        next();
    });
    
    notifications.on('connection', (socket) => {
        console.log('Notification client connected');
        
        socket.on('subscribe', (userId) => {
            socket.join(`user:${userId}`);
        });
        
        socket.on('disconnect', () => {
            console.log('Notification client disconnected');
        });
    });
    
    const chat = io.of('/chat');
    
    chat.on('connection', (socket) => {
        console.log('Chat client connected');
        
        socket.on('join-room', (roomId) => {
            socket.join(roomId);
        });
        
        socket.on('send-message', (data) => {
            io.of('/chat').to(data.roomId).emit('new-message', data);
        });
    });
};

// ============================================================================
// Export router and socket handler
// ============================================================================

module.exports = {
    router,
    socketRoutes
};