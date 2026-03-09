```markdown
# Internal Contribution Guidelines

This document outlines the workflow for contributing to this internal project.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://internal-git/spoke-tool.git
cd spoke-tool

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o bin/ ./cmd/...
```

## 📋 Development Workflow

### 1. Create a Feature Branch
```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes
- Follow existing code patterns
- Add tests for new functionality
- Update documentation as needed

### 3. Run Local Checks
```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Run linter (if configured)
golangci-lint run
```

### 4. Commit Changes
```bash
git add .
git commit -m "Add brief description of changes"
```

### 5. Push and Create Pull Request
```bash
git push origin feature/your-feature-name
# Create PR in internal Git interface
```

## ✅ Code Review Checklist

- [ ] Tests pass
- [ ] New code has tests
- [ ] Documentation updated
- [ ] No breaking changes without discussion
- [ ] Follows existing patterns

## 🧪 Testing Requirements

- Unit tests for all new functions
- Integration tests for API changes
- Run full test suite before PR:
  ```bash
  ./scripts/run-tests.sh
  ```

## 📚 Documentation

- Update README.md for user-facing changes
- Add godoc comments for exported functions
- Update API.md for API changes

## 🔒 Security Notes

- No secrets in code (use environment variables)
- Validate all inputs
- Log security-relevant events

## 🏷️ Versioning

We follow [Semantic Versioning](https://semver.org/):
- MAJOR.MINOR.PATCH (e.g., 1.2.3)

## 📞 Need Help?

- Check existing issues
- Ask in the team channel
- Tag a maintainer for review

---

*Internal Use Only - Last Updated: 2024*
```

## ✅ **What this simplified version provides:**

| Section | Purpose |
|---------|---------|
| **Quick Start** | Get running in 30 seconds |
| **Workflow** | Simple step-by-step process |
| **Checklist** | What reviewers look for |
| **Testing** | Minimum test requirements |
| **Documentation** | What needs updating |
| **Security** | Basic security reminders |

