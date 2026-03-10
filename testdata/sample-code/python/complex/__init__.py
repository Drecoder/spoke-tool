"""
Python Complex Examples Package

A comprehensive example demonstrating a production-ready API application
with models, services, configuration, and utilities.
"""

__version__ = "1.0.0"
__author__ = "Spoke Tool Team"

from .api import app
from .config import settings
from .models import User, Product, Order, Base
from .services import UserService, ProductService, OrderService

__all__ = [
    "app",
    "settings",
    "User",
    "Product",
    "Order",
    "Base",
    "UserService",
    "ProductService",
    "OrderService",
]