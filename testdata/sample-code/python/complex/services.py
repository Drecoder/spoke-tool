"""
Business Logic Services Module

Contains service layer classes that handle business logic,
coordinate between models, and interact with external services.
"""

from typing import List, Optional, Dict, Any
from datetime import datetime, timedelta
import logging
from functools import wraps

from sqlalchemy.orm import Session
from sqlalchemy import and_, or_
import stripe
import boto3

from .models import User, Product, Order, OrderItem, Review
from .config import settings


logger = logging.getLogger(__name__)


# ============================================================================
# Decorators
# ============================================================================

def log_service_call(func):
    """Decorator to log service method calls."""
    @wraps(func)
    def wrapper(*args, **kwargs):
        logger.info(f"Calling {func.__name__}")
        try:
            result = func(*args, **kwargs)
            logger.info(f"{func.__name__} completed successfully")
            return result
        except Exception as e:
            logger.error(f"{func.__name__} failed: {str(e)}")
            raise
    return wrapper


def transactional(func):
    """Decorator to handle database transactions."""
    @wraps(func)
    def wrapper(self, *args, **kwargs):
        try:
            result = func(self, *args, **kwargs)
            self.db.commit()
            return result
        except Exception:
            self.db.rollback()
            raise
    return wrapper


# ============================================================================
# Base Service
# ============================================================================

class BaseService:
    """Base service class with common functionality."""
    
    def __init__(self, db: Session):
        self.db = db
    
    def commit(self):
        """Commit the current transaction."""
        self.db.commit()
    
    def rollback(self):
        """Rollback the current transaction."""
        self.db.rollback()
    
    def flush(self):
        """Flush the current session."""
        self.db.flush()


# ============================================================================
# User Service
# ============================================================================

class UserService(BaseService):
    """Service for user-related operations."""
    
    @log_service_call
    def get_user(self, user_id: int) -> Optional[User]:
        """Get user by ID."""
        return self.db.query(User).filter(User.id == user_id).first()
    
    @log_service_call
    def get_user_by_email(self, email: str) -> Optional[User]:
        """Get user by email."""
        return self.db.query(User).filter(User.email == email).first()
    
    @log_service_call
    def get_user_by_username(self, username: str) -> Optional[User]:
        """Get user by username."""
        return self.db.query(User).filter(User.username == username).first()
    
    @log_service_call
    @transactional
    def create_user(self, user_data: Dict[str, Any]) -> User:
        """Create a new user."""
        user = User(**user_data)
        self.db.add(user)
        self.db.flush()
        return user
    
    @log_service_call
    @transactional
    def update_user(self, user_id: int, updates: Dict[str, Any]) -> Optional[User]:
        """Update an existing user."""
        user = self.get_user(user_id)
        if not user:
            return None
        
        for key, value in updates.items():
            if hasattr(user, key) and value is not None:
                setattr(user, key, value)
        
        self.db.flush()
        return user
    
    @log_service_call
    @transactional
    def delete_user(self, user_id: int) -> bool:
        """Soft delete a user."""
        user = self.get_user(user_id)
        if not user:
            return False
        
        user.soft_delete()
        self.db.flush()
        return True
    
    @log_service_call
    def list_users(
        self,
        skip: int = 0,
        limit: int = 100,
        role: Optional[str] = None,
        status: Optional[str] = None,
        search: Optional[str] = None
    ) -> List[User]:
        """List users with optional filters."""
        query = self.db.query(User).filter(User.deleted_at.is_(None))
        
        if role:
            query = query.filter(User.role == role)
        
        if status:
            query = query.filter(User.status == status)
        
        if search:
            query = query.filter(
                or_(
                    User.email.ilike(f"%{search}%"),
                    User.username.ilike(f"%{search}%"),
                    User.first_name.ilike(f"%{search}%"),
                    User.last_name.ilike(f"%{search}%")
                )
            )
        
        return query.offset(skip).limit(limit).all()
    
    @log_service_call
    def authenticate(self, email: str, password: str) -> Optional[User]:
        """Authenticate user (mock implementation)."""
        user = self.get_user_by_email(email)
        if not user:
            return None
        
        # In real implementation, verify password hash
        if password != "password123":  # Mock validation
            return None
        
        if user.status != "active":
            return None
        
        user.last_login_at = datetime.utcnow()
        user.login_count += 1
        self.db.commit()
        
        return user
    
    @log_service_call
    def get_user_orders(self, user_id: int) -> List[Order]:
        """Get all orders for a user."""
        return self.db.query(Order).filter(
            Order.user_id == user_id
        ).order_by(Order.created_at.desc()).all()
    
    @log_service_call
    def get_user_reviews(self, user_id: int) -> List[Review]:
        """Get all reviews by a user."""
        return self.db.query(Review).filter(
            Review.user_id == user_id
        ).order_by(Review.created_at.desc()).all()


# ============================================================================
# Product Service
# ============================================================================

class ProductService(BaseService):
    """Service for product-related operations."""
    
    @log_service_call
    def get_product(self, product_id: int) -> Optional[Product]:
        """Get product by ID."""
        return self.db.query(Product).filter(
            Product.id == product_id,
            Product.deleted_at.is_(None)
        ).first()
    
    @log_service_call
    def get_product_by_sku(self, sku: str) -> Optional[Product]:
        """Get product by SKU."""
        return self.db.query(Product).filter(
            Product.sku == sku.upper(),
            Product.deleted_at.is_(None)
        ).first()
    
    @log_service_call
    def get_product_by_slug(self, slug: str) -> Optional[Product]:
        """Get product by slug."""
        return self.db.query(Product).filter(
            Product.slug == slug,
            Product.deleted_at.is_(None)
        ).first()
    
    @log_service_call
    @transactional
    def create_product(self, product_data: Dict[str, Any]) -> Product:
        """Create a new product."""
        product = Product(**product_data)
        self.db.add(product)
        self.db.flush()
        return product
    
    @log_service_call
    @transactional
    def update_product(self, product_id: int, updates: Dict[str, Any]) -> Optional[Product]:
        """Update an existing product."""
        product = self.get_product(product_id)
        if not product:
            return None
        
        for key, value in updates.items():
            if hasattr(product, key) and value is not None:
                setattr(product, key, value)
        
        self.db.flush()
        return product
    
    @log_service_call
    @transactional
    def delete_product(self, product_id: int) -> bool:
        """Soft delete a product."""
        product = self.get_product(product_id)
        if not product:
            return False
        
        product.soft_delete()
        self.db.flush()
        return True
    
    @log_service_call
    def list_products(
        self,
        skip: int = 0,
        limit: int = 20,
        category_id: Optional[int] = None,
        min_price: Optional[float] = None,
        max_price: Optional[float] = None,
        in_stock: Optional[bool] = None,
        search: Optional[str] = None,
        sort_by: str = "created_at",
        sort_desc: bool = True
    ) -> List[Product]:
        """List products with filters and sorting."""
        query = self.db.query(Product).filter(Product.deleted_at.is_(None))
        
        if category_id:
            query = query.filter(Product.categories.any(id=category_id))
        
        if min_price is not None:
            query = query.filter(Product.price >= min_price)
        
        if max_price is not None:
            query = query.filter(Product.price <= max_price)
        
        if in_stock is not None:
            if in_stock:
                query = query.filter(Product.stock > 0)
            else:
                query = query.filter(Product.stock == 0)
        
        if search:
            query = query.filter(
                or_(
                    Product.name.ilike(f"%{search}%"),
                    Product.description.ilike(f"%{search}%")
                )
            )
        
        # Apply sorting
        sort_column = getattr(Product, sort_by, Product.created_at)
        if sort_desc:
            query = query.order_by(sort_column.desc())
        else:
            query = query.order_by(sort_column.asc())
        
        return query.offset(skip).limit(limit).all()
    
    @log_service_call
    @transactional
    def update_stock(
        self,
        product_id: int,
        quantity: int,
        operation: str = "set"
    ) -> Optional[Product]:
        """Update product stock."""
        product = self.get_product(product_id)
        if not product:
            return None
        
        if operation == "set":
            product.stock = quantity
        elif operation == "increment":
            product.stock += quantity
        elif operation == "decrement":
            if product.stock < quantity:
                raise ValueError("Insufficient stock")
            product.stock -= quantity
        else:
            raise ValueError(f"Invalid operation: {operation}")
        
        self.db.flush()
        return product
    
    @log_service_call
    @transactional
    def reserve_stock(self, product_id: int, quantity: int) -> bool:
        """Reserve stock for a pending order."""
        product = self.get_product(product_id)
        if not product:
            return False
        
        product.reserve_stock(quantity)
        self.db.flush()
        return True
    
    @log_service_call
    @transactional
    def release_stock(self, product_id: int, quantity: int) -> bool:
        """Release reserved stock."""
        product = self.get_product(product_id)
        if not product:
            return False
        
        product.release_stock(quantity)
        self.db.flush()
        return True


# ============================================================================
# Order Service
# ============================================================================

class OrderService(BaseService):
    """Service for order-related operations."""
    
    def __init__(self, db: Session):
        super().__init__(db)
        self.product_service = ProductService(db)
        self.user_service = UserService(db)
    
    @log_service_call
    def get_order(self, order_id: int) -> Optional[Order]:
        """Get order by ID."""
        return self.db.query(Order).filter(Order.id == order_id).first()
    
    @log_service_call
    def get_order_by_number(self, order_number: str) -> Optional[Order]:
        """Get order by order number."""
        return self.db.query(Order).filter(Order.order_number == order_number).first()
    
    @log_service_call
    @transactional
    def create_order(
        self,
        user_id: int,
        items: List[Dict[str, Any]],
        shipping_address_id: int,
        billing_address_id: Optional[int] = None
    ) -> Order:
        """Create a new order."""
        # Validate user exists
        user = self.user_service.get_user(user_id)
        if not user:
            raise ValueError(f"User {user_id} not found")
        
        # Calculate totals and validate stock
        subtotal = 0
        order_items = []
        
        for item in items:
            product = self.product_service.get_product(item["product_id"])
            if not product:
                raise ValueError(f"Product {item['product_id']} not found")
            
            quantity = item.get("quantity", 1)
            if product.stock < quantity:
                raise ValueError(f"Insufficient stock for product {product.id}")
            
            # Create order item
            item_total = float(product.price) * quantity
            subtotal += item_total
            
            order_items.append({
                "product_id": product.id,
                "quantity": quantity,
                "price": float(product.price),
                "total": item_total,
                "product_snapshot": {
                    "id": product.id,
                    "sku": product.sku,
                    "name": product.name,
                    "price": float(product.price)
                }
            })
            
            # Reserve stock
            product.reserve_stock(quantity)
        
        # Create order
        order = Order(
            order_number=self._generate_order_number(),
            user_id=user_id,
            subtotal=subtotal,
            tax=subtotal * 0.1,  # 10% tax
            shipping_cost=10.00,  # Flat shipping rate
            total=subtotal * 1.1 + 10.00,
            shipping_address_id=shipping_address_id,
            billing_address_id=billing_address_id or shipping_address_id
        )
        
        self.db.add(order)
        self.db.flush()
        
        # Create order items
        for item_data in order_items:
            item = OrderItem(
                order_id=order.id,
                **item_data
            )
            self.db.add(item)
        
        self.db.flush()
        return order
    
    @log_service_call
    def _generate_order_number(self) -> str:
        """Generate a unique order number."""
        from datetime import datetime
        import random
        
        timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
        random_suffix = f"{random.randint(1000, 9999)}"
        return f"ORD-{timestamp}-{random_suffix}"
    
    @log_service_call
    @transactional
    def update_order_status(self, order_id: int, status: str) -> Optional[Order]:
        """Update order status."""
        order = self.get_order(order_id)
        if not order:
            return None
        
        order.status = status
        
        if status == "cancelled":
            # Release reserved stock
            for item in order.items:
                self.product_service.release_stock(item.product_id, item.quantity)
        
        self.db.flush()
        return order
    
    @log_service_call
    @transactional
    def cancel_order(self, order_id: int) -> Optional[Order]:
        """Cancel an order."""
        order = self.get_order(order_id)
        if not order:
            return None
        
        if not order.can_be_cancelled:
            raise ValueError("Order cannot be cancelled")
        
        order.status = "cancelled"
        
        # Release reserved stock
        for item in order.items:
            self.product_service.release_stock(item.product_id, item.quantity)
        
        self.db.flush()
        return order
    
    @log_service_call
    def list_user_orders(
        self,
        user_id: int,
        skip: int = 0,
        limit: int = 20,
        status: Optional[str] = None
    ) -> List[Order]:
        """List orders for a user."""
        query = self.db.query(Order).filter(Order.user_id == user_id)
        
        if status:
            query = query.filter(Order.status == status)
        
        return query.order_by(Order.created_at.desc()).offset(skip).limit(limit).all()
    
    @log_service_call
    def get_order_stats(self, user_id: int) -> Dict[str, Any]:
        """Get order statistics for a user."""
        orders = self.db.query(Order).filter(Order.user_id == user_id).all()
        
        total_spent = sum(float(o.total) for o in orders)
        order_count = len(orders)
        
        status_counts = {}
        for order in orders:
            status_counts[order.status] = status_counts.get(order.status, 0) + 1
        
        return {
            "total_orders": order_count,
            "total_spent": total_spent,
            "average_order_value": total_spent / order_count if order_count else 0,
            "status_breakdown": status_counts
        }


# ============================================================================
# External Service Integrations
# ============================================================================

class PaymentService:
    """Service for payment processing."""
    
    def __init__(self):
        self.stripe_api_key = settings.STRIPE_API_KEY
        
        if self.stripe_api_key:
            stripe.api_key = self.stripe_api_key
    
    def process_payment(self, amount: float, currency: str, payment_method: str) -> Dict[str, Any]:
        """Process a payment (mock implementation)."""
        # In real implementation, this would call Stripe API
        return {
            "success": True,
            "payment_id": f"pay_{datetime.now().timestamp()}",
            "amount": amount,
            "currency": currency,
            "status": "succeeded"
        }
    
    def refund_payment(self, payment_id: str, amount: Optional[float] = None) -> Dict[str, Any]:
        """Refund a payment."""
        return {
            "success": True,
            "refund_id": f"ref_{datetime.now().timestamp()}",
            "payment_id": payment_id,
            "amount": amount
        }


class EmailService:
    """Service for sending emails."""
    
    def __init__(self):
        self.smtp_host = settings.SMTP_HOST
        self.smtp_port = settings.SMTP_PORT
        self.smtp_user = settings.SMTP_USER
        self.smtp_password = settings.SMTP_PASSWORD
    
    def send_welcome_email(self, email: str, name: str) -> bool:
        """Send welcome email to new user."""
        # Mock implementation
        logger.info(f"Sending welcome email to {name} at {email}")
        return True
    
    def send_order_confirmation(self, email: str, order: Order) -> bool:
        """Send order confirmation email."""
        logger.info(f"Sending order confirmation to {email} for order {order.order_number}")
        return True
    
    def send_password_reset(self, email: str, token: str) -> bool:
        """Send password reset email."""
        logger.info(f"Sending password reset email to {email} with token {token}")
        return True


class StorageService:
    """Service for file storage operations."""
    
    def __init__(self):
        self.s3_client = None
        if settings.AWS_ACCESS_KEY_ID:
            self.s3_client = boto3.client(
                "s3",
                aws_access_key_id=settings.AWS_ACCESS_KEY_ID,
                aws_secret_access_key=settings.AWS_SECRET_ACCESS_KEY,
                region_name=settings.AWS_REGION
            )
    
    def upload_file(self, file_data: bytes, filename: str, content_type: str) -> str:
        """Upload a file to storage."""
        # Mock implementation
        url = f"https://{settings.S3_BUCKET}.s3.amazonaws.com/{filename}"
        logger.info(f"Uploading file to {url}")
        return url
    
    def delete_file(self, url: str) -> bool:
        """Delete a file from storage."""
        logger.info(f"Deleting file from {url}")
        return True