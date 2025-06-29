# Test Update Plan for Muse Backend

## Current Status

- **Total Coverage:** 10.9% (needs improvement)
- **Schema Changes:** Migration 003 requires test updates
- **Pipeline:** Well-structured, professional setup

## Priority 1: Update Existing Tests (This Week)

### 1. Repository Tests Update

- [ ] Update `user_repository_test.go` for new schema
- [ ] Remove outdated tests for deleted tables (`artists`, `albums`, `tracks`)
- [ ] Add tests for new `user_preferences` table
- [ ] Update `playlist_tracks` tests for Spotify ID approach
- [ ] Add tests for optimized `reviews` table with Spotify IDs

### 2. Add Missing Core Tests

- [ ] Database connection tests (`internal/database`)
- [ ] Models validation tests (`internal/models`)
- [ ] GraphQL resolver tests (currently 0% coverage)
- [ ] Spotify integration tests with better coverage

### 3. Integration Test Updates

- [ ] Update GraphQL integration tests for new schema
- [ ] Add tests for Spotify ID-based operations
- [ ] Test new playlist functionality with Spotify tracks

## Priority 2: Expand Test Coverage (Next 2 Weeks)

### Target Coverage Goals

- **Overall:** 70%+ (from current 10.9%)
- **Repository Layer:** 80%+ (from current 9.4%)
- **GraphQL Resolvers:** 60%+ (from current 0%)
- **Business Logic:** 85%+

### New Test Areas

- [ ] End-to-end GraphQL operations
- [ ] Authentication flows with JWT
- [ ] Spotify OAuth callback handling
- [ ] Error handling scenarios
- [ ] Performance benchmarks for new schema

## Priority 3: Pipeline Enhancements (Month 3)

### CI/CD Improvements

- [ ] Add test result reporting
- [ ] Implement test performance monitoring
- [ ] Add automated test database seeding
- [ ] Set up coverage thresholds (fail if below 60%)

### Quality Gates

- [ ] Require 70% coverage for PR approval
- [ ] Add integration test requirements
- [ ] Performance regression detection

## Commands to Run Tests

```bash
# Update schema-dependent tests
make test-unit                    # Fast unit tests
make test                        # Full integration tests
make test-coverage               # Generate coverage report

# CI simulation
make ci-checks                   # Run what CI runs locally

# Performance testing
make bench                       # Run benchmarks
```

## Schema Alignment Checklist

### ✅ Updated for Migration 003

- [x] Removed `artists`, `albums`, `tracks` table references
- [x] Updated to use Spotify IDs instead of local data
- [x] Added `user_preferences` table support

### ❌ Still Needs Updates

- [ ] Repository tests for new playlist_tracks structure
- [ ] Tests for Spotify ID-based reviews
- [ ] Integration tests for user preferences
- [ ] GraphQL resolver tests for all new functionality

## Success Metrics

By completion, we should have:

- **70%+ overall test coverage**
- **All repository methods tested**
- **GraphQL resolvers covered**
- **CI pipeline passing consistently**
- **No schema-related test failures**
- **Performance benchmarks established**

## Timeline

- **Week 1:** Update existing tests for new schema
- **Week 2-3:** Add missing core tests, reach 50% coverage
- **Week 4:** Reach 70% coverage target
- **Month 2:** Add performance tests and benchmarks
- **Month 3:** Pipeline enhancements and quality gates
