# Phase 0: Testing - Quick Start Guide

**Duration:** Weeks 1-6
**Goal:** Production-ready test coverage (80%+) before open source release
**Related:** [open-source-strategy.md](./open-source-strategy.md) | [test-implementation-examples.md](./test-implementation-examples.md)

## Week 1: Day-by-Day Plan

### Monday: Setup and Agent Tests (Critical Path 1/2)

**Morning: Project Setup**
```bash
# 1. Install testify
go get github.com/stretchr/testify

# 2. Create test structure
mkdir -p internal/agent/testdata
mkdir -p internal/deploy/testdata
mkdir -p internal/docs/testdata
mkdir -p internal/config/testdata
mkdir -p internal/discovery/testdata
mkdir -p test/e2e

# 3. Create Makefile
cat > Makefile << 'EOF'
.PHONY: test test-cover
test:
	go test -v ./...
test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total
EOF
```

**Afternoon: Agent Test Data**
```bash
cd internal/agent/testdata

# Create valid-agent.md (see test-implementation-examples.md)
cat > valid-agent.md << 'EOF'
---
name: test-agent
version: "1.0.0"
description: A test agent for unit tests
---

This is the agent content.
EOF

# Create invalid-frontmatter.md
# Create missing-version.md
# Create no-frontmatter.md
# Create empty-content.md
# (See test-implementation-examples.md for full files)
```

**Evening: Write Agent Tests**
```bash
cd internal/agent

# Copy agent_test.go from test-implementation-examples.md
# Target: 20 tests, 95% coverage

go test -v
go test -cover
```

**Goal:** End of day with agent package at 95% coverage

---

### Tuesday: Deploy Tests (Critical Path 2/2)

**Morning: Deploy Test Data**
```bash
cd internal/deploy/testdata

# Create setup.go helper (see test-implementation-examples.md)
# Functions: CreateTestAgent, CreateProjectWithAgents
```

**Afternoon: Write Deploy Tests**
```bash
cd internal/deploy

# Copy deploy_test.go from test-implementation-examples.md
# Target: 15 tests, 90% coverage

go test -v
go test -cover
```

**Evening: Integration Test**
```bash
# Create integration_test.go
# Test full deployment workflow with real agent files

go test -v -tags=integration
```

**Goal:** End of day with deploy package at 90% coverage

---

### Wednesday: CI/CD Setup

**Morning: GitHub Actions**
```bash
mkdir -p .github/workflows

# Create test.yml
cat > .github/workflows/test.yml << 'EOF'
name: Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: go mod download
      - run: go test -v -race -coverprofile=coverage.txt ./...
      - uses: codecov/codecov-action@v4
        with:
          file: ./coverage.txt

  integration-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: go test -v -tags=integration ./...
EOF

# Commit and push
git add .github/workflows/test.yml
git commit -m "Add CI/CD for testing"
git push
```

**Afternoon: Codecov Setup**
```bash
# 1. Go to https://codecov.io
# 2. Sign in with GitHub
# 3. Enable lando-labs/cami repository
# 4. Add badge to README
```

**Evening: Coverage Enforcement**
```bash
# Create check-coverage.sh script
mkdir -p scripts
cat > scripts/check-coverage.sh << 'EOF'
#!/bin/bash
set -e
go test -coverprofile=coverage.out ./...
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "Coverage: ${COVERAGE}%"
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "❌ Below 80% target"
    exit 1
fi
echo "✅ Meets 80% target"
EOF

chmod +x scripts/check-coverage.sh
```

**Goal:** CI/CD running on every push

---

### Thursday: Docs Tests (High Priority 1/3)

**Morning: Docs Test Data**
```bash
cd internal/docs/testdata

# Create test CLAUDE.md files
cat > claude-no-section.md << 'EOF'
# My Project

Existing project documentation.
EOF

cat > claude-with-section.md << 'EOF'
# My Project

<!-- CAMI-MANAGED: DEPLOYED-AGENTS -->
## Deployed Agents

Old content here.
<!-- /CAMI-MANAGED: DEPLOYED-AGENTS -->
EOF
```

**Afternoon: Write Docs Tests**
```bash
cd internal/docs

# Create docs_test.go
# Test UpdateCLAUDEmd, scanDeployedAgents, mergeContent
# Target: 15 tests, 90% coverage

go test -v -cover
```

**Goal:** Docs package at 90% coverage

---

### Friday: Config Tests (High Priority 2/3)

**All Day: Config Package**
```bash
cd internal/config

# Create config_test.go
# Test Load, Save, AddDeployLocation, RemoveDeployLocation
# Target: 12 tests, 85% coverage

go test -v -cover
```

**Review Day's Progress:**
```bash
# Check overall coverage
make test-cover

# Should see:
# internal/agent:  95%
# internal/deploy: 90%
# internal/docs:   90%
# internal/config: 85%
```

**Goal:** 4 critical packages fully tested

---

## Week 2: Complete Core Testing

### Monday: Discovery Tests (High Priority 3/3)

**All Day: Discovery Package**
```bash
cd internal/discovery

# Create scan_test.go
# Test agent scanning and version comparison
# Target: 10 tests, 85% coverage

go test -v -cover
```

---

### Tuesday: CLI Tests (Medium Priority)

**Morning: CLI Test Structure**
```bash
cd internal/cli

# Create test files for each command
touch list_test.go
touch deploy_test.go
touch scan_test.go
touch locations_test.go
```

**Afternoon: Write CLI Tests**
```bash
# Focus on command logic, not Cobra integration
# Target: 15 tests, 70% coverage

go test -v ./internal/cli/
```

---

### Wednesday: E2E CLI Tests

**All Day: End-to-End Tests**
```bash
cd test/e2e

# Create cli_test.go
# Test: list, deploy, scan commands
# Build CLI, run commands, verify output

go test -v -tags=e2e ./test/e2e/
```

---

### Thursday: TUI Tests (Low Priority)

**Morning: Basic TUI Tests**
```bash
cd internal/tui

# Create tui_test.go
# Test model initialization and key handlers
# Target: 10 tests, 50% coverage (low priority)

go test -v -cover
```

**Afternoon: Review Overall Coverage**
```bash
make test-cover

# Target progress:
# Overall: 70-75% (on track for 80%)
```

---

### Friday: Week 2 Review and Adjustments

**Morning: Gap Analysis**
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Open coverage.html in browser
# Identify untested functions
```

**Afternoon: Fill Gaps**
- Add tests for red areas in coverage report
- Focus on critical untested branches
- Aim for 80% overall

**Evening: Documentation**
```bash
# Create TESTING.md
cat > TESTING.md << 'EOF'
# Testing Guide

## Running Tests

```bash
make test           # Unit tests
make test-cover     # With coverage
make test-all       # All tests
```

## Test Structure

- `internal/*/testdata/` - Test fixtures
- `internal/*_test.go` - Unit tests
- `internal/*_integration_test.go` - Integration tests
- `test/e2e/` - End-to-end tests

## Coverage Targets

- Overall: 80%
- Critical packages: 85-95%

## Writing Tests

See [reference/test-implementation-examples.md](reference/test-implementation-examples.md)
EOF
```

---

## Week 3-4: Remaining Packages

### Week 3 Focus Areas

**Monday-Tuesday: MCP Server Tests**
- Test tool handlers in `cmd/cami-mcp/main.go`
- Mock MCP protocol interactions
- Verify JSON responses

**Wednesday-Thursday: Remaining Packages**
- Any untested packages in `internal/`
- Edge cases and error paths
- Concurrent access tests (if applicable)

**Friday: Integration Testing**
- Full workflow tests
- Multi-agent deployments
- Source management (if Phase 2 started)

### Week 4: Polish and Documentation

**Monday-Tuesday: Refactoring for Testability**
- Inject dependencies (filesystem, etc.)
- Extract interfaces for mocking
- Improve error handling

**Wednesday: Performance Testing**
```bash
go test -bench=. ./...
```

**Thursday-Friday: Documentation**
- Update README with testing info
- Contributing guide includes testing requirements
- Architecture diagrams

---

## Week 5: Final Push to 80%

### Daily Routine

**Each Morning:**
```bash
# Check current coverage
make test-cover

# Identify lowest coverage packages
go tool cover -func=coverage.out | sort -k3 -n
```

**Each Afternoon:**
- Write tests for lowest coverage areas
- Focus on critical paths first
- Run tests: `go test -v ./...`

**Each Evening:**
- Review coverage progress
- Document any testing challenges
- Plan next day's focus

### Targets by End of Week 5

- [ ] `internal/agent`: 95%+
- [ ] `internal/deploy`: 90%+
- [ ] `internal/docs`: 90%+
- [ ] `internal/config`: 85%+
- [ ] `internal/discovery`: 85%+
- [ ] `internal/cli`: 70%+
- [ ] `internal/tui`: 50%+
- [ ] `cmd/*`: 40%+
- [ ] **Overall: 80%+**

---

## Week 6: QA and Cross-Platform Testing

### Monday: Cross-Platform CI

**Add macOS and Windows to CI:**
```yaml
# .github/workflows/test.yml
jobs:
  unit-tests:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    # ... rest of job
```

### Tuesday: Manual Testing on Each Platform

**macOS:**
```bash
make test-all
./cami list
./cami deploy --agents architect --location ~/test
```

**Linux (Ubuntu):**
```bash
make test-all
# Test with different terminals (bash, zsh, fish)
```

**Windows (via WSL or VM):**
```bash
make test-all
# Test path handling (Windows vs Unix paths)
```

### Wednesday: Performance Testing

```bash
# Test with large agent library (100+ agents)
# Test deployment to multiple locations
# Test scan on large project

# Benchmark
go test -bench=. -benchmem ./...
```

### Thursday: Security Review

```bash
# Check for hardcoded secrets
grep -r "password\|secret\|token" --include="*.go"

# Check file permissions in tests
# Verify no arbitrary file access

# Run gosec (Go security checker)
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...
```

### Friday: Final Review

**Morning: Coverage Report**
```bash
make test-cover
# Verify 80%+ coverage
# Generate final coverage badge
```

**Afternoon: Documentation Review**
- README.md accurate
- TESTING.md complete
- All examples work
- Contributing guide clear

**Evening: Tag Pre-Release**
```bash
git tag v0.2.1-tested
git push origin v0.2.1-tested
```

---

## Daily Commands Reference

### Run Tests
```bash
make test              # Fast, unit tests only
make test-cover        # With coverage report
make test-integration  # Integration tests
make test-e2e          # End-to-end tests
make test-all          # Everything
```

### Check Coverage
```bash
make test-cover
open coverage.html     # View in browser

# Check specific package
cd internal/agent
go test -cover
```

### Debug Test Failures
```bash
# Verbose output
go test -v ./internal/agent/

# Run specific test
go test -v -run TestLoadAgent_ValidFrontmatter ./internal/agent/

# With race detector
go test -race ./...
```

### Watch Tests (Optional)
```bash
# Install nodemon
npm install -g nodemon

# Watch and re-run tests
nodemon --exec 'go test -v ./internal/agent/' --watch 'internal/agent/*.go'
```

---

## Progress Tracking

### Week 1 Checklist
- [ ] Install testify
- [ ] Create test structure
- [ ] Agent tests (95% coverage)
- [ ] Deploy tests (90% coverage)
- [ ] CI/CD setup
- [ ] Codecov integration

### Week 2 Checklist
- [ ] Docs tests (90% coverage)
- [ ] Config tests (85% coverage)
- [ ] Discovery tests (85% coverage)
- [ ] CLI tests (70% coverage)
- [ ] E2E tests working

### Week 3-4 Checklist
- [ ] MCP server tests
- [ ] Remaining packages
- [ ] Refactoring for testability
- [ ] Performance benchmarks
- [ ] Documentation updated

### Week 5 Checklist
- [ ] 80%+ overall coverage achieved
- [ ] All critical paths tested
- [ ] Coverage enforcement in CI
- [ ] Badge in README

### Week 6 Checklist
- [ ] Cross-platform tests pass
- [ ] Manual QA on macOS/Linux/Windows
- [ ] Performance acceptable
- [ ] Security review complete
- [ ] Pre-release tagged (v0.2.1-tested)

---

## Success Criteria

By end of Week 6, you should have:

✅ **80%+ test coverage** verified by CI
✅ **All critical packages** at target coverage
✅ **CI/CD passing** on every commit
✅ **Cross-platform** tests passing
✅ **Documentation** complete (README, TESTING.md)
✅ **Pre-release tagged** (v0.2.1-tested)
✅ **Ready for Phase 1** (repository split)

---

## Common Issues and Solutions

### Issue: Coverage below target
**Solution:**
```bash
# Find uncovered code
go tool cover -html=coverage.out

# Add tests for red/yellow areas
# Focus on if/else branches
# Test error conditions
```

### Issue: Tests flaky (pass/fail randomly)
**Solution:**
```bash
# Run with race detector
go test -race ./...

# Check for:
# - Shared state between tests
# - Time-dependent logic
# - Filesystem cleanup
```

### Issue: Tests slow
**Solution:**
```bash
# Profile tests
go test -cpuprofile=cpu.out ./...
go tool pprof cpu.out

# Use t.Parallel() for independent tests
# Mock expensive operations
# Use t.TempDir() for filesystem tests
```

### Issue: CI fails but local passes
**Solution:**
```bash
# Check OS-specific issues
# Verify file paths (Unix vs Windows)
# Check Go version matches CI
# Run in Docker to reproduce CI environment
```

---

## Next Steps After Week 6

Once you achieve 80%+ coverage and all tests pass:

1. **Code Review:** Review all test code with team
2. **Merge to Main:** Merge testing branch
3. **Announce:** Share coverage achievement with team
4. **Begin Phase 1:** Start repository split (Week 7)
5. **Maintain:** Keep coverage at 80%+ for all new code

---

**You've got this! Week 1 is the hardest. Once you have the pattern down, the rest flows naturally.**

**Remember:** Tests are documentation. Write tests that explain how the code should work.
