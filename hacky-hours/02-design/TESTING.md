# Testing

## Definition of Done

A cut plan is **correct** if and only if:

1. Every required piece appears in exactly one assignment
2. No stock piece is over-allocated (assigned pieces + kerf ≤ stock size)
3. Waste percentage equals `(used stock − required pieces) / used stock`
4. On-hand stock is fully attempted before any purchasable stock is used
5. If `--no-rotate` is set, no assignment has `Rotated: true`
6. All pieces that could not be fit are reported — none are silently dropped

A cut plan is **optimal** if no alternative valid plan produces lower waste for the same inputs. Optimality is verified against known reference cases (see below).

---

## Test Strategy

### 1. Unit tests — solver invariants (property-based)

Use Go's built-in fuzz testing (`testing/fuzz`, available since Go 1.18) to generate random valid inputs and assert the invariants above hold for every output.

Invariants to assert after every `Solve()` call:

```go
// Every required piece is assigned
assignedLabels := collectLabels(plan.Assignments)
assert.ElementsMatch(t, requiredLabels(need), assignedLabels)

// No stock piece is over-allocated
for _, stock := range usedStock {
    assert.LessOrEqual(t, totalAssignedSize(stock, plan) + totalKerf(stock, plan), stock.Size())
}

// Waste is non-negative
assert.GreaterOrEqual(t, plan.WastePct, 0.0)
assert.LessOrEqual(t, plan.WastePct, 1.0)

// On-hand stock used before purchasable
if anyOnHandRemaining(stock, plan) {
    assert.Empty(t, plan.Purchased)
}
```

### 2. Table-driven tests — known reference cases

For small inputs where the optimal answer is known by hand, assert that the solver finds it.

**1D reference cases:**

| Stock | Need | Expected waste |
|-------|------|---------------|
| `[96"]` | `[48", 48"]` | 0% |
| `[96"]` | `[48", 48", 1"]` | needs 2 boards; waste ≈ 49% |
| `[96"]` | `[36", 36", 36"]` | waste = 12" / 96" = 12.5% |
| `[96", 48" offcut]` | `[48", 48"]` | 0% (offcut used first) |
| `[96"]` | `[97"]` | error: piece larger than stock |

**2D reference cases:**

| Stock | Need | Expected waste |
|-------|------|---------------|
| `[48×96]` | `[48×48, 48×48]` | 0% |
| `[48×96]` | `[24×96, 24×96]` | 0% |
| `[48×96]` | `[30×30, 30×30, 30×30]` | needs 1 sheet; waste ≈ 44% |
| `[48×96]` | `[30×48, 48×30]` | 0% with rotation, ~37% without |

### 3. Table-driven tests — input parsing

Verify that flag strings and YAML files parse to the correct model types.

```
"96"          → StockPiece{Length: 96}
"96:3"        → StockPiece{Length: 96, Count: 3}
"96:3:onhand" → StockPiece{Length: 96, Count: 3, OnHand: true}
"48x96"       → StockPiece{Width: 48, Height: 96}
"48x96:onhand"→ StockPiece{Width: 48, Height: 96, OnHand: true}
```

Invalid inputs that must return errors:

```
""            → error
"0"           → error (zero-size stock)
"-5"          → error (negative size)
"48x"         → error (malformed 2D)
"48x96x12"    → error (three dimensions not supported)
```

### 4. Golden file tests — output formatting

For one known `CutPlan`, assert that each formatter produces byte-for-byte identical output to a stored golden file. Run with `-update` flag to regenerate golden files when output format intentionally changes.

```
testdata/
  golden/
    simple_1d.ascii.txt
    simple_1d.text.txt
    simple_1d.json
    simple_2d.ascii.txt
    simple_2d.text.txt
    simple_2d.json
```

### 5. Edge case tests

Directly test the documented edge cases from BUSINESS_LOGIC.md:

| Input | Expected behavior |
|-------|------------------|
| Required piece > all stock | Error listing which piece(s), no partial plan |
| Empty cut list | Error: nothing to do |
| Empty stock | Error: no stock provided |
| Required piece == stock size exactly | Valid assignment, 0% waste, no saw cut annotation |
| All stock on-hand, all fits | `Purchased` is empty |
| On-hand insufficient | `Purchased` lists minimum additional stock |
| Duplicate required sizes | Both treated as separate pieces, both assigned |
| Stock smaller than smallest required | Reported as unusable offcut, not used |

---

## Tools

| Tool | Purpose | License |
|------|---------|---------|
| Go standard `testing` | Unit, table-driven, fuzz tests | BSD ✓ |
| `testing/fuzz` (Go 1.18+) | Property-based / fuzz testing | BSD ✓ |
| `github.com/stretchr/testify` | Assertions and require helpers | MIT ✓ |

No external test runner needed. `go test ./...` runs everything.

---

## Test File Layout

```
internal/
  input/
    flags_test.go
    file_test.go
    validate_test.go
  solver/
    solver1d/
      dp_test.go          ← reference cases + fuzz
    solver2d/
      guillotine_test.go  ← reference cases + fuzz
  output/
    ascii_test.go         ← golden files
    text_test.go
    json_test.go
testdata/
  golden/                 ← golden output files
  fuzz/                   ← fuzz corpus seeds
    solver1d/
    solver2d/
```

---

## CI

Run on every PR:

```bash
go test ./...
go test -fuzz=FuzzSolve1D -fuzztime=30s ./internal/solver/solver1d/
go test -fuzz=FuzzSolve2D -fuzztime=30s ./internal/solver/solver2d/
```

Fuzz time is kept short in CI; longer runs can be done locally when changing solver logic.
