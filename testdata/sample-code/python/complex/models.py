"""
Database Models Module

SQLAlchemy models for the application with relationships,
validation, and business logic methods.
"""

from datetime import datetime
from typing import List, Optional
from enum import Enum as PyEnum

from sqlalchemy import (
    Column, Integer, String, Float, Boolean, DateTime,
    ForeignKey, Text, Enum, Table, Index, UniqueConstraint,
    CheckConstraint, Numeric
)
from sqlalchemy.orm import relationship, declarative_base, validates
from sqlalchemy.ext.hybrid import hybrid_property
from sqlalchemy.dialects.postgresql import JSONB, UUID
import uuid

Base = declarative_base()


# ============================================================================
# Enums
# ============================================================================

class UserRole(str, PyEnum):
    """User role enumeration."""
    ADMIN = "admin"
    MANAGER = "manager"
    USER = "user"
    GUEST = "guest"


class UserStatus(str, PyEnum):
    """User status enumeration."""
    ACTIVE = "active"
    INACTIVE = "inactive"
    SUSPENDED = "suspended"
    PENDING = "pending"


class OrderStatus(str, PyEnum):
    """Order status enumeration."""
    PENDING = "pending"
    PROCESSING = "processing"
    CONFIRMED = "confirmed"
    SHIPPED = "shipped"
    DELIVERED = "delivered"
    CANCELLED = "cancelled"
    REFUNDED = "refunded"


class PaymentStatus(str, PyEnum):
    """Payment status enumeration."""
    PENDING = "pending"
    PAID = "paid"
    FAILED = "failed"
    REFUNDED = "refunded"


class PaymentMethod(str, PyEnum):
    """Payment method enumeration."""
    CREDIT_CARD = "credit_card"
    PAYPAL = "paypal"
    BANK_TRANSFER = "bank_transfer"
    CASH = "cash"


# ============================================================================
# Association Tables
# ============================================================================

# Many-to-many relationship between users and roles
user_roles = Table(
    "user_roles",
    Base.metadata,
    Column("user_id", Integer, ForeignKey("users.id"), primary_key=True),
    Column("role_id", Integer, ForeignKey("roles.id"), primary_key=True)
)

# Many-to-many relationship between products and categories
product_categories = Table(
    "product_categories",
    Base.metadata,
    Column("product_id", Integer, ForeignKey("products.id"), primary_key=True),
    Column("category_id", Integer, ForeignKey("categories.id"), primary_key=True)
)


# ============================================================================
# Models
# ============================================================================

class TimestampMixin:
    """Mixin to add created_at and updated_at timestamps."""
    created_at = Column(DateTime, default=datetime.utcnow, nullable=False)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)


class SoftDeleteMixin:
    """Mixin to add soft delete capability."""
    deleted_at = Column(DateTime, nullable=True)
    
    @property
    def is_deleted(self) -> bool:
        """Check if record is deleted."""
        return self.deleted_at is not None
    
    def soft_delete(self):
        """Soft delete the record."""
        self.deleted_at = datetime.utcnow()


class User(Base, TimestampMixin, SoftDeleteMixin):
    """User model."""
    __tablename__ = "users"
    
    id = Column(Integer, primary_key=True)
    uuid = Column(UUID(as_uuid=True), default=uuid.uuid4, unique=True, nullable=False)
    email = Column(String(255), unique=True, nullable=False, index=True)
    username = Column(String(100), unique=True, nullable=False, index=True)
    password_hash = Column(String(255), nullable=False)
    first_name = Column(String(100), nullable=False)
    last_name = Column(String(100), nullable=False)
    phone = Column(String(20), nullable=True)
    avatar_url = Column(String(500), nullable=True)
    bio = Column(Text, nullable=True)
    role = Column(Enum(UserRole), default=UserRole.USER, nullable=False)
    status = Column(Enum(UserStatus), default=UserStatus.PENDING, nullable=False)
    email_verified = Column(Boolean, default=False)
    phone_verified = Column(Boolean, default=False)
    last_login_at = Column(DateTime, nullable=True)
    login_count = Column(Integer, default=0)
    
    # Preferences
    language = Column(String(10), default="en")
    timezone = Column(String(50), default="UTC")
    theme = Column(String(20), default="light")
    notifications_enabled = Column(Boolean, default=True)
    
    # Metadata
    metadata = Column(JSONB, nullable=True, default={})
    
    # Relationships
    addresses = relationship("Address", back_populates="user", cascade="all, delete-orphan")
    orders = relationship("Order", back_populates="user")
    reviews = relationship("Review", back_populates="user")
    
    __table_args__ = (
        Index("idx_user_email_status", "email", "status"),
        Index("idx_user_created", "created_at"),
    )
    
    @hybrid_property
    def full_name(self) -> str:
        """Get user's full name."""
        return f"{self.first_name} {self.last_name}"
    
    @validates("email")
    def validate_email(self, key, email):
        """Validate email format."""
        if "@" not in email:
            raise ValueError("Invalid email format")
        return email.lower()
    
    @validates("phone")
    def validate_phone(self, key, phone):
        """Validate phone number if provided."""
        if phone and not phone.replace("+", "").replace("-", "").isdigit():
            raise ValueError("Invalid phone number format")
        return phone
    
    def to_dict(self) -> dict:
        """Convert user to dictionary (excluding sensitive data)."""
        return {
            "id": self.id,
            "uuid": str(self.uuid),
            "email": self.email,
            "username": self.username,
            "first_name": self.first_name,
            "last_name": self.last_name,
            "full_name": self.full_name,
            "phone": self.phone,
            "avatar_url": self.avatar_url,
            "bio": self.bio,
            "role": self.role.value if self.role else None,
            "status": self.status.value if self.status else None,
            "email_verified": self.email_verified,
            "phone_verified": self.phone_verified,
            "language": self.language,
            "timezone": self.timezone,
            "theme": self.theme,
            "notifications_enabled": self.notifications_enabled,
            "created_at": self.created_at.isoformat() if self.created_at else None,
            "last_login_at": self.last_login_at.isoformat() if self.last_login_at else None
        }


class Address(Base, TimestampMixin):
    """Address model."""
    __tablename__ = "addresses"
    
    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    address_type = Column(String(50), nullable=False)  # home, work, shipping, billing
    is_default = Column(Boolean, default=False)
    
    # Address fields
    street = Column(String(200), nullable=False)
    street2 = Column(String(200), nullable=True)
    city = Column(String(100), nullable=False)
    state = Column(String(100), nullable=False)
    postal_code = Column(String(20), nullable=False)
    country = Column(String(100), nullable=False)
    
    # Additional info
    phone = Column(String(20), nullable=True)
    instructions = Column(Text, nullable=True)
    
    # Relationships
    user = relationship("User", back_populates="addresses")
    orders_shipping = relationship("Order", foreign_keys="[Order.shipping_address_id]")
    orders_billing = relationship("Order", foreign_keys="[Order.billing_address_id]")
    
    __table_args__ = (
        Index("idx_address_user", "user_id"),
        Index("idx_address_country", "country", "city"),
    )
    
    @property
    def full_address(self) -> str:
        """Get full address as string."""
        parts = [self.street]
        if self.street2:
            parts.append(self.street2)
        parts.extend([self.city, self.state, self.postal_code, self.country])
        return ", ".join(parts)


class Category(Base, TimestampMixin):
    """Product category model."""
    __tablename__ = "categories"
    
    id = Column(Integer, primary_key=True)
    name = Column(String(200), unique=True, nullable=False)
    slug = Column(String(200), unique=True, nullable=False)
    description = Column(Text, nullable=True)
    parent_id = Column(Integer, ForeignKey("categories.id"), nullable=True)
    image_url = Column(String(500), nullable=True)
    is_active = Column(Boolean, default=True)
    sort_order = Column(Integer, default=0)
    
    # Relationships
    parent = relationship("Category", remote_side=[id], backref="children")
    products = relationship("Product", secondary=product_categories, back_populates="categories")
    
    __table_args__ = (
        Index("idx_category_slug", "slug"),
        Index("idx_category_parent", "parent_id"),
    )
    
    @validates("slug")
    def validate_slug(self, key, slug):
        """Validate slug format."""
        if not slug.replace("-", "").isalnum():
            raise ValueError("Slug can only contain alphanumeric characters and hyphens")
        return slug.lower()


class Product(Base, TimestampMixin, SoftDeleteMixin):
    """Product model."""
    __tablename__ = "products"
    
    id = Column(Integer, primary_key=True)
    sku = Column(String(100), unique=True, nullable=False)
    name = Column(String(500), nullable=False)
    slug = Column(String(500), unique=True, nullable=False)
    description = Column(Text, nullable=False)
    short_description = Column(String(500), nullable=True)
    price = Column(Numeric(10, 2), nullable=False)
    compare_at_price = Column(Numeric(10, 2), nullable=True)
    cost = Column(Numeric(10, 2), nullable=True)
    stock = Column(Integer, default=0, nullable=False)
    reserved_stock = Column(Integer, default=0)
    weight = Column(Float, nullable=True)
    weight_unit = Column(String(10), default="g")
    
    # Status
    is_active = Column(Boolean, default=True)
    is_featured = Column(Boolean, default=False)
    is_digital = Column(Boolean, default=False)
    
    # Media
    images = Column(JSONB, nullable=True, default=[])
    featured_image = Column(String(500), nullable=True)
    
    # SEO
    meta_title = Column(String(200), nullable=True)
    meta_description = Column(String(500), nullable=True)
    meta_keywords = Column(JSONB, nullable=True, default=[])
    
    # Statistics
    view_count = Column(Integer, default=0)
    sold_count = Column(Integer, default=0)
    rating_avg = Column(Float, default=0.0)
    rating_count = Column(Integer, default=0)
    
    # Relationships
    categories = relationship("Category", secondary=product_categories, back_populates="products")
    variants = relationship("ProductVariant", back_populates="product", cascade="all, delete-orphan")
    reviews = relationship("Review", back_populates="product")
    
    __table_args__ = (
        Index("idx_product_sku", "sku"),
        Index("idx_product_slug", "slug"),
        Index("idx_product_price", "price"),
        Index("idx_product_status", "is_active", "is_featured"),
        CheckConstraint("stock >= 0", name="check_stock_positive"),
        CheckConstraint("price > 0", name="check_price_positive"),
    )
    
    @hybrid_property
    def available_stock(self) -> int:
        """Get available stock (stock - reserved)."""
        return self.stock - self.reserved_stock
    
    @validates("sku")
    def validate_sku(self, key, sku):
        """Validate SKU format."""
        if not sku.strip():
            raise ValueError("SKU cannot be empty")
        return sku.upper()
    
    def reduce_stock(self, quantity: int):
        """Reduce stock by quantity."""
        if self.available_stock < quantity:
            raise ValueError(f"Insufficient stock. Available: {self.available_stock}, Requested: {quantity}")
        self.stock -= quantity
        self.sold_count += quantity
    
    def reserve_stock(self, quantity: int):
        """Reserve stock for pending orders."""
        if self.available_stock < quantity:
            raise ValueError(f"Insufficient stock. Available: {self.available_stock}, Requested: {quantity}")
        self.reserved_stock += quantity
    
    def release_stock(self, quantity: int):
        """Release reserved stock."""
        self.reserved_stock = max(0, self.reserved_stock - quantity)
    
    def to_dict(self) -> dict:
        """Convert product to dictionary."""
        return {
            "id": self.id,
            "sku": self.sku,
            "name": self.name,
            "slug": self.slug,
            "description": self.description,
            "short_description": self.short_description,
            "price": float(self.price),
            "compare_at_price": float(self.compare_at_price) if self.compare_at_price else None,
            "stock": self.stock,
            "available_stock": self.available_stock,
            "is_active": self.is_active,
            "is_featured": self.is_featured,
            "images": self.images,
            "featured_image": self.featured_image,
            "rating_avg": float(self.rating_avg),
            "rating_count": self.rating_count,
            "categories": [c.name for c in self.categories],
            "created_at": self.created_at.isoformat() if self.created_at else None
        }


class ProductVariant(Base, TimestampMixin):
    """Product variant model."""
    __tablename__ = "product_variants"
    
    id = Column(Integer, primary_key=True)
    product_id = Column(Integer, ForeignKey("products.id"), nullable=False)
    sku = Column(String(100), unique=True, nullable=False)
    name = Column(String(200), nullable=False)
    price = Column(Numeric(10, 2), nullable=False)
    stock = Column(Integer, default=0)
    options = Column(JSONB, nullable=False)  # e.g., {"size": "M", "color": "red"}
    images = Column(JSONB, nullable=True, default=[])
    
    # Relationships
    product = relationship("Product", back_populates="variants")
    
    __table_args__ = (
        Index("idx_variant_sku", "sku"),
        Index("idx_variant_product", "product_id"),
    )


class Order(Base, TimestampMixin):
    """Order model."""
    __tablename__ = "orders"
    
    id = Column(Integer, primary_key=True)
    order_number = Column(String(100), unique=True, nullable=False)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    status = Column(Enum(OrderStatus), default=OrderStatus.PENDING, nullable=False)
    
    # Addresses
    shipping_address_id = Column(Integer, ForeignKey("addresses.id"), nullable=True)
    billing_address_id = Column(Integer, ForeignKey("addresses.id"), nullable=True)
    
    # Totals
    subtotal = Column(Numeric(10, 2), nullable=False)
    tax = Column(Numeric(10, 2), default=0)
    shipping_cost = Column(Numeric(10, 2), default=0)
    discount = Column(Numeric(10, 2), default=0)
    total = Column(Numeric(10, 2), nullable=False)
    
    # Payment
    payment_status = Column(Enum(PaymentStatus), default=PaymentStatus.PENDING)
    payment_method = Column(Enum(PaymentMethod), nullable=True)
    payment_id = Column(String(255), nullable=True)
    payment_details = Column(JSONB, nullable=True)
    
    # Shipping
    shipping_method = Column(String(100), nullable=True)
    tracking_number = Column(String(255), nullable=True)
    estimated_delivery = Column(DateTime, nullable=True)
    
    # Notes
    customer_notes = Column(Text, nullable=True)
    admin_notes = Column(Text, nullable=True)
    
    # Metadata
    metadata = Column(JSONB, nullable=True, default={})
    
    # Relationships
    user = relationship("User", back_populates="orders")
    items = relationship("OrderItem", back_populates="order", cascade="all, delete-orphan")
    shipping_address = relationship("Address", foreign_keys=[shipping_address_id])
    billing_address = relationship("Address", foreign_keys=[billing_address_id])
    
    __table_args__ = (
        Index("idx_order_number", "order_number"),
        Index("idx_order_user", "user_id", "status"),
        Index("idx_order_created", "created_at"),
    )
    
    def calculate_total(self):
        """Calculate order total."""
        self.total = self.subtotal + self.tax + self.shipping_cost - self.discount
    
    @property
    def can_be_cancelled(self) -> bool:
        """Check if order can be cancelled."""
        return self.status in [OrderStatus.PENDING, OrderStatus.PROCESSING]
    
    @property
    def can_be_refunded(self) -> bool:
        """Check if order can be refunded."""
        return self.payment_status == PaymentStatus.PAID and self.status not in [OrderStatus.REFUNDED]


class OrderItem(Base, TimestampMixin):
    """Order item model."""
    __tablename__ = "order_items"
    
    id = Column(Integer, primary_key=True)
    order_id = Column(Integer, ForeignKey("orders.id"), nullable=False)
    product_id = Column(Integer, ForeignKey("products.id"), nullable=False)
    variant_id = Column(Integer, ForeignKey("product_variants.id"), nullable=True)
    
    quantity = Column(Integer, nullable=False)
    price = Column(Numeric(10, 2), nullable=False)
    total = Column(Numeric(10, 2), nullable=False)
    
    # Snapshot of product details at time of order
    product_snapshot = Column(JSONB, nullable=False)
    
    # Relationships
    order = relationship("Order", back_populates="items")
    product = relationship("Product")
    variant = relationship("ProductVariant")
    
    __table_args__ = (
        Index("idx_order_item_order", "order_id"),
        Index("idx_order_item_product", "product_id"),
    )
    
    @validates("quantity")
    def validate_quantity(self, key, quantity):
        """Validate quantity."""
        if quantity <= 0:
            raise ValueError("Quantity must be positive")
        return quantity


class Review(Base, TimestampMixin):
    """Product review model."""
    __tablename__ = "reviews"
    
    id = Column(Integer, primary_key=True)
    product_id = Column(Integer, ForeignKey("products.id"), nullable=False)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    order_id = Column(Integer, ForeignKey("orders.id"), nullable=True)
    
    rating = Column(Integer, nullable=False)
    title = Column(String(200), nullable=True)
    content = Column(Text, nullable=False)
    pros = Column(JSONB, nullable=True, default=[])
    cons = Column(JSONB, nullable=True, default=[])
    images = Column(JSONB, nullable=True, default=[])
    
    is_verified = Column(Boolean, default=False)
    is_approved = Column(Boolean, default=False)
    helpful_count = Column(Integer, default=0)
    unhelpful_count = Column(Integer, default=0)
    
    # Relationships
    product = relationship("Product", back_populates="reviews")
    user = relationship("User", back_populates="reviews")
    
    __table_args__ = (
        Index("idx_review_product", "product_id", "is_approved"),
        Index("idx_review_user", "user_id"),
        UniqueConstraint("product_id", "user_id", name="unique_user_product_review"),
        CheckConstraint("rating >= 1 AND rating <= 5", name="check_rating_range"),
    )
    
    @validates("rating")
    def validate_rating(self, key, rating):
        """Validate rating."""
        if not 1 <= rating <= 5:
            raise ValueError("Rating must be between 1 and 5")
        return rating