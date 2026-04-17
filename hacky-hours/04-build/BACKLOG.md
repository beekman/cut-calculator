# Backlog

## MVP — 1D Cutting Optimizer

- [x] Project scaffold (Go module, package layout, CI)
- [x] Core model types (`StockPiece`, `RequiredPiece`, `CutPlan`, `Assignment`)
- [x] CLI flag parsing (`--stock`, `--need`, `--kerf`, `--output`)
- [x] Input validation and error reporting
- [x] 1D solver: bounded knapsack DP
- [x] 1D solver: branch-and-bound assignment
- [x] Mixed stock: on-hand priority + purchasable recommendation
- [x] Plain text output formatter
- [x] Reference case tests + edge case tests
