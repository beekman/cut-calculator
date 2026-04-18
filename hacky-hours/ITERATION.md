# Iteration — v1.1.0 Post-Release

## Captured

- `cmd/cut-calculator/main.go` was never committed to git. The v1.1.0 tag does not
  include `cmd/`, so `go install github.com/beekman/cut-calculator/cmd/cut-calculator@latest`
  fails with "does not contain package". The binary entrypoint is missing from the
  published module.

## Triage

| Item | Category | Notes |
|------|----------|-------|
| `cmd/` untracked; `go install` broken | **Hotfix** | Commit `cmd/`, patch release |

## Status

- [ ] Hotfix shipped
