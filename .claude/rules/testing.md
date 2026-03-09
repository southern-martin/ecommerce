---
paths:
  - "**/*_test.go"
  - "tests/**/*.go"
---
# Testing Conventions

## Commands
```bash
make test                                    # all services + pkg (race detection)
make test SVC=product                        # single service
GOWORK=off go test ./tests/...               # integration/E2E (separate module)
make integration-test                        # integration tests
make e2e-test                                # E2E tests
make load-test-smoke                         # k6 smoke test
```

## Guidelines
- Tests module (`tests/`) is NOT in go.work — always use `GOWORK=off`
- Table-driven tests, `t.Run()` subtests
- Mock at repository interface level
- CI pipeline: lint → unit → integration → security scan → Docker build
