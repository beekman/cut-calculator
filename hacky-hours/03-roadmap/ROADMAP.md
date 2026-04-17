# Roadmap

## MVP — 1D Cutting Optimizer

The smallest version that proves the core value: given a list of lengths to cut
and the boards you have (plus what you can buy), produce the optimal cutting plan
with minimum waste.

**Done when:** A user can run `cut-calculator --stock 96:3:onhand --need 36:4 --need 48:2 --kerf 0.125`
and get a correct, human-readable cut list back.

### Tasks
- [ ] Project scaffold (Go module, package layout, CI)
- [ ] Core model types (`StockPiece`, `RequiredPiece`, `CutPlan`, `Assignment`)
- [ ] CLI flag parsing (`--stock`, `--need`, `--kerf`, `--output`)
- [ ] Input validation and error reporting
- [ ] 1D solver: bounded knapsack DP
- [ ] 1D solver: branch-and-bound assignment
- [ ] Mixed stock: on-hand priority + purchasable recommendation
- [ ] Plain text output formatter
- [ ] Reference case tests + edge case tests

---

## V1 — Complete Product

Adds 2D sheet cutting, visual output, file-based input, and JSON export.

**Done when:** A user can describe a pegboard job in a YAML file, run the tool,
and get an ASCII diagram showing how to cut each sheet.

### Tasks
- [ ] YAML file input parser
- [ ] Flag + file merge logic
- [ ] 2D solver: recursive guillotine packing
- [ ] 2D solver: rotation support + `--no-rotate` flag
- [ ] ASCII diagram output (1D and 2D)
- [ ] JSON output formatter
- [ ] ASCII becomes default output; plain text via `--output text`
- [ ] Fuzz tests (1D and 2D solvers)
- [ ] Golden file tests (all output formats)

---

## V2+ — Future

- Wallpaper pattern-repeat matching
