"""
Hello World Module

Simple hello world examples demonstrating basic Python syntax.
"""


def greet(name: str = "World") -> str:
    """
    Return a greeting message.
    
    Args:
        name: Name to greet (defaults to "World")
    
    Returns:
        Greeting string
        
    Examples:
        >>> greet("Alice")
        'Hello, Alice!'
        >>> greet()
        'Hello, World!'
    """
    return f"Hello, {name}!"


def greet_formal(title: str, first_name: str, last_name: str) -> str:
    """
    Return a formal greeting with title and full name.
    
    Args:
        title: Title (Mr., Mrs., Dr., etc.)
        first_name: First name
        last_name: Last name
    
    Returns:
        Formal greeting string
    """
    return f"Good day, {title} {first_name} {last_name}!"


def greet_many(names: list) -> list:
    """
    Greet multiple people.
    
    Args:
        names: List of names to greet
    
    Returns:
        List of greetings
    """
    return [greet(name) for name in names]


def main() -> None:
    """Main function to demonstrate greetings."""
    print(greet())
    print(greet("Python"))
    print(greet_formal("Dr.", "Jane", "Smith"))
    
    team = ["Alice", "Bob", "Charlie"]
    for greeting in greet_many(team):
        print(greeting)


if __name__ == "__main__":
    main()