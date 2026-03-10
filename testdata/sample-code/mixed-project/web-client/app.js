/**
 * Web Client Application
 * 
 * Modern frontend application for the mixed-language project.
 * Features:
 * - API integration with Go backend
 * - Real-time updates
 * - State management
 * - Component architecture
 * - Form validation
 * - Error handling
 */

// ============================================================================
// Application State
// ============================================================================

const AppState = {
    // User state
    user: null,
    token: localStorage.getItem('token'),
    
    // UI state
    loading: false,
    error: null,
    notifications: [],
    
    // Data state
    products: [],
    orders: [],
    users: [],
    
    // Settings
    theme: localStorage.getItem('theme') || 'light',
    language: localStorage.getItem('language') || 'en',
    
    // Update listeners
    listeners: new Set(),
    
    // Subscribe to state changes
    subscribe(callback) {
        this.listeners.add(callback);
        return () => this.listeners.delete(callback);
    },
    
    // Notify all listeners
    notify() {
        this.listeners.forEach(callback => callback(this));
    },
    
    // Update state
    setState(updates) {
        Object.assign(this, updates);
        this.notify();
    },
    
    // Reset state
    reset() {
        this.user = null;
        this.token = null;
        this.loading = false;
        this.error = null;
        this.products = [];
        this.orders = [];
        this.users = [];
        localStorage.removeItem('token');
        this.notify();
    }
};

// ============================================================================
// API Client
// ============================================================================

class APIClient {
    constructor(baseURL) {
        this.baseURL = baseURL || 'http://localhost:8080/api/v1';
        this.timeout = 30000; // 30 seconds
        this.retries = 3;
        this.retryDelay = 1000;
    }
    
    getHeaders() {
        const headers = {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
        };
        
        if (AppState.token) {
            headers['Authorization'] = `Bearer ${AppState.token}`;
        }
        
        return headers;
    }
    
    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const headers = this.getHeaders();
        
        let lastError;
        for (let i = 0; i < this.retries; i++) {
            try {
                const controller = new AbortController();
                const timeoutId = setTimeout(() => controller.abort(), this.timeout);
                
                const response = await fetch(url, {
                    ...options,
                    headers: { ...headers, ...options.headers },
                    signal: controller.signal
                });
                
                clearTimeout(timeoutId);
                
                if (!response.ok) {
                    const error = await response.json().catch(() => ({}));
                    throw new APIError(
                        error.message || `HTTP ${response.status}`,
                        response.status,
                        error
                    );
                }
                
                return await response.json();
                
            } catch (error) {
                lastError = error;
                
                if (error.name === 'AbortError') {
                    throw new APIError('Request timeout', 408);
                }
                
                if (i < this.retries - 1 && this.shouldRetry(error)) {
                    await new Promise(resolve => setTimeout(resolve, this.retryDelay * Math.pow(2, i)));
                    continue;
                }
                
                throw error;
            }
        }
        
        throw lastError;
    }
    
    shouldRetry(error) {
        // Retry on network errors and 5xx status codes
        if (error instanceof APIError) {
            return error.status >= 500 || error.status === 429;
        }
        return true;
    }
    
    // Auth endpoints
    async login(email, password) {
        return this.request('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ email, password })
        });
    }
    
    async register(userData) {
        return this.request('/auth/register', {
            method: 'POST',
            body: JSON.stringify(userData)
        });
    }
    
    async logout() {
        return this.request('/auth/logout', { method: 'POST' });
    }
    
    async refreshToken() {
        return this.request('/auth/refresh', { method: 'POST' });
    }
    
    // User endpoints
    async getCurrentUser() {
        return this.request('/users/me');
    }
    
    async getUsers(params = {}) {
        const query = new URLSearchParams(params).toString();
        return this.request(`/users${query ? '?' + query : ''}`);
    }
    
    async getUser(id) {
        return this.request(`/users/${id}`);
    }
    
    async updateUser(id, data) {
        return this.request(`/users/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }
    
    async deleteUser(id) {
        return this.request(`/users/${id}`, { method: 'DELETE' });
    }
    
    // Product endpoints
    async getProducts(params = {}) {
        const query = new URLSearchParams(params).toString();
        return this.request(`/products${query ? '?' + query : ''}`);
    }
    
    async getProduct(id) {
        return this.request(`/products/${id}`);
    }
    
    async createProduct(product) {
        return this.request('/products', {
            method: 'POST',
            body: JSON.stringify(product)
        });
    }
    
    async updateProduct(id, product) {
        return this.request(`/products/${id}`, {
            method: 'PUT',
            body: JSON.stringify(product)
        });
    }
    
    async deleteProduct(id) {
        return this.request(`/products/${id}`, { method: 'DELETE' });
    }
    
    // Order endpoints
    async getOrders(params = {}) {
        const query = new URLSearchParams(params).toString();
        return this.request(`/orders${query ? '?' + query : ''}`);
    }
    
    async getOrder(id) {
        return this.request(`/orders/${id}`);
    }
    
    async createOrder(order) {
        return this.request('/orders', {
            method: 'POST',
            body: JSON.stringify(order)
        });
    }
    
    async updateOrderStatus(id, status) {
        return this.request(`/orders/${id}/status`, {
            method: 'PUT',
            body: JSON.stringify({ status })
        });
    }
}

class APIError extends Error {
    constructor(message, status, data = {}) {
        super(message);
        this.name = 'APIError';
        this.status = status;
        this.data = data;
    }
}

// Initialize API client
const API = new APIClient();

// ============================================================================
// Router
// ============================================================================

class Router {
    constructor() {
        this.routes = new Map();
        this.currentRoute = null;
        this.params = {};
        this.query = {};
        
        window.addEventListener('popstate', () => this.handleRoute());
    }
    
    addRoute(path, handler) {
        this.routes.set(path, handler);
    }
    
    navigate(path, replace = false) {
        if (replace) {
            window.history.replaceState({}, '', path);
        } else {
            window.history.pushState({}, '', path);
        }
        this.handleRoute();
    }
    
    handleRoute() {
        const path = window.location.pathname;
        const query = Object.fromEntries(new URLSearchParams(window.location.search));
        
        // Find matching route
        for (const [routePath, handler] of this.routes) {
            const params = this.matchRoute(routePath, path);
            if (params) {
                this.currentRoute = routePath;
                this.params = params;
                this.query = query;
                handler(params, query);
                return;
            }
        }
        
        // 404 - Not found
        this.currentRoute = null;
        this.params = {};
        this.query = query;
        this.routes.get('404')?.();
    }
    
    matchRoute(routePath, actualPath) {
        const routeParts = routePath.split('/').filter(p => p);
        const actualParts = actualPath.split('/').filter(p => p);
        
        if (routeParts.length !== actualParts.length) return null;
        
        const params = {};
        
        for (let i = 0; i < routeParts.length; i++) {
            if (routeParts[i].startsWith(':')) {
                params[routeParts[i].slice(1)] = actualParts[i];
            } else if (routeParts[i] !== actualParts[i]) {
                return null;
            }
        }
        
        return params;
    }
    
    getLink(path) {
        return (e) => {
            e.preventDefault();
            this.navigate(path);
        };
    }
}

// Initialize router
const router = new Router();

// ============================================================================
// Component Base Class
// ============================================================================

class Component {
    constructor(props = {}) {
        this.props = props;
        this.state = {};
        this.element = null;
        this.children = [];
        this.unsubscribe = null;
    }
    
    setState(updates) {
        this.state = { ...this.state, ...updates };
        this.update();
    }
    
    mount(container) {
        this.element = this.render();
        container.appendChild(this.element);
        this.onMount();
        
        // Subscribe to state changes
        this.unsubscribe = AppState.subscribe(() => this.update());
        
        return this.element;
    }
    
    unmount() {
        if (this.unsubscribe) {
            this.unsubscribe();
        }
        this.onUnmount();
        this.children.forEach(child => child.unmount());
    }
    
    update() {
        if (this.element && this.element.parentNode) {
            const newElement = this.render();
            this.element.replaceWith(newElement);
            this.element = newElement;
            this.onUpdate();
        }
    }
    
    // Lifecycle hooks
    onMount() {}
    onUnmount() {}
    onUpdate() {}
    
    // Render method - to be overridden
    render() {
        throw new Error('Render method must be implemented');
    }
    
    // Utility methods
    createElement(tag, attributes = {}, ...children) {
        const element = document.createElement(tag);
        
        Object.entries(attributes).forEach(([key, value]) => {
            if (key === 'className') {
                element.className = value;
            } else if (key === 'style' && typeof value === 'object') {
                Object.assign(element.style, value);
            } else if (key.startsWith('on') && typeof value === 'function') {
                element.addEventListener(key.slice(2).toLowerCase(), value);
            } else {
                element.setAttribute(key, value);
            }
        });
        
        children.flat().forEach(child => {
            if (child instanceof Component) {
                this.children.push(child);
                child.mount(element);
            } else if (child instanceof Node) {
                element.appendChild(child);
            } else if (child !== null && child !== undefined) {
                element.appendChild(document.createTextNode(child));
            }
        });
        
        return element;
    }
    
    createFragment(...children) {
        const fragment = document.createDocumentFragment();
        children.flat().forEach(child => {
            if (child instanceof Component) {
                this.children.push(child);
                child.mount(fragment);
            } else if (child instanceof Node) {
                fragment.appendChild(child);
            } else if (child !== null && child !== undefined) {
                fragment.appendChild(document.createTextNode(child));
            }
        });
        return fragment;
    }
}

// ============================================================================
// UI Components
// ============================================================================

class Button extends Component {
    render() {
        const {
            type = 'button',
            variant = 'primary',
            size = 'md',
            disabled = false,
            loading = false,
            onClick,
            children
        } = this.props;
        
        const className = `btn btn-${variant} btn-${size} ${loading ? 'loading' : ''}`;
        
        return this.createElement(
            'button',
            {
                type,
                className,
                disabled: disabled || loading,
                onClick
            },
            loading ? 'Loading...' : children
        );
    }
}

class Input extends Component {
    constructor(props) {
        super(props);
        this.state = {
            value: props.value || '',
            error: null,
            touched: false
        };
    }
    
    handleChange(e) {
        const value = e.target.value;
        this.setState({ value, touched: true });
        
        // Validate
        if (this.props.validate) {
            const error = this.props.validate(value);
            this.setState({ error });
        }
        
        if (this.props.onChange) {
            this.props.onChange(value);
        }
    }
    
    handleBlur() {
        this.setState({ touched: true });
        if (this.props.onBlur) {
            this.props.onBlur(this.state.value);
        }
    }
    
    render() {
        const {
            type = 'text',
            label,
            placeholder,
            required = false,
            disabled = false,
            className = ''
        } = this.props;
        
        const { value, error, touched } = this.state;
        const showError = touched && error;
        
        const inputClass = `form-control ${showError ? 'is-invalid' : ''} ${className}`;
        
        return this.createElement(
            'div',
            { className: 'form-group' },
            label && this.createElement('label', {}, label),
            this.createElement('input', {
                type,
                value,
                placeholder,
                required,
                disabled,
                className: inputClass,
                onChange: (e) => this.handleChange(e),
                onBlur: () => this.handleBlur()
            }),
            showError && this.createElement('div', { className: 'invalid-feedback' }, error)
        );
    }
}

class Modal extends Component {
    constructor(props) {
        super(props);
        this.state = { isOpen: props.isOpen || false };
    }
    
    open() {
        this.setState({ isOpen: true });
        document.body.style.overflow = 'hidden';
    }
    
    close() {
        this.setState({ isOpen: false });
        document.body.style.overflow = '';
        if (this.props.onClose) {
            this.props.onClose();
        }
    }
    
    handleOverlayClick(e) {
        if (e.target === e.currentTarget) {
            this.close();
        }
    }
    
    handleEscape(e) {
        if (e.key === 'Escape' && this.state.isOpen) {
            this.close();
        }
    }
    
    onMount() {
        document.addEventListener('keydown', (e) => this.handleEscape(e));
    }
    
    onUnmount() {
        document.removeEventListener('keydown', (e) => this.handleEscape(e));
        document.body.style.overflow = '';
    }
    
    render() {
        const { title, children, size = 'md' } = this.props;
        const { isOpen } = this.state;
        
        if (!isOpen) return this.createElement('div', { style: { display: 'none' } });
        
        const modalClass = `modal-content modal-${size}`;
        
        return this.createElement(
            'div',
            {
                className: 'modal-overlay',
                onClick: (e) => this.handleOverlayClick(e)
            },
            this.createElement(
                'div',
                { className: modalClass },
                this.createElement(
                    'div',
                    { className: 'modal-header' },
                    this.createElement('h3', {}, title),
                    this.createElement('button', {
                        className: 'modal-close',
                        onClick: () => this.close()
                    }, '×')
                ),
                this.createElement(
                    'div',
                    { className: 'modal-body' },
                    children
                ),
                this.props.footer && this.createElement(
                    'div',
                    { className: 'modal-footer' },
                    this.props.footer
                )
            )
        );
    }
}

class Toast extends Component {
    constructor(props) {
        super(props);
        this.state = { visible: true };
    }
    
    onMount() {
        if (this.props.autoClose !== false) {
            setTimeout(() => this.close(), this.props.duration || 5000);
        }
    }
    
    close() {
        this.setState({ visible: false });
        if (this.props.onClose) {
            setTimeout(() => this.props.onClose(), 300);
        }
    }
    
    render() {
        const { type = 'info', message } = this.props;
        const { visible } = this.state;
        
        const toastClass = `toast toast-${type} ${visible ? 'show' : 'hide'}`;
        
        return this.createElement(
            'div',
            { className: toastClass },
            this.createElement('div', { className: 'toast-icon' }, this.getIcon(type)),
            this.createElement('div', { className: 'toast-message' }, message),
            this.createElement('button', {
                className: 'toast-close',
                onClick: () => this.close()
            }, '×')
        );
    }
    
    getIcon(type) {
        const icons = {
            success: '✓',
            error: '✗',
            warning: '⚠',
            info: 'ℹ'
        };
        return icons[type] || icons.info;
    }
}

class ToastContainer extends Component {
    constructor(props) {
        super(props);
        this.state = { toasts: [] };
    }
    
    add(toast) {
        const id = Date.now() + Math.random();
        this.setState({
            toasts: [...this.state.toasts, { id, ...toast }]
        });
    }
    
    remove(id) {
        this.setState({
            toasts: this.state.toasts.filter(t => t.id !== id)
        });
    }
    
    render() {
        const { position = 'top-right' } = this.props;
        
        return this.createElement(
            'div',
            { className: `toast-container toast-${position}` },
            ...this.state.toasts.map(toast => 
                new Toast({
                    ...toast,
                    onClose: () => this.remove(toast.id)
                })
            )
        );
    }
}

// ============================================================================
// Page Components
// ============================================================================

class LoginPage extends Component {
    constructor(props) {
        super(props);
        this.state = {
            email: '',
            password: '',
            loading: false,
            errors: {}
        };
    }
    
    async handleSubmit(e) {
        e.preventDefault();
        
        this.setState({ loading: true, errors: {} });
        
        try {
            const response = await API.login(this.state.email, this.state.password);
            
            AppState.setState({
                user: response.user,
                token: response.token
            });
            
            localStorage.setItem('token', response.token);
            router.navigate('/dashboard');
            
        } catch (error) {
            this.setState({
                errors: { form: error.message },
                loading: false
            });
        }
    }
    
    render() {
        return this.createElement(
            'div',
            { className: 'login-container' },
            this.createElement(
                'form',
                {
                    className: 'login-form',
                    onSubmit: (e) => this.handleSubmit(e)
                },
                this.createElement('h2', {}, 'Login'),
                
                this.state.errors.form && this.createElement(
                    'div',
                    { className: 'alert alert-error' },
                    this.state.errors.form
                ),
                
                new Input({
                    type: 'email',
                    label: 'Email',
                    value: this.state.email,
                    required: true,
                    onChange: (value) => this.setState({ email: value }),
                    validate: (value) => {
                        if (!value) return 'Email is required';
                        if (!/\S+@\S+\.\S+/.test(value)) return 'Invalid email format';
                        return null;
                    }
                }),
                
                new Input({
                    type: 'password',
                    label: 'Password',
                    value: this.state.password,
                    required: true,
                    onChange: (value) => this.setState({ password: value }),
                    validate: (value) => {
                        if (!value) return 'Password is required';
                        if (value.length < 6) return 'Password must be at least 6 characters';
                        return null;
                    }
                }),
                
                new Button({
                    type: 'submit',
                    variant: 'primary',
                    size: 'lg',
                    loading: this.state.loading,
                    disabled: this.state.loading
                }, 'Login'),
                
                this.createElement(
                    'p',
                    { className: 'text-center' },
                    "Don't have an account? ",
                    this.createElement(
                        'a',
                        {
                            href: '/register',
                            onClick: router.getLink('/register')
                        },
                        'Register'
                    )
                )
            )
        );
    }
}

class RegisterPage extends Component {
    constructor(props) {
        super(props);
        this.state = {
            name: '',
            email: '',
            password: '',
            confirmPassword: '',
            loading: false,
            errors: {}
        };
    }
    
    async handleSubmit(e) {
        e.preventDefault();
        
        // Validate passwords match
        if (this.state.password !== this.state.confirmPassword) {
            this.setState({ errors: { confirmPassword: 'Passwords do not match' } });
            return;
        }
        
        this.setState({ loading: true, errors: {} });
        
        try {
            const response = await API.register({
                name: this.state.name,
                email: this.state.email,
                password: this.state.password
            });
            
            AppState.setState({
                user: response.user,
                token: response.token
            });
            
            localStorage.setItem('token', response.token);
            router.navigate('/dashboard');
            
        } catch (error) {
            this.setState({
                errors: { form: error.message },
                loading: false
            });
        }
    }
    
    render() {
        return this.createElement(
            'div',
            { className: 'register-container' },
            this.createElement(
                'form',
                {
                    className: 'register-form',
                    onSubmit: (e) => this.handleSubmit(e)
                },
                this.createElement('h2', {}, 'Register'),
                
                this.state.errors.form && this.createElement(
                    'div',
                    { className: 'alert alert-error' },
                    this.state.errors.form
                ),
                
                new Input({
                    label: 'Name',
                    value: this.state.name,
                    required: true,
                    onChange: (value) => this.setState({ name: value })
                }),
                
                new Input({
                    type: 'email',
                    label: 'Email',
                    value: this.state.email,
                    required: true,
                    onChange: (value) => this.setState({ email: value })
                }),
                
                new Input({
                    type: 'password',
                    label: 'Password',
                    value: this.state.password,
                    required: true,
                    onChange: (value) => this.setState({ password: value })
                }),
                
                new Input({
                    type: 'password',
                    label: 'Confirm Password',
                    value: this.state.confirmPassword,
                    required: true,
                    onChange: (value) => this.setState({ confirmPassword: value }),
                    error: this.state.errors.confirmPassword
                }),
                
                new Button({
                    type: 'submit',
                    variant: 'primary',
                    size: 'lg',
                    loading: this.state.loading,
                    disabled: this.state.loading
                }, 'Register'),
                
                this.createElement(
                    'p',
                    { className: 'text-center' },
                    'Already have an account? ',
                    this.createElement(
                        'a',
                        {
                            href: '/login',
                            onClick: router.getLink('/login')
                        },
                        'Login'
                    )
                )
            )
        );
    }
}

class DashboardPage extends Component {
    constructor(props) {
        super(props);
        this.state = {
            stats: null,
            recentProducts: [],
            recentOrders: [],
            loading: true
        };
    }
    
    async onMount() {
        await this.loadData();
    }
    
    async loadData() {
        try {
            const [products, orders] = await Promise.all([
                API.getProducts({ limit: 5 }),
                API.getOrders({ limit: 5 })
            ]);
            
            this.setState({
                recentProducts: products.data || [],
                recentOrders: orders.data || [],
                stats: {
                    totalProducts: products.meta?.total || 0,
                    totalOrders: orders.meta?.total || 0,
                    pendingOrders: orders.data?.filter(o => o.status === 'pending').length || 0
                },
                loading: false
            });
        } catch (error) {
            console.error('Failed to load dashboard data:', error);
            this.setState({ loading: false });
        }
    }
    
    render() {
        if (this.state.loading) {
            return this.createElement('div', { className: 'loading' }, 'Loading...');
        }
        
        return this.createElement(
            'div',
            { className: 'dashboard' },
            this.createElement('h1', {}, 'Dashboard'),
            
            // Stats cards
            this.createElement(
                'div',
                { className: 'stats-grid' },
                this.createElement(
                    'div',
                    { className: 'stat-card' },
                    this.createElement('h3', {}, 'Total Products'),
                    this.createElement('p', {}, this.state.stats.totalProducts)
                ),
                this.createElement(
                    'div',
                    { className: 'stat-card' },
                    this.createElement('h3', {}, 'Total Orders'),
                    this.createElement('p', {}, this.state.stats.totalOrders)
                ),
                this.createElement(
                    'div',
                    { className: 'stat-card' },
                    this.createElement('h3', {}, 'Pending Orders'),
                    this.createElement('p', {}, this.state.stats.pendingOrders)
                )
            ),
            
            // Recent products
            this.createElement(
                'div',
                { className: 'recent-section' },
                this.createElement('h2', {}, 'Recent Products'),
                this.createElement(
                    'div',
                    { className: 'table-container' },
                    this.createElement(
                        'table',
                        { className: 'table' },
                        this.createElement(
                            'thead',
                            {},
                            this.createElement(
                                'tr',
                                {},
                                this.createElement('th', {}, 'Name'),
                                this.createElement('th', {}, 'Price'),
                                this.createElement('th', {}, 'Stock'),
                                this.createElement('th', {}, 'Actions')
                            )
                        ),
                        this.createElement(
                            'tbody',
                            {},
                            ...this.state.recentProducts.map(product =>
                                this.createElement(
                                    'tr',
                                    { key: product.id },
                                    this.createElement('td', {}, product.name),
                                    this.createElement('td', {}, `$${product.price}`),
                                    this.createElement('td', {}, product.stock),
                                    this.createElement(
                                        'td',
                                        {},
                                        this.createElement(
                                            'a',
                                            {
                                                href: `/products/${product.id}`,
                                                onClick: router.getLink(`/products/${product.id}`)
                                            },
                                            'View'
                                        )
                                    )
                                )
                            )
                        )
                    )
                )
            ),
            
            // Recent orders
            this.createElement(
                'div',
                { className: 'recent-section' },
                this.createElement('h2', {}, 'Recent Orders'),
                this.createElement(
                    'div',
                    { className: 'table-container' },
                    this.createElement(
                        'table',
                        { className: 'table' },
                        this.createElement(
                            'thead',
                            {},
                            this.createElement(
                                'tr',
                                {},
                                this.createElement('th', {}, 'Order #'),
                                this.createElement('th', {}, 'Date'),
                                this.createElement('th', {}, 'Status'),
                                this.createElement('th', {}, 'Total'),
                                this.createElement('th', {}, 'Actions')
                            )
                        ),
                        this.createElement(
                            'tbody',
                            {},
                            ...this.state.recentOrders.map(order =>
                                this.createElement(
                                    'tr',
                                    { key: order.id },
                                    this.createElement('td', {}, order.orderNumber),
                                    this.createElement('td', {}, new Date(order.createdAt).toLocaleDateString()),
                                    this.createElement(
                                        'td',
                                        {},
                                        this.createElement(
                                            'span',
                                            { className: `status-badge status-${order.status}` },
                                            order.status
                                        )
                                    ),
                                    this.createElement('td', {}, `$${order.total}`),
                                    this.createElement(
                                        'td',
                                        {},
                                        this.createElement(
                                            'a',
                                            {
                                                href: `/orders/${order.id}`,
                                                onClick: router.getLink(`/orders/${order.id}`)
                                            },
                                            'View'
                                        )
                                    )
                                )
                            )
                        )
                    )
                )
            )
        );
    }
}

class ProductsPage extends Component {
    constructor(props) {
        super(props);
        this.state = {
            products: [],
            loading: true,
            filters: {
                page: 1,
                limit: 20,
                search: '',
                category: '',
                minPrice: '',
                maxPrice: ''
            },
            pagination: null,
            showCreateModal: false
        };
    }
    
    async onMount() {
        await this.loadProducts();
    }
    
    async loadProducts() {
        this.setState({ loading: true });
        
        try {
            const response = await API.getProducts(this.state.filters);
            
            this.setState({
                products: response.data || [],
                pagination: response.meta,
                loading: false
            });
        } catch (error) {
            console.error('Failed to load products:', error);
            this.setState({ loading: false });
        }
    }
    
    handleFilterChange(key, value) {
        this.setState(
            { filters: { ...this.state.filters, [key]: value, page: 1 } },
            () => this.loadProducts()
        );
    }
    
    handlePageChange(page) {
        this.setState(
            { filters: { ...this.state.filters, page } },
            () => this.loadProducts()
        );
    }
    
    handleCreateProduct(product) {
        // Implementation would open modal and handle creation
        console.log('Create product:', product);
    }
    
    render() {
        const { products, loading, filters, pagination } = this.state;
        
        return this.createElement(
            'div',
            { className: 'products-page' },
            this.createElement('h1', {}, 'Products'),
            
            // Filters
            this.createElement(
                'div',
                { className: 'filters' },
                new Input({
                    placeholder: 'Search products...',
                    value: filters.search,
                    onChange: (value) => this.handleFilterChange('search', value)
                }),
                
                new Input({
                    placeholder: 'Category',
                    value: filters.category,
                    onChange: (value) => this.handleFilterChange('category', value)
                }),
                
                new Input({
                    type: 'number',
                    placeholder: 'Min Price',
                    value: filters.minPrice,
                    onChange: (value) => this.handleFilterChange('minPrice', value)
                }),
                
                new Input({
                    type: 'number',
                    placeholder: 'Max Price',
                    value: filters.maxPrice,
                    onChange: (value) => this.handleFilterChange('maxPrice', value)
                }),
                
                new Button({
                    variant: 'primary',
                    onClick: () => this.setState({ showCreateModal: true })
                }, 'Create Product')
            ),
            
            // Products table
            this.createElement(
                'div',
                { className: 'table-container' },
                loading ? 'Loading...' : this.createElement(
                    'table',
                    { className: 'table' },
                    this.createElement(
                        'thead',
                        {},
                        this.createElement(
                            'tr',
                            {},
                            this.createElement('th', {}, 'SKU'),
                            this.createElement('th', {}, 'Name'),
                            this.createElement('th', {}, 'Category'),
                            this.createElement('th', {}, 'Price'),
                            this.createElement('th', {}, 'Stock'),
                            this.createElement('th', {}, 'Actions')
                        )
                    ),
                    this.createElement(
                        'tbody',
                        {},
                        ...products.map(product =>
                            this.createElement(
                                'tr',
                                { key: product.id },
                                this.createElement('td', {}, product.sku),
                                this.createElement('td', {}, product.name),
                                this.createElement('td', {}, product.category),
                                this.createElement('td', {}, `$${product.price}`),
                                this.createElement(
                                    'td',
                                    { className: product.stock < 10 ? 'text-danger' : '' },
                                    product.stock
                                ),
                                this.createElement(
                                    'td',
                                    {},
                                    this.createElement(
                                        'a',
                                        {
                                            href: `/products/${product.id}`,
                                            onClick: router.getLink(`/products/${product.id}`)
                                        },
                                        'View'
                                    ),
                                    ' | ',
                                    this.createElement(
                                        'a',
                                        {
                                            href: `/products/${product.id}/edit`,
                                            onClick: router.getLink(`/products/${product.id}/edit`)
                                        },
                                        'Edit'
                                    )
                                )
                            )
                        )
                    )
                )
            ),
            
            // Pagination
            pagination && pagination.totalPages > 1 && this.createElement(
                'div',
                { className: 'pagination' },
                this.createElement(
                    'button',
                    {
                        disabled: filters.page === 1,
                        onClick: () => this.handlePageChange(filters.page - 1)
                    },
                    'Previous'
                ),
                this.createElement(
                    'span',
                    {},
                    `Page ${filters.page} of ${pagination.totalPages}`
                ),
                this.createElement(
                    'button',
                    {
                        disabled: filters.page === pagination.totalPages,
                        onClick: () => this.handlePageChange(filters.page + 1)
                    },
                    'Next'
                )
            ),
            
            // Create product modal
            this.state.showCreateModal && new Modal({
                title: 'Create Product',
                size: 'lg',
                isOpen: true,
                onClose: () => this.setState({ showCreateModal: false }),
                footer: this.createElement(
                    'div',
                    {},
                    new Button({
                        variant: 'secondary',
                        onClick: () => this.setState({ showCreateModal: false })
                    }, 'Cancel'),
                    new Button({
                        variant: 'primary',
                        onClick: () => this.handleCreateProduct()
                    }, 'Create')
                )
            }, this.createElement('p', {}, 'Product creation form would go here'))
        );
    }
}

// ============================================================================
// App Component
// ============================================================================

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            currentPage: null,
            showNav: true
        };
        
        // Initialize routes
        this.setupRoutes();
    }
    
    setupRoutes() {
        // Public routes
        router.addRoute('/', () => this.navigateTo(LoginPage));
        router.addRoute('/login', () => this.navigateTo(LoginPage));
        router.addRoute('/register', () => this.navigateTo(RegisterPage));
        
        // Protected routes (require authentication)
        router.addRoute('/dashboard', () => this.requireAuth(() => this.navigateTo(DashboardPage)));
        router.addRoute('/products', () => this.requireAuth(() => this.navigateTo(ProductsPage)));
        router.addRoute('/products/:id', (params) => {
            this.requireAuth(() => this.navigateTo(ProductDetailPage, { id: params.id }));
        });
        router.addRoute('/orders', () => this.requireAuth(() => this.navigateTo(OrdersPage)));
        router.addRoute('/orders/:id', (params) => {
            this.requireAuth(() => this.navigateTo(OrderDetailPage, { id: params.id }));
        });
        
        // 404 route
        router.addRoute('404', () => this.navigateTo(NotFoundPage));
    }
    
    requireAuth(callback) {
        if (AppState.user) {
            callback();
        } else {
            router.navigate('/login');
        }
    }
    
    navigateTo(Page, props = {}) {
        this.setState({ currentPage: new Page(props) });
    }
    
    handleLogout() {
        AppState.reset();
        router.navigate('/login');
    }
    
    render() {
        const { currentPage } = this.state;
        
        return this.createElement(
            'div',
            { className: 'app' },
            
            // Header
            AppState.user && this.createElement(
                'header',
                { className: 'header' },
                this.createElement(
                    'nav',
                    { className: 'nav' },
                    this.createElement(
                        'a',
                        {
                            href: '/dashboard',
                            onClick: router.getLink('/dashboard')
                        },
                        'Dashboard'
                    ),
                    this.createElement(
                        'a',
                        {
                            href: '/products',
                            onClick: router.getLink('/products')
                        },
                        'Products'
                    ),
                    this.createElement(
                        'a',
                        {
                            href: '/orders',
                            onClick: router.getLink('/orders')
                        },
                        'Orders'
                    ),
                    this.createElement(
                        'span',
                        { className: 'user-info' },
                        `Welcome, ${AppState.user?.name}`,
                        this.createElement(
                            'button',
                            {
                                className: 'btn-logout',
                                onClick: () => this.handleLogout()
                            },
                            'Logout'
                        )
                    )
                )
            ),
            
            // Main content
            this.createElement(
                'main',
                { className: 'main' },
                currentPage
            ),
            
            // Toast container
            new ToastContainer({ position: 'top-right' })
        );
    }
}

// ============================================================================
// Application Styles
// ============================================================================

const styles = `
    * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
    }
    
    body {
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
        line-height: 1.6;
        color: #333;
        background: #f5f5f5;
    }
    
    .app {
        min-height: 100vh;
        display: flex;
        flex-direction: column;
    }
    
    .header {
        background: #fff;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        padding: 1rem 2rem;
    }
    
    .nav {
        display: flex;
        align-items: center;
        gap: 2rem;
        max-width: 1200px;
        margin: 0 auto;
    }
    
    .nav a {
        color: #333;
        text-decoration: none;
        padding: 0.5rem 1rem;
        border-radius: 4px;
        transition: background 0.2s;
    }
    
    .nav a:hover {
        background: #f0f0f0;
    }
    
    .user-info {
        margin-left: auto;
        display: flex;
        align-items: center;
        gap: 1rem;
    }
    
    .btn-logout {
        padding: 0.5rem 1rem;
        background: #dc3545;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 0.9rem;
    }
    
    .btn-logout:hover {
        background: #c82333;
    }
    
    .main {
        flex: 1;
        max-width: 1200px;
        margin: 2rem auto;
        padding: 0 2rem;
    }
    
    .login-container,
    .register-container {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 400px;
    }
    
    .login-form,
    .register-form {
        background: #fff;
        padding: 2rem;
        border-radius: 8px;
        box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        width: 100%;
        max-width: 400px;
    }
    
    .login-form h2,
    .register-form h2 {
        margin-bottom: 2rem;
        text-align: center;
    }
    
    .form-group {
        margin-bottom: 1rem;
    }
    
    .form-group label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
    }
    
    .form-control {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 1rem;
    }
    
    .form-control:focus {
        outline: none;
        border-color: #007bff;
        box-shadow: 0 0 0 2px rgba(0,123,255,0.25);
    }
    
    .form-control.is-invalid {
        border-color: #dc3545;
    }
    
    .invalid-feedback {
        color: #dc3545;
        font-size: 0.875rem;
        margin-top: 0.25rem;
    }
    
    .btn {
        display: inline-block;
        padding: 0.75rem 1.5rem;
        border: none;
        border-radius: 4px;
        font-size: 1rem;
        cursor: pointer;
        transition: all 0.2s;
    }
    
    .btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }
    
    .btn-primary {
        background: #007bff;
        color: white;
    }
    
    .btn-primary:hover:not(:disabled) {
        background: #0056b3;
    }
    
    .btn-secondary {
        background: #6c757d;
        color: white;
    }
    
    .btn-secondary:hover:not(:disabled) {
        background: #545b62;
    }
    
    .btn-danger {
        background: #dc3545;
        color: white;
    }
    
    .btn-danger:hover:not(:disabled) {
        background: #c82333;
    }
    
    .btn-sm {
        padding: 0.25rem 0.5rem;
        font-size: 0.875rem;
    }
    
    .btn-lg {
        padding: 1rem 2rem;
        font-size: 1.125rem;
    }
    
    .btn.loading {
        position: relative;
        color: transparent;
    }
    
    .btn.loading::after {
        content: '';
        position: absolute;
        width: 1rem;
        height: 1rem;
        top: 50%;
        left: 50%;
        margin-top: -0.5rem;
        margin-left: -0.5rem;
        border: 2px solid white;
        border-top-color: transparent;
        border-radius: 50%;
        animation: spin 0.6s linear infinite;
    }
    
    @keyframes spin {
        to { transform: rotate(360deg); }
    }
    
    .alert {
        padding: 1rem;
        border-radius: 4px;
        margin-bottom: 1rem;
    }
    
    .alert-error {
        background: #f8d7da;
        border: 1px solid #f5c6cb;
        color: #721c24;
    }
    
    .alert-success {
        background: #d4edda;
        border: 1px solid #c3e6cb;
        color: #155724;
    }
    
    .table-container {
        background: #fff;
        border-radius: 8px;
        overflow-x: auto;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    
    .table {
        width: 100%;
        border-collapse: collapse;
    }
    
    .table th,
    .table td {
        padding: 1rem;
        text-align: left;
        border-bottom: 1px solid #eee;
    }
    
    .table th {
        background: #f8f9fa;
        font-weight: 600;
    }
    
    .table tr:hover {
        background: #f8f9fa;
    }
    
    .text-center {
        text-align: center;
    }
    
    .text-danger {
        color: #dc3545;
        font-weight: 600;
    }
    
    .loading {
        text-align: center;
        padding: 2rem;
        color: #666;
    }
    
    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 1.5rem;
        margin-bottom: 2rem;
    }
    
    .stat-card {
        background: #fff;
        padding: 1.5rem;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    
    .stat-card h3 {
        font-size: 0.9rem;
        color: #666;
        margin-bottom: 0.5rem;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }
    
    .stat-card p {
        font-size: 2rem;
        font-weight: 600;
        color: #333;
    }
    
    .recent-section {
        margin-bottom: 2rem;
    }
    
    .recent-section h2 {
        margin-bottom: 1rem;
        font-size: 1.25rem;
    }
    
    .filters {
        background: #fff;
        padding: 1rem;
        border-radius: 8px;
        margin-bottom: 1rem;
        display: flex;
        gap: 1rem;
        flex-wrap: wrap;
        align-items: flex-end;
    }
    
    .filters .form-group {
        flex: 1;
        min-width: 200px;
        margin-bottom: 0;
    }
    
    .pagination {
        margin-top: 2rem;
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 1rem;
    }
    
    .pagination button {
        padding: 0.5rem 1rem;
        border: 1px solid #ddd;
        background: white;
        border-radius: 4px;
        cursor: pointer;
    }
    
    .pagination button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
    
    .pagination button:hover:not(:disabled) {
        background: #f0f0f0;
    }
    
    .status-badge {
        display: inline-block;
        padding: 0.25rem 0.5rem;
        border-radius: 4px;
        font-size: 0.875rem;
        font-weight: 500;
    }
    
    .status-pending {
        background: #fff3cd;
        color: #856404;
    }
    
    .status-paid {
        background: #d4edda;
        color: #155724;
    }
    
    .status-shipped {
        background: #cce5ff;
        color: #004085;
    }
    
    .status-delivered {
        background: #d1ecf1;
        color: #0c5460;
    }
    
    .status-cancelled {
        background: #f8d7da;
        color: #721c24;
    }
    
    /* Modal styles */
    .modal-overlay {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0,0,0,0.5);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 1000;
    }
    
    .modal-content {
        background: white;
        border-radius: 8px;
        max-width: 90%;
        max-height: 90vh;
        overflow-y: auto;
        box-shadow: 0 4px 20px rgba(0,0,0,0.2);
    }
    
    .modal-sm { width: 400px; }
    .modal-md { width: 600px; }
    .modal-lg { width: 800px; }
    .modal-xl { width: 1000px; }
    
    .modal-header {
        padding: 1rem;
        border-bottom: 1px solid #eee;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
    
    .modal-header h3 {
        margin: 0;
    }
    
    .modal-close {
        background: none;
        border: none;
        font-size: 1.5rem;
        cursor: pointer;
        padding: 0 0.5rem;
    }
    
    .modal-body {
        padding: 1rem;
    }
    
    .modal-footer {
        padding: 1rem;
        border-top: 1px solid #eee;
        display: flex;
        justify-content: flex-end;
        gap: 0.5rem;
    }
    
    /* Toast styles */
    .toast-container {
        position: fixed;
        z-index: 2000;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }
    
    .toast-top-right {
        top: 1rem;
        right: 1rem;
    }
    
    .toast-top-left {
        top: 1rem;
        left: 1rem;
    }
    
    .toast-bottom-right {
        bottom: 1rem;
        right: 1rem;
    }
    
    .toast-bottom-left {
        bottom: 1rem;
        left: 1rem;
    }
    
    .toast {
        background: white;
        border-radius: 4px;
        padding: 1rem;
        display: flex;
        align-items: center;
        gap: 0.75rem;
        box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        min-width: 300px;
        max-width: 400px;
        animation: slideIn 0.3s ease;
    }
    
    .toast.show {
        opacity: 1;
        transform: translateX(0);
    }
    
    .toast.hide {
        opacity: 0;
        transform: translateX(100%);
        transition: all 0.3s ease;
    }
    
    @keyframes slideIn {
        from {
            opacity: 0;
            transform: translateX(100%);
        }
        to {
            opacity: 1;
            transform: translateX(0);
        }
    }
    
    .toast-success {
        border-left: 4px solid #28a745;
    }
    
    .toast-error {
        border-left: 4px solid #dc3545;
    }
    
    .toast-warning {
        border-left: 4px solid #ffc107;
    }
    
    .toast-info {
        border-left: 4px solid #17a2b8;
    }
    
    .toast-icon {
        font-size: 1.25rem;
        width: 24px;
        height: 24px;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    
    .toast-message {
        flex: 1;
        word-break: break-word;
    }
    
    .toast-close {
        background: none;
        border: none;
        font-size: 1.25rem;
        cursor: pointer;
        padding: 0 0.25rem;
    }
`;

// ============================================================================
// Application Initialization
// ============================================================================

// Add styles to document
const styleElement = document.createElement('style');
styleElement.textContent = styles;
document.head.appendChild(styleElement);

// Initialize app
const app = new App();
document.addEventListener('DOMContentLoaded', () => {
    const root = document.getElementById('root');
    if (!root) {
        const rootDiv = document.createElement('div');
        rootDiv.id = 'root';
        document.body.appendChild(rootDiv);
        app.mount(rootDiv);
    } else {
        app.mount(root);
    }
    
    // Handle initial route
    router.handleRoute();
});

// Export for testing
export {
    App,
    AppState,
    API,
    router,
    Component,
    Button,
    Input,
    Modal,
    Toast,
    ToastContainer,
    LoginPage,
    RegisterPage,
    DashboardPage,
    ProductsPage
};