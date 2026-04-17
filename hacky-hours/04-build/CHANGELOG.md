# Changelog

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
