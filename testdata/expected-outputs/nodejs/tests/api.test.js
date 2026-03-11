/**
 * API Tests
 * 
 * Tests for RESTful API endpoints with authentication, validation,
 * error handling, and edge cases.
 */

const request = require('supertest');
const app = require('../src/app');
const { User, Product, Order } = require('../src/models');
const { generateToken } = require('../src/utils/auth');
const { connectDB, disconnectDB } = require('../src/config/database');

// ============================================================================
// Test Setup and Teardown
// ============================================================================

beforeAll(async () => {
    await connectDB(process.env.TEST_DB_URL);
});

afterAll(async () => {
    await disconnectDB();
});

beforeEach(async () => {
    // Clear all collections before each test
    await User.deleteMany({});
    await Product.deleteMany({});
    await Order.deleteMany({});
    
    // Create test users
    const hashedPassword = await bcrypt.hash('password123', 10);
    testUser = await User.create({
        username: 'testuser',
        email: 'test@example.com',
        password: hashedPassword,
        role: 'user'
    });
    
    adminUser = await User.create({
        username: 'admin',
        email: 'admin@example.com',
        password: hashedPassword,
        role: 'admin'
    });
    
    // Create test products
    testProduct = await Product.create({
        name: 'Test Product',
        price: 99.99,
        description: 'A test product',
        category: 'electronics',
        stock: 100,
        sku: 'TEST-001'
    });
    
    // Generate auth tokens
    userToken = generateToken(testUser);
    adminToken = generateToken(adminUser);
});

// ============================================================================
// Helper Functions
// ============================================================================

const createTestProduct = async (overrides = {}) => {
    const productData = {
        name: 'New Product',
        price: 49.99,
        description: 'A new product',
        category: 'books',
        stock: 50,
        sku: 'TEST-' + Date.now(),
        ...overrides
    };
    return await Product.create(productData);
};

const createTestOrder = async (userId, productId, overrides = {}) => {
    const orderData = {
        user: userId,
        items: [{
            product: productId,
            quantity: 2,
            price: 99.99
        }],
        total: 199.98,
        status: 'pending',
        shippingAddress: {
            street: '123 Test St',
            city: 'Test City',
            zipCode: '12345'
        },
        ...overrides
    };
    return await Order.create(orderData);
};

// ============================================================================
// Authentication & Authorization Tests
// ============================================================================

describe('Authentication & Authorization', () => {
    
    describe('POST /api/auth/register', () => {
        test('should register a new user with valid data', async () => {
            const userData = {
                username: 'newuser',
                email: 'newuser@example.com',
                password: 'SecurePass123!'
            };
            
            const response = await request(app)
                .post('/api/auth/register')
                .send(userData)
                .expect(201);
            
            expect(response.body).toHaveProperty('user');
            expect(response.body.user).toHaveProperty('id');
            expect(response.body.user.username).toBe(userData.username);
            expect(response.body.user.email).toBe(userData.email);
            expect(response.body.user).not.toHaveProperty('password');
            expect(response.body).toHaveProperty('token');
            
            // Verify user was saved to database
            const savedUser = await User.findOne({ email: userData.email });
            expect(savedUser).toBeTruthy();
            expect(savedUser.username).toBe(userData.username);
        });
        
        test('should return 400 with invalid email format', async () => {
            const userData = {
                username: 'newuser',
                email: 'invalid-email',
                password: 'SecurePass123!'
            };
            
            const response = await request(app)
                .post('/api/auth/register')
                .send(userData)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('email');
        });
        
        test('should return 400 with weak password', async () => {
            const userData = {
                username: 'newuser',
                email: 'newuser@example.com',
                password: '123'
            };
            
            const response = await request(app)
                .post('/api/auth/register')
                .send(userData)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('password');
        });
        
        test('should return 400 with missing required fields', async () => {
            const response = await request(app)
                .post('/api/auth/register')
                .send({ username: 'newuser' })
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('required');
        });
        
        test('should return 409 when email already exists', async () => {
            const userData = {
                username: 'anotheruser',
                email: 'test@example.com', // Already exists
                password: 'SecurePass123!'
            };
            
            const response = await request(app)
                .post('/api/auth/register')
                .send(userData)
                .expect(409);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('already exists');
        });
        
        test('should return 409 when username already exists', async () => {
            const userData = {
                username: 'testuser', // Already exists
                email: 'unique@example.com',
                password: 'SecurePass123!'
            };
            
            const response = await request(app)
                .post('/api/auth/register')
                .send(userData)
                .expect(409);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('already exists');
        });
        
        test('should return 400 with invalid username format', async () => {
            const userData = {
                username: 'user@name!',
                email: 'user@example.com',
                password: 'SecurePass123!'
            };
            
            const response = await request(app)
                .post('/api/auth/register')
                .send(userData)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('username');
        });
        
        test('should return 400 with very long username', async () => {
            const userData = {
                username: 'a'.repeat(51),
                email: 'user@example.com',
                password: 'SecurePass123!'
            };
            
            const response = await request(app)
                .post('/api/auth/register')
                .send(userData)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('username');
        });
    });
    
    describe('POST /api/auth/login', () => {
        test('should login with valid credentials', async () => {
            const loginData = {
                email: 'test@example.com',
                password: 'password123'
            };
            
            const response = await request(app)
                .post('/api/auth/login')
                .send(loginData)
                .expect(200);
            
            expect(response.body).toHaveProperty('token');
            expect(response.body).toHaveProperty('user');
            expect(response.body.user.email).toBe(loginData.email);
        });
        
        test('should return 401 with invalid password', async () => {
            const loginData = {
                email: 'test@example.com',
                password: 'wrongpassword'
            };
            
            const response = await request(app)
                .post('/api/auth/login')
                .send(loginData)
                .expect(401);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('Invalid');
        });
        
        test('should return 401 with non-existent email', async () => {
            const loginData = {
                email: 'nonexistent@example.com',
                password: 'password123'
            };
            
            const response = await request(app)
                .post('/api/auth/login')
                .send(loginData)
                .expect(401);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('Invalid');
        });
        
        test('should return 400 with missing fields', async () => {
            const response = await request(app)
                .post('/api/auth/login')
                .send({ email: 'test@example.com' })
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('required');
        });
        
        test('should return 400 with invalid email format', async () => {
            const loginData = {
                email: 'invalid',
                password: 'password123'
            };
            
            const response = await request(app)
                .post('/api/auth/login')
                .send(loginData)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('email');
        });
        
        test('should return 429 after too many failed attempts', async () => {
            // Make multiple failed attempts
            for (let i = 0; i < 6; i++) {
                await request(app)
                    .post('/api/auth/login')
                    .send({ email: 'test@example.com', password: 'wrong' });
            }
            
            // Next attempt should be rate limited
            const response = await request(app)
                .post('/api/auth/login')
                .send({ email: 'test@example.com', password: 'wrong' })
                .expect(429);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('Too many');
        });
    });
    
    describe('POST /api/auth/logout', () => {
        test('should logout authenticated user', async () => {
            const response = await request(app)
                .post('/api/auth/logout')
                .set('Authorization', `Bearer ${userToken}`)
                .expect(200);
            
            expect(response.body).toHaveProperty('message');
            expect(response.body.message).toContain('success');
        });
        
        test('should return 401 without token', async () => {
            const response = await request(app)
                .post('/api/auth/logout')
                .expect(401);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('token');
        });
        
        test('should return 401 with invalid token', async () => {
            const response = await request(app)
                .post('/api/auth/logout')
                .set('Authorization', 'Bearer invalid.token.here')
                .expect(401);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('Invalid');
        });
        
        test('should return 401 with expired token', async () => {
            // Create expired token
            const expiredToken = generateToken(testUser, { expiresIn: '-1h' });
            
            const response = await request(app)
                .post('/api/auth/logout')
                .set('Authorization', `Bearer ${expiredToken}`)
                .expect(401);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('expired');
        });
    });
    
    describe('GET /api/auth/me', () => {
        test('should get current user profile', async () => {
            const response = await request(app)
                .get('/api/auth/me')
                .set('Authorization', `Bearer ${userToken}`)
                .expect(200);
            
            expect(response.body).toHaveProperty('user');
            expect(response.body.user.email).toBe('test@example.com');
            expect(response.body.user.username).toBe('testuser');
        });
        
        test('should return 401 without token', async () => {
            const response = await request(app)
                .get('/api/auth/me')
                .expect(401);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 404 for deleted user', async () => {
            // Delete the user
            await User.findByIdAndDelete(testUser._id);
            
            const response = await request(app)
                .get('/api/auth/me')
                .set('Authorization', `Bearer ${userToken}`)
                .expect(404);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('not found');
        });
    });
});

// ============================================================================
// User API Tests
// ============================================================================

describe('User API', () => {
    
    describe('GET /api/users', () => {
        test('should get all users as admin', async () => {
            const response = await request(app)
                .get('/api/users')
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(200);
            
            expect(Array.isArray(response.body)).toBe(true);
            expect(response.body.length).toBe(2); // testUser and adminUser
            expect(response.body[0]).not.toHaveProperty('password');
        });
        
        test('should return 403 for non-admin users', async () => {
            const response = await request(app)
                .get('/api/users')
                .set('Authorization', `Bearer ${userToken}`)
                .expect(403);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('Forbidden');
        });
        
        test('should support pagination', async () => {
            // Create additional users
            for (let i = 0; i < 15; i++) {
                await User.create({
                    username: `user${i}`,
                    email: `user${i}@example.com`,
                    password: 'password123'
                });
            }
            
            const response = await request(app)
                .get('/api/users?page=2&limit=5')
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(200);
            
            expect(response.body).toHaveProperty('data');
            expect(response.body).toHaveProperty('pagination');
            expect(response.body.data.length).toBe(5);
            expect(response.body.pagination.page).toBe(2);
            expect(response.body.pagination.total).toBeGreaterThan(15);
        });
        
        test('should support filtering', async () => {
            const response = await request(app)
                .get('/api/users?role=admin')
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(200);
            
            expect(response.body.length).toBe(1);
            expect(response.body[0].role).toBe('admin');
        });
        
        test('should support sorting', async () => {
            const response = await request(app)
                .get('/api/users?sort=username&order=desc')
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(200);
            
            expect(response.body[0].username).toBe('testuser');
        });
    });
    
    describe('GET /api/users/:id', () => {
        test('should get user by id', async () => {
            const response = await request(app)
                .get(`/api/users/${testUser._id}`)
                .set('Authorization', `Bearer ${userToken}`)
                .expect(200);
            
            expect(response.body).toHaveProperty('user');
            expect(response.body.user.email).toBe('test@example.com');
        });
        
        test('should return 404 for non-existent user', async () => {
            const fakeId = '507f1f77bcf86cd799439011';
            
            const response = await request(app)
                .get(`/api/users/${fakeId}`)
                .set('Authorization', `Bearer ${userToken}`)
                .expect(404);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('not found');
        });
        
        test('should return 400 for invalid user id', async () => {
            const response = await request(app)
                .get('/api/users/invalid-id')
                .set('Authorization', `Bearer ${userToken}`)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('Invalid');
        });
        
        test('should return 403 when user tries to access another user', async () => {
            const response = await request(app)
                .get(`/api/users/${adminUser._id}`)
                .set('Authorization', `Bearer ${userToken}`)
                .expect(403);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('Forbidden');
        });
    });
    
    describe('PUT /api/users/:id', () => {
        test('should update user profile', async () => {
            const updates = {
                username: 'updateduser',
                email: 'updated@example.com'
            };
            
            const response = await request(app)
                .put(`/api/users/${testUser._id}`)
                .set('Authorization', `Bearer ${userToken}`)
                .send(updates)
                .expect(200);
            
            expect(response.body).toHaveProperty('user');
            expect(response.body.user.username).toBe('updateduser');
            expect(response.body.user.email).toBe('updated@example.com');
            
            // Verify in database
            const updatedUser = await User.findById(testUser._id);
            expect(updatedUser.username).toBe('updateduser');
        });
        
        test('should not update password via this endpoint', async () => {
            const updates = {
                password: 'newpassword123'
            };
            
            const response = await request(app)
                .put(`/api/users/${testUser._id}`)
                .set('Authorization', `Bearer ${userToken}`)
                .send(updates)
                .expect(200);
            
            // Password should not be changed
            const updatedUser = await User.findById(testUser._id);
            expect(await bcrypt.compare('password123', updatedUser.password)).toBe(true);
        });
        
        test('should return 409 when updating to existing email', async () => {
            const updates = {
                email: 'admin@example.com'
            };
            
            const response = await request(app)
                .put(`/api/users/${testUser._id}`)
                .set('Authorization', `Bearer ${userToken}`)
                .send(updates)
                .expect(409);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('already exists');
        });
        
        test('should return 400 with invalid email', async () => {
            const updates = {
                email: 'not-an-email'
            };
            
            const response = await request(app)
                .put(`/api/users/${testUser._id}`)
                .set('Authorization', `Bearer ${userToken}`)
                .send(updates)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
        });
    });
    
    describe('DELETE /api/users/:id', () => {
        test('should delete user as admin', async () => {
            const response = await request(app)
                .delete(`/api/users/${testUser._id}`)
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(200);
            
            expect(response.body).toHaveProperty('message');
            
            // Verify deletion
            const deletedUser = await User.findById(testUser._id);
            expect(deletedUser).toBeNull();
        });
        
        test('should return 403 for non-admin users', async () => {
            const response = await request(app)
                .delete(`/api/users/${testUser._id}`)
                .set('Authorization', `Bearer ${userToken}`)
                .expect(403);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 404 for non-existent user', async () => {
            const fakeId = '507f1f77bcf86cd799439011';
            
            const response = await request(app)
                .delete(`/api/users/${fakeId}`)
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(404);
            
            expect(response.body).toHaveProperty('error');
        });
    });
});

// ============================================================================
// Product API Tests
// ============================================================================

describe('Product API', () => {
    
    describe('GET /api/products', () => {
        test('should get all products (public)', async () => {
            const response = await request(app)
                .get('/api/products')
                .expect(200);
            
            expect(Array.isArray(response.body)).toBe(true);
            expect(response.body.length).toBe(1); // testProduct
            expect(response.body[0].name).toBe('Test Product');
        });
        
        test('should support filtering by category', async () => {
            await createTestProduct({ category: 'books' });
            await createTestProduct({ category: 'electronics' });
            
            const response = await request(app)
                .get('/api/products?category=electronics')
                .expect(200);
            
            expect(response.body.length).toBe(2); // original + new electronics
            response.body.forEach(p => {
                expect(p.category).toBe('electronics');
            });
        });
        
        test('should support price range filtering', async () => {
            await createTestProduct({ price: 10 });
            await createTestProduct({ price: 50 });
            await createTestProduct({ price: 100 });
            
            const response = await request(app)
                .get('/api/products?minPrice=30&maxPrice=80')
                .expect(200);
            
            expect(response.body.length).toBe(1);
            expect(response.body[0].price).toBe(50);
        });
        
        test('should support search by name', async () => {
            await createTestProduct({ name: 'Laptop Pro' });
            await createTestProduct({ name: 'Laptop Air' });
            await createTestProduct({ name: 'Tablet' });
            
            const response = await request(app)
                .get('/api/products?search=Laptop')
                .expect(200);
            
            expect(response.body.length).toBe(2);
            response.body.forEach(p => {
                expect(p.name).toContain('Laptop');
            });
        });
        
        test('should return empty array when no products match', async () => {
            const response = await request(app)
                .get('/api/products?category=nonexistent')
                .expect(200);
            
            expect(Array.isArray(response.body)).toBe(true);
            expect(response.body.length).toBe(0);
        });
    });
    
    describe('GET /api/products/:id', () => {
        test('should get product by id', async () => {
            const response = await request(app)
                .get(`/api/products/${testProduct._id}`)
                .expect(200);
            
            expect(response.body).toHaveProperty('product');
            expect(response.body.product.name).toBe('Test Product');
            expect(response.body.product.price).toBe(99.99);
        });
        
        test('should return 404 for non-existent product', async () => {
            const fakeId = '507f1f77bcf86cd799439011';
            
            const response = await request(app)
                .get(`/api/products/${fakeId}`)
                .expect(404);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 400 for invalid product id', async () => {
            const response = await request(app)
                .get('/api/products/invalid-id')
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
        });
    });
    
    describe('POST /api/products', () => {
        test('should create new product as admin', async () => {
            const newProduct = {
                name: 'New Product',
                price: 149.99,
                description: 'Brand new product',
                category: 'electronics',
                stock: 50,
                sku: 'NEW-001'
            };
            
            const response = await request(app)
                .post('/api/products')
                .set('Authorization', `Bearer ${adminToken}`)
                .send(newProduct)
                .expect(201);
            
            expect(response.body).toHaveProperty('product');
            expect(response.body.product.name).toBe(newProduct.name);
            expect(response.body.product).toHaveProperty('id');
            
            // Verify in database
            const savedProduct = await Product.findById(response.body.product.id);
            expect(savedProduct).toBeTruthy();
            expect(savedProduct.sku).toBe(newProduct.sku);
        });
        
        test('should return 403 for non-admin users', async () => {
            const newProduct = {
                name: 'New Product',
                price: 149.99,
                category: 'electronics'
            };
            
            const response = await request(app)
                .post('/api/products')
                .set('Authorization', `Bearer ${userToken}`)
                .send(newProduct)
                .expect(403);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 401 without token', async () => {
            const newProduct = {
                name: 'New Product',
                price: 149.99
            };
            
            const response = await request(app)
                .post('/api/products')
                .send(newProduct)
                .expect(401);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 400 with missing required fields', async () => {
            const response = await request(app)
                .post('/api/products')
                .set('Authorization', `Bearer ${adminToken}`)
                .send({ name: 'Incomplete' })
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('required');
        });
        
        test('should return 400 with negative price', async () => {
            const newProduct = {
                name: 'New Product',
                price: -10,
                category: 'electronics'
            };
            
            const response = await request(app)
                .post('/api/products')
                .set('Authorization', `Bearer ${adminToken}`)
                .send(newProduct)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('price');
        });
        
        test('should return 409 with duplicate SKU', async () => {
            const newProduct = {
                name: 'Another Product',
                price: 199.99,
                category: 'electronics',
                sku: 'TEST-001' // Already exists
            };
            
            const response = await request(app)
                .post('/api/products')
                .set('Authorization', `Bearer ${adminToken}`)
                .send(newProduct)
                .expect(409);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('SKU already exists');
        });
    });
    
    describe('PUT /api/products/:id', () => {
        test('should update product as admin', async () => {
            const updates = {
                name: 'Updated Product',
                price: 79.99,
                stock: 75
            };
            
            const response = await request(app)
                .put(`/api/products/${testProduct._id}`)
                .set('Authorization', `Bearer ${adminToken}`)
                .send(updates)
                .expect(200);
            
            expect(response.body).toHaveProperty('product');
            expect(response.body.product.name).toBe('Updated Product');
            expect(response.body.product.price).toBe(79.99);
            expect(response.body.product.stock).toBe(75);
            
            // Verify in database
            const updatedProduct = await Product.findById(testProduct._id);
            expect(updatedProduct.name).toBe('Updated Product');
        });
        
        test('should return 403 for non-admin users', async () => {
            const updates = { price: 49.99 };
            
            const response = await request(app)
                .put(`/api/products/${testProduct._id}`)
                .set('Authorization', `Bearer ${userToken}`)
                .send(updates)
                .expect(403);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 404 for non-existent product', async () => {
            const fakeId = '507f1f77bcf86cd799439011';
            
            const response = await request(app)
                .put(`/api/products/${fakeId}`)
                .set('Authorization', `Bearer ${adminToken}`)
                .send({ name: 'Updated' })
                .expect(404);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 400 with invalid update data', async () => {
            const updates = { price: -5 };
            
            const response = await request(app)
                .put(`/api/products/${testProduct._id}`)
                .set('Authorization', `Bearer ${adminToken}`)
                .send(updates)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
        });
    });
    
    describe('DELETE /api/products/:id', () => {
        test('should delete product as admin', async () => {
            const response = await request(app)
                .delete(`/api/products/${testProduct._id}`)
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(200);
            
            expect(response.body).toHaveProperty('message');
            
            // Verify deletion
            const deletedProduct = await Product.findById(testProduct._id);
            expect(deletedProduct).toBeNull();
        });
        
        test('should return 403 for non-admin users', async () => {
            const response = await request(app)
                .delete(`/api/products/${testProduct._id}`)
                .set('Authorization', `Bearer ${userToken}`)
                .expect(403);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 404 for non-existent product', async () => {
            const fakeId = '507f1f77bcf86cd799439011';
            
            const response = await request(app)
                .delete(`/api/products/${fakeId}`)
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(404);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 409 if product has orders', async () => {
            await createTestOrder(testUser._id, testProduct._id);
            
            const response = await request(app)
                .delete(`/api/products/${testProduct._id}`)
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(409);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('has orders');
        });
    });
});

// ============================================================================
// Order API Tests
// ============================================================================

describe('Order API', () => {
    
    let testOrder;
    
    beforeEach(async () => {
        testOrder = await createTestOrder(testUser._id, testProduct._id);
    });
    
    describe('GET /api/orders', () => {
        test('should get user orders', async () => {
            const response = await request(app)
                .get('/api/orders')
                .set('Authorization', `Bearer ${userToken}`)
                .expect(200);
            
            expect(Array.isArray(response.body)).toBe(true);
            expect(response.body.length).toBe(1);
            expect(response.body[0]._id.toString()).toBe(testOrder._id.toString());
        });
        
        test('should get all orders as admin', async () => {
            const response = await request(app)
                .get('/api/orders')
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(200);
            
            expect(Array.isArray(response.body)).toBe(true);
            expect(response.body.length).toBe(1);
        });
        
        test('should support filtering by status', async () => {
            await createTestOrder(testUser._id, testProduct._id, { status: 'shipped' });
            await createTestOrder(testUser._id, testProduct._id, { status: 'delivered' });
            
            const response = await request(app)
                .get('/api/orders?status=pending')
                .set('Authorization', `Bearer ${userToken}`)
                .expect(200);
            
            expect(response.body.length).toBe(1);
            expect(response.body[0].status).toBe('pending');
        });
        
        test('should support date range filtering', async () => {
            const response = await request(app)
                .get('/api/orders?from=2024-01-01&to=2024-12-31')
                .set('Authorization', `Bearer ${userToken}`)
                .expect(200);
            
            expect(Array.isArray(response.body)).toBe(true);
        });
    });
    
    describe('GET /api/orders/:id', () => {
        test('should get order by id', async () => {
            const response = await request(app)
                .get(`/api/orders/${testOrder._id}`)
                .set('Authorization', `Bearer ${userToken}`)
                .expect(200);
            
            expect(response.body).toHaveProperty('order');
            expect(response.body.order._id.toString()).toBe(testOrder._id.toString());
            expect(response.body.order.items).toHaveLength(1);
            expect(response.body.order.total).toBe(199.98);
        });
        
        test('should return 404 for non-existent order', async () => {
            const fakeId = '507f1f77bcf86cd799439011';
            
            const response = await request(app)
                .get(`/api/orders/${fakeId}`)
                .set('Authorization', `Bearer ${userToken}`)
                .expect(404);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 403 when accessing another user order', async () => {
            const response = await request(app)
                .get(`/api/orders/${testOrder._id}`)
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(403);
            
            expect(response.body).toHaveProperty('error');
        });
    });
    
    describe('POST /api/orders', () => {
        test('should create new order', async () => {
            const newOrder = {
                items: [{
                    productId: testProduct._id,
                    quantity: 3
                }],
                shippingAddress: {
                    street: '456 New St',
                    city: 'New City',
                    zipCode: '67890'
                }
            };
            
            const response = await request(app)
                .post('/api/orders')
                .set('Authorization', `Bearer ${userToken}`)
                .send(newOrder)
                .expect(201);
            
            expect(response.body).toHaveProperty('order');
            expect(response.body.order.items).toHaveLength(1);
            expect(response.body.order.total).toBe(99.99 * 3);
            expect(response.body.order.status).toBe('pending');
            
            // Verify stock was reduced
            const updatedProduct = await Product.findById(testProduct._id);
            expect(updatedProduct.stock).toBe(97); // Was 100, sold 3
        });
        
        test('should return 400 with insufficient stock', async () => {
            const newOrder = {
                items: [{
                    productId: testProduct._id,
                    quantity: 200 // More than available
                }],
                shippingAddress: {
                    street: '123 Test St',
                    city: 'Test City',
                    zipCode: '12345'
                }
            };
            
            const response = await request(app)
                .post('/api/orders')
                .set('Authorization', `Bearer ${userToken}`)
                .send(newOrder)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('Insufficient stock');
        });
        
        test('should return 400 with invalid product', async () => {
            const newOrder = {
                items: [{
                    productId: '507f1f77bcf86cd799439011',
                    quantity: 1
                }],
                shippingAddress: {
                    street: '123 Test St',
                    city: 'Test City',
                    zipCode: '12345'
                }
            };
            
            const response = await request(app)
                .post('/api/orders')
                .set('Authorization', `Bearer ${userToken}`)
                .send(newOrder)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('Product not found');
        });
        
        test('should return 400 with missing shipping address', async () => {
            const newOrder = {
                items: [{
                    productId: testProduct._id,
                    quantity: 1
                }]
            };
            
            const response = await request(app)
                .post('/api/orders')
                .set('Authorization', `Bearer ${userToken}`)
                .send(newOrder)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('shipping address');
        });
        
        test('should handle multiple items in order', async () => {
            const product2 = await createTestProduct({ sku: 'TEST-002' });
            
            const newOrder = {
                items: [{
                    productId: testProduct._id,
                    quantity: 2
                }, {
                    productId: product2._id,
                    quantity: 1
                }],
                shippingAddress: {
                    street: '123 Test St',
                    city: 'Test City',
                    zipCode: '12345'
                }
            };
            
            const response = await request(app)
                .post('/api/orders')
                .set('Authorization', `Bearer ${userToken}`)
                .send(newOrder)
                .expect(201);
            
            expect(response.body.order.items).toHaveLength(2);
            expect(response.body.order.total).toBe(99.99 * 2 + product2.price);
        });
    });
    
    describe('PUT /api/orders/:id/status', () => {
        test('should update order status as admin', async () => {
            const response = await request(app)
                .put(`/api/orders/${testOrder._id}/status`)
                .set('Authorization', `Bearer ${adminToken}`)
                .send({ status: 'shipped' })
                .expect(200);
            
            expect(response.body).toHaveProperty('order');
            expect(response.body.order.status).toBe('shipped');
            
            // Verify in database
            const updatedOrder = await Order.findById(testOrder._id);
            expect(updatedOrder.status).toBe('shipped');
        });
        
        test('should return 403 for non-admin users', async () => {
            const response = await request(app)
                .put(`/api/orders/${testOrder._id}/status`)
                .set('Authorization', `Bearer ${userToken}`)
                .send({ status: 'shipped' })
                .expect(403);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 400 with invalid status', async () => {
            const response = await request(app)
                .put(`/api/orders/${testOrder._id}/status`)
                .set('Authorization', `Bearer ${adminToken}`)
                .send({ status: 'invalid-status' })
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 404 for non-existent order', async () => {
            const fakeId = '507f1f77bcf86cd799439011';
            
            const response = await request(app)
                .put(`/api/orders/${fakeId}/status`)
                .set('Authorization', `Bearer ${adminToken}`)
                .send({ status: 'shipped' })
                .expect(404);
            
            expect(response.body).toHaveProperty('error');
        });
    });
    
    describe('POST /api/orders/:id/cancel', () => {
        test('should cancel pending order', async () => {
            const response = await request(app)
                .post(`/api/orders/${testOrder._id}/cancel`)
                .set('Authorization', `Bearer ${userToken}`)
                .expect(200);
            
            expect(response.body).toHaveProperty('order');
            expect(response.body.order.status).toBe('cancelled');
            
            // Verify stock was restored
            const updatedProduct = await Product.findById(testProduct._id);
            expect(updatedProduct.stock).toBe(102); // Was 100, sold 2, now restored
        });
        
        test('should not cancel shipped order', async () => {
            testOrder.status = 'shipped';
            await testOrder.save();
            
            const response = await request(app)
                .post(`/api/orders/${testOrder._id}/cancel`)
                .set('Authorization', `Bearer ${userToken}`)
                .expect(400);
            
            expect(response.body).toHaveProperty('error');
            expect(response.body.error).toContain('cannot be cancelled');
        });
        
        test('should return 404 for non-existent order', async () => {
            const fakeId = '507f1f77bcf86cd799439011';
            
            const response = await request(app)
                .post(`/api/orders/${fakeId}/cancel`)
                .set('Authorization', `Bearer ${userToken}`)
                .expect(404);
            
            expect(response.body).toHaveProperty('error');
        });
        
        test('should return 403 when cancelling another user order', async () => {
            const response = await request(app)
                .post(`/api/orders/${testOrder._id}/cancel`)
                .set('Authorization', `Bearer ${adminToken}`)
                .expect(403);
            
            expect(response.body).toHaveProperty('error');
        });
    });
});

// ============================================================================
// Health Check Tests
// ============================================================================

describe('Health Check', () => {
    
    test('GET /health should return OK', async () => {
        const response = await request(app)
            .get('/health')
            .expect(200);
        
        expect(response.body).toHaveProperty('status');
        expect(response.body.status).toBe('OK');
        expect(response.body).toHaveProperty('timestamp');
        expect(response.body).toHaveProperty('uptime');
    });
    
    test('GET /health/db should return database status', async () => {
        const response = await request(app)
            .get('/health/db')
            .expect(200);
        
        expect(response.body).toHaveProperty('database');
        expect(response.body.database).toBe('connected');
    });
    
    test('GET /health/redis should return cache status', async () => {
        const response = await request(app)
            .get('/health/redis')
            .expect(200);
        
        expect(response.body).toHaveProperty('redis');
        expect(response.body.redis).toBe('connected');
    });
});

// ============================================================================
// Error Handling Tests
// ============================================================================

describe('Error Handling', () => {
    
    test('should handle 404 for unknown routes', async () => {
        const response = await request(app)
            .get('/api/unknown-route')
            .expect(404);
        
        expect(response.body).toHaveProperty('error');
        expect(response.body.error).toContain('Not found');
    });
    
    test('should handle JSON parsing errors', async () => {
        const response = await request(app)
            .post('/api/auth/login')
            .set('Content-Type', 'application/json')
            .send('{invalid json}')
            .expect(400);
        
        expect(response.body).toHaveProperty('error');
    });
    
    test('should handle large payloads', async () => {
        const largeData = { data: 'x'.repeat(10 * 1024 * 1024) }; // 10MB
        
        const response = await request(app)
            .post('/api/test/large')
            .send(largeData)
            .expect(413); // Payload too large
    });
    
    test('should handle concurrent requests', async () => {
        const promises = [];
        for (let i = 0; i < 10; i++) {
            promises.push(
                request(app)
                    .get('/api/products')
                    .expect(200)
            );
        }
        
        const responses = await Promise.all(promises);
        responses.forEach(response => {
            expect(response.status).toBe(200);
        });
    });
});
```

## ✅ **What this API test file covers:**

| Test Suite | Description |
|------------|-------------|
| **Authentication** | Register, login, logout, profile |
| **Authorization** | Role-based access (user vs admin) |
| **User API** | CRUD operations with pagination, filtering |
| **Product API** | Public/private endpoints, validation |
| **Order API** | Order creation, status updates, cancellation |
| **Health Checks** | Service status endpoints |
| **Error Handling** | 400, 401, 403, 404, 409 responses |
| **Edge Cases** | Invalid data, duplicates, concurrency |

## 🎯 **Key Testing Patterns:**

- ✅ **Setup/Teardown** - Database cleanup between tests
- ✅ **Authentication** - Token generation and validation
- ✅ **Authorization** - Role-based access control
- ✅ **CRUD Operations** - Create, read, update, delete
- ✅ **Pagination** - Page/limit parameters
- ✅ **Filtering** - Query parameter filtering
- ✅ **Validation** - Input validation errors
- ✅ **Error Responses** - Proper HTTP status codes
- ✅ **Business Logic** - Stock management, order processing