"""
Utility Functions Module

Common utility functions used throughout the application.
"""

import re
import hashlib
import json
import logging
from datetime import datetime, timedelta
from typing import Any, Dict, List, Optional, Union
from functools import wraps
from time import time
import jwt

from .config import settings


logger = logging.getLogger(__name__)


# ============================================================================
# Timing and Performance
# ============================================================================

def timer(func):
    """Decorator to measure function execution time."""
    @wraps(func)
    def wrapper(*args, **kwargs):
        start = time()
        result = func(*args, **kwargs)
        end = time()
        logger.debug(f"{func.__name__} took {end - start:.4f} seconds")
        return result
    return wrapper


def rate_limit(max_calls: int, period: int):
    """Decorator to rate limit function calls."""
    def decorator(func):
        calls = []
        
        @wraps(func)
        def wrapper(*args, **kwargs):
            now = time()
            # Remove old calls
            while calls and calls[0] < now - period:
                calls.pop(0)
            
            if len(calls) >= max_calls:
                raise Exception(f"Rate limit exceeded. Max {max_calls} calls per {period} seconds")
            
            calls.append(now)
            return func(*args, **kwargs)
        return wrapper
    return decorator


# ============================================================================
# Validation Utilities
# ============================================================================

def validate_email(email: str) -> bool:
    """Validate email format."""
    pattern = r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$"
    return bool(re.match(pattern, email))


def validate_phone(phone: str) -> bool:
    """Validate phone number format."""
    # Simple validation - only digits, +, -, spaces
    cleaned = re.sub(r"[\s\-\(\)]", "", phone)
    return cleaned.startswith("+") and cleaned[1:].isdigit() or cleaned.isdigit()


def validate_url(url: str) -> bool:
    """Validate URL format."""
    pattern = r"^https?://(?:[-\w.]|(?:%[\da-fA-F]{2}))+[^\s]*$"
    return bool(re.match(pattern, url))


def validate_password(password: str, min_length: int = 8) -> tuple[bool, str]:
    """
    Validate password strength.
    Returns (is_valid, message)
    """
    if len(password) < min_length:
        return False, f"Password must be at least {min_length} characters long"
    
    if not re.search(r"[A-Z]", password):
        return False, "Password must contain at least one uppercase letter"
    
    if not re.search(r"[a-z]", password):
        return False, "Password must contain at least one lowercase letter"
    
    if not re.search(r"\d", password):
        return False, "Password must contain at least one number"
    
    if not re.search(r"[!@#$%^&*(),.?\":{}|<>]", password):
        return False, "Password must contain at least one special character"
    
    return True, "Password is valid"


# ============================================================================
# Hashing and Security
# ============================================================================

def hash_password(password: str) -> str:
    """Hash a password (simple implementation - use bcrypt in production)."""
    return hashlib.sha256(password.encode()).hexdigest()


def verify_password(password: str, hashed: str) -> bool:
    """Verify a password against its hash."""
    return hash_password(password) == hashed


def generate_token(data: Dict[str, Any], expires_in: Optional[int] = None) -> str:
    """Generate a JWT token."""
    payload = data.copy()
    payload["iat"] = datetime.utcnow()
    
    if expires_in:
        payload["exp"] = datetime.utcnow() + timedelta(seconds=expires_in)
    else:
        payload["exp"] = datetime.utcnow() + timedelta(hours=24)
    
    return jwt.encode(payload, settings.SECRET_KEY, algorithm=settings.JWT_ALGORITHM)


def verify_token(token: str) -> Optional[Dict[str, Any]]:
    """Verify a JWT token."""
    try:
        payload = jwt.decode(
            token,
            settings.SECRET_KEY,
            algorithms=[settings.JWT_ALGORITHM]
        )
        return payload
    except jwt.ExpiredSignatureError:
        logger.warning("Token has expired")
        return None
    except jwt.InvalidTokenError as e:
        logger.warning(f"Invalid token: {e}")
        return None


# ============================================================================
# JSON Utilities
# ============================================================================

class DateTimeEncoder(json.JSONEncoder):
    """JSON encoder that handles datetime objects."""
    
    def default(self, obj):
        if isinstance(obj, datetime):
            return obj.isoformat()
        return super().default(obj)


def to_json(obj: Any, pretty: bool = False) -> str:
    """Convert object to JSON string."""
    if pretty:
        return json.dumps(obj, cls=DateTimeEncoder, indent=2, ensure_ascii=False)
    return json.dumps(obj, cls=DateTimeEncoder, ensure_ascii=False)


def from_json(json_str: str) -> Any:
    """Parse JSON string to Python object."""
    return json.loads(json_str)


# ============================================================================
# String Utilities
# ============================================================================

def slugify(text: str) -> str:
    """Convert text to URL-friendly slug."""
    # Convert to lowercase and replace spaces with hyphens
    text = text.lower().strip()
    text = re.sub(r"[^\w\s-]", "", text)
    text = re.sub(r"[\s_-]+", "-", text)
    return text.strip("-")


def truncate(text: str, max_length: int, suffix: str = "...") -> str:
    """Truncate text to maximum length."""
    if len(text) <= max_length:
        return text
    return text[:max_length - len(suffix)] + suffix


def camel_to_snake(name: str) -> str:
    """Convert camelCase to snake_case."""
    pattern = re.compile(r"(?<!^)(?=[A-Z])")
    return pattern.sub("_", name).lower()


def snake_to_camel(name: str, uppercase_first: bool = False) -> str:
    """Convert snake_case to camelCase."""
    components = name.split("_")
    if uppercase_first:
        return "".join(x.title() for x in components)
    return components[0] + "".join(x.title() for x in components[1:])


# ============================================================================
# Date and Time Utilities
# ============================================================================

def now() -> datetime:
    """Get current UTC datetime."""
    return datetime.utcnow()


def format_date(date: datetime, format: str = "%Y-%m-%d") -> str:
    """Format date as string."""
    return date.strftime(format)


def format_datetime(dt: datetime, format: str = "%Y-%m-%d %H:%M:%S") -> str:
    """Format datetime as string."""
    return dt.strftime(format)


def parse_date(date_str: str, format: str = "%Y-%m-%d") -> Optional[datetime]:
    """Parse date from string."""
    try:
        return datetime.strptime(date_str, format)
    except ValueError:
        return None


def time_ago(dt: datetime) -> str:
    """Get human-readable time difference."""
    now = datetime.utcnow()
    diff = now - dt
    
    if diff.days > 365:
        years = diff.days // 365
        return f"{years} year{'s' if years > 1 else ''} ago"
    if diff.days > 30:
        months = diff.days // 30
        return f"{months} month{'s' if months > 1 else ''} ago"
    if diff.days > 0:
        return f"{diff.days} day{'s' if diff.days > 1 else ''} ago"
    if diff.seconds > 3600:
        hours = diff.seconds // 3600
        return f"{hours} hour{'s' if hours > 1 else ''} ago"
    if diff.seconds > 60:
        minutes = diff.seconds // 60
        return f"{minutes} minute{'s' if minutes > 1 else ''} ago"
    return f"{diff.seconds} second{'s' if diff.seconds > 1 else ''} ago"


# ============================================================================
# Pagination Utilities
# ============================================================================

def paginate(items: List[Any], page: int = 1, per_page: int = 20) -> Dict[str, Any]:
    """Paginate a list of items."""
    total = len(items)
    start = (page - 1) * per_page
    end = start + per_page
    
    return {
        "items": items[start:end],
        "page": page,
        "per_page": per_page,
        "total": total,
        "total_pages": (total + per_page - 1) // per_page,
        "has_prev": page > 1,
        "has_next": end < total
    }


def paginate_query(query, page: int = 1, per_page: int = 20):
    """Paginate a SQLAlchemy query."""
    total = query.count()
    items = query.offset((page - 1) * per_page).limit(per_page).all()
    
    return {
        "items": items,
        "page": page,
        "per_page": per_page,
        "total": total,
        "total_pages": (total + per_page - 1) // per_page
    }


# ============================================================================
# Retry Utilities
# ============================================================================

def retry(max_attempts: int = 3, delay: float = 1.0, backoff: float = 2.0):
    """Decorator to retry a function on failure."""
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            current_delay = delay
            last_exception = None
            
            for attempt in range(max_attempts):
                try:
                    return func(*args, **kwargs)
                except Exception as e:
                    last_exception = e
                    if attempt < max_attempts - 1:
                        logger.warning(
                            f"Attempt {attempt + 1}/{max_attempts} failed: {e}. "
                            f"Retrying in {current_delay}s"
                        )
                        time.sleep(current_delay)
                        current_delay *= backoff
            
            raise last_exception
        return wrapper
    return decorator


# ============================================================================
# Cache Utilities
# ============================================================================

class TTLCache:
    """Simple time-to-live cache."""
    
    def __init__(self, ttl: int = 300):
        self.cache = {}
        self.ttl = ttl
    
    def get(self, key: str) -> Optional[Any]:
        """Get value from cache."""
        if key in self.cache:
            value, timestamp = self.cache[key]
            if time() - timestamp < self.ttl:
                return value
            del self.cache[key]
        return None
    
    def set(self, key: str, value: Any):
        """Set value in cache."""
        self.cache[key] = (value, time())
    
    def delete(self, key: str):
        """Delete value from cache."""
        if key in self.cache:
            del self.cache[key]
    
    def clear(self):
        """Clear all cache."""
        self.cache.clear()


cache = TTLCache()


def cached(ttl: int = 300):
    """Decorator to cache function results."""
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            # Create cache key from function name and arguments
            key = f"{func.__name__}:{str(args)}:{str(kwargs)}"
            
            # Try to get from cache
            result = cache.get(key)
            if result is not None:
                logger.debug(f"Cache hit for {key}")
                return result
            
            # Call function and cache result
            logger.debug(f"Cache miss for {key}")
            result = func(*args, **kwargs)
            cache.set(key, result)
            return result
        return wrapper
    return decorator


# ============================================================================
# Logging Utilities
# ============================================================================

def setup_logging(level: str = "INFO", log_file: Optional[str] = None):
    """Setup logging configuration."""
    handlers = [logging.StreamHandler()]
    
    if log_file:
        handlers.append(logging.FileHandler(log_file))
    
    logging.basicConfig(
        level=getattr(logging, level.upper()),
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        handlers=handlers
    )


class LoggerMixin:
    """Mixin to add logger to classes."""
    
    @property
    def logger(self):
        if not hasattr(self, "_logger"):
            self._logger = logging.getLogger(
                f"{self.__class__.__module__}.{self.__class__.__name__}"
            )
        return self._logger