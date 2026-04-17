# Changelog

## v1.0.0 — Complete Product: 2D Solver, YAML Input, ASCII Diagrams
Released: 2026-04-16

V1 ships the full described product. A user can now describe a pegboard,
deck, or sheet-goods job in a YAML file, run the tool, and get an ASCII
diagram showing how to cut each board or sheet.

Features:
- **2D cutting solver** — recursive guillotine packing with branch-and-bound pruning
- **Rotation support** — pieces may be rotated 90° for better fit; disable with `--no-rotate`
- **YAML file input** — describe jobs in a file with `-f myproject.yaml`
- **Flag + file merge** — file provides the base; flags add stock/need and override settings
- **ASCII diagram output** — default output; visual grid for 2D, labeled bar for 1D
- **Mode auto-detection** — 1D or 2D selected automatically based on input dimensions
- **Output formats** — ASCII (default), plain text (`--output text`), JSON (`--output json`)
- **Fuzz tests** — property-based invariant testing for both solvers (57k+ execs clean)
- **Golden file tests** — byte-for-exact output regression tests for all 3 formats × 2 modes
- **CI fuzz runs** — 30-second fuzz runs added to GitHub Actions on every push

---

## v0.1.0 — MVP: 1D Cutting Optimizer
Released: 2026-04-16

First release. Given a list of required lengths and available stock
(including on-hand offcuts), produces an optimal cutting plan with
minimum waste.

Features:
- 1D cutting stock solver (knapsack DP + greedy assignment)
- Mixed stock: on-hand pieces used before purchasable stock
- Kerf width support (`--kerf` flag)
- Plain text and JSON output (`--output text|json`)
- CLI flags: `--stock`, `--need`, `--kerf`, `--output`
