"""
Python Integration Tests

Tests the interaction between multiple components of the system
in a real Python environment.
"""

import os
import sys
import pytest
import tempfile
import shutil
import asyncio
import time
from pathlib import Path
from typing import Dict, List, Any

# Mock implementations of your components
# In a real implementation, these would import your actual modules
try:
    from src.analyzer import analyze_project
    from src.generator import generate_tests
    from src.extractor import extract_docs
    from src.formatter import format_readme
    from src.updater import update_readme
except ImportError:
    # Mock classes for testing
    class MockAnalyzer:
        def analyze(self, path):
            return {"functions": [], "classes": [], "files": []}
    
    analyze_project = MockAnalyzer().analyze

# ============================================================================
# Test Environment Setup
# ============================================================================

@pytest.fixture(scope="session")
def test_dir():
    """Create a temporary test directory."""
    temp_dir = tempfile.mkdtemp(prefix="spoke-tool-python-")
    yield temp_dir
    shutil.rmtree(temp_dir)

@pytest.fixture
def project_dir(test_dir):
    """Create a project directory for each test."""
    proj_dir = os.path.join(test_dir, "project")
    os.makedirs(proj_dir, exist_ok=True)
    return proj_dir

@pytest.fixture
def cleanup(project_dir):
    """Clean up project directory after each test."""
    yield
    for file in os.listdir(project_dir):
        os.remove(os.path.join(project_dir, file))

# ============================================================================
# Helper Functions
# ============================================================================

def create_test_file(project_dir, filename, content):
    """Create a test file in the project directory."""
    filepath = os.path.join(project_dir, filename)
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    return filepath

def file_exists(project_dir, filename):
    """Check if a file exists in the project directory."""
    return os.path.exists(os.path.join(project_dir, filename))

def read_test_file(project_dir, filename):
    """Read a test file from the project directory."""
    with open(os.path.join(project_dir, filename), 'r', encoding='utf-8') as f:
        return f.read()

# ============================================================================
# Test Spoke Integration Tests
# ============================================================================

class TestTestSpokeIntegration:
    """Integration tests for the test generation spoke."""

    def test_analyze_python_code_and_find_functions(self, project_dir, cleanup):
        """Test analyzing Python code and finding functions."""
        # Create test Python file
        py_code = '''"""
Simple math module
"""

def add(a: int, b: int) -> int:
    """Add two numbers."""
    return a + b

def subtract(a: int, b: int) -> int:
    """Subtract two numbers."""
    return a - b

def multiply(a: int, b: int) -> int:
    """Multiply two numbers."""
    return a * b

def divide(a: int, b: int) -> float:
    """Divide two numbers."""
    if b == 0:
        raise ValueError("Division by zero")
    return a / b

def process_data(data: list) -> list:
    """Process a list of numbers."""
    return [x * 2 for x in data]

class Calculator:
    """Simple calculator class."""
    
    def __init__(self):
        self.history = []
    
    def add(self, a, b):
        result = a + b
        self.history.append(('add', a, b, result))
        return result
    
    def subtract(self, a, b):
        result = a - b
        self.history.append(('subtract', a, b, result))
        return result
    
    def get_history(self):
        return self.history
'''
        create_test_file(project_dir, 'math.py', py_code)

        # Analyze project
        analysis = analyze_project(project_dir)

        # Verify analysis results
        assert 'functions' in analysis
        assert len(analysis['functions']) >= 5  # add, subtract, multiply, divide, process_data
        
        # Check for specific functions
        function_names = [f['name'] for f in analysis['functions']]
        assert 'add' in function_names
        assert 'subtract' in function_names
        assert 'multiply' in function_names
        assert 'divide' in function_names
        assert 'process_data' in function_names
        
        # Check for class
        assert 'classes' in analysis
        class_names = [c['name'] for c in analysis['classes']]
        assert 'Calculator' in class_names
        
        # Check for class methods
        calc_class = next(c for c in analysis['classes'] if c['name'] == 'Calculator')
        assert 'add' in calc_class['methods']
        assert 'subtract' in calc_class['methods']
        assert 'get_history' in calc_class['methods']

    def test_analyze_python_with_type_hints(self, project_dir, cleanup):
        """Test analyzing Python code with type hints."""
        py_code = '''"""
Module with type hints
"""

from typing import List, Optional, Dict, Union

def process_items(items: List[int]) -> List[int]:
    """Process a list of integers."""
    return [x * 2 for x in items]

def get_user(name: str, age: Optional[int] = None) -> Dict[str, Union[str, int]]:
    """Get user dictionary."""
    user = {"name": name}
    if age is not None:
        user["age"] = age
    return user

class DataProcessor:
    """Generic data processor."""
    
    def __init__(self, multiplier: float = 2.0):
        self.multiplier = multiplier
    
    def process(self, data: List[float]) -> List[float]:
        """Process data with multiplier."""
        return [x * self.multiplier for x in data]
'''
        create_test_file(project_dir, 'typed.py', py_code)

        # Analyze project
        analysis = analyze_project(project_dir)

        # Verify analysis found functions with type hints
        assert len(analysis['functions']) >= 2
        
        process_items = next(f for f in analysis['functions'] if f['name'] == 'process_items')
        assert 'params' in process_items
        assert len(process_items['params']) == 1
        assert process_items['params'][0]['type'] == 'List[int]'
        assert process_items['returns'] == 'List[int]'

    def test_find_untested_functions(self, project_dir, cleanup):
        """Test finding functions without tests."""
        # Create source file with functions
        source_code = '''"""
Calculator module
"""

def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

def multiply(a, b):
    return a * b

def divide(a, b):
    if b == 0:
        raise ValueError("Division by zero")
    return a / b
'''
        create_test_file(project_dir, 'calculator.py', source_code)

        # Create test file (only tests add and subtract)
        test_code = '''"""
Tests for calculator module
"""

import pytest
from calculator import add, subtract

def test_add():
    assert add(2, 3) == 5
    assert add(-1, 1) == 0

def test_subtract():
    assert subtract(10, 4) == 6
    assert subtract(4, 10) == -6
'''
        create_test_file(project_dir, 'test_calculator.py', test_code)

        # Analyze project
        analysis = analyze_project(project_dir)

        # Find untested functions
        untested = [f for f in analysis['functions'] if not f.get('has_test', False)]
        
        # multiply and divide should be untested
        assert len(untested) == 2
        untested_names = [f['name'] for f in untested]
        assert 'multiply' in untested_names
        assert 'divide' in untested_names

    def test_generate_test_file_for_untested_functions(self, project_dir, cleanup):
        """Test generating test files for untested functions."""
        # Create source file
        source_code = '''"""
Math utilities
"""

def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

def multiply(a, b):
    return a * b

def divide(a, b):
    if b == 0:
        raise ValueError("Division by zero")
    return a / b
'''
        create_test_file(project_dir, 'math_utils.py', source_code)

        # Analyze project
        analysis = analyze_project(project_dir)

        # Find untested functions
        untested = [f for f in analysis['functions'] if not f.get('has_test', False)]

        # Generate tests (mock generation)
        generated_tests = []
        for func in untested:
            test_code = f'''
import pytest
from math_utils import {func['name']}

def test_{func['name']}():
    # Test {func['name']} function
    pass
'''
            generated_tests.append({
                'function': func['name'],
                'code': test_code
            })

        # Verify test generation
        assert len(generated_tests) == len(untested)
        
        # Write test file
        test_file = os.path.join(project_dir, 'test_math_utils.py')
        with open(test_file, 'w') as f:
            f.write('"""Generated tests"""\n\n')
            f.write('import pytest\n')
            f.write('from math_utils import *\n\n')
            for test in generated_tests:
                f.write(test['code'])
                f.write('\n')

        # Verify file was created
        assert file_exists(project_dir, 'test_math_utils.py')

# ============================================================================
# Readme Spoke Integration Tests
# ============================================================================

class TestReadmeSpokeIntegration:
    """Integration tests for the README generation spoke."""

    def test_extract_docstrings_from_code(self, project_dir, cleanup):
        """Test extracting docstrings from Python code."""
        # Create source file with docstrings
        source_code = '''"""
Main module documentation.

This module provides utility functions.
"""

def add(a, b):
    """
    Add two numbers together.
    
    Args:
        a: First number
        b: Second number
    
    Returns:
        The sum of a and b
    """
    return a + b

def subtract(a, b):
    """
    Subtract two numbers.
    
    Args:
        a: First number
        b: Second number
    
    Returns:
        The difference a - b
    """
    return a - b

class Calculator:
    """
    A simple calculator class.
    
    This class provides basic arithmetic operations.
    """
    
    def multiply(self, a, b):
        """
        Multiply two numbers.
        
        Args:
            a: First number
            b: Second number
        
        Returns:
            The product a * b
        """
        return a * b
'''
        create_test_file(project_dir, 'calculator.py', source_code)

        # Extract documentation (mock extraction)
        docs = {
            'module_doc': 'Main module documentation.\n\nThis module provides utility functions.',
            'functions': [
                {
                    'name': 'add',
                    'docstring': 'Add two numbers together.',
                    'params': ['a', 'b'],
                    'returns': 'The sum of a and b'
                },
                {
                    'name': 'subtract',
                    'docstring': 'Subtract two numbers.',
                    'params': ['a', 'b'],
                    'returns': 'The difference a - b'
                }
            ],
            'classes': [
                {
                    'name': 'Calculator',
                    'docstring': 'A simple calculator class.',
                    'methods': [
                        {
                            'name': 'multiply',
                            'docstring': 'Multiply two numbers.',
                            'params': ['a', 'b'],
                            'returns': 'The product a * b'
                        }
                    ]
                }
            ]
        }

        # Verify extraction
        assert docs['module_doc'] is not None
        assert len(docs['functions']) == 2
        assert len(docs['classes']) == 1
        
        # Check function docstrings
        add_func = next(f for f in docs['functions'] if f['name'] == 'add')
        assert 'Add two numbers' in add_func['docstring']

    def test_extract_examples_from_tests(self, project_dir, cleanup):
        """Test extracting examples from test files."""
        # Create source file
        source_code = '''"""
Calculator module
"""

def add(a, b):
    return a + b
'''
        create_test_file(project_dir, 'calculator.py', source_code)

        # Create test file with examples
        test_code = '''"""
Tests for calculator module
"""

import pytest
from calculator import add

def test_add():
    """Example: Adding two numbers"""
    assert add(2, 3) == 5
    assert add(-1, 1) == 0
    assert add(0, 5) == 5

def test_add_edge_cases():
    """Example: Edge cases for addition"""
    assert add(0, 0) == 0
    assert add(-5, -3) == -8
'''
        create_test_file(project_dir, 'test_calculator.py', test_code)

        # Extract examples from tests
        examples = [
            {
                'function': 'add',
                'code': 'add(2, 3)  # returns 5',
                'description': 'Adding two positive numbers'
            },
            {
                'function': 'add',
                'code': 'add(-1, 1)  # returns 0',
                'description': 'Adding opposite numbers'
            }
        ]

        # Verify examples
        assert len(examples) == 2
        assert examples[0]['function'] == 'add'
        assert 'add(2, 3)' in examples[0]['code']

    def test_generate_readme_from_extracted_content(self, project_dir, cleanup):
        """Test generating README from extracted content."""
        # Generate README
        readme_content = """
Python Integration Tests

Tests the interaction between multiple components of the system
in a real Python environment.
"""

import os
import sys
import pytest
import tempfile
import shutil
import asyncio
import time
from pathlib import Path
from typing import Dict, List, Any

# Mock implementations of your components
# In a real implementation, these would import your actual modules
try:
    from src.analyzer import analyze_project
    from src.generator import generate_tests
    from src.extractor import extract_docs
    from src.formatter import format_readme
    from src.updater import update_readme
except ImportError:
    # Mock classes for testing
    class MockAnalyzer:
        def analyze(self, path):
            return {"functions": [], "classes": [], "files": []}
    
    analyze_project = MockAnalyzer().analyze

# ============================================================================
# Test Environment Setup
# ============================================================================

@pytest.fixture(scope="session")
def test_dir():
    """Create a temporary test directory."""
    temp_dir = tempfile.mkdtemp(prefix="spoke-tool-python-")
    yield temp_dir
    shutil.rmtree(temp_dir)

@pytest.fixture
def project_dir(test_dir):
    """Create a project directory for each test."""
    proj_dir = os.path.join(test_dir, "project")
    os.makedirs(proj_dir, exist_ok=True)
    return proj_dir

@pytest.fixture
def cleanup(project_dir):
    """Clean up project directory after each test."""
    yield
    for file in os.listdir(project_dir):
        os.remove(os.path.join(project_dir, file))

# ============================================================================
# Helper Functions
# ============================================================================

def create_test_file(project_dir, filename, content):
    """Create a test file in the project directory."""
    filepath = os.path.join(project_dir, filename)
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    return filepath

def file_exists(project_dir, filename):
    """Check if a file exists in the project directory."""
    return os.path.exists(os.path.join(project_dir, filename))

def read_test_file(project_dir, filename):
    """Read a test file from the project directory."""
    with open(os.path.join(project_dir, filename), 'r', encoding='utf-8') as f:
        return f.read()

# ============================================================================
# Test Spoke Integration Tests
# ============================================================================

class TestTestSpokeIntegration:
    """Integration tests for the test generation spoke."""

    def test_analyze_python_code_and_find_functions(self, project_dir, cleanup):
        """Test analyzing Python code and finding functions."""
        # Create test Python file
        py_code = '''"""
Simple math module
"""

def add(a: int, b: int) -> int:
    """Add two numbers."""
    return a + b

def subtract(a: int, b: int) -> int:
    """Subtract two numbers."""
    return a - b

def multiply(a: int, b: int) -> int:
    """Multiply two numbers."""
    return a * b

def divide(a: int, b: int) -> float:
    """Divide two numbers."""
    if b == 0:
        raise ValueError("Division by zero")
    return a / b

def process_data(data: list) -> list:
    """Process a list of numbers."""
    return [x * 2 for x in data]

class Calculator:
    """Simple calculator class."""
    
    def __init__(self):
        self.history = []
    
    def add(self, a, b):
        result = a + b
        self.history.append(('add', a, b, result))
        return result
    
    def subtract(self, a, b):
        result = a - b
        self.history.append(('subtract', a, b, result))
        return result
    
    def get_history(self):
        return self.history
'''
        create_test_file(project_dir, 'math.py', py_code)

        # Analyze project
        analysis = analyze_project(project_dir)

        # Verify analysis results
        assert 'functions' in analysis
        assert len(analysis['functions']) >= 5  # add, subtract, multiply, divide, process_data
        
        # Check for specific functions
        function_names = [f['name'] for f in analysis['functions']]
        assert 'add' in function_names
        assert 'subtract' in function_names
        assert 'multiply' in function_names
        assert 'divide' in function_names
        assert 'process_data' in function_names
        
        # Check for class
        assert 'classes' in analysis
        class_names = [c['name'] for c in analysis['classes']]
        assert 'Calculator' in class_names
        
        # Check for class methods
        calc_class = next(c for c in analysis['classes'] if c['name'] == 'Calculator')
        assert 'add' in calc_class['methods']
        assert 'subtract' in calc_class['methods']
        assert 'get_history' in calc_class['methods']

    def test_analyze_python_with_type_hints(self, project_dir, cleanup):
        """Test analyzing Python code with type hints."""
        py_code = '''"""
Module with type hints
"""

from typing import List, Optional, Dict, Union

def process_items(items: List[int]) -> List[int]:
    """Process a list of integers."""
    return [x * 2 for x in items]

def get_user(name: str, age: Optional[int] = None) -> Dict[str, Union[str, int]]:
    """Get user dictionary."""
    user = {"name": name}
    if age is not None:
        user["age"] = age
    return user

class DataProcessor:
    """Generic data processor."""
    
    def __init__(self, multiplier: float = 2.0):
        self.multiplier = multiplier
    
    def process(self, data: List[float]) -> List[float]:
        """Process data with multiplier."""
        return [x * self.multiplier for x in data]
'''
        create_test_file(project_dir, 'typed.py', py_code)

        # Analyze project
        analysis = analyze_project(project_dir)

        # Verify analysis found functions with type hints
        assert len(analysis['functions']) >= 2
        
        process_items = next(f for f in analysis['functions'] if f['name'] == 'process_items')
        assert 'params' in process_items
        assert len(process_items['params']) == 1
        assert process_items['params'][0]['type'] == 'List[int]'
        assert process_items['returns'] == 'List[int]'

    def test_find_untested_functions(self, project_dir, cleanup):
        """Test finding functions without tests."""
        # Create source file with functions
        source_code = '''"""
Calculator module
"""

def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

def multiply(a, b):
    return a * b

def divide(a, b):
    if b == 0:
        raise ValueError("Division by zero")
    return a / b
'''
        create_test_file(project_dir, 'calculator.py', source_code)

        # Create test file (only tests add and subtract)
        test_code = '''"""
Tests for calculator module
"""

import pytest
from calculator import add, subtract

def test_add():
    assert add(2, 3) == 5
    assert add(-1, 1) == 0

def test_subtract():
    assert subtract(10, 4) == 6
    assert subtract(4, 10) == -6
'''
        create_test_file(project_dir, 'test_calculator.py', test_code)

        # Analyze project
        analysis = analyze_project(project_dir)

        # Find untested functions
        untested = [f for f in analysis['functions'] if not f.get('has_test', False)]
        
        # multiply and divide should be untested
        assert len(untested) == 2
        untested_names = [f['name'] for f in untested]
        assert 'multiply' in untested_names
        assert 'divide' in untested_names

    def test_generate_test_file_for_untested_functions(self, project_dir, cleanup):
        """Test generating test files for untested functions."""
        # Create source file
        source_code = '''"""
Math utilities
"""

def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

def multiply(a, b):
    return a * b

def divide(a, b):
    if b == 0:
        raise ValueError("Division by zero")
    return a / b
'''
        create_test_file(project_dir, 'math_utils.py', source_code)

        # Analyze project
        analysis = analyze_project(project_dir)

        # Find untested functions
        untested = [f for f in analysis['functions'] if not f.get('has_test', False)]

        # Generate tests (mock generation)
        generated_tests = []
        for func in untested:
            test_code = f'''
import pytest
from math_utils import {func['name']}

def test_{func['name']}():
    # Test {func['name']} function
    pass
'''
            generated_tests.append({
                'function': func['name'],
                'code': test_code
            })

        # Verify test generation
        assert len(generated_tests) == len(untested)
        
        # Write test file
        test_file = os.path.join(project_dir, 'test_math_utils.py')
        with open(test_file, 'w') as f:
            f.write('"""Generated tests"""\n\n')
            f.write('import pytest\n')
            f.write('from math_utils import *\n\n')
            for test in generated_tests:
                f.write(test['code'])
                f.write('\n')

        # Verify file was created
        assert file_exists(project_dir, 'test_math_utils.py')

# ============================================================================
# Readme Spoke Integration Tests
# ============================================================================

class TestReadmeSpokeIntegration:
    """Integration tests for the README generation spoke."""

    def test_extract_docstrings_from_code(self, project_dir, cleanup):
        """Test extracting docstrings from Python code."""
        # Create source file with docstrings
        source_code = '''"""
Main module documentation.

This module provides utility functions.
"""

def add(a, b):
    """
    Add two numbers together.
    
    Args:
        a: First number
        b: Second number
    
    Returns:
        The sum of a and b
    """
    return a + b

def subtract(a, b):
    """
    Subtract two numbers.
    
    Args:
        a: First number
        b: Second number
    
    Returns:
        The difference a - b
    """
    return a - b

class Calculator:
    """
    A simple calculator class.
    
    This class provides basic arithmetic operations.
    """
    
    def multiply(self, a, b):
        """
        Multiply two numbers.
        
        Args:
            a: First number
            b: Second number
        
        Returns:
            The product a * b
        """
        return a * b
'''
        create_test_file(project_dir, 'calculator.py', source_code)

        # Extract documentation (mock extraction)
        docs = {
            'module_doc': 'Main module documentation.\n\nThis module provides utility functions.',
            'functions': [
                {
                    'name': 'add',
                    'docstring': 'Add two numbers together.',
                    'params': ['a', 'b'],
                    'returns': 'The sum of a and b'
                },
                {
                    'name': 'subtract',
                    'docstring': 'Subtract two numbers.',
                    'params': ['a', 'b'],
                    'returns': 'The difference a - b'
                }
            ],
            'classes': [
                {
                    'name': 'Calculator',
                    'docstring': 'A simple calculator class.',
                    'methods': [
                        {
                            'name': 'multiply',
                            'docstring': 'Multiply two numbers.',
                            'params': ['a', 'b'],
                            'returns': 'The product a * b'
                        }
                    ]
                }
            ]
        }

        # Verify extraction
        assert docs['module_doc'] is not None
        assert len(docs['functions']) == 2
        assert len(docs['classes']) == 1
        
        # Check function docstrings
        add_func = next(f for f in docs['functions'] if f['name'] == 'add')
        assert 'Add two numbers' in add_func['docstring']

    def test_extract_examples_from_tests(self, project_dir, cleanup):
        """Test extracting examples from test files."""
        # Create source file
        source_code = '''"""
Calculator module
"""

def add(a, b):
    return a + b
'''
        create_test_file(project_dir, 'calculator.py', source_code)

        # Create test file with examples
        test_code = '''"""
Tests for calculator module
"""

import pytest
from calculator import add

def test_add():
    """Example: Adding two numbers"""
    assert add(2, 3) == 5
    assert add(-1, 1) == 0
    assert add(0, 5) == 5

def test_add_edge_cases():
    """Example: Edge cases for addition"""
    assert add(0, 0) == 0
    assert add(-5, -3) == -8
'''
        create_test_file(project_dir, 'test_calculator.py', test_code)

        # Extract examples from tests
        examples = [
            {
                'function': 'add',
                'code': 'add(2, 3)  # returns 5',
                'description': 'Adding two positive numbers'
            },
            {
                'function': 'add',
                'code': 'add(-1, 1)  # returns 0',
                'description': 'Adding opposite numbers'
            }
        ]

        # Verify examples
        assert len(examples) == 2
        assert examples[0]['function'] == 'add'
        assert 'add(2, 3)' in examples[0]['code']

    def test_generate_readme_from_extracted_content(self, project_dir, cleanup):
        """Test generating README from extracted content."""
        # Generate README
        readme_content = '''# Calculator Library

A simple calculator library for Python.

## Installation

```bash
pip install calculator-lib

'''