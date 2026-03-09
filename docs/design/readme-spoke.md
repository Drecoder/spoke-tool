# Readme Spoke Design

## 🎯 Purpose

The Readme Spoke automatically generates and maintains README documentation based on code and tests. It extracts examples from passing tests, summarizes APIs, and assembles well-structured documentation—all while preserving any manually written content.

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Readme Spoke                             │
│                                                               │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐   │
│  │  Extractor   │───▶│  Summarizer  │───▶│  Formatter   │   │
│  │              │    │              │    │              │   │
│  │ • Test Parse │    │ • API Sig    │    │ • Markdown   │   │
│  │ • Example    │    │ • Desc Gen   │    │ • Section    │   │
│  │ • Comment    │    │ • Param Doc  │    │ • Template   │   │
│  └──────────────┘    └──────────────┘    └──────┬───────┘   │
│                                                  │           │
│                                          ┌───────▼───────┐   │
│                                          │    Merger     │   │
│                                          │               │   │
│                                          │ • Preserve    │   │
│                                          │ • Overlay     │   │
│                                          │ • Validate    │   │
│                                          └───────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## 🔍 Component Details

### 1. Extractor

The extractor pulls documentation-worthy content from code and tests.

**Responsibilities:**
- Parse test files for working examples
- Extract function signatures
- Read existing comments and docstrings
- Identify usage patterns
- Filter out failing tests (only use passing tests)

**Extraction Sources:**

| Source | Content | Priority |
|--------|---------|----------|
| Passing Tests | Working examples | High |
| Function Signatures | API structure | High |
| Code Comments | Developer intent | Medium |
| Existing README | Preserved content | Medium |
| Error Messages | What not to do | Low |

**Example Extraction:**

```go
// From Test File
func TestCalculateDiscount(t *testing.T) {
    result := CalculateDiscount(100, true)
    assert.Equal(t, 90.0, result)
}

// Extracted Example
{
    "language": "go",
    "code": "discount := CalculateDiscount(100, true) // returns 90.0",
    "description": "Calculate 10% discount for members",
    "is_test": true
}
```

### 2. Summarizer

The summarizer creates human-readable API documentation.

**Responsibilities:**
- Generate function descriptions
- Document parameters and return values
- Explain complex logic simply
- Create usage guidelines
- Link related functions

**Summarization Flow:**

```
Function Signature
    ↓
Analyze Name and Parameters
    ↓
Read Existing Comments
    ↓
Generate Description (Gemma 2B)
    ↓
Format in Language Style
    ↓
Output Documentation
```

**Language-Specific Output:**

#### Go (Godoc)
```go
// CalculateDiscount applies member discount to price.
// Members receive 10% off, non-members pay full price.
//
// Parameters:
//   - price: original price (can be negative)
//   - isMember: true if customer is a member
//
// Returns:
//   - discounted price as float64
//
// Example:
//   discount := CalculateDiscount(100.0, true) // returns 90.0
func CalculateDiscount(price float64, isMember bool) float64
```

#### Node.js (JSDoc)
```javascript
/**
 * Calculates discount for members.
 * @param {number} price - Original price
 * @param {boolean} isMember - Whether customer is a member
 * @returns {number} Discounted price
 * @example
 * const discount = calculateDiscount(100, true); // 90
 */
function calculateDiscount(price, isMember) {
    return isMember ? price * 0.9 : price;
}
```

#### Python (Google Style)
```python
def calculate_discount(price: float, is_member: bool) -> float:
    """Calculate discount for members.
    
    Members receive 10% off, non-members pay full price.
    
    Args:
        price: Original price (can be negative)
        is_member: Whether customer is a member
    
    Returns:
        Discounted price as float
    
    Example:
        >>> calculate_discount(100, True)
        90.0
        >>> calculate_discount(100, False)
        100.0
    """
    return price * 0.9 if is_member else price
```

### 3. Formatter

The formatter assembles documentation into proper README sections.

**Responsibilities:**
- Structure markdown content
- Apply language-appropriate formatting
- Create tables for parameters
- Format code blocks
- Add badges and links

**README Sections:**

```
Title
├── Badges (CI, coverage, version)
├── Description
├── Installation
├── Quick Start
├── API Reference
├── Examples
├── Contributing
└── License
```

**Section Templates:**

#### Installation Section
```markdown
## Installation

### Go
```bash
go get github.com/user/project
```

### Node.js
```bash
npm install project-name
```

### Python
```bash
pip install project-name
```
```

#### API Reference Section
```markdown
## API Reference

### `CalculateDiscount(price, isMember)`

Calculates discount for members.

| Parameter | Type | Description |
|-----------|------|-------------|
| `price` | `float64` | Original price |
| `isMember` | `bool` | Member status |

**Returns:** `float64` - Discounted price

**Example:**
```go
discount := CalculateDiscount(100, true) // 90.0
```
```

### 4. Merger

The merger combines generated content with existing manual content.

**Responsibilities:**
- Parse existing README
- Identify manually written sections
- Preserve custom content
- Merge generated sections
- Validate links and formatting

**Merge Strategy:**

```
Existing README
    ↓
Parse Sections
    ↓
For Each Section:
    ├── If Manual → Preserve
    └── If Generated → Update
    ↓
Add Missing Sections
    ↓
Validate Formatting
    ↓
Write New README
```

**Content Classification:**

```go
type ContentType int

const (
    Manual ContentType = iota  // Preserve exactly
    Generated                   // Can update
    Template                    // Fill in
    Append                      // Add to end
)
```

## 🔄 Workflow States

```
                    ┌─────────────┐
                    │ Tests Pass  │
                    └──────┬──────┘
                           ▼
                    ┌─────────────┐
                    │  Extract    │
                    │  Examples   │
                    └──────┬──────┘
                           ▼
                    ┌─────────────┐
                    │  Summarize  │
                    │    APIs     │
                    └──────┬──────┘
                           ▼
                    ┌─────────────┐
                    │   Format    │
                    │  Sections   │
                    └──────┬──────┘
                           ▼
                    ┌─────────────┐
                    │ Read Existing│
                    │   README    │
                    └──────┬──────┘
                           ▼
                    ┌─────────────┐
                    │    Merge    │
                    └──────┬──────┘
                           ▼
                    ┌─────────────┐
                    │  Validate   │
                    └──────┬──────┘
                           ▼
                    ┌─────────────┐
                    │    Write    │
                    └─────────────┘
```

## 📝 SLM Prompts

### Example Extraction Prompt
```
Extract a working code example from this passing test.

Language: {language}
Test Code:
{testCode}

Source Function:
{functionCode}

Return ONLY a clean example showing how to use the function.
Include a brief comment explaining what it does.
Format according to {language} conventions.
```

### API Summary Prompt
```
Write clear documentation for this function.

Language: {language}
Function:
{functionCode}

Test Examples:
{examples}

Generate documentation including:
- Brief description of what it does
- Parameter descriptions
- Return value description
- One simple example

Use {docFormat} format.
Return ONLY the documentation block.
```

## 🧪 Example Walkthrough

### Input: Passing Tests
```go
// From test file
func TestCalculateDiscount(t *testing.T) {
    // Member test
    result := CalculateDiscount(100, true)
    assert.Equal(t, 90.0, result)
    
    // Non-member test
    result = CalculateDiscount(100, false)
    assert.Equal(t, 100.0, result)
    
    // Edge case
    result = CalculateDiscount(0, true)
    assert.Equal(t, 0.0, result)
}
```

### Extracted Examples
```
- CalculateDiscount(100, true)  // returns 90.0 for members
- CalculateDiscount(100, false) // returns 100.0 for non-members
- CalculateDiscount(0, true)     // returns 0.0 (edge case)
```

### Generated Documentation
````markdown
## API Reference

### `CalculateDiscount(price, isMember)`

Calculates member discount on prices.

**Parameters:**
- `price` (float64): Original price amount
- `isMember` (bool): Whether customer is a member

**Returns:** (float64) Discounted price

**Examples:**
```go
// Member gets 10% off
discount := CalculateDiscount(100, true)
fmt.Println(discount) // 90.0

// Non-member pays full price
discount := CalculateDiscount(100, false)
fmt.Println(discount) // 100.0

// Zero price works too
discount := CalculateDiscount(0, true)
fmt.Println(discount) // 0.0
```