"""
pytest configuration file for the test suite.

This file contains shared fixtures, hooks, and configuration
for all tests in the project.
"""

import os
import sys
import json
import pytest
import tempfile
import shutil
import asyncio
from pathlib import Path
from typing import Dict, List, Any, Generator, AsyncGenerator
from datetime import datetime, timedelta
from unittest.mock import Mock, patch, MagicMock

# Add project root to Python path
sys.path.insert(0, str(Path(__file__).parent.parent))

# ============================================================================
# pytest Configuration
# ============================================================================

def pytest_configure(config):
    """Configure pytest with custom markers and settings."""
    # Register custom markers
    config.addinivalue_line(
        "markers",
        "slow: marks tests as slow (deselect with '-m \"not slow\"')"
    )
    config.addinivalue_line(
        "markers",
        "integration: marks tests as integration tests (require database)"
    )
    config.addinivalue_line(
        "markers",
        "unit: marks tests as unit tests (fast, isolated)"
    )
    config.addinivalue_line(
        "markers",
        "async: marks tests as async tests"
    )
    config.addinivalue_line(
        "markers",
        "benchmark: marks tests as benchmarks"
    )
    config.addinivalue_line(
        "markers",
        "smoke: marks tests as smoke tests (run first)"
    )
    config.addinivalue_line(
        "markers",
        "flaky: marks tests as flaky (allow retries)"
    )

def pytest_collection_modifyitems(config, items):
    """Modify test items after collection."""
    # Add `asyncio` marker to all async tests
    for item in items:
        if asyncio.iscoroutinefunction(item.function):
            item.add_marker(pytest.mark.asyncio)
        
        # Add `slow` marker to tests with certain names
        if any(marker in item.name for marker in ['benchmark', 'performance', 'load']):
            item.add_marker(pytest.mark.slow)

    # Sort tests: smoke tests first, then by name
    items.sort(key=lambda x: (
        0 if 'smoke' in x.keywords else 1,
        x.name
    ))

# ============================================================================
# Command Line Options
# ============================================================================

def pytest_addoption(parser):
    """Add custom command line options."""
    parser.addoption(
        "--integration",
        action="store_true",
        default=False,
        help="Run integration tests"
    )
    parser.addoption(
        "--slow",
        action="store_true",
        default=False,
        help="Run slow tests"
    )
    parser.addoption(
        "--db",
        action="store_true",
        default=False,
        help="Run database tests"
    )
    parser.addoption(
        "--api-key",
        action="store",
        default=None,
        help="API key for external service tests"
    )
    parser.addoption(
        "--test-data",
        action="store",
        default=None,
        help="Path to test data file"
    )

# ============================================================================
# Test Environment Setup/Teardown
# ============================================================================

@pytest.fixture(scope="session")
def test_env(request):
    """Set up test environment variables."""
    # Store original environment
    original_env = dict(os.environ)
    
    # Set test environment variables
    os.environ.update({
        "APP_ENV": "testing",
        "DEBUG": "true",
        "TESTING": "true",
        "DATABASE_URL": "sqlite:///:memory:",
        "REDIS_URL": "redis://localhost:6379/1",
        "API_KEY": request.config.getoption("--api-key") or "test-api-key",
        "JWT_SECRET": "test-secret-key",
        "LOG_LEVEL": "ERROR"
    })
    
    yield
    
    # Restore original environment
    os.environ.clear()
    os.environ.update(original_env)

@pytest.fixture(scope="session")
def test_dir():
    """Create a temporary directory for test files."""
    dirpath = tempfile.mkdtemp(prefix="pytest-")
    yield Path(dirpath)
    shutil.rmtree(dirpath)

# ============================================================================
# Fixtures for Common Test Data
# ============================================================================

@pytest.fixture
def sample_data() -> Dict[str, Any]:
    """Provide sample data for tests."""
    return {
        "id": 1,
        "name": "Test Item",
        "description": "This is a test item",
        "price": 19.99,
        "quantity": 100,
        "active": True,
        "tags": ["test", "sample", "example"],
        "metadata": {
            "created_at": "2024-01-01T00:00:00",
            "updated_at": "2024-01-01T00:00:00",
            "version": "1.0.0"
        },
        "nested": {
            "level1": {
                "level2": {
                    "value": "deeply nested"
                }
            }
        }
    }

@pytest.fixture
def sample_list() -> List[int]:
    """Provide a sample list of integers."""
    return [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

@pytest.fixture
def sample_string() -> str:
    """Provide a sample string."""
    return "The quick brown fox jumps over the lazy dog"

@pytest.fixture
def sample_users() -> List[Dict[str, Any]]:
    """Provide sample user data."""
    return [
        {
            "id": 1,
            "username": "alice",
            "email": "alice@example.com",
            "role": "admin",
            "active": True,
            "created_at": "2024-01-01"
        },
        {
            "id": 2,
            "username": "bob",
            "email": "bob@example.com",
            "role": "user",
            "active": True,
            "created_at": "2024-01-02"
        },
        {
            "id": 3,
            "username": "charlie",
            "email": "charlie@example.com",
            "role": "user",
            "active": False,
            "created_at": "2024-01-03"
        }
    ]

@pytest.fixture
def sample_products() -> List[Dict[str, Any]]:
    """Provide sample product data."""
    return [
        {
            "id": 101,
            "name": "Laptop",
            "price": 999.99,
            "category": "electronics",
            "in_stock": True,
            "stock_count": 50
        },
        {
            "id": 102,
            "name": "Mouse",
            "price": 29.99,
            "category": "electronics",
            "in_stock": True,
            "stock_count": 200
        },
        {
            "id": 103,
            "name": "Keyboard",
            "price": 79.99,
            "category": "electronics",
            "in_stock": False,
            "stock_count": 0
        },
        {
            "id": 104,
            "name": "Monitor",
            "price": 299.99,
            "category": "electronics",
            "in_stock": True,
            "stock_count": 30
        }
    ]

@pytest.fixture
def sample_orders(sample_users, sample_products) -> List[Dict[str, Any]]:
    """Provide sample order data."""
    return [
        {
            "id": 1001,
            "user_id": 1,
            "items": [
                {"product_id": 101, "quantity": 1, "price": 999.99},
                {"product_id": 102, "quantity": 2, "price": 29.99}
            ],
            "total": 1059.97,
            "status": "delivered",
            "created_at": "2024-01-15",
            "shipping_address": {
                "street": "123 Main St",
                "city": "Anytown",
                "zip": "12345"
            }
        },
        {
            "id": 1002,
            "user_id": 2,
            "items": [
                {"product_id": 103, "quantity": 1, "price": 79.99}
            ],
            "total": 79.99,
            "status": "pending",
            "created_at": "2024-01-16",
            "shipping_address": {
                "street": "456 Oak Ave",
                "city": "Othertown",
                "zip": "67890"
            }
        }
    ]

# ============================================================================
# Fixtures for File Operations
# ============================================================================

@pytest.fixture
def temp_file(test_dir) -> Path:
    """Create a temporary file with sample content."""
    file_path = test_dir / "test.txt"
    file_path.write_text("Hello, World!\nThis is a test file.\n")
    return file_path

@pytest.fixture
def temp_json_file(test_dir, sample_data) -> Path:
    """Create a temporary JSON file."""
    file_path = test_dir / "data.json"
    with open(file_path, 'w') as f:
        json.dump(sample_data, f, indent=2)
    return file_path

@pytest.fixture
def temp_csv_file(test_dir) -> Path:
    """Create a temporary CSV file."""
    file_path = test_dir / "data.csv"
    content = """id,name,age,city
1,Alice,30,New York
2,Bob,25,Los Angeles
3,Charlie,35,Chicago
4,Diana,28,Houston
"""
    file_path.write_text(content)
    return file_path

@pytest.fixture
def temp_image_file(test_dir) -> Path:
    """Create a temporary image file (1x1 pixel PNG)."""
    file_path = test_dir / "test.png"
    # Minimal 1x1 PNG file content
    png_data = bytes([
        0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,  # PNG signature
        0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,  # IHDR chunk
        0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,  # 1x1 image
        0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89,  # ... etc
        0x00, 0x00, 0x00, 0x0D, 0x49, 0x44, 0x41, 0x54,
        0x78, 0x9C, 0x63, 0x60, 0x00, 0x01, 0x00, 0x00, 0x05, 0x00,
        0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
        0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82
    ])
    file_path.write_bytes(png_data)
    return file_path

# ============================================================================
# Database Fixtures
# ============================================================================

@pytest.fixture
def db_connection():
    """Provide a mock database connection."""
    connection = MagicMock()
    cursor = MagicMock()
    connection.cursor.return_value = cursor
    cursor.fetchall.return_value = []
    cursor.fetchone.return_value = None
    return connection

@pytest.fixture
async def async_db_connection():
    """Provide an async mock database connection."""
    connection = MagicMock()
    cursor = MagicMock()
    connection.cursor.return_value = cursor
    cursor.fetchall.return_value = []
    cursor.fetchone.return_value = None
    return connection

@pytest.fixture
def mock_db_session():
    """Provide a mock SQLAlchemy session."""
    session = MagicMock()
    session.query.return_value.filter.return_value.all.return_value = []
    session.query.return_value.filter.return_value.first.return_value = None
    session.add.return_value = None
    session.commit.return_value = None
    session.rollback.return_value = None
    return session

# ============================================================================
# API/HTTP Fixtures
# ============================================================================

@pytest.fixture
def mock_requests():
    """Mock requests library for HTTP calls."""
    with patch('requests.get') as mock_get, \
         patch('requests.post') as mock_post, \
         patch('requests.put') as mock_put, \
         patch('requests.delete') as mock_delete:
        
        # Configure mock responses
        mock_response = MagicMock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"status": "success"}
        mock_response.text = "OK"
        
        mock_get.return_value = mock_response
        mock_post.return_value = mock_response
        mock_put.return_value = mock_response
        mock_delete.return_value = mock_response
        
        yield {
            'get': mock_get,
            'post': mock_post,
            'put': mock_put,
            'delete': mock_delete
        }

@pytest.fixture
def mock_httpx():
    """Mock httpx library for async HTTP calls."""
    with patch('httpx.AsyncClient') as mock_client:
        mock_response = MagicMock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"status": "success"}
        mock_response.text = "OK"
        
        client_instance = MagicMock()
        client_instance.get.return_value = mock_response
        client_instance.post.return_value = mock_response
        client_instance.put.return_value = mock_response
        client_instance.delete.return_value = mock_response
        client_instance.__aenter__.return_value = client_instance
        
        mock_client.return_value = client_instance
        
        yield mock_client

@pytest.fixture
def api_client():
    """Provide a test API client (FastAPI test client)."""
    try:
        from fastapi.testclient import TestClient
        from app.main import app
        
        client = TestClient(app)
        return client
    except ImportError:
        pytest.skip("FastAPI not installed")

# ============================================================================
# Cache Fixtures
# ============================================================================

@pytest.fixture
def mock_redis():
    """Mock Redis client."""
    redis_mock = MagicMock()
    redis_mock.get.return_value = None
    redis_mock.set.return_value = True
    redis_mock.delete.return_value = True
    redis_mock.exists.return_value = False
    redis_mock.expire.return_value = True
    redis_mock.ttl.return_value = -1
    return redis_mock

@pytest.fixture
def redis_client():
    """Provide a real Redis client for integration tests (requires Redis)."""
    try:
        import redis
        client = redis.Redis(
            host=os.getenv('REDIS_HOST', 'localhost'),
            port=int(os.getenv('REDIS_PORT', 6379)),
            db=int(os.getenv('REDIS_TEST_DB', 15)),
            decode_responses=True
        )
        client.ping()
        yield client
        client.flushdb()
        client.close()
    except (ImportError, redis.ConnectionError):
        pytest.skip("Redis not available")

# ============================================================================
# Async Fixtures
# ============================================================================

@pytest.fixture
def event_loop():
    """Create an event loop for async tests."""
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    yield loop
    loop.close()

@pytest.fixture
async def async_task_queue():
    """Provide an async task queue."""
    queue = asyncio.Queue()
    yield queue
    # Clean up any pending tasks
    while not queue.empty():
        try:
            await asyncio.wait_for(queue.get(), timeout=1.0)
        except asyncio.TimeoutError:
            break

@pytest.fixture
def background_tasks():
    """Track background tasks for cleanup."""
    tasks = []
    yield tasks
    # Cancel any remaining tasks
    for task in tasks:
        if not task.done():
            task.cancel()

# ============================================================================
# Time/Date Fixtures
# ============================================================================

@pytest.fixture
def mock_datetime():
    """Mock datetime module for consistent time-based tests."""
    with patch('datetime.datetime') as mock_dt:
        fixed_date = datetime(2024, 1, 15, 12, 0, 0)
        mock_dt.now.return_value = fixed_date
        mock_dt.utcnow.return_value = fixed_date
        mock_dt.side_effect = lambda *args, **kw: datetime(*args, **kw)
        yield mock_dt

@pytest.fixture
def mock_time():
    """Mock time module for consistent timing tests."""
    with patch('time.time') as mock_time:
        mock_time.return_value = 1705312800.0  # 2024-01-15 12:00:00 UTC
        yield mock_time

@pytest.fixture
def date_range():
    """Provide a range of dates for testing."""
    start = datetime(2024, 1, 1)
    end = datetime(2024, 1, 7)
    dates = [start + timedelta(days=i) for i in range(8)]
    return dates

# ============================================================================
# Mock Objects for Testing
# ============================================================================

@pytest.fixture
def mock_logger():
    """Provide a mock logger."""
    logger = MagicMock()
    logger.info = MagicMock()
    logger.debug = MagicMock()
    logger.warning = MagicMock()
    logger.error = MagicMock()
    logger.critical = MagicMock()
    return logger

@pytest.fixture
def mock_config():
    """Provide a mock configuration object."""
    class MockConfig:
        def __init__(self):
            self.debug = True
            self.testing = True
            self.database_url = "sqlite:///:memory:"
            self.redis_url = "redis://localhost:6379/1"
            self.secret_key = "test-secret-key"
            self.api_key = "test-api-key"
            
        def get(self, key, default=None):
            return getattr(self, key, default)
    
    return MockConfig()

@pytest.fixture
def mock_cache():
    """Provide a mock cache object."""
    cache = {}
    
    class MockCache:
        def get(self, key):
            return cache.get(key)
        
        def set(self, key, value, ttl=None):
            cache[key] = value
            return True
        
        def delete(self, key):
            if key in cache:
                del cache[key]
                return True
            return False
        
        def clear(self):
            cache.clear()
        
        def has(self, key):
            return key in cache
        
        def __len__(self):
            return len(cache)
    
    return MockCache()

# ============================================================================
# Fixtures for Error Cases
# ============================================================================

@pytest.fixture
def raises_timeout():
    """Create a function that raises TimeoutError after a delay."""
    async def _raises_timeout(delay=0.1):
        await asyncio.sleep(delay)
        raise TimeoutError("Operation timed out")
    return _raises_timeout

@pytest.fixture
def raises_value_error():
    """Create a function that raises ValueError."""
    def _raises_value_error(message="Invalid value"):
        raise ValueError(message)
    return _raises_value_error

@pytest.fixture
def raises_exception():
    """Create a function that raises a generic exception."""
    def _raises_exception(message="Something went wrong"):
        raise Exception(message)
    return _raises_exception

# ============================================================================
# Fixtures for Test Data Generation
# ============================================================================

@pytest.fixture
def random_string():
    """Generate random strings for testing."""
    import random
    import string
    
    def _random_string(length=10):
        return ''.join(random.choices(string.ascii_letters + string.digits, k=length))
    
    return _random_string

@pytest.fixture
def random_email(random_string):
    """Generate random email addresses."""
    def _random_email():
        return f"{random_string(8)}@{random_string(5)}.com"
    return _random_email

@pytest.fixture
def random_int():
    """Generate random integers."""
    import random
    
    def _random_int(min_val=0, max_val=100):
        return random.randint(min_val, max_val)
    
    return _random_int

# ============================================================================
# Performance Testing Fixtures
# ============================================================================

@pytest.fixture
def timer():
    """Simple timer for performance tests."""
    import time
    
    class Timer:
        def __init__(self):
            self.start_time = None
            self.end_time = None
        
        def start(self):
            self.start_time = time.time()
            return self
        
        def stop(self):
            self.end_time = time.time()
            return self
        
        @property
        def elapsed(self):
            if self.start_time and self.end_time:
                return self.end_time - self.start_time
            return 0
        
        def __enter__(self):
            self.start()
            return self
        
        def __exit__(self, *args):
            self.stop()
    
    return Timer

@pytest.fixture
def memory_profiler():
    """Simple memory profiler for tests."""
    import tracemalloc
    
    class MemoryProfiler:
        def __init__(self):
            self.snapshot1 = None
            self.snapshot2 = None
        
        def start(self):
            tracemalloc.start()
            self.snapshot1 = tracemalloc.take_snapshot()
            return self
        
        def stop(self):
            self.snapshot2 = tracemalloc.take_snapshot()
            tracemalloc.stop()
            return self
        
        @property
        def diff(self):
            if self.snapshot1 and self.snapshot2:
                return self.snapshot2.compare_to(self.snapshot1, 'lineno')
            return []
        
        def __enter__(self):
            self.start()
            return self
        
        def __exit__(self, *args):
            self.stop()
    
    return MemoryProfiler()

# ============================================================================
# Fixtures for Parameterized Testing
# ============================================================================

@pytest.fixture(params=[
    (5, 3, 8),
    (-1, 1, 0),
    (0, 0, 0),
    (100, 200, 300),
    (-5, -3, -8)
], ids=[
    "positive",
    "mixed",
    "zero",
    "large",
    "negative"
])
def add_test_cases(request):
    """Parameterized test cases for addition."""
    return request.param

@pytest.fixture(params=[
    (10, 2, 5),
    (9, 3, 3),
    (7, 2, 3.5),
    (10, 0, None)  # Division by zero case
], ids=[
    "exact",
    "integer",
    "float",
    "zero_divisor"
])
def divide_test_cases(request):
    """Parameterized test cases for division."""
    return request.param

# ============================================================================
# Fixtures for External Service Mocking
# ============================================================================

@pytest.fixture
def mock_smtp():
    """Mock SMTP server for email testing."""
    with patch('smtplib.SMTP') as mock_smtp:
        server_instance = MagicMock()
        server_instance.sendmail.return_value = {}
        server_instance.__enter__.return_value = server_instance
        mock_smtp.return_value = server_instance
        yield mock_smtp

@pytest.fixture
def mock_s3():
    """Mock S3 client for file upload testing."""
    with patch('boto3.client') as mock_boto:
        s3_instance = MagicMock()
        s3_instance.upload_file.return_value = None
        s3_instance.download_file.return_value = None
        s3_instance.generate_presigned_url.return_value = "https://mock-s3-url.com/file"
        mock_boto.return_value = s3_instance
        yield mock_boto

@pytest.fixture
def mock_stripe():
    """Mock Stripe API for payment testing."""
    with patch('stripe.Charge') as mock_charge, \
         patch('stripe.Customer') as mock_customer, \
         patch('stripe.Subscription') as mock_subscription:
        
        mock_charge.create.return_value = {'id': 'ch_mock', 'status': 'succeeded'}
        mock_customer.create.return_value = {'id': 'cus_mock'}
        mock_subscription.create.return_value = {'id': 'sub_mock'}
        
        yield {
            'charge': mock_charge,
            'customer': mock_customer,
            'subscription': mock_subscription
        }

# ============================================================================
# Fixtures for WebSocket Testing
# ============================================================================

@pytest.fixture
async def websocket_client():
    """Provide a WebSocket test client."""
    try:
        from fastapi.testclient import TestClient
        from app.main import app
        
        client = TestClient(app)
        with client.websocket_connect("/ws") as websocket:
            yield websocket
    except ImportError:
        pytest.skip("WebSocket support not available")

# ============================================================================
# Fixtures for GraphQL Testing
# ============================================================================

@pytest.fixture
def graphql_client():
    """Provide a GraphQL test client."""
    try:
        from graphene.test import Client
        from app.graphql.schema import schema
        
        client = Client(schema)
        return client
    except ImportError:
        pytest.skip("GraphQL not available")

# ============================================================================
# Cleanup Fixtures
# ============================================================================

@pytest.fixture(autouse=True)
def cleanup_files():
    """Clean up any files created during tests."""
    files_before = set(os.listdir('.'))
    yield
    files_after = set(os.listdir('.'))
    new_files = files_after - files_before
    
    # Clean up new files
    for file in new_files:
        try:
            if os.path.isfile(file):
                os.remove(file)
            elif os.path.isdir(file):
                shutil.rmtree(file)
        except (OSError, PermissionError):
            pass  # Ignore cleanup errors

@pytest.fixture(autouse=True)
def reset_mocks():
    """Reset all mocks after each test."""
    yield
    patch.stopall()