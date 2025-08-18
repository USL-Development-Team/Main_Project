# Template-Handler Contract Testing Strategy

## Overview

This document outlines a comprehensive testing strategy designed to prevent template-handler data contract violations like GitHub Issue #35, where missing fields in handler data structures cause runtime template failures.

## The Problem

GitHub Issue #35 demonstrated a critical architectural flaw:
- Handlers use anonymous structs without type safety
- Templates expect specific data fields at runtime
- No validation prevents mismatched contracts
- Failures only surface in production when users access forms

## Testing Strategy Architecture

### 1. **Startup Validation Tests** (Fail Fast)

**Purpose**: Catch template-handler mismatches before the application starts

**Implementation**: `internal/usl/handlers/startup_validation_test.go`

```go
// Run at application startup
validator := NewTemplateContractValidator()
if err := validator.ValidateAllContracts(); err != nil {
    log.Fatalf("Template validation failed: %v", err)
}
```

**Integration Points**:
- Application startup sequence
- Docker container health checks
- CI/CD deployment gates

### 2. **Template-Handler Contract Tests** (Regression Prevention)

**Purpose**: Systematically validate all handler-template combinations

**Implementation**: `internal/usl/handlers/template_contract_test.go`

**Key Features**:
- Tests all 22+ form handlers
- Validates required field presence
- Detects exact Issue #35 scenario
- Uses reflection for dynamic validation

**Example**:
```go
contractTests := []struct {
    name           string
    handlerFunc    func(w http.ResponseWriter, r *http.Request)
    templateName   string
    requiredFields []string
}{
    {
        name:         "NewTrackerForm_Contract",
        handlerFunc:  handler.NewTrackerForm,
        templateName: "tracker-new-page",
        requiredFields: []string{"Title", "CurrentPage", "Tracker", "Errors"},
    },
}
```

### 3. **Property-Based Testing** (Edge Case Discovery)

**Purpose**: Generate random data structures to find unexpected failures

**Implementation**: `internal/usl/handlers/property_based_template_test.go`

**Capabilities**:
- Generates 100+ random data combinations
- Tests field presence matrix (2^n combinations)
- Security vulnerability detection
- Performance under load testing

### 4. **Integration Testing** (End-to-End Validation)

**Purpose**: Test complete workflows in CI/CD pipeline

**Implementation**: `test/template_validation_integration_test.go`

**Test Levels**:
- Pre-commit: Fast validation (< 5 seconds)
- CI Pipeline: Comprehensive validation
- Post-deploy: Full regression suite

## CI/CD Integration Strategy

### Pre-Commit Hooks

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running template validation..."
go test -run TestCIPreCommitValidation ./test/
if [ $? -ne 0 ]; then
    echo "❌ Template validation failed - commit blocked"
    exit 1
fi
echo "✅ Template validation passed"
```

### GitHub Actions / CI Pipeline

```yaml
# .github/workflows/template-validation.yml
name: Template Validation

on: [push, pull_request]

jobs:
  template-validation:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      
      - name: Run Template Contract Tests
        run: |
          go test -v ./internal/usl/handlers/ -run TestTemplateHandlerContracts
          
      - name: Run Startup Validation
        run: |
          go test -v ./internal/usl/handlers/ -run TestStartupTemplateValidation
          
      - name: Run Property-Based Tests  
        run: |
          go test -v ./internal/usl/handlers/ -run TestPropertyBasedTemplateValidation
          
      - name: Integration Test Suite
        run: |
          go test -v ./test/ -run TestFullTemplateValidationSuite

  deploy-gate:
    needs: template-validation
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - name: Deploy Application
        run: |
          echo "✅ Template validation passed - deploying..."
          # Deployment commands here
```

### Deployment Health Checks

```go
// cmd/server/main.go
func main() {
    // Validate templates before starting server
    if err := validateTemplatesAtStartup(); err != nil {
        log.Fatalf("Template validation failed: %v", err)
    }
    
    // Start server only after validation passes
    startServer()
}

func validateTemplatesAtStartup() error {
    validator := handlers.NewTemplateContractValidator()
    return validator.ValidateAllContracts()
}
```

## Implementation Checklist

### Phase 1: Immediate Protection (1-2 days)

- [ ] Create `startup_validation_test.go`
- [ ] Add template contract definitions for critical forms
- [ ] Integrate startup validation into application boot
- [ ] Add pre-commit hook for fast validation

### Phase 2: Comprehensive Coverage (1 week)

- [ ] Implement `template_contract_test.go` for all 22+ handlers
- [ ] Create property-based testing suite
- [ ] Add CI/CD pipeline integration
- [ ] Document all template-handler contracts

### Phase 3: Advanced Validation (2 weeks)

- [ ] Add security testing for XSS/injection
- [ ] Implement performance baseline testing
- [ ] Create automated contract generation
- [ ] Add monitoring and alerting

## Test Categories and When to Run Them

| Test Type | Speed | When to Run | Purpose |
|-----------|-------|-------------|---------|
| Startup Validation | Fast (< 1s) | Every boot | Fail fast on broken contracts |
| Contract Tests | Medium (5-10s) | Pre-commit, CI | Validate all handlers |
| Property-Based | Slow (30s+) | CI, Nightly | Find edge cases |
| Integration | Medium (10s) | CI, Deploy | End-to-end validation |

## Monitoring and Alerting

### Production Monitoring

```go
// Add to your monitoring setup
func (h *MigrationHandler) renderTemplate(w http.ResponseWriter, templateName TemplateName, data any) {
    start := time.Now()
    
    err := h.templates.ExecuteTemplate(w, string(templateName), data)
    
    // Monitor template rendering failures
    if err != nil {
        metrics.TemplateRenderingErrors.WithLabelValues(string(templateName)).Inc()
        logger.Error("Template rendering failed", 
            "template", templateName,
            "error", err,
            "data_type", reflect.TypeOf(data))
    }
    
    metrics.TemplateRenderingDuration.WithLabelValues(string(templateName)).Observe(time.Since(start).Seconds())
}
```

### Alerting Rules

```yaml
# prometheus/alerts.yml
groups:
  - name: template-validation
    rules:
      - alert: TemplateRenderingFailures
        expr: rate(template_rendering_errors_total[5m]) > 0.01
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Template rendering failures detected"
          description: "Template {{ $labels.template }} failing at {{ $value }} errors/sec"
```

## Benefits of This Strategy

### Immediate Benefits
- **Prevents Issue #35 type bugs** from reaching production
- **Faster development** with early error detection
- **Higher confidence** in template changes

### Long-term Benefits  
- **Systematic contract validation** across all forms
- **Automated regression prevention** 
- **Performance monitoring** of template rendering
- **Security vulnerability detection**

### Developer Experience
- **Clear error messages** when contracts are broken
- **Fast feedback loop** via pre-commit hooks
- **Comprehensive documentation** of template requirements

## Migration Strategy for Existing Code

### Step 1: Audit Current State
```bash
# Find all handlers using anonymous structs
grep -r "struct {" internal/usl/handlers/ 

# Find all template files
find templates/ -name "*.html"
```

### Step 2: Create Template Contracts
For each template, document expected data structure:

```go
// Document what tracker-new.html expects
type TrackerNewPageData struct {
    Title       string                `json:"title"`
    CurrentPage string                `json:"current_page"` 
    Tracker     *usl.USLUserTracker   `json:"tracker"`      // This was missing!
    Errors      map[string]string     `json:"errors"`
}
```

### Step 3: Convert Anonymous Structs
Replace anonymous structs with typed view models:

```go
// Before (vulnerable to Issue #35)
data := struct {
    Title       string
    CurrentPage string  
    Errors      map[string]string
}{...}

// After (type-safe)
data := TrackerNewPageData{
    Title:       "New Tracker",
    CurrentPage: "trackers",
    Tracker:     &usl.USLUserTracker{}, // Explicit field
    Errors:      make(map[string]string),
}
```

## Troubleshooting Common Issues

### "Template field not found" Errors
1. Check template contract definition
2. Verify handler provides all required fields
3. Run `TestTemplateHandlerContracts` for specific handler

### Slow Test Performance
1. Use `TestCIPreCommitValidation` for fast feedback
2. Run full suite only in CI
3. Consider parallel test execution

### False Positives in Property-Based Testing
1. Adjust failure rate thresholds
2. Add more specific expected failure patterns
3. Review generated data validity

## Future Enhancements

### Automated Contract Generation
- Parse templates to extract field requirements
- Generate TypeScript interfaces for frontend
- Auto-update documentation

### Advanced Security Testing
- SQL injection pattern detection
- CSRF token validation
- Input sanitization verification

### Performance Optimization
- Template compilation caching
- Lazy loading of large test suites
- Parallel test execution

---

**Remember**: The goal is not just to fix Issue #35, but to create a systematic approach that prevents an entire class of template-handler contract violations. This testing strategy transforms runtime failures into compile-time or startup-time failures, dramatically improving system reliability.