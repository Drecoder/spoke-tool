/**
 * Models Module
 * 
 * Mongoose models for the Node.js API including:
 * - User management with roles and authentication
 * - Product catalog with categories and variants
 * - Order processing with items and status tracking
 * - Reviews and ratings
 * - Inventory management
 */

const mongoose = require('mongoose');
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const crypto = require('crypto');
const validator = require('validator');

// ============================================================================
// Schema Plugins
// ============================================================================

/**
 * Soft delete plugin
 */
const softDelete = (schema) => {
    schema.add({
        deletedAt: Date,
        isDeleted: { type: Boolean, default: false }
    });

    schema.pre(/^find/, function() {
        if (!this.getQuery().includeDeleted) {
            this.where({ isDeleted: false });
        }
    });

    schema.methods.softDelete = async function() {
        this.isDeleted = true;
        this.deletedAt = new Date();
        await this.save();
    };

    schema.methods.restore = async function() {
        this.isDeleted = false;
        this.deletedAt = undefined;
        await this.save();
    };
};

/**
 * Timestamp plugin
 */
const timestamp = (schema) => {
    schema.add({
        createdAt: { type: Date, default: Date.now },
        updatedAt: { type: Date, default: Date.now }
    });

    schema.pre('save', function(next) {
        this.updatedAt = Date.now();
        next();
    });

    schema.pre('updateOne', function(next) {
        this.set({ updatedAt: Date.now() });
        next();
    });
};

/**
 * Version plugin (optimistic locking)
 */
const versionLock = (schema) => {
    schema.add({
        __v: { type: Number, select: false }
    });

    schema.pre('save', function(next) {
        this.increment();
        next();
    });
};

// ============================================================================
// User Model
// ============================================================================

const userSchema = new mongoose.Schema({
    email: {
        type: String,
        required: [true, 'Email is required'],
        unique: true,
        lowercase: true,
        trim: true,
        validate: {
            validator: validator.isEmail,
            message: 'Please provide a valid email'
        }
    },
    password: {
        type: String,
        required: [true, 'Password is required'],
        minlength: [8, 'Password must be at least 8 characters'],
        select: false
    },
    name: {
        type: String,
        required: [true, 'Name is required'],
        trim: true,
        maxlength: [100, 'Name cannot exceed 100 characters']
    },
    role: {
        type: String,
        enum: ['user', 'admin', 'manager', 'moderator'],
        default: 'user'
    },
    avatar: {
        type: String,
        default: 'default-avatar.png'
    },
    bio: {
        type: String,
        maxlength: [500, 'Bio cannot exceed 500 characters']
    },
    phone: {
        type: String,
        validate: {
            validator: function(v) {
                return !v || validator.isMobilePhone(v, 'any');
            },
            message: 'Please provide a valid phone number'
        }
    },
    address: {
        street: String,
        city: String,
        state: String,
        zipCode: String,
        country: String
    },
    preferences: {
        newsletter: { type: Boolean, default: true },
        notifications: { type: Boolean, default: true },
        language: { type: String, default: 'en', enum: ['en', 'es', 'fr', 'de'] },
        theme: { type: String, default: 'light', enum: ['light', 'dark'] }
    },
    social: {
        github: String,
        twitter: String,
        linkedin: String
    },
    lastLogin: Date,
    passwordChangedAt: Date,
    passwordResetToken: String,
    passwordResetExpires: Date,
    emailVerified: { type: Boolean, default: false },
    emailVerificationToken: String,
    twoFactorEnabled: { type: Boolean, default: false },
    twoFactorSecret: { type: String, select: false },
    loginAttempts: { type: Number, default: 0 },
    lockUntil: Date,
    active: { type: Boolean, default: true }
}, {
    toJSON: { virtuals: true },
    toObject: { virtuals: true }
});

// Apply plugins
userSchema.plugin(softDelete);
userSchema.plugin(timestamp);
userSchema.plugin(versionLock);

// Virtual for full name (if first/last name were separate)
userSchema.virtual('initials').get(function() {
    return this.name
        .split(' ')
        .map(word => word[0])
        .join('')
        .toUpperCase()
        .slice(0, 2);
});

// Virtual for user's orders
userSchema.virtual('orders', {
    ref: 'Order',
    localField: '_id',
    foreignField: 'userId'
});

// Virtual for user's reviews
userSchema.virtual('reviews', {
    ref: 'Review',
    localField: '_id',
    foreignField: 'userId'
});

// Indexes
userSchema.index({ email: 1 });
userSchema.index({ role: 1 });
userSchema.index({ 'address.country': 1, 'address.city': 1 });

// Pre-save middleware - hash password
userSchema.pre('save', async function(next) {
    if (!this.isModified('password')) return next();
    
    try {
        const salt = await bcrypt.genSalt(12);
        this.password = await bcrypt.hash(this.password, salt);
        
        // Update passwordChangedAt if not new user
        if (!this.isNew) {
            this.passwordChangedAt = Date.now() - 1000; // Subtract 1s to ensure token is created after
        }
        
        next();
    } catch (error) {
        next(error);
    }
});

// Instance methods
userSchema.methods.comparePassword = async function(candidatePassword) {
    return await bcrypt.compare(candidatePassword, this.password);
};

userSchema.methods.generateAuthToken = function() {
    return jwt.sign(
        { 
            id: this._id,
            email: this.email,
            role: this.role 
        },
        process.env.JWT_SECRET,
        { expiresIn: process.env.JWT_EXPIRES_IN }
    );
};

userSchema.methods.generateRefreshToken = function() {
    return jwt.sign(
        { id: this._id },
        process.env.REFRESH_TOKEN_SECRET,
        { expiresIn: process.env.REFRESH_TOKEN_EXPIRES_IN }
    );
};

userSchema.methods.createPasswordResetToken = function() {
    const resetToken = crypto.randomBytes(32).toString('hex');
    
    this.passwordResetToken = crypto
        .createHash('sha256')
        .update(resetToken)
        .digest('hex');
    
    this.passwordResetExpires = Date.now() + 10 * 60 * 1000; // 10 minutes
    
    return resetToken;
};

userSchema.methods.createEmailVerificationToken = function() {
    const verificationToken = crypto.randomBytes(32).toString('hex');
    
    this.emailVerificationToken = crypto
        .createHash('sha256')
        .update(verificationToken)
        .digest('hex');
    
    return verificationToken;
};

userSchema.methods.incrementLoginAttempts = function() {
    // Reset attempts if lock has expired
    if (this.lockUntil && this.lockUntil < Date.now()) {
        return this.updateOne({
            $set: { loginAttempts: 1 },
            $unset: { lockUntil: 1 }
        });
    }
    
    // Increment attempts
    const updates = { $inc: { loginAttempts: 1 } };
    
    // Lock account after 5 failed attempts
    if (this.loginAttempts + 1 >= 5 && !this.isLocked) {
        updates.$set = { lockUntil: Date.now() + 2 * 60 * 60 * 1000 }; // 2 hours
    }
    
    return this.updateOne(updates);
};

userSchema.methods.isLocked = function() {
    return !!(this.lockUntil && this.lockUntil > Date.now());
};

// Static methods
userSchema.statics.findByEmail = function(email) {
    return this.findOne({ email: email.toLowerCase() });
};

userSchema.statics.findByCredentials = async function(email, password) {
    const user = await this.findOne({ email: email.toLowerCase() }).select('+password');
    
    if (!user) {
        return null;
    }
    
    const isMatch = await user.comparePassword(password);
    return isMatch ? user : null;
};

// Create User model
const User = mongoose.model('User', userSchema);

// ============================================================================
// Category Model
// ============================================================================

const categorySchema = new mongoose.Schema({
    name: {
        type: String,
        required: [true, 'Category name is required'],
        unique: true,
        trim: true,
        maxlength: [100, 'Category name cannot exceed 100 characters']
    },
    slug: {
        type: String,
        required: true,
        unique: true,
        lowercase: true
    },
    description: {
        type: String,
        maxlength: [500, 'Description cannot exceed 500 characters']
    },
    parent: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'Category'
    },
    image: String,
    icon: String,
    isActive: {
        type: Boolean,
        default: true
    },
    sortOrder: {
        type: Number,
        default: 0
    },
    meta: {
        title: String,
        description: String,
        keywords: [String]
    }
}, {
    toJSON: { virtuals: true },
    toObject: { virtuals: true }
});

categorySchema.plugin(timestamp);

// Virtual for subcategories
categorySchema.virtual('subcategories', {
    ref: 'Category',
    localField: '_id',
    foreignField: 'parent'
});

// Virtual for products in this category
categorySchema.virtual('products', {
    ref: 'Product',
    localField: '_id',
    foreignField: 'category'
});

// Pre-save middleware to generate slug
categorySchema.pre('save', function(next) {
    if (!this.isModified('name')) return next();
    
    this.slug = this.name
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, '-')
        .replace(/^-|-$/g, '');
    
    next();
});

const Category = mongoose.model('Category', categorySchema);

// ============================================================================
// Product Model
// ============================================================================

const productSchema = new mongoose.Schema({
    name: {
        type: String,
        required: [true, 'Product name is required'],
        trim: true,
        maxlength: [200, 'Product name cannot exceed 200 characters']
    },
    slug: {
        type: String,
        required: true,
        unique: true,
        lowercase: true
    },
    sku: {
        type: String,
        required: [true, 'SKU is required'],
        unique: true,
        uppercase: true
    },
    description: {
        type: String,
        required: [true, 'Description is required']
    },
    shortDescription: {
        type: String,
        maxlength: [300, 'Short description cannot exceed 300 characters']
    },
    price: {
        type: Number,
        required: [true, 'Price is required'],
        min: [0, 'Price cannot be negative']
    },
    compareAtPrice: {
        type: Number,
        min: [0, 'Compare at price cannot be negative']
    },
    cost: {
        type: Number,
        min: [0, 'Cost cannot be negative']
    },
    category: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'Category',
        required: [true, 'Category is required']
    },
    brand: String,
    tags: [String],
    attributes: {
        type: Map,
        of: mongoose.Schema.Types.Mixed
    },
    variants: [{
        sku: { type: String, required: true },
        name: String,
        price: Number,
        compareAtPrice: Number,
        options: Map,
        stock: { type: Number, default: 0 },
        images: [String]
    }],
    images: [String],
    featuredImage: String,
    stock: {
        type: Number,
        required: true,
        default: 0,
        min: 0
    },
    reservedStock: {
        type: Number,
        default: 0,
        min: 0
    },
    availableStock: {
        type: Number,
        virtual: true,
        get: function() {
            return this.stock - this.reservedStock;
        }
    },
    weight: {
        value: Number,
        unit: { type: String, enum: ['g', 'kg', 'lb', 'oz'], default: 'g' }
    },
    dimensions: {
        length: Number,
        width: Number,
        height: Number,
        unit: { type: String, enum: ['cm', 'in'], default: 'cm' }
    },
    seo: {
        title: String,
        description: String,
        keywords: [String]
    },
    isActive: {
        type: Boolean,
        default: true
    },
    isFeatured: {
        type: Boolean,
        default: false
    },
    isDigital: {
        type: Boolean,
        default: false
    },
    downloadUrl: String,
    requiresShipping: {
        type: Boolean,
        default: true
    },
    ratings: {
        average: { type: Number, default: 0 },
        count: { type: Number, default: 0 }
    },
    soldCount: {
        type: Number,
        default: 0
    },
    viewCount: {
        type: Number,
        default: 0
    }
}, {
    toJSON: { virtuals: true },
    toObject: { virtuals: true }
});

productSchema.plugin(timestamp);
productSchema.plugin(softDelete);

// Indexes for search and filtering
productSchema.index({ name: 'text', description: 'text', tags: 'text' });
productSchema.index({ category: 1, isActive: 1 });
productSchema.index({ price: 1 });
productSchema.index({ brand: 1 });
productSchema.index({ 'ratings.average': -1 });
productSchema.index({ createdAt: -1 });

// Virtual for reviews
productSchema.virtual('reviews', {
    ref: 'Review',
    localField: '_id',
    foreignField: 'productId'
});

// Pre-save middleware for slug generation
productSchema.pre('save', function(next) {
    if (this.isModified('name')) {
        this.slug = this.name
            .toLowerCase()
            .replace(/[^a-z0-9]+/g, '-')
            .replace(/^-|-$/g, '');
    }
    next();
});

// Instance methods
productSchema.methods.reduceStock = async function(quantity) {
    if (this.stock < quantity) {
        throw new Error('Insufficient stock');
    }
    this.stock -= quantity;
    this.soldCount += quantity;
    await this.save();
};

productSchema.methods.reserveStock = async function(quantity) {
    if (this.stock - this.reservedStock < quantity) {
        throw new Error('Insufficient available stock');
    }
    this.reservedStock += quantity;
    await this.save();
};

productSchema.methods.releaseStock = async function(quantity) {
    this.reservedStock = Math.max(0, this.reservedStock - quantity);
    await this.save();
};

productSchema.methods.updateRating = async function() {
    const reviews = await Review.find({ productId: this._id, isApproved: true });
    
    if (reviews.length === 0) {
        this.ratings = { average: 0, count: 0 };
    } else {
        const sum = reviews.reduce((acc, review) => acc + review.rating, 0);
        this.ratings = {
            average: sum / reviews.length,
            count: reviews.length
        };
    }
    
    await this.save();
};

// Static methods
productSchema.statics.findBySku = function(sku) {
    return this.findOne({ sku: sku.toUpperCase() });
};

productSchema.statics.findFeatured = function(limit = 10) {
    return this.find({ isFeatured: true, isActive: true })
        .limit(limit)
        .sort('-createdAt');
};

const Product = mongoose.model('Product', productSchema);

// ============================================================================
// Order Model
// ============================================================================

const orderItemSchema = new mongoose.Schema({
    productId: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'Product',
        required: true
    },
    name: {
        type: String,
        required: true
    },
    sku: {
        type: String,
        required: true
    },
    price: {
        type: Number,
        required: true,
        min: 0
    },
    quantity: {
        type: Number,
        required: true,
        min: 1
    },
    total: {
        type: Number,
        required: true,
        min: 0
    },
    variant: {
        sku: String,
        name: String,
        options: Map
    }
});

const orderSchema = new mongoose.Schema({
    orderNumber: {
        type: String,
        required: true,
        unique: true
    },
    userId: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'User',
        required: true
    },
    email: {
        type: String,
        required: true,
        lowercase: true
    },
    items: [orderItemSchema],
    subtotal: {
        type: Number,
        required: true,
        min: 0
    },
    tax: {
        type: Number,
        default: 0,
        min: 0
    },
    shippingCost: {
        type: Number,
        default: 0,
        min: 0
    },
    discount: {
        type: Number,
        default: 0,
        min: 0
    },
    total: {
        type: Number,
        required: true,
        min: 0
    },
    status: {
        type: String,
        enum: [
            'pending',
            'processing',
            'confirmed',
            'shipped',
            'delivered',
            'cancelled',
            'refunded',
            'failed'
        ],
        default: 'pending'
    },
    paymentStatus: {
        type: String,
        enum: ['pending', 'paid', 'failed', 'refunded'],
        default: 'pending'
    },
    paymentMethod: {
        type: String,
        enum: ['credit_card', 'paypal', 'bank_transfer', 'cash'],
        required: true
    },
    paymentId: String,
    paymentDetails: mongoose.Schema.Types.Mixed,
    shippingAddress: {
        name: String,
        addressLine1: { type: String, required: true },
        addressLine2: String,
        city: { type: String, required: true },
        state: String,
        zipCode: { type: String, required: true },
        country: { type: String, required: true },
        phone: String
    },
    billingAddress: {
        name: String,
        addressLine1: String,
        addressLine2: String,
        city: String,
        state: String,
        zipCode: String,
        country: String
    },
    shippingMethod: {
        carrier: String,
        service: String,
        trackingNumber: String,
        trackingUrl: String,
        estimatedDelivery: Date
    },
    notes: String,
    adminNotes: String,
    ipAddress: String,
    userAgent: String,
    couponCode: String,
    couponDiscount: Number,
    giftMessage: String,
    isGift: {
        type: Boolean,
        default: false
    }
}, {
    timestamps: true,
    toJSON: { virtuals: true },
    toObject: { virtuals: true }
});

// Indexes
orderSchema.index({ orderNumber: 1 });
orderSchema.index({ userId: 1, createdAt: -1 });
orderSchema.index({ status: 1 });
orderSchema.index({ 'shippingAddress.country': 1 });
orderSchema.index({ createdAt: 1 });

// Pre-save middleware to generate order number
orderSchema.pre('save', async function(next) {
    if (!this.orderNumber) {
        const date = new Date();
        const year = date.getFullYear().toString().slice(-2);
        const month = (date.getMonth() + 1).toString().padStart(2, '0');
        const day = date.getDate().toString().padStart(2, '0');
        
        // Get count of orders today for sequential number
        const count = await Order.countDocuments({
            createdAt: {
                $gte: new Date(date.setHours(0, 0, 0, 0)),
                $lt: new Date(date.setHours(23, 59, 59, 999))
            }
        });
        
        this.orderNumber = `ORD-${year}${month}${day}-${(count + 1).toString().padStart(4, '0')}`;
    }
    next();
});

// Pre-save middleware to calculate totals
orderSchema.pre('save', function(next) {
    // Recalculate item totals
    this.items.forEach(item => {
        item.total = item.price * item.quantity;
    });
    
    // Recalculate subtotal
    this.subtotal = this.items.reduce((acc, item) => acc + item.total, 0);
    
    // Recalculate total
    this.total = this.subtotal + this.tax + this.shippingCost - this.discount;
    
    next();
});

// Instance methods
orderSchema.methods.canBeCancelled = function() {
    return ['pending', 'processing'].includes(this.status);
};

orderSchema.methods.canBeRefunded = function() {
    return ['paid', 'shipped', 'delivered'].includes(this.paymentStatus);
};

orderSchema.methods.markAsPaid = async function(paymentId, paymentDetails = {}) {
    this.paymentStatus = 'paid';
    this.paymentId = paymentId;
    this.paymentDetails = paymentDetails;
    this.status = 'processing';
    await this.save();
};

orderSchema.methods.markAsShipped = async function(trackingNumber, carrier) {
    this.status = 'shipped';
    this.shippingMethod = {
        ...this.shippingMethod,
        trackingNumber,
        carrier,
        trackingUrl: `https://${carrier}.com/track/${trackingNumber}`
    };
    await this.save();
};

orderSchema.methods.markAsDelivered = async function() {
    this.status = 'delivered';
    await this.save();
};

orderSchema.methods.cancel = async function(reason) {
    if (!this.canBeCancelled()) {
        throw new Error('Order cannot be cancelled');
    }
    
    this.status = 'cancelled';
    this.adminNotes = reason;
    await this.save();
    
    // Restore stock for cancelled order
    for (const item of this.items) {
        await Product.findByIdAndUpdate(item.productId, {
            $inc: { stock: item.quantity }
        });
    }
};

// Static methods
orderSchema.statics.findByUser = function(userId) {
    return this.find({ userId }).sort('-createdAt');
};

orderSchema.statics.findByDateRange = function(startDate, endDate) {
    return this.find({
        createdAt: {
            $gte: startDate,
            $lte: endDate
        }
    });
};

orderSchema.statics.getSalesReport = async function(startDate, endDate) {
    const pipeline = [
        {
            $match: {
                createdAt: { $gte: startDate, $lte: endDate },
                paymentStatus: 'paid'
            }
        },
        {
            $group: {
                _id: { $dateToString: { format: '%Y-%m-%d', date: '$createdAt' } },
                totalSales: { $sum: '$total' },
                orderCount: { $sum: 1 },
                averageOrderValue: { $avg: '$total' }
            }
        },
        { $sort: { _id: 1 } }
    ];
    
    return this.aggregate(pipeline);
};

const Order = mongoose.model('Order', orderSchema);

// ============================================================================
// Review Model
// ============================================================================

const reviewSchema = new mongoose.Schema({
    productId: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'Product',
        required: true
    },
    userId: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'User',
        required: true
    },
    orderId: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'Order'
    },
    rating: {
        type: Number,
        required: true,
        min: 1,
        max: 5
    },
    title: {
        type: String,
        maxlength: 100
    },
    content: {
        type: String,
        required: true,
        maxlength: 1000
    },
    pros: [String],
    cons: [String],
    images: [String],
    isVerifiedPurchase: {
        type: Boolean,
        default: false
    },
    isApproved: {
        type: Boolean,
        default: false
    },
    helpfulVotes: {
        type: Number,
        default: 0
    },
    unhelpfulVotes: {
        type: Number,
        default: 0
    },
    votes: [{
        userId: { type: mongoose.Schema.Types.ObjectId, ref: 'User' },
        type: { type: String, enum: ['helpful', 'unhelpful'] }
    }],
    response: {
        content: String,
        respondedBy: { type: mongoose.Schema.Types.ObjectId, ref: 'User' },
        respondedAt: Date
    }
}, {
    timestamps: true
});

// Indexes
reviewSchema.index({ productId: 1, createdAt: -1 });
reviewSchema.index({ rating: 1 });
reviewSchema.index({ isApproved: 1 });

// Ensure one review per user per product
reviewSchema.index({ productId: 1, userId: 1 }, { unique: true });

// Pre-save middleware to set verified purchase flag
reviewSchema.pre('save', async function(next) {
    if (!this.orderId) {
        // Check if user has purchased this product
        const order = await Order.findOne({
            userId: this.userId,
            'items.productId': this.productId,
            paymentStatus: 'paid'
        });
        
        this.isVerifiedPurchase = !!order;
    }
    next();
});

// Post-save middleware to update product rating
reviewSchema.post('save', async function() {
    const Product = mongoose.model('Product');
    const product = await Product.findById(this.productId);
    await product.updateRating();
});

reviewSchema.post('remove', async function() {
    const Product = mongoose.model('Product');
    const product = await Product.findById(this.productId);
    await product.updateRating();
});

// Instance methods
reviewSchema.methods.vote = async function(userId, type) {
    const existingVote = this.votes.find(v => v.userId.toString() === userId.toString());
    
    if (existingVote) {
        // Remove existing vote
        if (existingVote.type === 'helpful') this.helpfulVotes--;
        if (existingVote.type === 'unhelpful') this.unhelpfulVotes--;
        
        this.votes = this.votes.filter(v => v.userId.toString() !== userId.toString());
        
        // Add new vote if different
        if (existingVote.type !== type) {
            this.votes.push({ userId, type });
            if (type === 'helpful') this.helpfulVotes++;
            if (type === 'unhelpful') this.unhelpfulVotes++;
        }
    } else {
        // Add new vote
        this.votes.push({ userId, type });
        if (type === 'helpful') this.helpfulVotes++;
        if (type === 'unhelpful') this.unhelpfulVotes++;
    }
    
    await this.save();
};

reviewSchema.methods.respond = async function(content, userId) {
    this.response = {
        content,
        respondedBy: userId,
        respondedAt: new Date()
    };
    await this.save();
};

const Review = mongoose.model('Review', reviewSchema);

// ============================================================================
// Cart Model (for guest users)
// ============================================================================

const cartItemSchema = new mongoose.Schema({
    productId: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'Product',
        required: true
    },
    quantity: {
        type: Number,
        required: true,
        min: 1
    },
    variant: Map,
    addedAt: {
        type: Date,
        default: Date.now
    }
});

const cartSchema = new mongoose.Schema({
    userId: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'User',
        sparse: true,
        unique: true
    },
    sessionId: {
        type: String,
        index: true,
        sparse: true
    },
    items: [cartItemSchema],
    couponCode: String,
    couponDiscount: Number,
    expiresAt: {
        type: Date,
        default: () => new Date(Date.now() + 30 * 24 * 60 * 60 * 1000) // 30 days
    }
}, {
    timestamps: true
});

// Index for cleanup job
cartSchema.index({ expiresAt: 1 }, { expireAfterSeconds: 0 });

// Virtual for cart totals
cartSchema.virtual('subtotal').get(function() {
    return this.items.reduce(async (acc, item) => {
        const product = await Product.findById(item.productId);
        return acc + (product?.price || 0) * item.quantity;
    }, 0);
});

// Instance methods
cartSchema.methods.addItem = async function(productId, quantity = 1, variant = null) {
    const product = await Product.findById(productId);
    if (!product) {
        throw new Error('Product not found');
    }
    
    const existingItem = this.items.find(item => 
        item.productId.toString() === productId.toString() &&
        JSON.stringify(item.variant) === JSON.stringify(variant)
    );
    
    if (existingItem) {
        existingItem.quantity += quantity;
    } else {
        this.items.push({ productId, quantity, variant });
    }
    
    await this.save();
};

cartSchema.methods.updateQuantity = function(productId, quantity, variant = null) {
    const item = this.items.find(item => 
        item.productId.toString() === productId.toString() &&
        JSON.stringify(item.variant) === JSON.stringify(variant)
    );
    
    if (!item) {
        throw new Error('Item not found in cart');
    }
    
    if (quantity <= 0) {
        this.items = this.items.filter(i => i !== item);
    } else {
        item.quantity = quantity;
    }
    
    return this.save();
};

cartSchema.methods.clear = function() {
    this.items = [];
    return this.save();
};

const Cart = mongoose.model('Cart', cartSchema);

// ============================================================================
// Inventory Model
// ============================================================================

const inventorySchema = new mongoose.Schema({
    productId: {
        type: mongoose.Schema.Types.ObjectId,
        ref: 'Product',
        required: true,
        unique: true
    },
    sku: {
        type: String,
        required: true,
        unique: true
    },
    quantity: {
        type: Number,
        required: true,
        default: 0,
        min: 0
    },
    reserved: {
        type: Number,
        default: 0,
        min: 0
    },
    reorderPoint: {
        type: Number,
        default: 10
    },
    reorderQuantity: {
        type: Number,
        default: 100
    },
    location: {
        warehouse: String,
        aisle: String,
        shelf: String,
        bin: String
    },
    supplier: {
        name: String,
        sku: String,
        leadTime: Number, // days
        minimumOrder: Number,
        price: Number
    },
    movements: [{
        type: { type: String, enum: ['in', 'out', 'adjust', 'return'] },
        quantity: Number,
        reference: String,
        userId: { type: mongoose.Schema.Types.ObjectId, ref: 'User' },
        notes: String,
        timestamp: { type: Date, default: Date.now }
    }]
}, {
    timestamps: true
});

// Indexes
inventorySchema.index({ sku: 1 });
inventorySchema.index({ 'supplier.name': 1 });

// Virtual for available stock
inventorySchema.virtual('available').get(function() {
    return this.quantity - this.reserved;
});

// Virtual for low stock alert
inventorySchema.virtual('isLowStock').get(function() {
    return this.available <= this.reorderPoint;
});

// Instance methods
inventorySchema.methods.receive = async function(quantity, reference, userId, notes) {
    this.quantity += quantity;
    this.movements.push({
        type: 'in',
        quantity,
        reference,
        userId,
        notes
    });
    await this.save();
};

inventorySchema.methods.ship = async function(quantity, reference, userId, notes) {
    if (this.available < quantity) {
        throw new Error('Insufficient available stock');
    }
    
    this.reserved -= quantity;
    this.quantity -= quantity;
    this.movements.push({
        type: 'out',
        quantity,
        reference,
        userId,
        notes
    });
    await this.save();
};

inventorySchema.methods.reserve = async function(quantity, reference) {
    if (this.available < quantity) {
        throw new Error('Insufficient available stock');
    }
    
    this.reserved += quantity;
    this.movements.push({
        type: 'reserve',
        quantity,
        reference
    });
    await this.save();
};

inventorySchema.methods.release = async function(quantity, reference) {
    this.reserved = Math.max(0, this.reserved - quantity);
    this.movements.push({
        type: 'release',
        quantity,
        reference
    });
    await this.save();
};

inventorySchema.methods.adjust = async function(newQuantity, reason, userId) {
    const difference = newQuantity - this.quantity;
    this.quantity = newQuantity;
    this.movements.push({
        type: 'adjust',
        quantity: difference,
        reference: reason,
        userId,
        notes: `Adjusted from ${this.quantity - difference} to ${newQuantity}`
    });
    await this.save();
};

// Static methods
inventorySchema.statics.getLowStock = function() {
    return this.aggregate([
        {
            $addFields: {
                available: { $subtract: ['$quantity', '$reserved'] }
            }
        },
        {
            $match: {
                $expr: { $lte: ['$available', '$reorderPoint'] }
            }
        },
        {
            $lookup: {
                from: 'products',
                localField: 'productId',
                foreignField: '_id',
                as: 'product'
            }
        },
        { $unwind: '$product' }
    ]);
};

const Inventory = mongoose.model('Inventory', inventorySchema);

// ============================================================================
// Export all models
// ============================================================================

module.exports = {
    User,
    Category,
    Product,
    Order,
    Review,
    Cart,
    Inventory
};